<template>
  <n-card class="container mt-4 mb-5">
    <h2 class="mb-4">题目上传</h2>
    
    <ProblemBasicInfo 
      v-model="problemData.basicInfo" 
    />
    
    <ProblemEditor 
      v-model="problemData.description" 
    />
    
    <TestCasesUpload 
      v-model="problemData.testCases" 
      ref="testCasesRef"
    />
    
    <ProblemConstraints 
      v-model="problemData.constraints" 
    />
    
    <SubmitActions 
      @saveDraft="handleSaveDraft"
      @submitForReview="handleSubmitForReview"
    />
  </n-card>
</template>

<script setup>
import ProblemBasicInfo from '@/components/ProblemUpload/ProblemBasicInfo.vue';
import ProblemEditor from '@/components/ProblemUpload/ProblemEditor.vue';
import TestCasesUpload from '@/components/ProblemUpload/TestCasesUpload.vue';
import ProblemConstraints from '@/components/ProblemUpload/ProblemConstraints.vue';
import SubmitActions from '@/components/ProblemUpload/SubmitActions.vue';
import { ref, onMounted } from 'vue'
import { useDialog } from 'naive-ui'
import api from '@/api'

// 主数据存储
const problemData = ref({
  basicInfo: {
    title: '',
    difficulty: '',
    tags: []
  },
  description: {
    description: '',
    inputFormat: '',
    outputFormat: '',
    hint: ''
  },
  testCases: [],
  constraints: {
    cpp: { timeLimit: 1000, memoryLimit: 256, maxSubmissions: 10 },
    python: { timeLimit: 1000, memoryLimit: 256, maxSubmissions: 10 },
    go: { timeLimit: 1000, memoryLimit: 256, maxSubmissions: 10 }
  }
});

const testCasesRef = ref(null); 
const dialog = useDialog()

const validateProblemData = () => { 
  const d = problemData.value

  // 基础信息
  if (!d.basicInfo.title) return "题目标题不能为空"
  if (!d.basicInfo.difficulty) return "题目难度不能为空"
  if (d.basicInfo.tags.length === 0) return "至少选择一个标签"

  // 描述
  if (!d.description.description) return "题目描述不能为空"
  if (!d.description.inputFormat) return "输入格式不能为空"
  if (!d.description.outputFormat) return "输出格式不能为空"

  // 测试用例
  if (d.testCases.length === 0) return "至少要有一个测试用例"

  // 约束 (可以根据需求细化)
  for (const [lang, c] of Object.entries(d.constraints)) {
    if (!c.timeLimit || !c.memoryLimit) {
      return `${lang} 的时间/内存限制不能为空`
    }
  }

  return null // 通过
}

// 移除可能导致递归的深层watch，只在需要时手动处理

const handleSaveDraft = () => {
  const draft = {
    data: problemData.value,
    savedAt: Date.now()
  }
    localStorage.setItem('problem-draft', JSON.stringify(draft))
  console.log("Save:", problemData.value)
  alert('草稿已保存到本地')
}

// 页面加载时恢复
onMounted(() => {
  const draft = localStorage.getItem('problem-draft')
  if (draft) {
    const { data, savedAt } = JSON.parse(draft)
    if (Date.now() - savedAt < 60 * 60 * 1000) { // 1小时
      problemData.value = data
    } else {
      localStorage.removeItem('problem-draft') // 超时清理
    }
    }
  console.log("problemData", problemData.value)
})


const handleSubmitForReview = async () => {
  await testCasesRef.value?.matchTestCases(); 

  const errorMsg = validateProblemData()
  if (errorMsg) {
    dialog.error({
      title: '校验失败',
      content: errorMsg,
      positiveText: '确定'
    })
    return
  }

  

  dialog.success({
    title: '校验成功',
    content: '所有必填项已填写，可以提交',
    positiveText: '上传',
    onPositiveClick: async () => {
      try {
        const resp = await api.uploadProblem(problemData.value)
        if (resp.code == 0) {
          msg = resp.message
        }
      } catch { }
    }
  })

  localStorage.removeItem('problem-draft')
};
</script>