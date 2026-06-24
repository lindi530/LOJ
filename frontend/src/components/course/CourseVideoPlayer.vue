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
          <p class="text-secondary mb-4">后端没有返回可播放的视频地址。</p>
          <RouterLink class="btn btn-primary" :to="{ name: 'CourseDetail', params: { course_id: courseId }, hash: '#content' }">
            返回课程内容
          </RouterLink>
        </div>
      </div>
    </div>
  </section>

  <section v-else class="card border-0 shadow-sm overflow-hidden">
    <div class="ratio ratio-16x9 bg-dark">
      <div>
        <video
          ref="videoRef"
          class="w-100 h-100 object-fit-contain bg-dark"
          :poster="video.poster || ''"
          controls
          playsinline
          preload="metadata"
        >
          当前浏览器不支持视频播放。
        </video>

        <div
          v-if="qualityOptions.length > 1"
          class="position-absolute top-0 end-0 z-1 m-2 m-sm-3"
          @click.stop
        >
          <select
            v-model="selectedQuality"
            class="form-select form-select-sm w-auto bg-dark text-light border-secondary shadow-sm"
            aria-label="视频清晰度"
          >
            <option
              v-for="option in qualityOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </option>
          </select>
        </div>
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

const qualityOrder = {
  '480P': 1,
  '720P': 2,
  '1080P': 3
}

const loading = ref(false)
const errorMessage = ref('')
const video = ref(createEmptyVideo())
const videoRef = ref(null)
const selectedQuality = ref('')
const apiBaseUrl = process.env.VUE_APP_API_BASE || 'http://localhost:8080'

let hlsPlayer = null
let HlsConstructor = null
let loadRequestId = 0
let sourceRequestId = 0

const qualityOptions = computed(() => video.value.profiles)
const currentProfile = computed(() => (
  qualityOptions.value.find((profile) => profile.value === selectedQuality.value) || null
))
const currentSource = computed(() => currentProfile.value?.source || '')

function createEmptyVideo() {
  return {
    poster: '',
    profiles: []
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

function resolveMediaUrl(mediaUrl) {
  const normalizedMediaUrl = String(mediaUrl || '').trim()

  if (!normalizedMediaUrl) {
    return ''
  }

  if (
    /^(https?:)?\/\//.test(normalizedMediaUrl) ||
    normalizedMediaUrl.startsWith('data:') ||
    normalizedMediaUrl.startsWith('blob:')
  ) {
    return normalizedMediaUrl
  }

  return normalizedMediaUrl.startsWith('/')
    ? new URL(normalizedMediaUrl, getApiOrigin()).href
    : ''
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

function normalizeResolution(value) {
  const text = String(value || '').trim().toUpperCase()
  const match = text.match(/(480|720|1080)/)

  return match ? `${match[1]}P` : text
}

function formatQualityLabel(resolution) {
  const match = resolution.match(/(480|720|1080)P/)

  return match ? `${match[1]}p` : resolution || '默认'
}

function normalizeProfiles(payload) {
  const profiles = Array.isArray(payload?.profiles) ? payload.profiles : []
  const uniqueProfiles = new Map()

  profiles.forEach((profile) => {
    const source = resolveMediaUrl(profile.play_url)

    if (!source) {
      return
    }

    const resolution = normalizeResolution(profile.resolution)
    const value = resolution || source

    uniqueProfiles.set(value, {
      value,
      source,
      label: formatQualityLabel(resolution)
    })
  })

  return [...uniqueProfiles.values()].sort((left, right) => (
    (qualityOrder[left.value] || 0) - (qualityOrder[right.value] || 0)
  ))
}

function getDefaultQuality(profiles) {
  if (profiles.length === 0) {
    return ''
  }

  return profiles.find((profile) => profile.value === '720P')?.value || profiles[profiles.length - 1].value
}

function normalizeVideo(data) {
  const payload = data && typeof data === 'object' ? data : {}

  return {
    poster: resolveMediaUrl(payload.cover_url),
    profiles: normalizeProfiles(payload)
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

function isHlsKeyRequest(requestUrl) {
  const rawUrl = String(requestUrl || '').trim()

  if (!rawUrl) {
    return false
  }

  try {
    const { pathname } = new URL(rawUrl, getAppOrigin())

    return pathname.includes('/api/course/video/') && pathname.endsWith('/hls-key')
  } catch {
    return rawUrl.includes('/api/course/video/') && rawUrl.includes('/hls-key')
  }
}

function attachHlsKeyAuthHeader(xhr, requestUrl) {
  if (!isHlsKeyRequest(requestUrl)) {
    return
  }

  const accessToken = localStorage.getItem('accessToken')
  if (accessToken) {
    xhr.setRequestHeader('Authorization', `Bearer ${accessToken}`)
  }
}

async function getHlsConstructor() {
  if (!HlsConstructor) {
    const hlsModule = await import('hls.js')

    HlsConstructor = hlsModule.default
  }

  return HlsConstructor
}

function playSafely(videoElement) {
  const playResult = videoElement.play()

  if (playResult?.catch) {
    playResult.catch(() => {})
  }
}

function restorePlayback(videoElement, playbackTime, shouldPlay) {
  const restore = () => {
    if (playbackTime > 0) {
      videoElement.currentTime = playbackTime
    }

    if (shouldPlay) {
      playSafely(videoElement)
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

  const playbackTime = Number(videoElement.currentTime || 0)
  const shouldPlay = !videoElement.paused && !videoElement.ended

  if (isHlsSource(source)) {
    const Hls = await getHlsConstructor()

    if (requestId !== sourceRequestId || currentSource.value !== source) {
      return
    }

    if (Hls.isSupported()) {
      hlsPlayer = new Hls({
        xhrSetup(xhr, url) {
          xhr.withCredentials = true
          attachHlsKeyAuthHeader(xhr, url)
        }
      })
      hlsPlayer.loadSource(source)
      hlsPlayer.attachMedia(videoElement)
      hlsPlayer.on(Hls.Events.MANIFEST_PARSED, () => {
        restorePlayback(videoElement, playbackTime, shouldPlay)
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
  restorePlayback(videoElement, playbackTime, shouldPlay)
}

function resetPlayer(nextErrorMessage = '') {
  destroyHlsPlayer()
  video.value = createEmptyVideo()
  selectedQuality.value = ''
  errorMessage.value = nextErrorMessage
}

async function loadVideo() {
  const requestId = loadRequestId + 1
  loadRequestId = requestId

  if (!props.courseId || !props.chapterId || !props.videoId) {
    loading.value = false
    resetPlayer('课程、章节或视频参数缺失。')
    return
  }

  loading.value = true
  errorMessage.value = ''
  destroyHlsPlayer()

  try {
    const resp = await api.getCourseChapterVideo(props.courseId, props.chapterId, props.videoId)

    if (requestId !== loadRequestId) {
      return
    }

    const nextVideo = normalizeVideo(getResponseData(resp))

    video.value = nextVideo
    selectedQuality.value = getDefaultQuality(nextVideo.profiles)
  } catch (error) {
    if (requestId !== loadRequestId) {
      return
    }

    resetPlayer(error?.response?.data?.message || error?.message || '视频信息加载失败，请稍后重试。')
  } finally {
    if (requestId === loadRequestId) {
      loading.value = false
      await applyVideoSource()
    }
  }
}

onMounted(loadVideo)
watch(() => [props.courseId, props.chapterId, props.videoId], loadVideo)
watch(selectedQuality, () => {
  if (!loading.value) {
    applyVideoSource()
  }
}, { flush: 'post' })
onBeforeUnmount(destroyHlsPlayer)
</script>
