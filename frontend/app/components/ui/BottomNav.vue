<template>
  <nav
    class="fixed bottom-0 start-0 end-0 mx-auto bg-white border-t border-gray-200 z-50"
    style="max-width: 430px"
  >
    <div
      class="flex items-center justify-around h-16 pb-[env(safe-area-inset-bottom)]"
    >
      <NuxtLink
        v-for="item in items"
        :key="item.to"
        :to="item.to"
        class="flex flex-col items-center justify-center flex-1 py-2 text-xs transition-colors duration-200"
        :class="[
          isActive(item.to)
            ? 'text-emerald-600 font-bold'
            : 'text-gray-500 hover:text-gray-700',
        ]"
      >
        <span class="text-xl mb-1">{{ getIcon(item.icon) }}</span>
        <span>{{ item.label }}</span>
      </NuxtLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
interface NavItem {
  label: string
  icon: string
  to: string
}

defineProps<{
  items: NavItem[]
}>()

const route = useRoute()

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(to + '/')
}

function getIcon(icon: string): string {
  const icons: Record<string, string> = {
    users: '👥',
    food: '🍽️',
    chart: '📊',
    chat: '💬',
    user: '👤',
    plan: '📋',
    tracking: '✏️',
    medications: '💊',
  }
  return icons[icon] || '📄'
}
</script>
