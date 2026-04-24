<script setup lang="ts">
import type { AdminStats } from '~/types/admin'

defineProps<{
  stats: AdminStats | null
  loading?: boolean
}>()

const metricItems = [
  { key: 'total_nutritionists', label: 'کل متخصصان' },
  { key: 'active_nutritionists', label: 'متخصصان فعال' },
  { key: 'inactive_nutritionists', label: 'متخصصان غیرفعال' },
  { key: 'total_clients', label: 'کل مراجعان' },
  { key: 'total_foods', label: 'کل غذاها' },
  { key: 'active_diet_plans', label: 'برنامه های فعال' },
] as const
</script>

<template>
  <section class="grid">
    <p v-if="loading" class="placeholder">در حال دریافت آمار...</p>
    <p v-else-if="!stats" class="placeholder">آماری برای نمایش وجود ندارد.</p>

    <article v-for="item in metricItems" v-else :key="item.key" class="card">
      <span>{{ item.label }}</span>
      <strong>{{ stats[item.key] }}</strong>
    </article>
  </section>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.card,
.placeholder {
  border: 1px solid #d3dce2;
  border-radius: 14px;
  background: #fff;
  padding: 12px;
}

.card {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card span {
  color: #53606a;
  font-size: 0.92rem;
}

.card strong {
  font-size: 1.25rem;
}

.placeholder {
  grid-column: 1 / -1;
  margin: 0;
}
</style>