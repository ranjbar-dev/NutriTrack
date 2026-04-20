<script setup lang="ts">
import type { MealResponse } from '~/types/plan.types'
import { useNutritionComputed } from '~/composables/useNutritionComputed'

const props = defineProps<{
  meal: MealResponse
  planId: string
  clientId: string
  dayId: string
  editable?: boolean
  isFirst?: boolean
  isLast?: boolean
}>()

const emit = defineEmits<{
  delete: [id: string]
  'move-up': [id: string]
  'move-down': [id: string]
}>()

const { mealTotals } = useNutritionComputed()
const totals = computed(() => mealTotals(props.meal))
</script>

<template>
  <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
    <NuxtLink
      :to="`/nutritionist/clients/${clientId}/plans/${planId}/days/${dayId}/meals/${meal.id}`"
      class="block p-4 hover:bg-gray-50 transition-colors"
    >
      <div class="flex items-center justify-between gap-2">
        <span class="font-semibold text-gray-800 text-sm">{{ meal.title }}</span>
        <span
          v-if="meal.scheduled_time"
          class="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full shrink-0"
        >
          {{ meal.scheduled_time.slice(0, 5) }}
        </span>
      </div>
      <div class="mt-2">
        <PlanNutritionBadges :totals="totals" compact />
      </div>
      <div class="flex items-center justify-between mt-1">
        <span class="text-xs text-gray-400">{{ meal.options.length }} گزینه</span>
        <span class="text-gray-300 text-xs">‹</span>
      </div>
    </NuxtLink>
    <div v-if="editable" class="flex items-center gap-2 px-4 py-2 border-t border-gray-50 bg-gray-50">
      <button
        :disabled="isFirst"
        class="text-xs text-gray-500 hover:text-gray-800 disabled:opacity-30 disabled:cursor-not-allowed transition-colors px-2 py-1"
        aria-label="انتقال به بالا"
        @click.prevent="emit('move-up', meal.id)"
      >
        ↑
      </button>
      <button
        :disabled="isLast"
        class="text-xs text-gray-500 hover:text-gray-800 disabled:opacity-30 disabled:cursor-not-allowed transition-colors px-2 py-1"
        aria-label="انتقال به پایین"
        @click.prevent="emit('move-down', meal.id)"
      >
        ↓
      </button>
      <span class="flex-1" />
      <button
        class="text-xs text-red-400 hover:text-red-600 transition-colors"
        :aria-label="`حذف وعده ${meal.title}`"
        @click.prevent="emit('delete', meal.id)"
      >
        حذف
      </button>
    </div>
  </div>
</template>
