<script setup lang="ts">
import type { DietPlanFlat } from '~/types/diet-authoring'

defineProps<{
  activePlan: DietPlanFlat | null
  archivedCount: number
  latestMessageAt: string | null
  latestLabAt: string | null
}>()
</script>

<template>
  <section class="grid">
    <article class="panel">
      <h3>برنامه فعال</h3>
      <p v-if="activePlan">{{ activePlan.title || 'بدون عنوان' }}</p>
      <p v-else>برنامه فعالی ثبت نشده است.</p>
    </article>

    <article class="panel">
      <h3>آرشیو برنامه</h3>
      <p>{{ archivedCount }} برنامه</p>
    </article>

    <article class="panel">
      <h3>آخرین پیام</h3>
      <p>{{ latestMessageAt || 'بدون پیام اخیر' }}</p>
    </article>

    <article class="panel">
      <h3>آخرین آزمایش</h3>
      <p>{{ latestLabAt || 'بدون ثبت آزمایش' }}</p>
    </article>
  </section>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.panel {
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 10px;
}

h3 {
  margin: 0 0 6px;
  font-size: 0.86rem;
}

p {
  margin: 0;
  color: #52616b;
  font-size: 0.86rem;
}
</style>
