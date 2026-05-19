<template>
  <div class="battle-game-container">
    <!-- 编辑器区域 -->
    <Edit 
      :room-id="roomId"
      :problem-id="problemId"
    />
    
    <!-- 游戏结束提示（覆盖在编辑器上） -->
    <div 
      v-if="gameOver" 
      class="game-over-overlay"
    >
      <div class="game-over-card">
        <!-- 动态显示胜利/失败图标 -->
        <div class="result-icon" :class="{ 'win-icon': win, 'lose-icon': !win }">
          {{ win ? '🏆' : '💔' }}
        </div>
        
        <!-- 动态显示胜利/失败标题 -->
        <h2 class="game-over-title" :class="{ 'win-title': win, 'lose-title': !win }">
          {{ win ? '恭喜你获胜！' : '很遗憾，你失败了' }}
        </h2>
        
        <!-- 动态显示胜利/失败提示 -->
        <div class="game-over-footer">
          <p>{{ win ? '你率先完成了题目' : '对手比你更快完成' }}</p>
        </div>
        
        <button 
          class="back-to-menu-btn"
          @click="handleBackToMenu"
        >
          返回菜单
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Edit from '@/views/ProblemView.vue'
import { useWebSocketContext } from '@/composables/useWebSocket.js'

const props = defineProps({
  roomId: {
    type: String,
    required: true,
  },
  problemId: {
    type: Number,
    required: true,
  },
});

const emit = defineEmits(['back-to-menu'])

const handleBackToMenu = () => {
  emit('back-to-menu');
}

const { registersaberResultCallback } = useWebSocketContext()

const gameOver = ref(false)
const overMessage = ref("")
const win = ref(null)


const unregister = registersaberResultCallback((msg) => { 
  overMessage.value = msg.content
  gameOver.value = true

  win.value = msg.win
  console.log("registersaberResultCallback: ", msg)
})

</script>

<style scoped>
.battle-game-container {
  height: 100%;
  width: 100%;
  overflow: hidden;
  position: relative;
  z-index: 1;
}

/* 游戏结束遮罩层 */
.game-over-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10; /* 确保在编辑器上方 */
  animation: fadeIn 0.5s ease-out;
}

/* 提示卡片 */
.game-over-card {
  background-color: white;
  border-radius: 12px;
  padding: 32px;
  text-align: center;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  animation: popUp 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

/* 结果图标样式 */
.result-icon {
  font-size: 64px;
  margin-bottom: 16px;
  animation: bounce 1s ease-in-out;
}

/* 胜利图标颜色 */
.win-icon {
  color: #FFD700; /* 金色 */
  text-shadow: 0 0 10px rgba(255, 215, 0, 0.5);
}

/* 失败图标颜色 */
.lose-icon {
  color: #666; /* 灰色 */
}

/* 标题样式 */
.game-over-title {
  font-size: 24px;
  margin-bottom: 12px;
  font-weight: bold;
}

/* 胜利标题颜色 */
.win-title {
  color: #D4AF37; /* 金色 */
}

/* 失败标题颜色 */
.lose-title {
  color: #333;
}

/* 消息样式 */
.game-over-message {
  font-size: 18px;
  color: #666;
  margin-bottom: 24px;
  min-height: 24px;
  padding: 0 16px;
}

/* 底部提示 */
.game-over-footer {
  font-size: 14px;
  color: #999;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

/* 返回菜单按钮 */
.back-to-menu-btn {
  margin-top: 16px;
  padding: 10px 24px;
  background-color: #337ab7;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 100%;
  max-width: 200px;
}

.back-to-menu-btn:hover {
  background-color: #286090;
  transform: translateY(-2px);
}

.back-to-menu-btn:active {
  transform: translateY(0);
}

/* 动画效果 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes popUp {
  0% { transform: scale(0.8); opacity: 0; }
  70% { transform: scale(1.05); }
  100% { transform: scale(1); opacity: 1; }
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-15px); }
}
</style>
