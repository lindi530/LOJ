package problem_service

import (
	"GO1/global"
	"GO1/models/problem_model"
	"GO1/models/ws_model"
	"GO1/pkg/constants"
	"GO1/service/ws_service"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func RunExample(code, language, input string, memoryLimit, timeLimit int, message *ws_model.EditStatus) problem_model.RunResult {
	ws_service.WsHub.SendEditData(message, constants.JudgeStatusPending)
	// 1. 创建临时目录
	dir := ""
	if tempDirEnv := os.Getenv("TEMP_DIR"); tempDirEnv != "" {
		dir = tempDirEnv
	}
	tempDir, err := os.MkdirTemp(dir, "ojcode-*")
	if err != nil {
		return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusSystemError}
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			global.Logger.Warn("remove judge temp dir failed: ", tempDir, " err: ", err)
		}
	}()

	// 2. 确定文件名、命令和镜像
	var codeFileName, compileCmd, runCmd, image, runtimeSetupCmd string

	switch language {
	case "cpp":
		codeFileName = "code.cpp"
		compileCmd = "g++ /app/code.cpp -o /app/a.out"
		runtimeSetupCmd = "cp /app/a.out /tmp/a.out && chmod +x /tmp/a.out"
		runCmd = "/tmp/a.out < /app/input.txt > /app/output.txt 2> /app/error.txt"
		image = "gcc:15"
	case "python":
		codeFileName = "code.py"
		runCmd = "python3 /app/code.py < /app/input.txt > /app/output.txt 2> /app/error.txt"
		image = "python:3.12"
	case "java":
		codeFileName = "Main.java"
		compileCmd = "javac /app/Main.java"
		runCmd = "java -cp /app Main < /app/input.txt > /app/output.txt 2> /app/error.txt"
		image = "openjdk:latest"
	case "go":
		codeFileName = "code.go"
		compileCmd = "go build -o /app/a.out /app/code.go"
		runtimeSetupCmd = "cp /app/a.out /tmp/a.out && chmod +x /tmp/a.out"
		runCmd = "/tmp/a.out < /app/input.txt > /app/output.txt 2> /app/error.txt"
		image = "golang:latest"
	default:
		return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusSystemError}
	}

	// 3. 写入文件
	err = os.WriteFile(filepath.Join(tempDir, codeFileName), []byte(code), 0644)
	if err != nil {
		return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusSystemError}
	}

	err = os.WriteFile(filepath.Join(tempDir, "input.txt"), []byte(input), 0644)
	if err != nil {
		return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusSystemError}
	}

	// 4. 构建 docker 命令
	if compileCmd != "" {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusCompiling)
		compileContainerName := judgeContainerName(tempDir, "compile")
		compileDockerCmd := buildDockerShellCmd(image, tempDir, timeLimitedJudgeCommand(compileCmd, judgeCompileTimeLimitSeconds), compileContainerName, 0)
		global.Logger.Info("docker compile command: ", compileDockerCmd)
		compileStderr, err := runDockerCommand(compileDockerCmd, judgeCompileTimeout(), compileContainerName)
		if err != nil {
			if judgeStatusFromDockerError(err, compileStderr) == constants.JudgeStatusSystemError {
				return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusSystemError}
			}
			return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusCompileError}
		}
	}

	if runtimeSetupCmd != "" {
		runCmd = runtimeSetupCmd + " && " + runCmd
	}
	runContainerName := judgeContainerName(tempDir, "run")
	dockerCmd := buildDockerCmd(image, tempDir, runCmd, runContainerName, memoryLimit, timeLimit)
	global.Logger.Info("docker command: ", dockerCmd)
	ws_service.WsHub.SendEditData(message, constants.JudgeStatusRunning)
	stderr, err := runDockerCommand(dockerCmd, judgeRunTimeout(timeLimit, 1), runContainerName)
	if err != nil {
		return problem_model.RunResult{Passed: false, Error: judgeStatusFromDockerError(err, stderr)}
	}

	// 7. 读取输出文件
	outBytes, err := readJudgeOutput(filepath.Join(tempDir, "output.txt"))
	if err != nil {
		if errors.Is(err, errJudgeOutputLimitExceeded) {
			return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusOutputLimitExceeded}
		}
		return problem_model.RunResult{Passed: false, Error: constants.JudgeStatusSystemError}
	}

	ws_service.WsHub.SendEditData(message, constants.JudgeStatusAccepted)

	return problem_model.RunResult{
		Passed: true,
		Output: strings.TrimSpace(string(outBytes)),
	}
}

func buildDockerCmd(image, tempDir, runCmd, containerName string, memoryLimit, timeLimit int) []string {
	return buildDockerShellCmd(image, tempDir, limitedJudgeCommand(runCmd, timeLimit), containerName, memoryLimit)
}

//  config 中写了配置
// isolate linux
//func RunExample(code, language, input string, memoryLimit, timeLimit int, message *ws_model.MessageWs) problem_model.RunResult {
//	ws_service.WsHub.CodeStateWs(message, "Pending")
//
//	// 1. 创建临时目录和文件
//	tempDir, err := os.MkdirTemp("", "ojcode-*")
//	if err != nil {
//		return problem_model.RunResult{Passed: false, Error: "Cannot create temp directory"}
//	}
//	defer os.RemoveAll(tempDir)
//
//	// 2. isolate sandbox 初始化
//	boxID := "0"
//	initCmd := exec.Command("isolate", "--box-id="+boxID, "--init")
//	if err := initCmd.Run(); err != nil {
//		return problem_model.RunResult{Passed: false, Error: "Failed to init isolate"}
//	}
//
//	boxPath := "/var/local/lib/isolate/" + boxID + "/box"
//
//	// 3. 写入代码和输入
//	languageConfig := global.Config.Language.Compile(language)
//	if languageConfig == nil {
//		return problem_model.RunResult{Passed: false, Error: "不支持该语言"}
//	}
//
//	codeFileName := languageConfig.SourceFile
//	compileCmd := languageConfig.CompileCmd // 如果不为空就编译
//
//	if err := os.WriteFile(filepath.Join(boxPath, codeFileName), []byte(code), 0644); err != nil {
//		return problem_model.RunResult{Passed: false, Error: "Cannot write code file"}
//	}
//	if err := os.WriteFile(filepath.Join(boxPath, "input.txt"), []byte(input), 0644); err != nil {
//		return problem_model.RunResult{Passed: false, Error: "Cannot write input file"}
//	}
//
//	// 4. 编译（如果需要）
//	if len(compileCmd) != 0 {
//		cmd := append([]string{
//			"--box-id=" + boxID,
//			"--run",
//			"--",
//		}, compileCmd...)
//		compile := exec.Command("isolate", cmd...)
//		var compileErr bytes.Buffer
//		compile.Stderr = &compileErr
//		if err := compile.Run(); err != nil {
//			return problem_model.RunResult{Passed: false, Error: "Compile Error: " + compileErr.String()}
//		}
//	}
//
//	// 5. 运行主程序
//	ws_service.WsHub.CodeStateWs(message, "Running")
//
//	outputFile := filepath.Join(boxPath, "output.txt")
//	runCmd := strings.Join(languageConfig.RunCmd, " ")
//	runArgs := []string{
//		"--box-id=" + boxID,
//		"--stdin=input.txt",
//		"--stdout=output.txt",
//		"--stderr=error.txt",
//		"--time", fmt.Sprintf("%d", timeLimit),
//		"--mem", fmt.Sprintf("%d", memoryLimit*1024), // isolate 以 KB 为单位
//		"--run", "--", "sh", "-c", runCmd,
//	}
//	runCmdExec := exec.Command("isolate", runArgs...)
//	if err := runCmdExec.Run(); err != nil {
//		return problem_model.RunResult{Passed: false, Error: "Runtime Error"}
//	}
//
//	// 6. 读取输出
//	outBytes, err := os.ReadFile(outputFile)
//	if err != nil {
//		return problem_model.RunResult{Passed: false, Error: "Cannot read output"}
//	}
//
//	// 7. 清理 sandbox
//	_ = exec.Command("isolate", "--box-id="+boxID, "--cleanup").Run()
//
//	ws_service.WsHub.CodeStateWs(message, "Finished")
//
//	return problem_model.RunResult{
//		Passed: true,
//		Output: strings.TrimSpace(string(outBytes)),
//	}
//}
