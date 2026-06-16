<template>
  <section>
    <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3">
      <div>
        <h2 class="h5 fw-bold mb-1">章节视频</h2>
        <p class="text-secondary small mb-0">页面进入后自动创建上传任务，选择视频后按后端分块规则上传。</p>
      </div>
      <span class="badge text-bg-light border">
        {{ uploadId ? `任务 ${uploadId}` : '创建任务中' }}
      </span>
    </div>

    <div v-if="taskError" class="alert alert-danger d-flex align-items-center gap-3">
      <i class="bi bi-exclamation-triangle-fill" aria-hidden="true"></i>
      <div class="flex-grow-1">{{ taskError }}</div>
      <button class="btn btn-outline-danger btn-sm" type="button" @click="createUploadTask">重试</button>
    </div>

    <label class="video-dropzone d-block p-4 text-center" :for="inputId">
      <input
        :id="inputId"
        ref="fileInputRef"
        class="visually-hidden"
        type="file"
        accept="video/*"
        :disabled="isBusy"
        @change="handleFileChange"
      >
      <span class="d-flex flex-column align-items-center gap-2">
        <i class="bi bi-cloud-arrow-up fs-1 text-primary" aria-hidden="true"></i>
        <span class="fw-semibold">{{ file ? file.name : '选择视频文件' }}</span>
        <span class="text-secondary small">
          <template v-if="file">
            {{ formatFileSize(file.size) }} · {{ chunkCount }} 个分块 · 每块 {{ chunkSizeLabel }}
          </template>
          <template v-else>
            支持常见视频格式，上传完成后才能提交章节。
          </template>
        </span>
      </span>
    </label>

    <div class="d-flex flex-column flex-sm-row justify-content-between align-items-sm-center gap-2 mt-3">
      <div class="text-secondary small">
        <i class="bi bi-info-circle me-1" aria-hidden="true"></i>
        {{ statusText }}
      </div>
      <div class="d-flex flex-wrap gap-2">
        <button
          class="btn btn-outline-secondary"
          type="button"
          :disabled="isBusy || !file"
          @click="clearFile"
        >
          <i class="bi bi-arrow-counterclockwise me-2" aria-hidden="true"></i>
          重选
        </button>
        <button
          class="btn btn-primary"
          type="button"
          :disabled="!canStartUpload"
          @click="startUpload"
        >
          <span v-if="isBusy" class="spinner-border spinner-border-sm me-2" aria-hidden="true"></span>
          <i v-else class="bi bi-play-circle me-2" aria-hidden="true"></i>
          {{ uploadButtonText }}
        </button>
      </div>
    </div>

    <div class="progress mt-3" role="progressbar" :aria-valuenow="progressPercent" aria-valuemin="0" aria-valuemax="100">
      <div class="progress-bar" :class="progressClass" :style="{ width: `${progressPercent}%` }">
        {{ progressPercent }}%
      </div>
    </div>

    <div v-if="videoOriginPath" class="alert alert-success d-flex gap-3 mt-3 mb-0">
      <i class="bi bi-check-circle-fill mt-1" aria-hidden="true"></i>
      <div>
        <div class="fw-semibold">视频上传完成</div>
        <div class="small text-break">{{ videoOriginPath }}</div>
      </div>
    </div>

    <div v-if="uploadError" class="alert alert-danger d-flex gap-3 mt-3 mb-0">
      <i class="bi bi-exclamation-triangle-fill mt-1" aria-hidden="true"></i>
      <div>{{ uploadError }}</div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import api from '@/api'
import { createMd5 } from '@/utils/md5'

const emit = defineEmits(['uploaded', 'uploading-change'])

const concurrency = 3
const bytesPerMb = 1024 * 1024
const defaultChunkSizeBytes = bytesPerMb
const inputId = `video-file-${Math.random().toString(36).slice(2, 8)}`
const fileInputRef = ref(null)
const file = ref(null)
const uploadId = ref('')
const chunkSizeBytes = ref(defaultChunkSizeBytes)
const status = ref('initializing')
const taskError = ref('')
const uploadError = ref('')
const uploadedBytes = ref(0)
const uploadingBytesTotal = ref(0)
const hashingPercent = ref(0)
const videoOriginPath = ref('')
const videoId = ref(null)

const chunkCount = computed(() => (file.value ? Math.ceil(file.value.size / chunkSizeBytes.value) : 0))
const chunkSizeLabel = computed(() => formatFileSize(chunkSizeBytes.value))
const isBusy = computed(() => ['checking', 'uploading', 'hashing', 'merging'].includes(status.value))
const canStartUpload = computed(() => Boolean(file.value && uploadId.value && !isBusy.value && status.value !== 'success'))

const progressPercent = computed(() => {
  if (status.value === 'success') {
    return 100
  }

  if (status.value === 'hashing') {
    return Math.min(90, 80 + Math.round(hashingPercent.value * 0.1))
  }

  if (status.value === 'merging') {
    return 95
  }

  if (status.value === 'checking') {
    return 5
  }

  if (status.value === 'uploading' && uploadingBytesTotal.value > 0) {
    return Math.min(80, Math.round((uploadedBytes.value / uploadingBytesTotal.value) * 80))
  }

  return 0
})

const statusText = computed(() => {
  const textMap = {
    initializing: '正在创建上传任务',
    ready: uploadId.value ? '上传任务已就绪，请选择视频文件。' : '等待上传任务创建。',
    checking: '正在确认服务端需要上传的分块。',
    uploading: '正在并发上传需要的分块。',
    hashing: '正在计算原始文件 MD5。',
    merging: '正在请求服务端合并视频。',
    success: '视频已上传完成，可以提交章节。',
    error: uploadError.value || taskError.value || '上传失败，请重试。'
  }

  return textMap[status.value] || '等待操作。'
})

const uploadButtonText = computed(() => {
  if (isBusy.value) {
    return '上传中'
  }

  if (status.value === 'success') {
    return '重新上传'
  }

  return '开始上传'
})

const progressClass = computed(() => {
  if (status.value === 'success') {
    return 'bg-success'
  }

  if (status.value === 'error') {
    return 'bg-danger'
  }

  return ''
})

watch(isBusy, (value) => emit('uploading-change', value), { immediate: true })

function extractResponseData(resp, fallbackMessage) {
  if (resp?.code !== undefined && resp.code !== 0) {
    throw new Error(resp.message || fallbackMessage)
  }

  return resp?.data ?? {}
}

function extractTaskData(resp) {
  const data = extractResponseData(resp, '上传任务创建失败')
  const chunkSizeMb = Number(data.chunk_size ?? data.chunkSize ?? 0)
  const explicitChunkSizeMb = Number(data.chunk_size_mb ?? data.chunkSizeMb ?? 0)
  const normalizedChunkSizeMb = chunkSizeMb > 0 ? chunkSizeMb : explicitChunkSizeMb

  return {
    uploadId: data.upload_id ?? data.uploadId ?? data.id ?? '',
    chunkSizeBytes: normalizedChunkSizeMb > 0 ? normalizedChunkSizeMb * bytesPerMb : 0
  }
}

async function createUploadTask() {
  status.value = 'initializing'
  taskError.value = ''

  try {
    const task = extractTaskData(await api.createVideoUploadTask())
    uploadId.value = task.uploadId
    chunkSizeBytes.value = Number.isFinite(task.chunkSizeBytes) && task.chunkSizeBytes > 0
      ? task.chunkSizeBytes
      : defaultChunkSizeBytes

    if (!uploadId.value) {
      throw new Error('后端未返回 upload_id')
    }

    status.value = 'ready'
  } catch (error) {
    taskError.value = error?.response?.data?.message || error?.message || '上传任务创建失败'
    status.value = 'error'
  }
}

function handleFileChange(event) {
  const [selectedFile] = event.target.files || []

  file.value = selectedFile || null
  videoOriginPath.value = ''
  videoId.value = null
  uploadError.value = ''
  uploadedBytes.value = 0
  uploadingBytesTotal.value = 0
  hashingPercent.value = 0
  emit('uploaded', '')

  if (selectedFile && uploadId.value) {
    status.value = 'ready'
  }
}

function clearFile() {
  file.value = null
  videoOriginPath.value = ''
  videoId.value = null
  uploadError.value = ''
  uploadedBytes.value = 0
  uploadingBytesTotal.value = 0
  hashingPercent.value = 0
  emit('uploaded', '')

  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }

  status.value = uploadId.value ? 'ready' : 'initializing'
}

function buildAllChunkIds() {
  return Array.from({ length: chunkCount.value }, (_, index) => index)
}

function getChunkBlob(chunkId) {
  const start = chunkId * chunkSizeBytes.value
  const end = Math.min(start + chunkSizeBytes.value, file.value.size)
  return file.value.slice(start, end)
}

async function calculateBlobMd5(blob) {
  const md5 = createMd5()
  const buffer = await readAsArrayBuffer(blob)

  md5.append(buffer)
  return md5.end()
}

async function getChunkUploadInfo(chunkId, md5) {
  const resp = await api.checkVideoChunk(uploadId.value, {
    md5,
    chunk_id: chunkId
  })

  return extractResponseData(resp, '分块上传地址获取失败')
}

function normalizeUploadHeaders(headers) {
  return Object.fromEntries(
    Object.entries(headers || {}).filter(([, value]) => value !== undefined && value !== null)
  )
}

async function uploadChunk(chunkId, chunkBlob, chunkInfo) {
  if (!chunkInfo?.upload_url) {
    throw new Error('后端未返回分块上传地址')
  }

  const response = await fetch(chunkInfo.upload_url, {
    method: chunkInfo.method || 'PUT',
    headers: normalizeUploadHeaders(chunkInfo.headers),
    body: chunkBlob
  })

  if (!response.ok) {
    throw new Error(`分块 ${chunkId + 1} 上传失败 (${response.status})`)
  }
}

async function processChunk(chunkId) {
  const chunkBlob = getChunkBlob(chunkId)
  const md5 = await calculateBlobMd5(chunkBlob)
  const chunkInfo = await getChunkUploadInfo(chunkId, md5)

  if (chunkInfo.exist !== true) {
    await uploadChunk(chunkId, chunkBlob, chunkInfo)
  }

  uploadedBytes.value += chunkBlob.size
}

async function uploadChunks(chunkIds) {
  let cursor = 0

  async function worker() {
    while (cursor < chunkIds.length) {
      const chunkId = chunkIds[cursor]
      cursor += 1
      await processChunk(chunkId)
    }
  }

  const workerCount = Math.min(concurrency, chunkIds.length)
  await Promise.all(Array.from({ length: workerCount }, worker))
}

function readAsArrayBuffer(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = () => resolve(reader.result)
    reader.onerror = () => reject(reader.error)
    reader.readAsArrayBuffer(blob)
  })
}

async function calculateFileMd5() {
  const md5 = createMd5()
  const totalChunks = chunkCount.value

  for (let chunkId = 0; chunkId < totalChunks; chunkId += 1) {
    const buffer = await readAsArrayBuffer(getChunkBlob(chunkId))
    md5.append(buffer)
    hashingPercent.value = Math.round(((chunkId + 1) / totalChunks) * 100)
  }

  return md5.end()
}

function extractVideoAsset(data) {
  if (typeof data === 'string') {
    return {
      origin_path: data,
      id: null
    }
  }

  return {
    origin_path: data?.origin_path || '',
    id: data?.id ?? data?.video_id ?? data?.videoId ?? null
  }
}

async function finishUpload(md5) {
  const resp = await api.finishVideoUpload(uploadId.value, {
    md5,
    file_name: file.value.name,
    chunk_count: chunkCount.value
  })
  const video = extractVideoAsset(extractResponseData(resp, '视频合并失败'))

  if (!video.origin_path) {
    throw new Error('后端未返回视频路径')
  }

  if (video.id === null || video.id === undefined || video.id === '') {
    throw new Error('Backend did not return video id')
  }

  return video
}

async function startUpload() {
  if (!file.value || !uploadId.value) {
    return
  }

  uploadError.value = ''
  videoOriginPath.value = ''
  videoId.value = null
  emit('uploaded', '')

  try {
    status.value = 'checking'
    const chunkIds = buildAllChunkIds()

    status.value = 'uploading'
    uploadedBytes.value = 0
    uploadingBytesTotal.value = file.value.size

    if (chunkIds.length > 0) {
      await uploadChunks(chunkIds)
    }

    status.value = 'hashing'
    hashingPercent.value = 0
    const md5 = await calculateFileMd5()

    status.value = 'merging'
    const video = await finishUpload(md5)
    videoOriginPath.value = video.origin_path
    videoId.value = video.id
    status.value = 'success'
    emit('uploaded', {
      origin_path: videoOriginPath.value,
      id: videoId.value
    })
  } catch (error) {
    uploadError.value = error?.response?.data?.message || error?.message || '视频上传失败，请重试。'
    status.value = 'error'
  }
}

function formatFileSize(bytes) {
  if (!bytes) {
    return '0B'
  }

  const units = ['B', 'KB', 'MB', 'GB']
  let size = bytes
  let unitIndex = 0

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex += 1
  }

  return `${size.toFixed(size >= 10 || unitIndex === 0 ? 0 : 1)}${units[unitIndex]}`
}

onMounted(createUploadTask)
</script>

<style scoped>
.video-dropzone {
  border: var(--bs-border-width) dashed var(--bs-border-color);
  border-radius: var(--bs-border-radius-lg);
  background-color: var(--bs-body-bg);
  cursor: pointer;
  transition: border-color .15s ease, background-color .15s ease;
}

.video-dropzone:hover,
.video-dropzone:focus-within {
  border-color: var(--bs-primary);
  background-color: var(--bs-primary-bg-subtle);
}
</style>
