<template>
  <div class="layout-container">
    <!-- 匹配成功特效提示 -->
    <MatchSuccessOverlay 
      :show="showMatchSuccess" 
      :is-fading-out="isFadingOut"
    />
    
    <!-- 房间号弹窗 -->
    <RoomNumberModal
      :show="showRoomModal"
      :room-number="roomNumber"
      :formatted-countdown="formattedRoomCountdown"
      :countdown-progress="roomCountdownProgress"
      @close="handleRoomModalClose"
      @copy-room-number="copyRoomNumber"
      @cancel-invitation="handleCancelInvitation"
    />
    
    <!-- 自定义过期提示弹窗 -->
    <div class="expire-alert" v-if="showExpireAlert">
      <div class="alert-backdrop" @click="showExpireAlert = false"></div>
      <div class="alert-content">
        <div class="alert-icon">{{ alertIcon }}</div>
        <h3 class="alert-title">{{ alertTitle }}</h3>
        <p class="alert-message">{{ alertMessage }}</p>
        <button class="alert-btn" @click="showExpireAlert = false">
          确定
        </button>
      </div>
    </div>
    
    <!-- 顶部标题区域 -->
    <TopBar :battle-type="battleType" />
    
    <!-- 中间玩家+VS区域 -->
    <div class="middle-area" ref="middleAreaRef">
      <!-- 左侧玩家区域 -->
      <PlayerCard 
        :player="leftPlayer" 
        :is-hovered="isLeftHovered" 
        :is-matching="isMatching"
        :box-style="playerBoxStyle"
        @mouse-enter="handleLeftEnter"
        @mouse-leave="handleLeftLeave"
      />
      
      <!-- 中间VS区域及匹配时间 -->
      <VSCard 
        :battle-type="battleType"
        :is-matching="isMatching"
        :match-time="formattedMatchTime"
        :found-players="foundPlayers"

      />
      
      <!-- 右侧玩家区域 -->
      <PlayerCard 
        v-if="checkBattleType()"
        :player="rightPlayer" 
        :is-hovered="isRightHovered" 
        :is-matching="isMatching"
        :box-style="playerBoxStyle"
        @mouse-enter="handleRightEnter"
        @mouse-leave="handleRightLeave"

      />
      <ReceiveInvite  v-else
        :problem-id="problemId"
        :is-loading="isMatching"
        @invite-friend="handleInviteFriend"
        @join-room="handleJoinRoom"
      />
    </div>
    
    <!-- 底部按钮区域 -->
    <BottomControls 
      :is-matching="isMatching"
      :battle-type="battleType"
      @match-or-cancel="handleMatchOrCancel"
      @back-to-menu="handleBackToMenu"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watchEffect, computed, watch } from 'vue';
import api from '@/api/index.js'
import { useStore } from 'vuex';
import { useWebSocketContext } from '@/composables/useWebSocket.js'

// 导入子组件
import MatchSuccessOverlay from './MatchSuccessOverlay.vue';
import TopBar from './TopBar.vue';
import PlayerCard from './PlayerCard.vue';
import VSCard from './VSCard.vue';
import BottomControls from './BottomControls.vue';
import RoomNumberModal from './RoomNumberModal.vue';
import ReceiveInvite from './ReceiveInvite.vue';


// 接收外部传入的对战类型
const props = defineProps({
  battleType: {
    type: String,
    required: true
  },
  postId: {
    type: [String, Number],
    required: true
  },
});

// 定义对外暴露的事件
const emit = defineEmits(['back-to-menu', 'update_match_success', 'to-battle-game', 'set-room-id']);

// 匹配状态管理
const isMatching = ref(false);
const matchTimeSeconds = ref(0);
const foundPlayers = ref(0);
let matchTimer = null;
let playerFoundTimer = null;

// 匹配成功特效相关变量
const showMatchSuccess = ref(false);
const isFadingOut = ref(false);

// 房间号弹窗相关变量
const showRoomModal = ref(false);
const roomNumber = ref('');
const roomCountdownSeconds = ref(180); // 倒计时时间，修改为3分钟
let roomCountdownTimer = null;
const inRoom = ref(false); // 新增：标记是否在房间内

// 新增：过期提示弹窗内容控制
const showExpireAlert = ref(false);
const cancelExpire = ref(false)
const autoAlertTitle = ref('已取消邀请')
const autoAlertMessage = ref('该房间已销毁，请重新创建房间');
const alertTitle = ref('房间已过期');
const alertMessage = ref('抱歉，房间已超过有效期，请重新创建房间');
const alertIcon = ref('⏰');


const checkBattleType = () => {
  return props.battleType === "天人之战";
}

// 格式化房间倒计时
const formattedRoomCountdown = computed(() => {
  const minutes = Math.floor(roomCountdownSeconds.value / 60);
  const seconds = roomCountdownSeconds.value % 60;
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
});

// 房间倒计时进度
const roomCountdownProgress = computed(() => {
  return (roomCountdownSeconds.value / 180) * 100;
});


const store = useStore();
// 左右玩家数据
const leftPlayer = ref({
  "avatar": computed(() => store.getters['user/userAvatar']),
  "name": computed(() => store.getters['user/userName']),
  "level": "",
  "rating": 0,
  "total_matches": 0,
  "wins": 0,
});

onMounted(async() => { 
  const resp = await api.getUserSaberStats(props.postId)
  if (resp.code === 0) {
    let data = resp.data
    
    leftPlayer.value.rating = data.rating
    leftPlayer.value.level = data.level
    leftPlayer.value.total_matches = data.total_matches
    leftPlayer.value.wins = data.wins
  } else {
    // 处理错误
  }
})

const rightPlayer = ref({
  "avatar": "",
  "name": "",
  "level": "",
  "rating": 0,
  "total_matches": 0,
  "wins": 0,
});

const { registerMatchCallback } = useWebSocketContext()


const unregister = registerMatchCallback((msg) => {
  console.log("匹配成功")
  handleMatchSuccess(msg)
})

const stopAll = () => { 
  // 关闭所有弹窗
    showMatchSuccess.value = false;
    showRoomModal.value = false;
    showExpireAlert.value = false;
    
    // 停止所有计时器
    stopMatchProcess();
    stopRoomCountdown();
    
    // 重置匹配状态
    isMatching.value = false;
}

const handleMatchSuccess = (msg) => {
  emit('set-room-id', roomNumber.value)
  stopAll()
  
  rightPlayer.value.avatar = msg.opponent.avatar
  rightPlayer.value.name = msg.opponent.user_name
  rightPlayer.value.level = msg.opponent.level
  rightPlayer.value.rating = msg.opponent.rating
  rightPlayer.value.total_matches = msg.opponent.total_matches
  rightPlayer.value.wins = msg.opponent.wins

  foundPlayers.value = 2

  // 触发匹配成功特效
  triggerMatchSuccessEffect();
  // 延迟发送匹配成功事件，等待特效展示完成
  setTimeout(() => {
    emit('to-battle-game', msg.room_id, msg.problem_id)
  }, 5000);
}

// 触发匹配成功特效
const triggerMatchSuccessEffect = () => {
  showMatchSuccess.value = true;
  isFadingOut.value = false;
  stopMatchProcess();
  
  // 2.5秒后开始淡出效果
  setTimeout(() => {
    isFadingOut.value = true;
  }, 4000);
  
  // 3秒后隐藏整个提示
  setTimeout(() => {
    showMatchSuccess.value = false;
  }, 5000);
}


// 格式化匹配时间为 MM:SS
const formattedMatchTime = computed(() => {
  const minutes = Math.floor(matchTimeSeconds.value / 60);
  const seconds = matchTimeSeconds.value % 60;
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
});

// 返回菜单事件
const handleBackToMenu = () => {
  stopMatchProcess();
  stopRoomCountdown();
  emit('back-to-menu');
  console.log('返回菜单按钮被点击');
};

// 显示提示弹窗
const showAlert = (title, message, icon = '⏰') => {
  alertTitle.value = title;
  alertMessage.value = message;
  alertIcon.value = icon;
  showExpireAlert.value = true;
};

// 开始匹配事件
const handleMatchOrCancel = async () => {
  if (isMatching.value) { 
    stopMatchProcess();
    if (props.battleType === '好友对战') {
      stopRoomCountdown();
      showRoomModal.value = false;
    }
    return;
  }
  
  foundPlayers.value = 1;
  isMatching.value = true;

  try {
    // 根据对战类型调用不同的API
    
      // 天人之战的API
      const resp = await api.match();
      if (resp.code === 0) {
        console.log("发送匹配请求成功！");

        // 启动匹配计时器
        matchTimeSeconds.value = 0;
        matchTimer = setInterval(() => {
          matchTimeSeconds.value++;
        }, 1000);
      } else {
        console.log("发送匹配请求失败：", resp.message);
        showAlert('匹配失败', resp.message || '无法发送匹配请求，请稍后再试');
        // 失败时重置状态
        isMatching.value = false;
      }
    
  } catch (error) {
    console.error("请求后端数据失败：", error);
    showAlert('网络错误', '与服务器连接失败，请检查网络');
    // 异常时重置状态
    isMatching.value = false;
  }
};

// 房间倒计时控制
const startRoomCountdown = () => {
  // 重置倒计时
  cancelExpire.value = false
  roomCountdownSeconds.value = 180; // 3分钟
  
  // 清除可能存在的旧计时器
  if (roomCountdownTimer) {
    clearInterval(roomCountdownTimer);
  }
  
  // 启动新计时器
  roomCountdownTimer = setInterval(() => {
    roomCountdownSeconds.value--;
    
    // 倒计时结束
    if (roomCountdownSeconds.value <= 0) {
      handleRoomExpired();
    }
  }, 1000);
};

// 停止房间倒计时
const stopRoomCountdown = () => {
  if (roomCountdownTimer) {
    clearInterval(roomCountdownTimer);
    roomCountdownTimer = null;
  }
};

// 房间过期处理
const handleRoomExpired = async () => {
  // 避免重复调用
  if (!roomNumber.value || roomCountdownSeconds.value <= 0) return;
  
  stopRoomCountdown();
  showRoomModal.value = false;
  isMatching.value = false;
  inRoom.value = false;
  
  try {
    const resp = await api.expireRoom(roomNumber.value)
    console.log("expire Room: ", resp)
    if (resp.code === 0) {
      console.log('通知房间过期成功');
      if (cancelExpire.value) {
        cancelExpire.value = false
        showAlert(autoAlertTitle.value, autoAlertMessage.value);
      } else {
        showAlert('房间已过期', '抱歉，房间已超过有效期，请重新创建房间');
      }
      
    }
  } catch (err) {
    console.error('通知房间过期失败:', err);
    showAlert('房间已过期', '抱歉，房间已超过有效期，请重新创建房间');
  }
};

// 房间弹窗关闭处理
const handleRoomModalClose = () => {
  showRoomModal.value = false;
  // 不停止倒计时，因为用户可能只是暂时关闭窗口
};

// 复制房间号
const copyRoomNumber = () => {
  navigator.clipboard.writeText(roomNumber.value).then(() => {
    // 显示复制成功的提示
    showAlert('复制成功', '房间号已成功复制到剪贴板', '✓');
    // 3秒后关闭
  }).catch(err => {
    console.error('无法复制文本: ', err);
    // 显示复制失败提示
    showAlert('复制失败', '请手动复制房间号', '⚠️');
  });
};

// 取消邀请处理 - 修复核心：只调用一次过期处理并确保计时器已停止
const handleCancelInvitation = () => {
  showRoomModal.value = false;
  cancelExpire.value = true
  // 先强制停止所有计时器，避免残留回调
  stopRoomCountdown();
  stopMatchProcess();
  // 只调用一次房间过期处理
  handleRoomExpired();
};

const handleReceiveInvite = () => { 
  showRoomModal.value = false;
  stopRoomCountdown();
  stopMatchProcess();
}

// 停止匹配流程（清除计时器）
const stopMatchProcess = () => {
  isMatching.value = false;
  if (matchTimer) {
    clearInterval(matchTimer);
    matchTimer = null;
  }
  if (playerFoundTimer) {
    clearTimeout(playerFoundTimer);
    playerFoundTimer = null;
  }
};

// 玩家区域悬停状态管理
const middleAreaRef = ref(null);
const playerBoxStyle = ref({ width: '0px', height: '0px' });
const isLeftHovered = ref(false);
const isRightHovered = ref(false);

// 左侧玩家悬停进入
const handleLeftEnter = () => {
  if (!isMatching.value) {
    isLeftHovered.value = true;
    console.log('左侧玩家区域悬停进入');
  }
};

// 左侧玩家悬停离开
const handleLeftLeave = () => {
  isLeftHovered.value = false;
  console.log('左侧玩家区域悬停离开');
};

// 右侧玩家悬停进入
const handleRightEnter = () => {
  if (!isMatching.value) {
    isRightHovered.value = true;
    console.log('右侧玩家区域悬停进入');
  }
};

// 右侧玩家悬停离开
const handleRightLeave = () => {
  isRightHovered.value = false;
  console.log('右侧玩家区域悬停离开');
};

// 邀请好友（创建房间）
const handleInviteFriend = async () => {
  if (isMatching.value) return;
  
  try {
    isMatching.value = true;
    const resp = await api.inviteFriend();
    if (resp.code === 0) {
      console.log("创建邀请成功！");
      roomNumber.value = resp.data.room_id;
      emit('set-room-id', roomNumber.value)
      showRoomModal.value = true;
      startRoomCountdown();
      inRoom.value = true;
    } else {
      showAlert('创建失败', resp.message || '无法创建邀请，请稍后再试', '⚠️');
      isMatching.value = false;
    }
  } catch (error) {
    console.error("创建邀请失败：", error);
    showAlert('创建失败', '网络错误，无法创建邀请', '⚠️');
    isMatching.value = false;
  }
};

// 加入房间
const handleJoinRoom = async (roomId) => {
  if (isMatching.value || !roomId) return;
  console.log("roomId: ", roomId)
  try {
    isMatching.value = true;
    const resp = await api.joinRoom(roomId);
    if (resp.code === 0) {
      console.log("加入房间成功！");
      roomNumber.value = roomId;
      emit('set-room-id', roomId)
      startRoomCountdown();
      inRoom.value = true;
    } else {
      showAlert('加入失败', resp.message || '房间不存在或已过期', '⚠️');
      isMatching.value = false;
    }
  } catch (error) {
    console.error("加入房间失败：", error);
    showAlert('加入失败', '网络错误，无法加入房间', '⚠️');
    isMatching.value = false;
  }
};

// 计算玩家框尺寸（自适应中间区域大小）
const calculatePlayerBoxSize = () => {
  if (!middleAreaRef.value) return;
  const middleArea = middleAreaRef.value;
  const middleWidth = middleArea.clientWidth;
  const middleHeight = middleArea.clientHeight;
  const vsAreaTotalWidth = 100; // VS区域宽度
  
  // 计算单个玩家区域最大可用宽度
  const maxAvailableWidthPerPlayer = (middleWidth - vsAreaTotalWidth) / 2;
  // 按 3:4 宽高比计算高度
  const heightBasedOnWidth = maxAvailableWidthPerPlayer * 4 / 3;
  // 中间区域最大可用高度（90%）
  const maxAvailableHeight = middleHeight * 0.9;
  
  let finalWidth, finalHeight;
  // 优先按宽度计算，若高度超出则按高度反算宽度
  if (heightBasedOnWidth <= maxAvailableHeight) {
    finalWidth = maxAvailableWidthPerPlayer;
    finalHeight = heightBasedOnWidth;
  } else {
    finalHeight = maxAvailableHeight;
    finalWidth = finalHeight * 3 / 4;
  }
  
  // 限制最大尺寸（避免过大）
  finalWidth = Math.min(270, finalWidth);
  finalHeight = Math.min(360, finalHeight);
  playerBoxStyle.value = { width: `${finalWidth}px`, height: `${finalHeight}px` };
};

// 组件挂载时初始化尺寸，并监听窗口 resize
onMounted(() => {
  calculatePlayerBoxSize();
  const handleResize = () => calculatePlayerBoxSize();
  window.addEventListener('resize', handleResize);
  
  // 组件卸载时清理
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize);
    stopMatchProcess(); // 确保计时器被清除
    stopRoomCountdown(); // 确保房间倒计时被清除
    unregister(); // 取消WebSocket回调注册
  });
  
  // 监听中间区域变化，重新计算尺寸
  watchEffect(() => middleAreaRef.value && calculatePlayerBoxSize());
});
</script>

<style scoped>
/* 全局容器样式 */
.layout-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
  background: transparent;
  pointer-events: auto;
  position: relative;
}

/* 中间玩家+VS区域 */
.middle-area {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  box-sizing: border-box;
  position: relative;
  z-index: 1;
}

/* 房间过期提示样式 */
.expire-alert {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1100; /* 确保在其他弹窗上方 */
  pointer-events: none;
}

.alert-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(2px);
  pointer-events: auto;
  animation: fadeIn 0.3s ease;
}

.alert-content {
  position: relative;
  width: 90%;
  max-width: 350px;
  background-color: #1e293b;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(58, 134, 255, 0.3);
  animation: popIn 0.3s ease-out;
  pointer-events: auto;
}

.alert-icon {
  font-size: 48px;
  margin-bottom: 16px;
  color: #fde68a;
}

.alert-title {
  margin: 0 0 12px 0;
  color: #e0e7ff;
  font-size: 20px;
  font-weight: 600;
}

.alert-message {
  margin: 0 0 24px 0;
  color: #94a3b8;
  font-size: 14px;
  line-height: 1.5;
}

.alert-btn {
  background-color: rgba(58, 134, 255, 0.8);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px 24px;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
  font-size: 16px;
}

.alert-btn:hover {
  background-color: rgba(37, 99, 235, 0.9);
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(58, 134, 255, 0.3);
}

/* 动画效果 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
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
</style>
    