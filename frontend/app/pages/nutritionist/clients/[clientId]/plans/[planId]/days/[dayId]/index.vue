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
const store = usePlanBuilderStore()
const { dayTotals } = useNutritionComputed()

const showExercises = ref(true)

const day = computed(() =>
  store.currentPlan?.days.find(item => item.id === dayId) ?? null,
)

const meals = computed(() =>
  [...(day.value?.meals ?? [])].sort((a, b) => a.display_order - b.display_order),
)

async function moveMealUp(mealId: string) {
  const index = meals.value.findIndex(meal => meal.id === mealId)
  if (index <= 0) return
  const current = meals.value[index]
  if (!current) return
  await store.reorderMeal(planId, dayId, mealId, current.display_order - 1)
}

async function moveMealDown(mealId: string) {
  const index = meals.value.findIndex(meal => meal.id === mealId)
  if (index === -1 || index >= meals.value.length - 1) return
  const current = meals.value[index]
  if (!current) return
  await store.reorderMeal(planId, dayId, mealId, current.display_order + 1)
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
          { label: 'جزئیات برنامه', to: `/nutritionist/clients/${clientId}/plans/${planId}` },
          { label: day ? `روز ${day.day_number}` : 'روز' },
        ]"
      />
    </div>

    <div v-if="store.loading" class="space-y-3 p-4">
      <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-white shadow-sm" />
    </div>

    <div v-else-if="day" class="space-y-4 p-4">
      <section class="rounded-2xl bg-white p-4 shadow-sm">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h1 class="font-bold text-gray-800">روز {{ day.day_number }}</h1>
            <p v-if="day.label" class="mt-1 text-sm text-gray-500">{{ day.label }}</p>
          </div>
          <PlanNutritionBadges :totals="dayTotals(day)" compact />
        </div>
      </section>

      <section class="space-y-3">
        <div class="flex items-center justify-between gap-3">
          <h2 class="font-bold text-gray-800">وعده‌ها</h2>
          <NuxtLink
            :to="`/nutritionist/clients/${clientId}/plans/${planId}/days/${dayId}/meals/new`"
            class="rounded-xl bg-emerald-600 px-4 py-2 text-sm font-medium text-white"
          >
            وعده جدید
          </NuxtLink>
        </div>

        <PlanMealCard
          v-for="(meal, index) in meals"
          :key="meal.id"
          :meal="meal"
          :plan-id="planId"
          :client-id="clientId"
          :day-id="dayId"
          :editable="store.currentPlan?.status === 'draft'"
          :is-first="index === 0"
          :is-last="index === meals.length - 1"
          @delete="store.deleteMeal(planId, dayId, $event)"
          @move-up="moveMealUp"
          @move-down="moveMealDown"
        />
      </section>

      <section class="rounded-2xl bg-white p-4 shadow-sm">
        <button
          type="button"
          class="flex w-full items-center justify-between gap-3"
          @click="showExercises = !showExercises"
        >
          <h2 class="font-bold text-gray-800">تمرین‌ها</h2>
          <span class="text-gray-300 transition-transform" :class="showExercises ? 'rotate-180' : ''">⌄</span>
        </button>

        <div v-show="showExercises" class="mt-4 space-y-2">
          <PlanExerciseCard
            v-for="exercise in day.exercises"
            :key="exercise.id"
            :exercise="exercise"
            :editable="store.currentPlan?.status === 'draft'"
            @delete="store.deleteExercise(planId, dayId, $event)"
          />
          <p v-if="!day.exercises.length" class="text-sm text-gray-400">
            تمرینی برای این روز ثبت نشده است
          </p>
        </div>
      </section>
    </div>

    <p v-else class="py-16 text-center text-sm text-gray-400">
      روز موردنظر یافت نشد
    </p>
  </div>
</template>
