<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>

<template>
  <section class="state-card error" role="alert">
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>
  </section>
</template>

<style scoped>
.state-card {
  padding: var(--space-5);
  border-radius: var(--radius-lg);
  background: #fff;
}

.error {
  border: 1px solid var(--color-danger);
}

h2 {
  margin: 0 0 var(--space-2);
  color: var(--color-danger);
}

p {
  margin: 0;
  color: var(--color-muted);
}
</style>
