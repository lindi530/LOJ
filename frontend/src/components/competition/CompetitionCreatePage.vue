<template>
  <main>
    <div class="container py-4 py-lg-5">
      <div class="row justify-content-center">
        <div class="col-xl-9">
          <header class="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
            <div>
              <h1 class="h3 fw-semibold mb-2">创建竞赛</h1>
              <p class="text-secondary mb-0">设置赛程并选择题目，发布一场新的算法竞赛。</p>
            </div>
            <RouterLink class="btn btn-outline-secondary" :to="{ name: 'Competition' }">
              <i class="bi bi-arrow-left me-2" aria-hidden="true"></i>
              返回竞赛列表
            </RouterLink>
          </header>

          <form
            ref="formRef"
            class="card border-0 shadow-sm needs-validation"
            novalidate
            @submit.prevent="handleSubmit"
          >
            <div class="card-body p-4 p-lg-5">
              <section class="mb-4">
                <h2 class="h5 fw-semibold d-flex align-items-center gap-2 mb-4">
                  <span class="badge rounded-pill text-bg-warning">1</span>
                  基本信息
                </h2>

                <div class="mb-4">
                  <label class="form-label fw-medium" for="competition-name">竞赛名称</label>
                  <input
                    id="competition-name"
                    v-model="form.name"
                    type="text"
                    class="form-control"
                    :class="{ 'is-invalid': submitted && !form.name.trim() }"
                    placeholder="例如：第1场算法竞赛"
                    required
                  >
                  <div class="invalid-feedback">请输入竞赛名称。</div>
                </div>

                <fieldset>
                  <legend class="form-label fw-medium">可见范围</legend>
                  <div class="row g-3">
                    <div class="col-md-6">
                      <label class="border rounded p-3 d-flex gap-3 h-100">
                        <input
                          v-model="form.visibility"
                          class="form-check-input mt-1"
                          type="radio"
                          :value="true"
                          name="visibility"
                        >
                        <span>
                          <span class="fw-medium d-block">
                            <i class="bi bi-globe2 me-2" aria-hidden="true"></i>公开竞赛
                          </span>
                          <small class="text-secondary">所有用户都可以看到并参加</small>
                        </span>
                      </label>
                    </div>
                    <div class="col-md-6">
                      <label class="border rounded p-3 d-flex gap-3 h-100">
                        <input
                          v-model="form.visibility"
                          class="form-check-input mt-1"
                          type="radio"
                          :value="false"
                          name="visibility"
                        >
                        <span>
                          <span class="fw-medium d-block">
                            <i class="bi bi-eye-slash me-2" aria-hidden="true"></i>非公开竞赛
                          </span>
                          <small class="text-secondary">仅保留给指定参与场景</small>
                        </span>
                      </label>
                    </div>
                  </div>
                </fieldset>

                <div v-if="!form.visibility" class="mt-4">
                  <label class="form-label fw-medium" for="competition-password">竞赛密码</label>
                  <div class="input-group">
                    <span class="input-group-text bg-white">
                      <i class="bi bi-lock" aria-hidden="true"></i>
                    </span>
                    <input
                      id="competition-password"
                      v-model="form.password"
                      type="password"
                      class="form-control"
                      :class="{ 'is-invalid': submitted && !form.password.trim() }"
                      placeholder="请输入参赛密码"
                      autocomplete="new-password"
                      required
                    >
                    <div class="invalid-feedback">非公开竞赛必须设置密码。</div>
                  </div>
                  <div class="form-text">参赛者需要输入此密码才能加入竞赛。</div>
                </div>
              </section>

              <hr class="my-4">

              <section class="mb-4">
                <h2 class="h5 fw-semibold d-flex align-items-center gap-2 mb-4">
                  <span class="badge rounded-pill text-bg-warning">2</span>
                  时间安排
                </h2>

                <div class="row g-3">
                  <div class="col-md-6">
                    <label class="form-label fw-medium" for="competition-start-time">开始时间</label>
                    <input
                      id="competition-start-time"
                      v-model="form.startTime"
                      type="datetime-local"
                      class="form-control"
                      step="60"
                      required
                    >
                    <div class="invalid-feedback">请选择开始时间。</div>
                  </div>

                  <div class="col-md-6">
                    <label class="form-label fw-medium d-flex flex-wrap align-items-center gap-2" for="competition-end-time">
                      结束时间
                      <span v-if="durationText" class="badge text-bg-light border fw-normal">
                        持续 {{ durationText }}
                      </span>
                    </label>
                    <input
                      id="competition-end-time"
                      v-model="form.endTime"
                      type="datetime-local"
                      class="form-control"
                      :class="{ 'is-invalid': endTimeConfirmed && scheduleInvalid }"
                      step="60"
                      required
                      @change="confirmEndTime"
                    >
                    <div v-if="endTimeConfirmed && scheduleInvalid" class="invalid-feedback d-block">
                      结束时间必须晚于开始时间。
                    </div>
                    <div v-else class="invalid-feedback">请选择结束时间。</div>
                  </div>
                </div>
              </section>

              <hr class="my-4">

              <section class="mb-4">
                <h2 class="h5 fw-semibold d-flex align-items-center gap-2 mb-4">
                  <span class="badge rounded-pill text-bg-warning">3</span>
                  竞赛公告
                </h2>
                <label class="form-label fw-medium" for="competition-announcement">公告内容</label>
                <textarea
                  id="competition-announcement"
                  v-model="form.announcement"
                  class="form-control"
                  :class="{ 'is-invalid': submitted && !form.announcement.trim() }"
                  rows="4"
                  placeholder="例如：这是一场刺激的算法竞赛"
                  required
                ></textarea>
                <div class="invalid-feedback">请输入竞赛公告。</div>
              </section>

              <hr class="my-4">

              <section>
                <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3">
                  <h2 class="h5 fw-semibold d-flex align-items-center gap-2 mb-0">
                    <span class="badge rounded-pill text-bg-warning">4</span>
                    选择题目
                  </h2>
                  <span class="text-secondary small">已选 {{ form.problemIds.length }} 题</span>
                </div>

                <label class="form-label fw-medium" for="problem-search">搜索并添加题目</label>
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
                      id="problem-search"
                      v-model="problemKeyword"
                      type="search"
                      class="form-control"
                      :inputmode="problemSearchMode === 'id' ? 'numeric' : 'search'"
                      :placeholder="problemSearchMode === 'id' ? '输入题号' : '输入题目名称'"
                      autocomplete="off"
                      @keydown.esc="problemKeyword = ''"
                    >
                  </div>

                  <div v-if="problemKeyword.trim()" class="list-group mt-1 shadow-sm search-results">
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
                      {{ problemsError }}，点击重新加载
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
                    <span class="badge rounded-pill text-bg-warning selected-problem-number">
                      {{ problemNumber(index) }}
                    </span>
                    <span class="text-secondary">#{{ problem.id }}</span>
                    <span class="fw-medium flex-grow-1">{{ problem.title }}</span>
                    <span v-if="problem.level" class="badge text-bg-light fw-normal">{{ problem.level }}</span>
                    <div class="btn-group btn-group-sm" role="group" :aria-label="`${problem.title} 排序`">
                      <button
                        type="button"
                        class="btn btn-outline-secondary"
                        :disabled="index === 0"
                        :aria-label="`上移题目 ${problem.title}`"
                        @click="moveProblem(index, -1)"
                      >
                        上移
                      </button>
                      <button
                        type="button"
                        class="btn btn-outline-secondary"
                        :disabled="index === selectedProblems.length - 1"
                        :aria-label="`下移题目 ${problem.title}`"
                        @click="moveProblem(index, 1)"
                      >
                        下移
                      </button>
                    </div>
                    <button type="button" class="btn btn-outline-secondary btn-sm" @click="removeProblem(problem.id)">
                      移除
                    </button>
                  </div>
                </div>

                <p v-if="submitted && form.problemIds.length === 0" class="text-danger small mt-2 mb-0">
                  请至少选择一道竞赛题目。
                </p>
              </section>

              <div v-if="submitError" class="alert alert-danger mt-4 mb-0" role="alert">
                {{ submitError }}
              </div>
            </div>

            <div class="card-footer bg-white border-top p-4 d-flex flex-column-reverse flex-sm-row justify-content-end gap-2">
              <RouterLink class="btn btn-outline-secondary px-4" :to="{ name: 'Competition' }">取消</RouterLink>
              <button class="btn btn-warning px-4 text-white" type="submit" :disabled="submitting">
                <span v-if="submitting" class="spinner-border spinner-border-sm me-2" aria-hidden="true"></span>
                {{ submitting ? '创建中...' : '创建竞赛' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useDialog } from 'naive-ui'
import { useRouter } from 'vue-router'
import api from '@/api'

const router = useRouter()
const dialog = useDialog()
const formRef = ref(null)
const submitted = ref(false)
const submitting = ref(false)
const submitError = ref('')
const endTimeConfirmed = ref(false)
const problems = ref([])
const problemsLoading = ref(false)
const problemsError = ref('')
const problemKeyword = ref('')
const problemSearchMode = ref('id')

const form = reactive({
  name: '',
  visibility: true,
  password: '',
  startTime: '',
  endTime: '',
  announcement: '',
  problemIds: []
})

const selectedProblems = computed(() => (
  form.problemIds
    .map((problemId) => problems.value.find((problem) => problem.id === problemId))
    .filter(Boolean)
))

const searchResults = computed(() => {
  const keyword = problemKeyword.value.trim()

  if (!keyword) {
    return []
  }

  return problems.value.filter((problem) => {
    if (form.problemIds.includes(problem.id)) {
      return false
    }

    if (problemSearchMode.value === 'id') {
      return String(problem.id).startsWith(keyword)
    }

    return problem.title.toLowerCase().includes(keyword.toLowerCase())
  })
})

const scheduleInvalid = computed(() => {
  const startTime = parseLocalTime(form.startTime)
  const endTime = parseLocalTime(form.endTime)
  return Boolean(startTime && endTime && endTime <= startTime)
})

const durationText = computed(() => {
  if (!endTimeConfirmed.value || scheduleInvalid.value) {
    return ''
  }

  const startTime = parseLocalTime(form.startTime)
  const endTime = parseLocalTime(form.endTime)
  if (!startTime || !endTime) {
    return ''
  }

  const totalMinutes = Math.floor((endTime.getTime() - startTime.getTime()) / 60000)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60

  if (hours && minutes) {
    return `${hours} 小时 ${minutes} 分钟`
  }

  if (hours) {
    return `${hours} 小时`
  }

  return `${minutes} 分钟`
})

function parseLocalTime(value) {
  if (!value) {
    return null
  }

  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? null : date
}

function toRFC3339(value) {
  const date = parseLocalTime(value)
  if (!date) {
    return ''
  }

  const offset = -date.getTimezoneOffset()
  const sign = offset >= 0 ? '+' : '-'
  const offsetHours = pad(Math.floor(Math.abs(offset) / 60))
  const offsetMinutes = pad(Math.abs(offset) % 60)

  return [
    `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`,
    `T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`,
    `${sign}${offsetHours}:${offsetMinutes}`
  ].join('')
}

function pad(value) {
  return String(value).padStart(2, '0')
}

function problemNumber(index) {
  let position = index + 1
  let number = ''

  while (position > 0) {
    position -= 1
    number = String.fromCharCode(65 + (position % 26)) + number
    position = Math.floor(position / 26)
  }

  return number
}

function confirmEndTime() {
  endTimeConfirmed.value = true
}

function addProblem(problem) {
  form.problemIds.push(problem.id)
  problemKeyword.value = ''
}

function removeProblem(problemId) {
  form.problemIds = form.problemIds.filter((id) => id !== problemId)
}

function moveProblem(index, direction) {
  const destination = index + direction

  if (destination < 0 || destination >= form.problemIds.length) {
    return
  }

  const currentProblemId = form.problemIds[index]
  form.problemIds[index] = form.problemIds[destination]
  form.problemIds[destination] = currentProblemId
}

function showSubmitError(message) {
  submitError.value = message
  dialog.error({
    title: '创建失败',
    content: message,
    positiveText: '确认'
  })
}

async function loadProblems() {
  problemsLoading.value = true
  problemsError.value = ''

  try {
    const resp = await api.getProblemList()
    if (resp.code === 0 && Array.isArray(resp.data)) {
      problems.value = resp.data
        .map((problem) => ({
          id: Number(problem.id),
          title: problem.title || '无标题',
          level: problem.level || ''
        }))
        .filter((problem) => Number.isInteger(problem.id))
    } else {
      problemsError.value = resp.message || '题目列表加载失败'
    }
  } catch (error) {
    problemsError.value = error?.message || '题目列表加载失败'
  } finally {
    problemsLoading.value = false
  }
}

async function handleSubmit() {
  submitted.value = true
  endTimeConfirmed.value = true
  submitError.value = ''
  formRef.value.classList.add('was-validated')

  if (!formRef.value.checkValidity()) {
    return
  }

  if (!form.name.trim()) {
    submitError.value = '请输入竞赛名称。'
    return
  }

  if (!form.visibility && !form.password.trim()) {
    submitError.value = '非公开竞赛必须设置密码。'
    return
  }

  if (scheduleInvalid.value) {
    submitError.value = '结束时间必须晚于开始时间。'
    return
  }

  if (!form.announcement.trim()) {
    submitError.value = '请输入竞赛公告。'
    return
  }

  if (form.problemIds.length === 0) {
    submitError.value = '请至少选择一道竞赛题目。'
    return
  }

  const payload = {
    name: form.name.trim(),
    visibility: form.visibility,
    password: form.visibility ? '' : form.password,
    start_time: toRFC3339(form.startTime),
    end_time: toRFC3339(form.endTime),
    announcement: form.announcement.trim(),
    problems: form.problemIds.map((problemId, index) => ({
      problem_number: problemNumber(index),
      problem_id: problemId
    }))
  }

  submitting.value = true
  try {
    const resp = await api.createCompetition(payload)
    if (resp.code === 0) {
      dialog.success({
        title: '创建成功',
        content: '竞赛已成功创建。',
        positiveText: '查看竞赛列表',
        closable: false,
        maskClosable: false,
        closeOnEsc: false,
        onPositiveClick: () => router.push({ name: 'Competition' })
      })
      return
    }

    showSubmitError(resp.message || '竞赛创建失败，请稍后重试。')
  } catch (error) {
    showSubmitError(error?.message || '竞赛创建失败，请稍后重试。')
  } finally {
    submitting.value = false
  }
}

onMounted(loadProblems)
</script>

<style scoped>
.search-results {
  max-height: 16rem;
  overflow-y: auto;
}

.selected-problem-number {
  min-width: 2.25rem;
}
</style>
