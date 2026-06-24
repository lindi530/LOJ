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
        </div>

      </div>

      <CourseVideoPlayer
        :course-id="courseId"
        :chapter-id="chapterId"
        :video-id="videoId"
      />
    </div>
  </main>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import CourseVideoPlayer from './CourseVideoPlayer.vue'

const route = useRoute()

const courseId = computed(() => route.params.course_id)
const chapterId = computed(() => route.params.chapter_id)
const videoId = computed(() => route.params.video_id)
</script>
