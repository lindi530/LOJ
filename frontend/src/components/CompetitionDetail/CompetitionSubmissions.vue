<template>
  <section class="competition-submissions">
    <div v-if="loading && records.length === 0" class="d-inline-flex align-items-center gap-2 m-0 p-4 text-secondary" role="status">
      <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
      提交记录加载中...
    </div>

    <div v-else-if="error" class="alert alert-danger d-inline-flex align-items-center gap-2 m-4 py-2 mb-0" role="alert">
      <i class="bi bi-exclamation-circle" aria-hidden="true"></i>
      {{ error }}
    </div>

    <template v-else>
      <div class="table-responsive submission-table d-none d-lg-block">
        <table class="table table-hover align-middle text-center mb-0">
          <thead>
            <tr>
              <th scope="col" class="text-nowrap">提交</th>
              <th scope="col" class="text-nowrap">用户</th>
              <th scope="col">
                <div class="submission-problem-heading">
                  <span>题目</span>
                  <details v-if="showProblemFilter" class="problem-number-filter">
                    <summary class="problem-number-filter__toggle" :class="{ 'problem-number-filter__toggle--active': problemFilter }">
                      <i class="bi bi-chevron-down" aria-hidden="true"></i>
                    </summary>
                    <div class="problem-number-filter__menu">
                      <button
                        type="button"
                        class="problem-number-filter__option"
                        :class="{ 'problem-number-filter__option--selected': problemFilter === '' }"
                        @click="selectProblemFilter('', $event)"
                      >
                        全部题目
                      </button>
                      <button
                        v-for="problemNumber in problemNumbers"
                        :key="problemNumber"
                        type="button"
                        class="problem-number-filter__option"
                        :class="{ 'problem-number-filter__option--selected': problemFilter === problemNumber }"
                        @click="selectProblemFilter(problemNumber, $event)"
                      >
                        {{ problemNumber }}
                      </button>
                    </div>
                  </details>
                </div>
              </th>
              <th scope="col" class="text-nowrap">结果</th>
              <th scope="col" class="text-nowrap">语言</th>
              <th scope="col" class="text-nowrap">耗时 / 内存</th>
              <th scope="col" class="text-nowrap">提交时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="records.length === 0">
              <td colspan="7">
                <div class="submission-empty d-flex flex-column align-items-center justify-content-center gap-1 p-4 text-center text-secondary">
                  <i class="bi bi-inboxes" aria-hidden="true"></i>
                  <strong class="text-dark">{{ problemFilter ? '当前题号暂无提交记录' : '暂无提交记录' }}</strong>
                  <span>{{ problemFilter ? '可以切换题号查看其他提交。' : '还没有任何判题结果。' }}</span>
                </div>
              </td>
            </tr>
            <tr v-for="record in records" :key="recordKey(record)">
              <td class="text-nowrap">
                <RouterLink
                  v-if="recordId(record) !== '-'"
                  class="submission-link"
                  :to="{ name: 'SubmissionDetail', params: { submission_id: recordId(record) } }"
                >
                  #{{ recordId(record) }}
                </RouterLink>
                <span v-else>-</span>
              </td>
              <td class="text-nowrap">
                <RouterLink
                  v-if="userId(record) !== '-' && userLabel(record) !== '-'"
                  class="submission-user-link"
                  :to="{ name: 'UserInfo', params: { userId: userId(record) } }"
                >
                  {{ userLabel(record) }}
                </RouterLink>
                <span v-else class="submission-user">
                  {{ userLabel(record) }}
                </span>
              </td>
              <td class="submission-table__problem">
                <RouterLink
                  v-if="competitionProblemRoute(record)"
                  class="submission-problem-link"
                  :to="competitionProblemRoute(record)"
                >
                  <strong class="submission-problem-title">{{ problemLabel(record) }}</strong>
                </RouterLink>
                <template v-else>
                  <strong class="submission-problem-title">{{ problemLabel(record) }}</strong>
                </template>
              </td>
              <td class="text-nowrap">
                <RouterLink
                  v-if="recordId(record) !== '-'"
                  class="submission-result submission-result-link"
                  :class="stateClass(resultLabel(record))"
                  :to="{ name: 'SubmissionDetail', params: { submission_id: recordId(record) } }"
                >
                  {{ resultLabel(record) }}
                </RouterLink>
                <span v-else class="submission-result" :class="stateClass(resultLabel(record))">
                  {{ resultLabel(record) }}
                </span>
              </td>
              <td class="text-nowrap">
                <span class="submission-language">{{ languageLabel(record) }}</span>
              </td>
              <td class="text-nowrap">
                <span class="submission-runtime">{{ durationLabel(record) }}</span>
                <span class="submission-divider">/</span>
                <span class="submission-memory">{{ memoryLabel(record) }}</span>
              </td>
              <td class="text-nowrap submission-time">{{ submittedAt(record) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="records.length > 0" class="d-grid gap-3 d-lg-none bg-body-tertiary p-3">
        <article v-for="record in records" :key="`card-${recordKey(record)}`" class="bg-white border rounded-2 shadow-sm p-3">
          <div class="d-flex align-items-start justify-content-between gap-3">
            <div>
              <RouterLink
                v-if="recordId(record) !== '-'"
                class="submission-link"
                :to="{ name: 'SubmissionDetail', params: { submission_id: recordId(record) } }"
              >
                #{{ recordId(record) }}
              </RouterLink>
              <span v-else class="fw-semibold">-</span>
              <div class="submission-card__user mt-1">
                <RouterLink
                  v-if="userId(record) !== '-' && userLabel(record) !== '-'"
                  class="submission-user-link"
                  :to="{ name: 'UserInfo', params: { userId: userId(record) } }"
                >
                  {{ userLabel(record) }}
                </RouterLink>
                <span v-else>{{ userLabel(record) }}</span>
              </div>
            </div>

            <RouterLink
              v-if="recordId(record) !== '-'"
              class="submission-result submission-result-link"
              :class="stateClass(resultLabel(record))"
              :to="{ name: 'SubmissionDetail', params: { submission_id: recordId(record) } }"
            >
              {{ resultLabel(record) }}
            </RouterLink>
            <span v-else class="submission-result" :class="stateClass(resultLabel(record))">
              {{ resultLabel(record) }}
            </span>
          </div>

          <div class="mt-3 border-top pt-3">
            <RouterLink
              v-if="competitionProblemRoute(record)"
              class="submission-problem-link"
              :to="competitionProblemRoute(record)"
            >
              <strong class="submission-problem-title">{{ problemLabel(record) }}</strong>
            </RouterLink>
            <template v-else>
              <strong class="submission-problem-title">{{ problemLabel(record) }}</strong>
            </template>
          </div>

          <div class="row row-cols-2 g-2 mt-3">
            <div class="col">
              <span class="d-block small text-secondary">语言</span>
              <strong class="submission-meta-value">{{ languageLabel(record) }}</strong>
            </div>
            <div class="col">
              <span class="d-block small text-secondary">耗时</span>
              <strong class="submission-meta-value">{{ durationLabel(record) }}</strong>
            </div>
            <div class="col">
              <span class="d-block small text-secondary">内存</span>
              <strong class="submission-meta-value">{{ memoryLabel(record) }}</strong>
            </div>
            <div class="col">
              <span class="d-block small text-secondary">时间</span>
              <strong class="submission-meta-value">{{ submittedAt(record) }}</strong>
            </div>
          </div>
        </article>
      </div>

      <div v-else class="submission-empty d-flex d-lg-none flex-column align-items-center justify-content-center gap-1 p-4 text-center text-secondary">
        <i class="bi bi-inboxes" aria-hidden="true"></i>
        <strong class="text-dark">{{ problemFilter ? '当前题号暂无提交记录' : '暂无提交记录' }}</strong>
        <span>{{ problemFilter ? '可以切换题号查看其他提交。' : '还没有任何判题结果。' }}</span>
      </div>

      <footer class="d-flex flex-column flex-sm-row align-items-sm-center justify-content-between gap-3 border-top bg-white px-4 py-3">
        <p class="mb-0 small text-secondary">显示 {{ shownStart }}-{{ shownEnd }} / {{ total }}</p>

        <nav v-if="totalPages > 1" aria-label="提交记录分页">
          <ul class="pagination pagination-sm mb-0">
            <li class="page-item" :class="{ disabled: currentPage === 1 || loading }">
              <button class="page-link d-inline-flex align-items-center gap-1" type="button" :disabled="currentPage === 1 || loading" @click="goToPage(currentPage - 1)">
                <i class="bi bi-chevron-left" aria-hidden="true"></i>
                上一页
              </button>
            </li>
            <li
              v-for="page in visiblePages"
              :key="page"
              class="page-item"
              :class="{ active: page === currentPage, disabled: loading }"
            >
              <button
                class="page-link"
                type="button"
                :disabled="loading"
                :aria-current="page === currentPage ? 'page' : undefined"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
            </li>
            <li class="page-item" :class="{ disabled: currentPage === totalPages || loading }">
              <button class="page-link d-inline-flex align-items-center gap-1" type="button" :disabled="currentPage === totalPages || loading" @click="goToPage(currentPage + 1)">
                下一页
                <i class="bi bi-chevron-right" aria-hidden="true"></i>
              </button>
            </li>
          </ul>
        </nav>
      </footer>
    </template>
  </section>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import api from '@/api'

const props = defineProps({
  competitionId: {
    type: [Number, String],
    required: true
  }
})

const pageSize = 20
const currentPage = ref(1)
const loading = ref(false)
const error = ref('')
const records = ref([])
const total = ref(0)
const problemFilter = ref('')
const problemNumbers = ref([])

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const shownStart = computed(() => records.value.length === 0 ? 0 : (currentPage.value - 1) * pageSize + 1)
const shownEnd = computed(() => Math.min((currentPage.value - 1) * pageSize + records.value.length, total.value))
const visiblePages = computed(() => {
  const count = totalPages.value
  const start = Math.max(1, Math.min(currentPage.value - 2, count - 4))
  const end = Math.min(count, start + 4)

  return Array.from({ length: end - start + 1 }, (_, index) => start + index)
})
const showProblemFilter = computed(() => problemNumbers.value.length > 0 || problemFilter.value !== '')
function value(input) {
  return input === undefined || input === null || input === '' ? '-' : input
}

function firstValue(record, keys) {
  const key = keys.find(item => record[item] !== undefined && record[item] !== null && record[item] !== '')
  return key ? record[key] : ''
}

function recordId(record) {
  return value(firstValue(record, ['id', 'submission_id', 'run_id']))
}

function recordKey(record) {
  return recordId(record) === '-' ? `${problemLabel(record)}-${formatTime(firstValue(record, ['submission_time', 'created_at']))}` : recordId(record)
}

function userLabel(record) {
  return value(firstValue(record, ['submission_user_name', 'user_name', 'username']))
}

function userId(record) {
  return value(firstValue(record, ['submission_user_id', 'user_id']))
}

function problemNumber(record) {
  return firstValue(record, ['problem_number', 'problem_index'])
}

function competitionProblemRoute(record) {
  const number = problemNumber(record)

  if (number === '') {
    return null
  }

  return {
    name: 'CompetitionProblem',
    params: {
      competition_id: props.competitionId,
      problem_number: number
    }
  }
}

function problemTitle(record) {
  return value(firstValue(record, ['problem_title', 'title']) || record.problem_id)
}

function problemLabel(record) {
  const number = problemNumber(record)
  const title = problemTitle(record)

  if (number && title !== '-') {
    return `${number}-${title}`
  }

  return value(number || title)
}

function resultLabel(record) {
  return value(firstValue(record, ['state', 'status', 'result'])) === '-' ? '评测中' : value(firstValue(record, ['state', 'status', 'result']))
}

function resultKind(state = '') {
  const text = String(state)
  const normalized = text.toLowerCase()

  if (normalized === 'ac' || normalized.includes('accepted') || text === '通过' || text === '答案正确') {
    return 'accepted'
  }

  if (
    normalized.includes('time limit') ||
    normalized.includes('memory limit') ||
    normalized === 'tle' ||
    normalized === 'mle' ||
    text.includes('超时') ||
    text.includes('超限')
  ) {
    return 'limited'
  }

  if (
    normalized.includes('wrong') ||
    normalized.includes('error') ||
    normalized.includes('rejected') ||
    normalized.includes('failed') ||
    text.includes('未通过') ||
    text.includes('答案错误') ||
    text.includes('编译错误') ||
    text.includes('运行错误')
  ) {
    return 'failed'
  }

  return 'judging'
}

function stateClass(state = '') {
  const kind = resultKind(state)

  if (kind === 'accepted') {
    return 'text-success'
  }

  if (kind === 'limited') {
    return 'text-warning'
  }

  if (kind === 'failed') {
    return 'text-danger'
  }

  return 'text-secondary'
}

function languageLabel(record) {
  return value(record.language || record.lang)
}

function durationLabel(record) {
  const duration = firstValue(record, ['exec_limit', 'exec_time', 'run_time', 'runtime', 'time'])
  return typeof duration === 'number' ? `${duration} ms` : value(duration)
}

function memoryLabel(record) {
  const memory = firstValue(record, ['memory_limit', 'memory_usage', 'memory', 'memory_used'])
  return typeof memory === 'number' ? `${memory} KB` : value(memory)
}

function submittedAt(record) {
  return formatTime(firstValue(record, ['submission_time', 'created_at', 'submit_time']))
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

function parseSubmissionPage(data) {
  if (Array.isArray(data)) {
    return {
      list: data,
      count: data.length,
      problemNumbers: []
    }
  }

  const list = data?.list || data?.records || data?.submissions || []
  if (!Array.isArray(list) && list && typeof list === 'object') {
    return {
      list: [list],
      count: 1,
      problemNumbers: normalizeProblemNumbers(data?.problem_numbers)
    }
  }

  if (data && typeof data === 'object' && 'submission_id' in data) {
    return {
      list: [data],
      count: 1,
      problemNumbers: normalizeProblemNumbers(data?.problem_numbers)
    }
  }

  return {
    list: Array.isArray(list) ? list : [],
    count: Number(data?.count ?? data?.total ?? data?.total_count ?? list.length) || 0,
    problemNumbers: normalizeProblemNumbers(data?.problem_numbers)
  }
}

function normalizeProblemNumbers(numbers) {
  return Array.isArray(numbers)
    ? numbers
      .filter(number => number !== undefined && number !== null && number !== '')
      .map(number => String(number))
    : []
}

function goToPage(page) {
  if (loading.value || page < 1 || page > totalPages.value || page === currentPage.value) {
    return
  }

  currentPage.value = page
}

function selectProblemFilter(problemNumber, event) {
  problemFilter.value = problemNumber
  event?.currentTarget?.closest('details')?.removeAttribute('open')
}

function resetAndLoadRecords() {
  if (currentPage.value === 1) {
    loadRecords()
    return
  }

  currentPage.value = 1
}

async function loadRecords() {
  const requestedCompetitionId = props.competitionId
  const requestedPage = currentPage.value
  const requestedProblemNumber = problemFilter.value
  loading.value = true
  error.value = ''

  try {
    const resp = await api.getCompetitionSubmissions(requestedCompetitionId, requestedPage, pageSize, requestedProblemNumber)
    if (
      String(requestedCompetitionId) !== String(props.competitionId) ||
      requestedPage !== currentPage.value ||
      requestedProblemNumber !== problemFilter.value
    ) {
      return
    }

    if (resp.code === 0) {
      const pageData = parseSubmissionPage(resp.data)
      records.value = pageData.list
      total.value = pageData.count
      problemNumbers.value = pageData.problemNumbers
      return
    }

    records.value = []
    total.value = 0
    problemNumbers.value = []
    error.value = resp.message || '提交记录加载失败，请稍后重试。'
  } catch (requestError) {
    if (
      String(requestedCompetitionId) === String(props.competitionId) &&
      requestedPage === currentPage.value &&
      requestedProblemNumber === problemFilter.value
    ) {
      records.value = []
      total.value = 0
      problemNumbers.value = []
      error.value = requestError?.message || '提交记录加载失败，请稍后重试。'
    }
  } finally {
    if (
      String(requestedCompetitionId) === String(props.competitionId) &&
      requestedPage === currentPage.value &&
      requestedProblemNumber === problemFilter.value
    ) {
      loading.value = false
    }
  }
}

watch(() => props.competitionId, () => {
  problemFilter.value = ''
  problemNumbers.value = []
  resetAndLoadRecords()
}, { immediate: true })

watch(problemFilter, resetAndLoadRecords)
watch(currentPage, loadRecords)
</script>

<style scoped>
.competition-submissions {
  color: #26354a;
  min-height: 11rem;
  background: linear-gradient(180deg, #fbfdfc 0%, #fff 8.5rem);
}

.submission-empty i {
  color: #22ae90;
  font-size: 2.1rem;
}

.submission-table th {
  padding: 0.95rem 1rem;
  color: #60728a;
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
  background: #f7fafc;
  border-bottom: 1px solid #e7eef3;
}

.submission-problem-heading {
  display: inline-flex;
  gap: 0.7rem;
  align-items: center;
  justify-content: center;
}

.problem-number-filter {
  position: relative;
  display: inline-flex;
}

.problem-number-filter__toggle {
  display: inline-flex;
  gap: 0.3rem;
  align-items: center;
  min-height: 1.95rem;
  padding: 0 0.55rem;
  color: #60728a;
  font-size: 0.78rem;
  font-weight: 700;
  list-style: none;
  cursor: pointer;
  background: #fff;
  border: 1px solid #d9e3eb;
  border-radius: 0.38rem;
}

.problem-number-filter__toggle::-webkit-details-marker {
  display: none;
}

.problem-number-filter__toggle--active {
  color: #13866c;
  border-color: #9edfd0;
}

.problem-number-filter[open] .problem-number-filter__toggle {
  color: #13866c;
  border-color: #22ae90;
  box-shadow: 0 0 0 3px rgba(34, 174, 144, 0.14);
}

.problem-number-filter[open] .problem-number-filter__toggle i {
  transform: rotate(180deg);
}

.problem-number-filter__toggle i {
  font-size: 0.7rem;
  transition: transform 140ms ease;
}

.problem-number-filter__menu {
  position: absolute;
  top: calc(100% + 0.35rem);
  left: 50%;
  z-index: 10;
  min-width: 7.2rem;
  overflow: hidden;
  padding: 0.28rem;
  background: #fff;
  border: 1px solid #dce6ee;
  border-radius: 0.45rem;
  box-shadow: 0 0.65rem 1.4rem rgba(38, 53, 74, 0.14);
  transform: translateX(-50%);
}

.problem-number-filter__option {
  display: block;
  width: 100%;
  min-height: 2rem;
  padding: 0 0.65rem;
  color: #43546b;
  font: inherit;
  font-size: 0.86rem;
  font-weight: 600;
  text-align: left;
  white-space: nowrap;
  background: transparent;
  border: 0;
  border-radius: 0.32rem;
}

.problem-number-filter__option:hover,
.problem-number-filter__option:focus {
  color: #13866c;
  background: #f1faf7;
  outline: none;
}

.problem-number-filter__option--selected {
  color: #087b61;
  background: #e7f8f3;
}

.submission-table td {
  padding: 1rem;
  color: #2b3b51;
  font-size: 0.92rem;
  border-color: #eef3f7;
}

.submission-table__problem {
  min-width: 8rem;
}

.submission-table__problem strong,
.submission-problem-title {
  display: block;
  color: #24364c;
  font-weight: 650;
}

.submission-problem-link {
  display: inline-block;
  text-decoration: none;
}

.submission-problem-link:hover .submission-problem-title {
  color: #174fb9;
  text-decoration: underline;
}

.submission-link {
  color: #256df2;
  font-weight: 700;
  text-decoration: none;
  font-variant-numeric: tabular-nums;
}

.submission-link:hover {
  color: #174fb9;
}

.submission-user,
.submission-card__user,
.submission-user-link {
  color: #52647b;
  font-weight: 600;
}

.submission-user-link,
.submission-result-link {
  text-decoration: none;
}

.submission-user-link:hover,
.submission-result-link:hover {
  text-decoration: underline;
}

.submission-result {
  min-width: 5.2rem;
  font-weight: 500;
}

.submission-language {
  color: #52647b;
  font-weight: 500;
}

.submission-runtime,
.submission-memory,
.submission-time {
  color: #52647b;
  font-variant-numeric: tabular-nums;
}

.submission-divider {
  margin: 0 0.34rem;
  color: #a2b0bf;
}

.submission-meta-value {
  display: block;
  overflow: hidden;
  margin-top: 0.12rem;
  color: #2d3d53;
  font-size: 0.9rem;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pagination {
  --bs-pagination-color: #13866c;
  --bs-pagination-active-bg: #22ae90;
  --bs-pagination-active-border-color: #22ae90;
  --bs-pagination-hover-color: #0f6e58;
}
</style>
