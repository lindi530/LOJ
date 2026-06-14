<template>
  <main class="py-4 py-lg-5">
    <div class="container-xl">
      <div class="d-flex flex-column flex-lg-row justify-content-between gap-3 align-items-lg-center mb-4">
        <div>
          <RouterLink
            class="btn btn-link link-secondary text-decoration-none px-0 mb-2"
            :to="{ name: 'CourseDetail', params: { course_id: courseId }, hash: '#content' }"
          >
            <i class="bi bi-arrow-left me-1" aria-hidden="true"></i>
            返回课程内容
          </RouterLink>
          <h1 class="h3 fw-bold mb-1">课程视频播放</h1>
          <p class="text-secondary mb-0">课程 {{ courseId }} / 章节 {{ chapterId }}</p>
        </div>

        <button
          class="btn btn-outline-primary align-self-start align-self-lg-center"
          type="button"
          :disabled="loading"
          @click="loadVideo"
        >
          <span v-if="loading" class="spinner-border spinner-border-sm me-2" aria-hidden="true"></span>
          <i v-else class="bi bi-arrow-clockwise me-2" aria-hidden="true"></i>
          重新加载
        </button>
      </div>

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

      <section v-else-if="!video.url" class="row justify-content-center">
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

      <CourseVideoPlayer v-else :video="video" />
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import CourseVideoPlayer from './CourseVideoPlayer.vue'

const route = useRoute()
const loading = ref(false)
const errorMessage = ref('')
const video = ref(createEmptyVideo())
const apiBaseUrl = process.env.VUE_APP_API_BASE || 'http://localhost:8080'

const courseId = computed(() => route.params.course_id)
const chapterId = computed(() => route.params.chapter_id)
const videoId = computed(() => route.params.video_id)

function createEmptyVideo() {
  return {
    id: '',
    title: '',
    description: '',
    url: '',
    poster: '',
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

function normalizeVideoUrl(url) {
  if (!url) {
    return ''
  }

  const normalizedUrl = String(url).trim().replace(/\\/g, '/')

  if (
    /^(https?:)?\/\//.test(normalizedUrl) ||
    normalizedUrl.startsWith('data:') ||
    normalizedUrl.startsWith('blob:')
  ) {
    return normalizedUrl
  }

  if (normalizedUrl.startsWith('/')) {
    return new URL(normalizedUrl, getApiOrigin()).href
  }

  return new URL(`/${normalizedUrl}`, getApiOrigin()).href
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
    return { url: data }
  }

  return data?.videoAsset || data?.video_asset || data?.video || data || {}
}

function normalizeVideo(data) {
  const payload = getVideoPayload(data)

  return {
    id: payload.id ?? payload.video_id ?? payload.videoId ?? '',
    title: payload.title || payload.name || `第 ${chapterId.value} 章视频`,
    description: payload.description || '',
    url: normalizeVideoUrl(payload.url || payload.video_url || payload.videoUrl || payload.src || payload.path || ''),
    poster: normalizeVideoUrl(payload.poster || payload.coverUrl || payload.cover_url || ''),
    duration: Number(payload.duration ?? payload.durationSeconds ?? payload.duration_seconds ?? 0),
    sizeBytes: Number(payload.sizeBytes ?? payload.size_bytes ?? payload.size ?? 0)
  }
}

async function loadVideo() {
  if (!courseId.value || !chapterId.value || !videoId.value) {
    video.value = createEmptyVideo()
    errorMessage.value = '课程、章节或视频参数缺失。'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const resp = await api.getCourseChapterVideo(courseId.value, chapterId.value, videoId.value)
    video.value = normalizeVideo(getResponseData(resp))
  } catch (error) {
    video.value = createEmptyVideo()
    errorMessage.value = error?.response?.data?.message || error?.message || '视频信息加载失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

onMounted(loadVideo)
watch(() => [courseId.value, chapterId.value, videoId.value], loadVideo)
</script>
