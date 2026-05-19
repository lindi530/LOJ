<template>
  <div class="player-container">
    <div class="effect-wrapper" :class="{ 'hover-active': isHovered && !isMatching }">
      <div class="player-box" 
           :style="boxStyle"
           @mouseenter="$emit('mouse-enter')" 
           @mouseleave="$emit('mouse-leave')">
        <!-- 玩家信息 -->
        <div class="player-info-container">
          <!-- 玩家头像 -->
          <div class="avatar-wrapper">
            <img :src="player.avatar" :alt="player.name" class="player-avatar">
            <div class="player-level-badge">{{ player.level }}</div>
          </div>
          
          <!-- 玩家基本信息 -->
          <div class="player-basic-info">
            <h3 class="player-name">{{ player.name }}</h3>
            <div class="player-rating">
              <i class="rating-icon">🏆</i>
              <span class="rating-value">{{ player.rating }}</span>
            </div>
          </div>
          
          <!-- 玩家统计数据 -->
          <div class="player-stats">
            <div class="stat-item">
              <span class="stat-label">场数</span>
              <span class="stat-value">{{ player.total_matches }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">胜率</span>
              <span class="stat-value">{{ player.wins }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  player: {
    type: Object,
    required: true,
    default: () => ({
      avatar: "",
      name: "",
      level: "",
      rating: 0,
      total_matches: 0,
      wins: 0
    })
  },
  isHovered: {
    type: Boolean,
    default: false
  },
  isMatching: {
    type: Boolean,
    default: false
  },
  boxStyle: {
    type: Object,
    required: true
  }
});

// 定义事件
const emit = defineEmits(['mouse-enter', 'mouse-leave']);
</script>

<style scoped>
.player-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0 10px;
  position: relative;
  z-index: 2;
  background: transparent;
}

/* 玩家区域样式 */
.effect-wrapper {
  padding: 4px;
  border-radius: 10px;
  transition: all 0.3s ease;
  position: relative;
  z-index: 3;
  background: transparent;
}

.player-box {
  background-color: rgba(15, 52, 96, 0.6);
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s ease;
  border: 1px solid rgba(58, 134, 255, 0.3);
  position: relative;
  z-index: 4;
  cursor: default;
}

/* 玩家信息容器 */
.player-info-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  padding: 8%;
  box-sizing: border-box;
  color: white;
  position: relative;
}

/* 头像区域 */
.avatar-wrapper {
  position: relative;
  width: 45%;
  aspect-ratio: 1/1;
  margin-bottom: 5%;
}

.player-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
  border: 2px solid rgba(58, 134, 255, 0.5);
  box-shadow: 0 0 10px rgba(58, 134, 255, 0.3);
  transition: all 0.3s ease;
}

/* 悬停时头像放大+增强阴影 */
.effect-wrapper.hover-active .player-avatar {
  transform: scale(1.05);
  box-shadow: 0 0 15px rgba(58, 134, 255, 0.5);
}

/* 等级徽章 */
.player-level-badge {
  position: absolute;
  bottom: -5%;
  right: -5%;
  background: linear-gradient(135deg, #3a86ff, #8338ec);
  color: white;
  font-size: clamp(0.7rem, 4vw, 0.9rem);
  font-weight: bold;
  padding: 2% 6%;
  border-radius: 50px;
  border: 2px solid rgba(15, 52, 96, 0.8);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  white-space: nowrap;
}

/* 玩家基本信息 */
.player-basic-info {
  text-align: center;
  width: 100%;
  margin-bottom: 5%;
}

.player-name {
  font-size: clamp(1rem, 6vw, 1.3rem);
  font-weight: 600;
  margin: 0 0 3% 0;
  color: #e0f2fe;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.player-rating {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #fde68a;
  font-size: clamp(0.8rem, 4vw, 1rem);
}

.rating-icon {
  font-size: 1.1em;
}

.rating-value {
  font-weight: 500;
}

/* 玩家统计数据 */
.player-stats {
  display: flex;
  justify-content: space-around;
  width: 100%;
  padding-top: 5%;
  border-top: 1px solid rgba(58, 134, 255, 0.2);
}

.stat-item {
  text-align: center;
  flex: 1;
}

.stat-label {
  display: block;
  font-size: clamp(0.7rem, 3vw, 0.8rem);
  color: #94a3b8;
  margin-bottom: 2px;
}

.stat-value {
  font-size: clamp(0.8rem, 4vw, 1rem);
  font-weight: 600;
  color: #e0e7ff;
}

/* 玩家区域悬停效果 */
.effect-wrapper::before {
  content: "";
  position: absolute;
  inset: 2px;
  border-radius: 8px;
  background: linear-gradient(45deg, rgba(58, 134, 255, 0.2), rgba(131, 56, 236, 0.15));
  opacity: 0.3;
  filter: blur(3px);
  z-index: 2;
  transition: all 0.3s ease;
}

.effect-wrapper.hover-active {
  transform: translateY(-2px) scale(1.01);
}

.effect-wrapper.hover-active::before {
  background: linear-gradient(45deg, #3a86ff, #8338ec);
  opacity: 0.6;
  filter: blur(6px);
  box-shadow: 0 0 10px rgba(58, 134, 255, 0.5);
}

.effect-wrapper.hover-active .player-box {
  border-color: rgba(58, 134, 255, 0.8);
  background-color: rgba(18, 69, 128, 0.7);
}

/* 响应式适配 */
@media (max-width: 768px) {
  .player-stats {
    flex-direction: column;
    gap: 8px;
  }
  
  .stat-item {
    display: flex;
    justify-content: space-between;
  }
}

@media (max-width: 576px) {
  .player-container { padding: 0 5px; }
  
  .avatar-wrapper {
    width: 50%;
  }
}
</style>
