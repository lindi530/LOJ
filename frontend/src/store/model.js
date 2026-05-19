import { defineStore } from 'pinia'

export const useModalStore = defineStore('modal', {
  state: () => ({
    showLogin: false,
    showChat: false
  }),
  actions: {
    open(modalName) {
      this[modalName] = true
    },
    close(modalName) {
      this[modalName] = false
    }
  }
})
