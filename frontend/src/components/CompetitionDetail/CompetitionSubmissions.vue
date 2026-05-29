<template>
  <div class="competition-submissions">
    <p v-if="loading && records.length === 0" class="d-inline-flex align-items-center gap-2 mb-0 text-secondary" role="status">
      <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
      提交记录加载中...
    </p>

    <div v-else-if="error" class="alert alert-danger d-inline-flex align-items-center gap-2 mb-0 py-2" role="alert">
      <i class="bi bi-exclamation-circle" aria-hidden="true"></i>
      {{ error }}
    </div>

    <p v-else-if="records.length === 0" class="mb-0 text-secondary">暂无提交记录。</p>

    <template v-else>
      <div class="table-responsive submission-table border rounded-3">
        <table class="table table-sm table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th scope="col" class="text-nowrap">提交编号</th>
              <th scope="col" class="text-nowrap">提交者</th>
              <th scope="col">题目</th>
              <th scope="col" class="text-nowrap">状态</th>
              <th scope="col" class="text-nowrap">语言</th>
              <th scope="col" class="text-nowrap">时间限制</th>
              <th scope="col" class="text-nowrap">内存限制</th>
              <th scope="col" class="text-nowrap">提交时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in records" :key="recordKey(record)">
              <td class="text-nowrap fw-semibold">
                <RouterLink
                  v-if="recordId(record) !== '-'"
                  class="link-primary text-decoration-none"
                  :to="{ name: 'SubmissionDetail', params: { submission_id: recordId(record) } }"
                >
                  #{{ recordId(record) }}
                </RouterLink>
                <span v-else>-</span>
              </td>
              <td class="text-nowrap">{{ userLabel(record) }}</td>
              <td class="submission-table__problem">{{ problemLabel(record) }}</td>
              <td class="text-nowrap">
                <span class="badge rounded-pill" :class="stateClass(resultLabel(record))">
                  {{ resultLabel(record) }}
                </span>
              </td>
              <td class="text-nowrap">{{ record.language || record.lang || '-' }}</td>
              <td class="text-nowrap">{{ durationLabel(record) }}</td>
              <td class="text-nowrap">{{ memoryLabel(record) }}</td>
              <td class="text-nowrap">{{ formatTime(firstValue(record, ['submission_time', 'created_at', 'submit_time'])) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="d-flex flex-column flex-sm-row align-items-sm-center justify-content-between gap-3 mt-3">
        <p class="mb-0 small text-secondary">
          第 {{ currentPage }} 页，显示 {{ shownStart }}-{{ shownEnd }} / {{ total }}
          <span v-if="loading" class="ms-2">
            <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
          </span>
        </p>

        <nav v-if="totalPages > 1" aria-label="提交记录分页">
          <ul class="pagination pagination-sm mb-0">
            <li class="page-item" :class="{ disabled: currentPage === 1 || loading }">
              <button class="page-link" type="button" @click="goToPage(currentPage - 1)">上一页</button>
            </li>
            <li
              v-for="page in visiblePages"
              :key="page"
              class="page-item"
              :class="{ active: page === currentPage, disabled: loading }"
            >
              <button class="page-link" type="button" @click="goToPage(page)">{{ page }}</button>
            </li>
            <li class="page-item" :class="{ disabled: currentPage === totalPages || loading }">
              <button class="page-link" type="button" @click="goToPage(currentPage + 1)">下一页</button>
            </li>
          </ul>
        </nav>
      </div>
    </template>
  </div>
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

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const shownStart = computed(() => records.value.length === 0 ? 0 : (currentPage.value - 1) * pageSize + 1)
const shownEnd = computed(() => Math.min((currentPage.value - 1) * pageSize + records.value.length, total.value))
const visiblePages = computed(() => {
  const count = totalPages.value
  const start = Math.max(1, Math.min(currentPage.value - 2, count - 4))
  const end = Math.min(count, start + 4)

  return Array.from({ length: end - start + 1 }, (_, index) => start + index)
})

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
  return value(firstValue(record, ['user_name', 'username', 'nickname', 'submission_user_id', 'user_id']))
}

function problemLabel(record) {
  const number = firstValue(record, ['problem_number', 'problem_index'])
  const title = firstValue(record, ['problem_title', 'title'])

  if (number && title) {
    return `${number} ${title}`
  }

  return value(number || title || record.problem_id)
}

function resultLabel(record) {
  return value(firstValue(record, ['state', 'status', 'result'])) === '-' ? '评测中' : value(firstValue(record, ['state', 'status', 'result']))
}

function stateClass(state = '') {
  const text = String(state)
  const normalized = text.toLowerCase()

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
    return 'text-bg-danger'
  }

  if (
    normalized.includes('time') ||
    normalized.includes('memory') ||
    normalized.includes('limit') ||
    text.includes('超时') ||
    text.includes('超限')
  ) {
    return 'text-bg-warning'
  }

  if (normalized === 'ac' || normalized.includes('accepted') || text === '通过' || text === '答案正确') {
    return 'text-bg-success'
  }

  return 'text-bg-secondary'
}

function durationLabel(record) {
  const duration = firstValue(record, ['exec_limit', 'exec_time', 'run_time', 'runtime', 'time'])
  return typeof duration === 'number' ? `${duration} ms` : value(duration)
}

function memoryLabel(record) {
  const memory = firstValue(record, ['memory_limit', 'memory_usage', 'memory', 'memory_used'])
  return typeof memory === 'number' ? `${memory} KB` : value(memory)
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
      count: data.length
    }
  }

  const list = data?.list || data?.records || data?.submissions || []
  if (!Array.isArray(list) && list && typeof list === 'object') {
    return {
      list: [list],
      count: 1
    }
  }

  if (data && typeof data === 'object' && 'submission_id' in data) {
    return {
      list: [data],
      count: 1
    }
  }

  return {
    list: Array.isArray(list) ? list : [],
    count: Number(data?.total ?? data?.total_count ?? list.length) || 0
  }
}

function goToPage(page) {
  if (loading.value || page < 1 || page > totalPages.value || page === currentPage.value) {
    return
  }

  currentPage.value = page
}

async function loadRecords() {
  const requestedCompetitionId = props.competitionId
  const requestedPage = currentPage.value
  loading.value = true
  error.value = ''

  try {
    const resp = await api.getCompetitionSubmissions(requestedCompetitionId, requestedPage, pageSize)
    if (String(requestedCompetitionId) !== String(props.competitionId) || requestedPage !== currentPage.value) {
      return
    }

    if (resp.code === 0) {
      const pageData = parseSubmissionPage(resp.data)
      console.log('Parsed submission page data:', pageData)
      records.value = pageData.list
      total.value = pageData.count
      return
    }

    records.value = []
    total.value = 0
    error.value = resp.message || '提交记录加载失败，请稍后重试。'
  } catch (requestError) {
    if (String(requestedCompetitionId) === String(props.competitionId) && requestedPage === currentPage.value) {
      records.value = []
      total.value = 0
      error.value = requestError?.message || '提交记录加载失败，请稍后重试。'
    }
  } finally {
    if (String(requestedCompetitionId) === String(props.competitionId) && requestedPage === currentPage.value) {
      loading.value = false
    }
  }
}

watch(() => props.competitionId, () => {
  if (currentPage.value === 1) {
    loadRecords()
    return
  }

  currentPage.value = 1
}, { immediate: true })

watch(currentPage, loadRecords)
</script>

<style scoped>
.competition-submissions {
  color: #26354a;
}

.submission-table {
  background: #fff;
}

.submission-table :deep(th) {
  color: #55677e;
  font-size: 0.9rem;
  font-weight: 600;
}

.submission-table :deep(td) {
  color: #26354a;
  font-size: 0.92rem;
}

.submission-table__problem {
  min-width: 8rem;
  font-weight: 600;
}

.badge {
  min-width: 4.5rem;
  font-weight: 500;
}

.pagination {
  --bs-pagination-color: #13866c;
  --bs-pagination-active-bg: #22ae90;
  --bs-pagination-active-border-color: #22ae90;
  --bs-pagination-hover-color: #0f6e58;
}
</style>
