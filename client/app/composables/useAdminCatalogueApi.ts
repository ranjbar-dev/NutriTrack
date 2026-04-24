import { useAsyncData } from '#app'
import type {
  AdminCreateFoodCategoryRequest,
  AdminFoodSearchQuery,
  AdminMedicationSearchQuery,
  FoodCategory,
  FoodItem,
  MedicationItem,
  PaginatedCatalogueResponse,
} from '~/types/catalogue'
import type { AdminMutationMessageResponse } from '~/types/admin'

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

export const useAdminCatalogueApi = () => {
  async function searchAdminFoods(query: AdminFoodSearchQuery = {}) {
    const key = `admin-foods-${query.page ?? 1}-${query.page_size ?? 20}-${query.q ?? ''}-${query.category_id ?? ''}`
    const suffix = withQuery({
      q: query.q?.trim(),
      category_id: query.category_id,
      page: query.page ?? 1,
      page_size: query.page_size ?? 20,
    })

    return useAsyncData(key, () =>
      $fetch<PaginatedCatalogueResponse<FoodItem>>(`${baseUrl}/admin/foods${suffix}`)
    )
  }

  async function forceDeleteFood(foodId: string) {
    return $fetch<void>(`${baseUrl}/admin/foods/${foodId}`, {
      method: 'DELETE',
    })
  }

  async function searchAdminMedications(query: AdminMedicationSearchQuery = {}) {
    const key = `admin-medications-${query.page ?? 1}-${query.page_size ?? 20}-${query.q ?? ''}`
    const suffix = withQuery({
      q: query.q?.trim(),
      page: query.page ?? 1,
      page_size: query.page_size ?? 20,
    })

    return useAsyncData(key, () =>
      $fetch<PaginatedCatalogueResponse<MedicationItem>>(`${baseUrl}/admin/medications${suffix}`)
    )
  }

  async function forceDeleteMedication(medicationId: string) {
    return $fetch<void>(`${baseUrl}/admin/medications/${medicationId}`, {
      method: 'DELETE',
    })
  }

  async function listFoodCategories() {
    return useAsyncData('admin-food-categories', () =>
      $fetch<{ data: FoodCategory[] }>(`${baseUrl}/food-categories`)
    )
  }

  async function createFoodCategory(payload: AdminCreateFoodCategoryRequest) {
    return $fetch<{ data: FoodCategory }>(`${baseUrl}/admin/food-categories`, {
      method: 'POST',
      body: payload,
    })
  }

  async function deleteFoodCategory(categoryId: string) {
    return $fetch<AdminMutationMessageResponse>(`${baseUrl}/admin/food-categories/${categoryId}`, {
      method: 'DELETE',
    })
  }

  return {
    searchAdminFoods,
    forceDeleteFood,
    searchAdminMedications,
    forceDeleteMedication,
    listFoodCategories,
    createFoodCategory,
    deleteFoodCategory,
  }
}