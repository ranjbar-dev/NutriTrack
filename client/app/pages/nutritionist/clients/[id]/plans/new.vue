<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, navigateTo } from '#imports'
import PlanPeriodFormCard from '~/components/nutritionist/PlanPeriodFormCard.vue'
import PlanDayEditor from '~/components/nutritionist/PlanDayEditor.vue'
import MealEditor from '~/components/nutritionist/MealEditor.vue'
import OptionEditor from '~/components/nutritionist/OptionEditor.vue'
import ExercisePrescriptionEditor from '~/components/nutritionist/ExercisePrescriptionEditor.vue'
import { useDietPlanAuthoringApi } from '~/composables/useDietPlanAuthoringApi'
import type { PlanDay, PlanMeal, PlanOption, PlanExercise, PlanPrescription } from '~/types/diet-authoring'

definePageMeta({
  layout: 'nutritionist',
})

const route = useRoute()
const clientId = computed(() => String(route.params.id || ''))
const api = useDietPlanAuthoringApi()
const planId = ref('')
const loading = ref(false)
const error = ref('')

const days = ref<PlanDay[]>([])
const meals = ref<PlanMeal[]>([])
const options = ref<PlanOption[]>([])
const exercises = ref<PlanExercise[]>([])
const prescriptions = ref<PlanPrescription[]>([])

async function createPlan(payload: { title?: string; start_date: string; end_date: string; notes?: string; daily_water_target_ml?: number }) {
  loading.value = true
  error.value = ''
  try {
    const data = await api.createPlan(clientId.value, payload)
    planId.value = data.id
  } catch {
    error.value = 'ایجاد برنامه انجام نشد.'
  }
  loading.value = false
}

async function addDay() {
  if (!planId.value) return
  const next = (days.value[days.value.length - 1]?.day_number || 0) + 1
  const data = await api.addDay(planId.value, { day_number: next })
  if (data.data) {
    days.value = [...days.value, data.data]
  }
}

async function removeDay(dayId: string) {
  if (!planId.value) return
  await api.deleteDay(planId.value, dayId)
  days.value = days.value.filter((d) => d.id !== dayId)
}

async function addMeal(dayId: string) {
  if (!planId.value) return
  const order = meals.value.filter((m) => m.day_id === dayId).length + 1
  const data = await api.addMeal(planId.value, dayId, { title: `وعده ${order}` })
  if (data.data) {
    meals.value = [...meals.value, data.data]
  }
}

async function removeMeal(mealId: string) {
  if (!planId.value) return
  const meal = meals.value.find((item) => item.id === mealId)
  if (!meal) return
  await api.deleteMeal(planId.value, meal.day_id, mealId)
  meals.value = meals.value.filter((m) => m.id !== mealId)
  options.value = options.value.filter((o) => o.meal_id !== mealId)
}

async function addOption(mealId: string) {
  if (!planId.value) return
  const meal = meals.value.find((item) => item.id === mealId)
  if (!meal) return
  const optionNumber = options.value.filter((o) => o.meal_id === mealId).length + 1
  const data = await api.addOption(planId.value, meal.day_id, mealId, { option_number: optionNumber })
  if (data.data) {
    options.value = [...options.value, data.data]
  }
}

async function removeOption(optionId: string) {
  if (!planId.value) return
  const option = options.value.find((item) => item.id === optionId)
  if (!option) return
  const meal = meals.value.find((item) => item.id === option.meal_id)
  if (!meal) return
  await api.deleteOption(planId.value, meal.day_id, meal.id, optionId)
  options.value = options.value.filter((o) => o.id !== optionId)
}

async function addItemFromFood(payload: { optionId: string; foodId: string; unit: string }) {
  if (!planId.value) return
  const option = options.value.find((item) => item.id === payload.optionId)
  if (!option) return
  const meal = meals.value.find((item) => item.id === option.meal_id)
  if (!meal) return

  const data = await api.addItem(planId.value, meal.day_id, meal.id, option.id, {
    food_id: payload.foodId,
    quantity: 1,
    unit: payload.unit,
  })

  if (!data.data) return
  const idx = options.value.findIndex((item) => item.id === option.id)
  if (idx < 0) return
  const current = options.value[idx]
  options.value[idx] = {
    ...current,
    items: [...current.items, data.data],
  }
  options.value = [...options.value]
}

async function addPrescriptionFromMedication(payload: { medicationId: string }) {
  if (!planId.value || days.value.length === 0) return
  const dayId = days.value[0].id
  const data = await api.addPrescription(planId.value, dayId, {
    medication_id: payload.medicationId,
    dosage: '1 عدد',
  })
  if (data.data) {
    prescriptions.value = [...prescriptions.value, data.data]
  }
}

async function removePrescription(prescriptionId: string) {
  if (!planId.value || days.value.length === 0) return
  const dayId = days.value[0].id
  await api.removePrescription(planId.value, dayId, prescriptionId)
  prescriptions.value = prescriptions.value.filter((item) => item.id !== prescriptionId)
}

async function removeExercise(exerciseId: string) {
  if (!planId.value || days.value.length === 0) return
  const dayId = days.value[0].id
  await api.removeExercise(planId.value, dayId, exerciseId)
  exercises.value = exercises.value.filter((item) => item.id !== exerciseId)
}

async function finishAndBack() {
  await navigateTo(`/nutritionist/clients/${clientId.value}`)
}
</script>

<template>
  <main class="page">
    <header>
      <h2>ایجاد برنامه جدید</h2>
      <p>ابتدا دوره برنامه را ثبت کنید، سپس ساختار روز/وعده/گزینه را کامل کنید.</p>
    </header>

    <p v-if="error" class="error">{{ error }}</p>

    <PlanPeriodFormCard :loading="loading" @submit="createPlan" />

    <section v-if="planId" class="builder">
      <div class="toolbar">
        <span>شناسه برنامه: {{ planId }}</span>
        <button type="button" @click="addDay">افزودن روز</button>
      </div>

      <PlanDayEditor :days="days" @add-meal="addMeal" @remove-day="removeDay" />

      <MealEditor :meals="meals" @add-option="addOption" @remove-meal="removeMeal" />

      <OptionEditor :options="options" @add-item-from-food="addItemFromFood" @remove-option="removeOption" />

      <ExercisePrescriptionEditor
        :exercises="exercises"
        :prescriptions="prescriptions"
        @add-prescription-from-medication="addPrescriptionFromMedication"
        @remove-exercise="removeExercise"
        @remove-prescription="removePrescription"
      />

      <button type="button" class="primary" @click="finishAndBack">اتمام و بازگشت به پرونده</button>
    </section>
  </main>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
}

.builder {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  border: 1px dashed #b8cad4;
  border-radius: 10px;
  padding: 8px;
  background: #f5fafc;
}

.error {
  color: #8b2121;
}

button {
  border: 1px solid #c8d2d8;
  border-radius: 8px;
  min-height: 36px;
  padding: 0 10px;
  background: #fff;
}

button.primary {
  border: none;
  background: #0f6b7a;
  color: #fff;
  font-weight: 700;
}
</style>
