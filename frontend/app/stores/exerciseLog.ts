import { defineStore } from 'pinia'
import type { ExerciseLogEntry } from '~/types/tracking.types'

export const useExerciseLogStore = defineStore('exerciseLog', () => {
  const todayLogs = ref<ExerciseLogEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchToday(date?: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const today = date ?? new Date().toISOString().slice(0, 10)
      todayLogs.value = await apiFetch<ExerciseLogEntry[]>(`/client/exercise-logs?date=${today}`)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری تمرین‌ها'
    }
    finally {
      loading.value = false
    }
  }

  async function logExercise(payload: { exercise_name: string; duration_minutes: number; calories_burned?: number; notes?: string }) {
    const { apiFetch } = useApi()
    const entry = await apiFetch<ExerciseLogEntry>('/client/exercise-logs', {
      method: 'POST',
      body: JSON.stringify({
        local_id: crypto.randomUUID(),
        date: new Date().toISOString().slice(0, 10),
        ...payload,
      }),
    })
    todayLogs.value.unshift(entry)
  }

  function $reset() {
    todayLogs.value = []
    error.value = null
    loading.value = false
  }

  return { todayLogs, loading, error, fetchToday, logExercise, $reset }
})
