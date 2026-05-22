<template>
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
</template>

<script setup>
defineProps({
  boardStats: {
    type: Array,
    default: () => []
  },
  quickLinks: {
    type: Array,
    default: () => []
  },
  activeAuthors: {
    type: Array,
    default: () => []
  }
})
</script>

<style scoped>
.posts-sidebar {
  display: grid;
  gap: 14px;
  align-self: start;
}

.side-panel {
  padding: 18px;
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.08);
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

@media (max-width: 1024px) {
  .posts-sidebar {
    position: static;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 760px) {
  .posts-sidebar {
    grid-template-columns: 1fr;
  }

  .side-panel {
    border-right: 0;
    border-left: 0;
    border-radius: 0;
  }
}
</style>
