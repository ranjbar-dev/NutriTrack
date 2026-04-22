<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>

<template>
  <section class="state-card" aria-live="polite">
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>
  </section>
</template>

<style scoped>
.state-card {
  padding: var(--space-5);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
}

h2 {
  margin: 0 0 var(--space-2);
}

p {
  margin: 0;
  color: var(--color-muted);
}
</style>
