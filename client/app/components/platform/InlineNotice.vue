<script setup lang="ts">
defineProps<{
  tone?: 'info' | 'success' | 'warning' | 'error'
  text: string
}>()
</script>

<template>
  <p class="inline-notice" :data-tone="tone ?? 'info'">{{ text }}</p>
</template>

<style scoped>
.inline-notice {
  margin: 0;
  padding: var(--space-3);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: 0.875rem;
}
</style>
