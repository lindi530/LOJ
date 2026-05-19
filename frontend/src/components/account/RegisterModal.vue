<template>
  <div
    class="modal fade"
    tabindex="-1"
    :class="{ show: visible }"
    :style="visible ? 'display: block;' : ''"
    @click.self="close"
  >
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">用户注册</h5>
          <button type="button" class="btn-close" @click="close" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submit">
            <div class="mb-3">
              <label for="regUsername" class="form-label">用户名</label>
              <input
                type="text"
                class="form-control"
                id="regUsername"
                v-model="form.username"
                required
              />
            </div>
            <div class="mb-3">
              <label for="regPassword" class="form-label">密码</label>
              <input
                type="password"
                class="form-control"
                id="regPassword"
                v-model="form.password"
                required
                placeholder="以大写字母开头且总长度大于4"
              />
            </div>
            <div class="mb-3">
              <label for="re_password" class="form-label">确认密码</label>
              <input
                type="password"
                class="form-control"
                id="re_password"
                v-model="form.re_password"
                required
                placeholder="再次输入密码"
              />
            </div>
            <div class="d-flex align-items-center">
              <button type="submit" class="btn btn-success">注册</button>
              <!-- 新增的错误提示区域 -->
              <transition name="fade">
                <div v-if="errorMessage" class="text-danger ms-3">
                  <i class="bi bi-exclamation-circle-fill me-1"></i>
                  {{ errorMessage }}
                </div>
              </transition>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import api from '@/api';
export default {
  name: 'RegisterModal',
  props: {
    visible: { type: Boolean, default: false }
  },
  data() {
    return { 
      form: { 
        username: '', 
        password: '', 
        re_password: '' 
      },
      errorMessage: ''
    };
  },
  methods: {
    close() {
      this.$emit('update:visible', false);
      this.reset();
    },
    async submit() {
      this.errorMessage = ''
      try {
        console.log("Hello 1")
        const {username, password, re_password} = this.form
        const regRes = await this.register({
          username,
          password,
          re_password
        })
        if(regRes) {
          await this.login({
            username,
            password
          })
          this.errorMessage = '注册成功，正在登录！'
          setTimeout(() => {
            this.close(); // 1秒后执行
          }, 1000);
        }
      } catch(err) {
        alert(err.response?.data?.message || '注册失败，请重试');
      }
    },
    async register(data) {
      console.log("Hello 2")
      try {
        const regRes = await api.register(data, {
          headers: {
            'Content-Type': 'application/json'
          }
        })
        if (regRes.code) {
          this.errorMessage = regRes.message
          return false
        }
        return true
      } catch(err) {
          alert(err.response?.data?.message || '注册失败，请重试');
      }
    },
    async login(data) {
      console.log("Hello 3")
      try {
        const loginRes = await api.login (data, {
          headers: {
            'Content-Type': 'application/json'
          }
        });
        if (loginRes.code) {
          this.errorMessage = loginRes.message
          return
        }
        const { accessToken, refreshToken, user } = loginRes.data;
        // 存储到 Vuex
        this.$store.commit('user/SET_ACCESSTOKEN', accessToken);
        this.$store.commit('user/SET_REFRESHTOKEN', refreshToken);
        this.$store.commit('user/SET_PROFILE', user);
        this.$emit('login-success', user);
      } catch(err) {
          alert(err.response?.data?.message || '登录失败，请重试');
      }
    },
    reset() {
      this.form.username = '';
      this.form.password = '';
      this.form.re_password = '';
      this.errorMessage = ''
    }
  }
};
</script>

<style scoped>
.modal.fade { transition: none; }
.modal.show { background-color: rgba(0, 0, 0, 0.5); }
</style>
