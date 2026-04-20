<script setup lang="ts">
import type { FoodEmbedded } from '~/types/plan.types'
import { usePlanBuilderStore } from '~/stores/planBuilder'

const emit = defineEmits<{
  select: [payload: { optionId: string, food_id: string, quantity: number, measurement_unit: string, notes?: string }]
}>()

type FoodSearchResponse = {
  data?: FoodEmbedded[]
}

const store = usePlanBuilderStore()
const { apiFetch } = useApi()

const searchInput = ref<HTMLInputElement | null>(null)
const selectedFood = ref<FoodEmbedded | null>(null)
const quantity = ref<number>(0)
const measurementUnit = ref<string>('gram')
const notes = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null

function resetSelection() {
  selectedFood.value = null
  quantity.value = 0
  measurementUnit.value = 'gram'
  notes.value = ''
}

function closeSheet() {
  store.closeFoodPicker()
  store.foodPicker.searchQuery = ''
  store.foodPicker.searchResults = []
  store.foodPicker.loading = false
  resetSelection()
}

async function searchFoods(query: string) {
  const trimmed = query.trim()
  if (!trimmed) {
    store.foodPicker.searchResults = []
    return
  }

  store.foodPicker.loading = true
  try {
    const data = await apiFetch<FoodSearchResponse>(`/foods?search=${encodeURIComponent(trimmed)}&limit=20`)
    store.foodPicker.searchResults = data.data ?? []
  }
  catch {
    store.foodPicker.searchResults = []
  }
  finally {
    store.foodPicker.loading = false
  }
}

function selectFood(food: FoodEmbedded) {
  selectedFood.value = food
  quantity.value = food.measurement_amount || 1
  measurementUnit.value = food.measurement_unit || 'gram'
  notes.value = ''
}

function submitSelection() {
  const optionId = store.foodPicker.targetOptionId
  if (!selectedFood.value || !optionId || quantity.value <= 0) {
    return
  }

  emit('select', {
    optionId,
    food_id: selectedFood.value.id,
    quantity: quantity.value,
    measurement_unit: measurementUnit.value,
    notes: notes.value || undefined,
  })
  resetSelection()
}

watch(() => store.foodPicker.searchQuery, (query) => {
  if (!store.foodPicker.open) {
    return
  }

  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    searchFoods(query)
  }, 300)
})

watch(() => store.foodPicker.open, async (open) => {
  if (!open) {
    return
  }

  resetSelection()
  await nextTick()
  searchInput.value?.focus()
})

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition name="slide-up">
      <div
        v-if="store.foodPicker.open"
        class="fixed inset-0 z-50 flex items-end"
        role="dialog"
        aria-modal="true"
        aria-label="انتخاب ماده غذایی"
      >
        <button
          type="button"
          class="absolute inset-0 bg-black/40"
          aria-label="بستن"
          @click="closeSheet"
        />

        <div class="relative z-10 w-full rounded-t-3xl bg-white px-4 pb-6 pt-4 shadow-2xl">
          <div class="mx-auto mb-4 h-1.5 w-16 rounded-full bg-gray-200" />
          <div class="flex items-center justify-between gap-3">
            <h2 class="text-base font-bold text-gray-800">انتخاب ماده غذایی</h2>
            <button type="button" class="text-sm text-gray-500" @click="closeSheet">بستن</button>
          </div>

          <input
            ref="searchInput"
            v-model="store.foodPicker.searchQuery"
            type="text"
            placeholder="جستجوی غذا..."
            class="mt-4 w-full rounded-2xl border border-gray-200 px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
          >

          <div class="mt-4 max-h-72 space-y-2 overflow-y-auto">
            <p v-if="store.foodPicker.loading" class="text-sm text-gray-400">
              در حال جستجو...
            </p>
            <button
              v-for="food in store.foodPicker.searchResults"
              :key="food.id"
              type="button"
              class="block w-full rounded-2xl border border-gray-100 px-4 py-3 text-start transition-colors hover:bg-gray-50"
              @click="selectFood(food)"
            >
              <p class="font-medium text-sm text-gray-800">{{ food.name }}</p>
              <p class="mt-1 text-xs text-gray-500">
                {{ food.calories }} کیلوکالری / {{ food.measurement_amount }} {{ food.measurement_unit }}
              </p>
            </button>

            <p
              v-if="!store.foodPicker.loading && store.foodPicker.searchQuery.trim() && !store.foodPicker.searchResults.length"
              class="text-sm text-gray-400"
            >
              نتیجه‌ای یافت نشد
            </p>
          </div>

          <div v-if="selectedFood" class="mt-4 rounded-2xl border border-emerald-100 bg-emerald-50 p-4">
            <p class="font-semibold text-sm text-emerald-800">{{ selectedFood.name }}</p>
            <div class="mt-3 space-y-3">
              <div>
                <label class="mb-1 block text-xs text-gray-600">مقدار</label>
                <input
                  v-model.number="quantity"
                  type="number"
                  min="0.1"
                  step="0.1"
                  class="w-full rounded-xl border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
                >
              </div>
              <div>
                <label class="mb-1 block text-xs text-gray-600">واحد</label>
                <input
                  v-model="measurementUnit"
                  type="text"
                  class="w-full rounded-xl border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
                >
              </div>
              <div>
                <label class="mb-1 block text-xs text-gray-600">یادداشت</label>
                <textarea
                  v-model="notes"
                  rows="2"
                  class="w-full rounded-xl border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>

            <button
              type="button"
              class="mt-4 w-full rounded-2xl bg-emerald-600 py-3 text-sm font-medium text-white transition-colors hover:bg-emerald-700"
              @click="submitSelection"
            >
              افزودن به گزینه
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
}
</style>
