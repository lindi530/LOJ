<template>
    <div class="mt-3">
        <textarea v-model="newComment" class="form-control mb-2" rows="3" placeholder="写下你的评论..."></textarea>
        <button @click="submitComment" class="btn btn-primary">发表评论</button>
      </div>
</template>

<script setup>
import api from '@/api'
import { ref } from 'vue'
const props = defineProps({
    postId: Number
})
const emit = defineEmits(['create-comment'])

const newComment = ref('')

async function submitComment() {
  if (!newComment.value.trim()) return

  const res = await api.createPostComment(props.postId, {
    content: newComment.value
  })

  const newC = res.data

  emit('create-comment', newC)

  newComment.value = ''
}
</script>