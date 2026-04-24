<script setup lang="ts">
definePageMeta({ layout: 'client' })

const links = [
  { to: '/client/tracking/food', label: 'ثبت غذا' },
  { to: '/client/tracking/water', label: 'ثبت آب' },
  { to: '/client/tracking/sleep', label: 'ثبت خواب' },
  { to: '/client/tracking/exercise', label: 'ثبت ورزش' },
  { to: '/client/tracking/medication', label: 'ثبت دارو' },
  { to: '/client/tracking/body', label: 'ثبت اندازه بدن' },
]
</script>

<template>
  <div class="tracking-index">
    <h1>ثبت های روزانه</h1>
    <div class="grid">
      <NuxtLink v-for="link in links" :key="link.to" :to="link.to">{{ link.label }}</NuxtLink>
    </div>
  </div>
</template>

<style scoped>
.tracking-index {
  padding: var(--spacing-lg);
  direction: rtl;
}

.grid {
  display: grid;
  gap: var(--spacing-sm);
}

a {
  display: block;
  text-decoration: none;
  padding: var(--spacing-md);
  border: 1px solid #d5e3ea;
  border-radius: 10px;
  color: #104952;
  background: #f8fcfe;
}
</style>
