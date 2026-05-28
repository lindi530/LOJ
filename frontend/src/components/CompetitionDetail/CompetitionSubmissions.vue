<template>
  <div class="competition-submissions">
    <p v-if="loading" class="competition-submissions__state" role="status">
      <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
      提交记录加载中...
    </p>
    <p v-else-if="error" class="competition-submissions__state competition-submissions__state--error" role="alert">
      {{ error }}
    </p>
    <p v-else-if="records.length === 0" class="competition-submissions__state">暂无提交记录。</p>

    <div v-else class="submission-table" role="table" aria-label="竞赛提交记录列表">
      <div class="submission-table__row submission-table__head" role="row">
        <span role="columnheader">题目</span>
        <span role="columnheader">状态</span>
        <span role="columnheader">得分</span>
        <span role="columnheader">语言</span>
        <span role="columnheader">运行时间</span>
        <span role="columnheader">提交时间</span>
      </div>
      <div v-for="record in records" :key="record.id" class="submission-table__row" role="row">
        <strong class="submission-table__problem" role="cell">{{ problemLabel(record) }}</strong>
        <span class="submission-state" :class="stateClass(record.state || record.status)" role="cell">
          {{ record.state || record.status || '评测中' }}
        </span>
        <span role="cell">{{ value(record.score) }}</span>
        <span role="cell">{{ record.language || '-' }}</span>
        <span role="cell">{{ record.exec_time || record.run_time || '-' }}</span>
        <time role="cell">{{ formatTime(record.created_at) }}</time>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import api from '@/api'

const props = defineProps({
  competitionId: {
    type: [Number, String],
    required: true
  }
})

const loading = ref(false)
const error = ref('')
const records = ref([])

function value(score) {
  return score === undefined || score === null ? '-' : score
}

function problemLabel(record) {
  return record.problem_number || record.problem_title || record.problem_id || '-'
}

function stateClass(state = '') {
  const normalized = state.toLowerCase()
  if (state === '通过' || normalized === 'accepted') {
    return 'submission-state--accepted'
  }

  if (state === '答案错误' || normalized === 'wrong answer' || normalized === 'wrong') {
    return 'submission-state--wrong'
  }

  return ''
}

function formatTime(input) {
  if (!input) {
    return '-'
  }

  const date = new Date(input)
  if (Number.isNaN(date.getTime())) {
    return input
  }

  return date.toLocaleString('zh-CN', { hour12: false })
}

async function loadRecords() {
  const requestedCompetitionId = props.competitionId
  loading.value = true
  error.value = ''

  try {
    const resp = await api.getCompetitionSubmissions(requestedCompetitionId)
    if (String(requestedCompetitionId) !== String(props.competitionId)) {
      return
    }

    if (resp.code === 0 && Array.isArray(resp.data)) {
      records.value = resp.data
      return
    }

    records.value = []
    error.value = resp.message || '提交记录加载失败，请稍后重试。'
  } catch (requestError) {
    if (String(requestedCompetitionId) === String(props.competitionId)) {
      records.value = []
      error.value = requestError?.message || '提交记录加载失败，请稍后重试。'
    }
  } finally {
    if (String(requestedCompetitionId) === String(props.competitionId)) {
      loading.value = false
    }
  }
}

watch(() => props.competitionId, loadRecords, { immediate: true })
</script>

<style scoped>
.competition-submissions__state {
  display: inline-flex;
  gap: 0.6rem;
  align-items: center;
  margin: 0;
  color: #74869b;
  font-size: 0.98rem;
}

.competition-submissions__state--error {
  color: #bf5454;
}

.submission-table {
  overflow: hidden;
  color: #26354a;
  border: 1px solid #e9eef4;
  border-radius: 0.65rem;
}

.submission-table__row {
  display: grid;
  grid-template-columns: minmax(4.5rem, 0.7fr) minmax(7.5rem, 1.1fr) minmax(4rem, 0.65fr) minmax(5rem, 0.8fr) minmax(6rem, 0.9fr) minmax(11rem, 1.4fr);
  gap: 0.85rem;
  align-items: center;
  min-height: 3.8rem;
  padding: 0 1.2rem;
  font-size: 0.92rem;
  background: #fff;
  border-bottom: 1px solid #edf1f5;
}

.submission-table__row:last-child {
  border-bottom: 0;
}

.submission-table__head {
  min-height: 3.25rem;
  color: #55677e;
  font-weight: 600;
  background: #f8fafc;
}

.submission-table__problem {
  color: #245fc9;
  font-weight: 600;
}

.submission-state {
  justify-self: start;
  padding: 0.27rem 0.65rem;
  color: #697b90;
  background: #f2f5f8;
  border-radius: 999px;
}

.submission-state--accepted {
  color: #14835d;
  background: #e7f8ef;
}

.submission-state--wrong {
  color: #c24141;
  background: #fdeeee;
}

@media (max-width: 767.98px) {
  .submission-table {
    overflow-x: auto;
  }

  .submission-table__row {
    min-width: 46rem;
  }
}
</style>
