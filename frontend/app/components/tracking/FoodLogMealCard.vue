<script setup lang="ts">
import type { MealResponse } from '~/types/plan.types'
import type { FoodLogEntry } from '~/types/tracking.types'

const props = defineProps<{
  meal: MealResponse
  log?: FoodLogEntry
}>()

const emit = defineEmits<{
  log: [optionId: string]
  skip: []
}>()

const selectedLabel = computed(() => {
  if (!props.log?.selected_option_id) return null
  const option = props.meal.options.find(item => item.id === props.log?.selected_option_id)
  return option?.label || `گزینه ${option?.option_number ?? ''}`.trim()
})
</script>

<template>
  <div class="rounded-2xl bg-white p-4 shadow-sm">
    <div class="mb-3 flex items-start justify-between gap-3">
      <div>
        <h2 class="text-sm font-semibold text-gray-800 text-start">{{ meal.title }}</h2>
        <p v-if="meal.scheduled_time" class="mt-1 text-xs text-gray-500 text-start">{{ meal.scheduled_time }}</p>
      </div>
      <span
        v-if="log"
        :class="[
          'rounded-full px-3 py-1 text-xs font-medium',
          log.is_skipped ? 'bg-rose-100 text-rose-700' : 'bg-emerald-100 text-emerald-700',
        ]"
      >
        {{ log.is_skipped ? 'رد شد' : (selectedLabel || 'ثبت شد') }}
      </span>
    </div>

    <div v-if="!log" class="space-y-2">
      <button
        v-for="option in meal.options"
        :key="option.id"
        type="button"
        class="w-full rounded-xl border border-gray-200 p-3 text-start transition hover:border-emerald-300 hover:bg-emerald-50"
        @click="emit('log', option.id)"
      >
        <p class="text-sm font-medium text-gray-800">
          {{ option.label || `گزینه ${option.option_number}` }}
        </p>
        <p class="mt-1 text-xs text-gray-500">
          {{ option.items.map(item => item.food.name).join('، ') || 'بدون آیتم' }}
        </p>
      </button>
      <button
        type="button"
        class="w-full rounded-xl border border-rose-200 bg-rose-50 p-3 text-sm font-medium text-rose-700"
        @click="emit('skip')"
      >
        رد کردن این وعده
      </button>
    </div>

    <p v-else class="text-xs text-gray-500 text-start">
      {{ log.notes || 'برای تغییر، ثبت جدید انجام دهید.' }}
    </p>
  </div>
</template>
