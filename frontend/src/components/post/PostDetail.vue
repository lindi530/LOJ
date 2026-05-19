<template>
  <div class="container mt-5">
    <n-card>
    <h2 class="text-center mb-3">{{ post.title }}</h2>

    <div class="d-flex justify-content-between align-items-center text-muted mb-4">
      <div class="d-flex align-items-center">
        <span class="me-1">作者：</span>
        <router-link
          :to="`/users/${post.author?.user_id}`"
          class="d-flex align-items-center text-decoration-none text-dark"
        >
          <img
            :src="post.author?.avatar"
            alt="avatar"
            class="rounded-circle me-2"
            width="32"
            height="32"
          />
          <strong>{{ post.author?.user_name }}</strong>
        </router-link>
      </div>
        
      <div>阅读量：{{ post.views || 0 }}</div>
    </div>

    <div class="bg-light p-4 rounded shadow-sm" style="margin-bottom: 12px;">
      <p class="mb-0" style="white-space: pre-wrap;">{{ post.content }}</p>
    </div>

    <div class="d-flex align-items-center">
      <!-- 美化后的点赞按钮 -->
      <button
        class="like-button"
        :class="{ 'liked': post.value?.liked }"
        @click="likePost"
      >
        <i class="bi bi-hand-thumbs-up"></i>
        <span class="like-text">{{ LikeMassage }}</span>
        <span class="like-count">{{ post.likes || 0 }}</span>
      </button>
    </div>
    <CommentList 
    :postId="postId" />
    </n-card>
  </div>
</template>

<script setup>
import { dialogError, dialogInfo } from '@/utils/dialog'
import { useDialog } from "naive-ui"
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import CommentList from '../comment/CommentList.vue'

const dialog = useDialog()
const route = useRoute()
const postId = route.params.post_id

const likeMsg = '取消点赞'
const UnLikeMsg = '点赞'

const post = ref({})

onMounted(async () => {
  // 获取帖子
  const resp = await api.getPostByPostId(postId)
  if (resp.code == 0) {
    post.value = resp.data
  } else {
    dialogError(dialog, "获取帖子信息失败", resp.message)
  }
  
})

const LikeMassage = computed(() => {
  return post.value.liked ? likeMsg : UnLikeMsg
})

const likePost = async () => {
  if (post.value.liked === false) {
    const resp = await api.likePost(postId)
    if (resp.code === 0) {
      post.value.likes = (post.value.likes + 1 || 0)
      post.value.liked = !post.value.liked
      dialogInfo(dialog, "点赞成功")
    } else {
      dialogError(dialog, "点赞失败", resp.message)
    }
  } else {
    const resp = await api.unLikePost(postId)
    if (resp.code === 0) {
      post.value.likes = (post.value.likes - 1 || 0)
      post.value.liked = !post.value.liked
      dialogInfo(dialog, "取消点赞成功")
    } else {
      dialogError(dialog, "取消点赞失败", resp.message)
    }
  }
}
</script>

<style scoped>
.container {
  max-width: 800px;
}

/* 点赞按钮样式 */
.like-button {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border: 1px solid #0d6efd;
  border-radius: 20px;
  background-color: transparent;
  color: #0d6efd;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.like-button:hover {
  background-color: rgba(13, 110, 253, 0.1);
  transform: translateY(-2px);
}

.like-button:active {
  transform: translateY(0);
}

/* 已点赞状态 */
.like-button.liked {
  background-color: #0d6efd;
  color: white;
  border-color: #0d6efd;
}

.like-button.liked:hover {
  background-color: #0b5ed7;
}

/* 图标样式 */
.like-button i {
  font-size: 16px;
  transition: transform 0.3s ease;
}

.like-button:hover i {
  transform: scale(1.1);
}

/* 点赞数量样式 */
.like-count {
  padding: 2px 8px;
  border-radius: 10px;
  background-color: rgba(13, 110, 253, 0.1);
  transition: all 0.3s ease;
}

.like-button.liked .like-count {
  background-color: rgba(255, 255, 255, 0.2);
}
</style>
