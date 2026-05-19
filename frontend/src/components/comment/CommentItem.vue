<template>
<!-- 左侧头像 -->
 <div class="d-flex border-bottom py-3">
    <router-link
    :to="`/users/${comment.author.user_id}`"
    class="text-decoration-none"
    >
    <img
        :src="comment.author.avatar"
        alt="avatar"
        class="rounded-circle me-3"
        width="50"
        height="50"
    />
    </router-link>

    <!-- 右侧评论内容 -->
    <div class="flex-grow-1">
    <!-- 第一行：用户名 + 时间 + 点赞 + 删除 -->
    <div class="d-flex justify-content-between align-items-center flex-wrap gap-2">
        <!-- 左侧：用户名、时间、点赞 -->
        <div class="d-flex align-items-center flex-wrap gap-2">
        <router-link
            :to="`/users/${comment.author.user_id}`"
            class="text-dark text-decoration-none"
        >
            <strong>{{ comment.author.user_name }}</strong>
        </router-link>

        <small class="text-muted">{{ formatDate(comment.created_at) }}</small>


        <n-button quaternary circle type="info" @click="likeComment">
            <template #icon>
              <img
              :src="iconSrc"
              alt="点赞"
              width="16"
              height="16"
              />
            </template>
            
        </n-button>

        <span class="text-secondary" style="font-size: 0.875rem;">
          {{ comment.likes }}
        </span>
        
        <!-- 点赞按钮 -->
    
        </div>

        <!-- 右侧：删除按钮 -->
        <button
        class="btn btn-sm btn-outline-danger px-1 py-0"
        style="font-size: 0.75rem;"
        @click="deleteComment(comment.id)"
        >
        <i class="bi bi-trash"></i>
        </button>

    </div>

    <!-- 第二行：评论内容 -->
    <p class="mb-0 mt-1">{{ comment.content }}</p>
    </div>
  </div>
</template>

<script setup>
import api from '@/api'
import { computed } from 'vue'
import { formatDate } from '@/utils/date'

const props  = defineProps({
  comment: Object,
})

console.log("comment: ", props.comment)

const emit = defineEmits(['update-comment', 'delete-comment'])

const iconUnlike = '/voteup-empty.png'
const iconLike = '/voteup.png'

const iconSrc = computed(() => {
  return props.comment.like ? iconLike : iconUnlike
})

async function likeComment() {
  if (props.comment.like === false) {
    const res = await api.likeComment(props.comment.id)
    if (res.code == 0) {
        emit('update-comment', {
        id: props.comment.id,
        like: true,
        likes: props.comment.likes + 1
      })
    }
  } else {
    const res = await api.unlikeComment(props.comment.id)
      if (res.code == 0) {
        emit('update-comment', {
        id: props.comment.id,
        like: false,
        likes: props.comment.likes - 1
      })
    }   
  } 

}

async function deleteComment(commentId) {
  const res = await api.deleteComment(commentId)
  console.log("res:", res)
  if (res.code == 0) {
    emit('delete-comment', commentId)
  }
}
</script>

<style scoped>
</style>