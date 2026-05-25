<template>
  <main class="pending-page">
    <div class="container py-5">
      <section class="card border-0 shadow-sm mx-auto text-center p-5 pending-card">
        <i class="bi pending-icon mb-3" :class="icon" aria-hidden="true"></i>
        <h1 class="h4 fw-semibold mb-2">{{ title }}</h1>
        <p class="text-secondary mb-4">{{ description }}</p>
        <RouterLink class="btn btn-outline-secondary mx-auto" :to="{ name: 'Competition' }">
          返回竞赛列表
        </RouterLink>
      </section>
    </div>
  </main>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    default: 'bi-hourglass-split'
  }
})
</script>

<style scoped>
.pending-page {
  min-height: calc(100vh - 60px);
}

.pending-card {
  max-width: 34rem;
  border-radius: 0.75rem;
}

.pending-icon {
  color: #ffa116;
  font-size: 2.75rem;
}
</style>
