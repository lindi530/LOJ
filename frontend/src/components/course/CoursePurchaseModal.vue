<template>
  <Teleport to="body">
    <div
      v-if="show"
      class="modal fade show d-block course-purchase-modal"
      tabindex="-1"
      role="dialog"
      aria-modal="true"
      aria-label="购买课程"
      aria-labelledby="coursePurchaseModalTitle"
    >
      <div class="modal-dialog modal-lg modal-dialog-centered">
        <div class="modal-content border-0 shadow">
          <div class="modal-header align-items-start gap-3">
            <div class="flex-grow-1 min-w-0">
              <div class="d-flex flex-column flex-sm-row justify-content-between align-items-start gap-2">
                <h1 id="coursePurchaseModalTitle" class="modal-title fs-5">购买课程</h1>
                <span
                  v-if="hasBalance"
                  class="badge text-bg-light border fw-semibold align-self-start align-self-sm-center"
                >
                  当前余额：{{ balanceText }}
                </span>
              </div>
              <p class="small text-secondary mb-0 text-break">{{ courseTitle }}</p>
            </div>
            <button
              class="btn-close"
              type="button"
              aria-label="关闭"
              @click="$emit('close')"
            ></button>
          </div>

          <div class="modal-body">
            <div v-if="loading" class="text-center py-5">
              <span class="spinner-border text-primary mb-3" aria-hidden="true"></span>
              <p class="mb-0 text-secondary">正在创建课程订单，请稍候。</p>
            </div>

            <div v-else-if="errorMessage" class="alert alert-danger mb-0" role="alert">
              <div class="d-flex gap-2">
                <i class="bi bi-exclamation-triangle-fill flex-shrink-0" aria-hidden="true"></i>
                <div>
                  <div class="fw-semibold">订单创建失败</div>
                  <div class="small mb-0">{{ errorMessage }}</div>
                </div>
              </div>
            </div>

            <div v-else-if="order">
              <div
                class="alert d-flex gap-2"
                :class="isActiveOrder ? 'alert-success' : 'alert-primary'"
                role="status"
              >
                <i
                  class="bi flex-shrink-0"
                  :class="isActiveOrder ? 'bi-check-circle-fill' : 'bi-wallet2'"
                  aria-hidden="true"
                ></i>
                <div>
                  <div class="fw-semibold">{{ statusTitle }}</div>
                  <div class="small mb-0">{{ statusDescription }}</div>
                </div>
              </div>

              <div v-if="payErrorMessage" class="alert alert-danger d-flex gap-2" role="alert">
                <i class="bi bi-exclamation-triangle-fill flex-shrink-0" aria-hidden="true"></i>
                <div>
                  <div class="fw-semibold">支付失败</div>
                  <div class="small mb-0">{{ payErrorMessage }}</div>
                </div>
              </div>

              <ul class="list-group list-group-flush border rounded-3">
                <li class="list-group-item d-flex justify-content-between gap-3">
                  <span class="text-secondary">支付金额</span>
                  <span class="fw-semibold text-end">{{ amountText }}</span>
                </li>
                <li
                  v-if="order.orderNo"
                  class="list-group-item d-flex flex-column flex-sm-row justify-content-between gap-2"
                >
                  <span class="text-secondary">订单号</span>
                  <span class="font-monospace text-break text-sm-end">{{ order.orderNo }}</span>
                </li>
                <li class="list-group-item d-flex justify-content-between gap-3">
                  <span class="text-secondary">订单状态</span>
                  <span class="badge rounded-pill" :class="statusBadgeClass">{{ statusText }}</span>
                </li>
              </ul>
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn btn-outline-secondary" type="button" @click="$emit('close')">
              关闭
            </button>
            <button
              v-if="errorMessage"
              class="btn btn-primary"
              type="button"
              :disabled="loading"
              @click="$emit('retry')"
            >
              重新请求
            </button>
            <button
              v-else-if="canPay"
              class="btn btn-primary"
              type="button"
              :disabled="payLoading"
              @click="$emit('pay')"
            >
              <span
                v-if="payLoading"
                class="spinner-border spinner-border-sm me-2"
                aria-hidden="true"
              ></span>
              <i v-else class="bi bi-wallet2 me-2" aria-hidden="true"></i>
              立刻支付
            </button>
            <button
              v-else-if="isActiveOrder"
              class="btn btn-primary"
              type="button"
              @click="$emit('go-learn')"
            >
              <i class="bi bi-play-circle me-2" aria-hidden="true"></i>
              进入学习
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="show" class="modal-backdrop fade show"></div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  courseTitle: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  payLoading: {
    type: Boolean,
    default: false
  },
  order: {
    type: Object,
    default: null
  },
  errorMessage: {
    type: String,
    default: ''
  },
  payErrorMessage: {
    type: String,
    default: ''
  }
})

defineEmits(['close', 'retry', 'pay', 'go-learn'])

const normalizedStatus = computed(() => Number(props.order?.status ?? -1))
const normalizedAmount = computed(() => Number(props.order?.amount ?? 0))
const normalizedBalance = computed(() => Number(props.order?.balance ?? 0))
const isActiveOrder = computed(() => normalizedStatus.value === 1)
const hasBalance = computed(() => props.order?.balance !== undefined && props.order?.balance !== null)
const canPay = computed(() => (
  normalizedStatus.value === 0 &&
  Boolean(props.order?.orderNo) &&
  !props.loading
))

const amountText = computed(() => {
  if (normalizedAmount.value <= 0) {
    return '0 虚拟币'
  }

  return `${normalizedAmount.value} 虚拟币`
})

const balanceText = computed(() => `${normalizedBalance.value} 虚拟币`)

const statusText = computed(() => {
  if (normalizedStatus.value === 1) {
    return '已开通'
  }

  if (normalizedStatus.value === 0) {
    return '待支付'
  }

  return '未知'
})

const statusBadgeClass = computed(() => (
  isActiveOrder.value ? 'text-bg-success' : 'text-bg-warning'
))

const statusTitle = computed(() => {
  if (isActiveOrder.value) {
    if (hasBalance.value) {
      return '支付成功'
    }

    return props.order?.free ? '课程已开通' : '课程已拥有'
  }

  return '确认支付'
})

const statusDescription = computed(() => {
  if (isActiveOrder.value) {
    if (hasBalance.value) {
      return '虚拟币支付已完成，可以进入课程内容学习。'
    }

    return '现在可以返回课程内容继续学习。'
  }

  return '请确认订单信息，使用虚拟币完成支付。'
})
</script>

<style scoped>
.course-purchase-modal {
  visibility: visible;
}
</style>
