<template>
  <article class="content-card">
    <nav class="content-card__tabs" aria-label="竞赛内容">
      <RouterLink
        class="content-card__tab"
        :to="tabTarget('')"
        :class="{ 'content-card__tab--active': activeTab === 'announcement' }"
        :aria-current="activeTab === 'announcement' ? 'page' : undefined"
      >
        公告
      </RouterLink>
      <RouterLink
        v-if="hasStarted"
        class="content-card__tab"
        :to="tabTarget('#problem')"
        :class="{ 'content-card__tab--active': activeTab === 'problems' }"
        :aria-current="activeTab === 'problems' ? 'page' : undefined"
      >
        题目
      </RouterLink>
      <RouterLink
        v-if="hasStarted"
        class="content-card__tab"
        :to="tabTarget('#submissions')"
        :class="{ 'content-card__tab--active': activeTab === 'submissions' }"
        :aria-current="activeTab === 'submissions' ? 'page' : undefined"
      >
        提交记录
      </RouterLink>
      <RouterLink
        v-if="hasStarted"
        class="content-card__tab"
        :to="tabTarget('#ranklist')"
        :class="{ 'content-card__tab--active': activeTab === 'ranking' }"
        :aria-current="activeTab === 'ranking' ? 'page' : undefined"
      >
        榜单
      </RouterLink>
    </nav>

    <section v-if="activeTab === 'announcement'" class="content-card__panel content-card__announcement" aria-label="竞赛公告内容">
      <p>{{ competition.announcement || '暂无公告内容。' }}</p>
    </section>

    <section v-else-if="activeTab === 'problems'" class="content-card__panel" aria-label="竞赛题目内容">
      <CompetitionProblems :competition-id="competition.id" />
    </section>

    <section v-else-if="activeTab === 'submissions'" class="content-card__panel" aria-label="竞赛提交记录">
      <CompetitionSubmissions :competition-id="competition.id" />
    </section>

    <section v-else class="content-card__panel" aria-label="竞赛榜单内容">
      <p v-if="rankingsLoading" class="content-card__state">榜单加载中...</p>
      <p v-else-if="rankingsError" class="content-card__state content-card__state--error">{{ rankingsError }}</p>
      <p v-else-if="rankings.length === 0" class="content-card__state">暂无排名数据。</p>
      <ol v-else class="content-card__ranking">
        <li v-for="item in rankings" :key="item.user_id || item.ranking">
          <strong class="content-card__place">{{ item.ranking }}</strong>
          <span>{{ item.user_name || '匿名选手' }}</span>
        </li>
      </ol>
    </section>

    <footer v-if="canCancel" class="content-card__registration">
      <span>
        <i class="bi bi-check-circle-fill" aria-hidden="true"></i>
        当前已报名本场竞赛
      </span>
      <button type="button" :disabled="canceling" @click="$emit('cancel')">
        <span v-if="canceling" class="spinner-border spinner-border-sm" aria-hidden="true"></span>
        {{ canceling ? '取消中...' : '取消报名' }}
      </button>
    </footer>
  </article>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import CompetitionProblems from './CompetitionProblems.vue'
import CompetitionSubmissions from './CompetitionSubmissions.vue'

const props = defineProps({
  competition: {
    type: Object,
    required: true
  },
  hasStarted: {
    type: Boolean,
    default: false
  },
  canCancel: {
    type: Boolean,
    default: false
  },
  canceling: {
    type: Boolean,
    default: false
  },
  rankings: {
    type: Array,
    default: () => []
  },
  rankingsLoading: {
    type: Boolean,
    default: false
  },
  rankingsError: {
    type: String,
    default: ''
  }
})

defineEmits(['cancel'])

const route = useRoute()
const router = useRouter()
const activeTab = computed(() => {
  if (props.hasStarted && route.hash === '#problem') {
    return 'problems'
  }

  if (props.hasStarted && route.hash === '#ranklist') {
    return 'ranking'
  }

  if (props.hasStarted && route.hash === '#submissions') {
    return 'submissions'
  }

  return 'announcement'
})

function tabTarget(hash) {
  return {
    name: 'CompetitionDetail',
    params: { competition_id: route.params.competition_id },
    hash
  }
}

watch(() => props.hasStarted, (hasStarted) => {
  if (!hasStarted && ['#problem', '#submissions', '#ranklist'].includes(route.hash)) {
    router.replace(tabTarget(''))
  }
}, { immediate: true })
</script>

<style scoped>
.content-card {
  overflow: hidden;
  background: #fff;
  border: 1px solid #e3eaf0;
  border-radius: 0.45rem;
  box-shadow: 0 0.42rem 1.45rem rgba(30, 45, 60, 0.05);
}

.content-card__tabs {
  display: flex;
  gap: 1.65rem;
  padding: 0 2.15rem;
  border-bottom: 1px solid #edf1f5;
}

.content-card__tab {
  position: relative;
  display: inline-flex;
  min-height: 4.6rem;
  align-items: center;
  padding: 0;
  color: #687a8f;
  font-size: 1.12rem;
  font-weight: 500;
  background: none;
  border: 0;
  text-decoration: none;
  transition: color 160ms ease;
}

.content-card__tab:hover,
.content-card__tab--active {
  color: #13866c;
}

.content-card__tab--active {
  font-weight: 600;
}

.content-card__tab--active::after {
  position: absolute;
  right: 0;
  bottom: -1px;
  left: 0;
  height: 0.18rem;
  content: "";
  background: #22ae90;
  border-radius: 999px 999px 0 0;
}

.content-card__panel {
  min-height: 9.4rem;
  padding: 2.15rem 2.15rem 2.4rem;
}

.content-card__announcement p {
  max-width: 46rem;
  margin: 0;
  color: #62748b;
  font-size: 1rem;
  line-height: 1.85;
  white-space: pre-line;
}

.content-card__state {
  margin: 0;
  color: #74869b;
  font-size: 0.98rem;
}

.content-card__state--error {
  color: #bf5454;
}

.content-card__ranking {
  display: grid;
  gap: 0.8rem;
  padding: 0;
  margin: 0;
  list-style: none;
}

.content-card__ranking li {
  display: flex;
  gap: 1rem;
  align-items: center;
  min-height: 3.15rem;
  padding: 0.65rem 1rem;
  color: #26354a;
  background: #f7fafb;
  border: 1px solid #edf1f5;
  border-radius: 0.52rem;
}

.content-card__place {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  color: #14886e;
  background: #e9f8f3;
  border-radius: 50%;
}

.content-card__registration {
  display: flex;
  gap: 1rem;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.7rem;
  color: #128068;
  font-size: 0.9rem;
  font-weight: 500;
  background: #f5fcf9;
  border-top: 1px solid #e2f2ed;
}

.content-card__registration span {
  display: inline-flex;
  gap: 0.45rem;
  align-items: center;
}

.content-card__registration button {
  display: inline-flex;
  gap: 0.42rem;
  align-items: center;
  min-height: 2.25rem;
  padding: 0 0.9rem;
  color: #5e7087;
  font-size: 0.87rem;
  font-weight: 500;
  background: #fff;
  border: 1px solid #d5dfe7;
  border-radius: 999px;
  transition: color 160ms ease, border-color 160ms ease, background 160ms ease;
}

.content-card__registration button:hover:not(:disabled) {
  color: #c44848;
  background: #fff8f8;
  border-color: #e7b8b8;
}

.content-card__registration button:disabled {
  opacity: 0.65;
}

.content-card__registration button:focus-visible {
  outline: 3px solid rgba(196, 72, 72, 0.18);
  outline-offset: 2px;
}

@media (max-width: 575.98px) {
  .content-card__tabs {
    gap: 1.3rem;
    padding: 0 1.15rem;
  }

  .content-card__panel {
    padding: 1.35rem 1.15rem 1.6rem;
  }

  .content-card__registration {
    align-items: stretch;
    flex-direction: column;
    padding: 1rem 1.15rem;
  }

  .content-card__registration button {
    justify-content: center;
  }
}
</style>
