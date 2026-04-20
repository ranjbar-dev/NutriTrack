<script setup lang="ts">
import type { MealOptionItemResponse } from '~/types/plan.types'
import { useNutritionComputed } from '~/composables/useNutritionComputed'

const props = defineProps<{
  item: MealOptionItemResponse
  editable?: boolean
}>()

const emit = defineEmits<{
  delete: [id: string]
}>()

const { itemTotals } = useNutritionComputed()
const totals = computed(() => itemTotals(props.item))
</script>

<template>
  <div class="rounded-xl border border-gray-100 bg-gray-50 px-3 py-3">
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0 flex-1">
        <p class="font-medium text-sm text-gray-800">
          {{ item.food.name }}
        </p>
        <p class="mt-1 text-xs text-gray-500">
          {{ item.quantity }} {{ item.measurement_unit }}
        </p>
        <div class="mt-2">
          <PlanNutritionBadges :totals="totals" compact />
        </div>
        <p v-if="item.notes" class="mt-2 text-xs text-gray-400">
          {{ item.notes }}
        </p>
      </div>

      <button
        v-if="editable"
        type="button"
        class="shrink-0 text-xs text-red-400 transition-colors hover:text-red-600"
        @click="emit('delete', item.id)"
      >
        حذف
      </button>
    </div>
  </div>
</template>
