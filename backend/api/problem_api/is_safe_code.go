package problem_api

import "strings"

func IsSafeCode(code, language string) bool {
	language = strings.ToLower(language)
	dangerMap := map[string][]string{
		"c":      {"system", "exec", "fork", "popen"},
		"cpp":    {"system", "exec", "fork", "popen"},
		"python": {"os.system", "subprocess", "eval", "exec", "shutil"},
		"java":   {"Runtime.getRuntime", "ProcessBuilder", "FilePermission"},
		"go":     {"os/exec", "syscall", "ioutil.WriteFile"},
	}

	if banned, ok := dangerMap[language]; ok {
		for _, keyword := range banned {
			if strings.Contains(code, keyword) {
				return false
			}
		}
	}
	return true
}
