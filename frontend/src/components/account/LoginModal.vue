<template>
  <n-modal :show="visible" @update:show="val => emit('update:visible', val)" preset="card" title="用户登录" @close="close" style="max-width: 400px;">
    <form @submit.prevent="login">
      <div class="mb-3">
        <label for="loginUsername" class="form-label">用户名</label>
        <input
          id="loginUsername"
          v-model="form.username"
          type="text"
          class="form-control"
          placeholder="请输入用户名"
          required
        />
      </div>
      <div class="mb-3">
        <label for="loginPassword" class="form-label">密码</label>
        <input
          id="loginPassword"
          v-model="form.password"
          type="password"
          class="form-control"
          placeholder="请输入密码"
          required
        />
      </div>
      <div class="d-flex align-items-center">
        <Button type="primary" html-type="submit">登录</Button>
        <transition name="fade">
          <div v-if="errorMessage" class="text-danger ms-3">
            <i class="bi bi-exclamation-circle-fill me-1"></i>
            {{ errorMessage }}
          </div>
        </transition>
      </div>
    </form>
  </n-modal>
</template>

<script setup>
import { ref, computed  } from 'vue'
import { useStore } from 'vuex'
import { initWebSocket } from '@/composables/useWebSocket'
import api from '@/api'

// props
const props = defineProps({
  visible: { type: Boolean, default: false }
})
const emit = defineEmits(['update:visible', 'login-success'])

// state
const form = ref({ username: '', password: '' })
const errorMessage = ref('')
const store = useStore()

// methods
const close = () => {
  emit('update:visible', false)
  form.value.username = ''
  form.value.password = ''
  errorMessage.value = ''
}

const login = async () => {
  errorMessage.value = ''
  try {
    const res = await api.login(
      { ...form.value },
      { headers: { 'Content-Type': 'application/json' } }
    )
    if (res.code) {
      errorMessage.value = res.message
      return
    }
    const { accessToken: at, refreshToken, user } = res.data
    console.log("loginData: ", res.data)
    // commit to store
    store.commit('user/SET_ACCESSTOKEN', at)
    store.commit('user/SET_REFRESHTOKEN', refreshToken)
    store.commit('user/SET_PROFILE', user)
    // init websocket after login

    initWebSocket(at)
    // emit success
    emit('login-success', user)
    close()
  } catch (err) {
    
    errorMessage.value = err.response?.data?.message || '登录失败，请重试'
  }
}
</script>

<style scoped>
.modal { background-color: rgba(0, 0, 0, 0.5); }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.text-danger { color: #dc3545; font-size: 0.875rem; }
</style>
