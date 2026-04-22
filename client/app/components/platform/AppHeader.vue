<script setup lang="ts">
import { useRtl } from '../../composables/useRtl'

defineProps<{
  title: string
  subtitle?: string
}>()

const { dir } = useRtl()
</script>

<template>
  <header class="app-header" :dir="dir">
    <h1>{{ title }}</h1>
    <p v-if="subtitle">{{ subtitle }}</p>
  </header>
</template>

<style scoped>
.app-header {
  padding: var(--space-4);
  background: var(--color-surface);
  border-bottom: 1px solid #dde6ea;
}

h1 {
  margin: 0;
  font-size: 1rem;
}

p {
  margin: var(--space-1) 0 0;
  color: var(--color-muted);
  font-size: 0.875rem;
}
</style>
