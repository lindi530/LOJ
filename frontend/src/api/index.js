import request from '@/utils/request';

export default {
  login(payload) {
    return request.post('/users/login', payload);
  },
  logout(formData) {
    return request.post('/users/logout', formData)
  },
  upLoadAvatar(userId, formData) {
    return request.post(`/users/${userId}/modify_avatar`, formData);
  },
  modifyUserProfile(userId, formData) {
    return request.patch(`/users/${userId}/profile`, formData);
  },
  getUserProfileInfo(userId) {
    return request.get(`/users/${userId}`)
  },
  getPosts(page, pageSize) {
    return request.get(`/posts?page=${page}&page_size=${pageSize}`)
  },
  getUserList() {
    return request.get('users/userlist')
  },
  getPostsByUserId(userId, page, pageSize) {
    return request.get(`/users/${userId}/posts?page=${page}&page_size=${pageSize}`);
  },
  createPost(postdata) {
    return request.post(`/users/posts/create`, postdata);
  },
  register(postdata) {
    return request.post('/users/register', postdata)
  },
  deletePost(userId, postId) {
    return request.delete(`/users/${userId}/posts/${postId}`)
  },
  checkFollowing(targetUserId) {
    return request.get(`/users/${targetUserId}/is_following`)
  },
  followUser(targetUserId) {
    return request.post(`/users/${targetUserId}/follow`)
  },
  unFollowUser(targetUserId) {
    return request.delete(`/users/${targetUserId}/follow`)
  },
  getPostByPostId(postId) {
    return request.get(`/posts/${postId}`)
  },
  getPostComments(postId) {
    return request.get(`/posts/${postId}/comments`)
  },
  createPostComment(postId, newPost) {
    return request.post(`/posts/${postId}/comments`, newPost)
  },
  likeComment(commentId) { 
    return request.post(`/comment/${commentId}/like`)
  },
  unlikeComment(commentId) {
    return request.post(`/comment/${commentId}/unlike`)
  },
  deleteComment(commentId) {
    return request.delete(`/posts/comments/${commentId}`)
  },
  validateAccessToken() {
    return request.get(`/auth/validate/access_token`)
  },
  validateRefreshToken(formData) {
    return request.post(`/auth/validate/refresh_token`, formData)
  },
  likePost(postId) {
    return request.post(`/posts/${postId}/like`)
  },
  unLikePost(postId) {
    return request.post(`/posts/${postId}/unlike`)
  },
  onlineState() {
    return request.post(`/users/online`)
  },
  getMessageByTargetId(targetId) { 
    return request.get(`/message/${targetId}`)
  },
  getFollowUserList() {
    return request.get(`/follow/userlist`)
  },
  getProblemDetail(problemID) {
    return request.get(`/problems/${problemID}`)
  },
  getProblemList() {
    return request.get(`/problems`)
  },
  submitCode(problemID, data) { 
    return request.post(`/problems/${problemID}/submit`, data)
  },
  submitExample(problemID, data) { 
    return request.post(`/problems/${problemID}/submit/example`, data)
  },
  match() {
    return request.post(`/match`)
  },
  cancelMatch() {
    return request.post(`/match/cancel`)
  },
  getUserSaberStats() { 
    return request.get(`/saber/info`)
  },
  saberSubmit(data) {
    return request.post(`/saber/submit`, data)
  },
  inviteFriend() {
    return request.post(`/saber/invite`)
  },
  expireRoom(roomID) {
    return request.post(`/saber/rooms/expire/${roomID}`)
  },
  joinRoom(roomID) {
    return request.post(`/saber/rooms/join/${roomID}`)
  },
  uploadProblem(data) {
    return request.post(`/upload/problem`, data)
  },
  getUserSubmissionList(userId, page, pageSize) {
    return request.get(`/users/${userId}/submissions?page=${page}&page_size=${pageSize}`);
  },
  getProblemSubmissionList(problemID) {
    return request.get(`/problems/${problemID}/submissions`)
  },
  getSubmissionDetail(submissionID) {
    return request.get(`/problems/submissions/${submissionID}`)
  }
  // 这里按需继续扩展接口
};