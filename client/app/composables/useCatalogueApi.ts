import { useAsyncData } from '#app'
import type {
  FoodCategory,
  FoodItem,
  FoodSearchQuery,
  MedicationItem,
  MedicationSearchQuery,
  PaginatedCatalogueResponse,
} from '~/types/catalogue'

const baseUrl = '/api/v1'

function buildQuery(params: Record<string, string | number | undefined>): string {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== '') {
      query.set(key, String(value))
    }
  })
  const serialized = query.toString()
  return serialized ? `?${serialized}` : ''
}

export const useCatalogueApi = () => {
  async function searchFoods(query: FoodSearchQuery = {}) {
    const key = `catalogue-foods-${query.q ?? ''}-${query.category_id ?? ''}-${query.page ?? 1}`
    const suffix = buildQuery({
      q: query.q?.trim(),
      category_id: query.category_id,
      page: query.page ?? 1,
      page_size: query.page_size ?? 20,
    })
    return useAsyncData(key, () =>
      $fetch<PaginatedCatalogueResponse<FoodItem>>(`${baseUrl}/foods${suffix}`)
    )
  }

  async function getFoodCategories() {
    return useAsyncData('catalogue-food-categories', () =>
      $fetch<{ data: FoodCategory[] }>(`${baseUrl}/food-categories`)
    )
  }

  async function searchMedications(query: MedicationSearchQuery = {}) {
    const key = `catalogue-medications-${query.q ?? ''}-${query.page ?? 1}`
    const suffix = buildQuery({
      q: query.q?.trim(),
      page: query.page ?? 1,
      page_size: query.page_size ?? 20,
    })

    return useAsyncData(key, () =>
      $fetch<PaginatedCatalogueResponse<MedicationItem>>(`${baseUrl}/medications${suffix}`)
    )
  }

  return {
    searchFoods,
    getFoodCategories,
    searchMedications,
  }
}
