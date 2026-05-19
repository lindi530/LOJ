<template>
  <div class="d-flex justify-content-end gap-3 mt-4">
    <n-button 
      type="info" 
      @click="handleSaveDraft"
    >
      保存草稿
    </n-button>
    <n-button 
      type="primary" 
      @click="handleSubmitForReview"
    >
      提交审核
    </n-button>
  </div>
</template>

<script setup>
import { NButton } from 'naive-ui';

// 定义事件
const emit = defineEmits(['saveDraft', 'submitForReview']);

// 处理保存草稿
const handleSaveDraft = () => {
  emit('saveDraft');
};

// 处理提交审核
const handleSubmitForReview = () => {
  emit('submitForReview');
};
</script>
