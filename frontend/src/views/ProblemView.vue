<template>
  <div
    class="problem-page container-fluid px-4 py-4"
    :class="pageThemeClass"
    :data-bs-theme="bootstrapTheme"
  >
    <div class="row g-4 justify-content-center">
      <div class="col-12" :class="roomId === undefined ? 'col-xxl-11' : ''">
        <div class="row g-4">
          <div class="col-12 col-xl-6">

            <n-card class="problem-surface" content-style="padding: 0;">
              <n-tabs
                type="line"
                size="large"
                :tabs-padding="20"
                pane-style="padding: 20px;"
              >
                <n-tab-pane name="题目信息">
                  <ProblemDetail
                    v-if="problem"
                    :problem="problem" 
                    :room-id="roomId"
                  />
                </n-tab-pane>
                <n-tab-pane name="提交记录"
                  v-if="activeTab">
                  <SubmissionList 
                    :problem-id="problemID"
                    :room-id="roomId"
                    class="p-4"
                  />
                </n-tab-pane>
              </n-tabs>
            </n-card>
          </div>

          <div class="col-12 col-xl-6">
            <CodeEditor 
              :problem-id="problemID"
              :room-id="roomId"
              :languages="problemLanguageSource"
              :theme="editorTheme"
              @update:theme="editorTheme = $event"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import ProblemDetail from '@/components/coding/ProblemDetail.vue'
import CodeEditor from '@/components/coding/CodeEditor.vue'
import SubmissionList from '@/components/coding/SubmissionList.vue'
import api from '@/api'

// 当前题目
const problem = ref(null)
const loading = ref(true)
const route = useRoute()
const editorTheme = ref('vs-light')
const bootstrapTheme = computed(() => editorTheme.value === 'vs-dark' ? 'dark' : 'light')
const pageThemeClass = computed(() => editorTheme.value === 'vs-dark' ? 'problem-page--dark' : 'problem-page--light')
const problemLanguageSource = computed(() => problem.value?.language || problem.value?.languages || problem.value?.constraints || [])


const props = defineProps({
  problemId: {
    type: Number,
    default: undefined,
  },
  roomId: {
    type: String,
    default: undefined,
  },

})

// 初始化：roomId 存在则 activeTab 为 false，否则为 true
const activeTab = ref(!props.roomId || props.roomId.trim() === '');

// 监听 roomId 变化，根据新值更新 activeTab
watch(
  () => props.roomId,
  (newRoomId) => {
    // 核心逻辑：新的 roomId 存在（非空且非纯空格）→ activeTab 为 false
    activeTab.value = !newRoomId || newRoomId.trim() === '';
  }
);

watch(editorTheme, (theme) => {
  document.documentElement.dataset.lojEditorTheme = theme === 'vs-dark' ? 'dark' : 'light'
}, { immediate: true })

onBeforeUnmount(() => {
  delete document.documentElement.dataset.lojEditorTheme
})

// 修复：统一转换为Number类型，避免类型不匹配
const problemID = ref(
  typeof route.params.problem_id !== 'undefined' 
    ? Number(route.params.problem_id) 
    : props.problemId
);

onMounted(async () => {
  try {
    const res = await api.getProblemDetail(problemID.value)
    if (res.code === 0) {
      problem.value = res.data
    } else {
      console.error('获取题目失败:', res.data)
    }
  } catch (err) {
    console.error('请求错误:', err)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.problem-page {
  min-height: calc(100vh - 60px);
  overflow-y: auto;
  transition: background-color 160ms ease, color 160ms ease;
}

.problem-page--light {
  color: #1f2937;
  background: #f6f8fb;
}

.problem-page--dark {
  color: #dbe7f5;
  background: #0b1220;
}

.problem-page--dark :deep(.n-card),
.problem-page--dark :deep(.n-tabs-pane-wrapper) {
  color: #dbe7f5;
  background-color: #111827;
  border-color: #243247;
}

.problem-page--dark :deep(.n-card-header),
.problem-page--dark :deep(.n-tabs-tab__label),
.problem-page--dark :deep(.n-tabs-nav-scroll-content) {
  color: #edf4ff;
}

.problem-page--dark :deep(.n-card .n-card),
.problem-page--dark :deep(.bg-body-tertiary) {
  background-color: #162235 !important;
  border-color: #2c3a50 !important;
}

.problem-page--dark :deep(.text-body-secondary) {
  color: #9fb0c5 !important;
}
</style>
