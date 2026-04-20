import { defineStore } from 'pinia'
import type { DailyDashboard } from '~/types/tracking.types'

export const useTrackingStore = defineStore('tracking', () => {
  const dashboard = ref<DailyDashboard | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchDailyDashboard(date?: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const today = date ?? new Date().toISOString().slice(0, 10)
      dashboard.value = await apiFetch<DailyDashboard>(`/client/tracking/daily?date=${today}`)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری خلاصه روزانه'
    }
    finally {
      loading.value = false
    }
  }

  const waterProgressPercent = computed(() => {
    const total = dashboard.value?.water_total_ml ?? 0
    const target = dashboard.value?.water_target_ml ?? 0
    if (!target) return 0
    return Math.min(100, Math.round((total / target) * 100))
  })

  function $reset() {
    dashboard.value = null
    error.value = null
    loading.value = false
  }

  return { dashboard, loading, error, waterProgressPercent, fetchDailyDashboard, $reset }
})
