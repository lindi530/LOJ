<template>
  <div id="app">
    <!-- 背景层 -->
    <div id="building"></div>

    <!-- 内容层 -->
    <n-config-provider :hljs="hljs">
      <n-dialog-provider>
        <n-message-provider>
          <NavBar />
          <router-view />
        </n-message-provider>
      </n-dialog-provider>
    </n-config-provider>
  </div>
</template>

<script setup>
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap/dist/js/bootstrap";
import "bootstrap-icons/font/bootstrap-icons.css";

import NavBar from "./components/NavBar";
import hljs from 'highlight.js/lib/core';
import python from "highlight.js/lib/languages/python";
import go from "highlight.js/lib/languages/go";
import cpp from "highlight.js/lib/languages/cpp";

hljs.registerLanguage("python", python);
hljs.registerLanguage("go", go);
hljs.registerLanguage("cpp", cpp);

// 背景图片
import bgImg from "@/assets/1.png";
</script>

<style scoped>
/* 背景层 */
#building {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;          /* 填满视口 */
  z-index: -1;            /* 放到内容下方 */
  background-image: url(~@/assets/back.jpg); /* 或使用 import 后动态绑定 */
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

/* 内容层 */
#app {
  position: relative;
  min-height: 100vh;      /* 保证内容至少撑开屏幕 */
  overflow: auto;         /* 页面可滚动 */
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
    padding-top: 60px; /* 导航栏高度，根据实际调整 */
  margin: 0; /* 清除默认 margin，避免影响布局 */
}


/* NavBar 等样式 */
nav {
  padding: 30px;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
}

nav a.router-link-exact-active {
  color: #42b983;
}
</style>
