<template>
  <div class="card mb-3">
    <!-- 卡片头部 -->
    <div class="card-header d-flex align-items-center justify-content-between">
      <!-- 左侧头像 -->
      <router-link
      :to="`/users/${post.author?.user_id}`"
      class="d-flex align-items-center text-decoration-none text-dark"
        >
        <div class="me-3">
          <img :src="post.author.avatar" alt="头像" class="rounded-circle" width="40" height="40">
        </div>
      </router-link>

      <!-- 中间用户名 + 时间 -->
      <div class="flex-grow-1">
         <div class="d-flex align-items-center">
          <router-link
            :to="`/users/${post.author?.user_id}`"
            class="text-decoration-none text-dark"
          >
            <strong>{{ post.author?.user_name }}</strong>
          </router-link>
        </div>
        <div class="text-muted small">{{ formatDate(post.created_at ) }}</div>
      </div>

      <!-- 右侧标题 -->
       <router-link
        :to="`/posts/${post.post_id}`"
        class="flex-grow-1 text-decoration-none text-body"
        >
        <div class="ms-3 text-end fw-bold">
{{ post.title }}
        </div>
      </router-link>
      
    </div>

    <!-- 卡片内容 -->
    <div class="card-body p-2">
      <pre class="bg-light p-3 rounded overflow-auto" style="max-height: 90px">
<router-link
  :to="`/posts/${post.post_id}`"
  class="flex-grow-1 text-decoration-none text-body"
>{{ post.content }}
</router-link>
      </pre>
    </div>
  </div>
</template>


<script setup>
import { formatDate } from '@/utils/date'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: monospace;
}
</style>
