import api from "@/api"
import { initWebSocket } from '@/composables/useWebSocket'

export async function restoreAuth(store) {
    const accessToken = localStorage.getItem('accessToken')
    const user = localStorage.getItem('user')
    if (!accessToken || !user) return
    console.log("restoreAuth !!!!!!!!!!!!!!!!")
    try {
        const userProfile = JSON.parse(user)
        const res1 = await api.validateAccessToken()
        if (res1.code == 0) {
            console.log("validateAccessToken success: ", accessToken)
            login(store, accessToken, userProfile)

            initWebSocket(accessToken)
        } else {
            const res2 = await api.validateRefreshToken({
                refresh_token: localStorage.getItem('refreshToken')
            })
            
            if (res2.code == 0) {
                console.log("refreshTOken 成功更新:", userProfile);
                const newAccessToken = res2.data.accessToken
                login(store, newAccessToken, userProfile)
                initWebSocket(newAccessToken)
            } else {
                console.log("refreshTOken 更新失败");
                logout(store)
            }
        }
    } catch (err) {
        // 清除状态
        logout(store)
    }
}

function login(store, accessToken, userProfile) {
    store.commit('user/SET_ACCESSTOKEN', accessToken)
    store.commit('user/SET_PROFILE', userProfile)
}

function logout(store) {
    console.log("LOGOUT")
    store.commit('user/LOGOUT')
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
}