<template>
  <main class="competition-detail-page">
    <div class="competition-detail-layout">
      <RouterLink class="back-link" :to="{ name: 'Competition' }">
        <i class="bi bi-arrow-left" aria-hidden="true"></i>
        返回竞赛列表
      </RouterLink>

      <div class="competition-detail-shell">
        <section v-if="loading" class="detail-state text-center" role="status">
          <span class="spinner-border" aria-hidden="true"></span>
          <p class="mb-0 mt-3">竞赛详情加载中...</p>
        </section>

        <section v-else-if="error" class="detail-state text-center" role="alert">
          <i class="bi bi-exclamation-circle detail-state__icon" aria-hidden="true"></i>
          <p class="mb-0">{{ error }}</p>
        </section>

        <template v-else-if="competition">
          <div class="competition-detail-cards">
            <CompetitionOverview
              :competition="competition"
              :current-time="now"
              :compact-running="isRunning && registered"
            />

            <section v-if="isUpcoming || (isRunning && !registered)" class="registration-stage" aria-label="竞赛报名">
              <button
                v-if="!registered"
                type="button"
                class="enter-button"
                :disabled="entering"
                @click="requestEntry"
              >
                <span v-if="entering" class="spinner-border spinner-border-sm" aria-hidden="true"></span>
                {{ entering ? '报名中...' : '立即报名' }}
              </button>

              <div v-else class="registration-success" role="status">
                <i class="bi bi-check-circle-fill" aria-hidden="true"></i>
                <span>您已成功报名，请留意公告与开赛时间。</span>
              </div>
            </section>

            <CompetitionContentCard
              :competition="competition"
              :has-started="canViewStartedContent"
              :can-cancel="isUpcoming && registered"
              :canceling="canceling"
              @cancel="requestCancellation"
            />
          </div>
        </template>
      </div>
    </div>

    <LoginModal v-model:visible="loginVisible" @login-success="continueAfterLogin" />
  </main>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import api from '@/api'
import LoginModal from '@/components/account/LoginModal.vue'
import { getCompetitionPhase } from '@/utils/competitionTime'
import CompetitionContentCard from './CompetitionContentCard.vue'
import CompetitionOverview from './CompetitionOverview.vue'

const props = defineProps({
  competition: {
    type: Object,
    default: null
  }
})

const route = useRoute()
const store = useStore()
const message = useMessage()
const competition = ref(null)
const loading = ref(false)
const error = ref('')
const registered = ref(false)
const entering = ref(false)
const canceling = ref(false)
const loginVisible = ref(false)
const actionAfterLogin = ref('')
const now = ref(Date.now())
let competitionRequestId = 0
let clockId = null

const competitionId = computed(() => route.params.competition_id)
const isLogin = computed(() => store.getters['user/isLogin'])
const startTime = computed(() => new Date(competition.value?.start_time).getTime())
const isUpcoming = computed(() => Number.isFinite(startTime.value) && now.value < startTime.value)
const hasStarted = computed(() => Number.isFinite(startTime.value) && now.value >= startTime.value)
const isRunning = computed(() => competition.value && getCompetitionPhase(competition.value, now.value) === 'running')
const canViewStartedContent = computed(() => hasStarted.value && (!isRunning.value || registered.value))

async function loadCompetition(detail) {
  const requestCompetitionId = String(competitionId.value || '')
  const requestId = ++competitionRequestId
  error.value = ''
  registered.value = false

  if (!requestCompetitionId) {
    loading.value = false
    competition.value = null
    error.value = '未获取到竞赛编号'
    return
  }

  if (detail && String(detail.id) === requestCompetitionId) {
    competition.value = { ...detail }
  } else {
    competition.value = null
  }

  await loadCompetitionDetail(requestCompetitionId, requestId)
}

async function loadCompetitionDetail(requestCompetitionId, requestId) {
  loading.value = !competition.value

  try {
    const resp = await api.getCompetition(requestCompetitionId)
    if (requestId !== competitionRequestId) {
      return
    }

    const detail = getCompetitionDetail(resp.data)
    if (resp.code !== 0 || !detail) {
      competition.value = null
      error.value = resp.message || '竞赛详情加载失败，请稍后重试'
      return
    }

    competition.value = { ...detail }
    await loadRegisteredStatus(requestCompetitionId)
  } catch (requestError) {
    if (requestId === competitionRequestId) {
      competition.value = null
      error.value = requestError?.message || '竞赛详情加载失败，请稍后重试'
    }
    return
  } finally {
    if (requestId === competitionRequestId) {
      loading.value = false
    }
  }
}

function getCompetitionDetail(data) {
  return data?.competition || data
}

async function loadRegisteredStatus(requestCompetitionId = competitionId.value) {
  registered.value = false
  if (!isLogin.value) {
    return
  }

  const resp = await api.hasEnteredCompetition(requestCompetitionId)
  if (resp.code !== 0) {
    throw new Error(resp.message || '报名状态加载失败')
  }

  registered.value = resp.data?.has_enter === true
}

function requireLogin(action) {
  if (isLogin.value) {
    return true
  }

  actionAfterLogin.value = action
  loginVisible.value = true
  return false
}

function requestEntry() {
  if (requireLogin('entry')) {
    enterCompetition()
  }
}

function requestCancellation() {
  if (requireLogin('cancel')) {
    cancelEntry()
  }
}

function continueAfterLogin() {
  const nextAction = actionAfterLogin.value
  actionAfterLogin.value = ''

  if (nextAction === 'entry') {
    enterCompetition()
  } else if (nextAction === 'cancel') {
    cancelEntry()
  }
}

async function enterCompetition() {
  if (entering.value || !competition.value) {
    return
  }

  entering.value = true
  try {
    const resp = await api.enterCompetition(competitionId.value)
    if (resp.code !== 0) {
      message.error(resp.message || '报名失败，请稍后重试。')
      return
    }

    registered.value = true
    changePlayerCount(1)
    message.success('报名成功')
  } catch (requestError) {
    message.error(requestError?.message || '报名失败，请稍后重试。')
  } finally {
    entering.value = false
  }
}

async function cancelEntry() {
  if (canceling.value || !competition.value) {
    return
  }

  canceling.value = true
  try {
    const resp = await api.cancelCompetitionEntry(competitionId.value)
    if (resp.code !== 0) {
      message.error(resp.message || '取消报名失败，请稍后重试。')
      return
    }

    registered.value = false
    changePlayerCount(-1)
    message.success('已取消报名')
  } catch (requestError) {
    message.error(requestError?.message || '取消报名失败，请稍后重试。')
  } finally {
    canceling.value = false
  }
}

function changePlayerCount(amount) {
  const currentCount = Number(competition.value.player_count)
  if (Number.isFinite(currentCount)) {
    competition.value.player_count = Math.max(0, currentCount + amount)
  }
}

watch([competitionId, () => props.competition], ([, detail]) => {
  loadCompetition(detail)
}, { immediate: true })

onMounted(() => {
  clockId = window.setInterval(() => {
    now.value = Date.now()
  }, 1000)
})

onUnmounted(() => {
  if (clockId) {
    window.clearInterval(clockId)
  }
})
</script>

<style scoped>
.competition-detail-page {
  min-height: calc(100vh - 60px);
  color: #172438;
}

.competition-detail-layout {
  display: grid;
  grid-template-columns: clamp(8.75rem, 12vw, 12rem) minmax(0, 76.75rem) clamp(8.75rem, 12vw, 12rem);
  gap: 1.25rem;
  align-items: start;
  justify-content: center;
  padding: 1.25rem 1.5rem 2.5rem;
}

.competition-detail-shell {
  grid-column: 2;
  min-width: 0;
}

.competition-detail-cards {
  display: grid;
  gap: 1.35rem;
}

.back-link {
  grid-column: 1;
  display: inline-flex;
  gap: 0.48rem;
  align-items: center;
  justify-self: start;
  margin-top: 0.8rem;
  color: #607289;
  font-size: 0.93rem;
  font-weight: 500;
  text-decoration: none;
  transition: color 160ms ease;
}

.back-link:hover {
  color: #17896f;
}

.detail-state {
  padding: 3.75rem 1.5rem;
  color: #697c92;
  background: #fff;
  border: 1px solid #e3eaf0;
  border-radius: 1rem;
  box-shadow: 0 0.4rem 1.4rem rgba(30, 45, 60, 0.05);
}

.detail-state__icon {
  display: block;
  margin-bottom: 1rem;
  color: #96a5b6;
  font-size: 2rem;
}

.registration-stage {
  display: flex;
  justify-content: center;
  margin: 0.65rem 0;
}

.enter-button {
  display: inline-flex;
  gap: 0.6rem;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 5rem;
  color: #fff;
  font-size: 1.5rem;
  font-weight: 600;
  background: #22ae90;
  border: 0;
  border-radius: 0.45rem;
  box-shadow: 0 0.7rem 1.35rem rgba(34, 174, 144, 0.25);
  transition: background 160ms ease, box-shadow 160ms ease, transform 160ms ease;
}

.enter-button:hover:not(:disabled) {
  background: #178d74;
  box-shadow: 0 0.85rem 1.5rem rgba(23, 141, 116, 0.3);
  transform: translateY(-1px);
}

.enter-button:disabled {
  opacity: 0.7;
}

.enter-button:focus-visible {
  outline: 3px solid rgba(34, 174, 144, 0.27);
  outline-offset: 2px;
}

.registration-success {
  display: inline-flex;
  gap: 0.6rem;
  align-items: center;
  width: 100%;
  justify-content: center;
  min-height: 4rem;
  padding: 0.94rem 1.25rem;
  color: #087b61;
  font-size: 0.95rem;
  font-weight: 500;
  background: #e9fbf5;
  border: 1px solid #c5f0e3;
  border-radius: 0.45rem;
}

.registration-success i {
  font-size: 1.15rem;
}

@media (max-width: 991.98px) {
  .competition-detail-layout {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
  }

  .competition-detail-shell {
    width: 100%;
  }

  .back-link {
    margin-top: 0;
  }

  .competition-detail-cards {
    gap: 1rem;
  }
}

@media (max-width: 575.98px) {
  .registration-stage {
    margin: 0.2rem 0;
  }

  .enter-button {
    min-height: 3.45rem;
    font-size: 1.08rem;
  }

  .registration-success {
    align-items: flex-start;
    border-radius: 0.85rem;
    line-height: 1.5;
  }
}
</style>
