<template>
  <section>
    <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3">
      <div>
        <h2 class="h5 fw-bold mb-1">{{ title }}</h2>
        <p v-if="description" class="text-secondary small mb-0">{{ description }}</p>
      </div>
      <span class="badge text-bg-light border">已选 {{ selectedProblems.length }} 题</span>
    </div>

    <label class="form-label fw-medium" :for="searchInputId">搜索并添加题目</label>
    <div class="mb-4">
      <div class="input-group">
        <select
          v-model="problemSearchMode"
          class="form-select flex-grow-0 w-auto"
          aria-label="选择题目搜索方式"
        >
          <option value="id">按题号</option>
          <option value="title">按题目名称</option>
        </select>
        <input
          :id="searchInputId"
          v-model="problemKeyword"
          type="search"
          class="form-control"
          :inputmode="problemSearchMode === 'id' ? 'numeric' : 'search'"
          :placeholder="problemSearchMode === 'id' ? '输入题号' : '输入题目名称'"
          autocomplete="off"
          @keydown.esc="problemKeyword = ''"
        >
      </div>

      <div v-if="problemKeyword.trim()" class="list-group mt-1 shadow-sm problem-search-results">
        <div v-if="problemsLoading" class="list-group-item text-secondary">
          <span class="spinner-border spinner-border-sm me-2" aria-hidden="true"></span>
          搜索中
        </div>
        <button
          v-else-if="problemsError"
          type="button"
          class="list-group-item list-group-item-action text-start text-secondary"
          @click="loadProblems"
        >
          {{ problemsError }}，点击重试
        </button>
        <div v-else-if="searchResults.length === 0" class="list-group-item text-secondary">
          没有找到可添加的题目。
        </div>
        <template v-else>
          <button
            v-for="problem in searchResults"
            :key="problem.id"
            type="button"
            class="list-group-item list-group-item-action d-flex align-items-center gap-3 text-start"
            @click="addProblem(problem)"
          >
            <span class="text-secondary">#{{ problem.id }}</span>
            <span class="fw-medium flex-grow-1">{{ problem.title }}</span>
            <span v-if="problem.level" class="badge text-bg-light fw-normal">{{ problem.level }}</span>
          </button>
        </template>
      </div>
    </div>

    <div v-if="selectedProblems.length === 0" class="alert alert-light border text-secondary mb-0">
      尚未选择题目，请使用上方搜索添加题目。
    </div>

    <div v-else class="list-group">
      <div
        v-for="(problem, index) in selectedProblems"
        :key="problem.id"
        class="list-group-item d-flex flex-wrap align-items-center gap-3"
      >
        <span class="badge rounded-pill text-bg-primary problem-order-badge">
          {{ index + 1 }}
        </span>
        <span class="text-secondary">#{{ problem.id }}</span>
        <span class="fw-medium flex-grow-1">{{ problem.title }}</span>
        <span v-if="problem.level" class="badge text-bg-light fw-normal">{{ problem.level }}</span>
        <div class="btn-group btn-group-sm" role="group" :aria-label="`${problem.title} 排序`">
          <button
            type="button"
            class="btn btn-outline-secondary"
            :disabled="index === 0"
            :title="`上移 ${problem.title}`"
            :aria-label="`上移 ${problem.title}`"
            @click="moveProblem(index, -1)"
          >
            <i class="bi bi-arrow-up" aria-hidden="true"></i>
          </button>
          <button
            type="button"
            class="btn btn-outline-secondary"
            :disabled="index === selectedProblems.length - 1"
            :title="`下移 ${problem.title}`"
            :aria-label="`下移 ${problem.title}`"
            @click="moveProblem(index, 1)"
          >
            <i class="bi bi-arrow-down" aria-hidden="true"></i>
          </button>
        </div>
        <button
          type="button"
          class="btn btn-outline-secondary btn-sm"
          :title="`移除 ${problem.title}`"
          :aria-label="`移除 ${problem.title}`"
          @click="removeProblem(problem.id)"
        >
          <i class="bi bi-x-lg" aria-hidden="true"></i>
        </button>
      </div>
    </div>

    <p v-if="submitted && required && selectedProblems.length === 0" class="text-danger small mt-2 mb-0">
      请至少选择一道题目。
    </p>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '@/api'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  title: {
    type: String,
    default: '选择题目'
  },
  description: {
    type: String,
    default: ''
  },
  submitted: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const problems = ref([])
const problemsLoading = ref(false)
const problemsError = ref('')
const problemKeyword = ref('')
const problemSearchMode = ref('id')
const searchInputId = `problem-search-${Math.random().toString(36).slice(2, 8)}`

const selectedProblems = computed(() => props.modelValue)
const selectedProblemIds = computed(() => selectedProblems.value.map((problem) => Number(problem.id)))

const searchResults = computed(() => {
  const keyword = problemKeyword.value.trim()

  if (!keyword) {
    return []
  }

  return problems.value.filter((problem) => {
    if (selectedProblemIds.value.includes(problem.id)) {
      return false
    }

    if (problemSearchMode.value === 'id') {
      return String(problem.id).startsWith(keyword)
    }

    return problem.title.toLowerCase().includes(keyword.toLowerCase())
  })
})

function normalizeProblem(problem) {
  return {
    id: Number(problem.id ?? problem.problem_id),
    title: problem.title || problem.name || '无标题',
    level: problem.level || problem.difficulty || ''
  }
}

function updateSelectedProblems(nextProblems) {
  emit('update:modelValue', nextProblems)
}

function addProblem(problem) {
  updateSelectedProblems([...selectedProblems.value, problem])
  problemKeyword.value = ''
}

function removeProblem(problemId) {
  updateSelectedProblems(selectedProblems.value.filter((problem) => problem.id !== problemId))
}

function moveProblem(index, direction) {
  const destination = index + direction

  if (destination < 0 || destination >= selectedProblems.value.length) {
    return
  }

  const nextProblems = [...selectedProblems.value]
  const currentProblem = nextProblems[index]
  nextProblems[index] = nextProblems[destination]
  nextProblems[destination] = currentProblem
  updateSelectedProblems(nextProblems)
}

async function loadProblems() {
  problemsLoading.value = true
  problemsError.value = ''

  try {
    const resp = await api.getProblemList()
    const data = resp?.data
    const list = Array.isArray(data) ? data : data?.list || []

    problems.value = list
      .map(normalizeProblem)
      .filter((problem) => Number.isInteger(problem.id))
  } catch (error) {
    problemsError.value = error?.response?.data?.message || error?.message || '题目列表加载失败'
  } finally {
    problemsLoading.value = false
  }
}

onMounted(loadProblems)
</script>

<style scoped>
.problem-search-results {
  max-height: 16rem;
  overflow-y: auto;
}

.problem-order-badge {
  min-width: 2.25rem;
}
</style>
