<script setup lang="ts"></script>

<template>
  <header class="site-header">
    <span class="site-header__title">VB Events</span>
  </header>
  <main>
    <RouterView />
  </main>
</template>

<style scoped>
.site-header {
  display: flex;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  background: #fff;
}

.site-header__title {
  font-weight: 700;
  font-size: 1.25rem;
  color: #0f766e;
}
</style>
