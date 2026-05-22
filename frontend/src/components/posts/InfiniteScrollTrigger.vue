<template>
  <div ref="triggerRef" class="load-trigger" aria-live="polite">
    <span v-if="loadingVisible" class="load-message">
      <i class="bi bi-lightning-charge-fill load-icon"></i>
      正在玩命加载中
    </span>
    <span v-else-if="hasMore" class="load-message">继续向下滚动加载更多</span>
    <span v-else class="load-message">没有更多帖子了</span>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  hasMore: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['load-more'])
const triggerRef = ref(null)
const loadingVisible = ref(false)
let observer = null
let hideTimer = null
let emitTimer = null
let loadingStartedAt = 0

function requestLoadMore() {
  if (props.loading || loadingVisible.value || !props.hasMore) {
    return
  }

  showLoading()
  emitTimer = window.setTimeout(() => {
    emitTimer = null
    if (!props.hasMore) {
      hideLoadingAfterMinimum()
      return
    }

    emit('load-more')
    nextTick(() => {
      if (!props.loading) {
        hideLoadingAfterMinimum()
      }
    })
  }, 1000)
}

function showLoading() {
  clearHideTimer()
  loadingStartedAt = Date.now()
  loadingVisible.value = true
}

function hideLoadingAfterMinimum() {
  clearHideTimer()
  const remaining = Math.max(1000 - (Date.now() - loadingStartedAt), 0)
  hideTimer = window.setTimeout(() => {
    if (!props.loading) {
      loadingVisible.value = false
    }
  }, remaining)
}

function clearHideTimer() {
  if (hideTimer) {
    window.clearTimeout(hideTimer)
    hideTimer = null
  }
}

function clearEmitTimer() {
  if (emitTimer) {
    window.clearTimeout(emitTimer)
    emitTimer = null
  }
}

function isTriggerNearViewport() {
  if (!triggerRef.value) {
    return false
  }

  const rect = triggerRef.value.getBoundingClientRect()
  return rect.top <= window.innerHeight + 160
}

function handleScrollFallback() {
  if (isTriggerNearViewport()) {
    requestLoadMore()
  }
}

function startObserver() {
  if (!triggerRef.value || !('IntersectionObserver' in window)) {
    window.addEventListener('scroll', handleScrollFallback, { passive: true })
    handleScrollFallback()
    return
  }

  observer = new IntersectionObserver(entries => {
    if (entries.some(entry => entry.isIntersecting)) {
      requestLoadMore()
    }
  }, {
    rootMargin: '160px 0px',
    threshold: 0.01
  })

  observer.observe(triggerRef.value)
}

onMounted(() => {
  startObserver()
})

onBeforeUnmount(() => {
  if (observer) {
    observer.disconnect()
  }
  clearHideTimer()
  clearEmitTimer()
  window.removeEventListener('scroll', handleScrollFallback)
})

watch(
  () => props.loading,
  loading => {
    if (loading) {
      if (!loadingVisible.value) {
        showLoading()
      }
      return
    }

    if (loadingVisible.value && !emitTimer) {
      hideLoadingAfterMinimum()
    }
  }
)

watch(
  () => props.hasMore,
  hasMore => {
    if (!hasMore && loadingVisible.value && !props.loading) {
      hideLoadingAfterMinimum()
    }
  }
)
</script>

<style scoped>
.load-trigger {
  display: grid;
  min-height: 58px;
  place-items: center;
  padding: 14px 16px;
  border-top: 1px solid #edf0f3;
  background: linear-gradient(90deg, #fbfcfd, #f1fbf7, #fbfcfd);
  color: #7a8797;
  font-size: 14px;
}

.load-message {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 24px;
}

.load-icon {
  color: #12a875;
  filter: drop-shadow(0 0 6px rgba(18, 168, 117, 0.35));
  animation: load-pulse 0.7s ease-in-out infinite alternate;
}

@keyframes load-pulse {
  from {
    opacity: 0.55;
    transform: translateY(1px) scale(0.92);
  }

  to {
    opacity: 1;
    transform: translateY(-1px) scale(1.08);
  }
}
</style>
