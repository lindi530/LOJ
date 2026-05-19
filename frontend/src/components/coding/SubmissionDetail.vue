<template>
    <div class="container mt-4 custom-container">
    <n-card size="large">

            <!-- 题目展示 -->
        <div class="problem-info mb-4">
          <h2 class="problem-title">{{ detailData.title }}</h2>
          <n-card>
            <p class="problem-content" v-html="detailData.description"></p>
          </n-card>
        </div>
        
        <!-- 标题区域 -->
        <div class="title">
          <router-link
            :to="`/users/${detailData?.user_id}`"
            class="text-decoration-none text-dark"
          >
          <span class="username" style="color: #4183c4;">{{ detailData.user_name }}</span>
          </router-link>
           提交的代码
        </div>
        <!-- 信息列表 -->
        <ul class="info-list">
            <li>提交时间：{{ formatTimeToYMDHMS(detailData.created_at) }}</li>
            <li>语言：{{ detailData.language }}</li>
            <li>代码长度：{{ codeLength }}</li>
            <li>运行时间：{{ detailData.exec_time }} ms</li>
            <li>占用内存：{{ detailData.memory_usage }} MB</li>
            <li>得分：{{ detailData.score }}</li>
            <li class="state">
                运行状态：
                <span :style="{ color: getStateColor(detailData.state) }">
                    {{ detailData.state }}
                </span>
            </li>
        </ul>
        <!-- 代码展示 -->
        <n-card  size="medium">
          <div style="overflow: auto">
            <n-code 
              v-if="detailData.language"
              :code="detailData.code" 
              :language="detailData.language" 
              show-line-numbers />
          </div>
        </n-card>
        

    </n-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
// 假设 api 模块，实际按项目结构调整
import api from '@/api'; 

const route = useRoute();
const submissionId = route.params.submission_id;
const detailData = ref({});

onMounted(async () => {
  try {
    const resp = await api.getSubmissionDetail(submissionId);
    console.log(resp.data)
    if (resp.code === 0) {
      detailData.value = resp.data,

    console.log("detailData: ", detailData.value)
    }
  } catch (error) {
    console.error('获取提交详情失败：', error);
  }
});

const codeLength = computed(() => {
    return detailData.value?.code?.length || 0
})

function getStateColor(state) {
  switch (state) {
    case 'Accepted':
      return 'green'
    case 'Wrong Answer':
      return 'red'
    case 'Runtime Error':
      return 'orange'
    case 'Running':
      return 'blue'
    default:
      return 'black'
  }
}

// 时间格式化（保持原逻辑）
const formatTimeToYMDHMS = (isoTime) => {
  if (!isoTime) return '';
  const date = new Date(isoTime);
  return [
    date.getFullYear(),
    String(date.getMonth() + 1).padStart(2, '0'),
    String(date.getDate()).padStart(2, '0')
  ].join('-') + ' ' + [
    String(date.getHours()).padStart(2, '0'),
    String(date.getMinutes()).padStart(2, '0'),
    String(date.getSeconds()).padStart(2, '0')
  ].join(':');
};
</script>

<style scoped>

.problem-title {
  font-size: 1.5rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.problem-content {
  font-size: 1rem;
  line-height: 1.5;
  color: #333;
}
.custom-container {
  max-width: 50%;   /* 宽度可调 */
  margin: 0 auto;     /* 居中 */
}

.submission-detail {
  padding: 16px;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin: 16px 0;
}
.title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 8px;
}
.info-list {
  list-style: none;
  padding: 0;
  margin: 0 0 12px 0;
}
.info-list li {
  margin: 4px 0;
  color: #666;
}
.code-box {
  background: #f8f8f8;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 12px;
  overflow: auto;
}
code {
  font-family: Consolas, Monaco, monospace;
  font-size: 14px;
  color: #333;
}
</style>