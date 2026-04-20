import { defineStore } from 'pinia'
import type { LogWaterPayload, WaterLogEntry } from '~/types/tracking.types'
import { sumWaterAmounts } from '~/utils/tracking'
import { useOfflineApi } from '~/composables/useOfflineApi'

export const useWaterLogStore = defineStore('waterLog', () => {
  const logs = ref<WaterLogEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const totalMl = computed(() => sumWaterAmounts(logs.value))

  async function fetchToday(date?: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const today = date ?? new Date().toISOString().slice(0, 10)
      logs.value = await apiFetch<WaterLogEntry[]>(`/client/water-logs?date=${today}`)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری مصرف آب'
    }
    finally {
      loading.value = false
    }
  }

  async function addWater(amountMl: number, loggedTime?: string) {
    const { clientPost } = useOfflineApi()
    const payload: LogWaterPayload = {
      local_id: crypto.randomUUID(),
      date: new Date().toISOString().slice(0, 10),
      amount_ml: amountMl,
      ...(loggedTime ? { logged_time: loggedTime } : {}),
    }
    const result = await clientPost<WaterLogEntry>('/client/water-logs', payload, { entityType: 'water_log' })
    if (!('queued' in result)) {
      logs.value.push(result)
    }
  }

  function $reset() {
    logs.value = []
    error.value = null
    loading.value = false
  }

  return { logs, totalMl, loading, error, fetchToday, addWater, $reset }
})
