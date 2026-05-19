<template>
  <div class="bg-white p-4 rounded shadow-sm">
      <h5 class="mb-3">评论</h5>
        <div v-if="comments.length">
          
            <CommentItem
              v-for="item in comments.comments"
              :key="item.id"
              :comment="item"
              @update-comment="handleUpdateComment"
              @delete-comment="handleDeleteComment"
            />
        </div>
                
        <CommentCreate 
          :postId="postId"
          @create-comment="handleNewComment"
        />
      
    </div>
</template>

<script setup>
import api from '@/api'
import CommentItem from './CommentItem.vue'
import CommentCreate from './CommentCreate.vue'
import { ref, onMounted, computed } from 'vue'

const props = defineProps({
  postId: Number
})

const comments = ref([])

onMounted(async() => { 
  // 获取评论
  const resp = await api.getPostComments(props.postId)
  console.log("resp: ", resp)
  comments.value = resp.data
})


function handleUpdateComment(updated) {
  const comment = comments.value.comments.find(c => c.id === updated.id)
  if (comment) {
    comment.like = updated.like
    comment.likes = updated.likes
  }
}

function handleNewComment(newC) {
  // 假设 newC 是新评论对象
  comments.value.comments.push(newC)
  comments.value.length += 1
}

function handleDeleteComment(commentId) {
  comments.value.comments = comments.value.comments.filter(
    comment => comment.id !== commentId
  );
  comments.value.length -= 1;
}

</script>
