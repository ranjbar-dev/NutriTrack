<script setup lang="ts">
export interface ClientNavItem {
  key: 'today' | 'history' | 'messages' | 'profile'
  label: string
  to: string
}

const props = defineProps<{
  items?: ClientNavItem[]
}>()

const defaults: ClientNavItem[] = [
  { key: 'today', label: 'امروز', to: '/client' },
  { key: 'history', label: 'تاریخچه', to: '/client/history/tracking' },
  { key: 'messages', label: 'پیام‌ها', to: '/client/messages' },
  { key: 'profile', label: 'پروفایل', to: '/client/profile' }
]

const navItems = computed(() => props.items ?? defaults)
</script>

<template>
  <nav class="bottom-nav" aria-label="ناوبری کاربر">
    <NuxtLink v-for="item in navItems" :key="item.key" :to="item.to">{{ item.label }}</NuxtLink>
  </nav>
</template>

<style scoped>
.bottom-nav {
  position: sticky;
  bottom: 0;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4) calc(var(--space-3) + var(--safe-bottom));
  background: var(--color-surface);
  border-top: 1px solid #dde6ea;
}

a {
  text-decoration: none;
  text-align: center;
  color: var(--color-text);
  font-size: 0.75rem;
}
</style>
