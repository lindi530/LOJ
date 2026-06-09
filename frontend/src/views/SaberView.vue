<template>
  <main class="saber-page">
    <section
      ref="stageRef"
      class="saber-page__stage"
      :style="{ backgroundImage: `url(${bgImage})` }"
    >
      <Menu
        v-if="currentView === 'menu'"
        :button-width="buttonWidth"
        :button-height="buttonHeight"
        :button-padding="buttonPadding"
        :button-font-size="buttonFontSize"
        :left-offset="menuLeftOffset"
        :top-offset="menuTopOffset"
        @heaven-battle="handleHeavenBattle"
        @friend-battle="handleFriendBattle"
      />

      <Battle
        v-if="currentView === 'battle'"
        :battle-type="currentBattleType"
        :post-id="postId"
        @set-room-id="handleSetRoomId"
        @back-to-menu="handleBackToMenu"
        @to-battle-game="handleToBattleGame"
      />

      <BattleGame
        v-if="currentView === 'match_success'"
        :room-id="roomId"
        :problem-id="problemId"
        @back-to-menu="handleBackToMenu"
      />
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import menuBg from '@/assets/background.jpg'
import Menu from '@/components/saber/menu.vue'
import Battle from '@/components/saber/BattleMatch.vue'
import BattleGame from '@/components/saber/BattleGame.vue'

const stageRef = ref(null)
const stageWidth = ref(1024)
const roomId = ref('')
const problemId = ref(null)
const currentView = ref('menu')
const currentBattleType = ref('')
const bgImage = ref(menuBg)
const postId = 0
let resizeObserver = null

const scaleRatio = computed(() => {
  const ratio = stageWidth.value / 1024
  return Math.max(0.65, Math.min(1.15, ratio))
})

const buttonWidth = computed(() => 180 * scaleRatio.value)
const buttonHeight = computed(() => 60 * scaleRatio.value)
const buttonPadding = computed(() => 12 * scaleRatio.value)
const buttonFontSize = computed(() => 20 * scaleRatio.value)
const menuLeftOffset = computed(() => 30 * scaleRatio.value)
const menuTopOffset = computed(() => 80 * scaleRatio.value)

function updateStageWidth() {
  if (!stageRef.value) return
  stageWidth.value = stageRef.value.getBoundingClientRect().width
}

function handleBackToMenu() {
  currentView.value = 'menu'
  bgImage.value = menuBg
}

function handleToBattleGame(room_id, problem_id) {
  roomId.value = room_id
  problemId.value = problem_id
  bgImage.value = menuBg
  currentView.value = 'match_success'
}

function handleHeavenBattle() {
  currentBattleType.value = '天人之战'
  currentView.value = 'battle'
  bgImage.value = menuBg
}

function handleFriendBattle() {
  currentBattleType.value = '好友对战'
  currentView.value = 'battle'
  bgImage.value = menuBg
}

function handleSetRoomId(value) {
  roomId.value = value
}

onMounted(() => {
  updateStageWidth()

  if (typeof ResizeObserver !== 'undefined' && stageRef.value) {
    resizeObserver = new ResizeObserver(updateStageWidth)
    resizeObserver.observe(stageRef.value)
  }

  window.addEventListener('resize', updateStageWidth)
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }

  window.removeEventListener('resize', updateStageWidth)
})
</script>

<style scoped>
.saber-page {
  width: 100%;
  height: calc(100vh - 60px);
  min-height: calc(100vh - 60px);
  box-sizing: border-box;
  overflow: hidden;
}

.saber-page__stage {
  position: relative;
  display: flex;
  width: 100%;
  height: 100%;
  min-height: inherit;
  overflow: hidden;
  box-sizing: border-box;
  background-position: center;
  background-size: cover;
}
</style>
