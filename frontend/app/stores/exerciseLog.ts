import { defineStore } from 'pinia'
import type { ExerciseLogEntry } from '~/types/tracking.types'
import { useOfflineApi } from '~/composables/useOfflineApi'

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
    const { clientPost } = useOfflineApi()
    const result = await clientPost<ExerciseLogEntry>('/client/exercise-logs', {
      local_id: crypto.randomUUID(),
      date: new Date().toISOString().slice(0, 10),
      ...payload,
    }, { entityType: 'exercise_log' })
    if (!('queued' in result)) {
      todayLogs.value.unshift(result)
    }
  }

  function $reset() {
    todayLogs.value = []
    error.value = null
    loading.value = false
  }

  return { todayLogs, loading, error, fetchToday, logExercise, $reset }
})
