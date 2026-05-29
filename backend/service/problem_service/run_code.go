package problem_service

import (
	"GO1/global"
	"GO1/models/problem_model"
	"GO1/models/ws_model"
	"GO1/pkg/constants"
	"GO1/service/ws_service"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunCode(userid int64, code, language string, testcases *[]problem_model.Example, memoryLimit, timeLimit int, message *ws_model.EditStatus) []problem_model.RunResult {
	dir := ""
	if tempDirEnv := os.Getenv("TEMP_DIR"); tempDirEnv != "" {
		dir = tempDirEnv
	}
	tempDir, err := os.MkdirTemp(dir, "ojcode-*")
	if err != nil {
		return []problem_model.RunResult{{Passed: false, Error: constants.JudgeStatusSystemError}}
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			global.Logger.Warn("remove judge temp dir failed: ", tempDir, " err: ", err)
		}
	}()

	var codeFileName, compileCmd, runCmd, image, runtimeSetupCmd string

	switch language {
	case "cpp":
		codeFileName = "code.cpp"
		// Native binaries cannot execute from the shared /app judge volume.
		compileCmd = "g++ /app/code.cpp -o /app/a.out"
		runtimeSetupCmd = "cp /app/a.out /tmp/a.out && chmod +x /tmp/a.out"
		runCmd = "/tmp/a.out"
		image = "gcc:15"
	case "python":
		codeFileName = "code.py"
		runCmd = "python3 /app/code.py"
		image = "python:3.12"
	case "java":
		codeFileName = "Main.java"
		compileCmd = "javac /app/Main.java"
		runCmd = "java -cp /app Main"
		image = "openjdk:latest"
	case "go":
		codeFileName = "code.go"
		compileCmd = "go build -o /app/a.out /app/code.go"
		runtimeSetupCmd = "cp /app/a.out /tmp/a.out && chmod +x /tmp/a.out"
		runCmd = "/tmp/a.out"
		image = "golang:latest"
	default:
		return []problem_model.RunResult{{Passed: false, Error: constants.JudgeStatusSystemError}}
	}

	// 写入代码文件
	codePath := filepath.Join(tempDir, codeFileName)
	err = os.WriteFile(codePath, []byte(code), 0644)
	if err != nil {
		return []problem_model.RunResult{{Passed: false, Error: constants.JudgeStatusSystemError}}
	}

	// 写入所有样例的输入文件
	for i, tc := range *testcases {
		inputPath := filepath.Join(tempDir, fmt.Sprintf("input%d.txt", i))
		err := os.WriteFile(inputPath, []byte(tc.Input), 0644)
		if err != nil {
			return []problem_model.RunResult{{Passed: false, Error: constants.JudgeStatusSystemError}}
		}
	}

	// 写入运行脚本 run.sh
	runScript := "#!/bin/sh\nset -e\n"
	if runtimeSetupCmd != "" {
		runScript += runtimeSetupCmd + "\n"
	}
	for i := range *testcases {
		runCommand := fmt.Sprintf("%s < /app/input%d.txt > /app/output%d.txt 2> /app/error%d.txt", runCmd, i, i, i)
		runScript += limitedJudgeCommand(runCommand, timeLimit) + "\n"
	}
	runScriptPath := filepath.Join(tempDir, "run.sh")
	err = os.WriteFile(runScriptPath, []byte(runScript), 0755)
	if err != nil {
		return []problem_model.RunResult{{Passed: false, Error: constants.JudgeStatusSystemError}}
	}

	// 构建 Docker 命令
	if compileCmd != "" {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusCompiling)
		compileContainerName := judgeContainerName(tempDir, "compile")
		compileDockerCmd := buildDockerShellCmd(image, tempDir, timeLimitedJudgeCommand(compileCmd, judgeCompileTimeLimitSeconds), compileContainerName, 0)
		compileStderr, err := runDockerCommand(compileDockerCmd, judgeCompileTimeout(), compileContainerName)
		if err != nil {
			if judgeStatusFromDockerError(err, compileStderr) == constants.JudgeStatusSystemError {
				return []problem_model.RunResult{{Passed: false, Error: constants.JudgeStatusSystemError}}
			}
			return []problem_model.RunResult{{Passed: false, Error: constants.JudgeStatusCompileError}}
		}
	}
	// 执行 Docker 命令
	ws_service.WsHub.SendEditData(message, constants.JudgeStatusRunning)
	runContainerName := judgeContainerName(tempDir, "run")
	dockerCmd := buildDockerShellCmd(image, tempDir, "sh /app/run.sh", runContainerName, memoryLimit)
	stderr, err := runDockerCommand(dockerCmd, judgeRunTimeout(timeLimit, len(*testcases)), runContainerName)
	if err != nil {
		return []problem_model.RunResult{{Passed: false, Error: judgeStatusFromDockerError(err, stderr)}}
	}

	// 读取每个样例的输出并比较
	results := make([]problem_model.RunResult, len(*testcases))
	for i, tc := range *testcases {
		outPath := filepath.Join(tempDir, fmt.Sprintf("output%d.txt", i))
		outBytes, err := readJudgeOutput(outPath)
		if err != nil {
			errStatus := constants.JudgeStatusSystemError
			if errors.Is(err, errJudgeOutputLimitExceeded) {
				errStatus = constants.JudgeStatusOutputLimitExceeded
			}
			results[i] = problem_model.RunResult{
				Passed: false,
				Error:  errStatus,
			}
			continue
		}
		output := strings.TrimSpace(string(outBytes))
		expected := strings.TrimSpace(tc.Output)
		passed := output == expected
		errStatus := ""
		if !passed {
			errStatus = constants.JudgeStatusWrongAnswer
		}
		results[i] = problem_model.RunResult{
			Passed:   passed,
			Output:   output,
			Expected: expected,
			Error:    errStatus,
		}
	}

	return results
}
