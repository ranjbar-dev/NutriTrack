<template>
  <div class="plan-page">
    <!-- Plan Header -->
    <div class="plan-header">
      <div class="header-content">
        <h1>برنامه غذایی</h1>
        <component :is="PlanContextBadge" :is_active="true" />
      </div>
      <div class="metadata">
        <p v-if="plan">
          شروع: {{ formatPersianDate(plan.start_date) }}
        </p>
        <p v-if="isStale" class="stale-marker">
          آخرین به روزرسانی: {{ formatTime(plan.updated_at) }}
        </p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading">
      <p>در حال بارگذاری...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="retry">تلاش دوباره</button>
    </div>

    <!-- Plan Content -->
    <div v-if="plan && !loading && !error" class="plan-content">
      <!-- Days and Meals -->
      <section v-for="day in plan.days" :key="`day-${day.day_of_week}`" class="day-section">
        <h2 class="day-title">{{ getDayName(day.day_of_week) }}</h2>

        <div v-for="meal in day.meals" :key="`meal-${meal.id}`" class="meal-card">
          <h3 class="meal-type">{{ getMealName(meal.meal_type) }}</h3>
          <p v-if="meal.notes" class="meal-notes">{{ meal.notes }}</p>

          <ul class="food-options">
            <li v-for="option in meal.options" :key="`option-${option.id}`" class="food-item">
              <span class="food-name">{{ option.food_name }}</span>
              <span class="food-quantity">{{ option.quantity_grams }}g</span>
              <span class="food-calories">{{ option.calories }}cal</span>
            </li>
          </ul>
        </div>
      </section>

      <!-- Water Target -->
      <section v-if="plan.daily_water_target_ml" class="water-section">
        <h2>هدف آب روزانه</h2>
        <p class="water-target">{{ plan.daily_water_target_ml }} میلی لیتر</p>
      </section>

      <!-- Exercises -->
      <section v-if="plan.exercises && plan.exercises.length > 0" class="exercises-section">
        <h2>فعالیت های ورزشی</h2>
        <div v-for="exercise in plan.exercises" :key="`exercise-${exercise.id}`" class="exercise-item">
          <h3>{{ exercise.name }}</h3>
          <p>مدت: {{ exercise.duration_minutes }} دقیقه</p>
          <p v-if="exercise.intensity">شدت: {{ exercise.intensity }}</p>
          <p v-if="exercise.notes" class="notes">{{ exercise.notes }}</p>
        </div>
      </section>

      <!-- Prescriptions -->
      <section v-if="plan.prescriptions && plan.prescriptions.length > 0" class="prescriptions-section">
        <h2>داروها</h2>
        <div v-for="med in plan.prescriptions" :key="`med-${med.id}`" class="prescription-item">
          <h3>{{ med.medication_name }}</h3>
          <p>تعداد دوز در روز: {{ med.doses_per_day }}</p>
          <p v-if="med.notes" class="notes">{{ med.notes }}</p>
        </div>
      </section>

      <!-- Plan Notes -->
      <section v-if="plan.notes" class="notes-section">
        <h2>یادداشت ها</h2>
        <p>{{ plan.notes }}</p>
      </section>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && !error && !plan" class="empty-state">
      <p>برنامه فعالی برای امروز پیدا نشد</p>
      <p class="help-text">برای دریافت برنامه، اینترنت را وصل کنید و صفحه را تازه سازی کنید.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useClientPlanApi } from '~/app/composables/useClientPlanApi'
import { usePersianFormat } from '~/app/composables/usePersianFormat'
import PlanContextBadge from '~/app/components/client/PlanContextBadge.vue'

const { getActivePlan } = useClientPlanApi()
const { formatPersianDate, formatTime } = usePersianFormat()

const plan = ref<any>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const isStale = computed(() => {
  if (!plan.value) return false
  // Consider cache older than 1 hour as stale
  const now = new Date()
  const updated = new Date(plan.value.updated_at)
  return (now.getTime() - updated.getTime()) > 3600000
})

const getDayName = (dayOfWeek: number) => {
  const days = ['یکشنبه', 'دوشنبه', 'سه شنبه', 'چهارشنبه', 'پنج شنبه', 'جمعه', 'شنبه']
  return days[dayOfWeek] || ''
}

const getMealName = (mealType: string) => {
  const meals: Record<string, string> = {
    breakfast: 'صبحانه',
    lunch: 'ناهار',
    dinner: 'شام',
    snack: 'میان وعده',
  }
  return meals[mealType] || mealType
}

const loadPlan = async () => {
  loading.value = true
  error.value = null
  try {
    plan.value = await getActivePlan()
    if (!plan.value) {
      error.value = 'برنامه ای یافت نشد'
    }
  } catch (err) {
    error.value = 'خطا در بارگذاری برنامه'
    console.error('Plan load error:', err)
  } finally {
    loading.value = false
  }
}

const retry = () => {
  loadPlan()
}

onMounted(() => {
  loadPlan()
})
</script>

<style scoped>
.plan-page {
  padding: var(--spacing-lg);
  direction: rtl;
}

.plan-header {
  margin-bottom: var(--spacing-2xl);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

h1 {
  font-size: 28px;
  font-weight: 600;
}

.metadata {
  font-size: 14px;
  color: #666;
}

.stale-marker {
  color: #ff9800;
  margin-top: var(--spacing-xs);
}

.day-section {
  margin-bottom: var(--spacing-lg);
}

.day-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: var(--spacing-md);
  border-bottom: 2px solid #0f6b7a;
  padding-bottom: var(--spacing-sm);
}

.meal-card {
  background: white;
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-md);
  border-radius: 8px;
  border-right: 4px solid #0f6b7a;
}

.meal-type {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: var(--spacing-sm);
}

.meal-notes {
  font-size: 14px;
  color: #666;
  margin-bottom: var(--spacing-sm);
  font-style: italic;
}

.food-options {
  list-style: none;
  padding: 0;
  margin: 0;
}

.food-item {
  display: flex;
  justify-content: space-between;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid #eee;
  font-size: 14px;
}

.food-name {
  flex: 1;
  font-weight: 500;
}

.food-quantity,
.food-calories {
  color: #666;
  margin-right: var(--spacing-sm);
}

.water-section,
.exercises-section,
.prescriptions-section,
.notes-section {
  background: white;
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  border-radius: 8px;
}

h2 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: var(--spacing-md);
}

.exercise-item,
.prescription-item {
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid #eee;
}

.exercise-item:last-child,
.prescription-item:last-child {
  border-bottom: none;
}

.notes {
  font-size: 14px;
  color: #666;
  margin-top: var(--spacing-sm);
}

.error-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-2xl) var(--spacing-lg);
  color: #666;
}

.loading {
  text-align: center;
  padding: var(--spacing-2xl);
  color: #666;
}
</style>
