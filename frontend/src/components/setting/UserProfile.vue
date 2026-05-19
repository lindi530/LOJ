<template>
  <div class="container my-4">
    <ModifyUserProfile
      ref="userForm"
      :user="userCopy"
      @submit="handleUserUpdate"
      @upload-avatar="handleAvatarUpload"
    />
  </div>
</template>

<script>
import ModifyUserProfile from './ModifyUserProfile.vue'
import api from '@/api'

export default {
  components: { ModifyUserProfile },
  name: "UserProfile",
  computed: {
    currentUser() {
      console.log("currentUser:", this.$store.getters['user/userInfo'])
      return this.$store.getters['user/userInfo']
    },
    userId() {
      return this.$store.getters['user/userId']
    }
  },
  data() {
    return {
      userCopy: {}
    }
  },
  watch: {
    currentUser: {
      immediate: true,
      handler(val) {
        this.userCopy = { ...val }
      }
    }
  },
  methods: {
    async handleUserUpdate(formData) {
      try {
        
        const upFormData = {
          "user_name": formData.user_name,
          "email": formData.email,
          "quote": formData.quote,
        }

        console.log("upFormData: ", upFormData)

        const res = await api.modifyUserProfile(this.userId, upFormData)
        console.log(res)
        if (res.code == 0) {
          this.$store.commit('user/MODIFT_PROFILE', upFormData);
        }
        
      } catch (err) {
        console.error(err)
      }
    },
    async handleAvatarUpload(file) {
      const formData = new FormData()
      formData.append('avatar', file)
      try {
        const res = await api.upLoadAvatar(this.userId, formData)
        const newUrl = res.data.avatar
        console.log("res: ", res)
        if (!newUrl) throw new Error('上传失败，服务器未返回头像地址')
        console.log("newUrl: ", newUrl)
        // 更新store
        this.$store.commit('user/SET_PROFILE_AVATAR', newUrl);
        // 更新子组件中的 avatar URL
        this.$refs.userForm.updateAvatarUrl(newUrl)
        console.log("子头像更新成功")
        // 同步更新子组件 userCopy 中的 avatar，便于下次提交
        this.userCopy.avatar = newUrl
        console.log("头像上传成功")
      } catch (err) {
        console.error('头像上传失败', err)
      }
    }
  }
}
</script>