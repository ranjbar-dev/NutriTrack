import { defineStore } from 'pinia'
import dayjs from 'dayjs'
import type { DietPlanListResponse, DietPlanResponse, DietPlanSummary, PlanDayResponse } from '~/types/plan.types'
import { db } from '~/db'

export const useClientPlanStore = defineStore('clientPlan', () => {
  const activePlan = ref<DietPlanResponse | null>(null)
  const myPlans = ref<DietPlanSummary[]>([])
  const activeDayNumber = ref<number>(1)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchActivePlan() {
    // D-06 Step 1: Read from IndexedDB cache first — unblock UI instantly
    try {
      const cached = await db.activePlan.get(1)
      if (cached) {
        activePlan.value = cached.data as DietPlanResponse
        initActiveDay()
      }
      else {
        // D-23: iOS eviction check — if we previously fetched but cache is empty
        const lastFetchMeta = await db.syncMeta.get('plan_last_fetch')
        if (lastFetchMeta) {
          const fetchedAt = new Date(lastFetchMeta.value)
          const oneHourAgo = new Date(Date.now() - 3600_000)
          if (fetchedAt > oneHourAgo) {
            // Storage was evicted — flag for OfflineBanner to show eviction notice
            await db.uiState.put({ key: 'eviction_detected', value: 'true' })
          }
        }
      }
    }
    catch {
      // Dexie unavailable (SSR or DB error) — skip cache, go to network
    }

    // D-06 Step 2: If offline, stop here — cached data (or empty state) is the answer
    if (typeof navigator !== 'undefined' && !navigator.onLine) {
      if (!activePlan.value) {
        error.value = 'داده‌ای در حافظه موجود نیست. پس از اتصال دوباره تلاش کنید.'
      }
      loading.value = false
      return
    }

    // D-06 Step 3: Fetch fresh data from API and update cache
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<DietPlanResponse>('/clients/me/active-plan')
      activePlan.value = data
      initActiveDay()
      // Persist fresh snapshot to Dexie (D-06)
      await db.activePlan.put({
        id: 1,
        plan_id: data.id,
        fetched_at: new Date().toISOString(),
        updated_hint: (data as { updated_at?: string }).updated_at ?? null,
        data,
      })
      await db.syncMeta.put({ key: 'plan_last_fetch', value: new Date().toISOString() })
      // Clear any previous eviction flag
      await db.uiState.delete('eviction_detected')
    }
    catch (e: unknown) {
      const err = e as { statusCode?: number; data?: { error?: string } }
      if (err.statusCode === 404) activePlan.value = null
      else if (!activePlan.value) {
        error.value = (err.data?.error) ?? 'خطا در بارگذاری برنامه'
      }
      // If we have cached data, silently ignore online fetch failure
    }
    finally {
      loading.value = false
    }
  }

  async function fetchMyPlans() {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<DietPlanListResponse>('/clients/me/plans')
      myPlans.value = data.data ?? []
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری تاریخچه برنامه‌ها'
    }
    finally {
      loading.value = false
    }
  }

  async function fetchPlanById(planId: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      return await apiFetch<DietPlanResponse>(`/diet-plans/${planId}`)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری برنامه'
      throw e
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
    myPlans.value = []
    activeDayNumber.value = 1
    loading.value = false
    error.value = null
  }

  return {
    activePlan,
    myPlans,
    activeDayNumber,
    activeDay,
    loading,
    error,
    fetchActivePlan,
    fetchMyPlans,
    fetchPlanById,
    initActiveDay,
    setActiveDay,
    $reset,
  }
})
