<template>
  <div class="vs-wrapper">
    <!-- 天人之战时显示VS内容 -->
    <template v-if="battleType === '天人之战'">
      <div class="match-status-bar" v-if="isMatching">
        <div class="status-content">
          <div class="status-row">
            <div class="loader-small">
              <div class="loader-dot"></div>
              <div class="loader-dot"></div>
              <div class="loader-dot"></div>
            </div>
            <span class="status-text">匹配中</span>
          </div>

          <div class="status-row">
            <div class="match-time-display">
              <span class="time-value">{{ matchTime }}</span>
            </div>
          </div>

          <div class="status-row">
            <div class="players-count">
              <span class="players-label">已找到:</span>
              <span class="players-value">{{ foundPlayers }}/2</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="vs-container">
        <div class="vs-background"></div>
        <div class="vs-text">VS</div>
      </div>
    </template>
    <template v-else>
    </template>

  </div>
</template>

<script setup>
const props = defineProps({
  isMatching: {
    type: Boolean,
    default: false
  },
  matchTime: {
    type: String,
    default: "00:00"
  },
  foundPlayers: {
    type: Number,
    default: 0
  },
  battleType: {  // 新增：接收对战类型参数
    type: String,
    required: true
  }
});

</script>

<style scoped>
/* VS区域样式 */
.vs-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  background: transparent;
}

/* 匹配状态条（匹配中显示） */
.match-status-bar {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
  padding: 3px 8px;
  background-color: rgba(15, 52, 96, 0.5);
  border-radius: 6px;
  border: 1px solid rgba(58, 134, 255, 0.3);
  white-space: nowrap;
  animation: slideDown 0.3s ease-out;
  z-index: 5;
  margin-bottom: 5px;
}

.status-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.status-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.loader-small {
  display: flex;
  gap: 4px;
}

@keyframes slideDown {
  from { transform: translateY(-15px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

/* 加载动画（三个点） */
.loader-small {
  display: flex;
  gap: 3px;
  align-items: center;
  height: 14px;
}

.loader-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #3a86ff;
  animation: bounceSmall 1.5s infinite ease-in-out;
}

.loader-dot:nth-child(2) {
  animation-delay: 0.2s;
  background-color: #8338ec;
}

.loader-dot:nth-child(3) {
  animation-delay: 0.4s;
  background-color: #4ade80;
}

@keyframes bounceSmall {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-3px); }
}

.status-text {
  color: #e0e7ff;
  font-size: 12px;
  white-space: nowrap;
}

/* 匹配时间显示 */
.match-time-display {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 1px 4px;
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.time-value {
  color: #409eff;
  font-size: 24px; 
  font-weight: 600;
  font-family: 'Courier New', monospace;
  min-width: 40px;
  text-align: center;
  text-shadow: 0 0 5px rgba(58, 134, 255, 0.3);
  animation: pulse 2s infinite alternate;
}

@keyframes pulse {
  from { opacity: 0.9; }
  to { opacity: 1; }
}

/* 玩家计数显示 */
.players-count {
  display: flex;
  align-items: center;
  gap: 3px;
}

.players-label {
  color: #94a3b8;
  font-size: 11px;
}

.players-value {
  color: #4ade80;
  font-weight: 600;
  font-size: 12px;
}

/* VS文字区域 */
.vs-container {
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  width: 80px;
  height: 80px;
  z-index: 4;
  background: transparent;
}

.vs-background {
  position: absolute;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, #ff3a3a 0%, #e94560 50%, rgba(233, 69, 96, 0.2) 100%);
  border-radius: 50%;
  box-shadow: 0 0 15px rgba(255, 58, 58, 0.6), 
              inset 0 0 10px rgba(255, 255, 255, 0.3);
  animation: vsPulse 2s infinite alternate;
  z-index: 1;
}

.vs-text {
  position: relative;
  color: white;
  font-size: 28px;
  font-weight: 800;
  text-shadow: 0 0 8px rgba(255, 255, 255, 0.8),
               0 0 12px rgba(255, 58, 58, 0.6);
  z-index: 2;
  letter-spacing: -1px;
  animation: vsFloat 3s infinite ease-in-out;
}

@keyframes vsPulse {
  0% { transform: scale(0.9); box-shadow: 0 0 15px rgba(255, 58, 58, 0.6); }
  100% { transform: scale(1.05); box-shadow: 0 0 25px rgba(255, 58, 58, 0.8); }
}

@keyframes vsFloat {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

/* 空白占位框（与VS容器大小一致） */
.empty-placeholder {
  /* 保持与VS组件完全相同的尺寸和间距 */
  width: 80px;
  height: 80px;
  margin-top: 28px; /* 匹配match-status-bar的高度，确保整体占位大小一致 */
  background-color: transparent; /* 透明背景 */
}

/* 响应式适配 */
@media (max-width: 768px) {
  .match-status-bar {
    gap: 6px;
    padding: 2px 6px;
  }
  
  .status-text, .players-value {
    font-size: 11px;
  }
  
  .players-label {
    font-size: 10px;
  }

  .empty-placeholder {
    margin-top: 24px;
  }
}

@media (max-width: 576px) {
  .vs-container, .empty-placeholder {
    width: 60px;
    height: 60px;
  }
  
  .vs-text {
    font-size: 22px;
  }

  .empty-placeholder {
    margin-top: 20px;
  }
}
</style>
    