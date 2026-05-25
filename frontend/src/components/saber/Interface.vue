<template>
  <AppModal
    :visible="showDialog"
    title="Saber"
    :close-on-backdrop="false"
    :draggable="true"
    :panel-style="{
      width: dialogSize.width + 'px',
      height: dialogSize.height + 'px',
      left: dialogLeft + 'px',
      top: dialogTop + 'px',
      position: 'absolute',
      '--modal-header-padding': headerPadding + 'px',
      '--modal-title-size': titleFontSize + 'px',
      '--modal-title-offset': titleMarginLeft + 'px',
      '--modal-close-size': closeBtnSize + 'px'
    }"
    @update:visible="closeDialog"
    @header-mousedown="startDrag"
  >
    <div
      class="saber-content"
      :style="{
        padding: bodyPadding + 'px',
        gap: buttonGap + 'px',
        backgroundImage: `url(${bgImage})`
      }"
    >
      <Menu
        v-if="currentView === 'menu'"
        :button-width="buttonWidth"
        :button-height="buttonHeight"
        :button-padding="buttonPadding"
        :button-font-size="buttonFontSize"
        @heaven-battle="handleHeavenBattle"
        @friend-battle="handleFriendBattle"
      />

      <Battle
        v-if="currentView === 'battle'"
        :battle-type="currentBattleType"
        :base-scale="scaleRatio"
        :battle-state="battleState"
        @set-room-id="handleSetRoomId"
        @back-to-menu="handleBackToMenu"
        @to-battle-game="handleToBattleGame"
        @update-match-success="updateMatchSuccess"
      />

      <BattleGame
        v-if="currentView === 'match_success'"
        :room-id="roomId"
        :problem-id="problemId"
        @back-to-menu="handleBackToMenu"
      />
    </div>

    <template #panel-extra>
      <div class="resize-handle" @mousedown="startResize"></div>
    </template>
  </AppModal>
</template>

<script setup>
import menuBg from '@/assets/background.jpg'; 
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
import AppModal from '@/components/modal/AppModal.vue';
import Menu from './menu.vue';
import Battle from './BattleMatch.vue'
import BattleGame from './BattleGame.vue';

// 接收显示状态props
const props = defineProps({
  visible: { 
    type: Boolean, 
    required: true,
    default: false
  }
})

// battle
const roomId = ref("")
const problemId = ref(null)
const battleState = ref(false)

const updateBattleState = ((val) => { 
  console.log("updateBattleState")
  battleState.value = val;
})

// 组件内部状态
const showDialog = ref(false);
const currentView = ref('menu'); // 视图状态：menu/battle
const currentBattleType = ref(''); // 记录当前对战类型
const emit = defineEmits(['update:visible'])

// 同步显示状态
watch(() => props.visible, (val) => {
  if (val) {
    showDialog.value = true;
    nextTick(() => {
      centerDialog();
    });
  } else {
    showDialog.value = false;
    currentView.value = 'menu'; // 隐藏时重置视图
  }
}, { immediate: true })

// 弹窗位置与尺寸配置
const dialogPositionRatio = ref({ x: 0.5, y: 0.5 });
const dialogLeft = ref(0);
const dialogTop = ref(0);
const baseRatio = 16 / 9;
const baseWidth = 1024;
const baseHeight = 576;
const minWidth = 320;
const minHeight = 260;
const minScaleRatio = 0.65;
const dialogSize = ref({ width: baseWidth, height: baseHeight });

// 元素尺寸配置
const elementSizes = ref({
  headerPadding: 6,
  titleFontSize: 16,
  titleMarginLeft: 16,
  closeBtnSize: 24,
  bodyPadding: 16,
  buttonWidth: 180,
  buttonHeight: 60,
  buttonPadding: 12,
  buttonFontSize: 20,
  buttonGap: 24
});

// 缩放后的元素尺寸
const headerPadding = ref(elementSizes.value.headerPadding);
const titleFontSize = ref(elementSizes.value.titleFontSize);
const titleMarginLeft = ref(elementSizes.value.titleMarginLeft)
const closeBtnSize = ref(elementSizes.value.closeBtnSize);
const bodyPadding = ref(elementSizes.value.bodyPadding);
const buttonWidth = ref(elementSizes.value.buttonWidth);
const buttonHeight = ref(elementSizes.value.buttonHeight);
const buttonPadding = ref(elementSizes.value.buttonPadding);
const buttonFontSize = ref(elementSizes.value.buttonFontSize);
const buttonGap = ref(elementSizes.value.buttonGap);

const updateElementSizes = (scaleRatio) => {
  const responsiveScale = Math.max(scaleRatio, minScaleRatio);

  headerPadding.value = elementSizes.value.headerPadding * responsiveScale;
  titleFontSize.value = Math.max(13, elementSizes.value.titleFontSize * responsiveScale);
  closeBtnSize.value = Math.max(24, elementSizes.value.closeBtnSize * responsiveScale);
  bodyPadding.value = elementSizes.value.bodyPadding * responsiveScale;
  buttonWidth.value = elementSizes.value.buttonWidth * responsiveScale;
  buttonHeight.value = elementSizes.value.buttonHeight * responsiveScale;
  buttonPadding.value = elementSizes.value.buttonPadding * responsiveScale;
  buttonFontSize.value = elementSizes.value.buttonFontSize * responsiveScale;
  titleMarginLeft.value = elementSizes.value.titleMarginLeft * responsiveScale;
  buttonGap.value = elementSizes.value.buttonGap * responsiveScale;
};

// 拖拽相关变量
const isDragging = ref(false);
const startX = ref(0);
const startY = ref(0);

// 拉伸相关变量
const isResizing = ref(false);
const startWidth = ref(0);
const startHeight = ref(0);

// 关闭弹窗
const closeDialog = () => {
  showDialog.value = false;
  emit('update:visible', false);
};

const bgImage = ref(menuBg);

const handleBackToMenu = () => { 
  battleState.value = false
  currentView.value = "menu";
  bgImage.value = menuBg;
}

const handleToBattleGame = (room_id, problem_id) => { 

  roomId.value = room_id
  problemId.value = problem_id
  bgImage.value = menuBg;

  currentView.value = "match_success";
}

// 按钮点击事件 - 切换到对战视图
const handleHeavenBattle = () => {
  currentBattleType.value = '天人之战';
  currentView.value = 'battle';
  bgImage.value = menuBg;
};

const handleFriendBattle = () => {
  currentBattleType.value = '好友对战';
  currentView.value = 'battle';
  bgImage.value = menuBg;
};

const handleSetRoomId = (val) => { 
  console.log("set RoomID", val)
  roomId.value = val
}

// 计算弹窗尺寸
const calculateSizes = () => {
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;
  
  const availableWidth = viewportWidth * 0.9;
  const availableHeight = viewportHeight * 0.9;
  
  let calcWidth, calcHeight;
  const heightByWidth = availableWidth / baseRatio;
  
  if (heightByWidth <= availableHeight) {
    calcWidth = availableWidth;
    calcHeight = heightByWidth;
  } else {
    calcHeight = availableHeight;
    calcWidth = availableHeight * baseRatio;
  }
  
  calcWidth = Math.max(minWidth, Math.min(baseWidth, calcWidth));
  calcHeight = Math.max(minHeight, Math.min(baseHeight, calcHeight));
  updateElementSizes(calcWidth / baseWidth);
  
  return { width: calcWidth, height: calcHeight };
};

// 计算弹窗位置
const calculatePositionFromRatio = () => {
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;
  const availableWidth = viewportWidth - dialogSize.value.width;
  const availableHeight = viewportHeight - dialogSize.value.height;
  
  dialogLeft.value = dialogPositionRatio.value.x * availableWidth;
  dialogTop.value = dialogPositionRatio.value.y * availableHeight;
};

// 居中弹窗
const centerDialog = () => {
  dialogPositionRatio.value = { x: 0.5, y: 0.5 };
  dialogSize.value = calculateSizes();
  calculatePositionFromRatio();
};

// 拖拽事件处理
const startDrag = (e) => {
  // 如果正在拉伸，则不启动拖拽
  if (isResizing.value) return;
  
  if (e.button !== 0) return;
  isDragging.value = true;
  startX.value = e.clientX;
  startY.value = e.clientY;
  document.body.style.cursor = 'move';
  e.stopPropagation();
  e.preventDefault();
};

const handleDragMove = (e) => {
  // 如果正在拉伸，则不处理拖拽
  if (isResizing.value) return;
  
  if (!isDragging.value) return;
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;
  const availableWidth = viewportWidth - dialogSize.value.width;
  const availableHeight = viewportHeight - dialogSize.value.height;
  
  dialogLeft.value = Math.max(0, Math.min(dialogLeft.value + (e.clientX - startX.value), availableWidth));
  dialogTop.value = Math.max(0, Math.min(dialogTop.value + (e.clientY - startY.value), availableHeight));
  
  dialogPositionRatio.value.x = availableWidth > 0 ? dialogLeft.value / availableWidth : 0.5;
  dialogPositionRatio.value.y = availableHeight > 0 ? dialogTop.value / availableHeight : 0.5;
  
  startX.value = e.clientX;
  startY.value = e.clientY;
};

const endDrag = () => {
  if (isDragging.value) {
    isDragging.value = false;
    document.body.style.cursor = '';
  }
};

// 拉伸事件处理
const startResize = (e) => {
  // 如果正在拖拽，则不启动拉伸
  if (isDragging.value) return;
  
  if (e.button !== 0) return;
  isResizing.value = true;
  startX.value = e.clientX;
  startY.value = e.clientY;
  startWidth.value = dialogSize.value.width;
  startHeight.value = dialogSize.value.height;
  document.body.style.cursor = 'se-resize';
  e.stopPropagation();
  e.preventDefault();
};

const handleResize = (e) => {
  if (!isResizing.value) return;
  
  // 计算新尺寸，确保不小于最小尺寸
  const newWidth = Math.max(minWidth, startWidth.value + (e.clientX - startX.value));
  const newHeight = Math.max(minHeight, startHeight.value + (e.clientY - startY.value));
  
  // 更新弹窗尺寸
  dialogSize.value.width = newWidth;
  dialogSize.value.height = newHeight;
  
  updateElementSizes(newWidth / baseWidth);
  
};

const endResize = () => {
  if (isResizing.value) {
    isResizing.value = false;
    document.body.style.cursor = '';
  }
};

const scaleRatio = ref(1);

// 窗口大小变化处理
const handleWindowResize = () => {
  if (showDialog.value && !isDragging.value && !isResizing.value) {
    dialogSize.value = calculateSizes();
    calculatePositionFromRatio();
  }
};

// 生命周期管理
onMounted(() => {
  document.addEventListener('mousemove', handleDragMove);
  document.addEventListener('mousemove', handleResize);  // 添加拉伸移动监听
  document.addEventListener('mouseup', endDrag);
  document.addEventListener('mouseup', endResize);  // 添加拉伸结束监听
  window.addEventListener('resize', handleWindowResize);
});

onUnmounted(() => {
  document.removeEventListener('mousemove', handleDragMove);
  document.removeEventListener('mousemove', handleResize);  // 移除拉伸移动监听
  document.removeEventListener('mouseup', endDrag);
  document.removeEventListener('mouseup', endResize);  // 移除拉伸结束监听
  window.removeEventListener('resize', handleWindowResize);
});
</script>

<style scoped>
.saber-content {
  position: relative;
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;
  box-sizing: border-box;
  background-position: center;
  background-size: cover;
}

.resize-handle {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 15px;
  height: 15px;
  background-color: #ccc;
  cursor: se-resize;
  border-top-left-radius: 4px;
  border: 1px solid #999;
  border-right: none;
  border-bottom: none;
  z-index: 10;
}

.resize-handle:hover {
  background-color: #aaa;
}
</style>
