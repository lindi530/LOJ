<template>
  <section class="position-relative text-center contest-banner mb-4 py-4 py-lg-5">
    <span class="trophy-mark mb-3">
      <i class="bi bi-trophy-fill" aria-hidden="true"></i>
    </span>
    <h1 class="h3 fw-semibold mb-2">LOJ 竞赛</h1>
    <p class="text-secondary mb-4 mb-md-0">挑战算法竞赛，与优秀选手一同成长</p>

    <button type="button" class="btn btn-warning px-4 py-2 fw-medium create-button" @click="$emit('create')">
      <i class="bi bi-plus-lg me-2" aria-hidden="true"></i>
      创建竞赛
    </button>
  </section>
</template>

<script setup>
defineEmits(['create'])
</script>

<style scoped>
.contest-banner {
  min-height: 11rem;
}

.trophy-mark {
  display: inline-flex;
  width: 4rem;
  height: 4rem;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #fff3de;
  color: #ffa116;
  font-size: 2rem;
}

.create-button {
  position: absolute;
  top: 1.5rem;
  right: 0;
}

.btn-warning {
  --bs-btn-bg: #ffa116;
  --bs-btn-border-color: #ffa116;
  --bs-btn-hover-bg: #ed940f;
  --bs-btn-hover-border-color: #ed940f;
  --bs-btn-color: #fff;
  --bs-btn-hover-color: #fff;
}

@media (max-width: 767.98px) {
  .create-button {
    position: static;
    width: 100%;
  }
}
</style>
