<template>
  <div class="container my-4">
    <n-card>
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
        <h4 class="mb-0">All Posts</h4>
        <h5 class="mb-0">第 {{ currentPage }} 页</h5>
      </div>

      <div v-if="posts && posts.length" >
        <div v-for="post in posts" :key="post.post_id">
          <SimplePost :post="post"/>
        </div>
      </div>
      <p v-else class="text-muted">还没有用户发送帖子</p>

      <n-space vertical class="vertical-space" >
        <n-pagination
          v-model:page="currentPage"
          :page-size="pageSize"
          :item-count="total"
          show-quick-jumper
          @update:page="fetchPosts"
        >
          <template #goto>
            跳转
          </template>
        </n-pagination>
      </n-space>
    </n-card>
  </div>
</template>

<script setup>
import { dialogError } from "@/utils/dialog"
import { useDialog  } from 'naive-ui'
import { ref, onMounted } from 'vue';
import SimplePost from '@/components/post/PostSimple.vue';
import api from '@/api';

// 响应式数据
const posts = ref([]);
const dialog = useDialog()
const currentPage = ref(1);
const pageSize = ref(10);  // 每页条数
const total = ref(0);      // 总数，后端返回

// 获取帖子
const fetchPosts = async () => {
  try {
    const resp = await api.getPosts(currentPage.value, pageSize.value);

    if (resp.code === 0) {
      posts.value = resp.data.list;        // 数据列表
      total.value = resp.data.total;      // 总数
    } else {
      dialogError(dialog, "提交记录请求发送成功", "服务器响应失败");
    }
  } catch (error) {
    dialogError(dialog, "提交记录请求发送失败", "");
  }
  window.scrollTo(0, 0);
};

// 格式化日期
const formattedDate = (dateStr) => {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric', month: 'short', day: 'numeric'
  });
};

// 组件挂载后执行
onMounted(() => {
  fetchPosts();
});
</script>

<style scoped>
.vertical-space {
  align-items: center;
  margin-top: auto;
}
</style>
