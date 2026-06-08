<template>
  <main
    class="problem-workspace"
    :class="workspaceThemeClass"
    :data-bs-theme="bootstrapTheme"
  >
    <section v-if="loading" class="problem-workspace__state" role="status">
      <span class="problem-workspace__loader" aria-hidden="true"></span>
      <span>题目加载中...</span>
    </section>
    <section v-else-if="error" class="problem-workspace__state problem-workspace__state--error" role="alert">
      <i class="bi bi-exclamation-circle" aria-hidden="true"></i>
      {{ error }}
    </section>

    <div v-else-if="problem" class="problem-workspace__split">
      <article class="statement" aria-label="题目内容">
        <header class="statement__header">
          <div class="statement__title-row">
            <div>
              <span class="statement__eyebrow">Problem #{{ problemID }}</span>
              <h1>{{ title }}</h1>
            </div>
            <div class="statement__labels" aria-label="题目标签">
              <span
                v-if="problem.level"
                class="problem-badge"
                :class="levelBadgeClass(problem.level)"
              >
                {{ problem.level }}
              </span>
              <span v-for="tag in problemTags" :key="tag" class="problem-badge problem-badge--tag">
                {{ tag }}
              </span>
            </div>
          </div>

          <div v-if="constraintCards.length" class="statement__constraints" aria-label="题目限制">
            <div v-for="constraint in constraintCards" :key="constraint.label" class="constraint">
              <span>{{ constraint.label }}</span>
              <strong>{{ constraint.value }}</strong>
            </div>
          </div>

          <nav v-if="canShowSubmissions" class="statement__tabs" aria-label="题目页面">
            <button
              type="button"
              :class="{ 'is-active': activePane === 'problem' }"
              @click="activePane = 'problem'"
            >
              题目信息
            </button>
            <button
              type="button"
              :class="{ 'is-active': activePane === 'submissions' }"
              @click="activePane = 'submissions'"
            >
              提交记录
            </button>
          </nav>
        </header>

        <div v-if="activePane === 'problem'" class="statement__body">
          <StatementSection title="题目描述" :content="problem.description" />
          <StatementSection title="输入描述" :content="problem.input_desc || problem.input_description" />
          <StatementSection title="输出描述" :content="problem.output_desc || problem.output_description" />

          <section v-for="(example, index) in examples" :key="index" class="statement-section">
            <h2>示例 {{ index + 1 }}</h2>
            <div class="sample">
              <div class="sample__block">
                <div class="sample__title">
                  <span>输入</span>
                  <button
                    class="btn btn-sm btn-link sample-copy-btn"
                    type="button"
                    :aria-label="`复制示例 ${index + 1} 输入`"
                    title="复制"
                    @click="copySample(example.input, `示例 ${index + 1} 输入`)"
                  >
                    <i class="bi bi-copy" aria-hidden="true"></i>
                  </button>
                </div>
                <pre>{{ example.input }}</pre>
              </div>
              <div class="sample__block">
                <div class="sample__title">
                  <span>输出</span>
                  <button
                    class="btn btn-sm btn-link sample-copy-btn"
                    type="button"
                    :aria-label="`复制示例 ${index + 1} 输出`"
                    title="复制"
                    @click="copySample(example.output, `示例 ${index + 1} 输出`)"
                  >
                    <i class="bi bi-copy" aria-hidden="true"></i>
                  </button>
                </div>
                <pre>{{ example.output }}</pre>
              </div>
            </div>
          </section>
        </div>

        <div v-else class="statement__body statement__body--submissions">
          <SubmissionList
            :problem-id="problemID"
            :room-id="roomId"
          />
        </div>
      </article>

      <CodeEditor
        :problem-id="problemID"
        :room-id="roomId"
        :initial-input="firstSampleInput"
        :languages="problemLanguageSource"
        :theme="editorTheme"
        @update:theme="editorTheme = $event"
      />
    </div>
  </main>
</template>

<script setup>
import { computed, defineComponent, h, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import CodeEditor from '@/components/coding/CodeEditor.vue'
import SubmissionList from '@/components/coding/SubmissionList.vue'
import api from '@/api'
import { extractLanguages, toMonacoLanguage } from '@/utils/languageOptions'

const StatementSection = defineComponent({
  props: {
    title: {
      type: String,
      required: true
    },
    content: {
      type: String,
      default: ''
    }
  },
  setup(props) {
    return () => props.content
      ? h('section', { class: 'statement-section' }, [
          h('h2', props.title),
          h('div', { class: 'statement-section__content', innerHTML: props.content })
        ])
      : null
  }
})

const props = defineProps({
  problemId: {
    type: Number,
    default: undefined
  },
  roomId: {
    type: String,
    default: undefined
  }
})

const route = useRoute()
const message = useMessage()
const problem = ref(null)
const loading = ref(true)
const error = ref('')
const editorTheme = ref('vs-light')
const activePane = ref('problem')
const problemID = ref(
  typeof route.params.problem_id !== 'undefined'
    ? Number(route.params.problem_id)
    : props.problemId
)

const bootstrapTheme = computed(() => editorTheme.value === 'vs-dark' ? 'dark' : 'light')
const workspaceThemeClass = computed(() => editorTheme.value === 'vs-dark' ? 'problem-workspace--dark' : 'problem-workspace--light')
const canShowSubmissions = computed(() => !props.roomId || props.roomId.trim() === '')
const title = computed(() => problem.value?.title || problem.value?.problem_title || '无标题')
const problemTags = computed(() => {
  const tags = problem.value?.tags
  if (Array.isArray(tags)) {
    return tags.filter(Boolean)
  }
  return tags ? [tags] : []
})
const problemLanguageSource = computed(() => problem.value?.language || problem.value?.languages || problem.value?.constraints || [])
const constraintLanguages = computed(() => {
  const languages = extractLanguages(problemLanguageSource.value)
  if (languages.length) {
    return languages.map(language => ({
      key: toMonacoLanguage(language),
      label: language
    }))
  }

  const constraints = problem.value?.constraints
  return constraints && typeof constraints === 'object' && !Array.isArray(constraints)
    ? Object.keys(constraints).map(language => ({ key: language, label: language }))
    : []
})
const timeLimit = computed(() => formatConstraint(['timeLimit', 'time_limit'], formatMilliseconds))
const memoryLimit = computed(() => formatConstraint(['memoryLimit', 'memory_limit'], value => `${value} MB`))
const submissionLimit = computed(() => formatConstraint(['maxSubmissions', 'max_submissions', 'submissionLimit', 'submission_limit'], value => `${value} 次`))
const constraintCards = computed(() => [
  { label: '时间限制', value: timeLimit.value },
  { label: '空间限制', value: memoryLimit.value },
  { label: '提交限制', value: submissionLimit.value }
].filter(constraint => constraint.value))
const examples = computed(() => {
  const data = problem.value?.examples
  return Array.isArray(data) ? data.filter(example => example.is_public !== false) : []
})
const firstSampleInput = computed(() => examples.value[0]?.input || '')

watch(canShowSubmissions, (visible) => {
  if (!visible) {
    activePane.value = 'problem'
  }
}, { immediate: true })

watch(editorTheme, (theme) => {
  document.documentElement.dataset.lojEditorTheme = theme === 'vs-dark' ? 'dark' : 'light'
}, { immediate: true })

onBeforeUnmount(() => {
  delete document.documentElement.dataset.lojEditorTheme
})

onMounted(loadProblem)

async function loadProblem() {
  loading.value = true
  error.value = ''

  try {
    const res = await api.getProblemDetail(problemID.value)
    if (res.code === 0) {
      problem.value = res.data
      return
    }

    error.value = res.message || '题目加载失败，请稍后重试。'
  } catch (requestError) {
    error.value = requestError?.message || '题目加载失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

async function copySample(text, label) {
  try {
    await navigator.clipboard.writeText(text || '')
    message.success(`${label}已复制`)
  } catch {
    message.error('复制失败，请手动选择内容复制')
  }
}

function levelBadgeClass(level) {
  switch (level) {
    case '简单':
      return 'problem-badge--easy'
    case '中等':
      return 'problem-badge--medium'
    case '困难':
      return 'problem-badge--hard'
    default:
      return 'problem-badge--default'
  }
}

function formatMilliseconds(value) {
  const numberValue = Number(value)
  if (!Number.isFinite(numberValue)) {
    return String(value)
  }

  return numberValue % 1000 === 0 ? `${numberValue / 1000} 秒` : `${numberValue} 毫秒`
}

function formatConstraint(fields, formatter) {
  const constraints = problem.value?.constraints
  if (!constraints) {
    return ''
  }

  const groups = new Map()
  constraintLanguages.value.forEach(({ key, label }) => {
    const constraint = findLanguageConstraint(constraints, key, label)
    if (!constraint) {
      return
    }

    const value = findConstraintValue(constraint, fields)
    if (value === undefined || value === null || value === '') {
      return
    }

    const formattedValue = formatter(value)
    const languages = groups.get(formattedValue) || []
    languages.push(label)
    groups.set(formattedValue, languages)
  })

  return [...groups.entries()]
    .map(([value, languages]) => `${languages.join('/')} ${value}`)
    .join('；')
}

function findLanguageConstraint(constraints, key, label) {
  if (Array.isArray(constraints)) {
    return constraints.find(constraint => {
      const language = String(constraint?.language || '').trim()
      return language === label || toMonacoLanguage(language) === key
    })
  }

  if (typeof constraints !== 'object') {
    return null
  }

  return constraints[key] || constraints[label] || constraints[toMonacoLanguage(label)] || null
}

function findConstraintValue(constraint, fields) {
  return fields.map(field => constraint?.[field]).find(value => value !== undefined && value !== null)
}
</script>

<style scoped>
.problem-workspace {
  --workspace-navy: #0b1220;
  --workspace-ink: #111827;
  --workspace-copy: #344256;
  --workspace-muted: #718096;
  --workspace-border: #e5ebf2;
  --workspace-cloud: #f2f6fb;
  --workspace-teal: #15b8a6;
  min-height: calc(100vh - 60px);
  overflow: hidden;
  color: var(--workspace-ink);
  font-family: Inter, "PingFang SC", "Microsoft YaHei", Arial, sans-serif;
  background: var(--workspace-cloud);
}

.problem-workspace--dark {
  --workspace-ink: #edf4ff;
  --workspace-copy: #c6d4e4;
  --workspace-muted: #91a2b8;
  --workspace-border: #263449;
  --workspace-cloud: #0b1220;
}

.problem-workspace__state {
  display: flex;
  gap: 0.8rem;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 60px);
  color: #607189;
  font-size: 0.95rem;
}

.problem-workspace__loader {
  width: 1.1rem;
  height: 1.1rem;
  border: 2px solid #d5e1ee;
  border-top-color: var(--workspace-teal);
  border-radius: 50%;
  animation: spin 700ms linear infinite;
}

.problem-workspace__state--error {
  color: #ba4b57;
}

.statement {
  min-width: 0;
  overflow-y: auto;
  color: var(--workspace-copy);
  background: #fff;
  border: 1px solid var(--workspace-border);
  border-radius: 1rem;
  box-shadow: 0 0.55rem 1.4rem rgba(15, 23, 42, 0.035);
}

.problem-workspace__split {
  display: grid;
  grid-template-columns: minmax(22rem, 0.86fr) minmax(29rem, 1.14fr);
  gap: 0.8rem;
  height: calc(100vh - 60px);
  padding: 0.8rem;
}

.problem-workspace--dark .statement,
.problem-workspace--dark .constraint,
.problem-workspace--dark .sample__block {
  background: #111827;
  border-color: var(--workspace-border);
}

.statement__header {
  padding: 2.2rem 2.15rem 1.2rem;
  border-bottom: 1px solid var(--workspace-border);
}

.statement__title-row,
.statement__labels,
.sample__title {
  display: flex;
  align-items: flex-start;
}

.statement__title-row {
  gap: 1rem;
  justify-content: space-between;
  margin-bottom: 1.45rem;
}

.statement__eyebrow {
  display: block;
  margin-bottom: 0.35rem;
  color: var(--workspace-teal);
  font-size: 0.74rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.statement__header h1 {
  margin: 0;
  color: var(--workspace-ink);
  font-size: clamp(1.55rem, 2.1vw, 1.9rem);
  font-weight: 660;
  line-height: 1.3;
}

.statement__labels {
  flex-wrap: wrap;
  gap: 0.45rem;
  justify-content: flex-end;
  max-width: 16rem;
}

.problem-badge {
  display: inline-flex;
  align-items: center;
  min-height: 1.7rem;
  padding: 0 0.62rem;
  color: #52657a;
  font-size: 0.72rem;
  font-weight: 650;
  background: #f4f7fb;
  border: 1px solid #e4ebf3;
  border-radius: 999px;
}

.statement__constraints,
.sample {
  display: grid;
  gap: 0.7rem;
}

.statement__constraints {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.constraint {
  min-height: 4.15rem;
  padding: 0.72rem 0.82rem;
  background: #f7f9fc;
  border: 1px solid #edf1f6;
  border-radius: 0.72rem;
}

.problem-workspace--dark .constraint {
  background: #162235;
}

.problem-badge--easy {
  color: #0f766e;
  background: rgba(21, 184, 166, 0.1);
  border-color: rgba(21, 184, 166, 0.28);
}

.problem-badge--medium {
  color: #9a6a11;
  background: rgba(245, 158, 11, 0.12);
  border-color: rgba(245, 158, 11, 0.28);
}

.problem-badge--hard {
  color: #b4233c;
  background: rgba(244, 63, 94, 0.1);
  border-color: rgba(244, 63, 94, 0.25);
}

.problem-badge--tag {
  color: #3b5f86;
  background: #eef6ff;
  border-color: #d7e7f7;
}

.constraint span {
  display: block;
  margin-bottom: 0.35rem;
  color: var(--workspace-muted);
  font-size: 0.72rem;
  font-weight: 560;
}

.constraint strong {
  display: block;
  overflow: hidden;
  color: #253247;
  font-size: 0.8rem;
  font-weight: 590;
  line-height: 1.4;
  text-overflow: ellipsis;
}

.problem-workspace--dark .constraint strong,
.problem-workspace--dark .sample__block pre,
.problem-workspace--dark :deep(.statement-section__content) {
  color: #dbe7f5;
}

.statement__tabs,
.statement__tabs button,
:deep(.statement-section h2) {
  display: flex;
}

.statement__tabs {
  gap: 1.2rem;
  margin: 1.1rem 0 -1.2rem;
}

.statement__tabs button {
  position: relative;
  min-height: 2.6rem;
  padding: 0;
  color: var(--workspace-muted);
  font-size: 0.86rem;
  font-weight: 590;
  background: transparent;
  border: 0;
}

.statement__tabs button.is-active {
  color: var(--workspace-ink);
}

.statement__tabs button.is-active::after {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  height: 2px;
  background: var(--workspace-teal);
  content: "";
}

:deep(.statement-section) {
  margin: 0 0 2rem;
}

:deep(.statement-section h2) {
  gap: 0.6rem;
  align-items: center;
  margin: 0 0 0.92rem;
  color: var(--workspace-ink);
  font-size: 1.04rem;
  font-weight: 640;
}

:deep(.statement-section h2::before) {
  width: 0.21rem;
  height: 1rem;
  background: var(--workspace-teal);
  border-radius: 99px;
  content: "";
}

:deep(.statement-section__content) {
  color: #45566d;
  font-size: 0.92rem;
  line-height: 1.88;
  white-space: pre-wrap;
}

:deep(.statement-section__content p) {
  margin: 0 0 0.6rem;
}

.statement__body {
  padding: 1.65rem 2.15rem 2.8rem;
}

.statement__body--submissions {
  padding: 1rem;
}

.sample__block {
  overflow: hidden;
  background: #f7f9fc;
  border: 1px solid #e8eef5;
  border-radius: 0.72rem;
}

.sample__title {
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  min-height: 2.3rem;
  padding: 0.4rem 0.75rem 0;
}

.sample__title span {
  color: var(--workspace-muted);
  font-size: 0.73rem;
  font-weight: 600;
}

.sample-copy-btn,
.problem-workspace__loader {
  border-radius: 999px;
}

.sample-copy-btn {
  display: inline-flex;
  width: 1.9rem;
  height: 1.9rem;
  align-items: center;
  justify-content: center;
  padding: 0;
  color: var(--workspace-muted);
  text-decoration: none;
}

.sample-copy-btn:hover {
  color: var(--workspace-teal);
  background: rgba(21, 184, 166, 0.1);
}

.sample__block pre {
  min-height: 3.2rem;
  padding: 0.58rem 0.86rem 0.88rem;
  margin: 0;
  color: #27364a;
  font: 0.85rem/1.68 "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  white-space: pre-wrap;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

  @media (max-width: 959.98px) {
  .problem-workspace {
    overflow-y: auto;
  }

  .problem-workspace__split {
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
    height: auto;
    min-height: calc(100vh - 60px);
    padding: 0.62rem;
  }

  .statement {
    max-height: none;
    overflow-y: visible;
    border-radius: 0.82rem;
  }

  .statement__header {
    padding: 1.35rem 1rem 1rem;
  }

  .statement__title-row {
    flex-direction: column;
    margin-bottom: 1.05rem;
  }

  .statement__labels {
    justify-content: flex-start;
    max-width: none;
  }

  .statement__header h1 {
    font-size: 1.36rem;
  }

  .statement__constraints {
    gap: 0.42rem;
  }

  .constraint {
    min-height: 3.7rem;
    padding: 0.55rem;
  }

  .constraint strong {
    font-size: 0.72rem;
  }

  .statement__body {
    padding: 1.1rem 1rem 1.45rem;
  }

  :deep(.statement-section) {
    margin-bottom: 1.3rem;
  }
}
</style>
