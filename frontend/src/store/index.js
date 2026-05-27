import { createStore } from 'vuex'
import user from './modules/user';
import ui from './modules/ui';

export default createStore({
  state: {
    selectedCompetition: null
  },
  getters: {
    selectedCompetition: state => state.selectedCompetition
  },
  mutations: {
    SET_SELECTED_COMPETITION(state, competition) {
      state.selectedCompetition = competition
    }
  },
  actions: {
  },
  modules: {
    ui,
    user
  },
})
