<script setup lang="ts">
import { useNutritionComputed } from '~/composables/useNutritionComputed'
import { usePlanBuilderStore } from '~/stores/planBuilder'

definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist', 'super_admin'],
})

const route = useRoute()
const clientId = route.params.clientId as string
const planId = route.params.planId as string
const dayId = route.params.dayId as string
const mealId = route.params.mealId as string
const store = usePlanBuilderStore()
const { mealTotals } = useNutritionComputed()

const day = computed(() =>
  store.currentPlan?.days.find(item => item.id === dayId) ?? null,
)

const meal = computed(() =>
  day.value?.meals.find(item => item.id === mealId) ?? null,
)

const options = computed(() =>
  [...(meal.value?.options ?? [])].sort((a, b) => a.option_number - b.option_number),
)

async function handleAddOption() {
  await store.addOption(planId, dayId, mealId)
}

function handleAddItem(optionId: string) {
  store.openFoodPicker(optionId)
}

async function handleSelectItem(payload: { optionId: string, food_id: string, quantity: number, measurement_unit: string, notes?: string }) {
  await store.addItem(planId, dayId, mealId, payload.optionId, {
    food_id: payload.food_id,
    quantity: payload.quantity,
    measurement_unit: payload.measurement_unit,
    notes: payload.notes,
  })
  store.closeFoodPicker()
}

async function handleDeleteItem(payload: { optionId: string, itemId: string }) {
  await store.deleteItem(planId, dayId, mealId, payload.optionId, payload.itemId)
}

async function handleDeleteOption(optionId: string) {
  await store.deleteOption(planId, dayId, mealId, optionId)
}

onMounted(() => store.loadPlan(planId))
onUnmounted(() => store.$reset())
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <div class="sticky top-0 z-10 border-b border-gray-100 bg-white px-4 py-3">
      <PlanBreadcrumb
        :items="[
          { label: 'برنامه‌ها', to: `/nutritionist/clients/${clientId}/plans` },
          { label: 'روز', to: `/nutritionist/clients/${clientId}/plans/${planId}/days/${dayId}` },
          { label: meal?.title ?? 'وعده' },
        ]"
      />
    </div>

    <div v-if="store.loading" class="space-y-3 p-4">
      <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-white shadow-sm" />
    </div>

    <div v-else-if="meal" class="space-y-4 p-4">
      <section class="rounded-2xl bg-white p-4 shadow-sm">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h1 class="font-bold text-gray-800">{{ meal.title }}</h1>
            <p v-if="meal.scheduled_time" class="mt-1 text-sm text-gray-500">
              ساعت {{ meal.scheduled_time.slice(0, 5) }}
            </p>
          </div>
          <PlanNutritionBadges :totals="mealTotals(meal)" compact />
        </div>
      </section>

      <section class="space-y-3">
        <div class="flex items-center justify-between gap-3">
          <h2 class="font-bold text-gray-800">گزینه‌ها</h2>
          <button
            v-if="store.currentPlan?.status === 'draft'"
            type="button"
            class="rounded-xl bg-emerald-600 px-4 py-2 text-sm font-medium text-white"
            @click="handleAddOption"
          >
            گزینه جدید
          </button>
        </div>

        <PlanOptionCard
          v-for="option in options"
          :key="option.id"
          :option="option"
          :editable="store.currentPlan?.status === 'draft'"
          @add-item="handleAddItem"
          @delete-item="handleDeleteItem"
          @delete-option="handleDeleteOption"
        />
      </section>

      <PlanFoodPickerSheet @select="handleSelectItem" />
    </div>

    <p v-else class="py-16 text-center text-sm text-gray-400">
      وعده موردنظر یافت نشد
    </p>
  </div>
</template>
