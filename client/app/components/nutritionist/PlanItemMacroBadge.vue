<script setup lang="ts">
import type { NutritionTotals } from '~/types/diet-authoring'

defineProps<{
  totals: NutritionTotals | null
}>()
</script>

<template>
  <div class="badge" v-if="totals">
    <span>{{ totals.calories }} kcal</span>
    <span>P {{ totals.protein }}</span>
    <span>C {{ totals.carbs }}</span>
    <span>F {{ totals.fat }}</span>
  </div>
</template>

<style scoped>
.badge {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  border: 1px dashed #9bb4c2;
  border-radius: 10px;
  padding: 6px;
  font-size: 0.75rem;
  background: #f7fbfd;
}
</style>
