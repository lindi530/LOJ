<template>
  <ProblemCodeEditor
    aria-label="竞赛代码提交"
    :initial-input="initialInput"
    :languages="languages"
    :theme="theme"
    :judge-mode="judgeMode"
    :can-run-example="canRunExample"
    :run-example-request="runExampleRequest"
    :submit-request="submitRequest"
    @update:theme="$emit('update:theme', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'
import api from '@/api'
import ProblemCodeEditor from '@/components/coding/ProblemCodeEditor.vue'

const props = defineProps({
  competitionId: {
    type: [Number, String],
    required: true
  },
  problemNumber: {
    type: String,
    required: true
  },
  problemId: {
    type: [Number, String],
    default: null
  },
  initialInput: {
    type: String,
    default: ''
  },
  languages: {
    type: [Array, Object, String],
    default: () => []
  },
  theme: {
    type: String,
    default: 'vs-light'
  },
  judgeMode: {
    type: String,
    default: 'ACM'
  }
})

defineEmits(['update:theme'])

const canRunExample = computed(() => props.problemId !== null && props.problemId !== undefined && props.problemId !== '')

function runExampleRequest(payload) {
  return api.submitExample(props.problemId, payload)
}

function submitRequest(payload) {
  return api.submitCompetitionProblem(props.competitionId, props.problemNumber, payload)
}
</script>
