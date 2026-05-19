<template>
  <n-drawer
    v-model:show="chatVisible"
    :default-width="1408"
    :placement="placement"
    resizable
  >
    <n-drawer-content title="聊天" style="height: 100%; overflow: auto;">
      <ChatAPP/>
    </n-drawer-content>
  </n-drawer>
  <SaberAPP v-model:visible="saberVisible" class="saber-fix"/>
  <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-gradient shadow-sm py-2"
    style="background-color: #493131;">

    <div class="container">
      <RouterLink class="navbar-brand fw-bold fs-4" to="/">My Blog</RouterLink>
      <button class="navbar-toggler" @click="isOpen = !isOpen" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>

      <div class="collapse navbar-collapse" :class="{ show: isOpen }">
        <ul class="navbar-nav me-auto mb-0">
          <li v-for="item in leftLinks" :key="item.to" class="nav-item">
            <button type="button" class="nav-link btn btn-link px-3" @click="handleLinkClick(item)">
              {{ item.label }}
            </button>
          </li>
        </ul>

        <ul class="navbar-nav ms-auto mb-0">
          <template v-if="isLogin">
            <li class="nav-item dropdown">
              <a
                class="nav-link dropdown-toggle user-menu"
                href="#"
                role="button"
                data-bs-toggle="dropdown"
                aria-expanded="false"
              >
                <img :src="userAvatar" alt="Avatar" class="rounded-circle me-2" width="28" height="28" />
                {{ userName }}
              </a>
              <ul class="dropdown-menu dropdown-menu-end shadow animated-dropdown">
                <li>
                  <RouterLink class="dropdown-item" :to="`/users/${userId}`">
                    <i class="bi bi-person me-2"></i> 用户信息
                  </RouterLink>
                </li>
                <li>
                  <RouterLink class="dropdown-item" to="/user/profile">
                    <i class="bi bi-gear me-2"></i> 设置
                  </RouterLink>
                </li>
                <li><hr class="dropdown-divider" /></li>
                <li>
                  <a class="dropdown-item text-danger" href="#" @click.prevent="logout">
                    <i class="bi bi-box-arrow-right me-2"></i> 登出
                  </a>
                </li>
              </ul>
            </li>
          </template>
          <template v-else>
            <li class="nav-item">
              <a href="#" class="nav-link px-3" @click.prevent="loginVisible = true">登录</a>
            </li>
            <li class="nav-item">
              <a href="#" class="nav-link px-3" @click.prevent="registerVisible = true">注册</a>
            </li>
          </template>
        </ul>
      </div>

      <LoginModal v-model:visible="loginVisible" @login-success="handleLoginSuccess" />
      <RegisterModal v-model:visible="registerVisible" />
    </div>
  </nav>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { closeWebSocket } from '@/composables/useWebSocket'
import LoginModal from './account/LoginModal.vue';
import RegisterModal from './account/RegisterModal.vue';
import ChatAPP from '@/views/UserChatView.vue'
import SaberAPP from './saber/Interface.vue'
import UserList from './list/UserList.vue';
import api from '@/api';

const router = useRouter();
const store = useStore();

const isOpen = ref(false);
const loginVisible = ref(false);
const registerVisible = ref(false);
const chatVisible = ref(false);
const placement = ref("right")
const saberVisible = ref(false)
const pendingRoute = ref(null);



const refreshToken = computed(() => store.getters['user/refreshToken'] || localStorage.refreshToken);
const accessToken = computed(() => store.getters['user/accessToken']);
const isLogin = computed(() => store.getters['user/isLogin']);
const userName = computed(() => store.getters['user/userName']);
const userId = computed(() => store.getters['user/userId']);
const userAvatar = computed(() => store.getters['user/userAvatar'] || '/default-avatar.png');
console.log("loginUserId: ", userId)
const leftLinks = computed(() => [
  { label: '首页', to: '/' },
  { label: '帖子列表', to: '/posts' },  
  { label: '用户信息', to: '/user-info' },
  { label: '用户列表', to: '/users/userList' },
  { label: '聊天', to: '/users/chat' },
  { label: '题库', to: '/problems' },
  { label: '对战', to: '/' },
  { label: '题目上传', to: '/upload/problem'},
  { label: '测试', to: '/test'}
]);

function handleLinkClick(item) {
  switch (item.label) { 
    case "用户信息":
      if (isLogin.value && userId.value) {
        router.push(`/users/${userId.value}`);
      } else {
        loginVisible.value = true;
        pendingRoute.value = `/users/${userId.value}`;
      }
      break
    case "聊天": 
      chatVisible.value = true;
      break
    case '对战':
      saberVisible.value = true;
      break
    default:
      router.push(item.to);
      break
  }
}

async function logout() {
  try {
    const resp = await api.logout({ refresh_token: refreshToken.value });
    
    if (resp.code === 0) {
      store.commit('user/SET_ACCESSTOKEN', '');
      store.commit('user/SET_REFRESHTOKEN', '');
      store.commit('user/SET_PROFILE', {});
      store.commit('user/LOGOUT');

      closeWebSocket();
    }
  } catch (err) {
    console.error("登出失败：", err);
  }
}

function handleLoginSuccess() {
  loginVisible.value = false;
  nextTick(() => {
    if (pendingRoute.value && userId.value) {
      router.push(`/users/${userId.value}`);
      pendingRoute.value = null;
    }
  });
}
</script>

<style scoped>
.user-menu {
  text-decoration: none !important;
  display: flex;
  align-items: center;
  font-weight: 500;
}

.animated-dropdown {
  animation: fadeInDown 0.2s ease-out;
}

.navbar-dark .nav-link,
.user-menu,
.dropdown-item {
  color: #ffffff !important;
}

.dropdown-menu {
  background-color: #ffffff;
  color: #333333;
}

.dropdown-item {
  color: #333333 !important;
}

.dropdown-item:hover {
  background-color: #f1f1f1;
  color: #000000 !important;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 关键修复：确保Saber组件的按钮不受父组件样式影响 */
:deep(.saber-fix) {
  position: relative;
  z-index: 9999 !important;
}

:deep(.saber-fix .content-text.btn) {
  pointer-events: auto !important;
  position: relative !important;
  z-index: 10000 !important;
}
</style>