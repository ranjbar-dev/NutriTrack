<script setup lang="ts">
import jalaali from 'jalaali-js'
import type { DietPlanResponse } from '~/types/plan.types'
import { useNutritionComputed } from '~/composables/useNutritionComputed'
import { useClientPlanStore } from '~/stores/clientPlan'

definePageMeta({
  layout: 'client',
  middleware: ['role-guard'],
  roles: ['client'],
})

const route = useRoute()
const planId = route.params.id as string
const store = useClientPlanStore()
const { optionTotals } = useNutritionComputed()

const plan = ref<DietPlanResponse | null>(null)
const activeDayNumber = ref(1)
const pageLoading = ref(true)

const activeDay = computed(() =>
  plan.value?.days.find(day => day.day_number === activeDayNumber.value) ?? null,
)

function toShamsi(isoDate?: string) {
  if (!isoDate) return '—'
  const parts = isoDate.split('-').map(Number)
  if (parts.length !== 3) return isoDate
  const [y, m, d] = parts
  const j = jalaali.toJalaali(y!, m!, d!)
  return `${j.jy}/${String(j.jm).padStart(2, '0')}/${String(j.jd).padStart(2, '0')}`
}

onMounted(async () => {
  pageLoading.value = true
  try {
    plan.value = await store.fetchPlanById(planId)
    activeDayNumber.value = plan.value?.days[0]?.day_number ?? 1
  }
  finally {
    pageLoading.value = false
  }
})

onUnmounted(() => store.$reset())
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <div class="sticky top-0 z-20 border-b border-gray-100 bg-white px-4 py-4">
      <PlanBreadcrumb
        :items="[
          { label: 'برنامه من', to: '/client/plan' },
          { label: 'آرشیو' },
        ]"
      />
    </div>

    <div v-if="pageLoading" class="space-y-3 p-4">
      <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-white shadow-sm" />
    </div>

    <div v-else-if="plan" class="space-y-4 pb-4">
      <div class="bg-white px-4 py-4 shadow-sm">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="text-sm text-gray-500">
              از {{ toShamsi(plan.start_date) }}
              <span v-if="plan.end_date"> تا {{ toShamsi(plan.end_date) }}</span>
            </p>
            <p v-if="plan.notes" class="mt-2 text-sm text-gray-600">{{ plan.notes }}</p>
          </div>
          <PlanStatusBadge :status="plan.status" />
        </div>
      </div>

      <div class="sticky top-[64px] z-10 bg-gray-50 px-4 py-2">
        <PlanDayTabBar
          :days="plan.days"
          :active-day-number="activeDayNumber"
          @update:active-day-number="activeDayNumber = $event"
        />
      </div>

      <div v-if="activeDay" class="space-y-4 px-4">
        <section class="rounded-2xl bg-white p-4 shadow-sm">
          <h2 class="font-bold text-gray-800">روز {{ activeDay.day_number }}</h2>
          <p v-if="activeDay.label" class="mt-1 text-sm text-gray-500">{{ activeDay.label }}</p>

          <div class="mt-4 space-y-4">
            <div
              v-for="meal in activeDay.meals"
              :key="meal.id"
              class="rounded-2xl border border-gray-100 p-4"
            >
              <div class="flex items-center justify-between gap-2">
                <h3 class="font-semibold text-sm text-gray-800">{{ meal.title }}</h3>
                <span v-if="meal.scheduled_time" class="rounded-full bg-gray-100 px-2 py-1 text-xs text-gray-500">
                  {{ meal.scheduled_time.slice(0, 5) }}
                </span>
              </div>

              <div class="mt-3 space-y-3">
                <div
                  v-for="option in meal.options"
                  :key="option.id"
                  class="rounded-xl bg-gray-50 p-3"
                >
                  <div class="flex items-center justify-between gap-2">
                    <p class="font-medium text-sm text-gray-700">گزینه {{ option.option_number }}</p>
                    <PlanNutritionBadges :totals="optionTotals(option)" compact />
                  </div>
                  <ul class="mt-3 space-y-2">
                    <li
                      v-for="item in option.items"
                      :key="item.id"
                      class="flex items-center justify-between gap-3 text-sm text-gray-600"
                    >
                      <span>{{ item.food.name }}</span>
                      <span class="shrink-0">{{ item.quantity }} {{ item.measurement_unit }}</span>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section v-if="activeDay.exercises.length" class="rounded-2xl bg-white p-4 shadow-sm">
          <h2 class="mb-3 font-bold text-gray-800">تمرین‌ها</h2>
          <div class="space-y-2">
            <PlanExerciseCard
              v-for="exercise in activeDay.exercises"
              :key="exercise.id"
              :exercise="exercise"
            />
          </div>
        </section>

        <section v-if="plan.medications.length" class="rounded-2xl bg-white p-4 shadow-sm">
          <h2 class="mb-3 font-bold text-gray-800">داروها</h2>
          <div class="space-y-3">
            <PlanMedicationCard
              v-for="medication in plan.medications"
              :key="medication.id"
              :medication="medication"
            />
          </div>
        </section>
      </div>
    </div>

    <p v-else class="py-16 text-center text-sm text-gray-400">
      برنامه موردنظر یافت نشد
    </p>
  </div>
</template>
