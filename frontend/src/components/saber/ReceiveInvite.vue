<template>
<div class="player-container">
  <div class="friend-invite-controls">
    <div class="room-input-group">
      <textarea 
        v-model="roomId" 
        placeholder="输入房间号加入对战"
        class="room-input"
        />
      <button 
        class="join-btn" 
        @click="handleJoin"
        :disabled="!roomId.trim() || isLoading"
      >
        加入房间
      </button>
    </div>
    <button 
      class="invite-btn" 
      @click="handleInvite"
      :disabled="isLoading"
    >
      {{ isLoading ? '处理中...' : '邀请好友' }}
    </button>
  </div>
</div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
  isLoading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['invite-friend', 'join-room']);

const roomId = ref('');

const handleInvite = () => {
  emit('invite-friend');
};

const handleJoin = () => {
  if (roomId.value.trim()) {
    emit('join-room', roomId.value.trim());
    roomId.value = '';
  }
};
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
.friend-invite-controls {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 80%;
  gap: 20px;
  padding: 20px;
}

.friend-invite-controls h3 {
  color: #e0e7ff;
  margin: 0;
  font-size: 20px;
}

.room-input-group {
  display: flex;
  width: 100%;
  gap: 10px;
}

.room-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid rgba(58, 134, 255, 0.5);
  border-radius: 6px;
  background-color: rgba(15, 52, 96, 0.6);
  color: #e0e7ff;
  font-size: 14px;
  outline: none;

}

.room-input::placeholder {
  color: #94a3b8;
}

.join-btn, .invite-btn {
  padding: 10px 16px;
  border-radius: 6px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.join-btn {
  background-color: rgba(58, 134, 255, 0.8);
  color: white;
}

.join-btn:hover:not(:disabled) {
  background-color: rgba(37, 99, 235, 0.9);
  transform: translateY(-2px);
}

.join-btn:disabled {
  background-color: rgba(58, 134, 255, 0.4);
  cursor: not-allowed;
}

.invite-btn {
  width: 100%;
  background-color: rgba(74, 222, 128, 0.8);
  color: #0f172a;
  font-size: 16px;
}

.invite-btn:hover:not(:disabled) {
  background-color: rgba(74, 222, 128, 1);
  transform: translateY(-2px);
}

.invite-btn:disabled {
  background-color: rgba(74, 222, 128, 0.4);
  cursor: not-allowed;
}
</style>
