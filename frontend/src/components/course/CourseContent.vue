<template>
  <div>
    <div class="d-flex flex-column flex-md-row gap-2 justify-content-between align-items-md-center mb-3">
      <div>
        <h2 class="h4 fw-bold mb-1">课程章节</h2>
        <p class="text-secondary mb-0">默认只展示章节标题和简介，点击章节后展开课时占位内容。</p>
      </div>
      <span class="badge text-bg-light border align-self-start align-self-md-center">
        {{ chapters.length }} 章
      </span>
    </div>

    <div :id="accordionId" class="accordion">
      <div
        v-for="chapter in chapters"
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
          >
            <span class="d-flex flex-column flex-md-row gap-1 gap-md-3 align-items-md-center flex-grow-1 pe-3">
              <span class="fw-bold">{{ chapter.title }}</span>
              <span class="small text-secondary">{{ chapter.summary }}</span>
            </span>
          </button>
        </h3>
        <div
          :id="chapterCollapseId(chapter)"
          class="accordion-collapse collapse"
          :data-bs-parent="`#${accordionId}`"
        >
          <div class="accordion-body bg-body-tertiary">
            <div class="list-group">
              <RouterLink
                v-for="(lesson, index) in chapter.lessons"
                :id="lessonAnchorId(chapter, index)"
                :key="lesson"
                class="list-group-item list-group-item-action d-flex justify-content-between align-items-center gap-3"
                :to="{
                  name: 'CourseDetail',
                  params: { course_id: courseId },
                  hash: `#${lessonAnchorId(chapter, index)}`
                }"
              >
                <span>
                  <span class="badge text-bg-secondary me-2">{{ index + 1 }}</span>
                  {{ lesson }}
                </span>
                <span class="small text-secondary text-nowrap">
                  占位
                  <i class="bi bi-chevron-right ms-1" aria-hidden="true"></i>
                </span>
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="chapters.length === 0" class="alert alert-secondary mt-3 mb-0">
      暂无章节内容。
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  courseId: {
    type: [Number, String],
    required: true
  },
  chapters: {
    type: Array,
    default: () => []
  }
})

const accordionId = `course-content-${props.courseId}`

function chapterCollapseId(chapter) {
  return `course-${props.courseId}-chapter-${chapter.id}`
}

function lessonAnchorId(chapter, index) {
  return `chapter-${chapter.id}-lesson-${index + 1}`
}
</script>
