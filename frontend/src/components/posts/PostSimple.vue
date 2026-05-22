<template>
  <article class="post-item">
    <div class="post-layout">
      <router-link
        :to="`/users/${post.author?.user_id}`"
        class="avatar-link"
      >
        <img
          :src="authorAvatar"
          :alt="`${authorName}的头像`"
          class="avatar"
        >
      </router-link>

      <div class="post-content-area">
        <header class="post-header">
          <div class="identity-row">
            <router-link
              :to="`/users/${post.author?.user_id}`"
              class="author-name"
            >
              {{ authorName }}
            </router-link>
            <span v-if="authorLevel" class="level-badge">LV.{{ authorLevel }}</span>
          </div>

          <div class="meta-row">
            <span v-if="formattedTime">{{ formattedTime }}</span>
            <span v-if="post.updated_at && post.updated_at !== post.created_at">已编辑</span>
            <span v-if="authorSchool">{{ authorSchool }}</span>
            <span v-if="authorTitle">{{ authorTitle }}</span>
          </div>
        </header>

        <router-link :to="`/posts/${post.post_id}`" class="post-link">
          <h2 class="post-title">{{ post.title || '无标题帖子' }}</h2>
          <p v-if="contentPreview" class="post-excerpt">
            {{ contentPreview }}
            <span v-if="isLongContent" class="read-more">查看更多</span>
          </p>
        </router-link>

        <div
          v-if="postImages.length"
          class="media-grid"
          :class="`media-count-${Math.min(postImages.length, 3)}`"
        >
          <router-link
            v-for="image in postImages.slice(0, 3)"
            :key="image"
            :to="`/posts/${post.post_id}`"
            class="media-link"
          >
            <img :src="image" alt="帖子图片" class="post-image">
          </router-link>
        </div>

        <footer class="post-toolbar">
          <button type="button" class="toolbar-item">
            <i class="bi bi-hand-thumbs-up"></i>
            <span>{{ post.likes || 0 }}</span>
          </button>

          <router-link :to="`/posts/${post.post_id}`" class="toolbar-item">
            <i class="bi bi-chat-square-dots"></i>
            <span>{{ commentCount }}</span>
          </router-link>
        </footer>
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

const authorName = computed(() => props.post.author?.user_name || '匿名用户')
const authorAvatar = computed(() => props.post.author?.avatar || '/default-avatar.svg')
const authorLevel = computed(() => props.post.author?.level || props.post.author?.lv || '')
const authorSchool = computed(() => props.post.author?.school || props.post.author?.university || '')
const authorTitle = computed(() => props.post.author?.title || props.post.author?.job || props.post.author?.position || '')

const commentCount = computed(() => {
  if (Array.isArray(props.post.comments)) {
    return props.post.comments.length
  }

  return props.post.comments || props.post.comment_count || props.post.comments_count || 0
})

const contentPreview = computed(() => {
  const content = props.post.content || ''
  return content.length > 132 ? `${content.slice(0, 132)}...` : content
})

const isLongContent = computed(() => (props.post.content || '').length > 132)

const formattedTime = computed(() => {
  if (!props.post.created_at) {
    return ''
  }

  const date = new Date(props.post.created_at)
  if (Number.isNaN(date.getTime())) {
    return props.post.created_at
  }

  const pad = value => String(value).padStart(2, '0')
  return `${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
})

const postImages = computed(() => {
  const source = props.post.images || props.post.pictures || props.post.photos || props.post.image || props.post.cover
  if (!source) {
    return []
  }

  const images = Array.isArray(source) ? source : [source]
  return images
    .map(item => (typeof item === 'string' ? item : item?.url || item?.src))
    .filter(Boolean)
})
</script>

<style scoped>
.post-item {
  position: relative;
  padding: 24px 8px 22px;
  background: #ffffff;
  border-bottom: 1px solid #edf0f3;
  transition: background-color 0.18s ease;
}

.post-item:hover {
  background: #fbfcfd;
}

.post-item:last-child {
  border-bottom: 0;
}

.post-layout {
  display: grid;
  grid-template-columns: 54px minmax(0, 1fr);
  gap: 14px;
}

.avatar-link {
  display: block;
  width: 54px;
  height: 54px;
  text-decoration: none;
}

.avatar {
  width: 54px;
  height: 54px;
  display: block;
  border: 2px solid #ffffff;
  border-radius: 50%;
  object-fit: cover;
  background: #eef2f6;
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.12);
}

.post-content-area {
  min-width: 0;
}

.post-header {
  display: flex;
  flex-direction: column;
  gap: 5px;
  min-width: 0;
}

.identity-row {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.author-name {
  max-width: 100%;
  overflow: hidden;
  color: #172033;
  font-size: 16px;
  font-weight: 700;
  line-height: 1.25;
  text-decoration: none;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.author-name:hover,
.post-link:hover .post-title {
  color: #12a875;
}

.level-badge {
  flex: 0 0 auto;
  padding: 1px 6px;
  color: #2f6df6;
  background: #e5efff;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 700;
  line-height: 1.45;
}

.meta-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  color: #98a2ad;
  font-size: 13px;
  line-height: 1.4;
}

.meta-row span + span::before {
  content: "·";
  margin: 0 7px;
  color: #c5ccd5;
}

.post-link {
  display: block;
  margin-top: 12px;
  color: inherit;
  text-decoration: none;
}

.post-title {
  margin: 0;
  color: #162033;
  font-size: 21px;
  font-weight: 750;
  line-height: 1.38;
}

.post-excerpt {
  margin: 8px 0 0;
  color: #273244;
  font-size: 17px;
  line-height: 1.72;
  white-space: pre-wrap;
  word-break: break-word;
}

.read-more {
  margin-left: 8px;
  color: #12a875;
  font-weight: 500;
  white-space: nowrap;
}

.media-grid {
  display: grid;
  gap: 8px;
  width: min(100%, 468px);
  margin-top: 14px;
}

.media-count-1 {
  width: min(100%, 220px);
  grid-template-columns: 1fr;
}

.media-count-2 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.media-count-3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.media-link {
  display: block;
  overflow: hidden;
  aspect-ratio: 1;
  border-radius: 8px;
  background: #edf2f7;
}

.post-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform 0.2s ease;
}

.media-link:hover .post-image {
  transform: scale(1.035);
}

.post-toolbar {
  display: flex;
  align-items: center;
  gap: 28px;
  margin-top: 18px;
  color: #87909b;
}

.toolbar-item {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  min-height: 30px;
  padding: 0;
  border: 0;
  background: transparent;
  color: inherit;
  font-size: 15px;
  line-height: 1;
  text-decoration: none;
  cursor: pointer;
}

.toolbar-item i {
  font-size: 18px;
}

.toolbar-item:hover {
  color: #12a875;
}

@media (max-width: 640px) {
  .post-item {
    padding: 20px 0 18px;
  }

  .post-layout {
    grid-template-columns: 44px minmax(0, 1fr);
    gap: 12px;
  }

  .avatar-link,
  .avatar {
    width: 44px;
    height: 44px;
  }

  .post-title {
    font-size: 19px;
  }

  .post-excerpt {
    font-size: 16px;
    line-height: 1.66;
  }

  .media-grid {
    width: min(100%, 360px);
  }
}
</style>
