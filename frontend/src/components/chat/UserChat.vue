<template>

        <div class="chat-layout">
          
          <div class="left-sidebar">
              <n-card size="large" bg-color="#f0f7ff"  
                :content-style="{ backgroundColor: '#f0f7ff' }"
              >
              <FollowList
                :followed-users="followedUsers"
                :selected-user-id="selectedUserId"
                @select="handleSelectUser"
              />
              </n-card>
          </div>
        
          <div class="right-chat">
            <ChatWindow
              :messages="messages"
              :user-id="userId"
              :user-avatar="userAvatar"
              :receiver-id="selectedUserId"
              :receiver-name="selectedUserName"
              :receiver-avatar="selectedUserAvatar"
              v-model:new-message="newMessage"
              @send="handleSendMessage"
              ref="chatRef"
            />
          </div>

        </div>
</template>

<script setup>
import { ref, computed, watch, onBeforeUnmount, onMounted } from 'vue'
import { useStore } from 'vuex'
import FollowList from '../chat/FollowList.vue'
import ChatWindow from '../chat/ChatWindow.vue'
import api from '@/api'
import { useWebSocketContext } from '@/composables/useWebSocket'


const store = useStore()
const { messageMap, followedUsers, selectedUserId, sendMessage } = useWebSocketContext()

const userId = computed(() => String(store.getters['user/userId']))
const userAvatar = computed(() => store.getters['user/userAvatar'] || '/default-avatar.png')



const newMessage = ref('')
const chatRef = ref(null)

function sameUserId(a, b) {
  return String(a) === String(b)
}

const messages = computed(() => messageMap.value[selectedUserId.value] || [])

const selectedUserAvatar = computed(() => {
  const u = followedUsers.value.find(u => sameUserId(u.user_id, selectedUserId.value))
  return u?.avatar || ''
})

const selectedUserName = computed(() => {
  const u = followedUsers.value.find(u => sameUserId(u.user_id, selectedUserId.value))
  return u?.user_name || ''
})

async function handleSelectUser(id) {
  selectedUserId.value = id
  const user = followedUsers.value.find(u => sameUserId(u.user_id, id))
  if (user) user.hasUnread = false

  if (!messageMap.value[id]) {
    try { 
      const res = await api.getMessageByTargetId(id)
      if (res.code === 0) {
        messageMap.value[id] = res.data
      }
    } catch (error) { 
      console.error('发送请求失败:', error);
    }
  }
}

function handleSendMessage() {
  
  if (!newMessage.value || !selectedUserId.value) return
  const msg = {
    type: 'chat',
    from: userId.value,
    to: selectedUserId.value,
    content: newMessage.value
  }
  messageMap.value[selectedUserId.value] = messageMap.value[selectedUserId.value] || []
  messageMap.value[selectedUserId.value].push(msg)
  console.log("msg: ", msg)
  sendMessage(msg)
  newMessage.value = ''
}

</script>

<style scoped>

.chat-layout {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.left-sidebar {
  width: 15vw;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.right-chat {
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;    /* ✅ 关键！允许子元素收缩 */
}
</style>
