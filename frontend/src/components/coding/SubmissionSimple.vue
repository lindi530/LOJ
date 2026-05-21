<template>
  <article class="submission-row">
    <router-link
      :to="submissionDetailRoute"
      target="_blank"
      class="run-id"
    >
      {{ runId }}
    </router-link>

    <router-link
      :to="problemRoute"
      target="_blank"
      class="problem-cell"
    >
      <span class="problem-title">{{ problemTitle }}</span>
    </router-link>

    <router-link
      :to="submissionDetailRoute"
      target="_blank"
      class="result-cell"
      :class="resultClass"
    >
      {{ resultLabel }}
    </router-link>

    <span class="metric-cell" data-label="运行时间(ms)">{{ execTimeLabel }}</span>
    <span class="metric-cell" data-label="使用内存(KB)">{{ memoryLabel }}</span>
    <span class="metric-cell" data-label="代码长度">{{ codeLengthLabel }}</span>
    <span class="metric-cell" data-label="使用语言">{{ languageLabel }}</span>
    <span class="time-cell" data-label="提交时间">{{ formatSubmissionDate(submission.created_at) }}</span>
  </article>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  author: {
    type: Object,
    required: true
  },
  submission: {
    type: Object,
    required: true
  }
})

const runId = computed(() => props.submission.id || props.submission.submission_id || '-')
const problemTitle = computed(() => props.submission.title || '未命名题目')
const languageLabel = computed(() => props.submission.language || props.submission.lang || '-')
const codeLengthLabel = computed(() => {
  const length = props.submission.code_length || props.submission.codeLength

  if (length !== undefined && length !== null) {
    return length
  }

  return props.submission.code ? props.submission.code.length : '-'
})

const execTimeLabel = computed(() => {
  const value = props.submission.exec_time ?? props.submission.execTime ?? props.submission.runtime
  return value === undefined || value === null || value === '' ? '-' : value
})

const memoryLabel = computed(() => {
  const value = props.submission.memory_usage ?? props.submission.memoryUsage ?? props.submission.memory
  return value === undefined || value === null || value === '' ? '-' : value
})

const resultLabel = computed(() => {
  const state = props.submission.state || props.submission.status || props.submission.result

  if (state) {
    return normalizeResult(state)
  }

  if (props.submission.score !== undefined && props.submission.score !== null) {
    return Number(props.submission.score) > 0 ? '答案正确' : '未通过'
  }

  return '查看详情'
})

const resultClass = computed(() => {
  const result = resultLabel.value

  if (['答案正确', 'Accepted', 'AC'].includes(result)) {
    return 'is-accepted'
  }

  if (result === '查看详情') {
    return 'is-pending'
  }

  return 'is-failed'
})

const problemRoute = computed(() => ({
  name: 'ProblemDetail',
  params: {
    problem_id: props.submission.problem_id
  }
}))

const submissionDetailRoute = computed(() => ({
  name: 'SubmissionDetail',
  params: {
    submission_id: runId.value
  }
}))

function normalizeResult(state) {
  const map = {
    accepted: '答案正确',
    ac: '答案正确',
    wrong_answer: '答案错误',
    wa: '答案错误',
    runtime_error: '运行错误',
    re: '运行错误',
    time_limit_exceeded: '运行超时',
    tle: '运行超时',
    memory_limit_exceeded: '内存超限',
    mle: '内存超限',
    compile_error: '编译错误',
    ce: '编译错误'
  }

  const key = String(state).trim().toLowerCase()
  return map[key] || state
}

function formatSubmissionDate(dateStr) {
  if (!dateStr) {
    return '-'
  }

  const date = new Date(dateStr)
  if (Number.isNaN(date.getTime())) {
    return dateStr
  }

  const pad = value => String(value).padStart(2, '0')
  return [
    date.getFullYear(),
    pad(date.getMonth() + 1),
    pad(date.getDate())
  ].join('-') + ` ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}
</script>

<style scoped>
.submission-row {
  display: grid;
  min-width: 780px;
  grid-template-columns: 72px minmax(112px, 1.25fr) 88px 92px 94px 72px 78px 150px;
  align-items: center;
  padding: 13px 16px;
  border-bottom: 1px dashed #e2e5ea;
  color: #1f2328;
  font-size: 14px;
  line-height: 1.35;
}

.submission-row:last-child {
  border-bottom: 0;
}

.run-id,
.problem-cell,
.result-cell {
  color: inherit;
  text-decoration: none;
}

.run-id {
  color: #2f6df6;
  font-family: Consolas, Monaco, monospace;
  font-size: 13px;
}

.problem-cell {
  display: flex;
  min-width: 0;
  align-items: center;
}

.problem-title {
  min-width: 0;
  overflow: hidden;
  color: #1f2328;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-cell {
  font-weight: 650;
}

.result-cell.is-accepted {
  color: #00a35a;
}

.result-cell.is-failed {
  color: #d03050;
}

.result-cell.is-pending {
  color: #637083;
}

.metric-cell,
.time-cell {
  color: #333b46;
}

.time-cell {
  color: #637083;
  font-size: 13px;
}

.run-id:hover,
.problem-cell:hover .problem-title,
.result-cell:hover {
  color: #00a35a;
}

@media (max-width: 900px) {
  .submission-row {
    min-width: 0;
    grid-template-columns: 1fr;
    gap: 10px;
    padding: 16px;
    margin-bottom: 12px;
    border: 1px solid #e6e9ee;
    border-radius: 6px;
    background: #ffffff;
  }

  .run-id::before {
    content: "运行ID ";
    color: #8a939d;
    font-family: Avenir, Helvetica, Arial, sans-serif;
  }

  .metric-cell,
  .time-cell {
    display: flex;
    justify-content: space-between;
    gap: 16px;
  }

  .metric-cell::before,
  .time-cell::before {
    content: attr(data-label);
    color: #8a939d;
  }
}
</style>
