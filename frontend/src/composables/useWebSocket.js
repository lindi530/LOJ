// useWebSocket.js
import { ref,  } from 'vue'
import api from '@/api'

let ws = null
let isConnected = false

const messageMap = ref({})
const followedUsers = ref([])
const selectedUserId = ref(null)

const messageCache = ref([])
let pendingMessages = ref([])
const submitCodeCallbacks = ref([])
const matchCallbacks = ref([])
const saberResultCallbacks = ref([])

function sameUserId(a, b) {
  return String(a) === String(b)
}

export function registersaberResultCallback(callback) { 
  saberResultCallbacks.value.push(callback)
  return () => {
    saberResultCallbacks.value = saberResultCallbacks.value.filter(cb => cb !== callback)
  }
}

export function registerSubmitCodeCallback(callback) {
  submitCodeCallbacks.value.push(callback)
  // 返回注销函数，避免组件卸载后仍执行
  return () => {
    submitCodeCallbacks.value = submitCodeCallbacks.value.filter(cb => cb !== callback)
  }
}

export function registerMatchCallback(callback) {
  matchCallbacks.value.push(callback)

  return () => {
    matchCallbacks.value = matchCallbacks.value.filter(cb => cb !== callback)
  }
}


export function initWebSocket(token) {
  if (ws && ws.readyState === WebSocket.OPEN) return

  const wsBaseUrl = process.env.VUE_APP_WS_BASE_URL;
  const wsUrl = `${wsBaseUrl}?token=${token}`;
  ws = new WebSocket(wsUrl);

  console.log("完成连接")

  getFollowedUsers()
  console.log("getFollowedUsers:", followedUsers.value)

  ws.onopen = () => {
    isConnected = true
    console.log('WebSocket连接成功')
    pendingMessages.value.forEach(msg => ws.send(JSON.stringify(msg)))
    pendingMessages.value = []
  }

  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data)
    console.log("onmessage", msg)
    switch (msg.type) { 
      case "match_success":
        handleMatchSuccess(msg)
        break;
      case "chat":
        handleChat(msg)
        break;
      case "edit_status":
        handleSubmitCode(msg)
        break;
      case "online_status":
        handleOnlineStatus(msg)
        break;
      case "saber_result":
        handleSaberResult(msg)
        break;
    }
  }
  // 链接关闭
  ws.onclose = () => {
    isConnected = false
    ws = null
    console.log('WebSocket连接关闭，尝试重连...')
    setTimeout(() => {
        if (localStorage.accessToken) {
          initWebSocket(localStorage.accessToken);
        }
      }, 3000);
  }

  // 连接错误
  ws.onerror = (error) => {
    console.error('WebSocket错误', error);
  };
}

export function closeWebSocket() {
  if (ws) {
    console.log("手动关闭 WebSocket 连接")
    ws.close() // 会触发 ws.onclose
    ws = null
  }

  isConnected = false
  messageMap.value = {}
  followedUsers.value = []
  selectedUserId.value = null
  pendingMessages.value = []
  messageCache.value = []
  submitCodeCallbacks.value = []
  matchCallbacks.value = []
  saberResultCallbacks.value = []
}

function handleSaberResult(msg) {
  saberResultCallbacks.value.forEach(callback => {
    console.log("saber result")
    callback(msg)
  })
}

function handleOnlineStatus(msg) {
  const user = followedUsers.value.find(user => sameUserId(user.user_id, msg.from))
  if (user) {
    user.online_state = msg.online_state
  }
}

function handleMatchSuccess(msg) { 
  matchCallbacks.value.forEach(callback => {
    callback(msg)
  })
}

function handleSubmitCode(msg) {
  submitCodeCallbacks.value.forEach(callback => {
    callback(msg)
  })
}

function handleChat(msg) {
  
  if (followedUsers.value.length === 0) {
    // 关注列表还没准备好，先缓存消息
      messageCache.value.push(msg)
      return
  }
  console.log("test: ", msg)
  processMessage(msg)
}

function processAllMessages() {
  messageCache.value.forEach(msg => processMessage(msg))
  messageCache.value = []  // 清空缓存
}
function processMessage(msg) {
  const from = msg.from
  if (!messageMap.value[from]) messageMap.value[from] = []
  messageMap.value[from].push(msg)

  console.log("processMessage: ", msg)
  const user = followedUsers.value.find(u => sameUserId(u.user_id, from))
  console.log("hasUnread", user)
  if (user) {
    user.hasUnread = true
  }
}

export function sendMessage(msg) {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    if (!ws) {
      initWebSocket(localStorage.accessToken)
    }

    // 缓存消息
    pendingMessages.value.push(msg)
    console.log('[WebSocket] not ready, message cached:', msg)
    return
  }

  console.log("Sending message:", JSON.stringify(msg))
  ws.send(JSON.stringify(msg))
}



export function useWebSocketContext() {
  return {
    initWebSocket,
    closeWebSocket,
    messageMap,
    followedUsers,
    selectedUserId,
    sendMessage,
    registerSubmitCodeCallback,
    registerMatchCallback,
    registersaberResultCallback,
  }
}

async function getFollowedUsers() {
  if (followedUsers.value.length != 0) return
  const resp = await api.getFollowUserList();
  if (resp.code === 0) {
    console.log("getFollowedUsers")
    console.log(resp.data.followUserList)
    followedUsers.value = resp.data.followUserList.map(u => ({ ...u, hasUnread: false }));

    processAllMessages()
  }
}
