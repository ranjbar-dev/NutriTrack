export interface FoodResponse {
  id: string
  name: string
  description?: string
  categories: string[]
  calories: number
  protein_g: number
  carbs_g: number
  fat_g: number
  fiber_g: number
  sugar_g: number
  sodium_mg: number
  measurement_unit: string
  measurement_amount: number
  is_active: boolean
  created_by: string
  creator_name: string
  created_at: string
  updated_at: string
}

export interface FoodListResponse {
  data: FoodResponse[]
  total: number
  page: number
  limit: number
  has_more: boolean
}

export interface CreateFoodPayload {
  name: string
  description?: string
  categories: string[]
  calories: number
  protein_g: number
  carbs_g: number
  fat_g: number
  fiber_g: number
  sugar_g: number
  sodium_mg: number
  measurement_unit: string
  measurement_amount: number
}

function normalizePersianText(value: string): string {
  return value
    .replace(/ي/g, 'ی')
    .replace(/ك/g, 'ک')
    .replace(/\s+/g, ' ')
    .trim()
}

export const useFoodStore = defineStore('food', () => {
  const foods = ref<FoodResponse[]>([])
  const currentFood = ref<FoodResponse | null>(null)
  const total = ref(0)
  const page = ref(1)
  const hasMore = ref(false)
  const loading = ref(false)
  const searchQuery = ref('')
  const selectedCategory = ref<string | null>(null)
  const showInactive = ref(false)

  function buildQueryParams() {
    const params = new URLSearchParams()
    if (searchQuery.value) {
      params.set('search', normalizePersianText(searchQuery.value))
    }
    if (selectedCategory.value) {
      params.set('category', selectedCategory.value)
    }
    params.set('is_active', showInactive.value ? 'false' : 'true')
    params.set('page', page.value.toString())
    params.set('limit', '20')
    return params.toString()
  }

  async function fetchFoods(reset = false) {
    loading.value = true
    try {
      if (reset) {
        page.value = 1
        foods.value = []
        hasMore.value = false
        total.value = 0
      }

      const { apiFetch } = useApi()
      const data = await apiFetch<FoodListResponse>(
        `/foods?${buildQueryParams()}`,
      )

      foods.value = reset ? data.data : [...foods.value, ...data.data]
      total.value = data.total
      hasMore.value = data.has_more
      return data
    }
    catch (error) {
      console.error('Failed to fetch foods', error)
      throw error
    }
    finally {
      loading.value = false
    }
  }

  async function fetchFood(id: string) {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<FoodResponse>(`/foods/${id}`)
      currentFood.value = data
      return data
    }
    catch (error) {
      console.error('Failed to fetch food', error)
      currentFood.value = null
      throw error
    }
    finally {
      loading.value = false
    }
  }

  async function createFood(payload: CreateFoodPayload) {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<FoodResponse>('/foods', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      return data
    }
    catch (error) {
      console.error('Failed to create food', error)
      throw error
    }
    finally {
      loading.value = false
    }
  }

  async function updateFood(id: string, payload: CreateFoodPayload) {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<FoodResponse>(`/foods/${id}`, {
        method: 'PUT',
        body: JSON.stringify(payload),
      })
      currentFood.value = data
      return data
    }
    catch (error) {
      console.error('Failed to update food', error)
      throw error
    }
    finally {
      loading.value = false
    }
  }

  async function deleteFood(id: string) {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      await apiFetch(`/foods/${id}`, { method: 'DELETE' })
      foods.value = foods.value.filter((food) => food.id !== id)
      total.value = Math.max(0, total.value - 1)
    }
    catch (error) {
      console.error('Failed to delete food', error)
      throw error
    }
    finally {
      loading.value = false
    }
  }

  async function loadMore() {
    if (loading.value || !hasMore.value) return
    page.value += 1
    await fetchFoods(false)
  }

  async function setSearch(query: string) {
    searchQuery.value = normalizePersianText(query)
    await fetchFoods(true)
  }

  async function setCategory(category: string | null) {
    selectedCategory.value = category
    await fetchFoods(true)
  }

  async function toggleInactive() {
    showInactive.value = !showInactive.value
    await fetchFoods(true)
  }

  return {
    foods,
    currentFood,
    total,
    page,
    hasMore,
    loading,
    searchQuery,
    selectedCategory,
    showInactive,
    fetchFoods,
    fetchFood,
    createFood,
    updateFood,
    deleteFood,
    loadMore,
    setSearch,
    setCategory,
    toggleInactive,
  }
})
