<template>
  <div class="room-number-modal" v-if="show">
    <div class="modal-overlay" @click="$emit('close')"></div>
    <div class="modal-content">
      <div class="modal-header">
        <h3>好友对战房间</h3>
      </div>
      
      <div class="modal-body">
        <p class="room-info">您的房间号为：</p>
        
        <!-- 优化长房间号显示（仅文字滚动） -->
        <div class="room-number-wrapper">
          <div class="room-number-outer">
            <div class="room-number">{{ roomNumber }}</div>
          </div>
        </div>
        
        <!-- 倒计时显示区域 -->
        <div class="countdown-section">
          <p class="countdown-text">房间将在 {{ formattedCountdown }} 后自动关闭</p>
          <div class="countdown-bar">
            <div 
              class="countdown-progress" 
              :style="{ width: countdownProgress + '%' }"
              :class="{ 'warning': countdownProgress < 30 }"
            ></div>
          </div>
        </div>
        
        <p class="invite-info">请将房间号发送给好友，等待好友加入...</p>
      </div>
      
      <!-- 按钮区域（复制和取消邀请在同一行） -->
      <div class="modal-footer">
        <button class="copy-btn" @click="$emit('copy-room-number')">
          <span>复制房间号</span>
        </button>
        <button class="cancel-btn" @click="$emit('cancel-invitation')">
          取消邀请
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  roomNumber: {
    type: String,
    default: ''
  },
  formattedCountdown: {
    type: String,
    default: '03:00'
  },
  countdownProgress: {
    type: Number,
    default: 100
  }
});

const emit = defineEmits([
  'close', 
  'copy-room-number',
  'cancel-invitation'
]);
</script>

<style scoped>
.room-number-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(3px);
}

.modal-content {
  position: relative;
  width: 90%;
  max-width: 500px;
  background-color: #1e293b;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(58, 134, 255, 0.3);
  animation: popIn 0.3s ease-out;
}

@keyframes popIn {
  from {
    transform: scale(0.9);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(58, 134, 255, 0.2);
}

.modal-header h3 {
  margin: 0;
  color: #e0e7ff;
  font-size: 18px;
  font-weight: 600;
}

.modal-body {
  text-align: center;
  margin-bottom: 20px;
}

.room-info {
  color: #94a3b8;
  margin: 8px 0 16px;
  font-size: 14px;
}

/* 房间号容器（仅文字滚动，背景固定） */
.room-number-wrapper {
  margin: 12px 0 20px;
  padding: 16px 8px;
  background-color: rgba(15, 52, 96, 0.6);
  border-radius: 8px;
  border: 1px dashed rgba(74, 222, 128, 0.5);
}

.room-number-outer {
  overflow-x: auto; /* 仅文字滚动 */
  -webkit-overflow-scrolling: touch;
  padding-bottom: 4px; /* 为滚动条预留空间 */
}

.room-number {
  font-size: 28px;
  font-weight: bold;
  letter-spacing: 2px;
  color: #4ade80;
  font-family: 'Courier New', monospace;
  text-shadow: 0 0 8px rgba(74, 222, 128, 0.3);
  white-space: nowrap; /* 防止换行 */
  display: inline-block; /* 确保宽度适应内容 */
}

.countdown-section {
  margin: 20px 0;
  padding: 10px;
  background-color: rgba(15, 52, 96, 0.3);
  border-radius: 6px;
}

.countdown-text {
  color: #e0e7ff;
  font-size: 14px;
  margin: 0 0 8px 0;
}

.countdown-bar {
  height: 6px;
  background-color: rgba(15, 52, 96, 0.5);
  border-radius: 3px;
  overflow: hidden;
}

.countdown-progress {
  height: 100%;
  background: linear-gradient(90deg, #4ade80, #3a86ff);
  transition: width 1s linear;
}

.countdown-progress.warning {
  background: linear-gradient(90deg, #fde68a, #ef4444);
  animation: pulse 1s infinite alternate;
}

@keyframes pulse {
  from { opacity: 0.8; }
  to { opacity: 1; }
}

.invite-info {
  color: #94a3b8;
  margin: 16px 0;
  font-size: 14px;
}

/* 底部按钮区域（两个按钮在同一行） */
.modal-footer {
  display: flex;
  justify-content: center;
  gap: 16px; /* 按钮之间的间距 */
  margin-top: 12px;
}

.copy-btn {
  background-color: rgba(58, 134, 255, 0.8);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 9px 20px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
  font-weight: 500;
}

.copy-btn:hover {
  background-color: rgba(37, 99, 235, 0.9);
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(58, 134, 255, 0.3);
}

.copy-icon {
  font-size: 16px;
}

.cancel-btn {
  background-color: rgba(233, 69, 96, 0.2);
  color: #e94560;
  border: 1px solid rgba(233, 69, 96, 0.4);
  border-radius: 6px;
  padding: 9px 24px;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
}

.cancel-btn:hover {
  background-color: rgba(233, 69, 96, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(233, 69, 96, 0.2);
}

/* 滚动条样式优化 */
.room-number-outer::-webkit-scrollbar {
  height: 4px;
}

.room-number-outer::-webkit-scrollbar-track {
  background: rgba(15, 52, 96, 0.3);
  border-radius: 4px;
}

.room-number-outer::-webkit-scrollbar-thumb {
  background: rgba(58, 134, 255, 0.5);
  border-radius: 4px;
}

/* 响应式调整 */
@media (max-width: 576px) {
  .modal-content {
    padding: 18px;
    max-width: 90%;
  }
  
  .room-number {
    font-size: 24px;
    letter-spacing: 1px;
  }
  
  .modal-footer {
    flex-direction: column;
    gap: 10px;
  }
  
  .copy-btn, .cancel-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
    