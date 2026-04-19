<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <slot />
    <UiBottomNav :items="navItems" />
  </div>
</template>

<script setup lang="ts">
const navItems = [
  { label: 'کلاینت‌ها', icon: 'users', to: '/nutritionist/clients' },
  { label: 'غذاها', icon: 'food', to: '/nutritionist/foods' },
  { label: 'داروها', icon: 'medications', to: '/nutritionist/medications' },
  { label: 'پیام‌ها', icon: 'chat', to: '/nutritionist/messages' },
  { label: 'پروفایل', icon: 'user', to: '/nutritionist/profile' },
]
</script>
