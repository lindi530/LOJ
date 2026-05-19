<template>
  <div 
    class="match-success-overlay" 
    v-if="show"
    :class="{ 'fade-out': isFadingOut }"
  >
    <div class="success-content">
      <div class="success-icon">✦</div>
      <div class="success-text">匹配成功！</div>
      <div class="success-subtext">即将进入对战...</div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  isFadingOut: {
    type: Boolean,
    default: false
  }
});
</script>

<style scoped>
/* 匹配成功特效样式 */
.match-success-overlay {
  position: absolute;
  top: 0px; /* 顶部标题栏高度 */
  bottom: 0px; /* 底部按钮栏高度 */
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
  opacity: 1;
  transition: opacity 0.5s ease-out;
  animation: pulseBackground 4s infinite alternate;

  margin: 10px;
}

.match-success-overlay.fade-out {
  opacity: 0;
}

.success-content {
  text-align: center;
  transform: scale(0.8);
  animation: popIn 0.5s forwards, float 4s ease-in-out;
}

.success-icon {
  font-size: 60px;
  color: #4ade80;
  margin-bottom: 20px;
  text-shadow: 0 0 15px rgba(74, 222, 128, 0.8);
  animation: spin 5s linear infinite;
}

.success-text {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 10px;
  text-shadow: 0 0 10px rgba(58, 134, 255, 0.8);
  
  /* 渐变背景设置 */
  background: linear-gradient(90deg, 
    rgba(58, 134, 255, 0) 0%,  /* 开始透明 */
    #3a86ff 25%,              /* 蓝色 */
    #8338ec 50%,              /* 紫色 */
    #4ade80 75%,              /* 绿色 */
    rgba(74, 222, 128, 0) 100% /* 结束透明 */
  );
  background-size: 200% auto;
  
  /* 文字裁剪 */
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  
  /* 应用动画 */
  animation: textGradientFlow 5s linear infinite;
}

.success-subtext {
  font-size: 20px;
  color: #e0e7ff;
  text-shadow: 0 0 5px rgba(58, 134, 255, 0.5);
}

/* 匹配成功特效动画 */
@keyframes popIn {
  from { transform: scale(0.5); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

@keyframes spin {
  from { transform: rotate(360deg); }
  to { transform: rotate(0deg); }
}

@keyframes pulseBackground {
  from { background-color: rgba(223,244,255, 0.5); }
  to { background-color: rgba(223,244,255, 0.8); }
}

@keyframes textGradientFlow {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 100% 50%;
  }
}

/* 响应式适配 */
@media (max-width: 768px) {
  .success-text {
    font-size: 36px;
  }
  
  .success-icon {
    font-size: 48px;
  }
  
  .success-subtext {
    font-size: 16px;
  }
}

@media (max-width: 576px) {
  .success-text {
    font-size: 28px;
  }
  
  .success-icon {
    font-size: 40px;
  }
}
</style>
