<template>
  <main class="competition-page">
    <div class="container py-4 py-lg-5">
      <CompetitionHeader @create="goToCreate" />

      <div class="row g-4 align-items-start">
        <div class="col-lg-8">
          <CompetitionList
            :ongoing-competitions="ongoingCompetitions"
            :ended-competitions="endedCompetitions"
            :ongoing-loading="ongoingLoading"
            :ended-loading="endedLoading"
            :ongoing-error="ongoingError"
            :ended-error="endedError"
            :ended-page="endedPage"
            :ended-page-size="ENDED_PAGE_SIZE"
            :ended-page-count="endedPageCount"
            @select="goToCompetition"
            @retry-ongoing="loadOngoingCompetitions"
            @retry-ended="loadEndedCompetitions"
            @update-ended-page="updateEndedPage"
          />
        </div>

        <aside class="col-lg-4">
          <CompetitionRankBoard
            :rankings="rankings"
            :loading="rankingsLoading"
            :error="rankingsError"
            @retry="loadRankings"
          />
        </aside>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import api from '@/api'
import CompetitionHeader from './CompetitionHeader.vue'
import CompetitionList from './CompetitionList.vue'
import CompetitionRankBoard from './CompetitionRankBoard.vue'

const ENDED_PAGE_SIZE = 8
const router = useRouter()
const store = useStore()
const ongoingCompetitions = ref([])
const endedCompetitions = ref([])
const rankings = ref([])
const ongoingLoading = ref(false)
const endedLoading = ref(false)
const rankingsLoading = ref(false)
const ongoingError = ref('')
const endedError = ref('')
const rankingsError = ref('')
const endedPage = ref(1)
const endedTotal = ref(0)
const endedHasExactTotal = ref(false)
const endedPageCount = computed(() => {
  const hasPossibleNextPage = endedCompetitions.value.length >= ENDED_PAGE_SIZE

  if (hasPossibleNextPage) {
    return Math.max(endedPage.value + 1, Math.ceil(endedTotal.value / ENDED_PAGE_SIZE))
  }

  if (endedHasExactTotal.value) {
    return Math.max(1, Math.ceil(endedTotal.value / ENDED_PAGE_SIZE))
  }

  return Math.max(1, endedPage.value)
})

async function loadOngoingCompetitions() {
  ongoingLoading.value = true
  ongoingError.value = ''

  try {
    const resp = await api.getCompetitionsByEndStatus(false)
    const list = getCompetitionList(resp.data)
    if (resp.code === 0 && list) {
      ongoingCompetitions.value = list
    } else {
      ongoingError.value = resp.message || '未结束竞赛加载失败'
    }
  } catch (error) {
    ongoingError.value = error?.message || '未结束竞赛加载失败'
  } finally {
    ongoingLoading.value = false
  }
}

async function loadEndedCompetitions() {
  endedLoading.value = true
  endedError.value = ''

  try {
    const resp = await api.getCompetitionsByEndStatus(true, endedPage.value, ENDED_PAGE_SIZE)
    const list = getCompetitionList(resp.data)
    if (resp.code === 0 && list) {
      endedCompetitions.value = list
      const total = getCompetitionTotal(resp.data)
      endedHasExactTotal.value = total !== null
      endedTotal.value = total ?? list.length
    } else {
      endedError.value = resp.message || '已结束竞赛加载失败'
    }
  } catch (error) {
    endedError.value = error?.message || '已结束竞赛加载失败'
  } finally {
    endedLoading.value = false
  }
}

function getCompetitionList(data) {
  if (Array.isArray(data)) {
    return data
  }

  if (Array.isArray(data?.list)) {
    return data.list
  }

  if (Array.isArray(data?.competitions)) {
    return data.competitions
  }

  return null
}

function getCompetitionTotal(data) {
  const total = Number(data?.total ?? data?.count ?? data?.total_count ?? data?.pagination?.total)

  if (Number.isFinite(total) && total >= 0) {
    return total
  }

  return null
}

function updateEndedPage(page) {
  if (endedLoading.value || page === endedPage.value) {
    return
  }

  endedPage.value = page
  loadEndedCompetitions()
}

async function loadRankings() {
  rankingsLoading.value = true
  rankingsError.value = ''

  try {
    const resp = await api.getCompetitionRanking()
    if (resp.code === 0 && Array.isArray(resp.data)) {
      rankings.value = resp.data
      console.log('Rankings loaded:', resp.data)
    } else {
      rankingsError.value = resp.message || 'Rank 分榜加载失败'
    }
  } catch (error) {
    rankingsError.value = error?.message || 'Rank 分榜加载失败'
  } finally {
    rankingsLoading.value = false
  }
}

function goToCreate() {
  router.push({ name: 'CompetitionCreate' })
}

function goToCompetition(competition) {
  store.commit('SET_SELECTED_COMPETITION', competition)
  router.push({
    name: 'CompetitionDetail',
    params: { competition_id: competition.id }
  })
}

onMounted(() => {
  loadOngoingCompetitions()
  loadEndedCompetitions()
  loadRankings()
})
</script>

<style scoped>
.competition-page {
  min-height: calc(100vh - 60px);
  color: #262626;
}
</style>
