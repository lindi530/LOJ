<template>
  <section class="feed-panel" aria-label="帖子列表">
    <div class="feed-toolbar">
      <span class="section-kicker">
        <i class="bi bi-chat-square-text"></i>
        LOJ 讨论区
      </span>

      <n-input
        :value="searchText"
        class="feed-search"
        placeholder="筛选帖子或作者"
        clearable
        @update:value="emit('update:searchText', $event)"
      >
        <template #prefix>
          <i class="bi bi-search"></i>
        </template>
      </n-input>
    </div>

    <div v-if="initialLoading" class="loading-state" aria-label="帖子加载中">
      <div v-for="item in 4" :key="item" class="skeleton-row">
        <span class="skeleton-avatar"></span>
        <span class="skeleton-copy"></span>
      </div>
    </div>

    <div v-else-if="posts.length" class="posts-list">
      <SimplePost
        v-for="post in posts"
        :key="post.post_id || post.id"
        :post="post"
      />
    </div>

    <div v-else class="empty-state">
      <i class="bi bi-inbox"></i>
      <strong>{{ emptyTitle }}</strong>
      <span>{{ emptyText }}</span>
      <button
        v-if="searchText"
        type="button"
        class="clear-filter"
        @click="emit('clear-search')"
      >
        清除筛选
      </button>
    </div>

    <InfiniteScrollTrigger
      v-if="hasLoadedPosts"
      :loading="loadingMore"
      :has-more="hasMore"
      @load-more="emit('load-more')"
    />
  </section>
</template>

<script setup>
import SimplePost from '@/components/posts/PostSimple.vue'
import InfiniteScrollTrigger from '@/components/posts/InfiniteScrollTrigger.vue'

defineProps({
  posts: {
    type: Array,
    default: () => []
  },
  hasLoadedPosts: {
    type: Boolean,
    default: false
  },
  initialLoading: {
    type: Boolean,
    default: false
  },
  loadingMore: {
    type: Boolean,
    default: false
  },
  hasMore: {
    type: Boolean,
    default: true
  },
  searchText: {
    type: String,
    default: ''
  },
  emptyTitle: {
    type: String,
    default: ''
  },
  emptyText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:searchText', 'clear-search', 'load-more'])
</script>

<style scoped>
.feed-panel {
  overflow: hidden;
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.08);
}

.feed-toolbar {
  display: grid;
  grid-template-columns: auto minmax(220px, 360px);
  gap: 14px;
  align-items: center;
  justify-content: space-between;
  padding: 16px 18px;
  border-bottom: 1px solid #edf0f3;
  background: #fbfcfd;
}

.section-kicker {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #12a875;
  font-size: 18px;
  font-weight: 800;
  flex: 0 0 auto;
  white-space: nowrap;
}

.feed-search {
  min-width: 0;
}

.posts-list {
  display: grid;
  gap: 12px;
  padding: 12px;
  background: #f3f6f5;
}

.posts-list :deep(.post-item) {
  padding: 24px 24px 22px;
  border: 1px solid #dde7e3;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
}

.posts-list :deep(.post-item:hover) {
  border-color: #b9d9cc;
  background: #f8fbfa;
}

.loading-state {
  display: grid;
  gap: 0;
}

.skeleton-row {
  display: grid;
  grid-template-columns: 54px minmax(0, 1fr);
  gap: 14px;
  padding: 24px;
  border-bottom: 1px solid #edf0f3;
}

.skeleton-avatar,
.skeleton-copy {
  display: block;
  background: linear-gradient(90deg, #eef2f6 0%, #f7f9fb 50%, #eef2f6 100%);
  background-size: 200% 100%;
  animation: loading-shimmer 1.2s ease-in-out infinite;
}

.skeleton-avatar {
  width: 54px;
  height: 54px;
  border-radius: 50%;
}

.skeleton-copy {
  height: 112px;
  border-radius: 8px;
}

.empty-state {
  display: grid;
  justify-items: center;
  gap: 8px;
  min-height: 280px;
  padding: 72px 24px;
  color: #8b96a5;
  text-align: center;
}

.empty-state i {
  color: #c4ccd6;
  font-size: 36px;
}

.empty-state strong {
  color: #344054;
  font-size: 18px;
}

.clear-filter {
  min-height: 34px;
  margin-top: 8px;
  padding: 0 14px;
  border: 1px solid #d7dee8;
  border-radius: 8px;
  background: #ffffff;
  color: #344054;
  font-weight: 700;
}

.clear-filter:hover {
  border-color: #12a875;
  color: #0f8f65;
}

@keyframes loading-shimmer {
  0% {
    background-position: 100% 0;
  }

  100% {
    background-position: -100% 0;
  }
}

@media (max-width: 760px) {
  .feed-panel {
    border-right: 0;
    border-left: 0;
    border-radius: 0;
  }

  .feed-toolbar {
    grid-template-columns: 1fr;
    padding: 14px 16px;
  }

  .section-kicker {
    font-size: 18px;
  }

  .posts-list :deep(.post-item) {
    padding: 20px 16px 18px;
  }

  .posts-list {
    gap: 10px;
    padding: 10px;
  }
}

@media (max-width: 520px) {
  .skeleton-row {
    grid-template-columns: 44px minmax(0, 1fr);
    padding: 20px 16px;
  }

  .skeleton-avatar {
    width: 44px;
    height: 44px;
  }

  .skeleton-copy {
    height: 126px;
  }
}
</style>
