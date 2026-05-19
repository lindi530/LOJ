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
    <n-input
      v-model:value="input"
      type="textarea"
      rows="3"
      placeholder=""
      class="mb-2"
    />
    <div class="d-flex justify-content-between mb-2">
      <n-button type="success" @click="runFunction">运行</n-button>
    </div>

    <n-card
      size="small"
      :segmented="{ content: true }"
      style="max-height: 120px; overflow: auto"
    >
      <pre class="result-text">{{ outputValue }}</pre>
    </n-card>
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
.d-flex {
  display: flex;
}
.justify-content-between {
  justify-content: space-between;
}
.mb-2 {
  margin-bottom: 8px;
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
</style>
