<template>
  <div class="competition-problems">
    <div v-if="loading" class="competition-problems__state" role="status">
      <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
      题目加载中...
    </div>

    <div v-else-if="error" class="competition-problems__state competition-problems__state--error" role="alert">
      <i class="bi bi-exclamation-circle" aria-hidden="true"></i>
      {{ error }}
    </div>

    <p v-else-if="problems.length === 0" class="competition-problems__state">暂无可展示题目。</p>

    <div v-else class="problem-table" role="table" aria-label="竞赛题目列表">
      <div class="problem-table__head" role="row">
        <span role="columnheader">题号</span>
        <span role="columnheader">题目名称</span>
        <span role="columnheader">状态</span>
        <span role="columnheader">通过数 / 提交数</span>
      </div>

      <RouterLink
        v-for="problem in problems"
        :key="problem.problem_id || problem.problem_number"
        class="problem-table__row"
        role="row"
        @click.capture="guardProblemNavigation"
        :to="{
          name: 'CompetitionProblem',
          params: {
            competition_id: competitionId,
            problem_number: problem.problem_number
          }
        }"
        :aria-label="`打开题目 ${problem.problem_number} ${problem.problem_title || '无标题'}`"
      >
        <span class="problem-table__number" role="cell">{{ problem.problem_number }}</span>
        <strong class="problem-table__title" role="cell">{{ problem.problem_title || '无标题' }}</strong>
        <span class="problem-status" :class="{ 'problem-status--passed': problem.has_access }" role="cell">
          <i class="bi" :class="problem.has_access ? 'bi-check-circle' : 'bi-circle'" aria-hidden="true"></i>
          {{ problem.has_access ? '已通过' : '未通过' }}
        </span>
        <span class="problem-table__submissions" role="cell">
          <strong>{{ count(problem.ac_count) }}</strong>
          <span>/</span>
          {{ count(problem.submit_count) }}
        </span>
      </RouterLink>
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
  },
  isLogin: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['request-login'])

const problems = ref([])
const loading = ref(false)
const error = ref('')

function count(value) {
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : 0
}

function guardProblemNavigation(event) {
  if (props.isLogin) {
    return
  }

  event.preventDefault()
  emit('request-login')
}

async function loadProblems() {
  const requestedCompetitionId = props.competitionId
  loading.value = true
  error.value = ''

  try {
    const resp = await api.getCompetitionProblems(requestedCompetitionId)
    if (String(requestedCompetitionId) !== String(props.competitionId)) {
      return
    }

    if (resp.code === 0 && Array.isArray(resp.data)) {
      problems.value = resp.data
      return
    }

    problems.value = []
    error.value = resp.message || '题目加载失败，请稍后重试。'
  } catch (requestError) {
    if (String(requestedCompetitionId) === String(props.competitionId)) {
      problems.value = []
      error.value = requestError?.message || '题目加载失败，请稍后重试。'
    }
  } finally {
    if (String(requestedCompetitionId) === String(props.competitionId)) {
      loading.value = false
    }
  }
}

watch(() => props.competitionId, loadProblems, { immediate: true })
</script>

<style scoped>
.competition-problems {
  color: #26354a;
}

.competition-problems__state {
  display: inline-flex;
  gap: 0.6rem;
  align-items: center;
  margin: 0;
  padding: 2.15rem;
  color: #74869b;
  font-size: 0.98rem;
}

.competition-problems__state--error {
  color: #bf5454;
}

.problem-table {
  overflow: hidden;
  border: 0;
  border-radius: 0;
}

.problem-table__head,
.problem-table__row {
  display: grid;
  grid-template-columns: minmax(4rem, 0.65fr) minmax(12rem, 2.1fr) minmax(7.5rem, 1.25fr) minmax(8.75rem, 1.25fr);
  gap: 1rem;
  align-items: center;
  padding: 0 1.55rem;
}

.problem-table__head {
  min-height: 3.6rem;
  color: #55677e;
  font-size: 0.93rem;
  font-weight: 600;
  background: #f8fafc;
  border-bottom: 1px solid #e9eef4;
}

.problem-table__row {
  min-height: 4.8rem;
  color: #26354a;
  text-decoration: none;
  background: #fff;
  border-bottom: 1px solid #edf1f5;
  transition: color 150ms ease, background 150ms ease;
}

.problem-table__row:last-child {
  border-bottom: 0;
}

.problem-table__row:hover {
  color: #145fe5;
  background: #fbfdff;
}

.problem-table__row:focus-visible {
  outline: 3px solid rgba(37, 109, 242, 0.2);
  outline-offset: -3px;
}

.problem-table__number {
  color: #256df2;
  font-size: 1.18rem;
  font-weight: 600;
}

.problem-table__title {
  overflow: hidden;
  font-size: 1rem;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.problem-status {
  display: inline-flex;
  gap: 0.42rem;
  align-items: center;
  justify-self: start;
  min-height: 2rem;
  padding: 0 0.7rem;
  color: #67788e;
  font-size: 0.9rem;
  font-weight: 500;
  background: #f6f8fa;
  border: 1px solid #e4eaf0;
  border-radius: 0.4rem;
}

.problem-status--passed {
  color: #16a05a;
  background: #edfbf2;
  border-color: #cdeed9;
}

.problem-table__submissions {
  color: #607289;
  font-size: 0.98rem;
  font-variant-numeric: tabular-nums;
}

.problem-table__submissions strong {
  color: #26354a;
  font-weight: 600;
}

.problem-table__submissions span {
  margin: 0 0.22rem;
  color: #9aa8b8;
}

@media (max-width: 767.98px) {
  .competition-problems__state {
    padding: 1.15rem;
  }

  .problem-table__head {
    display: none;
  }

  .problem-table__row {
    grid-template-columns: 2.3rem minmax(0, 1fr);
    gap: 0.55rem 0.8rem;
    padding: 0.95rem 1rem;
  }

  .problem-table__number {
    grid-row: span 2;
    align-self: start;
    padding-top: 0.1rem;
  }

  .problem-status {
    grid-column: 2;
  }

  .problem-table__submissions {
    grid-column: 2;
    margin-top: 0.15rem;
    font-size: 0.88rem;
  }

  .problem-table__submissions::before {
    margin-right: 0.45rem;
    color: #8998aa;
    content: "通过/提交";
  }
}
</style>
