import { useAsyncData } from '#app'
import type {
  AdminCreateNutritionistRequest,
  AdminMutationMessageResponse,
  AdminNutritionist,
  AdminNutritionistClientListResponse,
  AdminNutritionistListResponse,
  AdminNutritionistQuery,
  AdminNutritionistStatusRequest,
  AdminUpdateNutritionistRequest,
} from '~/types/admin'

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

export const useAdminNutritionistApi = () => {
  async function listNutritionists(query: AdminNutritionistQuery = {}) {
    const page = query.page ?? 1
    const pageSize = query.page_size ?? 20
    const q = query.q?.trim()
    const key = `admin-nutritionists-${page}-${pageSize}-${q ?? ''}`
    const suffix = withQuery({
      page,
      page_size: pageSize,
      q,
    })

    return useAsyncData(key, () =>
      $fetch<AdminNutritionistListResponse>(`${baseUrl}/admin/nutritionists${suffix}`)
    )
  }

  async function getNutritionist(nutritionistId: string) {
    return useAsyncData(`admin-nutritionist-${nutritionistId}`, () =>
      $fetch<AdminNutritionist>(`${baseUrl}/admin/nutritionists/${nutritionistId}`)
    )
  }

  async function createNutritionist(payload: AdminCreateNutritionistRequest) {
    return $fetch<AdminNutritionist>(`${baseUrl}/admin/nutritionists`, {
      method: 'POST',
      body: payload,
    })
  }

  async function updateNutritionist(nutritionistId: string, payload: AdminUpdateNutritionistRequest) {
    return $fetch<AdminNutritionist>(`${baseUrl}/admin/nutritionists/${nutritionistId}`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function setNutritionistStatus(
    nutritionistId: string,
    payload: AdminNutritionistStatusRequest,
  ) {
    return $fetch<AdminMutationMessageResponse>(`${baseUrl}/admin/nutritionists/${nutritionistId}/status`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function listNutritionistClients(nutritionistId: string, query: AdminNutritionistQuery = {}) {
    const page = query.page ?? 1
    const pageSize = query.page_size ?? 20
    const key = `admin-nutritionist-clients-${nutritionistId}-${page}-${pageSize}`
    const suffix = withQuery({
      page,
      page_size: pageSize,
    })

    return useAsyncData(key, () =>
      $fetch<AdminNutritionistClientListResponse>(
        `${baseUrl}/admin/nutritionists/${nutritionistId}/clients${suffix}`,
      )
    )
  }

  return {
    listNutritionists,
    getNutritionist,
    createNutritionist,
    updateNutritionist,
    setNutritionistStatus,
    listNutritionistClients,
  }
}