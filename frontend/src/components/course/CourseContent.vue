<template>
  <div>
    <div class="d-flex flex-column flex-md-row gap-2 justify-content-between align-items-md-center mb-3">
      <div>
        <h2 class="h4 fw-bold mb-1">课程章节</h2>
        <p class="text-secondary mb-0">默认只展示章节标题和简介，展开章节后加载视频与练习题。</p>
      </div>
      <div class="d-flex flex-wrap align-items-stretch gap-2 align-self-start align-self-md-center">
        <span class="badge text-bg-light border d-inline-flex align-items-center justify-content-center py-2 px-3 course-content-toolbar-item">
          {{ chapters.length }} 章
        </span>
        <RouterLink
          v-if="canCreateChapter"
          class="btn btn-primary btn-sm d-inline-flex align-items-center justify-content-center py-2 px-3 course-content-toolbar-item"
          :to="{ name: 'CourseChapterCreate', params: { course_id: courseId } }"
        >
          <i class="bi bi-plus-circle me-1" aria-hidden="true"></i>
          添加章节
        </RouterLink>
      </div>
    </div>

    <div v-if="chaptersLoading" class="d-flex align-items-center gap-2 text-secondary py-3">
      <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
      <span>正在加载章节...</span>
    </div>

    <div v-else-if="chaptersError" class="alert alert-danger d-flex flex-column flex-sm-row gap-3 justify-content-between align-items-sm-center" role="alert">
      <span>{{ chaptersError }}</span>
      <button class="btn btn-outline-danger btn-sm align-self-start align-self-sm-center" type="button" @click="loadChapters">
        重试
      </button>
    </div>

    <div v-else-if="chapters.length > 0" :id="accordionId" class="accordion">
      <div
        v-for="(chapter, index) in sortedChapters"
        :key="chapter.id"
        class="accordion-item"
      >
        <h3 class="accordion-header">
          <button
            class="accordion-button collapsed"
            type="button"
            data-bs-toggle="collapse"
            :data-bs-target="`#${chapterCollapseId(chapter)}`"
            aria-expanded="false"
            :aria-controls="chapterCollapseId(chapter)"
            @click="loadChapterDetail(chapter.id)"
          >
            <span class="d-flex flex-column flex-md-row gap-1 gap-md-3 align-items-md-center flex-grow-1 pe-3">
              <span class="d-inline-flex align-items-center gap-2">
                <span class="badge text-bg-light border">第 {{ chapter.sort || index + 1 }} 章</span>
                <span class="fw-bold">{{ chapter.title }}</span>
              </span>
              <span v-if="chapter.description" class="small text-secondary">{{ chapter.description }}</span>
            </span>
          </button>
        </h3>
        <div
          :id="chapterCollapseId(chapter)"
          class="accordion-collapse collapse"
          :data-bs-parent="`#${accordionId}`"
        >
          <div class="accordion-body bg-body-tertiary">
            <div v-if="chapterDetail(chapter.id).loading" class="d-flex align-items-center gap-2 text-secondary">
              <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
              <span>正在加载章节内容...</span>
            </div>

            <div v-else-if="chapterDetail(chapter.id).error" class="alert alert-warning d-flex flex-column flex-sm-row gap-3 justify-content-between align-items-sm-center mb-0" role="alert">
              <span>{{ chapterDetail(chapter.id).error }}</span>
              <button class="btn btn-outline-warning btn-sm align-self-start align-self-sm-center" type="button" @click="loadChapterDetail(chapter.id, true)">
                重试
              </button>
            </div>

            <div v-else-if="chapterDetail(chapter.id).loaded">
              <div v-if="chapterDetail(chapter.id).video?.id" class="card border border-secondary-subtle bg-secondary-subtle mb-3">
                <div class="card-body d-flex flex-column flex-sm-row justify-content-between align-items-sm-center gap-2 py-2 py-sm-3">
                  <span class="fw-semibold">该章节讲解视频</span>
                  <RouterLink
                    class="btn btn-outline-primary btn-sm d-inline-flex align-items-center justify-content-center text-center align-self-start align-self-sm-center"
                    :to="{ name: 'CourseChapterVideo', params: { course_id: courseId, chapter_id: chapter.id, video_id: chapterDetail(chapter.id).video.id } }"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <i class="bi bi-play-circle me-1" aria-hidden="true"></i>
                    观看视频
                  </RouterLink>
                </div>
              </div>

              <div v-if="chapterDetail(chapter.id).problems.length > 0" class="list-group">
                <RouterLink
                  v-for="(problem, problemIndex) in chapterDetail(chapter.id).problems"
                  :id="problemAnchorId(chapter, problem, problemIndex)"
                  :key="problem.id || problemIndex"
                  class="list-group-item list-group-item-action d-flex justify-content-between align-items-center gap-3"
                  :to="{ name: 'ProblemDetail', params: { problem_id: problem.id } }"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <span>
                    <span class="badge text-bg-secondary me-2">{{ problemIndex + 1 }}</span>
                    {{ problem.title }}
                  </span>
                  <span class="small text-secondary text-nowrap">
                    #{{ problem.id }}
                    <i class="bi bi-chevron-right ms-1" aria-hidden="true"></i>
                  </span>
                </RouterLink>
              </div>

              <div v-else class="alert alert-light border text-secondary mb-0">
                本章暂无练习题。
              </div>
            </div>

            <div v-else class="text-secondary">
              展开章节后加载内容。
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="alert alert-secondary mt-3 mb-0">
      暂无章节内容。
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import api from '@/api'

const props = defineProps({
  courseId: {
    type: [Number, String],
    required: true
  },
  canCreateChapter: {
    type: Boolean,
    default: false
  }
})

const chapters = ref([])
const chaptersLoading = ref(false)
const chaptersError = ref('')
const chapterDetails = reactive({})
const emptyChapterDetail = {
  loading: false,
  loaded: false,
  error: '',
  problems: [],
  video: null
}

const accordionId = computed(() => `course-content-${props.courseId}`)
const sortedChapters = computed(() => (
  [...chapters.value].sort((left, right) => (
    Number(left.sort || 0) - Number(right.sort || 0) ||
    Number(left.id || 0) - Number(right.id || 0)
  ))
))

function normalizeChapter(chapter) {
  return {
    id: chapter.id,
    title: chapter.title || '未命名章节',
    description: chapter.description || '',
    sort: Number(chapter.sort || 0)
  }
}

function normalizeProblem(problem) {
  return {
    id: problem.id,
    title: problem.title || '未命名题目'
  }
}

function normalizeChapterDetail(data) {
  const problems = Array.isArray(data?.problems) ? data.problems : []

  return {
    problems: problems.map(normalizeProblem),
    video: data?.video && typeof data.video === 'object' ? { id: data.video.id } : null
  }
}

function resetChapterDetails() {
  Object.keys(chapterDetails).forEach((chapterId) => {
    delete chapterDetails[chapterId]
  })
}

function ensureChapterDetail(chapterId) {
  const key = String(chapterId)

  if (!chapterDetails[key]) {
    chapterDetails[key] = { ...emptyChapterDetail }
  }

  return chapterDetails[key]
}

function chapterDetail(chapterId) {
  return chapterDetails[String(chapterId)] || emptyChapterDetail
}

async function loadChapters() {
  if (!props.courseId) {
    chapters.value = []
    return
  }

  chaptersLoading.value = true
  chaptersError.value = ''
  resetChapterDetails()

  try {
    const resp = await api.getCourseChapters(props.courseId)
    const data = resp?.data ?? resp
    const list = Array.isArray(data) ? data : Array.isArray(data?.list) ? data.list : []

    chapters.value = list.map(normalizeChapter)
  } catch (error) {
    chapters.value = []
    chaptersError.value = error?.response?.data?.message || error?.message || '章节列表加载失败，请稍后重试。'
  } finally {
    chaptersLoading.value = false
  }
}

async function loadChapterDetail(chapterId, force = false) {
  const detail = ensureChapterDetail(chapterId)

  if (!force && (detail.loading || detail.loaded)) {
    return
  }

  detail.loading = true
  detail.error = ''

  try {
    const resp = await api.getCourseChapterDetail(props.courseId, chapterId)
    const data = resp?.data ?? resp
    const normalized = normalizeChapterDetail(data)

    detail.problems = normalized.problems
    detail.video = normalized.video
    detail.loaded = true
  } catch (error) {
    detail.loaded = false
    detail.problems = []
    detail.video = null
    detail.error = error?.response?.data?.message || error?.message || '章节内容加载失败，请稍后重试。'
  } finally {
    detail.loading = false
  }
}

function chapterCollapseId(chapter) {
  return `course-${props.courseId}-chapter-${chapter.id}`
}

function problemAnchorId(chapter, problem, index) {
  return `chapter-${chapter.id}-problem-${problem.id || index + 1}`
}

watch(() => props.courseId, loadChapters, { immediate: true })
</script>

<style scoped>
.course-content-toolbar-item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 2.5rem;
  font-size: 0.875rem;
  line-height: 1;
  text-align: center;
  white-space: nowrap;
}
</style>
