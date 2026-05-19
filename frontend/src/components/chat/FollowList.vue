<template>
  <n-infinite-scroll style="height: 100vh" :distance="10" @load="handleLoad">
    <div
      v-for="user in followedUsers"
      :key="user.user_id"
      :class="['user-item-wrapper', isSelectedUser(user.user_id) ? 'selected' : '']"
      @click="$emit('select', user.user_id)"
    >
      <div class="user-item">
        <div class="user-info">
          <img :src="user.avatar" class="avatar" alt="avatar" />
          <span class="user-name">{{ user.user_name }}</span>
        </div>
        <div class="user-meta">
          <div v-if="user.hasUnread" class="unread-dot"></div>
          <div :class="['user-status', user.online_state ? 'online' : 'offline']">
            {{ user.online_state ? '在线' : '离线' }}
            <i :class="user.online_state ? 'bi bi-check-circle-fill' : 'bi bi-x-circle'"></i>
          </div>
        </div>
      </div>
    </div>
  </n-infinite-scroll>
</template>

<script setup>
import { ref } from "vue"

const props = defineProps({
  followedUsers: Array,
  selectedUserId: [Number, String],
})

const count = ref(6)

function isSelectedUser(userId) {
  return String(userId) === String(props.selectedUserId)
}

const handleLoad = () => {
  count.value += 1
}
</script>

<style scoped>
.follow-list {
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px 0;
}

.user-item-wrapper {
  cursor: pointer;
  padding: 8px 12px;
  transition: background-color 0.2s ease;
}

.user-item-wrapper:hover {
  background-color: #f5f5f5;
}

.user-item-wrapper.selected {
  background-color: #9ed048;
  color: #ffffff;
  border-radius: 8px; /* 圆角大小，可根据需要调整数值 */
}

.user-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.user-info {
  display: flex;
  align-items: center;
  min-width: 0;
  overflow: hidden;
  flex: 1;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 8px;
  flex-shrink: 0;
}

.user-name {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  flex: 1;
}

.user-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.unread-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: red;
}

.user-status {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  gap: 4px;
}

.user-status.online {
  color: #0eb840;
}

.user-status.offline {
  color: #999999;
}
</style>
