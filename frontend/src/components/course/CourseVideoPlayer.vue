<template>
  <section class="card border-0 shadow-sm overflow-hidden">
    <div class="ratio ratio-16x9 bg-dark">
      <video
        v-if="video.source"
        class="course-video-media"
        :src="video.source"
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

      <div v-if="video.origin_path" class="border rounded-3 bg-body-tertiary p-3 mt-3">
        <div class="small text-secondary mb-1">播放路径</div>
        <div class="text-break small">{{ video.origin_path }}</div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  video: {
    type: Object,
    required: true
  }
})

const durationText = computed(() => formatDuration(props.video.duration))
const sizeText = computed(() => formatFileSize(props.video.sizeBytes))

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
</script>

<style scoped>
.course-video-media {
  width: 100%;
  height: 100%;
  object-fit: contain;
  background-color: var(--bs-dark);
}
</style>
