<template>
  <section class="card border-0 shadow-sm rank-card">
    <div class="card-header bg-white border-bottom px-4 py-3">
      <h2 class="h5 fw-semibold mb-1">Rank 分榜</h2>
      <p class="small text-secondary mb-0">按竞赛积分排名，而非实时名次</p>
    </div>

    <div v-if="loading" class="list-group list-group-flush" aria-label="Rank 分榜加载中">
      <div v-for="index in 5" :key="index" class="list-group-item px-4 py-3">
        <p class="placeholder-glow mb-0">
          <span class="placeholder col-8"></span>
        </p>
      </div>
    </div>

    <div v-else-if="error" class="p-4 text-center">
      <p class="text-secondary small mb-3">{{ error }}</p>
      <button type="button" class="btn btn-outline-secondary btn-sm" @click="$emit('retry')">
        重新加载
      </button>
    </div>

    <div v-else-if="rankings.length === 0" class="card-body text-center py-5 text-secondary">
      暂无排名数据
    </div>

    <ol v-else class="list-group list-group-flush list-unstyled mb-0">
      <li
        v-for="item in rankings"
        :key="item.user_id"
        class="list-group-item d-flex align-items-center gap-3 px-4 py-3"
      >
        <span class="rank-number fw-semibold text-center" :class="rankClass(item.ranking)">
          {{ item.ranking }}
        </span>
        <img
          v-if="item.avatar"
          :src="item.avatar"
          :alt="`${item.user_name} 的头像`"
          class="rounded-circle rank-avatar"
        >
        <span v-else class="rounded-circle rank-avatar rank-avatar-fallback">
          {{ avatarInitial(item.user_name) }}
        </span>
        <span class="fw-medium text-truncate">{{ item.user_name || '匿名选手' }}</span>
      </li>
    </ol>
  </section>
</template>

<script setup>
defineProps({
  rankings: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  }
})

defineEmits(['retry'])

function rankClass(ranking) {
  return {
    1: 'rank-first',
    2: 'rank-second',
    3: 'rank-third'
  }[ranking] || 'text-secondary'
}

function avatarInitial(name) {
  return (name || '选').slice(0, 1).toUpperCase()
}
</script>

<style scoped>
.rank-card {
  border-radius: 0.75rem;
  overflow: hidden;
}

.rank-number {
  width: 1.5rem;
}

.rank-first {
  color: #ffa116;
}

.rank-second {
  color: #6c757d;
}

.rank-third {
  color: #b86b34;
}

.rank-avatar {
  width: 2.25rem;
  height: 2.25rem;
  object-fit: cover;
}

.rank-avatar-fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #f1f3f5;
  color: #6c757d;
}
</style>
