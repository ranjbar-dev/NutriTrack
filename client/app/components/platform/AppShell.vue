<script setup lang="ts">
import AppHeader from './AppHeader.vue'
import BottomNavClient from './BottomNavClient.vue'

export type PlatformRole = 'auth' | 'client' | 'nutritionist' | 'admin'

const props = defineProps<{
  role: PlatformRole
  title: string
  subtitle?: string
}>()

const showClientNav = computed(() => props.role === 'client')
</script>

<template>
  <section class="app-shell" :data-role="role">
    <AppHeader :title="title" :subtitle="subtitle" />

    <main class="shell-content">
      <slot />
    </main>

    <BottomNavClient v-if="showClientNav" />
  </section>
</template>

<style scoped>
.app-shell {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
}

.shell-content {
  flex: 1;
  padding: var(--space-4);
}
</style>
