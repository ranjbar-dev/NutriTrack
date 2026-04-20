import { defineStore } from 'pinia'
import type { FoodLogEntry, LogFoodPayload } from '~/types/tracking.types'
import { useOfflineApi } from '~/composables/useOfflineApi'

export const useFoodLogStore = defineStore('foodLog', () => {
  const todayLogs = ref<FoodLogEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchToday(date?: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const today = date ?? new Date().toISOString().slice(0, 10)
      todayLogs.value = await apiFetch<FoodLogEntry[]>(`/client/food-logs?date=${today}`)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری ثبت غذا'
    }
    finally {
      loading.value = false
    }
  }

  async function logFood(mealId: string, selectedOptionId: string, notes?: string) {
    const { clientPost } = useOfflineApi()
    const payload: LogFoodPayload = {
      local_id: crypto.randomUUID(),
      date: new Date().toISOString().slice(0, 10),
      meal_id: mealId,
      selected_option_id: selectedOptionId,
      is_skipped: false,
      ...(notes ? { notes } : {}),
    }
    const result = await clientPost<FoodLogEntry>('/client/food-logs', payload, { entityType: 'food_log' })
    if (!('queued' in result)) {
      upsertLocal(result)
    }
  }

  async function skipMeal(mealId: string) {
    const { clientPost } = useOfflineApi()
    const payload: LogFoodPayload = {
      local_id: crypto.randomUUID(),
      date: new Date().toISOString().slice(0, 10),
      meal_id: mealId,
      is_skipped: true,
    }
    const result = await clientPost<FoodLogEntry>('/client/food-logs', payload, { entityType: 'food_log' })
    if (!('queued' in result)) {
      upsertLocal(result)
    }
  }

  function upsertLocal(entry: FoodLogEntry) {
    const index = todayLogs.value.findIndex(log => log.meal_id === entry.meal_id)
    if (index >= 0) todayLogs.value[index] = entry
    else todayLogs.value.push(entry)
  }

  function getLogForMeal(mealId: string) {
    return todayLogs.value.find(log => log.meal_id === mealId)
  }

  function $reset() {
    todayLogs.value = []
    error.value = null
    loading.value = false
  }

  return { todayLogs, loading, error, fetchToday, logFood, skipMeal, getLogForMeal, $reset }
})
