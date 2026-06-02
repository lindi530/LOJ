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
	"strconv"
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

type judgeMetric struct {
	ExecTime    int
	MemoryUsage int
}

func measuredJudgeCommand(command string, timeLimitSeconds, memoryLimit int, metricPath string) string {
	cpuLimit := normalizeJudgeTimeLimitSeconds(timeLimitSeconds)
	timeLimitMS := normalizeJudgeTimeLimitMilliseconds(timeLimitSeconds)
	memoryLimitKB := 0
	if memoryLimit > 0 {
		memoryLimitKB = normalizeJudgeMemoryLimitKB(memoryLimit)
	}
	measuredCommand := fmt.Sprintf("ulimit -t %d; %s", cpuLimit, command)
	return strings.Join([]string{
		fmt.Sprintf("ulimit -f %d", judgeOutputFileBlocks),
		fmt.Sprintf("time_limit_ms=%d", timeLimitMS),
		fmt.Sprintf("memory_limit_kb=%d", memoryLimitKB),
		"read_cpu_usage_usec() {",
		"  if [ -r /sys/fs/cgroup/cpu.stat ]; then",
		"    while read -r name value _; do",
		"      if [ \"$name\" = \"usage_usec\" ]; then echo \"${value:-0}\"; return; fi",
		"    done < /sys/fs/cgroup/cpu.stat",
		"  fi",
		"  if [ -r /sys/fs/cgroup/cpuacct/cpuacct.usage ]; then",
		"    value=$(cat /sys/fs/cgroup/cpuacct/cpuacct.usage 2>/dev/null || echo 0)",
		"    echo $((value / 1000)); return",
		"  fi",
		"  echo 0",
		"}",
		"read_memory_usage_bytes() {",
		"  if [ -r /sys/fs/cgroup/memory.current ]; then cat /sys/fs/cgroup/memory.current; return; fi",
		"  if [ -r /sys/fs/cgroup/memory/memory.usage_in_bytes ]; then cat /sys/fs/cgroup/memory/memory.usage_in_bytes; return; fi",
		"  if [ -r /sys/fs/cgroup/memory.usage_in_bytes ]; then cat /sys/fs/cgroup/memory.usage_in_bytes; return; fi",
		"  echo 0",
		"}",
		fmt.Sprintf("metric_file=%s", shellQuote(metricPath)),
		"cpu_before=$(read_cpu_usage_usec)",
		"case \"$cpu_before\" in ''|*[!0-9]*) cpu_before=0 ;; esac",
		"max_memory=0",
		"timed_out=0",
		"sample_memory_usage() {",
		"  current_memory=$(read_memory_usage_bytes)",
		"  case \"$current_memory\" in ''|*[!0-9]*) current_memory=0 ;; esac",
		"  if [ \"$current_memory\" -gt \"$max_memory\" ]; then max_memory=$current_memory; fi",
		"}",
		"kill_judge_command() {",
		"  if [ -n \"$run_group_pid\" ]; then",
		"    kill -TERM -$run_group_pid 2>/dev/null || kill -TERM \"$run_pid\" 2>/dev/null",
		"    sleep 0.05",
		"    kill -KILL -$run_group_pid 2>/dev/null || kill -KILL \"$run_pid\" 2>/dev/null",
		"  else",
		"    kill -TERM \"$run_pid\" 2>/dev/null",
		"    sleep 0.05",
		"    kill -KILL \"$run_pid\" 2>/dev/null",
		"  fi",
		"}",
		"set +e",
		"if command -v setsid >/dev/null 2>&1; then",
		fmt.Sprintf("  setsid sh -c %s &", shellQuote(measuredCommand)),
		"  run_pid=$!",
		"  run_group_pid=$run_pid",
		"else",
		fmt.Sprintf("  sh -c %s &", shellQuote(measuredCommand)),
		"  run_pid=$!",
		"  run_group_pid=",
		"fi",
		"while kill -0 \"$run_pid\" 2>/dev/null; do",
		"  sample_memory_usage",
		"  cpu_now=$(read_cpu_usage_usec)",
		"  case \"$cpu_now\" in ''|*[!0-9]*) cpu_now=$cpu_before ;; esac",
		"  exec_time_ms=$(((cpu_now - cpu_before + 999) / 1000))",
		"  if [ \"$exec_time_ms\" -gt \"$time_limit_ms\" ]; then",
		"    timed_out=1",
		"    kill_judge_command",
		"    break",
		"  fi",
		"  sleep 0.01",
		"done",
		"wait \"$run_pid\"",
		"status=$?",
		"sample_memory_usage",
		"set -e",
		"cpu_after=$(read_cpu_usage_usec)",
		"case \"$cpu_after\" in ''|*[!0-9]*) cpu_after=$cpu_before ;; esac",
		"exec_time_ms=$(((cpu_after - cpu_before + 999) / 1000))",
		"if [ \"$exec_time_ms\" -lt 0 ]; then exec_time_ms=0; fi",
		"memory_usage_kb=$(((max_memory + 1023) / 1024))",
		"printf '%s %s\\n' \"$exec_time_ms\" \"$memory_usage_kb\" > \"$metric_file\"",
		"memory_limit_bytes=$((memory_limit_kb * 1024))",
		"if [ \"$timed_out\" -eq 1 ]; then status=124; fi",
		"if [ \"$exec_time_ms\" -gt \"$time_limit_ms\" ] && [ \"$status\" -ne 153 ]; then",
		"  if [ \"$status\" -eq 137 ] && [ \"$memory_limit_bytes\" -gt 0 ] && [ \"$max_memory\" -ge \"$memory_limit_bytes\" ]; then",
		"    status=137",
		"  else",
		"    status=124",
		"  fi",
		"fi",
		"if [ \"$status\" -ne 0 ]; then exit \"$status\"; fi",
	}, "\n")
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

func normalizeJudgeTimeLimitMilliseconds(timeLimit int) int {
	if timeLimit <= 0 {
		return judgeDefaultTimeLimitSeconds * 1000
	}
	if timeLimit > 60 {
		return timeLimit
	}
	return timeLimit * 1000
}

func normalizeJudgeMemoryLimitMB(memoryLimit int) int {
	if memoryLimit <= 0 {
		return judgeDefaultMemoryLimitMB
	}
	return memoryLimit
}

func normalizeJudgeMemoryLimitKB(memoryLimit int) int {
	return normalizeJudgeMemoryLimitMB(memoryLimit) * 1024
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
	if strings.Contains(stderr, "CPU time limit exceeded") {
		return constants.JudgeStatusTimeLimitExceeded
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		switch exitErr.ExitCode() {
		case 124, 152:
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

func readJudgeMetric(path string) judgeMetric {
	var metric judgeMetric
	data, err := os.ReadFile(path)
	if err != nil {
		return metric
	}

	fields := strings.Fields(string(data))
	if len(fields) < 2 {
		return metric
	}

	execTime, execTimeErr := strconv.Atoi(fields[0])
	if execTimeErr == nil {
		metric.ExecTime = execTime
	}

	memoryUsage, memoryErr := strconv.Atoi(fields[1])
	if memoryErr == nil {
		metric.MemoryUsage = memoryUsage
	}

	return metric
}

func readMaxJudgeMetric(tempDir string, testCaseCount int) judgeMetric {
	var maxMetric judgeMetric
	for i := 0; i < testCaseCount; i++ {
		metric := readJudgeMetric(filepath.Join(tempDir, fmt.Sprintf("meta%d.txt", i)))
		if metric.ExecTime > maxMetric.ExecTime {
			maxMetric.ExecTime = metric.ExecTime
		}
		if metric.MemoryUsage > maxMetric.MemoryUsage {
			maxMetric.MemoryUsage = metric.MemoryUsage
		}
	}
	return maxMetric
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
