<template>
  <main class="posts-page">
    <div class="posts-layout">
      <section class="feed-panel" aria-label="帖子列表">
        <div class="feed-toolbar">
          <span class="section-kicker">
            <i class="bi bi-chat-square-text"></i>
            LOJ 讨论区
          </span>

          <n-input
            v-model:value="searchText"
            class="feed-search"
            placeholder="筛选帖子或作者"
            clearable
          >
            <template #prefix>
              <i class="bi bi-search"></i>
            </template>
          </n-input>
        </div>

        <div v-if="loading" class="loading-state" aria-label="帖子加载中">
          <div v-for="item in 4" :key="item" class="skeleton-row">
            <span class="skeleton-avatar"></span>
            <span class="skeleton-copy"></span>
          </div>
        </div>

        <div v-else-if="displayedPosts.length" class="posts-list">
          <SimplePost
            v-for="post in displayedPosts"
            :key="post.post_id"
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
            @click="searchText = ''"
          >
            清除筛选
          </button>
        </div>

        <n-space v-if="total > pageSize" vertical class="pagination-wrap">
          <n-pagination
            v-model:page="currentPage"
            :page-size="pageSize"
            :item-count="total"
            show-quick-jumper
            @update:page="fetchPosts"
          >
            <template #goto>
              跳转
            </template>
          </n-pagination>
        </n-space>
      </section>

      <aside class="posts-sidebar" aria-label="帖子侧栏">
        <section class="side-panel stats-panel">
          <h2>讨论看板</h2>
          <div class="stats-grid">
            <div v-for="stat in boardStats" :key="stat.label" class="stat-item">
              <span>{{ stat.label }}</span>
              <strong>{{ stat.value }}</strong>
            </div>
          </div>
        </section>

        <section class="side-panel quick-panel">
          <h2>快捷入口</h2>
          <router-link
            v-for="item in quickLinks"
            :key="item.label"
            :to="item.to"
            class="quick-link"
          >
            <span :class="['quick-icon', item.theme]">
              <i :class="item.icon"></i>
            </span>
            <span>{{ item.label }}</span>
            <i class="bi bi-chevron-right"></i>
          </router-link>
        </section>

        <section v-if="activeAuthors.length" class="side-panel authors-panel">
          <h2>活跃作者</h2>
          <router-link
            v-for="author in activeAuthors"
            :key="author.key"
            :to="author.profilePath"
            class="author-row"
          >
            <img :src="author.avatar" :alt="`${author.name}的头像`">
            <span>
              <strong>{{ author.name }}</strong>
              <small>{{ author.count }} 篇帖子 · {{ author.likes }} 赞</small>
            </span>
          </router-link>
        </section>
      </aside>
    </div>
  </main>
</template>

<script setup>
import { dialogError } from '@/utils/dialog'
import { useDialog } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useStore } from 'vuex'
import SimplePost from '@/components/post/PostSimple.vue'
import api from '@/api'

const posts = ref([])
const dialog = useDialog()
const store = useStore()
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const searchText = ref('')

const userId = computed(() => store.getters['user/userId'])

const quickLinks = computed(() => {
  const links = [
    { label: '题库训练', to: '/problems', icon: 'bi bi-code-square', theme: 'blue' },
    { label: '用户列表', to: '/users/userList', icon: 'bi bi-people', theme: 'green' }
  ]

  if (userId.value) {
    links.unshift({
      label: '我的主页',
      to: `/users/${userId.value}`,
      icon: 'bi bi-person-circle',
      theme: 'amber'
    })
  }

  return links
})

const displayedPosts = computed(() => {
  const keyword = searchText.value.trim().toLowerCase()
  return posts.value.filter(post => {
    if (!keyword) {
      return true
    }

    const searchable = [
      post.title,
      post.content,
      post.author?.user_name
    ].filter(Boolean).join(' ').toLowerCase()

    return searchable.includes(keyword)
  })
})

const pageLikes = computed(() => posts.value.reduce((sum, post) => sum + Number(post.likes || 0), 0))
const pageViews = computed(() => posts.value.reduce((sum, post) => sum + Number(post.views || 0), 0))

const boardStats = computed(() => [
  { label: '浏览热度', value: formatNumber(pageViews.value) },
  { label: '互动热度', value: formatNumber(pageLikes.value) },
  { label: '活跃作者', value: formatNumber(activeAuthors.value.length) }
])

const activeAuthors = computed(() => {
  const authorMap = new Map()

  posts.value.forEach(post => {
    const authorId = post.author?.user_id || post.user_id || ''
    const key = authorId || `post-${post.post_id}`

    if (!authorMap.has(key)) {
      authorMap.set(key, {
        key,
        name: post.author?.user_name || '匿名用户',
        avatar: post.author?.avatar || '/default-avatar.svg',
        count: 0,
        likes: 0,
        profilePath: authorId ? `/users/${authorId}` : `/posts/${post.post_id}`
      })
    }

    const author = authorMap.get(key)
    author.count += 1
    author.likes += Number(post.likes || 0)
  })

  return [...authorMap.values()]
    .sort((a, b) => b.count - a.count || b.likes - a.likes)
    .slice(0, 4)
})

const emptyTitle = computed(() => (posts.value.length ? '没有匹配的帖子' : '还没有用户发布帖子'))
const emptyText = computed(() => (posts.value.length ? '换个关键词再试试' : '新的讨论会出现在这里'))

const fetchPosts = async () => {
  loading.value = true

  try {
    const resp = await api.getPosts(currentPage.value, pageSize.value)

    if (resp.code === 0) {
      posts.value = resp.data.list || resp.data.posts || []
      total.value = resp.data.total || 0
    } else {
      dialogError(dialog, '获取帖子失败', resp.message || '服务器响应失败')
    }
  } catch (error) {
    dialogError(dialog, '获取帖子失败', error?.message || '')
  } finally {
    loading.value = false
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function formatNumber(value) {
  return Number(value || 0).toLocaleString('zh-CN')
}

onMounted(fetchPosts)
</script>

<style scoped>
.posts-page {
  min-height: calc(100vh - 60px);
  padding: 28px 18px 52px;
  color: #172033;
}

.posts-layout {
  width: min(100%, 1180px);
  margin: 0 auto;
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

.posts-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 18px;
  align-items: start;
}

.feed-panel,
.side-panel {
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.08);
}

.feed-panel {
  overflow: hidden;
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

.feed-search {
  min-width: 0;
}

.posts-list {
  background: #ffffff;
}

.posts-list :deep(.post-item) {
  padding: 24px 24px 22px;
}

.posts-list :deep(.post-item:hover) {
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

.pagination-wrap {
  align-items: center;
  padding: 18px 16px 20px;
  border-top: 1px solid #edf0f3;
  background: #fbfcfd;
}

.posts-sidebar {
  display: grid;
  gap: 14px;
  align-self: start;
}

.side-panel {
  padding: 18px;
}

.side-panel h2 {
  margin: 0 0 14px;
  color: #1f2937;
  font-size: 16px;
  font-weight: 800;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.stat-item {
  min-height: 78px;
  padding: 12px;
  border: 1px solid #edf0f3;
  border-radius: 8px;
  background: #fbfcfd;
}

.stat-item span {
  display: block;
  color: #8b96a5;
  font-size: 12px;
  line-height: 1.4;
}

.stat-item strong {
  display: block;
  margin-top: 8px;
  color: #172033;
  font-size: 22px;
  font-weight: 800;
  line-height: 1;
}

.quick-panel {
  display: grid;
  gap: 10px;
}

.quick-link {
  display: grid;
  grid-template-columns: 34px minmax(0, 1fr) 16px;
  gap: 10px;
  align-items: center;
  min-height: 46px;
  padding: 6px 0;
  color: #344054;
  text-decoration: none;
}

.quick-link:hover {
  color: #0f8f65;
}

.quick-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 8px;
  font-size: 17px;
}

.quick-icon.amber {
  background: #fff4df;
  color: #b86e00;
}

.quick-icon.blue {
  background: #eaf2ff;
  color: #2f6df6;
}

.quick-icon.green {
  background: #e9f8f2;
  color: #0f8f65;
}

.authors-panel {
  display: grid;
  gap: 10px;
}

.author-row {
  display: grid;
  grid-template-columns: 38px minmax(0, 1fr);
  gap: 10px;
  align-items: center;
  min-height: 48px;
  color: #344054;
  text-decoration: none;
}

.author-row:hover strong {
  color: #0f8f65;
}

.author-row img {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  object-fit: cover;
  background: #eef2f6;
}

.author-row span {
  min-width: 0;
}

.author-row strong,
.author-row small {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.author-row strong {
  color: #172033;
  font-size: 14px;
  line-height: 1.35;
}

.author-row small {
  color: #8b96a5;
  font-size: 12px;
  line-height: 1.4;
}

@keyframes loading-shimmer {
  0% {
    background-position: 100% 0;
  }

  100% {
    background-position: -100% 0;
  }
}

@media (max-width: 1024px) {
  .posts-layout {
    grid-template-columns: minmax(0, 1fr);
  }

  .posts-sidebar {
    position: static;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 760px) {
  .posts-page {
    padding: 18px 0 36px;
  }

  .section-kicker {
    font-size: 18px;
  }

  .feed-panel,
  .side-panel {
    border-right: 0;
    border-left: 0;
    border-radius: 0;
  }

  .feed-toolbar {
    grid-template-columns: 1fr;
    padding: 14px 16px;
  }

  .posts-list :deep(.post-item) {
    padding: 20px 16px 18px;
  }

  .posts-sidebar {
    grid-template-columns: 1fr;
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
