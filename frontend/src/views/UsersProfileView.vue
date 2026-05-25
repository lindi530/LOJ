<template>
  <div class="profile-page">
    <div class="profile-shell">
      <UserProfileInfo class="profile-info-panel" :user="user"/>

      <UserSubmissionCalendar
        class="calendar-panel"
        :user-id="userId"
      />

      <main class="profile-main">
        <n-card class="activity-card" :bordered="false">
          <n-tabs type="line" class="profile-tabs">
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

        <NewPostForm
          v-if="canCreatePost"
          class="new-post-panel"
          @update-state="updateState"
        />
      </main>
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
import UserSubmissionCalendar from '../components/profile/UserSubmissionCalendar.vue';
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

function handlePostDeleted() {
  updatePosts.value = true;
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
.profile-page {
  min-height: calc(100vh - 60px);
  padding: 28px 0 56px;
}

.profile-shell {
  display: grid;
  grid-template-columns: minmax(240px, 286px) minmax(0, 1fr);
  grid-template-areas:
    "info calendar"
    ". main";
  gap: 18px 22px;
  width: min(1180px, calc(100% - 32px));
  margin: 0 auto;
  align-items: stretch;
}

.profile-info-panel {
  grid-area: info;
}

.calendar-panel {
  grid-area: calendar;
  align-self: stretch;
  min-width: 0;
}

.profile-main {
  grid-area: main;
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 18px;
}

.activity-card,
.new-post-panel {
  overflow: hidden;
  border: 1px solid #eef0f3;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(16, 24, 40, 0.06);
}

:deep(.activity-card > .n-card__content) {
  padding: 0;
}

:deep(.profile-tabs .n-tabs-nav) {
  padding: 10px 24px 0;
}

:deep(.profile-tabs .n-tab-pane) {
  padding: 10px 24px 24px;
}

:deep(.profile-tabs .n-tabs-tab) {
  font-size: 14px;
  font-weight: 600;
}

:deep(.new-post-panel > .n-card__content) {
  padding: 22px 24px 24px;
}

@media (max-width: 900px) {
  .profile-shell {
    grid-template-columns: 1fr;
    grid-template-areas:
      "info"
      "calendar"
      "main";
  }
}

@media (max-width: 576px) {
  .profile-page {
    padding: 18px 0 36px;
  }

  .profile-shell {
    width: min(100% - 20px, 1180px);
    gap: 14px;
  }

  :deep(.profile-tabs .n-tabs-nav) {
    padding: 0 16px;
  }

  :deep(.profile-tabs .n-tab-pane) {
    padding: 8px 16px 18px;
  }
}
</style>
