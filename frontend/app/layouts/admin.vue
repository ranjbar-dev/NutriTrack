<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <slot />
    <UiBottomNav :items="navItems" />
  </div>
</template>

<script setup lang="ts">
const navItems = [
  { label: 'متخصصین', icon: 'users', to: '/admin' },
  { label: 'غذاها', icon: 'food', to: '/admin/foods' },
  { label: 'آمار', icon: 'chart', to: '/admin/stats' },
]
</script>
