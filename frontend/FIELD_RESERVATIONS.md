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

## 用户主页提交记录

涉及前端位置：

- 用户主页内的提交记录 Tab：`/users/:userId`

涉及接口链接：

| 前端场景 | API 方法 | 接口链接 | 使用组件 |
| --- | --- | --- | --- |
| 用户主页提交记录 | `api.getUserSubmissionList(userId, page, pageSize)` | `GET /users/{userId}/submissions?page={page}&page_size={pageSize}` | `src/components/coding/SubmissionSimple.vue` |

接口返回的每一条提交数据会传给 `SubmissionSimple.vue` 的 `submission` 属性。以下字段已在前端做兼容预留：

| 字段路径 | 兼容字段名 | 用途 | 显示位置 | 缺失时表现 |
| --- | --- | --- | --- | --- |
| `submission.state` | `submission.status`、`submission.result` | 运行结果，如答案正确、答案错误、编译错误 | 提交记录表格的“运行结果”列 | 如果有 `score`，按分数兜底显示；否则显示“查看详情” |
| `submission.language` | `submission.lang` | 使用语言，如 `C++`、`Go`、`Java` | 提交记录表格的“使用语言”列 | 显示 `-` |
| `submission.exec_time` | `submission.execTime`、`submission.runtime` | 运行时间，单位 ms | 提交记录表格的“运行时间(ms)”列 | 显示 `-` |
| `submission.memory_usage` | `submission.memoryUsage`、`submission.memory` | 使用内存，单位 KB | 提交记录表格的“使用内存(KB)”列 | 显示 `-` |
| `submission.code_length` | `submission.codeLength` | 代码长度 | 提交记录表格的“代码长度”列 | 如果有 `code`，用 `code.length` 计算；否则显示 `-` |
| `submission.score` | - | 分数 / 判题得分 | 不单独显示，用作运行结果兜底判断 | 不参与兜底 |

说明：

- 题目列只显示 `submission.title`，不显示题号；`submission.problem_id` 仍用于点击题目名跳转题目详情页。
- 运行结果字段会把常见英文状态映射成中文展示，例如 `accepted` / `ac` 显示为“答案正确”，`wrong_answer` / `wa` 显示为“答案错误”。
- 如果后端没有返回运行结果、运行时间、内存或语言，前端会保留表格结构并显示 `-`，方便识别是数据缺失而不是布局问题。

推荐后端返回结构：

```json
{
  "id": 83652351,
  "problem_id": 1001,
  "title": "Two Sum",
  "state": "accepted",
  "language": "C++",
  "exec_time": 1563,
  "memory_usage": 81244,
  "score": 100,
  "code": "int main(){return 0;}",
  "created_at": "2026-05-20T19:15:46+08:00"
}
```

## 当前实现位置

字段读取逻辑集中在：

- `src/components/post/PostSimple.vue`
- `src/components/coding/SubmissionSimple.vue`

字段消费页面：

- `src/views/PostsListView.vue`
- `src/components/profile/UserPosts.vue`
- `src/components/profile/UserSubmissions.vue`
