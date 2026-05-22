<template>
  <main class="posts-page">
    <div class="posts-layout">
      <PostsFeedPanel
        v-model:search-text="searchText"
        :posts="displayedPosts"
        :has-loaded-posts="posts.length > 0"
        :initial-loading="initialLoading"
        :loading-more="loadingMore"
        :has-more="cursorReq.has_more"
        :empty-title="emptyTitle"
        :empty-text="emptyText"
        @clear-search="searchText = ''"
        @load-more="fetchPosts"
      />

      <PostsSidebar
        :board-stats="boardStats"
        :quick-links="quickLinks"
        :active-authors="activeAuthors"
      />
    </div>
  </main>
</template>

<script setup>
import { dialogError } from '@/utils/dialog'
import { useDialog } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useStore } from 'vuex'
import PostsFeedPanel from '@/components/posts/PostsFeedPanel.vue'
import PostsSidebar from '@/components/posts/PostsSidebar.vue'
import api from '@/api'

const posts = ref([])
const dialog = useDialog()
const store = useStore()
const initialLoading = ref(false)
const loadingMore = ref(false)
const searchText = ref('')
const cursorReq = ref({
  cursor: 0,
  has_more: true
})

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
    const key = authorId || `post-${getPostCursor(post)}`

    if (!authorMap.has(key)) {
      authorMap.set(key, {
        key,
        name: post.author?.user_name || '匿名用户',
        avatar: post.author?.avatar || '/default-avatar.svg',
        count: 0,
        likes: 0,
        profilePath: authorId ? `/users/${authorId}` : `/posts/${getPostCursor(post)}`
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
  if (initialLoading.value || loadingMore.value || !cursorReq.value.has_more) {
    return
  }

  const requestCursor = {
    cursor: cursorReq.value.cursor,
    has_more: cursorReq.value.has_more
  }
  const isInitialRequest = posts.value.length === 0 && requestCursor.cursor === 0

  if (isInitialRequest) {
    initialLoading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    const resp = await api.getCursorPosts(requestCursor)

    if (resp.code === 0) {
      const payload = normalizeCursorPayload(resp)
      posts.value = mergePosts(posts.value, payload.posts)
      cursorReq.value = buildNextCursorReq(payload.cursor, payload.posts, requestCursor.cursor)
    } else {
      dialogError(dialog, '获取帖子失败', resp.message || '服务器响应失败')
    }
  } catch (error) {
    dialogError(dialog, '获取帖子失败', error?.message || '')
  } finally {
    initialLoading.value = false
    loadingMore.value = false
  }
}

function normalizeCursorPayload(resp) {
  const data = resp?.data || {}
  const payload = data?.data || data?.Data || data
  const postsData = payload?.posts || payload?.Posts || []
  const cursorData = payload?.cursor || payload?.Cursor || {}

  return {
    posts: Array.isArray(postsData) ? postsData : [],
    cursor: cursorData || {}
  }
}

function buildNextCursorReq(cursorData, receivedPosts, previousCursor) {
  const responseCursor = normalizeNumber(cursorData?.cursor ?? cursorData?.Cursor)
  const nextCursor = responseCursor > 0
    ? responseCursor
    : resolveLastPostCursor(receivedPosts, previousCursor)
  const responseHasMore = normalizeBoolean(cursorData?.has_more ?? cursorData?.HasMore)

  return {
    cursor: nextCursor,
    has_more: responseHasMore && nextCursor !== previousCursor
  }
}

function mergePosts(currentPosts, nextPosts) {
  const seen = new Set(currentPosts.map(post => getPostCursor(post)).filter(Boolean))
  const mergedPosts = [...currentPosts]

  nextPosts.forEach(post => {
    const postCursor = getPostCursor(post)
    if (!postCursor || !seen.has(postCursor)) {
      mergedPosts.push(post)
    }
    if (postCursor) {
      seen.add(postCursor)
    }
  })

  return mergedPosts
}

function resolveLastPostCursor(receivedPosts, fallbackCursor) {
  for (let index = receivedPosts.length - 1; index >= 0; index -= 1) {
    const postCursor = getPostCursor(receivedPosts[index])
    if (postCursor > 0) {
      return postCursor
    }
  }

  return fallbackCursor
}

function getPostCursor(post) {
  return normalizeNumber(post?.post_id ?? post?.id ?? post?.PostID)
}

function normalizeNumber(value) {
  const numberValue = Number(value || 0)
  return Number.isFinite(numberValue) ? numberValue : 0
}

function normalizeBoolean(value) {
  if (typeof value === 'boolean') {
    return value
  }

  if (typeof value === 'string') {
    return value.toLowerCase() === 'true'
  }

  return Boolean(value)
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
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 18px;
  align-items: start;
  width: min(100%, 1180px);
  margin: 0 auto;
}

@media (max-width: 1024px) {
  .posts-layout {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 760px) {
  .posts-page {
    padding: 18px 0 36px;
  }
}
</style>
