package problem_model

type RunResult struct {
	Passed      bool   `json:"passed"`
	Output      string `json:"output"`
	Expected    string `json:"expected"`
	Error       string `json:"error"`
	ExecTime    int    `json:"exec_time"`    // ms
	MemoryUsage int    `json:"memory_usage"` // KB
}
