// store/modules/ui.js
export default {
  namespaced: true,
  state: () => ({
    showChatModal: false
  }),
  mutations: {
    setShowChatModal(state, val) {
      state.showChatModal = val
    }
  }
}