<template>
  <article class="competition-entry">
    <button
      type="button"
      class="competition-card"
      :class="{ 'competition-card--ended': isEnded }"
      @click="$emit('select', competition)"
    >
      <span class="competition-card__visual" aria-hidden="true">
        <span class="competition-card__trophy">
          <i class="bi bi-trophy-fill"></i>
        </span>
      </span>

      <span class="competition-card__content">
        <span class="competition-card__heading">
          <strong class="competition-card__title">{{ competition.name }}</strong>
          <span class="competition-card__access" :title="access.label" :aria-label="access.label">
            <i class="bi" :class="access.icon" aria-hidden="true"></i>
          </span>
          <span class="competition-card__status" :class="status.className">
            {{ status.label }}
          </span>
        </span>

        <span class="competition-card__schedule">
          <i class="bi bi-clock-fill" aria-hidden="true"></i>
          <span class="competition-card__label">比赛时间：</span>
          <span>{{ formatCompetitionSchedule(competition) }}</span>
          <span v-if="formatCompetitionDuration(competition)" class="competition-card__duration">
            (时长: {{ formatCompetitionDuration(competition) }})
          </span>
        </span>

        <span class="competition-card__details">
          <span class="competition-card__detail">
            <i class="bi bi-bookmark-star-fill" aria-hidden="true"></i>
            <span>
              <span class="competition-card__label">创建者：</span>
              {{ competition.creator_name || '未知创建者' }}
            </span>
          </span>
          <span class="competition-card__detail">
            <i class="bi bi-people-fill" aria-hidden="true"></i>
            <span>
              <span class="competition-card__label">参与人数：</span>
              {{ formatParticipantCount(competition.player_count) }}
            </span>
          </span>
        </span>
      </span>

      <span class="competition-card__action">
        <span class="competition-card__button">查看详情</span>
        <small v-if="status.hint">{{ status.hint }}</small>
      </span>
    </button>
  </article>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
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
  }
})

defineEmits(['select'])

const now = ref(Date.now())
let clockId = null

const status = computed(() => getStatus(props.competition, now.value))
const access = computed(() => getCompetitionAccess(props.competition))
const isEnded = computed(() => getCompetitionPhase(props.competition, now.value) === 'ended')

onMounted(() => {
  const endTime = new Date(props.competition.end_time).getTime()
  if (Number.isFinite(endTime) && now.value <= endTime) {
    clockId = window.setInterval(() => {
      now.value = Date.now()

      if (Number.isFinite(endTime) && now.value > endTime) {
        window.clearInterval(clockId)
        clockId = null
      }
    }, 1000)
  }
})

onUnmounted(() => {
  if (clockId) {
    window.clearInterval(clockId)
  }
})

function getStatus(competition, currentTime) {
  const phase = getCompetitionPhase(competition, currentTime)
  const endTime = new Date(competition.end_time).getTime()

  if (phase === 'pending') {
    return {
      label: '待定',
      hint: '时间待公布',
      className: 'competition-card__status--pending'
    }
  }

  if (phase === 'upcoming') {
    const startTime = new Date(competition.start_time).getTime()
    return {
      label: '报名中',
      hint: `距开始 ${formatCountdown(startTime - currentTime)}`,
      className: 'competition-card__status--upcoming'
    }
  }

  if (phase === 'running') {
    return {
      label: '进行中',
      hint: `距结束 ${formatCountdown(endTime - currentTime)}`,
      className: 'competition-card__status--running'
    }
  }

  return {
    label: '已结束',
    hint: '查看比赛结果',
    className: 'competition-card__status--ended'
  }
}
</script>

<style scoped>
.competition-entry {
  width: 100%;
}

.competition-card {
  position: relative;
  display: grid;
  grid-template-columns: 8.75rem minmax(0, 1fr) 7.75rem;
  gap: 1.25rem;
  align-items: center;
  width: 100%;
  min-height: 10.5rem;
  padding: 0 1.4rem 0 0;
  overflow: hidden;
  color: #26354a;
  text-align: left;
  background: #fff;
  border: 1px solid #e3eaf0;
  border-radius: 1rem;
  box-shadow: 0 0.3rem 1.25rem rgba(30, 45, 60, 0.05);
  transition: border-color 160ms ease, box-shadow 160ms ease, transform 160ms ease;
}

.competition-card::before {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  width: 0.32rem;
  content: "";
  background: #22ae90;
}

.competition-card--ended::before {
  background: #c6ced8;
}

.competition-card:hover {
  border-color: rgba(34, 174, 144, 0.48);
  box-shadow: 0 0.65rem 1.7rem rgba(19, 78, 69, 0.12);
  transform: translateY(-1px);
}

.competition-card:focus-visible {
  outline: 3px solid rgba(34, 174, 144, 0.33);
  outline-offset: 2px;
}

.competition-card__visual {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  align-self: stretch;
  min-height: 10.4rem;
  overflow: hidden;
  background:
    radial-gradient(circle at 70% 28%, rgba(255, 191, 70, 0.26), transparent 36%),
    linear-gradient(145deg, #eefcf8 0%, #dbfff3 100%);
  border-radius: 0.95rem 0 0 0.95rem;
}

.competition-card__visual::before,
.competition-card__visual::after {
  position: absolute;
  content: "";
  border-radius: 50%;
}

.competition-card__visual::before {
  right: -1.25rem;
  bottom: -1.7rem;
  width: 5.4rem;
  height: 5.4rem;
  background: #ffefc7;
}

.competition-card__visual::after {
  top: 0.8rem;
  left: 1rem;
  width: 0.58rem;
  height: 0.58rem;
  background: #22ae90;
}

.competition-card--ended .competition-card__visual {
  background:
    radial-gradient(circle at 70% 28%, rgba(166, 174, 184, 0.18), transparent 36%),
    linear-gradient(145deg, #f4f6f8 0%, #e9edf1 100%);
}

.competition-card--ended .competition-card__visual::before {
  background: #e1e6ec;
}

.competition-card--ended .competition-card__visual::after {
  background: #adb7c2;
}

.competition-card--ended .competition-card__trophy {
  color: #99a5b1;
  box-shadow: 0 0.6rem 1.1rem rgba(87, 101, 116, 0.11);
}

.competition-card__trophy {
  z-index: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 4.2rem;
  height: 4.2rem;
  color: #f3a20d;
  background: #fff;
  border-radius: 50%;
  box-shadow: 0 0.6rem 1.1rem rgba(199, 135, 19, 0.16);
  font-size: 1.95rem;
}

.competition-card__content,
.competition-card__heading,
.competition-card__schedule,
.competition-card__details,
.competition-card__detail,
.competition-card__action {
  display: flex;
}

.competition-card__content {
  min-width: 0;
  flex-direction: column;
  gap: 0.85rem;
}

.competition-card__heading {
  flex-wrap: wrap;
  align-items: center;
  gap: 0.65rem;
}

.competition-card__title {
  overflow: hidden;
  font-size: 1.22rem;
  font-weight: 600;
  line-height: 1.35;
  color: #172438;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.competition-card__status {
  display: inline-flex;
  align-items: center;
  height: 1.7rem;
  padding: 0 0.7rem;
  font-size: 0.8rem;
  font-weight: 600;
  border-radius: 0.38rem;
}

.competition-card__status--upcoming {
  color: #1476d4;
  background: #e8f3ff;
}

.competition-card__status--running {
  color: #068664;
  background: #dafbef;
}

.competition-card__status--ended {
  color: #fff;
  background: #b4b4b4;
}

.competition-card__status--pending {
  color: #ae7000;
  background: #fff1d5;
}

.competition-card__schedule {
  flex-wrap: wrap;
  align-items: center;
  gap: 0.38rem;
  color: #6f8095;
  font-size: 0.93rem;
  line-height: 1.5;
}

.competition-card__schedule i,
.competition-card__detail > i {
  width: 1.05rem;
  margin-right: 0.12rem;
  color: #92a0ae;
  text-align: center;
}

.competition-card__duration {
  color: #91a0b2;
}

.competition-card__details {
  flex-wrap: wrap;
  align-items: center;
  gap: 0.7rem 1.25rem;
  color: #6f8095;
  font-size: 0.91rem;
}

.competition-card__detail {
  gap: 0.45rem;
  align-items: center;
  white-space: nowrap;
}

.competition-card__label {
  margin-right: 0.18rem;
}

.competition-card__access {
  display: inline-flex;
  align-items: center;
  margin-left: -0.18rem;
  color: #7e8d9f;
  font-size: 1.22rem;
  line-height: 1.35;
}

.competition-card__access i {
  display: inline-flex;
  align-items: center;
  line-height: inherit;
}

.competition-card__access i::before {
  vertical-align: 0;
}

.competition-card__action {
  flex-direction: column;
  align-items: center;
  gap: 0.7rem;
  color: #7f8d9f;
}

.competition-card__button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 2.6rem;
  color: #fff;
  font-size: 0.98rem;
  font-weight: 400;
  background: #22ae90;
  border-radius: 999px;
  transition: background 160ms ease;
}

.competition-card:hover .competition-card__button {
  background: #1c967c;
}

.competition-card__action small {
  font-size: 0.78rem;
  text-align: center;
}

@media (max-width: 767.98px) {
  .competition-card {
    grid-template-columns: 3.4rem minmax(0, 1fr);
    gap: 0.85rem;
    padding: 1rem 1rem 1rem 0;
  }

  .competition-card::before {
    width: 0.25rem;
  }

  .competition-card__visual {
    align-self: start;
    height: 3.4rem;
    min-height: 0;
    border-radius: 0 0.7rem 0.7rem 0;
  }

  .competition-card__visual::before,
  .competition-card__visual::after {
    display: none;
  }

  .competition-card__trophy {
    width: 2.4rem;
    height: 2.4rem;
    font-size: 1.15rem;
  }

  .competition-card__title {
    font-size: 1.06rem;
  }

  .competition-card__access {
    font-size: 1.06rem;
  }

  .competition-card__details {
    align-items: flex-start;
    flex-direction: column;
    gap: 0.55rem;
  }

  .competition-card__action {
    grid-column: 2;
    flex-direction: row;
  }

  .competition-card__button {
    width: auto;
    min-height: 2.35rem;
    padding: 0 1.15rem;
  }
}
</style>
