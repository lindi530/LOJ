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
            <div class="d-flex flex-column flex-sm-row justify-content-between align-items-sm-start gap-3 mb-3">
              <div class="d-flex flex-wrap gap-2">
                <span class="badge course-meta-badge rounded-pill text-bg-primary">{{ course.level || '课程' }}</span>
                <span class="badge course-meta-badge rounded-pill text-bg-light border">{{ course.duration || '长期有效' }}</span>
                <span
                  class="badge course-meta-badge rounded-pill"
                  :class="course.isFree ? 'text-bg-success' : 'text-bg-warning'"
                >
                  {{ priceText }}
                </span>
              </div>

              <button
                class="btn btn-primary flex-shrink-0 align-self-start"
                type="button"
                :disabled="courseAccessLoading || (!hasCourseAccess && purchaseLoading)"
                @click="handleCourseAction"
              >
                <span
                  v-if="courseAccessLoading || (!hasCourseAccess && purchaseLoading)"
                  class="spinner-border spinner-border-sm me-2"
                  aria-hidden="true"
                ></span>
                <i
                  v-else
                  class="bi me-2"
                  :class="hasCourseAccess ? 'bi-play-circle' : 'bi-cart-check'"
                  aria-hidden="true"
                ></i>
                {{ courseAccessLoading ? '加载中' : (hasCourseAccess ? '进入课程' : '购买课程') }}
              </button>
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
                  <div class="fw-bold">{{ course.chapterCount || 0 }}</div>
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
          :can-create-chapter="canCreateChapter"
        />
      </div>
    </div>

    <CoursePurchaseModal
      :show="purchaseModalVisible"
      :course-title="course.title"
      :loading="purchaseLoading"
      :pay-loading="purchasePayLoading"
      :order="purchaseOrder"
      :error-message="purchaseError"
      :pay-error-message="purchasePayError"
      @close="closePurchaseModal"
      @retry="createCourseOrder"
      @pay="payCourseOrder"
      @go-learn="goToCourseLearningPage"
    />
  </section>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import api from '@/api'
import CourseContent from './CourseContent.vue'
import CourseIntro from './CourseIntro.vue'
import CoursePurchaseModal from './CoursePurchaseModal.vue'

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

const store = useStore()
const router = useRouter()
const purchaseModalVisible = ref(false)
const purchaseLoading = ref(false)
const purchasePayLoading = ref(false)
const purchaseError = ref('')
const purchasePayError = ref('')
const purchaseOrder = ref(null)
const courseOwned = ref(false)
const courseCreatorByAccess = ref(false)
const courseAccessLoading = ref(false)
let purchaseRequestId = 0
let purchasePayRequestId = 0
let courseAccessRequestId = 0

const hasAuthToken = computed(() => Boolean(
  store.getters['user/accessToken'] || localStorage.getItem('accessToken')
))
const currentUserId = computed(() => store.getters['user/userId'])
const courseCreatorId = computed(() => (
  props.course.createdBy ??
  props.course.created_by ??
  props.course.creatorId ??
  props.course.creator_id ??
  ''
))
const isCourseCreator = computed(() => (
  Boolean(currentUserId.value) &&
  Boolean(courseCreatorId.value) &&
  String(currentUserId.value) === String(courseCreatorId.value)
))
const canCreateChapter = computed(() => isCourseCreator.value || courseCreatorByAccess.value)
const hasCourseAccess = computed(() => (
  canCreateChapter.value ||
  courseOwned.value
))

const priceText = computed(() => {
  if (props.course.isFree) {
    return '免费'
  }

  return `￥${Number(props.course.price || 0).toFixed(2)}`
})

const resetPurchaseState = () => {
  purchaseError.value = ''
  purchasePayError.value = ''
  purchaseOrder.value = null
}

const goToCourseLearningPage = () => {
  purchaseModalVisible.value = false
  router.push({
    name: 'CourseDetail',
    params: { course_id: props.course.id },
    hash: '#content'
  })
}

const createCourseOrder = async () => {
  if (purchaseLoading.value || !props.course.id) {
    return
  }

  const requestId = ++purchaseRequestId
  purchaseLoading.value = true
  resetPurchaseState()

  try {
    const resp = await api.createCourseOrder(props.course.id)

    if (requestId !== purchaseRequestId) {
      return
    }

    if (resp?.code !== 0) {
      throw new Error(resp?.message || '订单创建失败，请稍后重试。')
    }

    purchaseOrder.value = resp?.data || {}
    if (purchaseOrder.value?.owned) {
      courseOwned.value = true
    }
  } catch (error) {
    if (requestId !== purchaseRequestId) {
      return
    }

    purchaseError.value = error?.response?.data?.message || error?.message || '订单创建失败，请稍后重试。'
  } finally {
    if (requestId === purchaseRequestId) {
      purchaseLoading.value = false
    }
  }
}

const payCourseOrder = async () => {
  const orderNo = purchaseOrder.value?.orderNo

  if (purchasePayLoading.value || !orderNo) {
    return
  }

  const requestId = ++purchasePayRequestId
  purchasePayLoading.value = true
  purchasePayError.value = ''

  try {
    const resp = await api.payCourseOrder(orderNo)

    if (requestId !== purchasePayRequestId) {
      return
    }

    if (resp?.code !== 0) {
      throw new Error(resp?.message || '支付失败，请稍后重试。')
    }

    purchaseOrder.value = {
      ...purchaseOrder.value,
      ...(resp?.data || {})
    }
    if (purchaseOrder.value?.owned) {
      courseOwned.value = true
    }
  } catch (error) {
    if (requestId !== purchasePayRequestId) {
      return
    }

    purchasePayError.value = error?.response?.data?.message || error?.message || '支付失败，请稍后重试。'
  } finally {
    if (requestId === purchasePayRequestId) {
      purchasePayLoading.value = false
    }
  }
}

const openPurchaseModal = () => {
  purchaseModalVisible.value = true
  createCourseOrder()
}

const handleCourseAction = () => {
  if (courseAccessLoading.value) {
    return
  }

  if (hasCourseAccess.value) {
    goToCourseLearningPage()
    return
  }

  openPurchaseModal()
}

const closePurchaseModal = () => {
  purchaseModalVisible.value = false
}

const loadCourseAccess = async () => {
  const requestId = ++courseAccessRequestId
  courseOwned.value = false
  courseCreatorByAccess.value = false

  if (!props.course.id || !hasAuthToken.value) {
    courseAccessLoading.value = false
    return
  }

  courseAccessLoading.value = true

  try {
    const resp = await api.getCourseAccess(props.course.id)

    if (requestId !== courseAccessRequestId) {
      return
    }

    if (resp?.code === 0) {
      const access = resp?.data || {}
      courseOwned.value = Boolean(access.owned)
      courseCreatorByAccess.value = Boolean(access.creator)
    }
  } catch (error) {
    if (requestId !== courseAccessRequestId) {
      return
    }

    courseOwned.value = false
    courseCreatorByAccess.value = false
  } finally {
    if (requestId === courseAccessRequestId) {
      courseAccessLoading.value = false
    }
  }
}

watch(() => props.course.id, () => {
  purchaseRequestId += 1
  purchasePayRequestId += 1
  purchaseLoading.value = false
  purchasePayLoading.value = false
  purchaseModalVisible.value = false
  resetPurchaseState()
})

watch([() => props.course.id, hasAuthToken], loadCourseAccess, { immediate: true })
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
