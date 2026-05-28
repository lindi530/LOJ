<template>
  <n-card size="small" :segmented="{ content: true }">
    <template #header>
      <div class="header-content">
        <span class="title-text">
          {{ testSample === true ? '样例测试状态' : '代码提交状态' }}：
        </span>

        <span 
          :class="statusClass"
          class="status-value"
        >
          {{ activeStatus }}
        </span>
      </div>
    </template>
    <div class="sample-test-grid">
      <label class="sample-test-field">
        <span class="sample-test-label">输入</span>
        <n-input
          v-model:value="input"
          type="textarea"
          rows="3"
          placeholder=""
          class="sample-test-control"
        />
      </label>

      <div class="sample-test-field">
        <span class="sample-test-label">输出</span>
        <n-card
          size="small"
          :segmented="{ content: true }"
          class="sample-test-output"
        >
          <pre class="result-text">{{ outputValue || '运行后将在这里显示输出。' }}</pre>
        </n-card>
      </div>
    </div>

    <div class="sample-test-actions mt-2">
      <n-button class="sample-run-button" type="success" @click="runFunction">
        <i class="bi bi-play-fill me-1" aria-hidden="true"></i>运行样例
      </n-button>
      <slot name="actions"></slot>
    </div>
  </n-card>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { NCard, NInput, NButton } from 'naive-ui'

const props = defineProps({
  testSample: Boolean,
  activeStatus: String,
  outputValue: String,
})

const emit = defineEmits(['handleTestSample', "handleActiveStatus", 'handleRunExample'])

const input = ref("")

// 样例提交
const runFunction = async () => {
  console.log("提交中...")
  emit('handleTestSample', true)
  emit('handleActiveStatus', "")
  emit('handleRunExample', input.value)
}


const statusClass = computed(() => {
  switch (props.activeStatus) {
    case 'Finished':
      return 'status-success';
    case 'Accepted':
      return 'status-success';
    case 'Running':
      return 'status-run';
    case 'Pending':
      return 'status-pending'
    case 'Time Limit Exceeded':
      return 'status-timeout'
    default:
      return 'status-error';
  }
});
</script>

<style scoped>
.result-text {
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  font-size: 13px;
}
.sample-test-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.75rem;
}

.sample-test-field {
  display: block;
  min-width: 0;
}

.sample-test-label {
  display: block;
  margin-bottom: 0.38rem;
  color: var(--bs-secondary-color, #6c757d);
  font-size: 0.78rem;
  font-weight: 600;
}

:deep(.sample-test-control textarea),
:deep(.sample-test-output .n-card__content) {
  height: 88px;
}

:deep(.sample-test-output .n-card__content) {
  overflow: auto;
}

.sample-run-button {
  min-width: 8.4rem;
}

.sample-test-actions {
  display: flex;
  gap: 0.62rem;
  justify-content: flex-end;
}

.n-card--success :deep(.n-card__header) {
  color: #2e7d32;
}
.n-card--fail :deep(.n-card__header) {
  color: #c62828;
}

.status-value {
  font-size: 15px;
}

/* 不同状态的颜色 */

.status-success {
  color: #18a058; /* 成功绿色 */
}
.status-error, .status-timeout {
  color: #d03050; /* 错误红色 */
}
.status-pending {
  color: #343c38;
}
.status-run {
  color: #337AB7; /* 运行蓝色 */
}

@media (max-width: 767.98px) {
  .sample-test-grid {
    grid-template-columns: 1fr;
  }

  .sample-test-actions {
    flex-direction: column;
  }

  .sample-test-actions :deep(.n-button) {
    width: 100%;
  }
}
</style>
