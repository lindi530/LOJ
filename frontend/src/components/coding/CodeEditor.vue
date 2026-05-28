<template>
  <n-card title="在线提交">
    <template #header-extra>
      <n-space align="center" style="gap: 12px">
        <n-select
          v-model:value="internalLang"
          :options="langOptions"
          style="width: 120px"
          size="small"
          :disabled="!langOptions.length"
        />
        <n-button
          size="small"
          tertiary
          circle
          :aria-label="editorTheme === 'vs-light' ? '切换深色主题' : '切换浅色主题'"
          :title="editorTheme === 'vs-light' ? '切换深色主题' : '切换浅色主题'"
          @click="toggleTheme"
        >
          <i class="bi" :class="editorTheme === 'vs-light' ? 'bi-moon-stars' : 'bi-sun'" aria-hidden="true"></i>
        </n-button>
      </n-space>
    </template>

    <!-- 编辑器容器：使用flex布局确保占满可用空间 -->
    <div ref="editorContainer" class="editor-container mb-3">
      <MonacoEditor
        v-if="editorReady"
        v-model:value="internalCode"
        :language="monacoLanguage"
        :theme="editorTheme"
        height="100%"
        :options="editorOptions"
      />
    </div>

    <div class="d-flex justify-content-start align-items-center mb-3">
      <n-button
        circle
        tertiary
        aria-label="重置代码"
        title="重置代码"
        @click="resetCode"
      >
        <i class="bi bi-arrow-counterclockwise" aria-hidden="true"></i>
      </n-button>
    </div>

    <SampleTest
      :activeStatus="activeStatus"
      :testSample="testSample"
      :outputValue="outputValue"
      @handleActiveStatus="handleActiveStatus"
      @handleTestSample="handleTestSample"
      @handleRunExample="handleRunExample"
    >
      <template #actions>
        <n-button
          type="primary"
          class="editor-submit-button"
          @click="submitCode"
          :disabled="isLoading"
          :class="{ 'opacity-50': isLoading }"
        >
          <i class="bi bi-cloud-arrow-up me-1" aria-hidden="true"></i>
          {{ isLoading ? '提交中...' : '提交代码' }}
        </n-button>
      </template>
    </SampleTest>
  </n-card>
</template>

<script setup>
import { computed, ref, watch, onMounted, onUnmounted, defineAsyncComponent } from 'vue'
import { NCard, NSelect, NButton, NSpace, useMessage } from 'naive-ui'
import { useIntersectionObserver } from '@vueuse/core'
import api from '@/api'
import SampleTest from './SampleTest.vue'
import { useWebSocketContext } from '@/composables/useWebSocket.js'
import { defaultCodeFor, toLanguageOptions, toMonacoLanguage } from '@/utils/languageOptions'

const MonacoEditor = defineAsyncComponent(() => import('monaco-editor-vue3'))

// props + emit
const props = defineProps({
  problemId: Number,
  roomId: String,
  languages: {
    type: [Array, Object, String],
    default: () => []
  },
  theme: {
    type: String,
    default: 'vs-light'
  }
})

const emit = defineEmits(['update:theme'])

const langOptions = computed(() => toLanguageOptions(props.languages))

const internalLang = ref('')
const internalCode = ref('')
const editorReady = ref(false)
const editorContainer = ref(null)
const userEdited = ref(false)
const message = useMessage()
const isEditorVisible = ref(true)
const testSample = ref(true)
const isLoading = ref(false)
const activeStatus = ref('')
const outputValue = ref('')
const { registerSubmitCodeCallback } = useWebSocketContext()
const editorTheme = computed(() => props.theme)
const monacoLanguage = computed(() => toMonacoLanguage(internalLang.value))

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
  internalCode.value = defaultCodeFor(internalLang.value)
  userEdited.value = false
  message.info('代码已重置')
}

function toggleTheme() {
  emit('update:theme', editorTheme.value === 'vs-light' ? 'vs-dark' : 'vs-light')
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
  internalCode.value = defaultCodeFor(newLang)
})

watch(langOptions, (options) => {
  if (!options.length) {
    internalLang.value = ''
    internalCode.value = ''
    return
  }

  if (!options.some(option => option.value === internalLang.value)) {
    internalLang.value = options[0].value
  }
}, { immediate: true })

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
  internalCode.value = defaultCodeFor(internalLang.value)
  
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
.editor-container {
  width: 100%;
  height: clamp(300px, 50vh, 520px);
}

:deep(.n-card__content) {
  display: flex;
  flex-direction: column;
}

.editor-submit-button {
  min-width: 8.4rem;
}
</style>
