module.exports = {
  devServer: {
    allowedHosts: 'all',
    proxy: {
      // 当请求以 /users 开头时，将请求代理到后端服务
      '/users': {
        target: process.env.VUE_APP_API_BASE,
        changeOrigin: true,
      },
      '/posts': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
};
