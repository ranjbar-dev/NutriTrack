import type { FoodRequest } from '~/types/foodRequest.types'

export const useFoodRequestStore = defineStore('foodRequest', () => {
  const requests = ref<FoodRequest[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchClientRequests() {
    const { apiFetch } = useApi()
    loading.value = true
    error.value = null
    try {
      requests.value = await apiFetch<FoodRequest[]>('/client/food-requests')
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در دریافت درخواست‌ها'
    } finally {
      loading.value = false
    }
  }

  async function fetchNutriRequests() {
    const { apiFetch } = useApi()
    loading.value = true
    error.value = null
    try {
      requests.value = await apiFetch<FoodRequest[]>('/nutritionist/food-requests')
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در دریافت درخواست‌ها'
    } finally {
      loading.value = false
    }
  }

  async function createRequest(foodName: string, description?: string): Promise<FoodRequest | null> {
    const { apiFetch } = useApi()
    try {
      const fr = await apiFetch<FoodRequest>('/client/food-requests', {
        method: 'POST',
        body: JSON.stringify({ food_name: foodName, description }),
      })
      requests.value.unshift(fr)
      return fr
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در ثبت درخواست'
      return null
    }
  }

  async function approve(requestId: string): Promise<FoodRequest | null> {
    const { apiFetch } = useApi()
    try {
      const fr = await apiFetch<FoodRequest>(`/nutritionist/food-requests/${requestId}/approve`, { method: 'PATCH' })
      const idx = requests.value.findIndex(r => r.id === requestId)
      if (idx !== -1) requests.value[idx] = fr
      return fr
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در تأیید درخواست'
      return null
    }
  }

  async function reject(requestId: string, reason?: string): Promise<FoodRequest | null> {
    const { apiFetch } = useApi()
    try {
      const fr = await apiFetch<FoodRequest>(`/nutritionist/food-requests/${requestId}/reject`, {
        method: 'PATCH',
        body: JSON.stringify({ rejection_reason: reason }),
      })
      const idx = requests.value.findIndex(r => r.id === requestId)
      if (idx !== -1) requests.value[idx] = fr
      return fr
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در رد درخواست'
      return null
    }
  }

  return { requests, loading, error, fetchClientRequests, fetchNutriRequests, createRequest, approve, reject }
})
