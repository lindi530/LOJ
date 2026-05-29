package problem_service

import (
	"GO1/pkg/constants"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	judgeDefaultTimeLimitSeconds = 3
	judgeCompileTimeLimitSeconds = 15
	judgeDefaultMemoryLimitMB    = 128
	judgeDockerStartupSeconds    = 10
	judgeDockerCleanupSeconds    = 5
	judgeDockerStderrLimitBytes  = 4096
	// ulimit -f uses 512-byte blocks on Linux. 32768 blocks caps each output file at 16 MiB.
	judgeOutputFileBlocks = 32768
	judgeOutputLimitBytes = judgeOutputFileBlocks * 512
)

var errJudgeOutputLimitExceeded = errors.New("judge output limit exceeded")

func limitedJudgeCommand(command string, timeLimitSeconds int) string {
	return fmt.Sprintf(
		"ulimit -f %d; timeout %ds sh -c %s",
		judgeOutputFileBlocks,
		normalizeJudgeTimeLimitSeconds(timeLimitSeconds),
		shellQuote(command),
	)
}

func timeLimitedJudgeCommand(command string, timeLimitSeconds int) string {
	return fmt.Sprintf(
		"timeout %ds sh -c %s",
		normalizeJudgeTimeLimitSeconds(timeLimitSeconds),
		shellQuote(command),
	)
}

func normalizeJudgeTimeLimitSeconds(timeLimit int) int {
	if timeLimit <= 0 {
		return judgeDefaultTimeLimitSeconds
	}
	if timeLimit > 60 {
		return (timeLimit + 999) / 1000
	}
	return timeLimit
}

func normalizeJudgeMemoryLimitMB(memoryLimit int) int {
	if memoryLimit <= 0 {
		return judgeDefaultMemoryLimitMB
	}
	return memoryLimit
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func buildDockerShellCmd(image, tempDir, command, containerName string, memoryLimit int) []string {
	dockerCmd := []string{
		"run", "--rm",
		"--name", containerName,
		"--network", "none",
	}
	if memoryLimit > 0 {
		limit := normalizeJudgeMemoryLimitMB(memoryLimit)
		dockerCmd = append(dockerCmd,
			fmt.Sprintf("--memory=%dm", limit),
			fmt.Sprintf("--memory-swap=%dm", limit),
		)
	}
	dockerCmd = append(dockerCmd,
		"-v", fmt.Sprintf("%s:/app", tempDir),
		image,
		"sh", "-c",
		command,
	)
	return dockerCmd
}

func judgeContainerName(tempDir, stage string) string {
	raw := filepath.Base(tempDir) + "-" + stage
	var name strings.Builder
	name.WriteString("loj-")
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			name.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			name.WriteRune(r)
		case r >= '0' && r <= '9':
			name.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			name.WriteRune(r)
		default:
			name.WriteByte('-')
		}
	}
	return name.String()
}

func judgeCompileTimeout() time.Duration {
	return time.Duration(judgeCompileTimeLimitSeconds+judgeDockerStartupSeconds) * time.Second
}

func judgeRunTimeout(timeLimit, testCaseCount int) time.Duration {
	if testCaseCount <= 0 {
		testCaseCount = 1
	}
	seconds := judgeDockerStartupSeconds + testCaseCount*(normalizeJudgeTimeLimitSeconds(timeLimit)+1)
	return time.Duration(seconds) * time.Second
}

func runDockerCommand(dockerCmd []string, timeout time.Duration, containerName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", dockerCmd...)
	stderr := &limitedBuffer{limit: judgeDockerStderrLimitBytes}
	cmd.Stderr = stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		removeDockerContainer(containerName)
		return stderr.String(), ctx.Err()
	}
	return stderr.String(), err
}

func removeDockerContainer(containerName string) {
	if containerName == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(judgeDockerCleanupSeconds)*time.Second)
	defer cancel()
	_ = exec.CommandContext(ctx, "docker", "rm", "-f", containerName).Run()
}

func judgeStatusFromDockerError(err error, stderr string) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return constants.JudgeStatusSystemError
	}
	if strings.Contains(stderr, "File size limit exceeded") {
		return constants.JudgeStatusOutputLimitExceeded
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		switch exitErr.ExitCode() {
		case 124:
			return constants.JudgeStatusTimeLimitExceeded
		case 137:
			return constants.JudgeStatusMemoryLimitExceeded
		case 153:
			return constants.JudgeStatusOutputLimitExceeded
		case 125, 126, 127:
			return constants.JudgeStatusSystemError
		default:
			return constants.JudgeStatusRuntimeError
		}
	}
	return constants.JudgeStatusSystemError
}

func readJudgeOutput(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.Size() > judgeOutputLimitBytes {
		return nil, errJudgeOutputLimitExceeded
	}
	return os.ReadFile(path)
}

type limitedBuffer struct {
	bytes.Buffer
	limit int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.limit <= 0 || b.Len() >= b.limit {
		return len(p), nil
	}
	remaining := b.limit - b.Len()
	if len(p) > remaining {
		_, _ = b.Buffer.Write(p[:remaining])
		return len(p), nil
	}
	_, _ = b.Buffer.Write(p)
	return len(p), nil
}
