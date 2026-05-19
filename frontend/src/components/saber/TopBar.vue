<template>
  <div class="top-area">
    <div class="top-content">
      <div class="page-title">{{ battleType }}</div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  battleType: {
    type: String,
    required: true
  }
});
</script>

<style scoped>
/* 顶部标题区域 */
.top-area {
  height: 80px;
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0 20px;
  background: transparent;
  z-index: 10;
}

.top-content {
  width: 100%;
  max-width: 1200px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: transparent;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(90deg, #3a86ff, #8338ec);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  letter-spacing: 1px;
  transition: font-size 0.3s ease;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .page-title {
    font-size: 26px;
  }
}

@media (max-width: 576px) {
  .top-area { 
    padding: 8px 10px;
    height: 60px;
  }
  
  .page-title {
    font-size: 22px;
  }
}
</style>
