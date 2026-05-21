<template>
  <aside class="profile-card">
    <div class="identity-block">
      <img
        :src="avatarSrc"
        :alt="`${displayName}的头像`"
        class="profile-avatar"
      >

      <h1 class="profile-name">{{ displayName }}</h1>
      <p class="profile-email">{{ emailLabel }}</p>
      <p class="profile-joined">加入于 {{ formattedDate(user.create_time) }}</p>

      <n-button
        v-if="shouldShowFollowButton"
        class="follow-button"
        size="medium"
        :type="isFollowing ? 'default' : 'primary'"
        @click="toggleFollow"
      >
        {{ isFollowing ? '取消关注' : '关注' }}
      </n-button>
    </div>

    <div class="follow-grid">
      <div>
        <strong>{{ followingCount }}</strong>
        <span>关注</span>
      </div>
      <div>
        <strong>{{ followerCount }}</strong>
        <span>粉丝</span>
      </div>
    </div>

    <section class="profile-section">
      <h2>个人简介</h2>
      <p class="profile-quote">{{ quoteLabel }}</p>
    </section>

    <section class="profile-section">
      <h2>社区统计</h2>
      <div class="stat-list">
        <div class="stat-row">
          <span>关注</span>
          <strong>{{ followingCount }}</strong>
        </div>
        <div class="stat-row">
          <span>粉丝</span>
          <strong>{{ followerCount }}</strong>
        </div>
      </div>
    </section>
  </aside>
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

const currentUserId = computed(() => store.getters['user/userId']);
const targetUserId = computed(() => parseInt(route.params.userId, 10));
const isFollowing = ref(false);

const displayName = computed(() => user.value?.user_name || '未命名用户');
const avatarSrc = computed(() => user.value?.avatar || '/default-avatar.svg');
const emailLabel = computed(() => user.value?.email || '未填写邮箱');
const quoteLabel = computed(() => user.value?.quote || '这个用户还没有留下签名');
const followingCount = computed(() => formatCount(user.value?.following_count));
const followerCount = computed(() => formatCount(user.value?.follower_count));

const loadUserInfo = async () => {
  await checkFollowing();
};

onMounted(() => {
  loadUserInfo();
});

watch(targetUserId, () => {
  loadUserInfo();
});

const checkFollowing = async () => {
  if (!targetUserId.value) {
    return;
  }

  const resp = await api.checkFollowing(targetUserId.value);
  if (resp.code === 0) {
    isFollowing.value = resp.data.isFollowing;
  }
};

const shouldShowFollowButton = computed(() => {
  return Number(currentUserId.value) !== targetUserId.value;
});

const toggleFollow = async () => {
  if (isFollowing.value) {
    const resp = await api.unFollowUser(targetUserId.value);
    if (resp.code === 0) {
      isFollowing.value = false;
      user.value.follower_count = Math.max(0, Number(user.value.follower_count || 0) - 1);
    }
  } else {
    const resp = await api.followUser(targetUserId.value);
    if (resp.code === 0) {
      isFollowing.value = true;
      user.value.follower_count = Number(user.value.follower_count || 0) + 1;
    }
  }
};

const formattedDate = (dateStr) => {
  if (!dateStr) {
    return '未知';
  }

  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) {
    return '未知';
  }

  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

function formatCount(value) {
  const count = Number(value || 0);

  if (count >= 10000) {
    return `${(count / 10000).toFixed(1)}万`;
  }

  return String(count);
}
</script>

<style scoped>
.profile-card {
  align-self: start;
  overflow: hidden;
  border: 1px solid #eef0f3;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(16, 24, 40, 0.06);
}

.identity-block {
  display: flex;
  padding: 24px 22px 20px;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.profile-avatar {
  width: 96px;
  height: 96px;
  display: block;
  border: 3px solid #ffffff;
  border-radius: 50%;
  background: #eef2f6;
  box-shadow: 0 10px 24px rgba(31, 35, 40, 0.14);
  object-fit: cover;
}

.profile-name {
  max-width: 100%;
  margin: 16px 0 4px;
  overflow: hidden;
  color: #1f2328;
  font-size: 22px;
  font-weight: 700;
  line-height: 1.25;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-email,
.profile-joined {
  max-width: 100%;
  margin: 0;
  overflow: hidden;
  color: #6a737d;
  font-size: 13px;
  line-height: 1.6;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.follow-button {
  width: 100%;
  margin-top: 18px;
  font-weight: 650;
}

.follow-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  border-top: 1px solid #f0f2f5;
  border-bottom: 1px solid #f0f2f5;
}

.follow-grid > div {
  display: flex;
  min-width: 0;
  padding: 16px 12px;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.follow-grid > div + div {
  border-left: 1px solid #f0f2f5;
}

.follow-grid strong {
  color: #1f2328;
  font-size: 18px;
  font-weight: 750;
  line-height: 1.2;
}

.follow-grid span {
  color: #8a939d;
  font-size: 12px;
}

.profile-section {
  padding: 18px 22px;
  border-bottom: 1px solid #f0f2f5;
}

.profile-section:last-child {
  border-bottom: 0;
}

.profile-section h2 {
  margin: 0 0 12px;
  color: #1f2328;
  font-size: 15px;
  font-weight: 700;
  line-height: 1.3;
}

.profile-quote {
  margin: 0;
  color: #4f5965;
  font-size: 14px;
  line-height: 1.65;
  word-break: break-word;
}

.stat-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: #6a737d;
  font-size: 14px;
}

.stat-row strong {
  color: #1f2328;
  font-weight: 700;
}
</style>
