<template>
  <section class="competition-rank-board bg-white">

    <div v-if="loading" class="p-4" role="status" aria-label="排行榜加载中">
      <div v-for="index in 5" :key="index" class="placeholder-glow d-flex align-items-center gap-3 py-2">
        <span class="placeholder rounded-circle rank-placeholder"></span>
        <span class="placeholder col-3"></span>
        <span class="placeholder col-2"></span>
        <span class="placeholder col-4"></span>
      </div>
    </div>

    <div v-else-if="error" class="alert alert-danger d-flex align-items-center justify-content-between gap-3 m-4" role="alert">
      <span class="d-inline-flex align-items-center gap-2">
        <i class="bi bi-exclamation-circle" aria-hidden="true"></i>
        {{ error }}
      </span>
      <button type="button" class="btn btn-outline-danger btn-sm flex-shrink-0" @click="loadRankings">
        重新加载
      </button>
    </div>

    <div v-else-if="rankRows.length === 0" class="d-flex flex-column align-items-center justify-content-center gap-2 p-5 text-center text-secondary">
      <i class="bi bi-bar-chart-line fs-2 text-success" aria-hidden="true"></i>
      <strong class="text-dark">暂无排名数据</strong>
    </div>

    <div v-else class="table-responsive competition-rank-table-wrap">
      <table
        class="table table-hover align-middle text-center mb-0 competition-rank-table"
        :style="tableStyle"
      >
        <caption class="visually-hidden">竞赛排行榜</caption>
        <colgroup>
          <col class="rank-col">
          <col class="user-col">
          <col class="score-col">
          <col class="penalty-col">
          <col
            v-for="problem in problemColumns"
            :key="`col-${problem.number}`"
            class="problem-col"
          >
        </colgroup>
        <thead>
          <tr>
            <th scope="col" class="rank-column">Rank</th>
            <th scope="col" class="user-column">选手</th>
            <th scope="col" class="score-column">Score</th>
            <th scope="col" class="penalty-column">Penalty</th>
            <th v-for="problem in problemColumns" :key="problem.number" scope="col" class="problem-column">
              <span class="problem-heading">
                <span class="problem-heading__number">{{ problem.number }}</span>
                <span class="problem-heading__stats">{{ problem.accepted }} / {{ problem.tried }}</span>
              </span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rankingRows" :key="row.row_key">
            <th scope="row" class="fw-semibold">
              <i v-if="rankIcon(row.competition_rank)" class="bi me-1" :class="rankIcon(row.competition_rank)" aria-hidden="true"></i>
              {{ row.competition_rank }}
            </th>
            <td class="user-cell">
              <RouterLink
                v-if="row.user_id"
                class="rank-user-link fw-semibold text-decoration-none"
                :to="{ name: 'UserInfo', params: { userId: row.user_id } }"
              >
                {{ row.user_name }}
              </RouterLink>
              <span v-else class="fw-semibold">{{ row.user_name }}</span>
            </td>
            <td class="fw-semibold text-success">{{ displayNumber(row.solved_count) }}</td>
            <td class="text-secondary">{{ displayNumber(row.penalty) }}</td>
            <td
              v-for="cell in row.cells"
              :key="`${row.row_key}-${cell.problemNumber}`"
              class="problem-cell"
              :class="cell.className"
            >
              <template v-if="cell.hasSubmission">
                <span
                  v-if="cell.acceptedTimeLabel"
                  class="problem-cell__time"
                >
                  {{ cell.acceptedTimeLabel }}
                </span>
                <span
                  v-if="cell.wrongCountLabel"
                  class="problem-cell__fail-count"
                >
                  {{ cell.wrongCountLabel }}
                </span>
              </template>
              <span v-else class="problem-cell__empty" aria-hidden="true"></span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import api from '@/api'
import { formatCountdown } from '@/utils/competitionTime'

const props = defineProps({
  competition: {
    type: Object,
    required: true
  }
})

const loading = ref(false)
const error = ref('')
const rankings = ref([])
const problemNumbers = ref([])
let requestId = 0

const competitionId = computed(() => props.competition?.id)
const rankRows = computed(() => rankings.value
  .map((row, index) => normalizeRankRow(row, index))
  .sort((left, right) => left.competition_rank - right.competition_rank))
const problemStatsByNumber = computed(() => {
  const stats = Object.create(null)

  problemNumbers.value.forEach(number => {
    stats[number] = {
      accepted: 0,
      tried: 0
    }
  })

  rankRows.value.forEach(row => {
    Object.entries(row.submissionMap).forEach(([number, submission]) => {
      const stat = stats[number]
      if (!stat) {
        return
      }

      stat.tried += 1
      if (submission.is_ac) {
        stat.accepted += 1
      }
    })
  })

  return stats
})
const problemColumns = computed(() => problemNumbers.value
  .map(number => {
    const stats = problemStatsByNumber.value[number] || { accepted: 0, tried: 0 }

    return {
      number,
      accepted: stats.accepted,
      tried: stats.tried
    }
  }))
const firstAcceptedTimesByProblem = computed(() => {
  const times = Object.create(null)

  rankRows.value.forEach(row => {
    row.submissions.forEach(submission => {
      const acceptedTime = acceptedTimestamp(submission)
      if (!submission.is_ac || submission.problem_number === '' || !Number.isFinite(acceptedTime)) {
        return
      }

      if (!Number.isFinite(times[submission.problem_number]) || acceptedTime < times[submission.problem_number]) {
        times[submission.problem_number] = acceptedTime
      }
    })
  })

  return times
})
const rankingRows = computed(() => rankRows.value.map(row => ({
  ...row,
  cells: problemColumns.value.map(problem => createProblemCell(row, problem.number))
})))
const tableColumnWidths = Object.freeze({
  rank: 4.4,
  user: 10.5,
  score: 5.2,
  penalty: 6.1,
  problem: 5.6
})
const tableStyle = computed(() => {
  const tableWidth = tableColumnWidths.rank +
    tableColumnWidths.user +
    tableColumnWidths.score +
    tableColumnWidths.penalty +
    (problemColumns.value.length * tableColumnWidths.problem)

  return {
    width: `${tableWidth}rem`,
    '--rank-col-width': `${tableColumnWidths.rank}rem`,
    '--user-col-width': `${tableColumnWidths.user}rem`,
    '--score-col-width': `${tableColumnWidths.score}rem`,
    '--penalty-col-width': `${tableColumnWidths.penalty}rem`,
    '--problem-col-width': `${tableColumnWidths.problem}rem`
  }
})

onMounted(loadRankings)

watch(competitionId, (nextId, previousId) => {
  if (String(nextId || '') !== String(previousId || '')) {
    loadRankings()
  }
})

async function loadRankings() {
  const currentCompetitionId = competitionId.value
  const currentRequestId = ++requestId

  if (!currentCompetitionId) {
    rankings.value = []
    problemNumbers.value = []
    loading.value = false
    error.value = '未获取到竞赛编号'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const resp = await api.getCompetitionRankList(currentCompetitionId)
    if (currentRequestId !== requestId || String(currentCompetitionId) !== String(competitionId.value)) {
      return
    }
    if (resp?.code !== undefined && resp.code !== 0) {
      rankings.value = []
      problemNumbers.value = []
      error.value = resp.message || '榜单加载失败，请稍后重试。'
      return
    }

    const payload = resp?.code === 0 ? resp.data : resp
    const parsed = parseRankingPayload(payload)
    rankings.value = parsed.list
    problemNumbers.value = parsed.problemNumbers
  } catch (requestError) {
    if (currentRequestId === requestId) {
      rankings.value = []
      problemNumbers.value = []
      error.value = requestError?.message || '榜单加载失败，请稍后重试。'
    }
  } finally {
    if (currentRequestId === requestId) {
      loading.value = false
    }
  }
}

function parseRankingPayload(payload) {
  if (Array.isArray(payload)) {
    return {
      list: payload,
      problemNumbers: []
    }
  }

  return {
    list: Array.isArray(payload?.ranking_list) ? payload.ranking_list : [],
    problemNumbers: normalizeProblemNumbers(payload?.problem_numbers)
  }
}

function normalizeProblemNumbers(numbers) {
  return Array.isArray(numbers)
    ? numbers
      .filter(number => number !== undefined && number !== null && number !== '')
      .map(number => String(number))
      .sort(compareProblemNumber)
    : []
}

function normalizeRankRow(row, index) {
  const submissions = Array.isArray(row?.submissions)
    ? row.submissions.map(normalizeSubmission)
    : []
  const rank = Number(row?.competition_rank ?? row?.ranking ?? index + 1)
  const competitionRank = Number.isFinite(rank) ? rank : index + 1
  const userId = row?.user_id ?? ''
  const userName = row?.user_name || '匿名选手'
  const submissionMap = submissions.reduce((map, submission) => {
    if (submission.problem_number !== '') {
      map[submission.problem_number] = submission
    }

    return map
  }, Object.create(null))

  return {
    ...row,
    competition_rank: competitionRank,
    user_id: userId,
    user_name: userName,
    solved_count: row?.solved_count ?? 0,
    penalty: row?.penalty ?? 0,
    submissions,
    submissionMap,
    row_key: userId || `${competitionRank}-${userName}`
  }
}

function normalizeSubmission(submission) {
  return {
    problem_number: submission?.problem_number === undefined || submission?.problem_number === null
      ? ''
      : String(submission.problem_number),
    is_ac: submission?.is_ac === true,
    first_ac_time: submission?.first_ac_time || '',
    wrong_count: normalizeWrongCount(submission)
  }
}

function compareProblemNumber(left, right) {
  return String(left).localeCompare(String(right), 'zh-CN', {
    numeric: true,
    sensitivity: 'base'
  })
}

function displayNumber(value) {
  const number = Number(value)
  return Number.isFinite(number) ? number.toLocaleString('zh-CN') : '-'
}

function rankIcon(rank) {
  return {
    1: 'bi-trophy-fill text-warning',
    2: 'bi-award-fill text-secondary',
    3: 'bi-award-fill rank-third'
  }[rank] || ''
}

function createProblemCell(row, problemNumber) {
  const submission = row.submissionMap[problemNumber]

  return {
    problemNumber,
    hasSubmission: Boolean(submission),
    className: problemCellClass(submission, problemNumber),
    acceptedTimeLabel: acceptedTimeLabel(submission),
    wrongCountLabel: wrongCountLabel(submission)
  }
}

function problemCellClass(submission, problemNumber) {
  if (!submission) {
    return ''
  }

  if (!submission.is_ac) {
    return 'problem-cell--failed'
  }

  return isProblemFirstAC(submission, problemNumber) ? 'problem-cell--first-ac' : 'problem-cell--accepted'
}

function isProblemFirstAC(submission, problemNumber) {
  const firstAcceptedTime = firstAcceptedTimesByProblem.value[problemNumber]
  const acceptedTime = acceptedTimestamp(submission)

  return Number.isFinite(firstAcceptedTime) &&
    Number.isFinite(acceptedTime) &&
    acceptedTime === firstAcceptedTime
}

function acceptedTimestamp(submission) {
  if (!submission?.first_ac_time) {
    return NaN
  }

  return new Date(submission.first_ac_time).getTime()
}

function acceptedTimeLabel(submission) {
  if (!submission?.is_ac || !submission.first_ac_time) {
    return ''
  }

  const acceptedTime = new Date(submission.first_ac_time).getTime()
  const startTime = new Date(props.competition?.start_time).getTime()
  if (Number.isFinite(acceptedTime) && Number.isFinite(startTime) && acceptedTime >= startTime) {
    return formatCountdown(acceptedTime - startTime)
  }

  if (Number.isFinite(acceptedTime)) {
    return new Date(submission.first_ac_time).toLocaleTimeString('zh-CN', { hour12: false })
  }

  return submission.first_ac_time
}

function normalizeWrongCount(submission) {
  const value = submission?.wrong_count ??
    submission?.wrongCount ??
    submission?.failed_count ??
    submission?.failedCount ??
    submission?.fail_count ??
    submission?.failCount ??
    submission?.wrong_attempts ??
    submission?.wrongAttempts
  const number = Number(value)

  return Number.isFinite(number) && number > 0 ? number : 0
}

function wrongCountLabel(submission) {
  return submission?.wrong_count > 0 ? `(-${submission.wrong_count})` : ''
}
</script>

<style scoped>
.competition-rank-board {
  color: #26354a;
}

.competition-rank-table-wrap {
  overflow-x: auto;
}

.competition-rank-table {
  table-layout: fixed;
  border-collapse: collapse;
}

.competition-rank-table th,
.competition-rank-table td {
  overflow: hidden;
  max-width: 0;
  white-space: nowrap;
  text-overflow: ellipsis;
  border: 1px solid #d8e2ec;
}

.competition-rank-table thead th {
  height: 4.45rem;
  padding: 0.55rem 0.65rem;
  color: #52647b;
  font-size: 0.82rem;
  vertical-align: middle;
  background: #f8fafc;
  border-color: #d8e2ec;
}

.competition-rank-table tbody td,
.competition-rank-table tbody th {
  padding: 0.75rem 0.65rem;
  border-color: #d8e2ec;
}

.rank-col {
  width: var(--rank-col-width);
}

.score-col {
  width: var(--score-col-width);
}

.penalty-col {
  width: var(--penalty-col-width);
}

.user-col {
  width: var(--user-col-width);
}

.problem-col {
  width: var(--problem-col-width);
}

.user-column {
  min-width: 0;
  text-align: center;
}

.rank-column,
.score-column,
.penalty-column,
.problem-column {
  text-align: center;
}

.problem-heading {
  display: inline-flex;
  max-width: 100%;
  min-height: 2.5rem;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.18rem;
  line-height: 1.15;
}

.problem-heading__number {
  display: block;
  max-width: 100%;
  overflow: hidden;
  color: #4d5f77;
  font-size: 0.92rem;
  font-weight: 700;
  text-overflow: ellipsis;
}

.problem-heading__stats {
  display: block;
  max-width: 100%;
  overflow: hidden;
  color: #697c92;
  font-size: 0.76rem;
  font-weight: 500;
  text-overflow: ellipsis;
}

.problem-cell {
  height: 4.35rem;
  color: #fff;
  font-weight: 600;
  vertical-align: middle;
}

.problem-cell__empty {
  display: block;
  min-height: 1.7rem;
}

.problem-cell--first-ac {
  --bs-table-bg: #055932;
  --bs-table-hover-bg: #055932;
  background: #055932;
}

.problem-cell--accepted {
  --bs-table-bg: #29863e;
  --bs-table-hover-bg: #29863e;
  background: #29863e;
}

.problem-cell--failed {
  --bs-table-bg: #a72837;
  --bs-table-hover-bg: #a72837;
  background: #a72837;
}

.problem-cell__time,
.problem-cell__fail-count {
  display: block;
  overflow: hidden;
  max-width: 100%;
  color: #fff;
  font-size: 0.78rem;
  line-height: 1.25;
  font-variant-numeric: tabular-nums;
  text-overflow: ellipsis;
}

.problem-cell__fail-count {
  margin-top: 0.16rem;
}

.user-cell {
  overflow: hidden;
  text-overflow: ellipsis;
}

.rank-user-link {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  color: #24364c;
  text-overflow: ellipsis;
  vertical-align: middle;
}

.rank-user-link:hover {
  color: #13866c;
}

.rank-placeholder {
  width: 2rem;
  height: 2rem;
}

.rank-third {
  color: #b86b34;
}

</style>
