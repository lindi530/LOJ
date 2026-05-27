<template>
  <section class="competition-overview" :class="{ 'competition-overview--compact': compactRunning }" aria-labelledby="competition-title">
    <header class="competition-overview__heading" :class="{ 'competition-overview__heading--compact': compactRunning }">
      <div class="competition-overview__title-row">
        <h1 id="competition-title">{{ competition.name }}</h1>
        <span v-if="compactRunning" class="competition-overview__status competition-overview__status--running">
          进行中
        </span>
      </div>

      <div v-if="!compactRunning" class="competition-overview__labels">
        <span class="competition-overview__status" :class="status.className">
          {{ status.label }}
        </span>
        <span class="competition-overview__access" :title="access.label">
          <i class="bi" :class="access.icon" aria-hidden="true"></i>
          <span>{{ access.label }}</span>
        </span>
      </div>

      <p v-else class="competition-overview__countdown" aria-live="polite">
        距离比赛结束
        <strong>{{ endCountdown }}</strong>
      </p>
    </header>

    <dl v-if="!compactRunning" class="competition-overview__facts">
      <div>
        <dt><i class="bi bi-clock" aria-hidden="true"></i>时间</dt>
        <dd>
          {{ formatCompetitionSchedule(competition) }}
          <span v-if="duration" class="competition-overview__duration">（时长: {{ duration }}）</span>
        </dd>
      </div>
      <div>
        <dt><i class="bi bi-person" aria-hidden="true"></i>创建者</dt>
        <dd>{{ competition.creator_name || '未知创建者' }}</dd>
      </div>
      <div>
        <dt><i class="bi bi-people" aria-hidden="true"></i>参赛人数</dt>
        <dd>{{ formatParticipantCount(competition.player_count) }}</dd>
      </div>
    </dl>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import {
  formatCompetitionDuration,
  formatCompetitionSchedule,
  formatCountdown,
  formatParticipantCount,
  getCompetitionAccess,
  getCompetitionPhase
} from '@/utils/competitionTime'

const props = defineProps({
  competition: {
    type: Object,
    required: true
  },
  currentTime: {
    type: Number,
    required: true
  },
  compactRunning: {
    type: Boolean,
    default: false
  }
})

const access = computed(() => getCompetitionAccess(props.competition))

const status = computed(() => {
  const phase = getCompetitionPhase(props.competition, props.currentTime)
  if (phase === 'pending') {
    return { label: '待定', className: 'competition-overview__status--pending' }
  }

  if (phase === 'upcoming') {
    return { label: '报名中', className: 'competition-overview__status--upcoming' }
  }

  if (phase === 'running') {
    return { label: '进行中', className: 'competition-overview__status--running' }
  }

  return { label: '已结束', className: 'competition-overview__status--ended' }
})

const duration = computed(() => formatCompetitionDuration(props.competition))
const endCountdown = computed(() => {
  const endTime = new Date(props.competition.end_time).getTime()
  return formatCountdown(endTime - props.currentTime)
})
</script>

<style scoped>
.competition-overview {
  padding: 2.45rem 2.6rem 0.95rem;
  overflow: hidden;
  background: #fff;
  border: 1px solid #e3eaf0;
  border-radius: 0.45rem;
  box-shadow: 0 0.42rem 1.45rem rgba(30, 45, 60, 0.05);
}

.competition-overview--compact {
  padding-bottom: 2.45rem;
}

.competition-overview__heading {
  margin-bottom: 1.55rem;
}

.competition-overview__heading--compact {
  display: flex;
  gap: 1.5rem;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0;
}

.competition-overview__title-row {
  display: flex;
  gap: 1.2rem;
  align-items: center;
  min-width: 0;
}

.competition-overview__heading h1 {
  margin: 0 0 0.85rem;
  color: #172438;
  font-size: clamp(1.75rem, 1.3rem + 1.8vw, 2.55rem);
  font-weight: 700;
  line-height: 1.25;
}

.competition-overview__heading--compact h1 {
  margin-bottom: 0;
}

.competition-overview__labels {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  align-items: center;
}

.competition-overview__access,
.competition-overview__status {
  display: inline-flex;
  gap: 0.32rem;
  align-items: center;
  height: 2.55rem;
  padding: 0 0.92rem;
  border-radius: 0.48rem;
  font-size: 0.98rem;
  font-weight: 600;
}

.competition-overview__countdown {
  display: inline-flex;
  flex-shrink: 0;
  gap: 0.7rem;
  align-items: baseline;
  margin: 0;
  color: #657990;
  font-size: 0.98rem;
}

.competition-overview__countdown strong {
  color: #087b61;
  font-size: 1.35rem;
  font-variant-numeric: tabular-nums;
}

.competition-overview__access {
  color: #687a8f;
  background: transparent;
}

.competition-overview__status--upcoming {
  color: #068664;
  background: #effbf7;
  border: 1px solid #22ae90;
}

.competition-overview__status--running {
  color: #068664;
  background: #dafbef;
}

.competition-overview__status--ended {
  color: #fff;
  background: #b4b4b4;
}

.competition-overview__status--pending {
  color: #ae7000;
  background: #fff1d5;
}

.competition-overview__duration {
  color: #8d9caf;
}

.competition-overview__facts {
  margin: 0;
}

.competition-overview__facts div {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 4.45rem;
  color: #566982;
  font-size: 1rem;
  border-top: 1px solid #edf1f5;
}

.competition-overview__facts dt {
  display: flex;
  gap: 1.25rem;
  align-items: center;
  font-weight: 400;
}

.competition-overview__facts dt i {
  width: 1.35rem;
  color: #556981;
  font-size: 1.45rem;
  text-align: center;
}

.competition-overview__facts dd {
  margin: 0;
  color: #26354a;
  font-size: 1.02rem;
  font-weight: 400;
  text-align: right;
}

@media (max-width: 575.98px) {
  .competition-overview {
    padding: 1.3rem 1.15rem 0.35rem;
  }

  .competition-overview--compact {
    padding-bottom: 1.3rem;
  }

  .competition-overview__heading {
    margin-bottom: 1rem;
  }

  .competition-overview__heading--compact {
    align-items: flex-start;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 0;
  }

  .competition-overview__title-row {
    flex-wrap: wrap;
    gap: 0.8rem;
  }

  .competition-overview__facts div {
    align-items: flex-start;
    flex-direction: column;
    gap: 0.6rem;
    justify-content: center;
    min-height: 4.8rem;
  }

  .competition-overview__facts dt {
    gap: 0.65rem;
  }

  .competition-overview__facts dd {
    padding-left: 2rem;
    line-height: 1.6;
    text-align: left;
  }
}
</style>
