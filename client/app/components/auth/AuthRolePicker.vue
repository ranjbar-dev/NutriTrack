<script setup lang="ts">
export interface AuthRoleEntry {
  id: 'client' | 'nutritionist' | 'admin'
  title: string
  description: string
  href: string
}

defineProps<{
  entries: AuthRoleEntry[]
}>()
</script>

<template>
  <nav class="role-picker" aria-label="انتخاب مسیر ورود">
    <NuxtLink
      v-for="entry in entries"
      :key="entry.id"
      class="role-card"
      :to="entry.href"
    >
      <h3>{{ entry.title }}</h3>
      <p>{{ entry.description }}</p>
    </NuxtLink>
  </nav>
</template>

<style scoped>
.role-picker {
  display: grid;
  gap: var(--space-3);
}

.role-card {
  display: grid;
  gap: var(--space-1);
  padding: var(--space-4);
  border-radius: 12px;
  border: 1px solid #c5d2da;
  background: #ffffff;
  color: inherit;
  text-decoration: none;
}

.role-card h3 {
  margin: 0;
}

.role-card p {
  margin: 0;
  color: #476272;
}

.role-card:focus-visible {
  outline: 2px solid #0f6b7a;
  outline-offset: 2px;
}
</style>
