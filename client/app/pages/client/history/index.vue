<script setup lang="ts">
definePageMeta({ layout: 'client' })
</script>

<template>
  <div class="history-index">
    <h1>تاریخچه</h1>
    <NuxtLink to="/client/history/tracking">تاریخچه ثبت ها و همگام سازی</NuxtLink>
    <NuxtLink to="/client/history/plans">برنامه های قبلی</NuxtLink>
  </div>
</template>

<style scoped>
.history-index {
  padding: var(--spacing-lg);
  display: grid;
  gap: var(--spacing-sm);
  direction: rtl;
}

a {
  display: block;
  text-decoration: none;
  color: #0f6b7a;
  background: #f4fafd;
  border: 1px solid #d2e1e7;
  padding: var(--spacing-md);
  border-radius: 10px;
}
</style>
