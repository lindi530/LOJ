<template>
  <div v-if="submissions.length">
    <SubmissionSimple
      v-for="submission in submissions"
      :author="author"
      :key="submission.submission_id"
      :submission="submission"
    />
  </div>
  <p v-else class="text-muted">还没有任何提交记录</p>

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
.vertical-space {
  align-items: center;
  margin-top: auto;
}
</style>
