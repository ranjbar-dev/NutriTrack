<script setup lang="ts">
import type { NutritionistClient } from '~/types/nutritionist-workspace'

defineProps<{
  client: NutritionistClient | null
}>()
</script>

<template>
  <section class="header-card">
    <p v-if="!client">در حال بارگذاری اطلاعات مرجع...</p>
    <template v-else>
      <h1>{{ client.full_name }}</h1>
      <p>{{ client.mobile }}</p>
      <div class="tags">
        <span class="tag" :class="client.is_active ? 'ok' : 'off'">{{ client.is_active ? 'فعال' : 'غیرفعال' }}</span>
        <span v-if="client.bmi !== null" class="tag">BMI: {{ client.bmi.toFixed(1) }}</span>
        <span v-if="client.weight !== null" class="tag">وزن: {{ client.weight }}kg</span>
      </div>
    </template>
  </section>
</template>

<style scoped>
.header-card {
  border: 1px solid #d4dce0;
  border-radius: 12px;
  background: #fff;
  padding: 14px;
}

h1 {
  margin: 0;
  font-size: 1.1rem;
}

p {
  margin: 6px 0;
  color: #58666f;
}

.tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.tag {
  background: #eff4f7;
  color: #32434f;
  border-radius: 12px;
  padding: 2px 8px;
  font-size: 0.78rem;
}

.tag.ok {
  background: #dff5e3;
  color: #1f6c2e;
}

.tag.off {
  background: #f8e3e3;
  color: #8b2121;
}
</style>
