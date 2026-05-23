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
        height="100%"
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
import { ref, watch, onMounted, onUnmounted, defineAsyncComponent } from 'vue'
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
const theme = ref('vs-light')
const userEdited = ref(false)
const message = useMessage()
const isEditorVisible = ref(true)
const testSample = ref(true)
const isLoading = ref(false)
const activeStatus = ref('')
const outputValue = ref('')
const { registerSubmitCodeCallback } = useWebSocketContext()

function defaultCode(lang) {
  return {
    cpp: `#include <bits/stdc++.h>\nusing namespace std;\nint main() {\n  return 0;\n}`,
    python: `class Solution:\n    def twoSum(self, nums, target):\n        pass`,
    go: `package main\n\nimport(\n  \"fmt\"\n)\n\nfunc main() {\n}`
  }[lang] || ''
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

onMounted(() => {
  internalCode.value = defaultCode(internalLang.value)
  
  const { stop } = useIntersectionObserver(
    editorContainer,
    ([{ isIntersecting }]) => {
      isEditorVisible.value = isIntersecting
      
      if (isIntersecting && !editorReady.value) {
        setTimeout(() => {
          editorReady.value = true
        }, 100)
      }
    }
  )
  
  onUnmounted(() => {
    stop()
    unregister()
  })
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
  height: clamp(300px, 50vh, 520px);
}

:deep(.n-card__content) {
  display: flex;
  flex-direction: column;
}
</style>
