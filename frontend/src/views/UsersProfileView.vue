<template>
  <div class="container my-4">
    <!-- 上半部分的 row -->
    <div class="row">
      <UserProfileInfo :user="user"/>
      <n-card class="col-md-9 mb-4">
        <n-tabs type="line" animated>
          <n-tab-pane name="posts" tab="帖子">
            <UserPosts 
              :update-posts="updatePosts"
              :user-id="userId"
              @post-deleted="handlePostDeleted"
              @update-state="updateState"
            />
          </n-tab-pane>
          <n-tab-pane name="submissions" tab="提交记录">
            <UserSubmissions 
              :author="user"
            />
          </n-tab-pane>
        </n-tabs>
      </n-card>
    </div>

    <!-- 下半部分也放入 row，并使用相同的栅格列（如 col-md-9） -->
    <div class="row">
        <NewPostForm
          v-if="canCreatePost"
          @update-state="updateState"
        />
    </div>
  </div>
</template>

<script setup>
import { dialogError } from '@/utils/dialog.js'
import { ref, computed, onMounted, watch } from 'vue';
import { useStore } from 'vuex';
import { useRoute } from 'vue-router';
import { useDialog }  from 'naive-ui'
import UserProfileInfo from '../components/profile/UserProfileInfo.vue';
import UserPosts from '../components/profile/UserPosts.vue';
import UserSubmissions from '../components/profile/UserSubmissions.vue';
import NewPostForm from '../components/profile/NewPostForm.vue';
import api from '@/api';

const user = ref({});
const updatePosts = ref(false);
const dialog = useDialog()
const route = useRoute();

const userId = computed(() => route.params.userId);
const store = useStore();

const loginUserId = computed(() => store.getters['user/userId']);

const canCreatePost = computed(() => {
  return userId.value == loginUserId.value;
})

function updateState(val) {
  updatePosts.value = val
}

function handlePostDeleted(deletedPostId) {
  posts.value = posts.value.filter(post => post.id !== deletedPostId);
}

onMounted(async () => {
  loadUserProfile()
});

watch(
  () => route.params.userId,
  (newUserId, oldUserId) => {
    if (newUserId !== oldUserId) {
      loadUserProfile();
    }
  }
);

const loadUserProfile = async () => {
  try {
    const profileResp = await api.getUserProfileInfo(userId.value);
    if (profileResp.code === 0) {
      user.value = profileResp.data;
    } else {
      dialogError(dialog, "用户信息请求发送成功", "服务器响应失败")
    }
  } catch {
    dialogError(dialog, "用户信息请求发送失败", "")
  }

  // try {
  //   const postsResp = await api.getPostsByUserId(userId.value);
  //   if (postsResp.code === 0) {
  //     posts.value = postsResp.data.posts;
  //   } else {
  //     dialogError(dialog, "帖子信息请求发送成功", "服务器响应失败")
  //   }
  // } catch { 
  //   dialogError(dialog, "帖子信息请求发送失败", "")
  // }
};

</script>

<style scoped>
.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
