<template>
  <section v-if="loading" class="row justify-content-center">
    <div class="col-lg-7">
      <div class="card border-0 shadow-sm text-center">
        <div class="card-body p-4 p-lg-5">
          <span class="spinner-border text-primary mb-3" aria-hidden="true"></span>
          <h2 class="h4 mb-2">视频加载中</h2>
          <p class="text-secondary mb-0">请稍候。</p>
        </div>
      </div>
    </div>
  </section>

  <section v-else-if="errorMessage" class="row justify-content-center">
    <div class="col-lg-7">
      <div class="card border-0 shadow-sm text-center">
        <div class="card-body p-4 p-lg-5">
          <i class="bi bi-exclamation-triangle display-5 text-warning" aria-hidden="true"></i>
          <h2 class="h4 mt-3 mb-2">视频加载失败</h2>
          <p class="text-secondary mb-4">{{ errorMessage }}</p>
          <button class="btn btn-primary" type="button" @click="loadVideo">
            重试
          </button>
        </div>
      </div>
    </div>
  </section>

  <section v-else-if="!currentSource" class="row justify-content-center">
    <div class="col-lg-7">
      <div class="card border-0 shadow-sm text-center">
        <div class="card-body p-4 p-lg-5">
          <i class="bi bi-camera-video-off display-5 text-secondary" aria-hidden="true"></i>
          <h2 class="h4 mt-3 mb-2">本章节暂无视频</h2>
          <p class="text-secondary mb-4">后端没有返回可播放的视频路径。</p>
          <RouterLink class="btn btn-primary" :to="{ name: 'CourseDetail', params: { course_id: courseId }, hash: '#content' }">
            返回课程内容
          </RouterLink>
        </div>
      </div>
    </div>
  </section>

  <section v-else class="card border-0 shadow-sm overflow-hidden">
    <div class="ratio ratio-16x9 bg-dark">
      <video
        v-if="currentSource"
        ref="videoRef"
        class="course-video-media"
        :poster="video.poster || ''"
        controls
        playsinline
        preload="metadata"
      >
        当前浏览器不支持视频播放。
      </video>
      <div v-else class="d-flex flex-column align-items-center justify-content-center gap-2 text-light">
        <i class="bi bi-play-circle display-5" aria-hidden="true"></i>
        <span>暂无可播放视频</span>
      </div>
    </div>

    <div class="card-body p-3 p-lg-4">
      <div class="d-flex flex-column flex-lg-row justify-content-between gap-3">
        <div>
          <span class="badge text-bg-primary mb-2">课程视频</span>
          <h1 class="h4 fw-bold mb-2">{{ video.title || '章节视频' }}</h1>
          <p class="text-secondary mb-0">
            {{ video.description || '本章节视频' }}
          </p>
        </div>

        <div class="d-flex flex-wrap align-content-start gap-2">
          <span v-if="video.id" class="badge text-bg-light border">ID {{ video.id }}</span>
          <span v-if="durationText" class="badge text-bg-light border">
            <i class="bi bi-clock me-1" aria-hidden="true"></i>
            {{ durationText }}
          </span>
          <span v-if="sizeText" class="badge text-bg-light border">
            <i class="bi bi-hdd me-1" aria-hidden="true"></i>
            {{ sizeText }}
          </span>
        </div>
      </div>

      <div
        v-if="qualityOptions.length > 1"
        class="d-flex flex-column flex-sm-row align-items-sm-center gap-2 mt-3"
      >
        <span class="small text-secondary">清晰度</span>
        <div class="btn-group btn-group-sm" role="group" aria-label="视频清晰度">
          <button
            v-for="option in qualityOptions"
            :key="option.value"
            class="btn"
            :class="selectedQuality === option.value ? 'btn-primary' : 'btn-outline-primary'"
            type="button"
            :disabled="selectedQuality === option.value"
            @click="selectQuality(option.value)"
          >
            {{ option.label }}
          </button>
        </div>
      </div>

      <div v-if="currentSource" class="border rounded-3 bg-body-tertiary p-3 mt-3">
        <div class="small text-secondary mb-1">当前播放路径</div>
        <div class="text-break small">{{ currentSource }}</div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import api from '@/api'

const props = defineProps({
  courseId: {
    type: [String, Number],
    required: true
  },
  chapterId: {
    type: [String, Number],
    required: true
  },
  videoId: {
    type: [String, Number],
    required: true
  }
})

const loading = ref(false)
const errorMessage = ref('')
const video = ref(createEmptyVideo())
const videoRef = ref(null)
const selectedQuality = ref('')
const apiBaseUrl = process.env.VUE_APP_API_BASE || 'http://localhost:8080'
let HlsConstructor = null
let hlsPlayer = null
let sourceRequestId = 0

const durationText = computed(() => formatDuration(video.value.duration))
const qualityOptions = computed(() => video.value.profiles)
const currentProfile = computed(() => (
  qualityOptions.value.find((profile) => profile.value === selectedQuality.value) || null
))
const currentSource = computed(() => currentProfile.value?.source || video.value.source)
const sizeText = computed(() => formatFileSize(currentProfile.value?.fileSize || video.value.sizeBytes))

function createEmptyVideo() {
  return {
    id: '',
    title: '',
    description: '',
    origin_path: '',
    source: '',
    poster: '',
    profiles: [],
    duration: 0,
    sizeBytes: 0
  }
}

function getAppOrigin() {
  return window.location?.origin || 'http://localhost:8080'
}

function getApiOrigin() {
  if (/^(https?:)?\/\//.test(apiBaseUrl)) {
    return new URL(apiBaseUrl).origin
  }

  return getAppOrigin()
}

function normalizeOriginPath(originPath) {
  if (!originPath) {
    return ''
  }

  const normalizedOriginPath = String(originPath).trim().replace(/\\/g, '/')

  if (
    /^(https?:)?\/\//.test(normalizedOriginPath) ||
    normalizedOriginPath.startsWith('data:') ||
    normalizedOriginPath.startsWith('blob:')
  ) {
    return normalizedOriginPath
  }

  if (normalizedOriginPath.startsWith('/')) {
    return new URL(normalizedOriginPath, getApiOrigin()).href
  }

  return new URL(`/${normalizedOriginPath}`, getApiOrigin()).href
}

function getResponseData(resp) {
  const body = resp?.data ?? resp

  if (body?.code === 1) {
    throw new Error(body.message || '视频信息加载失败')
  }

  if (body && typeof body === 'object' && 'data' in body && ('code' in body || 'message' in body)) {
    return body.data
  }

  return body
}

function getVideoPayload(data) {
  if (typeof data === 'string') {
    return { origin_path: data }
  }

  if (data?.videoAsset && typeof data.videoAsset === 'object') {
    return { ...data, ...data.videoAsset }
  }

  if (data?.video_asset && typeof data.video_asset === 'object') {
    return { ...data, ...data.video_asset }
  }

  if (data?.video && typeof data.video === 'object') {
    return { ...data, ...data.video }
  }

  return data || {}
}

function getProfileItems(payload) {
  const profileData = payload.profiles || payload.video_profiles || payload.videoProfiles || payload.qualities

  if (Array.isArray(profileData)) {
    return profileData
  }

  if (profileData && typeof profileData === 'object') {
    return Object.entries(profileData).map(([resolution, value]) => (
      typeof value === 'string'
        ? { resolution, playlist_path: value }
        : { resolution, ...value }
    ))
  }

  return []
}

function normalizeResolution(value) {
  const text = String(value || '').trim().toUpperCase()
  const match = text.match(/(480|720|1080)/)

  return match ? `${match[1]}P` : text
}

function formatQualityLabel(profile) {
  const resolution = normalizeResolution(profile.resolution || profile.quality || profile.label || profile.height)
  const match = resolution.match(/(480|720|1080)P/)

  return match ? `${match[1]}p` : resolution
}

function qualityRank(value) {
  const resolution = normalizeResolution(value)

  if (resolution === '480P') return 1
  if (resolution === '720P') return 2
  if (resolution === '1080P') return 3

  return 0
}

function normalizeProfiles(payload) {
  const uniqueProfiles = new Map()

  getProfileItems(payload).forEach((profile) => {
    const sourcePath = (
      profile.playlist_path ||
      profile.playlistPath ||
      profile.play_path ||
      profile.playPath ||
      profile.source ||
      profile.url ||
      profile.path ||
      ''
    )
    const source = normalizeOriginPath(sourcePath)

    if (!source) {
      return
    }

    const resolution = normalizeResolution(profile.resolution || profile.quality || profile.label || profile.height)
    const value = resolution || source

    uniqueProfiles.set(value, {
      value,
      label: formatQualityLabel(profile),
      source,
      fileSize: Number(profile.file_size ?? profile.fileSize ?? profile.size_bytes ?? profile.sizeBytes ?? 0),
      bitrate: Number(profile.bitrate || 0),
      width: Number(profile.width || 0),
      height: Number(profile.height || 0)
    })
  })

  return [...uniqueProfiles.values()].sort((left, right) => (
    qualityRank(left.value) - qualityRank(right.value)
  ))
}

function getDefaultQuality(profiles) {
  if (profiles.length === 0) {
    return ''
  }

  return profiles.find((profile) => profile.value === '720P')?.value || profiles[profiles.length - 1].value
}

function normalizeVideo(data) {
  const payload = getVideoPayload(data)
  const originPath = payload.origin_path || payload.originPath || ''
  const playPath = payload.play_path || payload.playPath || ''
  const profiles = normalizeProfiles(payload)

  return {
    id: payload.id ?? payload.video_id ?? payload.videoId ?? '',
    title: payload.title || payload.name || `第 ${props.chapterId} 章视频`,
    description: payload.description || '',
    origin_path: originPath,
    source: normalizeOriginPath(playPath || originPath),
    poster: normalizeOriginPath(payload.poster || payload.coverPath || payload.cover_path || payload.coverUrl || payload.cover_url || ''),
    profiles,
    duration: Number(payload.duration ?? payload.durationSeconds ?? payload.duration_seconds ?? 0),
    sizeBytes: Number(payload.sizeBytes ?? payload.size_bytes ?? payload.size ?? 0)
  }
}

function destroyHlsPlayer() {
  if (hlsPlayer) {
    hlsPlayer.destroy()
    hlsPlayer = null
  }
}

function isHlsSource(source) {
  return /\.m3u8($|\?)/i.test(source)
}

async function getHlsConstructor() {
  if (!HlsConstructor) {
    const hlsModule = await import('hls.js')

    HlsConstructor = hlsModule.default
  }

  return HlsConstructor
}

function restorePlayback(videoElement, currentTime, shouldPlay) {
  const restore = () => {
    if (currentTime > 0) {
      videoElement.currentTime = currentTime
    }

    if (shouldPlay) {
      const playResult = videoElement.play()

      if (playResult?.catch) {
        playResult.catch(() => {})
      }
    }
  }

  if (videoElement.readyState >= 1) {
    restore()
    return
  }

  videoElement.addEventListener('loadedmetadata', restore, { once: true })
}

async function applyVideoSource() {
  const requestId = sourceRequestId + 1
  sourceRequestId = requestId

  await nextTick()

  const videoElement = videoRef.value
  const source = currentSource.value

  destroyHlsPlayer()

  if (!videoElement || !source) {
    return
  }

  const currentTime = videoElement.currentTime || 0
  const shouldPlay = !videoElement.paused

  if (isHlsSource(source)) {
    const Hls = await getHlsConstructor()

    if (requestId !== sourceRequestId || currentSource.value !== source) {
      return
    }

    if (Hls.isSupported()) {
      hlsPlayer = new Hls()
      hlsPlayer.loadSource(source)
      hlsPlayer.attachMedia(videoElement)
      hlsPlayer.on(Hls.Events.MANIFEST_PARSED, () => {
        restorePlayback(videoElement, currentTime, shouldPlay)
      })
      hlsPlayer.on(Hls.Events.ERROR, (_, data) => {
        if (data?.fatal) {
          destroyHlsPlayer()
          errorMessage.value = '视频播放失败，请切换清晰度或稍后重试。'
        }
      })
      return
    }
  }

  videoElement.src = source
  videoElement.load()
  restorePlayback(videoElement, currentTime, shouldPlay)
}

function selectQuality(quality) {
  selectedQuality.value = quality
}

async function loadVideo() {
  if (!props.courseId || !props.chapterId || !props.videoId) {
    video.value = createEmptyVideo()
    errorMessage.value = '课程、章节或视频参数缺失。'
    return
  }

  let shouldApplySource = false
  loading.value = true
  errorMessage.value = ''

  try {
    const resp = await api.getCourseChapterVideo(props.courseId, props.chapterId, props.videoId)
    const nextVideo = normalizeVideo(getResponseData(resp))

    selectedQuality.value = getDefaultQuality(nextVideo.profiles)
    video.value = nextVideo
    shouldApplySource = Boolean(currentSource.value)
  } catch (error) {
    video.value = createEmptyVideo()
    selectedQuality.value = ''
    errorMessage.value = error?.response?.data?.message || error?.message || '视频信息加载失败，请稍后重试。'
  } finally {
    loading.value = false
    if (shouldApplySource) {
      await applyVideoSource()
    }
  }
}

function formatDuration(seconds) {
  const totalSeconds = Number(seconds || 0)

  if (!Number.isFinite(totalSeconds) || totalSeconds <= 0) {
    return ''
  }

  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const restSeconds = Math.floor(totalSeconds % 60)
  const paddedMinutes = hours > 0 ? String(minutes).padStart(2, '0') : String(minutes)
  const paddedSeconds = String(restSeconds).padStart(2, '0')

  return hours > 0
    ? `${hours}:${paddedMinutes}:${paddedSeconds}`
    : `${paddedMinutes}:${paddedSeconds}`
}

function formatFileSize(bytes) {
  const value = Number(bytes || 0)

  if (!Number.isFinite(value) || value <= 0) {
    return ''
  }

  const units = ['B', 'KB', 'MB', 'GB']
  let size = value
  let unitIndex = 0

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex += 1
  }

  return `${size.toFixed(size >= 10 || unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`
}

onMounted(loadVideo)
watch(() => [props.courseId, props.chapterId, props.videoId], loadVideo)
watch(currentSource, () => {
  if (!loading.value) {
    applyVideoSource()
  }
}, { flush: 'post' })
onBeforeUnmount(destroyHlsPlayer)
</script>

<style scoped>
.course-video-media {
  width: 100%;
  height: 100%;
  object-fit: contain;
  background-color: var(--bs-dark);
}
</style>
