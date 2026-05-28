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

    <template v-else-if="problem">
      <header class="workspace-header">
        <RouterLink
          class="workspace-header__back"
          :to="{ name: 'CompetitionDetail', params: { competition_id: competitionId } }"
        >
          <i class="bi bi-arrow-left" aria-hidden="true"></i>
          返回竞赛
        </RouterLink>
        <div class="workspace-header__identity">
          <span>{{ problemNumber }}</span>
          <strong>{{ title }}</strong>
        </div>
        <div class="workspace-header__status" role="status">
          <i aria-hidden="true"></i>
          竞赛作答中
        </div>
      </header>

      <div class="problem-workspace__split">
        <article class="statement" aria-label="题目内容">
          <header class="statement__header">
            <h1>{{ problemNumber }}. {{ title }}</h1>
            <div v-if="constraintCards.length" class="statement__constraints" aria-label="题目限制">
              <div v-for="constraint in constraintCards" :key="constraint.label" class="constraint">
                <span>{{ constraint.label }}</span>
                <strong>{{ constraint.value }}</strong>
              </div>
            </div>
          </header>

          <div class="statement__body">
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
        </article>

        <CompetitionCodeEditor
          :key="`${competitionId}:${problemNumber}`"
          :competition-id="competitionId"
          :problem-number="problemNumber"
          :problem-id="problemId"
          :initial-input="firstSampleInput"
          :languages="problemLanguageSource"
          :theme="editorTheme"
          judge-mode="ACM"
          @update:theme="editorTheme = $event"
        />
      </div>
    </template>
  </main>
</template>

<script setup>
import { computed, defineComponent, h, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import api from '@/api'
import CompetitionCodeEditor from './CompetitionCodeEditor.vue'
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

const route = useRoute()
const loading = ref(false)
const error = ref('')
const problem = ref(null)
const editorTheme = ref('vs-light')
const message = useMessage()
const competitionId = computed(() => String(route.params.competition_id))
const problemNumber = computed(() => String(route.params.problem_number))
const bootstrapTheme = computed(() => editorTheme.value === 'vs-dark' ? 'dark' : 'light')
const workspaceThemeClass = computed(() => editorTheme.value === 'vs-dark' ? 'problem-workspace--dark' : 'problem-workspace--light')
const title = computed(() => problem.value?.title || problem.value?.problem_title || '无标题')
const problemId = computed(() => problem.value?.problem_id ?? problem.value?.id ?? null)
const problemLanguageSource = computed(() => problem.value?.language || problem.value?.languages || problem.value?.constraints || [])
const timeLimit = computed(() => formatConstraint('timeLimit', formatMilliseconds))
const memoryLimit = computed(() => formatConstraint('memoryLimit', value => `${value} MB`))
const submissionLimit = computed(() => formatConstraint('maxSubmissions', value => `${value} 次`))
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
const constraintLanguages = computed(() => {
  const languages = extractLanguages(problemLanguageSource.value)
  if (languages.length) {
    return languages.map(language => ({
      key: toMonacoLanguage(language),
      label: language
    }))
  }

  const constraints = problem.value?.constraints
  return constraints && typeof constraints === 'object'
    ? Object.keys(constraints).map(language => ({ key: language, label: language }))
    : []
})

function formatMilliseconds(value) {
  return value % 1000 === 0 ? `${value / 1000} 秒` : `${value} 毫秒`
}

function formatConstraint(field, formatter) {
  const constraints = problem.value?.constraints
  if (!constraints) {
    return ''
  }

  const groups = new Map()
  constraintLanguages.value.forEach(({ key, label }) => {
    const value = constraints[key]?.[field]
    if (value === undefined || value === null) {
      return
    }

    const formattedValue = formatter(value)
    const languages = groups.get(formattedValue) || []
    languages.push(label)
    groups.set(formattedValue, languages)
  })

  return [...groups.entries()]
    .map(([value, languages]) => `${languages.join('/')} ${value}`)
    .join('，')
}

async function copySample(text, label) {
  try {
    await navigator.clipboard.writeText(text || '')
    message.success(`${label}已复制`)
  } catch (error) {
    message.error('复制失败，请手动选择内容复制')
  }
}

watch(editorTheme, (theme) => {
  document.documentElement.dataset.lojEditorTheme = theme === 'vs-dark' ? 'dark' : 'light'
}, { immediate: true })

onBeforeUnmount(() => {
  delete document.documentElement.dataset.lojEditorTheme
})

async function loadProblem() {
  const requestedCompetitionId = competitionId.value
  const requestedProblemNumber = problemNumber.value
  loading.value = true
  error.value = ''
  problem.value = null

  try {
    const resp = await api.getCompetitionProblem(requestedCompetitionId, requestedProblemNumber)
    if (requestedCompetitionId !== competitionId.value || requestedProblemNumber !== problemNumber.value) {
      return
    }

    if (resp?.code === 0 && resp.data && typeof resp.data === 'object') {
      problem.value = resp.data
      return
    }

    error.value = resp?.message || '题目加载失败，请稍后重试。'
  } catch (requestError) {
    if (requestedCompetitionId === competitionId.value && requestedProblemNumber === problemNumber.value) {
      error.value = requestError?.message || '题目加载失败，请稍后重试。'
    }
  } finally {
    if (requestedCompetitionId === competitionId.value && requestedProblemNumber === problemNumber.value) {
      loading.value = false
    }
  }
}

watch([competitionId, problemNumber], loadProblem, { immediate: true })
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
  position: fixed;
  z-index: 1050;
  inset: 0;
  min-height: 100vh;
  overflow: hidden;
  color: var(--workspace-ink);
  font-family: Inter, "PingFang SC", "Microsoft YaHei", Arial, sans-serif;
  background: var(--workspace-cloud);
}

.problem-workspace--light {
  --workspace-ink: #111827;
  --workspace-copy: #344256;
  --workspace-muted: #718096;
  --workspace-border: #e5ebf2;
  --workspace-cloud: #f2f6fb;
}

.problem-workspace--dark {
  --workspace-ink: #edf4ff;
  --workspace-copy: #c6d4e4;
  --workspace-muted: #91a2b8;
  --workspace-border: #263449;
  --workspace-cloud: #0b1220;
}

.problem-workspace--dark .statement,
.problem-workspace--dark .constraint,
.problem-workspace--dark .sample__block {
  background: #111827;
  border-color: var(--workspace-border);
}

.problem-workspace--dark .constraint {
  background: #162235;
}

.problem-workspace--dark .constraint strong,
.problem-workspace--dark .sample__block pre,
.problem-workspace--dark :deep(.statement-section__content) {
  color: #dbe7f5;
}

.problem-workspace--light .workspace-header {
  color: #475569;
  background: #fff;
  border-bottom-color: #dfe7f0;
  box-shadow: 0 0.55rem 1.35rem rgba(15, 23, 42, 0.06);
}

.problem-workspace--light .workspace-header__back,
.problem-workspace--light .workspace-header__status {
  color: #52657a;
}

.problem-workspace--light .workspace-header__back:hover {
  color: #0f766e;
}

.problem-workspace--light .workspace-header__identity strong {
  color: #111827;
}

.problem-workspace__state {
  display: flex;
  gap: 0.8rem;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
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

.workspace-header {
  display: grid;
  grid-template-columns: minmax(9.5rem, 1fr) minmax(16rem, auto) minmax(9.5rem, 1fr);
  gap: 1.2rem;
  align-items: center;
  height: 68px;
  padding: 0 1.35rem;
  color: #c6d1df;
  background: var(--workspace-navy);
  border-bottom: 1px solid #1c273a;
}

.workspace-header__back {
  display: inline-flex;
  gap: 0.62rem;
  align-items: center;
  justify-self: start;
  color: #b3bfd0;
  font-size: 0.9rem;
  font-weight: 500;
  text-decoration: none;
  transition: color 160ms ease;
}

.workspace-header__back i {
  font-size: 1rem;
}

.workspace-header__back:hover {
  color: #fff;
}

.workspace-header__back:focus-visible {
  outline: 2px solid rgba(21, 184, 166, 0.62);
  outline-offset: 6px;
  border-radius: 0.25rem;
}

.workspace-header__identity {
  display: inline-flex;
  gap: 0.75rem;
  align-items: center;
  justify-self: center;
  min-width: 0;
}

.workspace-header__identity span {
  color: #52dfcc;
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.1em;
}

.workspace-header__identity strong {
  overflow: hidden;
  max-width: min(37vw, 34rem);
  color: #eef4fc;
  font-size: 0.95rem;
  font-weight: 560;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.workspace-header__status {
  display: inline-flex;
  gap: 0.52rem;
  align-items: center;
  justify-self: end;
  color: #94a5bc;
  font-size: 0.83rem;
}

.workspace-header__status i {
  width: 0.46rem;
  height: 0.46rem;
  background: var(--workspace-teal);
  border-radius: 50%;
  box-shadow: 0 0 0 0.28rem rgba(21, 184, 166, 0.16);
}

.problem-workspace__split {
  display: grid;
  grid-template-columns: minmax(22rem, 0.86fr) minmax(29rem, 1.14fr);
  gap: 0.8rem;
  height: calc(100vh - 68px);
  padding: 0.8rem;
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

.statement__header {
  padding: 2.2rem 2.15rem 1.45rem;
  border-bottom: 1px solid var(--workspace-border);
}

.statement__header h1 {
  margin: 0 0 1.6rem;
  color: var(--workspace-ink);
  font-size: clamp(1.55rem, 2.1vw, 1.9rem);
  font-weight: 660;
  line-height: 1.3;
  letter-spacing: -0.035em;
}

.statement__constraints {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.7rem;
}

.constraint {
  min-height: 4.15rem;
  padding: 0.72rem 0.82rem;
  background: #f7f9fc;
  border: 1px solid #edf1f6;
  border-radius: 0.72rem;
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

.statement__body {
  padding: 1.65rem 2.15rem 2.8rem;
}

:deep(.statement-section) {
  margin: 0 0 2rem;
}

:deep(.statement-section h2) {
  display: flex;
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

.sample {
  display: grid;
  gap: 0.78rem;
}

.sample__block {
  overflow: hidden;
  background: #f7f9fc;
  border: 1px solid #e8eef5;
  border-radius: 0.72rem;
}

.sample__title {
  display: flex;
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

.sample-copy-btn {
  display: inline-flex;
  width: 1.9rem;
  height: 1.9rem;
  align-items: center;
  justify-content: center;
  padding: 0;
  color: var(--workspace-muted);
  border-radius: 999px;
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

  .workspace-header {
    position: sticky;
    z-index: 2;
    top: 0;
    grid-template-columns: auto minmax(0, 1fr);
    height: 58px;
    padding: 0 0.85rem;
  }

  .workspace-header__back {
    font-size: 0.8rem;
  }

  .workspace-header__back i {
    font-size: 1.05rem;
  }

  .workspace-header__identity {
    justify-self: start;
  }

  .workspace-header__identity strong {
    max-width: calc(100vw - 7rem);
  }

  .workspace-header__status {
    display: none;
  }

  .problem-workspace__split {
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
    height: auto;
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

  .statement__header h1 {
    margin-bottom: 1.05rem;
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

  :deep(.contest-editor) {
    min-height: 41rem;
  }
}
</style>
