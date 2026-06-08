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
              <form
                v-if="!registered"
                class="registration-panel"
                novalidate
                @submit.prevent="requestEntry"
              >
                <div class="registration-panel__body">
                  <fieldset v-if="requiresCompetitionPassword" class="fieldset p-0">
                    <legend class="fieldset-legend">竞赛密码</legend>
                    <div
                      class="input input-lg w-full registration-password-field"
                      :class="{
                        'input-error': competitionPasswordError,
                        'input-success': hasCompetitionPassword && !competitionPasswordError
                      }"
                    >
                      <i class="bi bi-lock-fill opacity-60" aria-hidden="true"></i>
                      <input
                        ref="competitionPasswordInput"
                        v-model="competitionPassword"
                        :type="competitionPasswordVisible ? 'text' : 'password'"
                        class="grow"
                        placeholder="请输入竞赛密码"
                        autocomplete="current-password"
                        :disabled="entering"
                        :aria-invalid="competitionPasswordError ? 'true' : 'false'"
                        :aria-describedby="competitionPasswordError ? 'competition-entry-password-error' : undefined"
                        @input="clearCompetitionPasswordError"
                      >
                      <i
                        v-if="hasCompetitionPassword"
                        class="bi bi-check-circle-fill registration-password-status"
                        aria-label="已输入密码"
                      ></i>
                      <button
                        type="button"
                        class="registration-password-toggle"
                        :disabled="entering || !hasCompetitionPassword"
                        :aria-label="competitionPasswordVisible ? '隐藏密码' : '显示密码'"
                        @click="toggleCompetitionPasswordVisible"
                      >
                        <i
                          class="bi"
                          :class="competitionPasswordVisible ? 'bi-eye-slash' : 'bi-eye'"
                          aria-hidden="true"
                        ></i>
                      </button>
                    </div>
                    <p
                      v-if="competitionPasswordError"
                      id="competition-entry-password-error"
                      class="label text-error"
                    >
                      {{ competitionPasswordError }}
                    </p>
                  </fieldset>

                  <button
                    type="submit"
                    class="btn btn-success registration-submit"
                    :disabled="entering"
                  >
                    <span v-if="entering" class="loading loading-spinner loading-sm" aria-hidden="true"></span>
                    <i v-else class="bi bi-play-circle-fill" aria-hidden="true"></i>
                    {{ entering ? '报名中...' : '立即报名' }}
                  </button>
                </div>
              </form>

              <div v-else class="alert alert-success registration-success" role="status">
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
const competitionPasswordInput = ref(null)
const competitionPassword = ref('')
const competitionPasswordError = ref('')
const competitionPasswordVisible = ref(false)
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
const requiresCompetitionPassword = computed(() => (
  Boolean(competition.value?.password_hash) ||
  competition.value?.visibility === false ||
  Number(competition.value?.visibility) === 0
))
const hasCompetitionPassword = computed(() => Boolean(competitionPassword.value.trim()))

async function loadCompetition(detail) {
  const requestCompetitionId = String(competitionId.value || '')
  const requestId = ++competitionRequestId
  error.value = ''
  registered.value = false
  competitionPassword.value = ''
  competitionPasswordError.value = ''
  competitionPasswordVisible.value = false

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

function clearCompetitionPasswordError() {
  if (competitionPasswordError.value) {
    competitionPasswordError.value = ''
  }
}

function toggleCompetitionPasswordVisible() {
  if (hasCompetitionPassword.value && !entering.value) {
    competitionPasswordVisible.value = !competitionPasswordVisible.value
  }
}

function validateCompetitionPassword() {
  if (!requiresCompetitionPassword.value) {
    competitionPasswordError.value = ''
    return true
  }

  if (competitionPassword.value.trim()) {
    competitionPasswordError.value = ''
    return true
  }

  competitionPasswordError.value = '请输入竞赛密码后再报名。'
  competitionPasswordInput.value?.focus()
  return false
}

function requestEntry() {
  if (!validateCompetitionPassword()) {
    return
  }

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

  if (!validateCompetitionPassword()) {
    return
  }

  const payload = requiresCompetitionPassword.value
    ? { password: competitionPassword.value.trim() }
    : undefined

  entering.value = true
  try {
    const resp = await api.enterCompetition(competitionId.value, payload)
    if (resp.code !== 0) {
      message.error(resp.message || '报名失败，请稍后重试。')
      return
    }

    registered.value = true
    competitionPassword.value = ''
    competitionPasswordVisible.value = false
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

.registration-panel {
  width: 100%;
  padding: 1rem;
  background: #fff;
  border: 1px solid #dfe8ee;
  border-radius: 0.5rem;
  box-shadow: 0 0.32rem 1rem rgba(30, 45, 60, 0.045);
}

.registration-panel__body {
  display: grid;
  gap: 0.85rem;
}

.registration-panel .fieldset {
  display: grid;
  gap: 0.5rem;
  margin: 0;
  min-width: 0;
}

.registration-panel .fieldset-legend {
  margin: 0;
  padding: 0;
  color: #172438;
  font-size: 1rem;
  font-weight: 600;
}

.registration-panel .input {
  min-height: 3rem;
  color: #172438;
  background: #fff;
  border-color: #cfd9e3;
  border-radius: 0.45rem;
  box-shadow: none;
}

.registration-password-field {
  gap: 0.45rem;
}

.registration-password-field > i:first-child {
  color: #748398;
}

.registration-password-field input {
  color: #172438;
  background: transparent;
}

.registration-password-field input::placeholder {
  color: #8b97a8;
}

.registration-password-status {
  color: #1ba37f;
  font-size: 1rem;
}

.registration-password-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  color: #607289;
  background: transparent;
  border: 0;
  border-radius: 0.35rem;
}

.registration-password-toggle:hover:not(:disabled),
.registration-password-toggle:focus-visible {
  color: #178d74;
  background: #edf8f4;
}

.registration-password-toggle:disabled {
  color: #a7b3c1;
  cursor: not-allowed;
  opacity: 0.6;
}

.registration-submit {
  width: 100%;
  min-height: 3.05rem;
  color: #fff;
  font-size: 1.08rem;
  font-weight: 600;
  --bs-btn-bg: #22ae90;
  --bs-btn-border-color: #22ae90;
  --bs-btn-hover-bg: #178d74;
  --bs-btn-hover-border-color: #178d74;
  --bs-btn-active-bg: #147762;
  --bs-btn-active-border-color: #147762;
  --bs-btn-disabled-bg: #22ae90;
  --bs-btn-disabled-border-color: #22ae90;
}

.registration-success {
  width: 100%;
  justify-content: center;
  border-radius: 0.5rem;
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

  .registration-success {
    align-items: flex-start;
    line-height: 1.5;
  }
}
</style>
