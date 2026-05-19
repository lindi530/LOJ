// src/utils/request.js
import axios from 'axios'
import store from '@/store'
import router from '@/router'
import api from '@/api'

const request = axios.create({
  baseURL: process.env.VUE_APP_API_BASE || 'http://localhost:8080',
  timeout: 10000
})

function isFormData(data) {
  return data instanceof FormData
}

// 请求拦截器
request.interceptors.request.use(config => {
  const accessToken = store.getters['user/accessToken'] || localStorage.getItem('accessToken')
  console.log("config.url: ", config.url)
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`
    if (!config.url.includes('/logout')) {
      refreshUserOnlineStatus(accessToken)
    }
  }
  if (isFormData(config.data)) {
    delete config.headers['Content-Type']
  }
  
  return config
}, error => Promise.reject(error))

let lastPing = 0
async function refreshUserOnlineStatus(token) {
  // 简单 ping，后端用来更新 Redis TTL
  const now = Date.now()
  if (now - lastPing < 1 * 1000) return
  lastPing = now

  if (!token) return
  await api.onlineState()
}

// 响应拦截器
request.interceptors.response.use(
  async res => {
    // 这里根据你后端业务自定义的 token 过期标志判断

    if (res.data.code === 1 && res.data.err_code === 3000) {
      // 模拟 401 状态时的 originalRequest 对象
      const originalRequest = res.config || {}

      // 用一个标记避免重复刷新
      if (!originalRequest._retry) {
        originalRequest._retry = true

        const refreshToken = store.getters['user/refreshToken'] || localStorage.getItem('refreshToken')
        if (!refreshToken) {
          refreshStore(store, '')
          router.push('/login')
          return Promise.reject(new Error('No refresh token, please login again'))
        }

        try {
          const refreshRes = await api.validateRefreshToken({ refresh_token: refreshToken })
          if (refreshRes?.code === 0 && refreshRes?.data?.accessToken) {
            const newAccessToken = refreshRes.data.accessToken
            refreshStore(store, newAccessToken)
            originalRequest._retry = true

            // 用 axios 实例重新发起请求
            return request(originalRequest)
          } else {
            refreshStore(store, '')
            router.push('/login')
            return Promise.reject(new Error('Refresh token invalid, please login again'))
          }
        } catch (e) {
          console.error('[REFRESH TOKEN ERROR]', e)
          refreshStore(store, '')
          router.push('/login')
          return Promise.reject(e)
        }
      } else {
        // 已经重试过，直接拒绝
        return Promise.reject(new Error('Token expired, please login again'))
      }
    }

    return res.data
  },
  async error => {
    const status = error.response?.status
    if (status === 404) {
      router.push('/404')
      return new Promise(() => {})
    }

    return Promise.reject(error)
  }
)

function refreshStore(store, accessToken) {
  store.commit('user/SET_ACCESSTOKEN', accessToken || '') // 这个store是异步的

  if (accessToken) {
    localStorage.setItem('accessToken', accessToken)
  } else {
    store.commit('user/SET_REFRESHTOKEN', '')
    store.commit('user/SET_PROFILE', {})
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    router.push('/login')
  }
}


export default request