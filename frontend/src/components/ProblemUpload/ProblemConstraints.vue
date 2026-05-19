<template>
  <n-card title="限制条件" class="mb-4">
    <!-- 表头行 -->
    <div class="table-header">
      <div class="col language-col">编程语言</div>
      <div class="col time-col">时间限制 (毫秒)</div>
      <div class="col memory-col">内存限制 (MB)</div>
      <div class="col submit-col">最大提交次数</div>
    </div>
    <!-- C++ 数据行 -->
    <div class="table-row">
      <div class="col language-col">C++</div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.cpp.timeLimit"
          :min="100"
          :step="100"
          @update:value="syncToParent"
        />
      </div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.cpp.memoryLimit"
          :min="1"
          :step="1"
          @update:value="syncToParent"
        />
      </div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.cpp.maxSubmissions"
          :min="1"
          :step="1"
          @update:value="syncToParent"
        />
      </div>
    </div>
    <!-- Python 数据行 -->
    <div class="table-row">
      <div class="col language-col">Python</div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.python.timeLimit"
          :min="100"
          :step="100"
          @update:value="syncToParent"
        />
      </div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.python.memoryLimit"
          :min="1"
          :step="1"
          @update:value="syncToParent"
        />
      </div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.python.maxSubmissions"
          :min="1"
          :step="1"
          @update:value="syncToParent"
        />
      </div>
    </div>
    <!-- Go 数据行 -->
    <div class="table-row">
      <div class="col language-col">Go</div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.go.timeLimit"
          :min="100"
          :step="100"
          @update:value="syncToParent"
        />
      </div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.go.memoryLimit"
          :min="1"
          :step="1"
          @update:value="syncToParent"
        />
      </div>
      <div class="col">
        <n-input-number
          v-model:value="localModel.go.maxSubmissions"
          :min="1"
          :step="1"
          @update:value="syncToParent"
        />
      </div>
    </div>
  </n-card>
</template>

<script setup>
import { reactive, watch } from 'vue';
import { NCard, NInputNumber } from 'naive-ui';

// 接收父组件传递的初始数据
const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      cpp: { timeLimit: 1000, memoryLimit: 256, maxSubmissions: 10 },
      python: { timeLimit: 1000, memoryLimit: 256, maxSubmissions: 10 },
      go: { timeLimit: 1000, memoryLimit: 256, maxSubmissions: 10 }
    })
  }
});

// 定义事件，用于向父组件传递更新后的数据
const emit = defineEmits(['update:modelValue']);

// 创建响应式本地数据，方便双向绑定
const localModel = reactive({ ...props.modelValue });

// 同步本地数据到父组件
const syncToParent = () => {
  emit('update:modelValue', { ...localModel });
};

// 监听父组件数据变化，同步到本地
watch(
  () => props.modelValue,
  (newVal) => {
    Object.assign(localModel, newVal);
  },
  { deep: true }
);
</script>

<style scoped>
/* 卡片样式 */
.n-card {
  --n-card-padding: 16px;
}

/* 表头和数据行通用的弹性布局设置 */
.table-header,
.table-row {
  display: flex;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
}

/* 表头样式 */
.table-header {
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  font-weight: bold;
}

/* 数据行样式 */
.table-row {
  border-bottom: 1px solid #e9ecef;
  padding: 8px 0;
}

/* 列通用样式，设置百分比宽度和居中对齐 */
.col {
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center;     /* 垂直居中 */
  padding: 0 8px;
  box-sizing: border-box;
}

/* 各列百分比宽度分配，可根据实际需求调整 */
.language-col {
  flex: 0 0 20%; /* 编程语言列占 20% 宽度 */
}
.time-col {
  flex: 0 0 25%; /* 时间限制列占 25% 宽度 */
}
.memory-col {
  flex: 0 0 25%; /* 内存限制列占 25% 宽度 */
}
.submit-col {
  flex: 0 0 30%; /* 最大提交次数列占 30% 宽度 */
}

/* 调整 Naive UI 输入框样式，使其在列中更好地居中并适应宽度 */
:deep(.n-input-number) {
  width: 80%; /* 让输入框占列宽的 80%，留出一些边距 */
  max-width: 120px; /* 限制最大宽度，避免在宽屏下列太宽时输入框过长 */
  margin: 0 auto; /* 水平居中 */
}
</style>