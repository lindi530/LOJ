package problem_service

import "GO1/models/problem_model"

func SummarizeRunResultMetrics(results []problem_model.RunResult) (execTime, memoryUsage int) {
	for _, result := range results {
		if result.ExecTime > execTime {
			execTime = result.ExecTime
		}
		if result.MemoryUsage > memoryUsage {
			memoryUsage = result.MemoryUsage
		}
	}
	return execTime, memoryUsage
}
