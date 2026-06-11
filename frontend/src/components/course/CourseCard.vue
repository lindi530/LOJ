<template>
  <RouterLink
    class="card border-0 shadow-sm text-reset text-decoration-none h-100 w-100"
    :to="{ name: 'CourseDetail', params: { course_id: course.id }, hash: '#intro' }"
  >
    <div class="ratio ratio-16x9 bg-body-tertiary rounded-top overflow-hidden">
      <img
        v-if="course.coverUrl"
        :src="course.coverUrl"
        class="w-100 h-100 object-fit-cover"
        :alt="course.title"
      >
      <div
        v-else
        class="d-flex align-items-center justify-content-center bg-primary-subtle text-primary-emphasis"
        aria-hidden="true"
      >
        <i class="bi bi-mortarboard-fill display-5"></i>
      </div>
    </div>

    <div class="card-body d-flex flex-column gap-3">
      <div class="d-flex justify-content-between gap-3 align-items-start">
        <div>
          <h2 class="h5 card-title fw-bold mb-2">{{ course.title }}</h2>
        </div>
        <span
          class="badge"
          :class="course.isFree ? 'text-bg-success' : 'text-bg-warning'"
        >
          {{ priceText }}
        </span>
      </div>

      <p class="card-text text-secondary text-break mb-0">{{ course.description }}</p>

      <div class="row g-2 mt-auto text-secondary small">
        <div class="col-6">
          <div class="border rounded-3 p-2 bg-body-tertiary">
            <i class="bi bi-play-circle me-1" aria-hidden="true"></i>
            {{ course.lessonCount }} 节课
          </div>
        </div>
        <div class="col-6">
          <div class="border rounded-3 p-2 bg-body-tertiary">
            <i class="bi bi-people me-1" aria-hidden="true"></i>
            {{ course.studentCount }} 人
          </div>
        </div>
      </div>
    </div>

    <div class="card-footer bg-body border-0 pt-0">
      <div class="d-grid">
        <div class="btn btn-outline-primary">
          查看课程
          <i class="bi bi-arrow-right ms-1" aria-hidden="true"></i>
        </div>
      </div>
    </div>
  </RouterLink>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  course: {
    type: Object,
    required: true
  }
})

const priceText = computed(() => {
  if (props.course.isFree) {
    return '免费'
  }

  return `￥${Number(props.course.price || 0).toFixed(2)}`
})
</script>
