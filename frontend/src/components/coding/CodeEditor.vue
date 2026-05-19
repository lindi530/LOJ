<template>
  <n-card title="在线提交">
    <template #header-extra>
      <n-space align="center" style="gap: 12px">
        <n-select
          v-model:value="internalLang"
          :options="langOptions"
          style="width: 120px"
          size="small"
        />
        <n-button size="small" @click="toggleTheme" tertiary>
          {{ theme === 'vs-light' ? '切换深色' : '切换浅色' }}
        </n-button>
      </n-space>
    </template>

    <!-- 编辑器容器：使用flex布局确保占满可用空间 -->
    <div ref="editorContainer" class="editor-container mb-3">
      <MonacoEditor
        v-if="editorReady"
        v-model:value="internalCode"
        :language="internalLang"
        :theme="theme"
        :height="editorHeight"
        :options="editorOptions"
      />
    </div>

    <div class="d-flex justify-content-between">
      <n-button @click="resetCode">重置代码</n-button>
      <n-button 
        type="primary" 
        @click="submitCode"
        :disabled="isLoading"
        :class="{ 'opacity-50 cursor-not-allowed': isLoading }"
      >
        {{ isLoading ? '提交中...' : '提交' }}
      </n-button>
    </div>

    <SampleTest
      :activeStatus="activeStatus"
      :testSample="testSample"
      :outputValue="outputValue"
      @handleActiveStatus="handleActiveStatus"
      @handleTestSample="handleTestSample"
      @handleRunExample="handleRunExample"
    />
  </n-card>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick, defineAsyncComponent } from 'vue'
import { NCard, NSelect, NButton, NSpace, useMessage } from 'naive-ui'
import { useIntersectionObserver } from '@vueuse/core'
import api from '@/api'
import SampleTest from './SampleTest.vue'
import { useWebSocketContext } from '@/composables/useWebSocket.js'

const MonacoEditor = defineAsyncComponent(() => import('monaco-editor-vue3'))

// props + emit
const props = defineProps({
  problemId: Number,
  roomId: String,
})

console.log("codeEditor roomId:", props.problemId, "  ", props.roomId );

const langOptions = [
  { label: 'C++', value: 'cpp' },
  { label: 'Python3', value: 'python' },
  { label: 'Go', value: 'go' }
]

const internalLang = ref('cpp')
const internalCode = ref('')
const editorReady = ref(false)
const editorContainer = ref(null)
const editorHeight = ref(0) // 初始化为0，由父容器计算决定
const theme = ref('vs-light')
const userEdited = ref(false)
const message = useMessage()
const isEditorVisible = ref(true)
const testSample = ref(true)
const isLoading = ref(false)
const activeStatus = ref('')
const outputValue = ref('')
const { registerSubmitCodeCallback } = useWebSocketContext()

// 存储 ResizeObserver 实例
let resizeObserver = null

function defaultCode(lang) {
  return {
    cpp: `#include <bits/stdc++.h>\nusing namespace std;\nint main() {\n  return 0;\n}`,
    python: `class Solution:\n    def twoSum(self, nums, target):\n        pass`,
    go: `package main\n\nimport(\n  \"fmt\"\n)\n\nfunc main() {\n}`
  }[lang] || ''
}

// 计算编辑器可用高度（减去其他元素占用的空间）
const calculateEditorHeight = () => {
  if (!editorContainer.value) return 0
  
  // 获取父容器高度
  const containerRect = editorContainer.value.parentElement.getBoundingClientRect()
  
  // 计算其他元素占用的高度（卡片头部、按钮区域、测试区域）
  const headerHeight = 60 // 卡片头部大致高度
  const buttonAreaHeight = 50 // 按钮区域高度
  const testAreaHeight = 200 // 测试区域大致高度
  const padding = 30 // 内边距总和
  
  // 计算可用高度
  const availableHeight = containerRect.height - headerHeight - buttonAreaHeight - testAreaHeight - padding
  return Math.max(300, availableHeight) // 确保最小高度为300px
}

// 监听容器尺寸变化并更新编辑器高度
const handleContainerResize = () => {
  if (editorContainer.value) {
    const newHeight = calculateEditorHeight()
    editorHeight.value = newHeight
  }
}

const handleTestSample = (value) => {
  testSample.value = value
}

const handleActiveStatus = (value) => { 
  activeStatus.value = value
}

const handleRunExample = async (value) => { 
  const resp = await api.submitExample(
    props.problemId, 
    {
      "input": value,
      "language": internalLang.value,
      "code": internalCode.value,
    }
  )
  if (resp.code === 0) {
    outputValue.value = resp.data.output;
  }
}

const unregister = registerSubmitCodeCallback((msg) => {
  if (msg.content) {
    handleActiveStatus(msg.content)
  }
})

function resetCode() {
  internalCode.value = defaultCode(internalLang.value)
  userEdited.value = false
  message.info('代码已重置')
}

function toggleTheme() {
  theme.value = theme.value === 'vs-light' ? 'vs-dark' : 'vs-light'
}

async function submitCode() {
  if (isLoading.value) return

  isLoading.value = true
  handleTestSample(false)
  handleActiveStatus("")

  try {
    console.log("props.roomId: ", props.roomId)
    if (typeof props.roomId !== 'undefined') {
      console.log("props.roomId")
      const payload = {
          room_id: props.roomId,
          language: internalLang.value,
          code: internalCode.value,
      }
      console.log("payload: ", payload)
      const resp = await api.saberSubmit(
        payload
       )
    } else {
      console.log("2")
      const resp = await api.submitCode(
      props.problemId,
      {
        "language": internalLang.value,
        "code": internalCode.value,
      })
    }
    
  } catch (error) {
    message.error('提交失败：' + error.message)
  } finally {
    isLoading.value = false
  }
}

watch(internalLang, (newLang) => {
  internalCode.value = defaultCode(newLang)
})

watch(internalCode, () => {
  if (!userEdited.value) {
    userEdited.value = true
  }
})

// 编辑器配置：启用自动布局
const editorOptions = {
  minimap: { enabled: false },
  automaticLayout: true, // 关键：启用自动布局适应容器变化
  scrollBeyondLastLine: false,
  fontSize: 14
}

onMounted(async () => {
  internalCode.value = defaultCode(internalLang.value)
  await nextTick()
  
  // 初始化时计算一次高度
  handleContainerResize()
  
  // 监听容器尺寸变化（包括父组件大小变化）
  if ('ResizeObserver' in window) {
    resizeObserver = new ResizeObserver(entries => {
      requestAnimationFrame(() => {
        handleContainerResize()
      })
    })
    
    if (editorContainer.value && editorContainer.value.parentElement) {
      // 监听父元素的尺寸变化（关键）
      resizeObserver.observe(editorContainer.value.parentElement)
    }
  }
  
  const { stop } = useIntersectionObserver(
    editorContainer,
    ([{ isIntersecting }]) => {
      isEditorVisible.value = isIntersecting
      
      if (isIntersecting && !editorReady.value) {
        setTimeout(() => {
          editorReady.value = true
          // 可见时再计算一次高度
          handleContainerResize()
        }, 100)
      }
    }
  )
  
  onUnmounted(() => {
    stop()
    unregister()
  })
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
})
</script>

<style scoped>
.d-flex {
  display: flex;
}
.justify-content-between {
  justify-content: space-between;
}
.mb-3 {
  margin-bottom: 12px;
}
.editor-container {
  width: 100%;
  /* 关键：使用100%高度而不是固定值 */
  height: 100%;
  min-height: 300px; /* 确保最小高度 */
}

/* 确保父容器可以被正确计算高度 */
:deep(.n-card__content) {
  display: flex;
  flex-direction: column;
  height: 100%;
}
</style>