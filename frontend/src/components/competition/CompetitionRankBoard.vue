<template>
  <section class="rank-board" aria-labelledby="competition-rank-title">
    <header class="rank-board__header">
      <div>
        <h2 id="competition-rank-title">Rank 榜</h2>
      </div>
      <span v-if="rankRows.length > 0" class="rank-board__count">
        {{ rankRows.length }} 人
      </span>
    </header>

    <div v-if="loading" class="rank-list" role="status" aria-label="Rank 分榜加载中">
      <div v-for="index in 5" :key="index" class="rank-entry rank-entry--loading">
        <span class="rank-skeleton rank-skeleton--rank"></span>
        <span class="rank-skeleton rank-skeleton--avatar"></span>
        <span class="rank-skeleton rank-skeleton--name"></span>
        <span class="rank-skeleton rank-skeleton--score"></span>
      </div>
    </div>

    <div v-else-if="error" class="rank-state">
      <p>{{ error }}</p>
      <button type="button" class="rank-retry" @click="emit('retry')">
        重新加载
      </button>
    </div>

    <div v-else-if="rankRows.length === 0" class="rank-state">
      <p>暂无排名数据</p>
    </div>

    <ol v-else class="rank-list">
      <li
        v-for="row in rankRows"
        :key="row.key"
        class="rank-entry"
        :class="{ 'rank-entry--clickable': row.hasProfileLink }"
      >
        <component
          :is="row.hasProfileLink ? RouterLink : 'div'"
          v-bind="row.hasProfileLink ? { to: userProfileRoute(row), 'aria-label': `查看 ${row.userName} 的个人信息` } : {}"
          class="rank-entry__content"
          :class="{ 'rank-entry__content--link': row.hasProfileLink }"
        >
          <span
            class="rank-entry__rank"
            :class="{ 'rank-entry__rank--top': row.ranking <= 3 }"
            :aria-label="`第 ${row.ranking} 名`"
          >
            {{ row.ranking }}
          </span>

          <img
            v-if="row.avatar"
            :src="row.avatar"
            :alt="`${row.userName} 的头像`"
            class="rank-entry__avatar"
          >
          <span v-else class="rank-entry__avatar rank-entry__avatar--fallback">
            {{ avatarInitial(row.userName) }}
          </span>

          <div class="rank-entry__main">
            <span class="rank-entry__name" :style="{ color: row.nameColor }">
              {{ row.userName }}
            </span>
            <span v-if="row.attendedText" class="rank-entry__sub">
              参赛 {{ row.attendedText }} 场
            </span>
          </div>

          <div class="rank-entry__score" aria-label="Rating">
            <span>Rating</span>
            <strong>{{ row.rankScoreText }}</strong>
          </div>
        </component>
      </li>
    </ol>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { getRankScoreColor } from '@/components/RankColor'

const props = defineProps({
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

const emit = defineEmits(['retry'])

const rankRows = computed(() => props.rankings.map(normalizeRankRow))

function normalizeRankRow(item, index) {
  const ranking = index + 1
  const userId = item?.user_id ?? null
  const userName = item?.user_name || '匿名选手'
  const rankScore = normalizeOptionalNumber(item?.user_rating)
  const attended = normalizeOptionalNumber(item?.user_competition_count)

  return {
    key: userId ?? `${ranking}-${userName}`,
    userId,
    hasProfileLink: userId !== null && userId !== undefined && userId !== '',
    ranking,
    userName,
    avatar: item?.user_avatar || '',
    nameColor: getRankScoreColor(rankScore ?? 0),
    rankScoreText: rankScore === null ? '-' : formatNumber(rankScore),
    attendedText: attended === null ? '' : formatNumber(attended)
  }
}

function userProfileRoute(row) {
  return {
    name: 'UserInfo',
    params: {
      userId: row.userId
    }
  }
}

function normalizeOptionalNumber(value) {
  const number = Number(value)
  return Number.isFinite(number) ? number : null
}

function formatNumber(value) {
  return Number(value).toLocaleString('zh-CN', {
    maximumFractionDigits: 1
  })
}

function avatarInitial(name) {
  return (name || '选').trim().slice(0, 1).toUpperCase()
}
</script>

<style scoped>
.rank-board {
  padding: 1rem;
  border: 1px solid rgba(226, 232, 240, 0.76);
  border-radius: 1rem;
  background:
    linear-gradient(180deg, rgba(248, 250, 252, 0.92), rgba(255, 255, 255, 0) 8rem),
    #ffffff;
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.1);
}

.rank-board__header {
  display: flex;
  margin-bottom: 0.85rem;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
}

.rank-board__eyebrow {
  margin: 0 0 0.15rem;
  color: #64748b;
  font-size: 0.68rem;
  font-weight: 700;
  line-height: 1;
  text-transform: uppercase;
}

.rank-board__header h2 {
  margin: 0;
  color: #0f172a;
  font-size: 1rem;
  font-weight: 700;
  line-height: 1.25;
}

.rank-board__count {
  display: inline-flex;
  padding: 0.28rem 0.55rem;
  border-radius: 999px;
  background: #f1f5f9;
  color: #475569;
  font-size: 0.74rem;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
}

.rank-list {
  display: flex;
  padding: 0;
  margin: 0;
  flex-direction: column;
  gap: 0.42rem;
  list-style: none;
}

.rank-entry {
  display: flex;
  min-height: 4.85rem;
  padding: 0.62rem 0.95rem;
  align-items: center;
  gap: 0.85rem;
  border: 1px solid rgba(226, 232, 240, 0.6);
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 14px 34px rgba(15, 23, 42, 0.08);
  transition: border-color 0.18s ease, background-color 0.18s ease, transform 0.18s ease;
}

.rank-entry:hover {
  border-color: #a1a1aa;
  background: #ffffff;
  transform: translateY(-1px);
}

.rank-entry:focus-within {
  outline: 3px solid rgba(255, 161, 22, 0.28);
  outline-offset: 2px;
}

.rank-entry--clickable {
  cursor: pointer;
}

.rank-entry__content {
  display: flex;
  width: 100%;
  min-width: 0;
  align-items: center;
  gap: 0.85rem;
  color: inherit;
  text-decoration: none;
}

.rank-entry__content--link:hover {
  color: inherit;
}

.rank-entry__rank {
  display: inline-flex;
  min-width: 1.45rem;
  max-width: 3rem;
  height: 1.25rem;
  padding: 0 0.35rem;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 999px;
  background: #f1f5f9;
  color: #334155;
  font-size: 0.72rem;
  font-weight: 800;
  line-height: 1;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

.rank-entry__rank--top {
  background: #fff7ed;
  color: #d97706;
}

.rank-entry__avatar {
  width: 3.25rem;
  height: 3.25rem;
  flex: 0 0 auto;
  border-radius: 50%;
  object-fit: cover;
  outline: 1px solid rgba(226, 232, 240, 0.85);
  outline-offset: -0.5px;
}

.rank-entry__avatar--fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
  color: #64748b;
  font-size: 0.76rem;
  font-weight: 700;
}

.rank-entry__main {
  display: flex;
  min-width: 0;
  flex: 1 1 auto;
  flex-direction: column;
  gap: 0.1rem;
}

.rank-entry__name {
  overflow: hidden;
  font-size: 0.9rem;
  font-weight: 700;
  line-height: 1.3;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-entry__sub {
  overflow: hidden;
  color: #64748b;
  font-size: 0.72rem;
  font-weight: 600;
  line-height: 1.25;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-entry__score {
  display: inline-flex;
  min-width: 4.2rem;
  flex: 0 0 auto;
  flex-direction: column;
  align-items: flex-end;
  justify-content: center;
  gap: 0.06rem;
  font-size: 0.72rem;
  line-height: 1.2;
}

.rank-entry__score span {
  color: #64748b;
  font-weight: 600;
}

.rank-entry__score strong {
  max-width: 4.8rem;
  overflow: hidden;
  color: #1e293b;
  font-size: 0.86rem;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

.rank-state {
  display: flex;
  min-height: 9.5rem;
  padding: 1rem;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.85rem;
  color: #64748b;
  text-align: center;
}

.rank-state p {
  margin: 0;
  font-size: 0.86rem;
  font-weight: 600;
}

.rank-retry {
  padding: 0.42rem 0.75rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.55rem;
  background: #ffffff;
  color: #334155;
  font-size: 0.78rem;
  font-weight: 700;
  line-height: 1;
}

.rank-retry:hover {
  border-color: #94a3b8;
  background: #f8fafc;
}

.rank-entry--loading {
  pointer-events: none;
}

.rank-skeleton {
  display: inline-flex;
  overflow: hidden;
  border-radius: 999px;
  background: linear-gradient(90deg, #f1f5f9 0%, #e2e8f0 50%, #f1f5f9 100%);
  background-size: 200% 100%;
  animation: rank-skeleton 1.2s ease-in-out infinite;
}

.rank-skeleton--rank {
  width: 1.45rem;
  height: 1.25rem;
}

.rank-skeleton--avatar {
  width: 3.25rem;
  height: 3.25rem;
}

.rank-skeleton--name {
  height: 0.85rem;
  flex: 1 1 auto;
  min-width: 5rem;
}

.rank-skeleton--score {
  width: 3.4rem;
  height: 0.85rem;
}

@keyframes rank-skeleton {
  0% {
    background-position: 100% 0;
  }

  100% {
    background-position: -100% 0;
  }
}

@media (max-width: 575.98px) {
  .rank-board {
    padding: 0.85rem;
  }

  .rank-entry {
    min-height: 4.45rem;
    padding: 0.55rem 0.7rem;
    gap: 0.65rem;
  }

  .rank-entry__content {
    gap: 0.65rem;
  }

  .rank-entry__avatar,
  .rank-skeleton--avatar {
    width: 3rem;
    height: 3rem;
  }

  .rank-entry__score {
    min-width: 3.45rem;
  }
}
</style>
