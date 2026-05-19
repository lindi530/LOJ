<template>
    <div v-if="posts.length">
      <SimplePost
        v-for="post in posts"
        :key="post.post_id"
        :post="post"
      />
    </div>

    <p v-else class="text-muted">还没有任何发帖记录</p>

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
 
</template>

<script setup>
import { dialogError } from "@/utils/dialog"
import { useDialog  } from 'naive-ui'
import { ref, onMounted, watch } from 'vue';
import api from '@/api';
import SimplePost from '../post/PostSimple.vue';

// props
const props = defineProps({
  userId: {
    type: String,
    required: true,
  },
  updatePosts: {
    type: Boolean,
    required: true,
  }
});
// emit
const emit = defineEmits(['post-deleted', 'update-state']);

// 响应式状态
const posts = ref([]);
const dialog = useDialog()
const currentPage = ref(1);
const pageSize = ref(2);  // 每页条数
const total = ref(0);      // 总数，后端返回

// 方法
// 获取帖子
const fetchPosts = async () => {
  try {
    const resp = await api.getPostsByUserId(props.userId, currentPage.value, pageSize.value);

    if (resp.code === 0) {
      emit('update-state', false)
      posts.value = resp.data.posts;        // 数据列表
      total.value = resp.data.total;      // 总数
    } else {
      dialogError(dialog, "提交记录请求发送成功", "服务器响应失败");
    }
  } catch (error) {
    dialogError(dialog, "提交记录请求发送失败", "");
  }
};

onMounted(fetchPosts)

watch(() => props.updatePosts, (newVal, oldVal) => {
  if (newVal) {
    fetchPosts(); // 你要调用的函数
  }
});
</script>


<style scoped>
.vertical-space {
  align-items: center;
  margin-top: auto;
}
</style>