<script setup lang="ts">
import type { NutritionTotals } from '~/types/plan.types'

defineProps<{ totals: NutritionTotals, compact?: boolean }>()
</script>

<template>
  <div :class="['flex items-center flex-wrap', compact ? 'gap-1' : 'gap-2']">
    <span :class="['font-semibold text-emerald-700', compact ? 'text-xs' : 'text-sm']">
      {{ totals.calories.toFixed(0) }} کیلوکالری
    </span>
    <span v-if="!compact" class="text-gray-300">|</span>
    <span :class="['text-gray-500', compact ? 'text-xs' : 'text-sm']">پ {{ totals.protein_g.toFixed(1) }}g</span>
    <span :class="['text-gray-500', compact ? 'text-xs' : 'text-sm']">ک {{ totals.carbs_g.toFixed(1) }}g</span>
    <span :class="['text-gray-500', compact ? 'text-xs' : 'text-sm']">چ {{ totals.fat_g.toFixed(1) }}g</span>
  </div>
</template>
