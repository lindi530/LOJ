<template>
  <div class="bottom-area">
    <button 
      v-if="battleType === '天人之战'"
      class="action-btn" 
      :class="isMatching ? 'matching-btn cancel-match-btn' : (battleType === '天人之战' ? 'friend-btn' : 'match-btn')" 
      @click="$emit('match-or-cancel')"
    >
      <i :class="isMatching ? 'icon-loading' : (battleType === '好友对战' ? 'icon-friend' : 'icon-game')"></i>
      <span>
        {{ isMatching 
          ? '取消匹配' 
          : (battleType === '好友对战' ? '邀请好友' : '开始匹配') 
        }}
      </span>
    </button>
    <button class="action-btn back-btn" @click="$emit('back-to-menu')">
      <i class="icon-back"></i>
      <span>返回菜单</span>
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  isMatching: {
    type: Boolean,
    default: false
  },
  battleType: {
    type: String,
    required: true
  }
});

// 定义事件
const emit = defineEmits(['match-or-cancel', 'back-to-menu']);
</script>

<style scoped>
/* 底部按钮区域 */
.bottom-area {
  height: 70px;
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 24px;
  padding: 0 20px;
  position: relative;
  z-index: 10;
  background: transparent;
}

/* 通用按钮样式 */
.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 28px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s ease;
  position: relative;
  z-index: 11;
  overflow: hidden;
}

/* 开始匹配按钮 */
.match-btn {
  background-color: rgba(58, 134, 255, 0.8);
  color: white;
  box-shadow: 0 4px 15px rgba(58, 134, 255, 0.2);
}

.match-btn:hover {
  background-color: rgba(37, 99, 235, 0.9);
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(58, 134, 255, 0.4);
}

/* 邀请好友按钮 */
.friend-btn {
  background-color: rgba(74, 222, 128, 0.8);
  color: #064e3b;
  box-shadow: 0 4px 15px rgba(74, 222, 128, 0.2);
}

.friend-btn:hover {
  background-color: rgba(34, 197, 94, 0.9);
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(74, 222, 128, 0.4);
}

/* 匹配中按钮（不可点击） */
.matching-btn {
  background-color: rgba(148, 163, 184, 0.8);
  color: #1e293b;
  box-shadow: 0 4px 15px rgba(148, 163, 184, 0.2);
  cursor: default;
}

.matching-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(148, 163, 184, 0.3);
}

/* 取消匹配按钮样式 */
.action-btn.cancel-match-btn {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
}
.action-btn.cancel-match-btn:hover {
  background-color: #ff7875;
}

/* 按钮图标 */
.icon-game::before { content: "🎮"; }
.icon-friend::before { content: "👥"; }
.icon-loading::before { content: "🔄"; animation: spin 1.5s linear infinite; }
.icon-back::before { content: "↩"; }

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 按钮点击反馈 */
.match-btn:active {
  transform: translateY(-1px);
  box-shadow: 0 3px 10px rgba(58, 134, 255, 0.3);
}

/* 返回菜单按钮 */
.back-btn {
  background-color: rgba(30, 58, 138, 0.8);
  color: #e0e7ff;
  box-shadow: 0 4px 15px rgba(30, 58, 138, 0.2);
}

.back-btn:hover {
  background-color: rgba(30, 64, 175, 0.9);
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(30, 58, 138, 0.4);
  color: white;
}

.back-btn:active {
  transform: translateY(-1px);
  box-shadow: 0 3px 10px rgba(30, 58, 138, 0.3);
}

/* 响应式适配 */
@media (max-width: 576px) {
  .bottom-area { 
    height: 60px; 
    gap: 12px; 
    padding: 0 10px;
  }
  
  .action-btn {
    padding: 8px 18px;
    font-size: 14px;
  }
}
</style>
