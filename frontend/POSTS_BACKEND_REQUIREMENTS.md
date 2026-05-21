# 帖子列表页后端配合需求

本文档记录 `PostsListView.vue` 新布局需要后端配合补齐的能力。每一项都标明对应的前端功能和原因，方便后续拆分接口任务。

## 1. 帖子列表支持全站排序

**后端功能**

`GET /posts` 增加排序参数：

```text
GET /posts?page=1&page_size=10&sort=latest
GET /posts?page=1&page_size=10&sort=hot
GET /posts?page=1&page_size=10&sort=views
```

建议含义：

- `latest`：按 `created_at DESC`
- `hot`：按综合热度排序，例如 `likes * 3 + views + comment_count * 2`
- `views`：按 `views DESC`

**对应前端功能**

[PostsListView.vue](src/views/PostsListView.vue) 顶部的排序切换：

- 最新
- 热门
- 高浏览

**理由**

现在前端只能对当前页已经拿到的帖子做本地排序。这样会出现“热门”只是在当前 10 条里热门，而不是全站热门的问题。排序下沉到后端后，分页和排序才能保持一致。

## 2. 帖子列表支持全站搜索

**后端功能**

`GET /posts` 增加搜索参数：

```text
GET /posts?page=1&page_size=10&keyword=动态规划
```

建议搜索范围：

- 帖子标题 `title`
- 帖子内容 `content`
- 作者昵称 `author.user_name`

**对应前端功能**

[PostsListView.vue](src/views/PostsListView.vue) 顶部搜索框：

```text
筛选帖子或作者
```

**理由**

现在搜索只筛选当前页的数据。用户搜索一个不在当前页、但存在于其它页的帖子时，会误以为没有结果。后端搜索可以返回全站匹配结果，并与分页一起工作。

## 3. 帖子列表返回评论数量

**后端功能**

`GET /posts` 的每条帖子增加字段：

```json
{
  "comment_count": 12
}
```

字段名可以兼容前端已有兜底：

- `comment_count`
- `comments_count`
- `comments`

推荐使用 `comment_count`，语义最清楚。

**对应前端功能**

[PostSimple.vue](src/components/post/PostSimple.vue) 底部评论入口显示评论数。

**理由**

帖子列表页不应该为了显示评论数量再额外请求每篇帖子的评论列表。后端直接返回聚合后的评论数量，可以减少请求次数，也能让列表展示更准确。

## 4. 帖子列表返回更完整的作者信息

**后端功能**

`GET /posts` 的 `author` 建议补充：

```json
{
  "author": {
    "user_id": 1,
    "user_name": "Alice",
    "avatar": "/images/avatar.png",
    "level": 3,
    "school": "LOJ University",
    "title": "算法爱好者"
  }
}
```

**对应前端功能**

[PostSimple.vue](src/components/post/PostSimple.vue) 作者区域已经预留展示：

- 等级 `LV.x`
- 学校
- 头衔 / 职业 / 身份

**理由**

当前后端 `AuthorInfo` 只有 `user_id`、`user_name`、`avatar`。前端虽然写好了展示位，但拿不到数据，页面信息密度会偏低。补齐作者资料后，列表会更接近牛客、LeetCode 讨论区的社区感。

## 5. 提供帖子侧栏数据接口

**后端功能**

新增一个侧栏聚合接口，例如：

```text
GET /posts/sidebar
```

建议返回：

```json
{
  "stats": {
    "views": 12031,
    "likes": 894,
    "active_authors": 36
  },
  "active_authors": [
    {
      "user_id": 1,
      "user_name": "Alice",
      "avatar": "/images/avatar.png",
      "post_count": 8,
      "like_count": 126
    }
  ],
  "hot_posts": [
    {
      "post_id": 10,
      "title": "这道 DP 题怎么优化？",
      "heat": 98
    }
  ]
}
```

**对应前端功能**

[PostsListView.vue](src/views/PostsListView.vue) 右侧栏：

- 讨论看板
- 活跃作者
- 后续可扩展热门帖子

**理由**

现在右侧栏数据是从当前页帖子临时推导的。例如“活跃作者”只代表当前页，而不是全站或近期活跃作者。侧栏接口可以让这些信息变成真实社区指标。

## 6. 提供更自然的发帖接口

**后端功能**

建议新增：

```text
POST /posts
```

请求体：

```json
{
  "title": "帖子标题",
  "content": "帖子内容"
}
```

用户身份从 token 中解析，不需要前端传 `user_id`。

**对应前端功能**

[PostsListView.vue](src/views/PostsListView.vue) 的“写帖子”入口，以及用户主页里的发帖表单。

**理由**

当前创建帖子接口是 `/users/posts/create`，并且前端创建时看起来会传 `user_id`。发帖属于当前登录用户的行为，后端从 token 获取用户更安全，也能避免伪造 `user_id`。

## 推荐实施优先级

1. 排序参数：直接影响“最新 / 热门 / 高浏览”的正确性。
2. 搜索参数：让搜索从当前页筛选升级为全站搜索。
3. 评论数量：让帖子卡片信息完整，且避免额外请求。
4. 作者信息补齐：提升社区列表的信息密度。
5. 侧栏聚合接口：让右侧看板和活跃作者更真实。
6. 新发帖接口：优化接口设计和安全性。
