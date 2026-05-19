<template>

<!-- 表格核心：还原 <thead> 和 <tbody> 的多列结构 -->
<table v-if="submissionRecords.length !== 0">
  <!-- 表头：明确 5 列（提交时间、得分、状态、运行时间、语言） -->
  <thead class="bg-gray-50">
    <tr>
      <th class="px-4 py-3 text-center  text-sm font-semibold text-gray-700">提交编号</th>
      <th class="px-4 py-3 text-center  text-sm font-semibold text-gray-700">得分</th>
      <th class="px-4 py-3 text-center  text-sm font-semibold text-gray-700">状态</th>
      <th class="px-4 py-3 text-center  text-sm font-semibold text-gray-700">运行时间</th>
      <th class="px-4 py-3 text-center  text-sm font-semibold text-gray-700">语言</th>
      <th class="px-4 py-3 text-center  text-sm font-semibold text-gray-700">提交时间</th>
    </tr>
  </thead>
  <!-- 内容：循环渲染每一条记录，拆分成 5 列 -->
  <tbody class="bg-white divide-y divide-gray-100">
    <tr 
      v-for="(record, index) in submissionRecords" 
      :key="record.id"
      class="hover:bg-gray-50 transition-colors fade-in"
      :style="{ animationDelay: `${index * 0.1}s` }"
      @click="viewRecord(record)"
    >
      <td class="px-4 py-4 whitespace-nowrap text-center text-sm text-gray-500">
        {{ record.id }}
      </td>

      <!-- 得分列 -->
      <td class="px-4 py-4 whitespace-nowrap text-center text-sm font-medium" :class="getScoreColor(record.score)">
        {{ record.score }}
      </td>
      <!-- 状态列 -->
      <td class="px-4 py-4 whitespace-nowrap">
        <!-- <span class="status-badge" :class="getStatusClass(record.state)">
          <i :class="getStatusIcon(record.state)" class="mr-1"></i>{{ record.state }}
        </span> -->
        <router-link
          :to="{
            name: 'SubmissionDetail',
            params: {
              submission_id: record.id,
            },
          }"
          target="_blank"
          class="status-badge"
          :class="getStatusClass(record.state)"
        >
          <i :class="getStatusIcon(record.state)" class="mr-1 text-center"></i>{{ record.state }}
        </router-link>
      </td>
      <!-- 运行时间列 -->
      <td class="px-4 py-4 whitespace-nowrap text-center text-sm text-gray-500">
        {{ record.exec_time }}
      </td>
      <!-- 语言列 -->
      <td class="px-4 py-4 whitespace-nowrap text-sm" >
        <span class="px-2 py-1 rounded text-center text-xs" :class="getLanguageClass(record.language)">
          {{ record.language }}
        </span>
      </td>
                        <!-- 提交时间列 -->
      <td class="px-4 py-4 whitespace-nowrap text-center  text-sm text-gray-500">
        {{ formatTimeToYMDHMS(record.created_at) }}
      </td>
    </tr>
  </tbody>
</table>
<n-card v-else title="暂无提交记录"></n-card>
</template>

<script setup>
import api from '@/api/index';
import { ref, onMounted  } from 'vue';

const props = defineProps({
  roomId: {
    type: Number,
    required: true
  },
  problemId: {
    type: Number,
    required: true
  },
  tab: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:tab'])

// 提交记录数据
const submissionRecords = ref([]);

onMounted(async () => {
  try {
    const response = await api.getProblemSubmissionList(props.problemId); 
    console.log("responser: ", response)
    if (response.code === 0) {
      submissionRecords.value = response.data;
    } else {
      console.error('请求失败:', response.msg || '未知错误');
    }
  } catch (err) {
    console.error('请求异常:', err);
  }
});

// 查看记录详情
const viewRecord = (record) => {
  console.log('查看提交记录:', record);
};

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

// 分数颜色（保持原逻辑）
const getScoreColor = (score) => {
  if (score === 100) return 'text-success';
  if (score > 0) return 'text-danger';
  return 'text-warning';
};

// 状态样式（保持原逻辑）
const getStatusClass = (status) => {
  switch (status) {
    case '通过':
    case 'Accepted':
      return 'bg-success/10 text-success';
    case '答案错误':
    case 'Wrong':
      return 'bg-danger/10 text-danger';
    case '时间超限':
      return 'bg-warning/10 text-warning';
    case '编译错误':
      return 'bg-orange-100 text-orange-800';
    case '部分通过':
      return 'bg-blue-100 text-blue-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

// 状态图标（保持原逻辑）
const getStatusIcon = (status) => {
  switch (status) {
    case '通过':
      return 'fa fa-check-circle';
    case '答案错误':
      return 'fa fa-times-circle';
    case '时间超限':
      return 'fa fa-clock-o';
    case '编译错误':
      return 'fa fa-exclamation-triangle';
    case '部分通过':
      return 'fa fa-check-circle-half-o';
    default:
      return 'fa fa-question-circle';
  }
};

// 语言样式（保持原逻辑）
const getLanguageClass = (language) => {
  if (!language) return 'bg-gray-100 text-gray-800';
  const lang = language.toLowerCase();
  const languageClasses = {
    'python': 'bg-blue-100 text-blue-800',
    'javascript': 'bg-purple-100 text-purple-800',
    'java': 'bg-green-100 text-green-800',
    'cpp': 'bg-red-100 text-red-800',
    'c': 'bg-red-100 text-red-800',
    'go': 'bg-blue-400 text-blue-900'
  };
  return languageClasses[lang] || 'bg-gray-100 text-gray-800';
};
</script>

<style scoped>
/* 保持原有样式，可根据需求微调列宽 */
.status-badge {
  @apply px-2 py-1 rounded-full text-xs font-medium;
}
.bg-blue-100 { background-color: rgba(59, 130, 246, 0.1); }
.text-blue-800 { color: rgb(30, 64, 175); }
.bg-purple-100 { background-color: rgba(139, 92, 246, 0.1); }
.text-purple-800 { color: rgb(76, 29, 149); }
.bg-green-100 { background-color: rgba(16, 185, 129, 0.1); }
.text-green-800 { color: rgb(22, 93, 66); }
.bg-red-100 { background-color: rgba(239, 68, 68, 0.1); }
.text-red-800 { color: rgb(127, 29, 29); }
.bg-blue-400 { background-color: rgba(59, 130, 246, 0.4); }
.text-blue-900 { color: rgb(23, 37, 84); }
.bg-gray-100 { background-color: rgba(243, 244, 246, 0.5); }
.text-gray-800 { color: rgb(31, 35, 40); }
.bg-orange-100 { background-color: rgba(249, 115, 22, 0.1); }
.text-orange-800 { color: rgb(144, 49, 5); }
.text-success { color: #00B42A; }
.text-danger { color: #F53F3F; }
.text-warning { color: #FF7D00; }
.fade-in {
  animation: fadeIn 0.3s ease-in-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

</style>