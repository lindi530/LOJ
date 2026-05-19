<template>
  <div
    class="container px-5 py-4"
    :class="{ 'custom-container': roomId === undefined }"
    style="height: 100%; overflow-y: auto;"
  >
    <div class="row justify-content-center">
      <div class="col-12">
        <div class="d-flex" style="gap: 24px;">
          <!-- 左侧题目区域 -->
          <div class="flex-fill" style="max-width: 50%;">

            <n-card content-style="padding: 0;">
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
                    class="p-4"
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

          <!-- 右侧编辑器+测试 -->
          <div class="flex-fill" style="max-width: 50%;">
            <CodeEditor 
              :problem-id="problemID"
              :room-id="roomId"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import ProblemDetail from '@/components/coding/ProblemDetail.vue'
import CodeEditor from '@/components/coding/CodeEditor.vue'
import SubmissionList from '@/components/coding/SubmissionList.vue'
import api from '@/api'

// 当前题目
const problem = ref(null)
const loading = ref(true)
const route = useRoute()


const props = defineProps({
  problemId: {
    type: Number,
    required: true,
  },
  roomId: {
    type: String,
    required: true,
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
.custom-container {
  max-width: 85%;   /* 宽度可调 */
  margin: 0 auto;     /* 居中 */
}

/* 移除全局样式污染，避免影响父组件 */
/* .page-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative; /* 确保在父级堆叠上下文中 
}

/* 确保内容区域不被裁剪 
.container-fluid {
  padding: 0;
  margin: 0;
  height: 100%;
} */
</style>