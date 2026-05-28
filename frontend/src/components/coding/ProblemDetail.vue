<template>
    <n-card :bordered="false" size="small" style="margin-bottom: 16px;">
      <template #header>
        <div class="d-flex justify-content-between align-items-center">
          <h2 class="mb-0">{{ problem.title }}</h2>
        </div>
        <!-- <button 
          v-if="roomId === undefined"
          class="mt-4 px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors"
          @click="switchToSubmissions"
        >
          提交记录
        </button> -->
      </template>
      <template #header-extra>
        <n-space>
          <n-tag :type="levelTagType(problem.level)" size="small">
            {{ problem.level }}
          </n-tag>
          <n-tag
            v-for="tag in problem.tags"
            :key="tag"
            type="info"
            size="small"
          >
            {{ tag }}
          </n-tag>
        </n-space>
      </template>
    </n-card>

    <!-- 描述模块 -->
    <n-card title="题目描述" size="small" class="mb-3">
      <div v-html="problem.description" />
    </n-card>

    <!-- 输入描述 -->
    <n-card title="输入描述" size="small" class="mb-3">
      <p>{{ problem.input_desc }}</p>
    </n-card>

    <!-- 输出描述 -->
    <n-card title="输出描述" size="small" class="mb-3">
      <p>{{ problem.output_desc }}</p>
    </n-card>

    <!-- 示例部分 -->
    <n-card title="示例" size="small">
      <div
        v-for="(ex, idx) in problem.examples"
        :key="idx"
        class="vstack gap-2 mb-3"
      >
        <div class="sample-block bg-body-tertiary border rounded-3">
          <div class="sample-header d-flex align-items-center justify-content-between">
            <span class="small fw-semibold text-body-secondary">输入 {{ idx + 1 }}</span>
            <button
              class="btn btn-sm btn-link sample-copy-btn"
              type="button"
              :aria-label="`复制示例 ${idx + 1} 输入`"
              title="复制"
              @click="copySample(ex.input, `示例 ${idx + 1} 输入`)"
            >
              <i class="bi bi-copy" aria-hidden="true"></i>
            </button>
          </div>
          <pre class="m-0 px-3 pb-3 pt-2"><code>{{ ex.input }}</code></pre>
        </div>
        <div class="sample-block bg-body-tertiary border rounded-3">
          <div class="sample-header d-flex align-items-center justify-content-between">
            <span class="small fw-semibold text-body-secondary">输出 {{ idx + 1 }}</span>
            <button
              class="btn btn-sm btn-link sample-copy-btn"
              type="button"
              :aria-label="`复制示例 ${idx + 1} 输出`"
              title="复制"
              @click="copySample(ex.output, `示例 ${idx + 1} 输出`)"
            >
              <i class="bi bi-copy" aria-hidden="true"></i>
            </button>
          </div>
          <pre class="m-0 px-3 pb-3 pt-2"><code>{{ ex.output }}</code></pre>
        </div>
      </div>
    </n-card>
</template>

<script setup>
import { useMessage } from 'naive-ui'

const message = useMessage()

const props = defineProps({
  roomId: {
    type: Number,
    default: undefined
  },
  problem: {
    type: Object,
    required: true
  },
  tab: {
    type: String,
    default: undefined
  }
})

console.log("roomId: ", props.roomId)

const emit = defineEmits(['update:tab'])

const switchToSubmissions = () => {
  // 触发事件，通知父组件更新activeTab为'submissions'
  emit('update:tab', 'submissions')
}

async function copySample(text, label) {
  try {
    await navigator.clipboard.writeText(text || '')
    message.success(`${label}已复制`)
  } catch (error) {
    message.error('复制失败，请手动选择内容复制')
  }
}

function levelTagType(level) {
  switch (level) {
    case '简单':
      return 'success'
    case '中等':
      return 'warning'
    case '困难':
      return 'error'
    default:
      return 'default'
  }
}
</script>

<style scoped>
.sample-header {
  min-height: 2.25rem;
  padding: 0.4rem 0.75rem 0;
}

.sample-copy-btn {
  display: inline-flex;
  width: 1.9rem;
  height: 1.9rem;
  align-items: center;
  justify-content: center;
  padding: 0;
  color: var(--bs-secondary-color);
  border-radius: 999px;
  text-decoration: none;
}

.sample-copy-btn:hover {
  color: var(--bs-primary);
  background: var(--bs-tertiary-bg);
}

.sample-block pre {
  font: 14px/1.6 "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  white-space: pre-wrap;
}
</style>
