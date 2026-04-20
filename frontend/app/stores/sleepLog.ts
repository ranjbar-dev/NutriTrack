import { defineStore } from 'pinia'
import type { SleepLogEntry, UpsertSleepPayload } from '~/types/tracking.types'

export const useSleepLogStore = defineStore('sleepLog', () => {
  const todayLog = ref<SleepLogEntry | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchToday(date?: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const today = date ?? new Date().toISOString().slice(0, 10)
      todayLog.value = await apiFetch<SleepLogEntry>(`/client/sleep-logs?date=${today}`)
    }
    catch (e: unknown) {
      const err = e as { statusCode?: number; data?: { error?: string } }
      if (err.statusCode === 404) todayLog.value = null
      else error.value = err.data?.error ?? 'خطا در بارگذاری خواب'
    }
    finally {
      loading.value = false
    }
  }

  async function upsertSleep(payload: Omit<UpsertSleepPayload, 'local_id'>) {
    const { apiFetch } = useApi()
    todayLog.value = await apiFetch<SleepLogEntry>('/client/sleep-logs', {
      method: 'POST',
      body: JSON.stringify({ ...payload, local_id: crypto.randomUUID() }),
    })
  }

  function $reset() {
    todayLog.value = null
    error.value = null
    loading.value = false
  }

  return { todayLog, loading, error, fetchToday, upsertSleep, $reset }
})
