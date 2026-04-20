<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'client' })

const planStore = useClientPlanStore()
const foodLogStore = useFoodLogStore()
const activeDay = computed(() => planStore.activeDay)

onMounted(async () => {
  if (!planStore.activePlan) await planStore.fetchActivePlan()
  await foodLogStore.fetchToday()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6">
    <h1 class="mb-4 text-lg font-bold text-gray-800 text-start">ثبت وعده‌های غذایی</h1>
    <div v-if="!planStore.activePlan" class="rounded-2xl bg-white p-6 text-center text-sm text-gray-500 shadow-sm">
      برنامه غذایی فعالی ندارید
    </div>
    <div v-else class="space-y-3">
      <FoodLogMealCard
        v-for="meal in activeDay?.meals ?? []"
        :key="meal.id"
        :meal="meal"
        :log="foodLogStore.getLogForMeal(meal.id)"
        @log="(optionId) => foodLogStore.logFood(meal.id, optionId)"
        @skip="() => foodLogStore.skipMeal(meal.id)"
      />
    </div>
  </div>
</template>
