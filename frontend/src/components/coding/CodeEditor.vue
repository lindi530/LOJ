<template>
  <ProblemCodeEditor
    aria-label="普通题库代码提交"
    :initial-input="initialInput"
    :languages="languages"
    :theme="theme"
    judge-mode="ACM"
    :can-run-example="canRunExample"
    :run-example-request="runExampleRequest"
    :submit-request="submitRequest"
    @update:theme="$emit('update:theme', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'
import api from '@/api'
import ProblemCodeEditor from './ProblemCodeEditor.vue'

const props = defineProps({
  problemId: Number,
  roomId: String,
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
  }
})

defineEmits(['update:theme'])

const canRunExample = computed(() => props.problemId !== null && props.problemId !== undefined && props.problemId !== '')

function runExampleRequest(payload) {
  return api.submitExample(props.problemId, payload)
}

function submitRequest(payload) {
  if (typeof props.roomId !== 'undefined') {
    return api.saberSubmit({
      room_id: props.roomId,
      language: payload.language,
      code: payload.code
    })
  }

  return api.submitCode(props.problemId, {
    language: payload.language,
    code: payload.code
  })
}
</script>
