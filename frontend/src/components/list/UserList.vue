<template>
  <n-space vertical class="container content-box">
    <n-card title="用户列表" size="medium">
      <input
        v-model="searchQuery"
        type="text"
        class="form-control mb-3"
        placeholder="搜用户名"
      />
      <ul class="list-group">
        <li
          v-for="(user, index) in filteredUsers"
          :key="user.user_id"
          :class="[
            'list-group-item d-flex justify-content-between align-items-center',
            index % 2 === 1 ? 'color-odd' : 'color-even'
          ]"
        >
          <router-link
            :to="`/users/${user.user_id}`"
            class="d-flex align-items-center text-decoration-none text-dark"
          >
            <img :src="user.avatar" class="rounded-circle me-3" alt="avatar" width="48" height="48" />
            <div class="d-flex align-items-center">
              <div class="fw-bold text-primary me-2">{{ user.user_name }}</div>
              <div class="d-flex align-items-center">
                <i
                  :class="[
                    'bi',
                    user.gender === 'male' ? 'bi-person-fill text-primary' : 'bi-person-fill text-danger'
                  ]"
                ></i>
              </div>
            </div>
          </router-link>

          <div
            :class="user.online_state ? 'text-success' : 'text-muted'"
            class="d-flex align-items-center"
          >
            {{ user.online_state ? '在线' : '离线' }}
            <i
              :class="user.online_state ? 'bi bi-check-circle-fill' : 'bi bi-x-circle'"
              class="ms-1"
            ></i>
          </div>
        </li>
      </ul>
    </n-card>
  </n-space>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import api from '@/api';

const users = ref([]);
const searchQuery = ref('');

onMounted(async () => {
  const resp = await api.getUserList();
  if (resp.code === 0) {
    users.value = resp.data.userList;
  }
});

const filteredUsers = computed(() => {
  const q = searchQuery.value.toLowerCase();
  return users.value
    .filter(user => user.user_name?.toLowerCase().includes(q))
    .sort((a, b) => (b.online_state === true ? 1 : 0) - (a.online_state === true ? 1 : 0));
});

</script>

<style>
.page-wrapper {
  display: flex;
  justify-content: center;
  padding: 2rem;
}

.content-box {
  width: 100%;
  max-width: 900px;
}

.color-odd {
  background-color: #d4f2e7;
}

.color-even {
  background-color: #d3e0f3;
}

.bi {
  font-size: 1.2rem;
}
</style>
