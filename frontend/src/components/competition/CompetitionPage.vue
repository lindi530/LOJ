<template>
  <main class="competition-page">
    <div class="container py-4 py-lg-5">
      <CompetitionHeader @create="goToCreate" />

      <div class="row g-4 align-items-start">
        <div class="col-lg-8">
          <CompetitionList
            :competitions="competitions"
            :loading="competitionsLoading"
            :error="competitionsError"
            @select="goToCompetition"
            @retry="loadCompetitions"
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
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import CompetitionHeader from './CompetitionHeader.vue'
import CompetitionList from './CompetitionList.vue'
import CompetitionRankBoard from './CompetitionRankBoard.vue'

const router = useRouter()
const competitions = ref([])
const rankings = ref([])
const competitionsLoading = ref(false)
const rankingsLoading = ref(false)
const competitionsError = ref('')
const rankingsError = ref('')

async function loadCompetitions() {
  competitionsLoading.value = true
  competitionsError.value = ''

  try {
    const resp = await api.getCompetitionList()
    if (resp.code === 0 && Array.isArray(resp.data)) {
      competitions.value = resp.data
    } else {
      competitionsError.value = resp.message || '竞赛列表加载失败'
    }
  } catch (error) {
    competitionsError.value = error?.message || '竞赛列表加载失败'
  } finally {
    competitionsLoading.value = false
  }
}

async function loadRankings() {
  rankingsLoading.value = true
  rankingsError.value = ''

  try {
    const resp = await api.getCompetitionRankList()
    if (resp.code === 0 && Array.isArray(resp.data)) {
      rankings.value = resp.data
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

function goToCompetition(competitionId) {
  router.push({
    name: 'CompetitionDetail',
    params: { competition_id: competitionId }
  })
}

onMounted(() => {
  loadCompetitions()
  loadRankings()
})
</script>

<style scoped>
.competition-page {
  min-height: calc(100vh - 60px);
  color: #262626;
}
</style>
