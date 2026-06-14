<template>
  <main class="py-4 py-lg-5">
    <div class="container-xl">
      <CourseDetailPanel
        v-if="isDetailPage && selectedCourse"
        :course="selectedCourse"
        :active-tab="activeTab"
      />

      <section v-else-if="coursesLoading" class="row justify-content-center">
        <div class="col-lg-7">
          <div class="card border-0 shadow-sm text-center">
            <div class="card-body p-4 p-lg-5">
              <span class="spinner-border text-primary mb-3" aria-hidden="true"></span>
              <h1 class="h4 mb-2">课程加载中</h1>
              <p class="text-secondary mb-0">正在获取课程列表，请稍候。</p>
            </div>
          </div>
        </div>
      </section>

      <section v-else-if="coursesError" class="row justify-content-center">
        <div class="col-lg-7">
          <div class="card border-0 shadow-sm text-center">
            <div class="card-body p-4 p-lg-5">
              <i class="bi bi-wifi-off display-5 text-secondary" aria-hidden="true"></i>
              <h1 class="h4 mt-3 mb-2">课程加载失败</h1>
              <p class="text-secondary mb-4">{{ coursesError }}</p>
              <button class="btn btn-primary" type="button" @click="loadCourses">
                重新加载
              </button>
            </div>
          </div>
        </div>
      </section>

      <section v-else-if="isDetailPage" class="row justify-content-center">
        <div class="col-lg-7">
          <div class="card border-0 shadow-sm text-center">
            <div class="card-body p-4 p-lg-5">
              <i class="bi bi-journal-x display-5 text-secondary" aria-hidden="true"></i>
              <h1 class="h4 mt-3 mb-2">课程不存在</h1>
              <p class="text-secondary mb-4">当前课程列表中没有找到该课程，请返回课程列表重新选择。</p>
              <RouterLink class="btn btn-primary" :to="{ name: 'Course' }">
                返回课程列表
              </RouterLink>
            </div>
          </div>
        </div>
      </section>

      <section v-else>
        <div class="card border-0 shadow-sm mb-4">
          <div class="card-body p-4 p-lg-5">
            <div class="row g-4 align-items-center">
              <div class="col-lg-8">
                <p class="text-primary fw-semibold mb-2">LOJ 课程</p>
                <h1 class="display-6 fw-bold mb-3">系统学习编程与算法</h1>
                <p class="lead text-secondary mb-0">
                  从语言基础到算法专题，按模块选择课程，循序渐进完成训练。
                </p>
              </div>
              <div class="col-lg-4">
                <div class="row g-3 text-center">
                  <div class="col-6">
                    <div class="border rounded-3 p-3 bg-body-tertiary">
                      <div class="fs-4 fw-bold">{{ total }}</div>
                      <div class="small text-secondary">课程数</div>
                    </div>
                  </div>
                  <div class="col-6">
                    <div class="border rounded-3 p-3 bg-body-tertiary">
                      <div class="fs-4 fw-bold">{{ totalStudentCount }}</div>
                      <div class="small text-secondary">学习人次</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="card border-0 shadow-sm">
          <div class="card-body p-4 p-lg-5">
            <div class="d-flex flex-column flex-md-row gap-2 justify-content-between align-items-md-center mb-4">
              <div>
                <h2 class="h4 fw-bold mb-1">全部课程</h2>
              </div>
              <RouterLink class="btn btn-primary align-self-start align-self-md-center" :to="{ name: 'CoursePublish' }">
                <i class="bi bi-plus-circle me-2" aria-hidden="true"></i>
                发布课程
              </RouterLink>
            </div>

            <div v-if="courses.length === 0" class="alert alert-light border text-secondary mb-0">
              暂无课程，请稍后再来查看。
            </div>

            <div v-else class="row g-4">
              <div
                v-for="course in courses"
                :key="course.id"
                class="col-12 col-md-6 col-xl-4 d-flex"
              >
                <CourseCard :course="course" />
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import CourseCard from '@/components/course/CourseCard.vue'
import CourseDetailPanel from '@/components/course/CourseDetailPanel.vue'
import api from '@/api'

const route = useRoute()

const courses = ref([])
const total = ref(0)
const coursesLoading = ref(false)
const coursesError = ref('')
const apiBaseUrl = process.env.VUE_APP_API_BASE || 'http://localhost:8080'

const getAppOrigin = () => (
  window.location?.origin || 'http://localhost:8080'
)

const getApiOrigin = () => {
  if (/^(https?:)?\/\//.test(apiBaseUrl)) {
    return new URL(apiBaseUrl).origin
  }

  return getAppOrigin()
}

const getAbsoluteBaseUrl = (baseUrl) => {
  const origin = getAppOrigin()
  const rawBaseUrl = baseUrl || origin
  const absoluteBaseUrl = /^(https?:)?\/\//.test(rawBaseUrl)
    ? rawBaseUrl
    : new URL(rawBaseUrl, origin).href

  return absoluteBaseUrl.endsWith('/') ? absoluteBaseUrl : `${absoluteBaseUrl}/`
}

const normalizeStoredImagePath = (coverUrl) => (
  String(coverUrl)
    .trim()
    .replace(/\\/g, '/')
    .replace(/^\/?upload\/images\//, '/images/')
    .replace(/^\/?images\//, '/images/')
)

const normalizeCoverUrl = (coverUrl) => {
  if (!coverUrl) {
    return ''
  }

  const normalizedCoverUrl = normalizeStoredImagePath(coverUrl)

  if (
    /^(https?:)?\/\//.test(normalizedCoverUrl) ||
    normalizedCoverUrl.startsWith('data:') ||
    normalizedCoverUrl.startsWith('blob:')
  ) {
    return normalizedCoverUrl
  }

  if (normalizedCoverUrl.startsWith('/')) {
    return new URL(normalizedCoverUrl, getApiOrigin()).href
  }

  return new URL(normalizedCoverUrl, getAbsoluteBaseUrl(apiBaseUrl)).href
}

const normalizeCourse = (course) => ({
  id: course.id,
  title: course.title || '未命名课程',
  description: course.description || '暂无课程描述',
  coverUrl: normalizeCoverUrl(course.coverUrl),
  price: Number(course.price || 0),
  isFree: Boolean(course.isFree),
  lessonCount: Number(course.lessonCount || 0),
  studentCount: Number(course.studentCount || 0),
  chapterCount: Number(course.chapterCount ?? course.chapter_count ?? 0),
  createdBy: course.createdBy ?? course.created_by ?? course.creatorId ?? course.creator_id ?? ''
})

const loadCourses = async () => {
  coursesLoading.value = true
  coursesError.value = ''

  try {
    const resp = await api.getCourseList()
    const data = resp?.data || {}
    const list = Array.isArray(data.list) ? data.list : []

    courses.value = list.map(normalizeCourse)
    total.value = Number(data.total ?? courses.value.length)
  } catch (error) {
    courses.value = []
    total.value = 0
    coursesError.value = error?.response?.data?.message || error?.message || '课程列表加载失败，请稍后重试。'
  } finally {
    coursesLoading.value = false
  }
}

const courseId = computed(() => route.params.course_id)
const isDetailPage = computed(() => Boolean(courseId.value))
const selectedCourse = computed(() => (
  courses.value.find((course) => String(course.id) === String(courseId.value))
))
const totalStudentCount = computed(() => (
  courses.value.reduce((sum, course) => sum + Number(course.studentCount || 0), 0)
))
const activeTab = computed(() => {
  if (route.hash === '#content' || route.hash.startsWith('#chapter-')) {
    return 'content'
  }

  return 'intro'
})

onMounted(loadCourses)
</script>
