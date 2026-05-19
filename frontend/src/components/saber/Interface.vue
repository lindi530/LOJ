<template>
  <div v-if="showDialog" class="saber-app-root">
    <!-- 将teleport内部的元素包裹在一个容器中，统一管理堆叠上下文 -->
    <teleport to="body">
      <div 
          :class="['modal-backdrop fade', showDialog ? 'show' : '']"
          :style="{ zIndex: 10000 }" 
        >
      <div class="modal-container">
        <!-- 弹窗内容 - 层级高于遮罩 -->
        <div class="modal-content-wrapper">
          <div 
            ref="dialogRef"
            :style="{
              width: dialogSize.width + 'px',
              height: dialogSize.height + 'px',
              left: dialogLeft + 'px',
              top: dialogTop + 'px',
              transform: 'none !important',
              overflow: 'visible'
            }"
            class="modal-dialog position-absolute bg-white rounded shadow-lg"
          >
            <!-- 弹窗头部 -->
            <div 
              ref="headerRef"
              class="modal-header custom-header cursor-move d-flex justify-content-between align-items-center"
              :style="{
                padding: headerPadding + 'px',
                minHeight: 'auto'
              }"
              @mousedown="startDrag"
            >
              <h5 class="modal-title text-dark mb-0" 
                  :style="[
                    {
                      fontSize: titleFontSize + 'px',
                      lineHeight: '1',
                      padding: '0'
                    },
                    { marginLeft: titleMarginLeft + 'px !important' }
                  ]">
                Saber
              </h5>
              <button 
                type="button" 
                class="btn-close btn-close-dark " 
                aria-label="Close"
                :style="{
                  width: closeBtnSize + 'px',
                  height: closeBtnSize + 'px'
                }"
                @click="closeDialog"
              ></button>
            </div>
            
            <!-- 弹窗内容区 -->
            <div class="modal-body content-with-bg d-flex align-items justify-content button" 
                :style="{
                  padding: bodyPadding + 'px',
                  position: 'relative',
                  zIndex: 1,
                  gap: buttonGap + 'px',
                  height: `calc(100% - ${headerActualHeight}px)`,
                  overflow: 'hidden',
                  backgroundImage: `url(${bgImage})`,
                  backgroundSize: 'cover',
                  backgroundPosition: 'center'
                }">
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
            
            <!-- 拉伸手柄 -->
            <div 
              class="resize-handle"
              @mousedown="startResize"
            ></div>
          </div>
        </div>
      </div>
       </div>
    </teleport>
  </div>
</template>

<script setup>
import BattleBg from '@/assets/1.png'; 
import menuBg from '@/assets/background.jpg'; 
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
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

// 新增：用于获取头部实际高度的变量
const headerRef = ref(null);
const headerActualHeight = ref(0);

// 同步显示状态
watch(() => props.visible, (val) => {
  if (val) {
    showDialog.value = true;
    nextTick(() => {
      centerDialog();
      updateHeaderHeight(); // 显示时计算头部高度
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
const dialogRef = ref(null);
const baseRatio = 16 / 9;
const baseWidth = 1024;
const baseHeight = 576;
const minWidth = 320;
const minHeight = 180;
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

// 新增：更新头部实际高度的方法
const updateHeaderHeight = () => {
  if (headerRef.value) {
    // 获取头部元素的实际高度（包括padding和border）
    headerActualHeight.value = headerRef.value.offsetHeight;
  }
};

const bgImage = ref(menuBg);

const handleBackToMenu = () => { 
  battleState.value = false
  currentView.value = "menu";
  bgImage.value = menuBg;
  nextTick(updateHeaderHeight); // 视图切换后重新计算高度
}

const handleToBattleGame = (room_id, problem_id) => { 

  roomId.value = room_id
  problemId.value = problem_id
  bgImage.value = menuBg;

  nextTick(() => {
    currentView.value = "match_success";
    updateHeaderHeight(); // 视图切换后重新计算高度
  });
}

// 按钮点击事件 - 切换到对战视图
const handleHeavenBattle = () => {
  currentBattleType.value = '天人之战';
  currentView.value = 'battle';
  bgImage.value = menuBg;
  nextTick(updateHeaderHeight); // 视图切换后重新计算高度
};

const handleFriendBattle = () => {
  currentBattleType.value = '好友对战';
  currentView.value = 'battle';
  bgImage.value = menuBg;
  nextTick(updateHeaderHeight); // 视图切换后重新计算高度
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
  const scaleRatio = calcWidth / baseWidth;
  
  // 按比例更新所有元素尺寸
  headerPadding.value = elementSizes.value.headerPadding * scaleRatio;
  titleFontSize.value = elementSizes.value.titleFontSize * scaleRatio;
  closeBtnSize.value = elementSizes.value.closeBtnSize * scaleRatio;
  bodyPadding.value = elementSizes.value.bodyPadding * scaleRatio;
  buttonWidth.value = elementSizes.value.buttonWidth * scaleRatio;
  buttonHeight.value = elementSizes.value.buttonHeight * scaleRatio;
  buttonPadding.value = elementSizes.value.buttonPadding * scaleRatio;
  buttonFontSize.value = elementSizes.value.buttonFontSize * scaleRatio;
  titleMarginLeft.value = elementSizes.value.titleMarginLeft * scaleRatio;
  buttonGap.value = elementSizes.value.buttonGap * scaleRatio;
  
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
  nextTick(updateHeaderHeight); // 居中后重新计算头部高度
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
  
  // 重新计算元素缩放比例
  const scaleRatio = newWidth / baseWidth;
  headerPadding.value = elementSizes.value.headerPadding * scaleRatio;
  titleFontSize.value = elementSizes.value.titleFontSize * scaleRatio;
  closeBtnSize.value = elementSizes.value.closeBtnSize * scaleRatio;
  bodyPadding.value = elementSizes.value.bodyPadding * scaleRatio;
  buttonWidth.value = elementSizes.value.buttonWidth * scaleRatio;
  buttonHeight.value = elementSizes.value.buttonHeight * scaleRatio;
  buttonPadding.value = elementSizes.value.buttonPadding * scaleRatio;
  buttonFontSize.value = elementSizes.value.buttonFontSize * scaleRatio;
  titleMarginLeft.value = elementSizes.value.titleMarginLeft * scaleRatio;
  buttonGap.value = elementSizes.value.buttonGap * scaleRatio;
  
  // 更新头部高度
  updateHeaderHeight();
};

const endResize = () => {
  if (isResizing.value) {
    isResizing.value = false;
    document.body.style.cursor = '';
  }
};

const scaleRatio = ref(1);

// 窗口大小变化处理：增加更新头部高度
const handleWindowResize = () => {
  if (showDialog.value && !isDragging.value && !isResizing.value) {
    dialogSize.value = calculateSizes();
    calculatePositionFromRatio();
    updateHeaderHeight(); // 窗口变化后重新计算头部高度
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

.modal-container {
  position: fixed;  /* 关键：创建独立的堆叠上下文 */
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  /* z-index: 9999; */
}

/* 保持原有遮罩样式不变 */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.15s linear;
  /* z-index: 1; */
}

.modal-backdrop.show {
  opacity: 1;
  pointer-events: auto;
}

/* 样式保持不变 */
.saber-app-root {
  position: absolute;
  width: 0;
  height: 0;
  overflow: hidden;
  /* z-index: -1; */
}

* {
  -webkit-text-size-adjust: none !important;
  text-size-adjust: none !important;
  box-sizing: border-box !important;
}

.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  /* z-index: 10; */
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.15s linear;
}

.modal-content-wrapper {
  position: relative;
  /* z-index: 2; 高于遮罩层 */
  width: 100%;
  height: 100%;
  pointer-events: none; /* 允许点击穿透到遮罩 */
}

.modal-backdrop.show {
  opacity: 1;
  pointer-events: auto;
}

.modal-dialog {
  margin: 0;
  transition: transform 0.3s ease-out;
  overflow: hidden;
  transform: none !important;

  pointer-events: auto;
  /* z-index: 3;  确保在内容包装器内层级最高 */
}

.modal-header {
  user-select: none;
  pointer-events: auto;
  margin: 0 !important;
  padding: 0 !important;
}

.custom-header {
  background-color: #F5F5F5;
}

.content-with-bg {
  height: 100%;
  margin: 0 !important;
  padding: 0 !important;
  position: relative;
  /* z-index: 1; */
}

.modal-body {
  position: relative;
  /* z-index: 1; 弹窗内部内容层级 */
}

/* 修复Edit组件可能的层级问题 */
.page-container {
  position: relative;
  /* z-index: 1; 确保在弹窗内容范围内 */
}

.button-container {
  flex-wrap: wrap;
}

:deep(.modal-title) {
  margin-left: var(--title-margin) !important;
}

:deep(.modal-header), :deep(.modal-body) {
  margin: 0 !important;
  padding: 0 !important;
}

/* 拉伸手柄样式 */
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

/* 鼠标悬停时的样式 */
.resize-handle:hover {
  background-color: #aaa;
}
</style>