<template>
  <n-card title="题目描述" class="mb-4">
    <div class="mt-3">
      <label class="form-label">题目描述</label>
      <n-input
        v-model:value="localDescription.description"
        type="textarea"
        rows="6"
        placeholder="请输入题目描述"
        @input="handleDescriptionChange"
      />
    </div>
    
    <div class="mt-3">
      <label class="form-label">输入描述</label>
      <n-input
        v-model:value="localDescription.inputFormat"
        type="textarea"
        rows="3"
        placeholder="请输入输入描述"
        @input="handleInputFormatChange"
      />
    </div>
    
    <div class="mt-3">
      <label class="form-label">输出描述</label>
      <n-input
        v-model:value="localDescription.outputFormat"
        type="textarea"
        rows="3"
        placeholder="请输入输出描述"
        @input="handleOutputFormatChange"
      />
    </div>
    
    <div class="mt-3">
      <label class="form-label">提示</label>
      <n-input
        v-model:value="localDescription.hint"
        type="textarea"
        rows="3"
        placeholder="请输入提示信息（可选）"
        @input="handleHintChange"
      />
    </div>
  </n-card>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
import { NCard, NInput } from 'naive-ui';

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      description: '',
      inputFormat: '',
      outputFormat: '',
      hint: ''
    })
  }
});

const emit = defineEmits(['update:modelValue']);

// 本地数据存储
const localDescription = ref({ ...props.modelValue });

// 显式处理各个输入框的变化
const handleDescriptionChange = (value) => {
  localDescription.value.description = value;
  updateParent();
};

const handleInputFormatChange = (value) => {
  localDescription.value.inputFormat = value;
  updateParent();
};

const handleOutputFormatChange = (value) => {
  localDescription.value.outputFormat = value;
  updateParent();
};

const handleHintChange = (value) => {
  localDescription.value.hint = value;
  updateParent();
};

// 统一更新父组件的方法
const updateParent = () => {
  nextTick(() => {
    emit('update:modelValue', { ...localDescription.value });
  });
};

// 监听父组件数据变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (JSON.stringify(newVal) !== JSON.stringify(localDescription.value)) {
      localDescription.value = { ...newVal };
    }
  },
  { deep: true }
);
</script>
