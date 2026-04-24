import { useAsyncData } from '#app'
import type {
  NutritionistClient,
  NutritionistClientListQuery,
  NutritionistClientListResponse,
  NutritionistClientStatusUpdateRequest,
} from '~/types/nutritionist-workspace'

const baseUrl = '/api/v1'

function withQuery(params: Record<string, string | number | undefined>): string {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== '') {
      query.set(key, String(value))
    }
  })
  const serialized = query.toString()
  return serialized ? `?${serialized}` : ''
}

export const useNutritionistClientApi = () => {
  async function listClients(query: NutritionistClientListQuery = {}) {
    const page = query.page ?? 1
    const pageSize = query.page_size ?? 20
    const q = query.q?.trim()
    const status = query.status
    const sort = query.sort
    const key = `nutritionist-clients-${page}-${pageSize}-${q ?? ''}-${status ?? ''}-${sort ?? ''}`
    const suffix = withQuery({
      page,
      page_size: pageSize,
      q,
      status,
      sort,
    })

    return useAsyncData(key, () => $fetch<NutritionistClientListResponse>(`${baseUrl}/clients${suffix}`))
  }

  async function getClientProfile(clientId: string) {
    return useAsyncData(`nutritionist-client-${clientId}`, () =>
      $fetch<NutritionistClient>(`${baseUrl}/clients/${clientId}`)
    )
  }

  async function updateClient(clientId: string, payload: Partial<NutritionistClient>) {
    return $fetch<NutritionistClient>(`${baseUrl}/clients/${clientId}`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function setClientStatus(clientId: string, payload: NutritionistClientStatusUpdateRequest) {
    return $fetch<NutritionistClient>(`${baseUrl}/clients/${clientId}/status`, {
      method: 'PATCH',
      body: payload,
    })
  }

  return {
    listClients,
    getClientProfile,
    updateClient,
    setClientStatus,
  }
}
