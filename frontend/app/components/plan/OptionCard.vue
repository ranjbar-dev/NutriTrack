<script setup lang="ts">
import type { MealOptionResponse } from '~/types/plan.types'
import { useNutritionComputed } from '~/composables/useNutritionComputed'

const props = defineProps<{
  option: MealOptionResponse
  editable?: boolean
}>()

const emit = defineEmits<{
  'add-item': [optionId: string]
  'delete-item': [{ optionId: string, itemId: string }]
  'delete-option': [optionId: string]
}>()

const expanded = ref(true)
const { optionTotals } = useNutritionComputed()
const totals = computed(() => optionTotals(props.option))
const panelId = computed(() => `option-panel-${props.option.id}`)
</script>

<template>
  <div class="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm">
    <button
      type="button"
      class="flex w-full items-center justify-between gap-3 px-4 py-4 text-start"
      :aria-expanded="expanded"
      :aria-controls="panelId"
      @click="expanded = !expanded"
    >
      <div>
        <p class="font-semibold text-sm text-gray-800">
          گزینه {{ option.option_number }}
          <span v-if="option.label" class="ms-1 text-gray-500">{{ option.label }}</span>
        </p>
        <div class="mt-2">
          <PlanNutritionBadges :totals="totals" compact />
        </div>
      </div>
      <span class="text-gray-300 transition-transform" :class="expanded ? 'rotate-180' : ''">⌄</span>
    </button>

    <div v-show="expanded" :id="panelId" class="border-t border-gray-100 px-4 py-4">
      <div v-if="option.items.length" class="space-y-3">
        <PlanFoodItemRow
          v-for="item in option.items"
          :key="item.id"
          :item="item"
          :editable="editable"
          @delete="emit('delete-item', { optionId: option.id, itemId: $event })"
        />
      </div>
      <p v-else class="text-sm text-gray-400">
        هنوز ماده غذایی اضافه نشده است
      </p>

      <div v-if="editable" class="mt-4 flex items-center gap-3">
        <button
          type="button"
          class="rounded-xl bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-700"
          @click="emit('add-item', option.id)"
        >
          افزودن ماده غذایی
        </button>
        <button
          type="button"
          class="text-sm text-red-400 transition-colors hover:text-red-600"
          @click="emit('delete-option', option.id)"
        >
          حذف گزینه
        </button>
      </div>
    </div>
  </div>
</template>
