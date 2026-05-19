const state = () => ({
  accessToken: '',
  refreshToken: '',
  profile: {} // 用户信息
});

const mutations = {
  SET_REFRESHTOKEN(state, token) {
    state.refreshToken = token;
    localStorage.setItem('refreshToken', token)
  },
  SET_ACCESSTOKEN(state, token) {
    state.accessToken = token;
    localStorage.setItem('accessToken', token)
  },
  SET_PROFILE(state, user) {
    state.profile = user;
    localStorage.setItem('user', JSON.stringify(user))
  },
  SET_PROFILE_AVATAR(state, avatar) {
    state.profile.avatar = avatar;
  },
  MODIFT_PROFILE(state, profile) {
    state.profile.user_name = profile.user_name;
    state.profile.email = profile.email;
    state.profile.quote = profile.quote;
  },
  LOGOUT(state) {
    state.accessToken = '';
    state.refreshToken = '';
    state.profile = {};
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }
};

const actions = {
  // logout({ commit }) {
  //   commit('LOGOUT');
  // }
};

const getters = {
  isLogin: state => !!state.accessToken,
  accessToken: state => state?.accessToken,
  refreshToken: state => state?.refreshToken,
  userName: state => state.profile?.user_name || '',
  userQuote: state => state.profile?.quote || '',
  userEmail: state => state.profile?.email || '',
  updateTime: state => state.profile?.update_time || '',
  userAvatar: state => state.profile?.avatar || 'https://cdn.acwing.com/media/user/profile/photo/89646_sm_28e4eb758d.jpg',
  userId: state => state.profile?.user_id || '', 
  userInfo: state => state.profile,
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
};
