<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <slot />
    <UiBottomNav :items="navItems" />
  </div>
</template>

<script setup lang="ts">
const navItems = [
  { label: 'برنامه', icon: 'plan', to: '/client/plan' },
  { label: 'ثبت', icon: 'tracking', to: '/client/tracking' },
  { label: 'پیام‌ها', icon: 'chat', to: '/client/messages' },
  { label: 'پروفایل', icon: 'user', to: '/client/profile' },
]
</script>
