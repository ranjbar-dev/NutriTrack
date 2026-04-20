import type { ClientListItem, ClientListResponse, ClientProfile } from '~/types/clientManagement.types'

export const useClientManagementStore = defineStore('clientManagement', () => {
  const clients = ref<ClientListItem[]>([])
  const total = ref(0)
  const currentPage = ref(1)
  const selectedClient = ref<ClientProfile | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchClients(params: {
    q?: string
    sort?: string
    active?: boolean | null
    page?: number
    limit?: number
  } = {}) {
    const { apiFetch } = useApi()
    loading.value = true
    error.value = null
    try {
      const qs = new URLSearchParams()
      if (params.q) qs.set('q', params.q)
      if (params.sort) qs.set('sort', params.sort)
      if (params.active !== null && params.active !== undefined) qs.set('active', String(params.active))
      if (params.page) qs.set('page', String(params.page))
      if (params.limit) qs.set('limit', String(params.limit))

      const data = await apiFetch<ClientListResponse>(`/nutritionist/clients?${qs}`)
      clients.value = data.items
      total.value = data.total
      currentPage.value = data.page
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در دریافت لیست مراجعین'
    } finally {
      loading.value = false
    }
  }

  async function fetchClientProfile(clientId: string) {
    const { apiFetch } = useApi()
    loading.value = true
    error.value = null
    try {
      selectedClient.value = await apiFetch<ClientProfile>(`/nutritionist/clients/${clientId}`)
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در دریافت پروفایل'
    } finally {
      loading.value = false
    }
  }

  async function setActive(clientId: string, active: boolean) {
    const { apiFetch } = useApi()
    const action = active ? 'activate' : 'deactivate'
    await apiFetch(`/nutritionist/clients/${clientId}/${action}`, { method: 'PATCH' })
    const item = clients.value.find(c => c.id === clientId)
    if (item) item.is_active = active
    if (selectedClient.value?.id === clientId) selectedClient.value.is_active = active
  }

  async function updateProfile(clientId: string, data: { date_of_birth?: string; height_cm?: number }) {
    const { apiFetch } = useApi()
    const updated = await apiFetch<ClientProfile>(`/nutritionist/clients/${clientId}/profile`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
    selectedClient.value = updated
    return updated
  }

  return { clients, total, currentPage, selectedClient, loading, error, fetchClients, fetchClientProfile, setActive, updateProfile }
})
