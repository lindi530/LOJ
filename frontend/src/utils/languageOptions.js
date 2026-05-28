const templateByMonacoLanguage = {
  cpp: '#include <bits/stdc++.h>\nusing namespace std;\n\nint main() {\n    ios::sync_with_stdio(false);\n    cin.tie(nullptr);\n\n    return 0;\n}\n',
  python: 'def main():\n    pass\n\n\nif __name__ == "__main__":\n    main()\n',
  go: 'package main\n\nimport (\n\t"bufio"\n\t"fmt"\n\t"os"\n)\n\nfunc main() {\n\tin := bufio.NewReader(os.Stdin)\n\t_ = in\n\tfmt.Println()\n}\n',
  java: 'public class Main {\n    public static void main(String[] args) {\n    }\n}\n'
}

const fileNameByMonacoLanguage = {
  cpp: 'main.cpp',
  python: 'main.py',
  go: 'main.go',
  java: 'Main.java'
}

export function toLanguageOptions(source) {
  return extractLanguages(source).map(language => ({
    label: language,
    value: language
  }))
}

export function extractLanguages(source) {
  const languages = []
  collectLanguages(source, languages)
  return [...new Set(languages.map(language => String(language).trim()).filter(Boolean))]
}

export function toMonacoLanguage(language) {
  const raw = String(language || '').trim().toLowerCase()
  if (!raw) return ''
  if (raw === 'cpp' || raw === 'cxx' || raw.includes('c++')) return 'cpp'
  if (raw === 'py' || raw.includes('python')) return 'python'
  if (raw === 'go' || raw.includes('golang')) return 'go'
  if (raw.includes('java')) return 'java'
  return raw
}

export function defaultCodeFor(language) {
  return templateByMonacoLanguage[toMonacoLanguage(language)] || ''
}

export function fileNameFor(language) {
  return fileNameByMonacoLanguage[toMonacoLanguage(language)] || 'main.txt'
}

function collectLanguages(value, target) {
  if (!value) return

  if (typeof value === 'string') {
    target.push(value)
    return
  }

  if (Array.isArray(value)) {
    value.forEach(item => collectLanguages(item, target))
    return
  }

  if (typeof value !== 'object') return

  if (typeof value.language === 'string') {
    target.push(value.language)
    return
  }

  if (Array.isArray(value.language)) {
    collectLanguages(value.language, target)
  }

  if (Array.isArray(value.languages)) {
    collectLanguages(value.languages, target)
  }

  if (Array.isArray(value.language_options)) {
    collectLanguages(value.language_options, target)
  }

  if (Array.isArray(value.constraints)) {
    collectLanguages(value.constraints, target)
    return
  }

  const keyedConstraints = value.constraints && typeof value.constraints === 'object'
    ? value.constraints
    : value

  Object.entries(keyedConstraints).forEach(([key, constraint]) => {
    if (!constraint || typeof constraint !== 'object') return
    if (typeof constraint.language === 'string') {
      target.push(constraint.language)
      return
    }
    if ('timeLimit' in constraint || 'memoryLimit' in constraint || 'time_limit' in constraint || 'memory_limit' in constraint) {
      target.push(key)
    }
  })
}
