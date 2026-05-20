# 前端字段预留记录

本文档记录前端页面已经预留但后端可能尚未返回的字段，方便后续接口对接时确认“哪个链接需要哪些字段，以及这些字段显示在什么位置”。

## 帖子列表卡片

涉及前端位置：

- 首页：`/`
- 帖子列表：`/posts`
- 用户主页内的帖子列表：`/users/:userId`

涉及接口链接：

| 前端场景 | API 方法 | 接口链接 | 使用组件 |
| --- | --- | --- | --- |
| 首页 / 帖子列表 | `api.getPosts(page, pageSize)` | `GET /posts?page={page}&page_size={pageSize}` | `src/components/post/PostSimple.vue` |
| 用户主页帖子列表 | `api.getPostsByUserId(userId, page, pageSize)` | `GET /users/{userId}/posts?page={page}&page_size={pageSize}` | `src/components/post/PostSimple.vue` |

接口返回的每一条帖子数据会传给 `PostSimple.vue` 的 `post` 属性。以下字段已在前端做兼容预留：

| 字段路径 | 兼容字段名 | 用途 | 显示位置 | 缺失时表现 |
| --- | --- | --- | --- | --- |
| `post.author.level` | `post.author.lv` | 作者等级 | 用户名右侧的蓝色 `LV.x` 标签 | 不显示等级标签 |
| `post.author.school` | `post.author.university` | 作者学校 | 发布时间后方的灰色作者信息行 | 不显示学校 |
| `post.author.title` | `post.author.job`、`post.author.position` | 作者身份 / 职位 | 学校后方的灰色作者信息行 | 不显示身份 / 职位 |
| `post.images` | `post.pictures`、`post.photos`、`post.image`、`post.cover` | 帖子图片 | 正文摘要下方的图片网格 | 不显示图片区 |
| `post.comments` | `post.comment_count`、`post.comments_count` | 评论数量 | 底部评论按钮右侧 | 显示 `0` |

图片字段格式支持：

```json
{
  "images": [
    "https://example.com/post-1.jpg",
    "https://example.com/post-2.jpg"
  ]
}
```

也支持对象数组：

```json
{
  "images": [
    { "url": "https://example.com/post-1.jpg" },
    { "src": "https://example.com/post-2.jpg" }
  ]
}
```

推荐后端返回结构：

```json
{
  "post_id": 1,
  "title": "在大学的你，过得开心吗",
  "content": "帖子正文内容",
  "created_at": "2026-05-19T15:30:00+08:00",
  "updated_at": "2026-05-19T15:40:00+08:00",
  "likes": 7,
  "comment_count": 3,
  "images": ["https://example.com/post.jpg"],
  "author": {
    "user_id": 399,
    "user_name": "戴森球399",
    "avatar": "https://example.com/avatar.jpg",
    "level": 2,
    "school": "贵州医科大学",
    "title": "前端工程师"
  }
}
```

## 当前实现位置

字段读取逻辑集中在：

- `src/components/post/PostSimple.vue`

字段消费页面：

- `src/views/PostsListView.vue`
- `src/components/profile/UserPosts.vue`
