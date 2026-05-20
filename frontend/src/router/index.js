import { createRouter, createWebHashHistory } from 'vue-router'

import PostsView from '../views/PostsListView.vue'
import UsersView from '../views/UsersProfileView.vue'
import NotFound from '@/components/NotFound.vue'
import UserList from '@/views/UserListView.vue'
import Test from '@/components/test/Test.vue'
import Chat from '@/views/UserChatView.vue'
import UserProfile from '@/components/setting/UserProfile.vue'
import PostDetail from '@/components/post/PostDetail.vue'
import Problem from '@/views/ProblemListView.vue'
import ProblemDetail from '@/views/ProblemView.vue'
import ProblemUpload from '@/views/ProblemUpload.vue'
import SubmissionDetail from '@/components/coding/SubmissionDetail.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: PostsView
  },
  {
    path: '/posts',
    name: 'Posts',
    component: PostsView
  },
  {
    path: '/users/userList',
    name: 'UserList',
    component: UserList
  },
  {
    path: '/users/chat',
    name: 'Chat',
    component: Chat,
    meta: { requiresAuth: true }
  },
  {
    path: '/users/:userId',
    name: 'UserInfo',
    component: UsersView
  },
  {
    path: '/user/profile',
    name: 'UserProfile',
    component: UserProfile
  },
  {
    path: '/posts/:post_id',
    name: 'PostDetail',
    component: PostDetail,
    props: true
  },
  {
    path: '/problems',
    name: 'Problem',
    component: Problem
  },
  {
    path: '/problems/:problem_id',
    name: 'ProblemDetail',
    component: ProblemDetail
  },
  {
    path: '/problems/submissions/:submission_id',
    name: 'SubmissionDetail',
    component: SubmissionDetail
  },
  {
    path: '/upload/problem',
    name: 'ProblemUpload',
    component: ProblemUpload
  },
  {
    path: '/test',
    name: 'test',
    component: Test,
    meta: { requiresAuth: true }
  },
  {
    path: '/404',
    name: '404',
    component: NotFound
  },
  { path: '/:pathMatch(.*)*', redirect: '/404' }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
