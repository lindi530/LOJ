<template>
  <n-card>
    <h5 class="mb-3">创建新帖子</h5>
    <form @submit.prevent="submit">
      <div class="mb-3">
        <label for="title" class="form-label">Title</label>
        <input v-model="post.title" type="text" id="title" class="form-control" required>
      </div>
      <div class="mb-3">
        <label for="content" class="form-label">Content</label>
        <textarea v-model="post.content" id="content" class="form-control" rows="5" required></textarea>
      </div>
      <button type="submit" class="btn btn-primary">Post</button>
    </form>
  </n-card>
</template>

<script setup>
import { ref } from 'vue';
import { dialogError } from '@/utils/dialog.js'
import { useDialog }  from 'naive-ui'

// 响应式数据
const post = ref({
  title: '',
  content: ''
});

// emit
const emit = defineEmits(['update-state']);

const dialog = useDialog()
// 提交方法
const submit = () => {
  const newPost = {
    ...post.value,
    status: 0,
    content: post.value.content.slice(0, 50) // 截取前50字符
  };

  addPost(newPost)

  // 重置表单
  post.value.title = '';
  post.value.content = '';
};

async function addPost(post) {
  try {
    const resp = await api.createPost({ user_id: userId.value, ...post });
    if (resp.code === 0) {
      emit('update-state', false)
      dialogInfo(dialog, resp.message)
    } else {
      dialogInfo(dialog, resp.message)
    }
  } catch {
    dialogError(dialog, "发送请求失败", "")
  }
}
</script>

