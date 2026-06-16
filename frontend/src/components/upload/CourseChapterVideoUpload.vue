<template>
  <section class="card border-0 shadow-sm mb-4">
    <div class="card-body p-4">
      <VideoChunkUpload
        @uploaded="handleVideoUploaded"
        @uploading-change="handleUploadingChange"
      />
    </div>
  </section>
</template>

<script setup>
import VideoChunkUpload from './VideoChunkUpload.vue'

defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      origin_path: '',
      id: null
    })
  },
  uploading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'update:uploading', 'uploaded', 'uploading-change'])

function normalizeVideoAsset(video) {
  if (typeof video === 'string') {
    return {
      origin_path: video,
      id: null
    }
  }

  return {
    origin_path: video?.origin_path || '',
    id: video?.id ?? video?.video_id ?? video?.videoId ?? null
  }
}

function handleVideoUploaded(video) {
  const videoAsset = normalizeVideoAsset(video)

  emit('update:modelValue', videoAsset)
  emit('uploaded', videoAsset)
}

function handleUploadingChange(uploading) {
  emit('update:uploading', uploading)
  emit('uploading-change', uploading)
}
</script>
