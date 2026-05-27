<template>
  <CompetitionDetailPage :competition="competition" />
</template>

<script setup>
import CompetitionDetailPage from '@/components/CompetitionDetail/CompetitionDetailPage.vue'

defineProps({
  competition: {
    type: Object,
    default: null
  }
})
</script>
