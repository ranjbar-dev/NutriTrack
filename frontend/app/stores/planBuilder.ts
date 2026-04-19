import { defineStore } from 'pinia'
import type { DietPlanResponse, FoodPickerState } from '~/types/plan.types'

export const usePlanBuilderStore = defineStore('planBuilder', () => {
  const currentPlan = ref<DietPlanResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const foodPicker = ref<FoodPickerState>({
    open: false,
    targetOptionId: null,
    searchQuery: '',
    searchResults: [],
    loading: false,
  })

  // Core load action — called after EVERY CRUD mutation (D-22/Pitfall 6)
  async function loadPlan(planId: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<DietPlanResponse>(`/diet-plans/${planId}`)
      currentPlan.value = data
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری برنامه'
    }
    finally {
      loading.value = false
    }
  }

  // Plan header mutations — call loadPlan after each
  async function createPlan(clientId: string, body: object): Promise<DietPlanResponse> {
    const { apiFetch } = useApi()
    return await apiFetch<DietPlanResponse>('/diet-plans', {
      method: 'POST',
      body: JSON.stringify({ client_id: clientId, ...body }),
    })
  }

  async function updatePlanHeader(planId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}`, { method: 'PATCH', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function activatePlan(planId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/activate`, { method: 'PATCH' })
    await loadPlan(planId)
  }

  async function deletePlan(planId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}`, { method: 'DELETE' })
    currentPlan.value = null
  }

  // Day mutations
  async function addDay(planId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days`, { method: 'POST', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function updateDay(planId: string, dayId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}`, { method: 'PUT', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function deleteDay(planId: string, dayId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}`, { method: 'DELETE' })
    await loadPlan(planId)
  }

  // Meal mutations
  async function addMeal(planId: string, dayId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals`, { method: 'POST', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function updateMeal(planId: string, dayId: string, mealId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}`, { method: 'PUT', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function deleteMeal(planId: string, dayId: string, mealId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}`, { method: 'DELETE' })
    await loadPlan(planId)
  }

  async function reorderMeal(planId: string, dayId: string, mealId: string, newOrder: number) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}/order`, {
      method: 'PATCH',
      body: JSON.stringify({ display_order: newOrder }),
    })
    await loadPlan(planId)
  }

  // Option mutations
  async function addOption(planId: string, dayId: string, mealId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}/options`, { method: 'POST' })
    await loadPlan(planId)
  }

  async function deleteOption(planId: string, dayId: string, mealId: string, optId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}/options/${optId}`, { method: 'DELETE' })
    await loadPlan(planId)
  }

  // Item mutations
  async function addItem(planId: string, dayId: string, mealId: string, optId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}/options/${optId}/items`, {
      method: 'POST',
      body: JSON.stringify(body),
    })
    await loadPlan(planId)
  }

  async function updateItem(planId: string, dayId: string, mealId: string, optId: string, itemId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}/options/${optId}/items/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify(body),
    })
    await loadPlan(planId)
  }

  async function deleteItem(planId: string, dayId: string, mealId: string, optId: string, itemId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/meals/${mealId}/options/${optId}/items/${itemId}`, { method: 'DELETE' })
    await loadPlan(planId)
  }

  // Exercise mutations
  async function addExercise(planId: string, dayId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/exercises`, { method: 'POST', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function updateExercise(planId: string, dayId: string, exId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/exercises/${exId}`, { method: 'PUT', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function deleteExercise(planId: string, dayId: string, exId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/days/${dayId}/exercises/${exId}`, { method: 'DELETE' })
    await loadPlan(planId)
  }

  // Medication mutations
  async function addMedication(planId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/medications`, { method: 'POST', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function updateMedication(planId: string, medId: string, body: object) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/medications/${medId}`, { method: 'PUT', body: JSON.stringify(body) })
    await loadPlan(planId)
  }

  async function deleteMedication(planId: string, medId: string) {
    const { apiFetch } = useApi()
    await apiFetch(`/diet-plans/${planId}/medications/${medId}`, { method: 'DELETE' })
    await loadPlan(planId)
  }

  // Food picker helpers
  function openFoodPicker(optionId: string) {
    foodPicker.value.open = true
    foodPicker.value.targetOptionId = optionId
    foodPicker.value.searchQuery = ''
    foodPicker.value.searchResults = []
  }

  function closeFoodPicker() {
    foodPicker.value.open = false
    foodPicker.value.targetOptionId = null
  }

  function $reset() {
    currentPlan.value = null
    loading.value = false
    error.value = null
    foodPicker.value = {
      open: false,
      targetOptionId: null,
      searchQuery: '',
      searchResults: [],
      loading: false,
    }
  }

  return {
    currentPlan,
    loading,
    error,
    foodPicker,
    loadPlan,
    createPlan,
    updatePlanHeader,
    activatePlan,
    deletePlan,
    addDay,
    updateDay,
    deleteDay,
    addMeal,
    updateMeal,
    deleteMeal,
    reorderMeal,
    addOption,
    deleteOption,
    addItem,
    updateItem,
    deleteItem,
    addExercise,
    updateExercise,
    deleteExercise,
    addMedication,
    updateMedication,
    deleteMedication,
    openFoodPicker,
    closeFoodPicker,
    $reset,
  }
})
