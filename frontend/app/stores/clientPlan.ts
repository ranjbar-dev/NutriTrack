import { defineStore } from 'pinia'
import dayjs from 'dayjs'
import type { DietPlanResponse, PlanDayResponse } from '~/types/plan.types'

export const useClientPlanStore = defineStore('clientPlan', () => {
  const activePlan = ref<DietPlanResponse | null>(null)
  const activeDayNumber = ref<number>(1)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchActivePlan() {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<DietPlanResponse>('/clients/me/active-plan')
      activePlan.value = data
      initActiveDay()
    }
    catch (e: unknown) {
      const err = e as { statusCode?: number; data?: { error?: string } }
      if (err.statusCode === 404) {
        activePlan.value = null
      }
      else {
        error.value = (err.data?.error) ?? 'خطا در بارگذاری برنامه'
      }
    }
    finally {
      loading.value = false
    }
  }

  function initActiveDay() {
    if (!activePlan.value) {
      activeDayNumber.value = 1
      return
    }
    const startDate = dayjs(activePlan.value.start_date)
    const today = dayjs()
    // day_number is 1-indexed: day 1 = start_date, day 2 = start_date + 1, etc.
    const diffDays = today.diff(startDate, 'day') + 1
    const totalDays = activePlan.value.days.length

    // Clamp to valid day range [1, totalDays]
    if (diffDays < 1) {
      activeDayNumber.value = 1
    }
    else if (totalDays > 0 && diffDays > totalDays) {
      activeDayNumber.value = totalDays
    }
    else {
      activeDayNumber.value = diffDays
    }
  }

  const activeDay = computed<PlanDayResponse | null>(() => {
    if (!activePlan.value) return null
    return activePlan.value.days.find(d => d.day_number === activeDayNumber.value) ?? null
  })

  function setActiveDay(dayNumber: number) {
    activeDayNumber.value = dayNumber
  }

  function $reset() {
    activePlan.value = null
    activeDayNumber.value = 1
    loading.value = false
    error.value = null
  }

  return {
    activePlan,
    activeDayNumber,
    activeDay,
    loading,
    error,
    fetchActivePlan,
    initActiveDay,
    setActiveDay,
    $reset,
  }
})
