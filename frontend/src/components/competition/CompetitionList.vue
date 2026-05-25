<template>
  <section class="card border-0 shadow-sm contest-card">
    <div class="card-header bg-white border-bottom d-flex justify-content-between align-items-center px-4 py-3">
      <h2 class="h5 fw-semibold mb-0">竞赛列表</h2>
      <span class="small text-secondary">全部竞赛</span>
    </div>

    <div v-if="loading" class="list-group list-group-flush" aria-label="竞赛列表加载中">
      <div v-for="index in 3" :key="index" class="list-group-item px-4 py-4">
        <p class="placeholder-glow mb-2">
          <span class="placeholder col-5"></span>
        </p>
        <p class="placeholder-glow mb-0">
          <span class="placeholder col-7 placeholder-sm"></span>
        </p>
      </div>
    </div>

    <div v-else-if="error" class="p-4 text-center">
      <p class="text-secondary mb-3">{{ error }}</p>
      <button type="button" class="btn btn-outline-secondary btn-sm" @click="$emit('retry')">
        重新加载
      </button>
    </div>

    <div v-else-if="competitions.length === 0" class="card-body text-center py-5 text-secondary">
      <i class="bi bi-calendar2-event fs-2 d-block mb-2" aria-hidden="true"></i>
      暂无竞赛
    </div>

    <div v-else class="list-group list-group-flush">
      <button
        v-for="competition in competitions"
        :key="competition.id"
        type="button"
        class="list-group-item list-group-item-action d-flex gap-3 px-4 py-3"
        @click="$emit('select', competition.id)"
      >
        <span class="contest-icon flex-shrink-0">
          <i class="bi bi-trophy" aria-hidden="true"></i>
        </span>

        <span class="flex-grow-1 text-start overflow-hidden">
          <span class="d-flex flex-wrap align-items-center gap-2 mb-1">
            <span class="fw-semibold text-dark text-truncate">{{ competition.name }}</span>
            <span class="badge rounded-pill" :class="getStatus(competition).classes">
              {{ getStatus(competition).label }}
            </span>
          </span>
          <span class="small text-secondary d-block mb-2">{{ formatSchedule(competition) }}</span>
          <span class="d-flex flex-wrap gap-2">
            <span class="badge text-bg-light fw-normal text-secondary">
              <i class="bi bi-person me-1" aria-hidden="true"></i>{{ competition.creator_name || '未知创建者' }}
            </span>
            <span class="badge text-bg-light fw-normal text-secondary">
              <i class="bi me-1" :class="getAccess(competition).icon" aria-hidden="true"></i>
              {{ getAccess(competition).label }}
            </span>
          </span>
        </span>

        <span class="participants text-end flex-shrink-0">
          <strong class="d-block fw-semibold text-dark">{{ formatNumber(competition.player_count) }}</strong>
          <small class="text-secondary">参赛者</small>
        </span>
      </button>
    </div>
  </section>
</template>

<script setup>
defineProps({
  competitions: {
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

defineEmits(['select', 'retry'])

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  month: 'short',
  day: 'numeric',
  weekday: 'short',
  hour: '2-digit',
  minute: '2-digit'
})

function formatSchedule(competition) {
  const start = formatTime(competition.start_time)
  const end = formatTime(competition.end_time)
  return start && end ? `${start} - ${end}` : '比赛时间待定'
}

function formatTime(value) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '' : dateFormatter.format(date)
}

function formatNumber(value) {
  return Number(value || 0).toLocaleString('zh-CN')
}

function getAccess(competition) {
  if (competition.password_hash) {
    return { label: '需密码', icon: 'bi-lock' }
  }

  if (Number(competition.visibility) === 1) {
    return { label: '公开竞赛', icon: 'bi-globe2' }
  }

  return { label: '非公开', icon: 'bi-eye-slash' }
}

function getStatus(competition) {
  const now = Date.now()
  const startTime = new Date(competition.start_time).getTime()
  const endTime = new Date(competition.end_time).getTime()

  if (!Number.isFinite(startTime) || !Number.isFinite(endTime)) {
    return { label: '待定', classes: 'text-bg-secondary' }
  }

  if (now < startTime) {
    return { label: '即将开始', classes: 'status-upcoming' }
  }

  if (now <= endTime) {
    return { label: '进行中', classes: 'status-running' }
  }

  return { label: '已结束', classes: 'text-bg-light text-secondary' }
}
</script>

<style scoped>
.contest-card {
  border-radius: 0.75rem;
  overflow: hidden;
}

.contest-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.75rem;
  height: 2.75rem;
  border-radius: 0.65rem;
  background: #fff4e3;
  color: #ffa116;
  font-size: 1.25rem;
}

.list-group-item-action:hover {
  background-color: #fafafa;
}

.status-upcoming {
  background: #fff4e3;
  color: #bd7100;
}

.status-running {
  background: #e4f7ec;
  color: #168143;
}

.participants {
  min-width: 3.9rem;
}
</style>
