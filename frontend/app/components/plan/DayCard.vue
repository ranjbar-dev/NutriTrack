<script setup lang="ts">
import type { PlanDayResponse } from '~/types/plan.types'
import { useNutritionComputed } from '~/composables/useNutritionComputed'

const props = defineProps<{
  day: PlanDayResponse
  planId: string
  clientId: string
  editable?: boolean
}>()

const emit = defineEmits<{ delete: [id: string] }>()

const { dayTotals } = useNutritionComputed()
const totals = computed(() => dayTotals(props.day))
const dayLabel = computed(() =>
  props.day.label
    ? `روز ${props.day.day_number} — ${props.day.label}`
    : `روز ${props.day.day_number}`,
)
</script>

<template>
  <NuxtLink
    :to="`/nutritionist/clients/${clientId}/plans/${planId}/days/${day.id}`"
    class="block bg-white rounded-2xl p-4 shadow-sm hover:shadow-md transition-shadow"
  >
    <div class="flex items-center justify-between gap-3">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2">
          <span class="font-semibold text-gray-800 text-sm">{{ dayLabel }}</span>
          <span class="bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded-full">
            {{ day.meals.length }} وعده
          </span>
        </div>
        <div class="mt-2">
          <PlanNutritionBadges :totals="totals" compact />
        </div>
      </div>
      <div class="flex items-center gap-2 shrink-0">
        <button
          v-if="editable"
          class="text-red-400 hover:text-red-600 text-xs transition-colors"
          :aria-label="`حذف روز ${day.day_number}`"
          @click.prevent="emit('delete', day.id)"
        >
          حذف
        </button>
        <span class="text-gray-300">‹</span>
      </div>
    </div>
  </NuxtLink>
</template>
