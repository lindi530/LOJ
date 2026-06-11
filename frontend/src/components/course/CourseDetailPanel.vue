<template>
  <section>
    <RouterLink
      class="btn btn-link link-secondary text-decoration-none px-0 mb-3"
      :to="{ name: 'Course' }"
    >
      <i class="bi bi-arrow-left me-1" aria-hidden="true"></i>
      返回课程列表
    </RouterLink>

    <div class="card border-0 shadow-sm mb-4">
      <div class="card-body p-3 p-lg-5">
        <div class="row g-4 align-items-center">
          <div class="col-lg-5">
            <div class="ratio ratio-16x9 rounded-3 overflow-hidden bg-body-tertiary">
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
                <i class="bi bi-mortarboard-fill display-4"></i>
              </div>
            </div>
          </div>

          <div class="col-lg-7">
            <div class="d-flex flex-wrap gap-2 mb-3">
              <span class="badge course-meta-badge rounded-pill text-bg-primary">{{ course.level || '课程' }}</span>
              <span class="badge course-meta-badge rounded-pill text-bg-light border">{{ course.duration || '长期有效' }}</span>
              <span
                class="badge course-meta-badge rounded-pill"
                :class="course.isFree ? 'text-bg-success' : 'text-bg-warning'"
              >
                {{ priceText }}
              </span>
            </div>

            <h1 class="display-6 fw-bold mb-3">{{ course.title }}</h1>
            <p class="lead text-secondary text-break mb-4 course-description-clamp">{{ course.description }}</p>

            <div class="row row-cols-1 row-cols-sm-2 row-cols-md-4 g-3 mb-4">
              <div class="col">
                <div class="border rounded-3 p-3 bg-body-tertiary">
                  <div class="fw-bold">{{ course.lessonCount }}</div>
                  <div class="small text-secondary">课时</div>
                </div>
              </div>
              <div class="col">
                <div class="border rounded-3 p-3 bg-body-tertiary">
                  <div class="fw-bold">{{ course.studentCount }}</div>
                  <div class="small text-secondary">学生</div>
                </div>
              </div>
              <div class="col">
                <div class="border rounded-3 p-3 bg-body-tertiary">
                  <div class="fw-bold">{{ course.chapters?.length || 0 }}</div>
                  <div class="small text-secondary">章节</div>
                </div>
              </div>
              <div class="col">
                <div class="border rounded-3 p-3 bg-body-tertiary">
                  <div class="fw-bold">永久</div>
                  <div class="small text-secondary">有效期</div>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>
    </div>

    <div class="card border-0 shadow-sm">
      <div class="card-header bg-body border-0 p-3 p-lg-4 pb-0">
        <ul class="nav nav-pills gap-2">
          <li class="nav-item">
            <RouterLink
              class="nav-link border"
              :class="activeTab === 'intro' ? 'active' : 'text-body bg-body'"
              :to="{ name: 'CourseDetail', params: { course_id: course.id }, hash: '#intro' }"
            >
              <i class="bi bi-info-circle me-1" aria-hidden="true"></i>
              课程介绍
            </RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink
              class="nav-link border"
              :class="activeTab === 'content' ? 'active' : 'text-body bg-body'"
              :to="{ name: 'CourseDetail', params: { course_id: course.id }, hash: '#content' }"
            >
              <i class="bi bi-list-task me-1" aria-hidden="true"></i>
              课程内容
            </RouterLink>
          </li>
        </ul>
      </div>

      <div class="card-body p-3 p-lg-4">
        <CourseIntro v-if="activeTab === 'intro'" :course="course" />
        <CourseContent
          v-else
          :course-id="course.id"
          :chapters="course.chapters || []"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import CourseContent from './CourseContent.vue'
import CourseIntro from './CourseIntro.vue'

const props = defineProps({
  course: {
    type: Object,
    required: true
  },
  activeTab: {
    type: String,
    default: 'intro'
  }
})

const priceText = computed(() => {
  if (props.course.isFree) {
    return '免费'
  }

  return `￥${Number(props.course.price || 0).toFixed(2)}`
})
</script>

<style scoped>
.course-description-clamp {
  display: -webkit-box;
  overflow: hidden;
  line-clamp: 3;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.course-meta-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 2rem;
  padding: 0.4rem 0.9rem;
  font-size: 0.95rem;
  font-weight: 600;
  line-height: 1.2;
}
</style>
