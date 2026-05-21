<template>
  <div v-if="submissions.length" class="submission-table">
    <div class="submission-head">
      <span>运行ID</span>
      <span>题目</span>
      <span>运行结果</span>
      <span>运行时间(ms)</span>
      <span>使用内存(KB)</span>
      <span>代码长度</span>
      <span>使用语言</span>
      <span>提交时间</span>
    </div>

    <div class="submission-list">
      <SubmissionSimple
        v-for="submission in submissions"
        :author="author"
        :key="submission.submission_id || submission.id"
        :submission="submission"
      />
    </div>
  </div>
  <p v-else class="empty-text">还没有任何提交记录</p>

  <n-space vertical class="vertical-space">
    <n-pagination
      v-model:page="currentPage"
      :page-size="pageSize"
      :item-count="total"
      show-quick-jumper
      @update:page="fetchSubmissions"
    >
      <template #goto>
        跳转
      </template>
    </n-pagination>
  </n-space>
</template>

<script setup>
import { dialogError } from "@/utils/dialog"
import { useDialog  } from 'naive-ui'
import { ref, onMounted } from "vue";
import SubmissionSimple from "../coding/SubmissionSimple.vue";
import api from "@/api"; // 假设你有个 api 封装

const props = defineProps({
  author: {
    type: Object,
    required: true
  }
})

const dialog = useDialog()
const currentPage = ref(1);
const pageSize = ref(2);  // 每页条数
const total = ref(0);      // 总数，后端返回
const submissions = ref([]);

// 获取分页数据
const fetchSubmissions = async () => {
  try {
    const resp = await api.getUserSubmissionList(props.author.user_id, currentPage.value, pageSize.value)

    if (resp.code === 0) {
      submissions.value = resp.data.list; // 数据列表
      total.value = resp.data.total;      // 总数
    } else {
      dialogError(dialog, "提交记录请求发送成功", "服务器响应失败");
    }
  } catch {
    dialogError(dialog, "提交记录请求发送失败", "");
  }
};

onMounted(fetchSubmissions);
</script>

<style scoped>
.submission-table {
  overflow-x: auto;
  border: 1px solid #e6e9ee;
  border-radius: 6px;
  background: #ffffff;
}

.submission-head {
  display: grid;
  min-width: 780px;
  grid-template-columns: 72px minmax(112px, 1.25fr) 88px 92px 94px 72px 78px 150px;
  align-items: center;
  gap: 0;
  padding: 0 16px;
  background: #f3f4f6;
  color: #6f7782;
  font-size: 13px;
  font-weight: 500;
  line-height: 48px;
}

.empty-text {
  margin: 24px 0;
  color: #8a939d;
  text-align: center;
}

.vertical-space {
  align-items: center;
  margin-top: 18px;
}

@media (max-width: 900px) {
  .submission-table {
    border: 0;
    background: transparent;
  }

  .submission-head {
    display: none;
  }

  .submission-list {
    overflow-x: visible;
  }
}
</style>
