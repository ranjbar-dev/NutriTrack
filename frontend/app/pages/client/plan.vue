<script setup lang="ts">
import jalaali from 'jalaali-js'
import { useNutritionComputed } from '~/composables/useNutritionComputed'
import { useClientPlanStore } from '~/stores/clientPlan'

definePageMeta({
  layout: 'client',
  middleware: ['role-guard'],
  roles: ['client'],
})

const store = useClientPlanStore()
const { optionTotals } = useNutritionComputed()
const activeTab = ref<'active' | 'history'>('active')
const pageLoading = ref(true)

function toShamsi(isoDate?: string) {
  if (!isoDate) return '—'
  const parts = isoDate.split('-').map(Number)
  if (parts.length !== 3) return isoDate
  const [y, m, d] = parts
  const j = jalaali.toJalaali(y!, m!, d!)
  return `${j.jy}/${String(j.jm).padStart(2, '0')}/${String(j.jd).padStart(2, '0')}`
}

const archivedPlans = computed(() =>
  store.myPlans.filter(plan => plan.status === 'archived'),
)

onMounted(async () => {
  pageLoading.value = true
  await Promise.all([
    store.fetchActivePlan(),
    store.fetchMyPlans(),
  ])
  pageLoading.value = false
})

onUnmounted(() => store.$reset())
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <div class="sticky top-0 z-20 border-b border-gray-100 bg-white px-4 py-4">
      <h1 class="text-lg font-bold text-gray-800">برنامه من</h1>
      <div class="mt-4 grid grid-cols-2 gap-2 rounded-2xl bg-gray-100 p-1">
        <button
          type="button"
          class="rounded-xl px-4 py-2 text-sm font-medium transition-colors"
          :class="activeTab === 'active' ? 'bg-white text-emerald-700 shadow-sm' : 'text-gray-500'"
          @click="activeTab = 'active'"
        >
          برنامه فعال
        </button>
        <button
          type="button"
          class="rounded-xl px-4 py-2 text-sm font-medium transition-colors"
          :class="activeTab === 'history' ? 'bg-white text-emerald-700 shadow-sm' : 'text-gray-500'"
          @click="activeTab = 'history'"
        >
          تاریخچه
        </button>
      </div>
    </div>

    <div v-if="pageLoading" class="space-y-3 p-4">
      <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-white shadow-sm" />
    </div>

    <div v-else-if="activeTab === 'active'">
      <div v-if="store.activePlan" class="space-y-4 pb-4">
        <div class="bg-white px-4 py-4 shadow-sm">
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-sm text-gray-500">
                از {{ toShamsi(store.activePlan.start_date) }}
                <span v-if="store.activePlan.end_date"> تا {{ toShamsi(store.activePlan.end_date) }}</span>
              </p>
              <p v-if="store.activePlan.notes" class="mt-2 text-sm text-gray-600">
                {{ store.activePlan.notes }}
              </p>
            </div>
            <PlanWaterBadge v-if="store.activePlan.daily_water_target_ml" :ml="store.activePlan.daily_water_target_ml" />
          </div>
        </div>

        <div class="sticky top-[86px] z-10 bg-gray-50 px-4 py-2">
          <PlanDayTabBar
            :days="store.activePlan.days"
            :active-day-number="store.activeDayNumber"
            @update:active-day-number="store.setActiveDay"
          />
        </div>

        <div v-if="store.activeDay" class="space-y-4 px-4">
          <section class="rounded-2xl bg-white p-4 shadow-sm">
            <div class="mb-4 flex items-center justify-between gap-2">
              <div>
                <h2 class="font-bold text-gray-800">روز {{ store.activeDay.day_number }}</h2>
                <p v-if="store.activeDay.label" class="text-sm text-gray-500">{{ store.activeDay.label }}</p>
              </div>
            </div>

            <div class="space-y-4">
              <div
                v-for="meal in store.activeDay.meals"
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

          <section
            v-if="store.activeDay.exercises.length"
            class="rounded-2xl bg-white p-4 shadow-sm"
          >
            <h2 class="mb-3 font-bold text-gray-800">تمرین‌ها</h2>
            <div class="space-y-2">
              <PlanExerciseCard
                v-for="exercise in store.activeDay.exercises"
                :key="exercise.id"
                :exercise="exercise"
              />
            </div>
          </section>

          <section
            v-if="store.activePlan.medications.length"
            class="rounded-2xl bg-white p-4 shadow-sm"
          >
            <h2 class="mb-3 font-bold text-gray-800">داروها</h2>
            <div class="space-y-3">
              <PlanMedicationCard
                v-for="medication in store.activePlan.medications"
                :key="medication.id"
                :medication="medication"
              />
            </div>
          </section>
        </div>
      </div>

      <div v-else class="px-4 py-16 text-center">
        <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-2xl">
          🍽️
        </div>
        <h2 class="mt-4 font-bold text-gray-800">برنامه‌ای فعال ندارید</h2>
        <p class="mt-2 text-sm text-gray-500">برای دریافت برنامه جدید با کارشناس تغذیه خود در تماس باشید.</p>
      </div>
    </div>

    <div v-else class="space-y-3 p-4">
      <NuxtLink
        v-for="plan in archivedPlans"
        :key="plan.id"
        :to="`/client/plans/${plan.id}`"
        class="block rounded-2xl bg-white p-4 shadow-sm"
      >
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="font-semibold text-sm text-gray-800">
              {{ toShamsi(plan.start_date) }}
              <span v-if="plan.end_date"> تا {{ toShamsi(plan.end_date) }}</span>
            </p>
            <p class="mt-1 text-xs text-gray-500">{{ plan.day_count }} روز</p>
          </div>
          <PlanStatusBadge :status="plan.status" />
        </div>
      </NuxtLink>

      <p v-if="!archivedPlans.length" class="py-12 text-center text-sm text-gray-400">
        تاریخچه‌ای برای نمایش وجود ندارد
      </p>
    </div>
  </div>
</template>
