<template>
  <div class="row g-4 align-items-start">
    <div class="col-lg-9">
      <div class="mb-0">
        <h2 class="h4 fw-bold mb-3">课程简介</h2>
        <p class="text-secondary mb-0">{{ course.description }}</p>
      </div>
    </div>

    <aside class="col-lg-3">
      <div class="list-group shadow-sm small">
        <div class="list-group-item d-flex justify-content-between gap-3">
          <span class="text-secondary">讲师</span>
          <strong>{{ course.teacher || 'LOJ' }}</strong>
        </div>
        <div class="list-group-item d-flex justify-content-between gap-3">
          <span class="text-secondary">难度</span>
          <strong>{{ course.level || '入门' }}</strong>
        </div>
        <div class="list-group-item d-flex justify-content-between gap-3">
          <span class="text-secondary">周期</span>
          <strong>{{ course.duration || '长期有效' }}</strong>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
defineProps({
  course: {
    type: Object,
    required: true
  }
})
</script>
