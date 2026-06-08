<template>
  <section
    class="contest-editor"
    :class="{ 'contest-editor--light': editorTheme === 'vs-light' }"
    :data-bs-theme="bootstrapTheme"
    :aria-label="ariaLabel"
  >
    <header class="contest-editor__toolbar">
      <div class="editor-file">
        <i class="bi bi-code-slash" aria-hidden="true"></i>
        <span>{{ currentFile }}</span>
        <small>{{ modeLabel }}</small>
      </div>
      <div class="editor-actions">
        <label class="sr-only" for="problem-code-language">选择提交语言</label>
        <select id="problem-code-language" v-model="language" class="form-select form-select-sm language-select" :disabled="!languageOptions.length">
          <option v-for="option in languageOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>
        <button
          class="toolbar-button toolbar-icon-button"
          type="button"
          :aria-label="editorTheme === 'vs-dark' ? '切换浅色主题' : '切换深色主题'"
          :title="editorTheme === 'vs-dark' ? '切换浅色主题' : '切换深色主题'"
          @click="toggleTheme"
        >
          <i class="bi" :class="editorTheme === 'vs-dark' ? 'bi-sun' : 'bi-moon-stars'" aria-hidden="true"></i>
        </button>
        <button
          class="toolbar-button toolbar-icon-button"
          type="button"
          aria-label="重置代码"
          title="重置代码"
          @click="resetCode"
        >
          <i class="bi bi-arrow-counterclockwise" aria-hidden="true"></i>
        </button>
      </div>
    </header>

    <div ref="editorContainer" class="contest-editor__canvas">
      <MonacoEditor
        v-if="editorReady"
        v-model:value="code"
        :language="monacoLanguage"
        :theme="editorTheme"
        height="100%"
        :options="editorOptions"
      />
      <div v-else class="editor-loading" role="status">
        <span aria-hidden="true"></span>
        正在准备编辑器
      </div>
    </div>

    <section class="test-console" aria-label="运行与提交结果">
      <header class="test-console__header">
        <div class="test-console__tabs" role="tablist" aria-label="结果面板">
          <button
            type="button"
            role="tab"
            :aria-selected="testSample"
            :class="{ 'is-active': testSample }"
            @click="testSample = true"
          >
            自定义测试
          </button>
          <button
            type="button"
            role="tab"
            :aria-selected="!testSample"
            :class="{ 'is-active': !testSample }"
            @click="testSample = false"
          >
            提交状态
          </button>
        </div>
        <span class="judge-status" :class="statusClass" role="status">{{ displayedStatus }}</span>
      </header>

      <div v-if="testSample" class="test-console__runner" role="tabpanel">
        <label class="console-field console-field--input">
          <span class="console-label">输入</span>
          <textarea v-model="sampleInput" rows="3" placeholder="输入自定义测试数据"></textarea>
        </label>
        <div class="console-field console-output">
          <span class="console-label">输出</span>
          <pre>{{ outputValue || '运行后将在这里显示输出。' }}</pre>
        </div>
      </div>

      <div v-else class="submission-status" role="tabpanel">
        <i class="bi bi-activity" aria-hidden="true"></i>
        <div>
          <strong :class="statusClass">{{ displayedStatus }}</strong>
          <p>{{ submissionMessage || '提交代码后，评测进度与结果会显示在这里。' }}</p>
        </div>
      </div>

      <div class="console-actions">
        <button class="run-button" type="button" :disabled="runningExample" @click="runExample">
          <span v-if="runningExample" class="button-spinner" aria-hidden="true"></span>
          <i v-else class="bi bi-play-fill" aria-hidden="true"></i>
          {{ runningExample ? '运行中' : '运行样例' }}
        </button>
        <button class="submit-button" type="button" :disabled="submitting" @click="submitCode">
          <span v-if="submitting" class="button-spinner" aria-hidden="true"></span>
          <i v-else class="bi bi-cloud-arrow-up" aria-hidden="true"></i>
          {{ submitting ? '提交中...' : '提交代码' }}
        </button>
      </div>
    </section>

  </section>
</template>

<script setup>
import { computed, defineAsyncComponent, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useIntersectionObserver } from '@vueuse/core'
import { useMessage } from 'naive-ui'
import { useWebSocketContext } from '@/composables/useWebSocket.js'
import { defaultCodeFor, fileNameFor, toLanguageOptions, toMonacoLanguage } from '@/utils/languageOptions'

const MonacoEditor = defineAsyncComponent(() => import('monaco-editor-vue3'))

const props = defineProps({
  ariaLabel: {
    type: String,
    default: '代码提交'
  },
  initialInput: {
    type: String,
    default: ''
  },
  languages: {
    type: [Array, Object, String],
    default: () => []
  },
  theme: {
    type: String,
    default: 'vs-light'
  },
  judgeMode: {
    type: String,
    default: 'ACM'
  },
  canRunExample: {
    type: Boolean,
    default: true
  },
  runExampleRequest: {
    type: Function,
    default: null
  },
  submitRequest: {
    type: Function,
    required: true
  }
})

const emit = defineEmits(['update:theme'])

const language = ref('')
const code = ref('')
const editorReady = ref(false)
const editorContainer = ref(null)
const testSample = ref(true)
const sampleInput = ref('')
const outputValue = ref('')
const activeStatus = ref('')
const submissionMessage = ref('')
const submitting = ref(false)
const runningExample = ref(false)
const message = useMessage()
const { registerSubmitCodeCallback } = useWebSocketContext()
let stopObserver = () => {}

const languageOptions = computed(() => toLanguageOptions(props.languages))
const editorTheme = computed(() => props.theme)
const bootstrapTheme = computed(() => editorTheme.value === 'vs-dark' ? 'dark' : 'light')
const monacoLanguage = computed(() => toMonacoLanguage(language.value))
const currentFile = computed(() => fileNameFor(language.value))
const modeLabel = computed(() => props.judgeMode ? `${props.judgeMode} 模式` : '')
const displayedStatus = computed(() => activeStatus.value || (testSample.value ? '待运行' : '等待提交'))
const statusClass = computed(() => {
  const status = activeStatus.value.toLowerCase()
  if (status.includes('accepted') || status.includes('finished') || status.includes('通过')) {
    return 'judge-status--success'
  }
  if (status.includes('running') || status.includes('compiling') || status.includes('运行') || status.includes('编译')) {
    return 'judge-status--running'
  }
  if (status.includes('pending') || status.includes('等待') || status.includes('提交中')) {
    return 'judge-status--pending'
  }
  if (activeStatus.value) {
    return 'judge-status--failure'
  }
  return ''
})
const editorOptions = {
  automaticLayout: true,
  fontSize: 14,
  lineHeight: 24,
  minimap: { enabled: false },
  padding: { top: 18, bottom: 12 },
  renderLineHighlight: 'line',
  scrollBeyondLastLine: false,
  smoothScrolling: true
}

const unregister = registerSubmitCodeCallback((statusMessage) => {
  if (!statusMessage.content) {
    return
  }

  testSample.value = false
  activeStatus.value = statusMessage.content
  submissionMessage.value = '评测状态已更新。'
})

function isRequestTimeout(error) {
  return error?.code === 'ECONNABORTED' || String(error?.message || '').toLowerCase().includes('timeout')
}

function clearExecutionFeedback() {
  activeStatus.value = ''
  outputValue.value = ''
  submissionMessage.value = ''
  testSample.value = true
}

function resetCode() {
  code.value = defaultCodeFor(language.value)
  clearExecutionFeedback()
  message.info('代码已重置')
}

function toggleTheme() {
  emit('update:theme', editorTheme.value === 'vs-dark' ? 'vs-light' : 'vs-dark')
}

function canRunExample() {
  return props.canRunExample && typeof props.runExampleRequest === 'function'
}

async function runExample() {
  if (runningExample.value) {
    return
  }

  if (!canRunExample()) {
    message.warning('当前题目暂不支持自定义测试')
    return
  }

  if (!code.value.trim()) {
    message.warning('请输入代码后再运行')
    return
  }

  runningExample.value = true
  testSample.value = true
  activeStatus.value = ''
  outputValue.value = ''
  try {
    const resp = await props.runExampleRequest({
      input: sampleInput.value,
      language: language.value,
      code: code.value
    })

    if (isFailureResponse(resp)) {
      const failure = resp.message || '运行失败，请稍后重试。'
      activeStatus.value = 'Failed'
      outputValue.value = failure
      message.error(failure)
      return
    }

    outputValue.value = resp?.data?.output ?? resp?.output ?? ''
  } catch (requestError) {
    if (isRequestTimeout(requestError)) {
      outputValue.value = '请求超时，运行状态仍以 WebSocket 推送为准。'
      message.warning('请求超时，运行状态将继续通过实时推送更新')
      return
    }

    const failure = requestError?.message || '运行失败，请稍后重试。'
    activeStatus.value = 'Failed'
    outputValue.value = failure
    message.error(failure)
  } finally {
    runningExample.value = false
  }
}

async function submitCode() {
  if (submitting.value) {
    return
  }

  if (!code.value.trim()) {
    message.warning('请输入代码后再提交')
    return
  }

  submitting.value = true
  testSample.value = false
  activeStatus.value = ''
  submissionMessage.value = '代码已发送，正在等待评测状态推送。'
  try {
    const resp = await props.submitRequest({
      language: language.value,
      code: code.value
    })

    if (isFailureResponse(resp)) {
      const failure = resp.message || '提交失败，请稍后重试。'
      activeStatus.value = 'Failed'
      submissionMessage.value = failure
      message.error(failure)
      return
    }

    submissionMessage.value = '提交成功，评测状态将通过实时推送更新。'
    message.success('提交成功')
  } catch (requestError) {
    if (isRequestTimeout(requestError)) {
      submissionMessage.value = '提交请求超时，评测状态仍以 WebSocket 推送为准。'
      message.warning('提交请求超时，评测状态将继续通过实时推送更新')
      return
    }

    const failure = requestError?.message || '提交失败，请稍后重试。'
    activeStatus.value = 'Failed'
    submissionMessage.value = failure
    message.error(failure)
  } finally {
    submitting.value = false
  }
}

function isFailureResponse(resp) {
  return resp && Object.prototype.hasOwnProperty.call(resp, 'code') && resp.code !== 0
}

watch(language, (selectedLanguage) => {
  code.value = defaultCodeFor(selectedLanguage)
  clearExecutionFeedback()
})

watch(languageOptions, (options) => {
  if (!options.length) {
    language.value = ''
    code.value = ''
    return
  }

  if (!options.some(option => option.value === language.value)) {
    language.value = options[0].value
  }
}, { immediate: true })

watch(() => props.initialInput, (nextInput) => {
  sampleInput.value = nextInput || ''
}, { immediate: true })

onMounted(() => {
  const observer = useIntersectionObserver(
    editorContainer,
    ([entry]) => {
      if (entry?.isIntersecting && !editorReady.value) {
        editorReady.value = true
      }
    }
  )
  stopObserver = observer.stop
})

onBeforeUnmount(() => {
  stopObserver()
  unregister()
})
</script>

<style scoped>
.contest-editor {
  display: flex;
  flex-direction: column;
  min-width: 0;
  height: 100%;
  overflow: hidden;
  color: #e5edf7;
  font-family: Inter, "PingFang SC", "Microsoft YaHei", Arial, sans-serif;
  background: #0b1220;
  border: 1px solid #182337;
  border-radius: 1rem;
  box-shadow: 0 0.65rem 1.8rem rgba(15, 23, 42, 0.14);
}

.contest-editor__toolbar {
  display: flex;
  gap: 1rem;
  align-items: center;
  justify-content: space-between;
  min-height: 62px;
  padding: 0 1.05rem;
  background: #0e1728;
  border-bottom: 1px solid #1d283b;
}

.editor-file {
  display: inline-flex;
  gap: 0.62rem;
  align-items: center;
  min-width: 0;
  color: #dce6f3;
  font-size: 0.86rem;
  font-weight: 540;
}

.editor-file i {
  color: #1cc7b2;
  font-size: 1rem;
}

.editor-file small {
  padding-left: 0.65rem;
  color: #71839c;
  font-size: 0.72rem;
  font-weight: 530;
  border-left: 1px solid #263349;
}

.editor-actions {
  display: inline-flex;
  gap: 0.45rem;
  align-items: center;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  overflow: hidden;
  clip-path: inset(50%);
  white-space: nowrap;
}

.language-select,
.toolbar-button {
  height: 36px;
  color: #cad5e5;
  font-size: 0.8rem;
  font-weight: 500;
  background: #131e31;
  border: 1px solid #223046;
  border-radius: 0.58rem;
  transition: border-color 150ms ease, background 150ms ease, color 150ms ease;
}

.language-select {
  min-width: 7.6rem;
  padding: 0 0.7rem;
}

.toolbar-button {
  display: inline-flex;
  gap: 0.42rem;
  align-items: center;
  justify-content: center;
  padding: 0 0.72rem;
}

.toolbar-icon-button {
  width: 36px;
  padding: 0;
}

.toolbar-button:hover {
  color: #eef5ff;
  background: #18263b;
  border-color: #304259;
}

.language-select:focus-visible,
.toolbar-button:focus-visible {
  outline: 2px solid rgba(21, 184, 166, 0.68);
  outline-offset: 1px;
}

.contest-editor__canvas {
  flex: 1 1 auto;
  min-height: 20rem;
  overflow: hidden;
  background: #0d1525;
}

.contest-editor--light .contest-editor__canvas {
  background: #fff;
}

.contest-editor--light {
  color: #1f2937;
  background: #fff;
  border-color: #dfe7f0;
  box-shadow: 0 0.65rem 1.8rem rgba(15, 23, 42, 0.08);
}

.contest-editor--light .contest-editor__toolbar,
.contest-editor--light .test-console {
  background: #f8fafc;
  border-color: #e2e8f0;
}

.contest-editor--light .editor-file {
  color: #263446;
}

.contest-editor--light .editor-file small,
.contest-editor--light .console-field span,
.contest-editor--light .console-output > span,
.contest-editor--light .judge-status,
.contest-editor--light .submission-status,
.contest-editor--light .submission-status p {
  color: #64748b;
}

.contest-editor--light .language-select,
.contest-editor--light .toolbar-button,
.contest-editor--light .console-field textarea,
.contest-editor--light .console-output pre {
  color: #1f2937;
  background: #fff;
  border-color: #d8e1ec;
}

.contest-editor--light .toolbar-button:hover {
  color: #0f766e;
  background: #eef6ff;
  border-color: #bcd5f3;
}

.contest-editor--light .test-console__tabs button.is-active {
  color: #0f172a;
}

.contest-editor:not(.contest-editor--light) :deep(.monaco-editor),
.contest-editor:not(.contest-editor--light) :deep(.monaco-editor-background),
.contest-editor:not(.contest-editor--light) :deep(.monaco-editor .margin) {
  background-color: #0d1525 !important;
}

.editor-loading {
  display: flex;
  gap: 0.7rem;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #75879e;
  font-size: 0.84rem;
}

.editor-loading span {
  width: 0.95rem;
  height: 0.95rem;
  border: 2px solid #243247;
  border-top-color: #15b8a6;
  border-radius: 50%;
  animation: editor-spin 700ms linear infinite;
}

.test-console {
  flex: 0 0 auto;
  background: #0d1627;
  border-top: 1px solid #202d43;
}

.test-console__header {
  display: flex;
  gap: 1rem;
  align-items: center;
  justify-content: space-between;
  height: 47px;
  padding: 0 1rem;
}

.test-console__tabs {
  display: inline-flex;
  gap: 1.1rem;
  height: 100%;
}

.test-console__tabs button {
  position: relative;
  color: #778aa2;
  font-size: 0.8rem;
  font-weight: 550;
  background: transparent;
  border: 0;
}

.test-console__tabs button.is-active {
  color: #f1f6fd;
}

.test-console__tabs button.is-active::after {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  height: 2px;
  background: #15b8a6;
  content: "";
}

.judge-status {
  color: #8293aa;
  font-size: 0.76rem;
  font-weight: 560;
}

.judge-status--success {
  color: #27c7a6;
}

.judge-status--running {
  color: #4ea4ff;
}

.judge-status--pending {
  color: #f1b74e;
}

.judge-status--failure {
  color: #fa7782;
}

.test-console__runner {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.72rem;
  align-items: start;
  padding: 0 1rem 0.72rem;
}

.console-field--input {
  min-width: 0;
}

.console-output {
  min-width: 0;
}

.console-label {
  display: block;
  margin-bottom: 0.38rem;
  color: #73859e;
  font-size: 0.69rem;
  font-weight: 600;
}

.console-field textarea,
.console-output pre {
  width: 100%;
  height: 88px;
  padding: 0.68rem 0.75rem;
  margin: 0;
  color: #dce6f4;
  font: 0.78rem/1.55 "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  background: #111c30;
  border: 1px solid #202e44;
  border-radius: 0.62rem;
}

.console-field textarea {
  resize: none;
}

.console-field textarea:focus-visible {
  outline: 2px solid rgba(21, 184, 166, 0.52);
  outline-offset: 0;
}

.console-output pre {
  overflow: auto;
  color: #b8c6d8;
  white-space: pre-wrap;
}

.console-actions {
  display: flex;
  gap: 0.62rem;
  justify-content: flex-end;
  padding: 0 1rem 0.92rem;
}

.run-button,
.submit-button {
  display: inline-flex;
  gap: 0.42rem;
  align-items: center;
  justify-content: center;
  color: #f4fffd;
  font-size: 0.84rem;
  font-weight: 590;
  background: #15b8a6;
  border: 1px solid #15b8a6;
  border-radius: 0.6rem;
  transition: background 150ms ease, border-color 150ms ease, transform 150ms ease;
}

.run-button {
  height: 40px;
  width: 8.4rem;
  padding: 0 0.95rem;
}

.submit-button {
  height: 40px;
  width: 8.4rem;
  padding: 0 1rem;
}

.run-button:hover:not(:disabled),
.submit-button:hover:not(:disabled) {
  background: #0fa795;
  border-color: #0fa795;
  transform: translateY(-1px);
}

.run-button:disabled,
.submit-button:disabled {
  cursor: not-allowed;
  opacity: 0.66;
}

.run-button:focus-visible,
.submit-button:focus-visible {
  outline: 3px solid rgba(21, 184, 166, 0.3);
  outline-offset: 2px;
}

.submission-status {
  display: flex;
  gap: 0.92rem;
  align-items: center;
  min-height: 101px;
  padding: 0.2rem 1.15rem 1rem;
  color: #8395ad;
}

.submission-status > i {
  color: #52657f;
  font-size: 1.2rem;
}

.submission-status strong {
  display: block;
  margin-bottom: 0.38rem;
  font-size: 0.92rem;
  font-weight: 600;
}

.submission-status p {
  margin: 0;
  color: #778aa3;
  font-size: 0.78rem;
}

.button-spinner {
  width: 0.85rem;
  height: 0.85rem;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: editor-spin 650ms linear infinite;
}

@keyframes editor-spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 1199.98px) {
  .toolbar-button i {
    font-size: 0.9rem;
  }
}

@media (max-width: 959.98px) {
  .contest-editor {
    height: auto;
    border-radius: 0.82rem;
  }

  .contest-editor__toolbar {
    flex-wrap: wrap;
    min-height: 94px;
    padding: 0.74rem;
  }

  .editor-actions {
    width: 100%;
  }

  .language-select {
    flex: 1;
  }

  .toolbar-button {
    font-size: 0.77rem;
  }

  .contest-editor__canvas {
    flex-basis: auto;
    height: 25rem;
  }

  .test-console__runner {
    grid-template-columns: 1fr;
  }

  .console-actions {
    flex-direction: column;
  }

  .run-button,
  .submit-button {
    width: 100%;
  }
}
</style>
