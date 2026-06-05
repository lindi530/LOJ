<template>
  <n-drawer
    v-model:show="chatVisible"
    :default-width="1408"
    :placement="placement"
    resizable
  >
    <n-drawer-content title="聊天" style="height: 100%; overflow: auto;">
      <ChatAPP />
    </n-drawer-content>
  </n-drawer>

  <SaberAPP v-model:visible="saberVisible" class="saber-fix" />

  <nav
    class="navbar navbar-expand-md navbar-light fixed-top loj-navbar"
  >
    <div class="container nav-container">
      <RouterLink class="navbar-brand loj-brand" to="/">LOJ</RouterLink>

      <button
        class="navbar-toggler"
        @click="isOpen = !isOpen"
        aria-label="Toggle navigation"
        :aria-expanded="isOpen"
      >
        <span class="navbar-toggler-icon"></span>
      </button>

      <div class="collapse navbar-collapse nav-collapse" :class="{ show: isOpen }">
        <ul class="navbar-nav nav-links me-auto mb-0">
          <li v-for="item in leftLinks" :key="item.label" class="nav-item">
            <button type="button" class="nav-link nav-link-button" @click="handleLinkClick(item)">
              {{ item.label }}
            </button>
          </li>
        </ul>

        <ul class="navbar-nav nav-actions ms-auto mb-0">
          <li class="nav-item">
            <button
              type="button"
              class="nav-icon-btn"
              aria-label="打开聊天"
              title="聊天"
              @click="openChat"
            >
              <i class="bi bi-chat-dots"></i>
            </button>
          </li>

          <template v-if="isLogin">
            <li class="nav-item">
              <n-dropdown
                trigger="click"
                placement="bottom-end"
                :options="userMenuOptions"
                :menu-props="userDropdownMenuProps"
                :render-icon="renderUserMenuIcon"
                :render-label="renderUserMenuLabel"
                @select="handleUserMenuSelect"
              >
              <button type="button" class="nav-link user-menu" aria-haspopup="menu">
                <img :src="userAvatar" alt="Avatar" class="rounded-circle user-avatar">
                <span class="user-name">{{ userName }}</span>
                <i class="bi bi-chevron-down user-menu-chevron"></i>
              </button>
              </n-dropdown>
            </li>
          </template>

          <template v-else>
            <li class="nav-item">
              <a href="#" class="nav-link auth-link" @click.prevent="loginVisible = true">登录</a>
            </li>
            <li class="nav-item">
              <a href="#" class="nav-link auth-link" @click.prevent="registerVisible = true">注册</a>
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
import { ref, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { closeWebSocket } from '@/composables/useWebSocket'
import LoginModal from './account/LoginModal.vue'
import RegisterModal from './account/RegisterModal.vue'
import ChatAPP from '@/views/UserChatView.vue'
import SaberAPP from './saber/Interface.vue'
import api from '@/api'

const router = useRouter()
const store = useStore()

const isOpen = ref(false)
const loginVisible = ref(false)
const registerVisible = ref(false)
const chatVisible = ref(false)
const placement = ref('right')
const saberVisible = ref(false)

const refreshToken = computed(() => store.getters['user/refreshToken'] || localStorage.refreshToken)
const isLogin = computed(() => store.getters['user/isLogin'])
const userName = computed(() => store.getters['user/userName'])
const userId = computed(() => store.getters['user/userId'])
const userAvatar = computed(() => store.getters['user/userAvatar'] || '/default-avatar.svg')
const userMenuOptions = computed(() => [
  { label: '用户信息', key: 'profile', icon: 'bi bi-person' },
  { label: '设置', key: 'settings', icon: 'bi bi-gear' },
  { type: 'divider', key: 'divider' },
  { label: '登出', key: 'logout', icon: 'bi bi-box-arrow-right', danger: true }
])
const userDropdownMenuProps = { class: 'loj-user-dropdown-menu' }

const leftLinks = computed(() => [
  { label: '帖子', to: '/' },
  { label: '用户列表', to: '/users/userList' },
  { label: '题库', to: '/problems' },
  { label: '竞赛', to: '/competition' },
  { label: '对战', to: '/' },
  { label: '题目上传', to: '/upload/problem' },
  { label: '测试', to: '/test' }
])

function handleLinkClick(item) {
  isOpen.value = false

  switch (item.label) {
    case '对战':
      saberVisible.value = true
      break
    default:
      router.push(item.to)
      break
  }
}

function openChat() {
  isOpen.value = false
  chatVisible.value = true
}

function handleUserMenuSelect(key) {
  switch (key) {
    case 'profile':
      router.push(`/users/${userId.value}`)
      break
    case 'settings':
      router.push('/user/profile')
      break
    case 'logout':
      logout()
      break
    default:
      break
  }
}

function renderUserMenuIcon(option) {
  if (!option.icon) return null

  return h('i', {
    class: [
      'loj-user-dropdown-icon',
      option.icon,
      option.danger ? 'is-danger' : ''
    ]
  })
}

function renderUserMenuLabel(option) {
  return h('span', {
    class: [
      'loj-user-dropdown-label',
      option.danger ? 'is-danger' : ''
    ]
  }, option.label)
}

async function logout() {
  try {
    const resp = await api.logout({ refresh_token: refreshToken.value })

    if (resp.code === 0) {
      store.commit('user/SET_ACCESSTOKEN', '')
      store.commit('user/SET_REFRESHTOKEN', '')
      store.commit('user/SET_PROFILE', {})
      store.commit('user/LOGOUT')
      closeWebSocket()
    }
  } catch (err) {
    console.error('登出失败：', err)
  }
}

function handleLoginSuccess() {
  loginVisible.value = false
}
</script>

<style scoped>
.loj-navbar {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.94), rgba(240, 247, 255, 0.9)),
    linear-gradient(90deg, rgba(37, 99, 235, 0.12), rgba(20, 184, 166, 0.12));
  border-bottom: 1px solid rgba(15, 23, 42, 0.1);
  box-shadow: 0 14px 34px rgba(15, 23, 42, 0.12);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  padding: 0.65rem 0 !important;
}

.nav-container,
.nav-collapse,
.nav-links,
.nav-actions {
  align-items: center;
}

.nav-container {
  gap: 0.75rem;
}

.loj-brand {
  color: #0f172a !important;
  font-size: 1.45rem;
  font-weight: 800;
  letter-spacing: 0;
  line-height: 1;
  text-decoration: none;
}

.loj-brand::before {
  content: "";
  display: inline-block;
  width: 0.7rem;
  height: 0.7rem;
  margin-right: 0.45rem;
  border-radius: 999px;
  background: linear-gradient(135deg, #2563eb, #14b8a6);
  box-shadow: 0 0 0 0.25rem rgba(37, 99, 235, 0.12);
  vertical-align: 0.05rem;
}

.navbar-toggler {
  border: 1px solid rgba(15, 23, 42, 0.12);
  border-radius: 0.75rem;
  background-color: rgba(255, 255, 255, 0.72);
  padding: 0.45rem 0.6rem;
}

.navbar-toggler:focus {
  box-shadow: 0 0 0 0.2rem rgba(37, 99, 235, 0.14);
}

.nav-collapse {
  gap: 0.75rem;
}

.nav-links {
  gap: 0.15rem;
}

.nav-actions {
  gap: 0.35rem;
}

.nav-link-button,
.auth-link,
.user-menu {
  display: flex;
  align-items: center;
  border-radius: 0.65rem;
  color: #334155 !important;
  font-size: 0.94rem;
  font-weight: 600;
  line-height: 1;
  text-decoration: none !important;
  transition: color 0.18s ease, background-color 0.18s ease, transform 0.18s ease;
}

.nav-link-button {
  border: 0;
  background: transparent;
  padding: 0.65rem 0.78rem;
}

.auth-link {
  padding: 0.65rem 0.78rem !important;
}

.nav-link-button:hover,
.auth-link:hover,
.user-menu:hover {
  background-color: rgba(37, 99, 235, 0.08);
  color: #0f766e !important;
  transform: translateY(-1px);
}

.nav-icon-btn {
  display: inline-flex;
  width: 2.5rem;
  height: 2.5rem;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(37, 99, 235, 0.18);
  border-radius: 999px;
  background-color: rgba(255, 255, 255, 0.88);
  box-shadow: 0 8px 18px rgba(37, 99, 235, 0.12);
  color: #1d4ed8;
  font-size: 1.1rem;
  transition: color 0.18s ease, background-color 0.18s ease, border-color 0.18s ease, transform 0.18s ease;
}

.nav-icon-btn:hover {
  border-color: rgba(20, 184, 166, 0.36);
  background-color: #eff6ff;
  color: #0f766e;
  transform: translateY(-1px);
}

.user-menu {
  gap: 0.4rem;
  padding: 0.32rem 0.56rem 0.32rem 0.36rem !important;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background-color: rgba(255, 255, 255, 0.72);
}

.user-avatar {
  width: 30px;
  height: 30px;
  border: 2px solid rgba(37, 99, 235, 0.25);
  box-shadow: 0 4px 10px rgba(15, 23, 42, 0.12);
  object-fit: cover;
}

.user-name {
  max-width: 8rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-menu-chevron {
  font-size: 0.72rem;
  opacity: 0.72;
}

.loj-user-dropdown {
  position: relative;
}

.loj-user-dropdown-menu {
  top: calc(100% + 0.45rem);
  right: 0;
  left: auto;
  z-index: 1060;
}

:global(.loj-user-dropdown-menu) {
  min-width: 10rem;
  padding: 0.45rem;
  border: 1px solid rgba(15, 23, 42, 0.1);
  border-radius: 0.85rem;
  background-color: rgba(255, 255, 255, 0.98);
  box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.15);
  color: #334155;
  animation: fadeInDown 0.2s ease-out;
  transform-origin: top right;
}

:global(.loj-user-dropdown-menu .n-dropdown-option) {
  margin: 0.05rem 0;
}

:global(.loj-user-dropdown-menu .n-dropdown-option-body) {
  min-height: auto;
  padding: 0.55rem 0.7rem;
  border-radius: 0.6rem;
  color: #334155;
  font-weight: 500;
}

:global(.loj-user-dropdown-menu .n-dropdown-option-body:hover),
:global(.loj-user-dropdown-menu .n-dropdown-option-body--pending) {
  background-color: #eef6ff;
  color: #0f766e;
}

:global(.loj-user-dropdown-menu .n-dropdown-divider) {
  margin: 0.35rem 0;
  background-color: rgba(15, 23, 42, 0.1);
}

:global(.loj-user-dropdown-icon) {
  margin-right: 0.5rem;
  font-size: 1rem;
}

:global(.loj-user-dropdown-label),
:global(.loj-user-dropdown-icon) {
  color: inherit;
}

:global(.loj-user-dropdown-label.is-danger),
:global(.loj-user-dropdown-icon.is-danger) {
  color: #dc3545;
}

:global(.loj-user-dropdown-menu .n-dropdown-option-body:hover .is-danger),
:global(.loj-user-dropdown-menu .n-dropdown-option-body--pending .is-danger) {
  color: #b02a37;
}

:global(html[data-loj-editor-theme="dark"] .loj-user-dropdown-menu) {
  color: #dbe7f5;
  background-color: #111827;
  border-color: #2c3a50;
}

:global(html[data-loj-editor-theme="dark"] .loj-user-dropdown-menu .n-dropdown-option-body) {
  color: #dbe7f5;
}

.animated-dropdown {
  animation: fadeInDown 0.2s ease-out;
}

.dropdown-menu {
  min-width: 10rem;
  padding: 0.45rem;
  border: 1px solid rgba(15, 23, 42, 0.1);
  border-radius: 0.85rem;
  background-color: rgba(255, 255, 255, 0.98);
  color: #334155;
}

.dropdown-item {
  padding: 0.55rem 0.7rem;
  border-radius: 0.6rem;
  color: #334155 !important;
  font-weight: 500;
}

.dropdown-item:hover {
  background-color: #eef6ff;
  color: #0f766e !important;
}

.dropdown-item.text-danger {
  color: #dc3545 !important;
}

.dropdown-item.text-danger:hover {
  background-color: rgba(220, 53, 69, 0.08);
  color: #b02a37 !important;
}

:global(html[data-loj-editor-theme="dark"] .loj-navbar) {
  background: #0b1220;
  border-bottom-color: #243247;
  box-shadow: 0 14px 34px rgba(0, 0, 0, 0.24);
}

:global(html[data-loj-editor-theme="dark"] .loj-brand),
:global(html[data-loj-editor-theme="dark"] .nav-link-button),
:global(html[data-loj-editor-theme="dark"] .auth-link),
:global(html[data-loj-editor-theme="dark"] .user-menu) {
  color: #dbe7f5 !important;
}

:global(html[data-loj-editor-theme="dark"] .nav-icon-btn),
:global(html[data-loj-editor-theme="dark"] .user-menu),
:global(html[data-loj-editor-theme="dark"] .navbar-toggler) {
  background-color: #111827;
  border-color: #2c3a50;
}

:global(html[data-loj-editor-theme="dark"] .dropdown-menu) {
  color: #dbe7f5;
  background-color: #111827;
  border-color: #2c3a50;
}

:global(html[data-loj-editor-theme="dark"] .dropdown-item) {
  color: #dbe7f5 !important;
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

@media (max-width: 767.98px) {
  .nav-collapse {
    margin-top: 0.75rem;
    padding: 0.75rem;
    border: 1px solid rgba(15, 23, 42, 0.08);
    border-radius: 0.95rem;
    background-color: rgba(255, 255, 255, 0.92);
    box-shadow: 0 14px 30px rgba(15, 23, 42, 0.1);
  }

  .nav-links,
  .nav-actions {
    align-items: stretch;
    gap: 0.25rem;
  }

  .nav-link-button,
  .auth-link,
  .user-menu {
    justify-content: flex-start;
  }

  .nav-icon-btn {
    width: 100%;
    height: 2.45rem;
    border-radius: 0.65rem;
  }

  .loj-user-dropdown-menu {
    position: static;
    width: 100%;
    margin-top: 0.35rem;
  }
}

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
