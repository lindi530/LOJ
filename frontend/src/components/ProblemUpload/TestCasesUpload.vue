<template>
  <n-card title="测试用例" class="mb-4">
    <!-- 左右卡片容器 -->
    <div style="display: flex; gap: 16px; width: 100%;">
      <!-- 左侧输入卡片 -->
      <n-card 
        style="flex: 1; min-width: 0;"
        bordered
        title="输入文件 (Input)"
      >
        <div style="display: flex; gap: 12px;">
          <p style="width: 100%" class="text-muted mb-3">上传测试用例的输入文件（请确保文件名正确，如input1.txt）</p>
          <n-button 
            v-if="inputFiles.length > 0" 
            @click="clearInputFiles"
            type="error" dashed
          >
            清空输入文件
          </n-button>
        </div>

        <n-upload
            v-model:file-list="inputFiles"
            @change="handleInputFileChange"
            @remove="handleInputFileRemove"
            :max-size="10 * 1024 * 1024"
            :multiple="true"
            type="text"
            style="flex: 1;"
        >
            <n-button type="primary" dashed style="width: 100%;">选择输入文件</n-button>
        </n-upload>

        <!-- 显示已选输入文件数量 -->
        <div v-if="inputFiles.length > 0" class="mt-2 text-sm">
            已选择 {{ inputFiles.length }} 个文件
        </div>
      </n-card>
      
      <!-- 右侧输出卡片 -->
      <n-card 
        style="flex: 1; min-width: 0;"
        bordered
        title="输出文件 (Output)"
      >

        <div style="display: flex; gap: 12px;">
          <p style="width: 100%" class="text-muted mb-3">上传测试用例的输出文件（请确保文件名正确，如output1.txt）</p>
          <n-button 
            v-if="outputFiles.length > 0" 
            @click="clearOutputFiles"
            type="error" dashed
          >
            清空输出文件
          </n-button>
        </div>

        <n-upload
            v-model:file-list="outputFiles"
            @change="handleOutputFileChange"
            @remove="handleOutputFileRemove"
            :max-size="10 * 1024 * 1024"
            :multiple="true"
            type="text"
            style="flex: 1;"
        >
            <n-button type="primary" dashed style="width: 100%;">选择输出文件</n-button>
        </n-upload>
        
        <!-- 显示已选输出文件数量 -->
        <div v-if="outputFiles.length > 0" class="mt-2 text-sm">
          已选择 {{ outputFiles.length }} 个文件
        </div>
      </n-card>
    </div>
    
    <div class="mt-3 flex items-center flex-wrap gap-2">
      <div style="display: flex; align-items: center;">
        <!-- <n-button 
            type="info" 
            @click="matchTestCases"
            :disabled="inputFiles.length === 0 || outputFiles.length === 0"
            class="mr-3"
        >
            匹配并添加测试用例
        </n-button> -->
        
        <!-- 匹配结果提示信息 -->
        <div v-if="matchResult.message" 
            :class="matchResult.type === 'success' ? 'text-success ml-2' : 'text-error ml-2'"
        >
            {{ matchResult.message }}
        </div>
      </div>
      
    </div>
  </n-card>
</template>

<script setup>
import { ref, watch, nextTick, onMounted  } from 'vue';
import { NUpload, NButton, NCard } from 'naive-ui';

// 定义属性和事件
const props = defineProps({
  modelValue: {
    type: Array,
  }
});

const emit = defineEmits(['update:modelValue']);

// 本地数据
const inputFiles = ref([])
const outputFiles = ref([])
const localTestCases = ref([])
const matchResult = ref({
  message: '',
  type: 'success'
});

// 处理文件上传
const handleInputFileChange = (fileList) => {
  const files = normalizeFileList(fileList);
  const newFiles = files.filter(file => 
    !inputFiles.value.some(existing => existing.name === file.name)
  );
  inputFiles.value = [...inputFiles.value, ...newFiles];
  matchResult.value = { message: '', type: 'success' };
};

const handleInputFileRemove = (fileData) => { 
    inputFiles.value.splice(fileData.index, 1);
}

const handleOutputFileRemove = (fileData) => { 
    outputFiles.value.splice(fileData.index, 1)
}

const handleOutputFileChange = (fileList) => {
  const files = normalizeFileList(fileList);
  const newFiles = files.filter(file => 
    !outputFiles.value.some(existing => existing.name === file.name)
  );
  outputFiles.value = [...outputFiles.value, ...newFiles];
  matchResult.value = { message: '', type: 'success' };
};

// 清空输入文件
const clearInputFiles = () => {
  inputFiles.value = [];
  localTestCases.value = [];
  matchResult.value = { message: '', type: 'success' };
  updateParent()
};

// 清空输出文件
const clearOutputFiles = () => {
  localTestCases.value = [];
  outputFiles.value = [];
  matchResult.value = { message: '', type: 'success' };
  updateParent()
};

// 标准化文件列表格式
const normalizeFileList = (fileList) => {
  if (fileList && fileList.fileList) {
    return fileList.fileList;
  }
  if (fileList && !Array.isArray(fileList)) {
    return [fileList];
  }
  return Array.isArray(fileList) ? fileList : [];
};

// 匹配测试用例
const matchTestCases = async () => {
  localTestCases.value = []
  matchResult.value = { message: '', type: 'success' };
  
  let matchedCount = 0;
  const invalidFiles = [];
  const unmatchedInputs = [];
  const unmatchedOutputs = [];
  const tempMatches = [];
  const matchedOutputIndices = [];
  let pendingReads = 0; // 跟踪待处理的文件读取操作
  
  // 遍历所有输入文件
  inputFiles.value.forEach(inputFile => {
    const inputFileObj = getFileObject(inputFile);
    if (!inputFileObj) {
      invalidFiles.push(inputFile.name || '未知文件');
      unmatchedInputs.push(inputFile.name || '未知文件');
      return;
    }
    
    const inputBase = getBaseName(inputFile.name);
    const inputNumber = extractNumbers(inputBase);
    let isMatched = false;
    
    // 在输出文件中查找匹配的文件
    outputFiles.value.forEach((outputFile, outputIndex) => {
      if (matchedOutputIndices.includes(outputIndex)) return;
      
      const outputFileObj = getFileObject(outputFile);
      if (!outputFileObj) {
        invalidFiles.push(outputFile.name || '未知文件');
        return;
      }
      
      const outputBase = getBaseName(outputFile.name);
      const outputNumber = extractNumbers(outputBase);
      
      if (inputNumber && inputNumber === outputNumber) {
        tempMatches.push({
          inputFile: inputFile.name,
          outputFile: outputFile.name,
          inputObj: inputFileObj,
          outputObj: outputFileObj
        });
        
        matchedCount++;
        isMatched = true;
        matchedOutputIndices.push(outputIndex);
      }
    });
    
    if (!isMatched) {
      unmatchedInputs.push(inputFile.name || '未知文件');
    }
  });
  
  // 收集未匹配的输出文件
  outputFiles.value.forEach((file, index) => {
    if (!matchedOutputIndices.includes(index)) {
      unmatchedOutputs.push(file.name || '未知文件');
    }
  });
  
  // 读取所有匹配文件的内容
  pendingReads = tempMatches.length;
  if (pendingReads === 0) {
    // 如果没有需要读取的文件，直接更新父组件
    updateParent();
  } else {
    tempMatches.forEach(match => {
      readFileContents(match, () => {
        pendingReads--;
        // 所有文件读取完成后再更新父组件
        if (pendingReads === 0) {
          updateParent();
        }
      });
    });
  }
  
  // 显示结果统计
  nextTick(() => {
    let message = `成功匹配 ${matchedCount} 组测试用例`;
    let isError = false;
    
    if (invalidFiles.length > 0) {
      message += `，无效文件: ${invalidFiles.join(', ')}`;
      isError = true;
    }
    if (unmatchedInputs.length > 0) {
      message += `，未匹配的输入文件: ${unmatchedInputs.join(', ')}`;
      isError = true;
    }
    if (unmatchedOutputs.length > 0) {
      message += `，未匹配的输出文件: ${unmatchedOutputs.join(', ')}`;
      isError = true;
    }
    
    matchResult.value = {
      message: message,
      type: isError ? 'error' : 'success'
    };
  });
};

// 读取文件内容的辅助函数，添加回调参数
const readFileContents = (match, onComplete) => {
  const inputReader = new FileReader();
  inputReader.onload = (e) => {
    const inputContent = e.target.result;
    
    const outputReader = new FileReader();
    outputReader.onload = (e2) => {
      const outputContent = e2.target.result;
      
      localTestCases.value.push({
        inputFile: match.inputFile,
        outputFile: match.outputFile,
        inputContent: inputContent,
        outputContent: outputContent,
        isPublic: false
      });
      
      // 调用回调通知完成
      onComplete();
    };
    
    outputReader.readAsText(match.outputObj);
  };
  
  inputReader.readAsText(match.inputObj);
};

// 安全获取文件对象的辅助函数
const getFileObject = (file) => {
  const checkProperties = (obj) => {
    if (!obj || typeof obj !== 'object') return null;
    if (obj instanceof Blob) return obj;
    
    const fileProperties = ['raw', 'file', 'originFileObj'];
    for (const prop of fileProperties) {
      if (prop in obj) {
        const result = checkProperties(obj[prop]);
        if (result) return result;
      }
    }
    
    return null;
  };
  
  return checkProperties(file);
};

// 获取文件名（不含扩展名）
const getBaseName = (fileName) => {
  if (!fileName) return '';
  const dotIndex = fileName.lastIndexOf('.');
  return dotIndex > -1 ? fileName.substring(0, dotIndex) : fileName;
};

// 提取文件名中的数字部分
const extractNumbers = (name) => {
  if (!name) return '';
  const match = name.match(/\d+/g);
  return match ? match.join('') : '';
};

// 更新父组件
const updateParent = () => {
  nextTick(() => {
    // 创建新数组引用以确保响应式更新
    const newTestCases = [...localTestCases.value];
    localTestCases.value = newTestCases;
    emit('update:modelValue', newTestCases);
  });
};

// 监听父组件数据变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (JSON.stringify(newVal) !== JSON.stringify(localTestCases.value)) {
      localTestCases.value = [...newVal];
    }
  },
  { deep: true }
);
watch(
    () => props.modelValue,
    (newVal, oldVal) => {
        const saved = newVal

        if (saved && Array.isArray(saved)) {
          localTestCases.value = [...saved]

            // 构建 n-upload 可识别的 file-list
          inputFiles.value = saved.map((tc, index) => ({
            id: `input${index}`,
            name: tc.inputFile,
            status: 'finished',
            file: new File([tc.inputContent], tc.inputFile, { type: 'text/plain' })
          }))

          outputFiles.value = saved.map((tc, index) => ({
            id: `output${index}`,
            name: tc.outputFile,
            status: 'finished',
            file: new File([tc.outputContent], tc.outputFile, { type: 'text/plain' })
          }))
        }
    },
    { deep: true } // 如果是对象或数组，使用 deep
)

defineExpose({
  matchTestCases
});
</script>
