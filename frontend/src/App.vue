<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

onMounted(() => {
  authStore.fetchUser()
})
</script>

<template>
  <header class="site-header">
    <span class="site-header__title">VB Events</span>
    <span v-if="authStore.user" class="site-header__user">{{ authStore.user.name }}</span>
  </header>
  <main>
    <RouterView />
  </main>
</template>

<style scoped>
.site-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  background: #fff;
}

.site-header__title {
  font-weight: 700;
  font-size: 1.25rem;
  color: #0f766e;
}

.site-header__user {
  font-weight: 600;
  color: #334155;
}
</style>
