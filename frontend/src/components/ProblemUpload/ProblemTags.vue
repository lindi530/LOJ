<template>
  <div class="mt-3">
    <label class="form-label mb-3">标签（可多选）</label>
    
    <!-- 所有可选标签区域 - 系统提供的标签 -->
    <div class="d-flex flex-wrap gap-2 mb-4">
      <n-tag
        v-for="tag in allTags"
        :key="tag"
        :type="isSelected(tag) ? 'primary' : 'default'"
        :class="{ 
          'cursor-pointer': true,
          'transition-all duration-200': true 
        }"
        @click="toggleTag(tag)"
      >
        {{ tag }}
      </n-tag>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
import { NTag, NInput } from 'naive-ui';

const props = defineProps({
  modelValue: {  // 接收父组件的v-model绑定
    type: Array,
    default: () => []
  },
  allTags: {  // 系统提供的标签列表
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['update:modelValue']);

// 本地状态与父组件同步
const localTags = ref([...props.modelValue]);

// 检查标签是否被选中
const isSelected = (tag) => {
  return localTags.value.includes(tag);
};

// 切换标签选中状态
const toggleTag = (tag) => {
  const index = localTags.value.indexOf(tag);
  if (index > -1) {
    // 取消选择 - 从已选列表中移除
    localTags.value.splice(index, 1);
  } else {
    // 选中 - 添加到已选列表
    localTags.value.push(tag);
  }
  updateParent();
};

// 更新父组件数据
const updateParent = () => {
  nextTick(() => {
    emit('update:modelValue', [...localTags.value]);
  });
};

// 监听父组件数据变化，保持同步
watch(
  () => props.modelValue,
  (newVal) => {
    if (JSON.stringify(newVal) !== JSON.stringify(localTags.value)) {
      localTags.value = [...newVal];
    }
  },
  { deep: true }
);
</script>

<style scoped>
/* 可选标签未选中状态鼠标悬停效果 */
:deep(.n-tag--default.cursor-pointer:hover) {
  background-color: #f0f0f0;
}

/* 已选标签在可选区域的样式 */
:deep(.n-tag--primary.cursor-pointer) {
  opacity: 1; /* 保持完全不透明，通过颜色区分选中状态 */
}
</style>
