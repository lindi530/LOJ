// 文件: components/HomePage.vue
<template>
  <div>
    <!-- Hero Section -->
    <section class="py-5 text-center bg-light">
      <div class="container">
        <h1 class="display-4 fw-bold">欢迎来到 My Blog</h1>
        <p class="lead text-muted">在这里分享你的故事，发现别人的视角</p>
        <router-link to="/posts" class="btn btn-primary btn-lg mt-3">查看所有帖子</router-link>
      </div>
    </section>

    <!-- Featured Posts -->
    <section class="py-5">
      <div class="container">
        <h2 class="mb-4">精选文章</h2>
        <div class="row">
          <div v-for="post in featuredPosts" :key="post.id" class="col-md-4 mb-4">
            <div class="card h-100 shadow-sm">
              <img :src="post.cover" class="card-img-top" alt="Post cover" />
              <div class="card-body d-flex flex-column">
                <h5 class="card-title">{{ post.title }}</h5>
                <p class="card-text text-truncate">{{ post.excerpt }}</p>
                <router-link :to="`/posts/${post.id}`" class="mt-auto btn btn-outline-primary">阅读全文</router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories -->
    <section class="py-5 bg-light">
      <div class="container">
        <h2 class="mb-4">热门分类</h2>
        <div class="d-flex flex-wrap gap-2">
          <router-link
            v-for="cat in categories"
            :key="cat.id"
            :to="`/category/${cat.id}`"
            class="btn btn-outline-secondary"
          >
            {{ cat.name }}
          </router-link>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="py-4 bg-dark text-white">
      <div class="container text-center">
        <small>&copy; {{ new Date().getFullYear() }} My Blog. 版权所有.</small>
      </div>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'HomePage',
  data() {
    return {
      featuredPosts: [],
      categories: []
    };
  },
  created() {
    this.fetchFeaturedPosts();
    this.fetchCategories();
  },
  methods: {
    fetchFeaturedPosts() {
      // TODO: 调用 API 获取精选文章，例如 /api/posts?featured=true
      this.featuredPosts = [
        { id: 1, title: '示例文章一', excerpt: '这是文章摘要...', cover: '/images/post1.jpg' },
        { id: 2, title: '示例文章二', excerpt: '这是文章摘要...', cover: '/images/post2.jpg' },
        { id: 3, title: '示例文章三', excerpt: '这是文章摘要...', cover: '/images/post3.jpg' }
      ];
    },
    fetchCategories() {
      // TODO: 调用 API 获取分类列表
      this.categories = [
        { id: 'tech', name: '技术' },
        { id: 'life', name: '生活' },
        { id: 'travel', name: '旅行' }
      ];
    }
  }
};
</script>

<style scoped>
.card-img-top {
  height: 200px;
  object-fit: cover;
}
.text-truncate {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
