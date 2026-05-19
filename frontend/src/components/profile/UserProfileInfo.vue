<template>
  <div class="col-md-3 mb-4">
    <div class="card shadow-sm">
      <div class="card-body text-center">
        <img
          :src=" user.avatar"
          alt="Avatar"
          class="rounded-circle mb-3"
          width="100"
          height="100"
          style="cursor: pointer"
        >
        <n-divider />
        <h5 h5 class="card-title">{{ user.user_name }}</h5>
        <n-divider />
        <p class="text-muted">{{ user.email }}</p>
        <n-divider />
        <p class="text-muted">{{ "个性签名：" + user.quote }}</p>
        <n-divider />
        <p>Joined: {{ formattedDate(user.create_time) }}</p>
        <n-divider />

        <div class="d-flex justify-content-between w-100 mb-2 text-center">
          <div class="flex-fill">
            <div class="text-muted small">关注数</div>
            <div class="fw-bold">{{ user.following_count }}</div>
          </div>
          <div class="flex-fill">
            <div class="text-muted small">粉丝数</div>
            <div class="fw-bold">{{ user.follower_count }}</div>
          </div>
        </div>
        <n-divider />

        <n-button
          v-if="shouldShowFollowButton"
          size="small"
          type="primary"
          @click="toggleFollow"
        >
          {{ isFollowing ? '取消关注' : '关注' }}
        </n-button>
      </div>
    </div>
  </div>
</template>


<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useStore } from 'vuex';
import { useRoute } from 'vue-router';
import api from '@/api';

const props = defineProps({
  user: {
    type: Object,
    required: true
  }
});

const route = useRoute();
const store = useStore();
const user = computed(() => props.user);

const currentUserId = computed(() => store.getters['user/userId']); // 当前登录用户 ID，可按实际情况调整
const targetUserId = computed(() => parseInt(route.params.userId));
const isFollowing = ref(false);


const loadUserInfo = async () => {
  // const resp = await api.getUserProfileInfo(targetUserId.value);
  // if (resp.code === 0) {
  //   user.value = resp.data;
  //   await checkFollowing(); // 👈 拉数据后同步关注状态
  // }

  // console.log("user: ", user.value)
  await checkFollowing();
};

onMounted(() => {    // 每次进入页面时，调用一次
  loadUserInfo();
});

// watch(
//   () => route.params.userId,
//   (newUserId, oldUserId) => {
//     if (newUserId !== oldUserId) {
//       loadUserInfo(parseInt(newUserId));
//     }
//   }
// );

const checkFollowing = async () => {
  const resp = await api.checkFollowing(targetUserId.value);
  if (resp.code === 0) {
    isFollowing.value = resp.data.isFollowing;
  }
};

const shouldShowFollowButton = computed(() => {
  return currentUserId.value != targetUserId.value;
});

// 切换关注状态
const toggleFollow = async () => {
  if (isFollowing.value) {
    const resp = await api.unFollowUser(targetUserId.value);
    if (resp.code === 0) {
      isFollowing.value = false;
      user.value.follower_count--;
    }
  } else {
    const resp = await api.followUser(targetUserId.value);
    if (resp.code === 0) {
      isFollowing.value = true;
      user.value.follower_count++;
    }
  }
};

// 格式化日期函数
const formattedDate = (dateStr) =>
  new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

function isEmptyObject(obj) {
  return obj && Object.keys(obj).length === 0 && obj.constructor === Object;
}
</script>