<template>
  <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-4">
    <div class="flex items-start justify-between gap-3">
      <div class="flex-1">
        <h3 class="text-lg font-bold text-gray-800">
          {{ food.name }}
        </h3>
      </div>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="text-gray-500 hover:text-emerald-600 transition-colors"
          aria-label="ویرایش غذا"
          @click="emit('edit', food.id)"
        >
          ✏️
        </button>
        <button
          type="button"
          class="text-gray-500 hover:text-red-600 transition-colors"
          aria-label="حذف غذا"
          @click="emit('delete', food.id)"
        >
          🗑️
        </button>
      </div>
    </div>

    <div class="flex flex-wrap gap-2 mt-3">
      <span
        v-for="category in food.categories"
        :key="category"
        class="px-2.5 py-1 text-xs rounded-full font-medium"
        :class="getCategoryClass(category)"
      >
        {{ getCategoryLabel(category) }}
      </span>
    </div>

    <p class="text-sm text-gray-600 mt-3">
      {{ toPersianDigits(food.calories) }} کالری | {{
        toPersianDigits(food.protein_g)
      }}
      گرم پروتئین
    </p>
  </div>
</template>

<script setup lang="ts">
import { toPersianDigits } from '~/utils/persian-digits'
import type { FoodResponse } from '~/stores/food'

defineProps<{
  food: FoodResponse
}>()

const emit = defineEmits<{
  edit: [id: string]
  delete: [id: string]
}>()

const categoryLabels: Record<string, string> = {
  breakfast: 'صبحانه',
  lunch: 'ناهار',
  dinner: 'شام',
  snack: 'میان‌وعده',
  fruit: 'میوه',
  beverage: 'نوشیدنی',
  supplement: 'مکمل',
  other: 'سایر',
}

const categoryClasses: Record<string, string> = {
  breakfast: 'bg-amber-100 text-amber-700',
  lunch: 'bg-emerald-100 text-emerald-700',
  dinner: 'bg-indigo-100 text-indigo-700',
  snack: 'bg-orange-100 text-orange-700',
  fruit: 'bg-rose-100 text-rose-700',
  beverage: 'bg-sky-100 text-sky-700',
  supplement: 'bg-purple-100 text-purple-700',
  other: 'bg-gray-100 text-gray-600',
}

function getCategoryLabel(category: string) {
  return categoryLabels[category] ?? category
}

function getCategoryClass(category: string) {
  return categoryClasses[category] ?? 'bg-gray-100 text-gray-600'
}
</script>
