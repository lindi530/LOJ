<template>
  <n-card title="题目基本信息" class="mb-4">
    <div class="row g-3">
      <div class="col-md-12">
        <label class="form-label">题目名称</label>
        <n-input
          v-model:value="localBasicInfo.title"
          placeholder="请输入题目名称"
          @input="handleTitleChange"
        />
      </div>
      
      <div class="col-md-6">
        <label class="form-label">难度级别</label>
        <n-select
          v-model:value="localBasicInfo.difficulty"
          placeholder="选择难度级别"
          style="width: 90%"
          :options="difficultyOptions"
          @update:value="handleDifficultyChange"
        />
      </div>

      <div class="col-md-6">
        <label class="form-label">题目分类</label>
        <n-select
          v-model:value="localBasicInfo.category"
          placeholder="选择题目分类"
          style="width: 90%"
          :options="categoryOptions"
          @update:value="handleFieldChange('category', $event)"
        />
      </div>
      
      <!-- 标签组件调用：修复props传递，使用all-tags与子组件匹配 -->
      <div class="col-md-12">
        <ProblemTags 
          v-model="localBasicInfo.tags"  
          :all-tags="tagSuggestions"   
        />
      </div>
    </div>
  </n-card>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
import { NCard, NInput, NSelect } from 'naive-ui';
import ProblemTags from './ProblemTags.vue';

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      title: '',
      difficulty: '',
      category: '',
      tags: []  // 确保初始包含tags属性
    })
  }
});

const emit = defineEmits(['update:modelValue']);

// 本地状态初始化
const localBasicInfo = ref({ ...props.modelValue });

// 分类选项
const categoryOptions = [
  { label: '算法', value: 'algorithm' },
  { label: '数据结构', value: 'data_structure' },
  { label: '数学', value: 'math' },
  { label: '字符串处理', value: 'string' },
  { label: '其他', value: 'other' }
];

// 系统标签列表（传递给ProblemTags组件）
const tagSuggestions = [
  '数组', '字符串', '链表', '树', '图', '动态规划', 
  '贪心算法', '回溯', '排序', '查找', '哈希表',
  '双指针', '递归', '分治', '位运算', '数学',
  'BFS', 'DFS', '堆', '栈', '队列', '滑动窗口'
];

// 难度选项
const difficultyOptions = [
  { label: '简单', value: '简单' },
  { label: '中等', value: '中等' },
  { label: '困难', value: '困难' }
];

// 通用字段变更处理
const handleFieldChange = (field, value) => {
  localBasicInfo.value[field] = value
  updateParent()
};

// 处理各个字段变化
const handleTitleChange = (value) => {
  localBasicInfo.value.title = value;
  updateParent();
};

const handleDifficultyChange = (value) => {
  localBasicInfo.value.difficulty = value;
  updateParent();
};

// 更新父组件
const updateParent = () => {
    nextTick(() => {
    emit('update:modelValue', { ...localBasicInfo.value });
  });
};


watch(
  () => localBasicInfo.value.tags,
  (newTags) => {
    updateParent();
  },
  { deep: true } // 深度监听数组变化
);

watch(
  () => props.modelValue,
  (newVal) => {
    if (JSON.stringify(newVal) !== JSON.stringify(localBasicInfo.value)) {
      localBasicInfo.value = { ...newVal };
    }
  },
  { deep: true }
);
</script>

<style scoped>
:deep(.n-select-menu) {
  z-index: 1000 !important;
}
</style>
