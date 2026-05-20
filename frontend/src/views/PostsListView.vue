<template>
  <main class="posts-page">
    <section class="posts-shell">
      <div class="posts-head">
        <div>
          <h1>帖子</h1>
          <p>第 {{ currentPage }} 页，共 {{ total }} 条</p>
        </div>
      </div>

      <div v-if="posts && posts.length" class="posts-list">
        <SimplePost
          v-for="post in posts"
          :key="post.post_id"
          :post="post"
        />
      </div>

      <div v-else class="empty-state">
        还没有用户发布帖子
      </div>

      <n-space v-if="total > pageSize" vertical class="vertical-space">
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
  </main>
</template>

<script setup>
import { dialogError } from '@/utils/dialog'
import { useDialog } from 'naive-ui'
import { ref, onMounted } from 'vue'
import SimplePost from '@/components/post/PostSimple.vue'
import api from '@/api'

const posts = ref([])
const dialog = useDialog()
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchPosts = async () => {
  try {
    const resp = await api.getPosts(currentPage.value, pageSize.value)

    if (resp.code === 0) {
      posts.value = resp.data.list || []
      total.value = resp.data.total || 0
    } else {
      dialogError(dialog, '获取帖子失败', resp.message || '服务器响应失败')
    }
  } catch (error) {
    dialogError(dialog, '获取帖子失败', error?.message || '')
  }

  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(fetchPosts)
</script>

<style scoped>
.posts-page {
  min-height: calc(100vh - 60px);
  padding: 28px 16px 48px;
}

.posts-shell {
  width: min(100%, 960px);
  margin: 0 auto;
  padding: 24px 28px 30px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 8px;
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.08);
}

.posts-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 18px;
  border-bottom: 1px solid #edf0f3;
}

.posts-head h1 {
  margin: 0;
  color: #172033;
  font-size: 28px;
  font-weight: 700;
  line-height: 1.25;
}

.posts-head p {
  margin: 6px 0 0;
  color: #8b96a5;
  font-size: 14px;
}

.posts-list {
  display: flex;
  flex-direction: column;
}

.empty-state {
  padding: 48px 0;
  color: #8b96a5;
  text-align: center;
}

.vertical-space {
  align-items: center;
  margin-top: 26px;
}

@media (max-width: 640px) {
  .posts-page {
    padding: 14px 0 32px;
  }

  .posts-shell {
    padding: 18px 16px 24px;
    border-right: 0;
    border-left: 0;
    border-radius: 0;
  }

  .posts-head h1 {
    font-size: 24px;
  }
}
</style>
