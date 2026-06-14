<template>
  <main class="py-4 py-lg-5">
    <div class="container-xl">
      <div class="d-flex flex-column flex-lg-row justify-content-between gap-3 align-items-lg-center mb-4">
        <div>
          <nav aria-label="章节创建路径" class="mb-2">
            <ol class="breadcrumb mb-0">
              <li class="breadcrumb-item">
                <RouterLink :to="{ name: 'Course' }">课程</RouterLink>
              </li>
              <li class="breadcrumb-item">
                <RouterLink :to="{ name: 'CourseDetail', params: { course_id: courseId }, hash: '#content' }">
                  课程内容
                </RouterLink>
              </li>
              <li class="breadcrumb-item active" aria-current="page">添加章节</li>
            </ol>
          </nav>
          <h1 class="h3 fw-bold mb-2">添加课程章节</h1>
          <p class="text-secondary mb-0">先上传章节视频，再填写章节信息与练习题目。</p>
        </div>
        <RouterLink
          class="btn btn-outline-secondary align-self-start align-self-lg-center"
          :to="{ name: 'CourseDetail', params: { course_id: courseId }, hash: '#content' }"
        >
          <i class="bi bi-arrow-left me-2" aria-hidden="true"></i>
          返回课程内容
        </RouterLink>
      </div>

      <form ref="formRef" class="row g-4 needs-validation" novalidate @submit.prevent="handleSubmit">
        <div class="col-lg-7">
          <section class="card border-0 shadow-sm mb-4">
            <div class="card-body p-4 p-lg-5">
              <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-4">
                <div>
                  <h2 class="h5 fw-bold mb-1">基础信息</h2>
                  <p class="text-secondary small mb-0">章节标题、简介和排序会直接进入课程目录。</p>
                </div>
                <span class="badge text-bg-primary">课程 #{{ courseId }}</span>
              </div>

              <div class="mb-4">
                <label class="form-label fw-medium" for="chapter-title">章节标题</label>
                <input
                  id="chapter-title"
                  v-model.trim="form.title"
                  class="form-control form-control-lg"
                  :class="{ 'is-invalid': submitted && !form.title }"
                  type="text"
                  maxlength="80"
                  placeholder="例如：二分查找与边界处理"
                  required
                >
                <div class="invalid-feedback">请输入章节标题。</div>
              </div>

              <div class="mb-4">
                <label class="form-label fw-medium" for="chapter-description">章节简介</label>
                <textarea
                  id="chapter-description"
                  v-model.trim="form.description"
                  class="form-control"
                  rows="4"
                  maxlength="300"
                  placeholder="简要说明本章节覆盖的知识点，可留空。"
                ></textarea>
                <div class="d-flex justify-content-end mt-1">
                  <div class="form-text">{{ form.description.length }}/300</div>
                </div>
              </div>

              <div class="row g-3">
                <div class="col-md-6">
                  <label class="form-label fw-medium" for="chapter-sort">章节排序</label>
                  <input
                    id="chapter-sort"
                    v-model.number="form.sort"
                    class="form-control"
                    :class="{ 'is-invalid': submitted && !sortValid }"
                    type="number"
                    min="1"
                    step="1"
                    required
                  >
                  <div class="invalid-feedback">排序请输入不小于 1 的整数。</div>
                </div>
                <div class="col-md-6">
                  <div class="border rounded p-3 h-100 bg-body-tertiary">
                    <div class="fw-semibold mb-1">提交状态</div>
                    <div class="small text-secondary">
                      {{ videoUploading ? '视频上传中，暂不可提交。' : '视频上传完成后即可提交章节。' }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </section>

          <section class="card border-0 shadow-sm">
            <div class="card-body p-4 p-lg-5">
              <ProblemSelector
                v-model="selectedProblems"
                title="章节练习"
                description="题目顺序会作为章节内练习排序上传。"
                required
                :submitted="submitted"
              />
            </div>
          </section>
        </div>

        <div class="col-lg-5">
          <section class="card border-0 shadow-sm mb-4">
            <div class="card-body p-4">
              <VideoChunkUpload
                @uploaded="handleVideoUploaded"
                @uploading-change="videoUploading = $event"
              />
            </div>
          </section>

          <section class="card border-0 shadow-sm">
            <div class="card-body p-4">
              <h2 class="h5 fw-bold mb-3">提交预览</h2>
              <div class="list-group list-group-flush mb-4">
                <div class="list-group-item px-0 d-flex justify-content-between gap-3">
                  <span class="text-secondary">课程 ID</span>
                  <span class="fw-semibold">#{{ courseId }}</span>
                </div>
                <div class="list-group-item px-0 d-flex justify-content-between gap-3">
                  <span class="text-secondary">章节排序</span>
                  <span class="fw-semibold">{{ normalizedSort }}</span>
                </div>
                <div class="list-group-item px-0 d-flex justify-content-between gap-3">
                  <span class="text-secondary">练习题目</span>
                  <span class="fw-semibold">{{ selectedProblems.length }} 题</span>
                </div>
                <div class="list-group-item px-0">
                  <div class="text-secondary mb-1">视频地址</div>
                  <div class="small text-break">{{ videoUrl || '等待视频上传完成' }}</div>
                </div>
              </div>

              <div v-if="submitError" class="alert alert-danger d-flex gap-3">
                <i class="bi bi-exclamation-triangle-fill mt-1" aria-hidden="true"></i>
                <div>{{ submitError }}</div>
              </div>

              <div class="d-grid gap-2">
                <button class="btn btn-primary" type="submit" :disabled="submitDisabled">
                  <span v-if="submitting" class="spinner-border spinner-border-sm me-2" aria-hidden="true"></span>
                  <i v-else class="bi bi-check2-circle me-2" aria-hidden="true"></i>
                  {{ submitting ? '提交中' : '提交章节' }}
                </button>
                <RouterLink
                  class="btn btn-outline-secondary"
                  :to="{ name: 'CourseDetail', params: { course_id: courseId }, hash: '#content' }"
                >
                  取消
                </RouterLink>
              </div>
            </div>
          </section>
        </div>
      </form>
    </div>
  </main>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { useDialog } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import ProblemSelector from '@/components/ProblemUpload/ProblemSelector.vue'
import VideoChunkUpload from '@/components/upload/VideoChunkUpload.vue'
import api from '@/api'

const route = useRoute()
const router = useRouter()
const dialog = useDialog()
const formRef = ref(null)
const submitted = ref(false)
const submitting = ref(false)
const submitError = ref('')
const selectedProblems = ref([])
const videoUploading = ref(false)
const videoAsset = ref({
  url: '',
  id: null
})

const form = reactive({
  title: '',
  description: '',
  sort: 1
})

const courseId = computed(() => Number(route.params.course_id))
const normalizedSort = computed(() => Number(form.sort || 1))
const sortValid = computed(() => Number.isInteger(normalizedSort.value) && normalizedSort.value >= 1)
const submitDisabled = computed(() => submitting.value || videoUploading.value)
const videoUrl = computed(() => videoAsset.value.url)
const hasVideoAsset = computed(() => {
  const id = videoAsset.value.id

  return Boolean(videoAsset.value.url && id !== null && id !== undefined && id !== '')
})

function normalizeVideoAsset(video) {
  if (typeof video === 'string') {
    return {
      url: video,
      id: null
    }
  }

  return {
    url: video?.url || video?.video_url || video?.videoUrl || '',
    id: video?.id ?? video?.video_id ?? video?.videoId ?? null
  }
}

function handleVideoUploaded(video) {
  videoAsset.value = normalizeVideoAsset(video)
}

function buildPayload() {
  const description = form.description.trim()

  return {
    course_id: courseId.value,
    title: form.title.trim(),
    description: description || null,
    sort: normalizedSort.value,
    problems: selectedProblems.value.map((problem, index) => ({
      id: Number(problem.id),
      sort: index + 1
    })),
    video: {
      url: videoAsset.value.url,
      id: videoAsset.value.id
    }
  }
}

function normalizeApiResult(resp) {
  if (resp?.data && typeof resp.data === 'object' && 'code' in resp.data) {
    return resp.data
  }

  return resp
}

function validateForm() {
  formRef.value?.classList.add('was-validated')

  if (!courseId.value) {
    submitError.value = '课程 ID 无效。'
    return false
  }

  if (!form.title.trim()) {
    submitError.value = '请输入章节标题。'
    return false
  }

  if (!sortValid.value) {
    submitError.value = '章节排序请输入不小于 1 的整数。'
    return false
  }

  if (selectedProblems.value.length === 0) {
    submitError.value = '请至少选择一道章节练习题。'
    return false
  }

  if (videoUploading.value) {
    submitError.value = '视频上传中，请等待上传完成后再提交。'
    return false
  }

  if (!hasVideoAsset.value) {
    submitError.value = '请先上传章节视频。'
    return false
  }

  return true
}

async function handleSubmit() {
  submitted.value = true
  submitError.value = ''

  if (!validateForm()) {
    return
  }

  submitting.value = true

  try {
    const resp = await api.createCourseChapter(buildPayload())
    const result = normalizeApiResult(resp)

    if (result?.code === 0) {
      dialog.success({
        title: '章节创建成功',
        content: result?.message || '章节已添加到课程内容中。',
        positiveText: '返回课程内容',
        closable: false,
        maskClosable: false,
        closeOnEsc: false,
        onPositiveClick: () => router.push({
          name: 'CourseDetail',
          params: { course_id: courseId.value },
          hash: '#content'
        })
      })
      return
    }

    submitError.value = result?.message || '章节创建失败，请稍后重试。'
  } catch (error) {
    submitError.value = error?.response?.data?.message || error?.message || '请求发送失败，请稍后重试。'
  } finally {
    submitting.value = false
  }
}
</script>
