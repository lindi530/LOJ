<template>
  <div class="container mt-5">
    <div class="row">
      <!-- 左侧头像 -->
      <div class="col-md-4 text-center">
        <div class="card">
          <img :src="form.avatar" class="card-img-top img-thumbnail" style="height: 300px; object-fit: cover;" />
          <div class="card-body">
            <input type="file" ref="fileInput" @change="onFileChange" accept="image/*" hidden />
            <button class="btn btn-primary" @click="triggerFileInput">更新头像</button>
          </div>
        </div>
      </div>

      <!-- 右侧用户信息 -->
      <div class="col-md-8">
        <div class="card">
          <div class="card-body">
            <h4 class="card-title mb-4">个人信息</h4>
            <form @submit.prevent="submitForm">
              <div class="mb-3">
                <label class="form-label">用户名：</label>
                <input v-model="form.user_name" type="text" class="form-control" />
              </div>
              <div class="mb-3">
                <label class="form-label">邮件地址：</label>
                <input v-model="form.email" type="email" class="form-control" />
              </div>
              <div class="mb-3">
                <label class="form-label">个性签名：</label>
                <textarea v-model="form.quote" class="form-control" rows="3" maxlength="50"></textarea>
                 <div class="form-text text-muted">
                   最多输入 50 个字符
                 </div>
                <div class="form-text text-end text-muted">
                  {{ form.quote?.length || 0}}/50 字
                </div>
              </div>
              <div class="mb-3">
                <label class="form-label">组织：</label>
                <input v-model="form.organization" type="text" class="form-control" />
              </div>
              <button class="btn btn-success" type="submit">更新信息</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ModifyUserProfile",
  props: {
    user: Object
  },
  emits: ['submit', 'upload-avatar'],
  data() {
    return {
      form: {}
    }
  },
  mounted() {
    this.form = { ...this.user }
  },
  methods: {
    triggerFileInput() {
      this.$refs.fileInput.click()
    },
    onFileChange(e) {
      const file = e.target.files[0]
      if (!file) return
      this.$emit('upload-avatar', file)  // 交给父组件处理上传
    },
    submitForm() {
      this.$emit('submit', { ...this.form })  // 提交表单
    },
    updateAvatarUrl(newUrl) {
      this.form.avatar = newUrl
    }
  },
  watch: {
    user: {
      deep: true,
      immediate: true,
      handler(val) {
        this.form = { ...val }
      }
    }
  }
}
</script>

<style scoped>
.card-img-top {
  max-width: 100%;
  border-radius: 8px;
}
</style>
