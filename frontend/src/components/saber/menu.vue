<template>
  <div class="button-container" :style="{ marginLeft: leftOffset + 'px', marginTop: topOffset + 'px' }">
    <button 
      type="button" 
      class="battle-btn heaven-battle"
      :style="{
        minWidth: buttonWidth + 'px',
        minHeight: buttonHeight + 'px',
        padding: buttonPadding + 'px',
        fontSize: buttonFontSize * 1.5 + 'px'
      }"
      @click="$emit('heaven-battle')"
    >
      天人之战
    </button>
    <button 
      type="button" 
      class="battle-btn friend-battle"
      :style="{
        minWidth: buttonWidth + 'px',
        minHeight: buttonHeight + 'px',
        padding: buttonPadding + 'px',
        fontSize: buttonFontSize * 1.5 + 'px'
      }"
      @click="$emit('friend-battle')"
    >
      好友对战
    </button>
  </div>
</template>


<script setup>
// 接收从父组件传递的样式参数
const props = defineProps({
  buttonWidth: {
    type: Number,
    required: true
  },
  buttonHeight: {
    type: Number,
    required: true
  },
  buttonPadding: {
    type: Number,
    required: true
  },
  buttonFontSize: {
    type: Number,
    required: true,
    validator: (value) => value >= 20 && value <= 36
  },
  leftOffset: {
    type: Number,
    default: 30
  },
  topOffset: {
    type: Number,
    default: 80
  }
});


// 定义要发射的事件
const emit = defineEmits(['heaven-battle', 'friend-battle']);
</script>


<style scoped>
/* 按钮容器 - 上下排列并左对齐 */
.button-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  position: relative;
  z-index: 101;
  align-items: flex-start;
}

/* 基础按钮样式 - 确保可点击 */
.battle-btn {
  background: none;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  padding: 24px 20px;
  position: relative;
  overflow: visible;
  pointer-events: auto !important;
  display: inline-block;
}

/* 更换为高清晰度字体，去除过多光效 */
.heaven-battle, .friend-battle {
  /* 更换为更清晰的无衬线字体组合 */
  font-family: 'Roboto', 'Segoe UI', 'Arial', sans-serif;
  font-weight: 700; /* 适中粗细，保证清晰 */
  letter-spacing: 0.8px;
  text-transform: uppercase;
  /* 简化文字阴影，仅保留轻微效果避免模糊 */
  text-shadow: 0 1px 2px rgba(0,0,0,0.3);
  line-height: 1.2;
  /* 字体平滑处理 */
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 天人之战按钮样式 - 更纯的颜色减少模糊 */
.heaven-battle {
  color: #64b5f6;
}

/* 好友对战按钮样式 - 更纯的颜色减少模糊 */
.friend-battle {
  color: #f5f5f5;
}

/* 悬停效果 */
.battle-btn:hover {
  transform: translateX(6px) scale(1.03);
  filter: brightness(1.2);
}

/* 点击效果 */
.battle-btn:active {
  transform: translateX(3px) scale(0.98);
  filter: brightness(0.9);
}

/* 按钮左侧装饰线 - 保留但简化 */
.battle-btn::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  opacity: 0.8;
  transition: all 0.2s ease;
}

.heaven-battle::before {
  background-color: #64b5f6;
  box-shadow: 0 0 4px #64b5f6; /* 简化阴影 */
}

.friend-battle::before {
  background-color: #f5f5f5;
  box-shadow: 0 0 4px #f5f5f5; /* 简化阴影 */
}

.battle-btn:hover::before {
  width: 6px;
  opacity: 1;
}

/* 确保没有父元素阻止点击 */
:deep(*) {
  pointer-events: inherit !important;
}
</style>
