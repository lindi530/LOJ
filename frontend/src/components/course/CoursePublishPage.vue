<template>
  <main class="py-4 py-lg-5">
    <div class="container-xl">
      <div class="d-flex flex-column flex-lg-row justify-content-between gap-3 align-items-lg-center mb-4">
        <div>
          <nav aria-label="课程发布路径" class="mb-2">
            <ol class="breadcrumb mb-0">
              <li class="breadcrumb-item">
                <RouterLink :to="{ name: 'Course' }">课程</RouterLink>
              </li>
              <li class="breadcrumb-item active" aria-current="page">发布课程</li>
            </ol>
          </nav>
          <h1 class="h3 fw-bold mb-2">发布课程</h1>
          <p class="text-secondary mb-0">整理封面、定价和发布状态，形成一份清晰可维护的课程草案。</p>
        </div>
        <RouterLink class="btn btn-outline-secondary align-self-start align-self-lg-center" :to="{ name: 'Course' }">
          <i class="bi bi-arrow-left me-2" aria-hidden="true"></i>
          返回课程列表
        </RouterLink>
      </div>

      <form ref="formRef" class="row g-4 needs-validation" novalidate @submit.prevent="handleSubmit">
        <div class="col-lg-5">
          <section class="card border-0 shadow-sm h-100">
            <div class="card-body p-4">
              <div class="d-flex justify-content-between align-items-center gap-3 mb-3">
                <div>
                  <h2 class="h5 fw-bold mb-1">课程封面</h2>
                  <p class="text-secondary small mb-0">建议使用 16:9 图片，列表和详情页展示更稳定。</p>
                </div>
                <span class="badge text-bg-light border">cover</span>
              </div>

              <label class="course-cover-dropzone d-block overflow-hidden" for="course-cover">
                <span class="ratio ratio-16x9">
                  <img
                    v-if="coverPreview"
                    class="course-cover-preview"
                    :src="coverPreview"
                    alt="课程封面预览"
                  >
                  <span v-else class="course-cover-placeholder d-flex flex-column align-items-center justify-content-center text-center p-4">
                    <i class="bi bi-image fs-1 mb-3" aria-hidden="true"></i>
                    <span class="fw-semibold">选择课程封面</span>
                    <span class="small">JPG / PNG / GIF / JPEG / WebP</span>
                  </span>
                </span>
              </label>
              <input
                id="course-cover"
                ref="coverInputRef"
                class="visually-hidden"
                type="file"
                accept=".jpg,.png,.gif,.jpeg,.webp,image/jpeg,image/png,image/gif,image/webp"
                @change="handleCoverChange"
              >
              <div v-if="submitted && !form.cover" class="text-danger small mt-2">
                请上传课程封面。
              </div>

              <div class="alert alert-light border d-flex gap-3 mt-4 mb-0">
                <i class="bi bi-info-circle text-primary mt-1" aria-hidden="true"></i>
                <div class="small text-secondary">
                  免费课程会自动将价格归零；付费课程需要填写不小于 0 的价格。
                </div>
              </div>
            </div>
          </section>
        </div>

        <div class="col-lg-7">
          <section class="card border-0 shadow-sm">
            <div class="card-body p-4 p-lg-5">
              <div class="d-flex flex-column flex-sm-row justify-content-between align-items-start align-items-sm-center gap-2 mb-4">
                <div>
                  <h2 class="h5 fw-bold mb-1">基础信息</h2>
                  <p class="text-secondary small mb-0">字段保持和课程数据模型一致，便于后续接入真实流程。</p>
                </div>
                <span class="badge" :class="form.isFree ? 'text-bg-success' : 'text-bg-warning'">
                  {{ form.isFree ? '免费课程' : '付费课程' }}
                </span>
              </div>

              <div class="mb-4">
                <label class="form-label fw-medium" for="course-title">课程标题</label>
                <input
                  id="course-title"
                  v-model.trim="form.title"
                  class="form-control form-control-lg"
                  :class="{ 'is-invalid': submitted && !form.title }"
                  type="text"
                  maxlength="60"
                  placeholder="例如：Go 语言算法训练营"
                  required
                >
                <div class="invalid-feedback">请输入课程标题。</div>
              </div>

              <div class="mb-4">
                <label class="form-label fw-medium" for="course-description">课程描述</label>
                <textarea
                  id="course-description"
                  v-model.trim="form.description"
                  class="form-control"
                  :class="{ 'is-invalid': submitted && !form.description }"
                  rows="5"
                  maxlength="300"
                  placeholder="用一段话说明课程适合谁、能学到什么。"
                  required
                ></textarea>
                <div class="d-flex justify-content-between gap-3 mt-1">
                  <div class="invalid-feedback d-block" :class="{ invisible: !submitted || form.description }">
                    请输入课程描述。
                  </div>
                  <div class="form-text text-nowrap">{{ form.description.length }}/300</div>
                </div>
              </div>

              <div class="row g-3 mb-4">
                <div class="col-md-6">
                  <label class="form-label fw-medium" for="course-price">价格</label>
                  <div class="input-group">
                    <span class="input-group-text">￥</span>
                    <input
                      id="course-price"
                      v-model.number="form.price"
                      class="form-control"
                      :class="{ 'is-invalid': submitted && !priceValid }"
                      type="number"
                      min="0"
                      step="0.01"
                      :disabled="form.isFree"
                      required
                    >
                    <div class="invalid-feedback">价格不能小于 0。</div>
                  </div>
                </div>

                <div class="col-md-6">
                  <label class="form-label fw-medium" for="course-sort">排序</label>
                  <input
                    id="course-sort"
                    v-model.number="form.sort"
                    class="form-control"
                    :class="{ 'is-invalid': submitted && !sortValid }"
                    type="number"
                    min="0"
                    step="1"
                    placeholder="数值越小越靠前"
                    required
                  >
                  <div class="invalid-feedback">排序请输入 0 或正整数。</div>
                </div>
              </div>

              <div class="row g-3 mb-4">
                <div class="col-md-6">
                  <label class="form-label fw-medium" for="course-status">状态</label>
                  <select id="course-status" v-model="form.status" class="form-select">
                    <option
                      v-for="option in statusOptions"
                      :key="option.value"
                      :value="option.value"
                    >
                      {{ option.label }}
                    </option>
                  </select>
                </div>

                <div class="col-md-6">
                  <div class="border rounded p-3 h-100 bg-body-tertiary">
                    <div class="form-check form-switch mb-1">
                      <input
                        id="course-is-free"
                        v-model="form.isFree"
                        class="form-check-input"
                        type="checkbox"
                        role="switch"
                      >
                      <label class="form-check-label fw-medium" for="course-is-free">免费课程</label>
                    </div>
                    <div class="small text-secondary">开启后价格字段会锁定为 ￥0.00。</div>
                  </div>
                </div>
              </div>

              <hr class="my-4">

              <div class="d-flex flex-column flex-sm-row justify-content-between gap-3 align-items-sm-center">
                <div class="text-secondary small">
                  当前状态：
                  <span class="fw-semibold text-body">{{ currentStatusLabel }}</span>
                  <span class="mx-2">·</span>
                  排序 {{ form.sort || 0 }}
                </div>
                <div class="d-flex gap-2">
                  <button class="btn btn-outline-secondary" type="button" :disabled="submitting" @click="resetForm">
                    <i class="bi bi-arrow-counterclockwise me-2" aria-hidden="true"></i>
                    重置
                  </button>
                  <button class="btn btn-primary" type="submit" :disabled="submitting">
                    <span
                      v-if="submitting"
                      class="spinner-border spinner-border-sm me-2"
                      aria-hidden="true"
                    ></span>
                    <i v-else class="bi bi-send-check me-2" aria-hidden="true"></i>
                    {{ submitting ? '发布中' : '发布课程' }}
                  </button>
                </div>
              </div>

              <div v-if="submitSuccess" class="alert alert-success d-flex gap-3 mt-4 mb-0">
                <i class="bi bi-check-circle-fill mt-1" aria-hidden="true"></i>
                <div>
                  <div class="fw-semibold">{{ submitMessage }}</div>
                  <div class="small">标题：{{ submitSnapshot.title }}；状态：{{ submitSnapshot.status }}。</div>
                </div>
              </div>

              <div v-if="submitError" class="alert alert-danger d-flex gap-3 mt-4 mb-0">
                <i class="bi bi-exclamation-triangle-fill mt-1" aria-hidden="true"></i>
                <div>
                  <div class="fw-semibold">课程发布失败</div>
                  <div class="small">{{ submitError }}</div>
                </div>
              </div>
            </div>
          </section>
        </div>
      </form>
    </div>
  </main>
</template>

<script setup>
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import api from '@/api'

const defaultForm = () => ({
  cover: null,
  title: '',
  description: '',
  price: 0,
  isFree: true,
  status: 0,
  sort: 0
})

const statusOptions = [
  { value: 0, label: '草稿' },
  { value: 1, label: '已发布' },
  { value: 2, label: '已下架' }
]

const formRef = ref(null)
const coverInputRef = ref(null)
const coverPreview = ref('')
const submitted = ref(false)
const submitting = ref(false)
const submitSuccess = ref(false)
const submitError = ref('')
const submitMessage = ref('')
const submitSnapshot = reactive({
  title: '',
  status: ''
})
const form = reactive(defaultForm())

const priceValid = computed(() => form.isFree || Number(form.price) >= 0)
const sortValid = computed(() => Number.isInteger(Number(form.sort)) && Number(form.sort) >= 0)
const currentStatusLabel = computed(() => (
  statusOptions.find((option) => option.value === form.status)?.label || '草稿'
))

watch(
  () => form.isFree,
  (isFree) => {
    if (isFree) {
      form.price = 0
    }
  }
)

const revokeCoverPreview = () => {
  if (coverPreview.value) {
    URL.revokeObjectURL(coverPreview.value)
    coverPreview.value = ''
  }
}

const handleCoverChange = (event) => {
  const [file] = event.target.files || []

  revokeCoverPreview()
  form.cover = file || null

  if (file) {
    coverPreview.value = URL.createObjectURL(file)
  }

  submitSuccess.value = false
  submitError.value = ''
}

const formValid = computed(() => (
  Boolean(form.cover) &&
  Boolean(form.title) &&
  Boolean(form.description) &&
  priceValid.value &&
  sortValid.value
))

const buildPayload = () => {
  const payload = new FormData()

  payload.append('cover', form.cover)
  payload.append('title', form.title)
  payload.append('description', form.description)
  payload.append('price', String(form.isFree ? 0 : Number(form.price || 0)))
  payload.append('isFree', String(form.isFree))
  payload.append('status', String(form.status))
  payload.append('sort', String(Number(form.sort || 0)))

  return payload
}

const handleSubmit = async () => {
  submitted.value = true
  submitSuccess.value = false
  submitError.value = ''
  submitMessage.value = ''

  if (!formValid.value) {
    return
  }

  submitting.value = true

  try {
    const resp = await api.courseCreate(buildPayload())

    if (resp?.code === 0) {
      submitSnapshot.title = form.title
      submitSnapshot.status = currentStatusLabel.value
      submitMessage.value = resp.message || '课程发布成功。'
      submitSuccess.value = true
      return
    }

    submitError.value = resp?.message || '请稍后重试。'
  } catch (error) {
    submitError.value = error?.response?.data?.message || error?.message || '请求发送失败，请稍后重试。'
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  Object.assign(form, defaultForm())
  submitted.value = false
  submitting.value = false
  submitSuccess.value = false
  submitError.value = ''
  submitMessage.value = ''
  revokeCoverPreview()

  if (coverInputRef.value) {
    coverInputRef.value.value = ''
  }
}

onBeforeUnmount(revokeCoverPreview)
</script>

<style scoped>
.course-cover-dropzone {
  border: var(--bs-border-width) dashed var(--bs-border-color);
  border-radius: var(--bs-border-radius-lg);
  background-color: var(--bs-body-bg);
  cursor: pointer;
  transition: border-color .15s ease, background-color .15s ease;
}

.course-cover-dropzone:hover,
.course-cover-dropzone:focus-within {
  border-color: var(--bs-primary);
  background-color: var(--bs-primary-bg-subtle);
}

.course-cover-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.course-cover-placeholder {
  color: var(--bs-secondary-color);
}
</style>
