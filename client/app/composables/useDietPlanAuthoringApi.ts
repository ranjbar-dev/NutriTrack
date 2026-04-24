import { useAsyncData } from '#app'
import type {
  ApiDataEnvelope,
  CreatePlanRequest,
  DietPlanFlat,
  DietPlanFull,
  PaginatedEnvelope,
  PlanDay,
  PlanDayCreateRequest,
  PlanExercise,
  PlanExerciseCreateRequest,
  PlanItem,
  PlanItemCreateRequest,
  PlanMeal,
  PlanMealCreateRequest,
  PlanOption,
  PlanOptionCreateRequest,
  PlanPrescription,
  PlanPrescriptionCreateRequest,
  UpdatePlanRequest,
} from '~/types/diet-authoring'

const baseUrl = '/api/v1'

function pageSuffix(page = 1, pageSize = 20): string {
  return `?page=${page}&page_size=${pageSize}`
}

export const useDietPlanAuthoringApi = () => {
  async function createPlan(clientId: string, payload: CreatePlanRequest) {
    return $fetch<DietPlanFlat>(`${baseUrl}/clients/${clientId}/plans`, {
      method: 'POST',
      body: payload,
    })
  }

  async function listClientPlans(clientId: string, page = 1, pageSize = 20) {
    return useAsyncData(`nutritionist-client-plans-${clientId}-${page}`, () =>
      $fetch<PaginatedEnvelope<DietPlanFlat>>(`${baseUrl}/clients/${clientId}/plans${pageSuffix(page, pageSize)}`)
    )
  }

  async function getPlan(planId: string) {
    return useAsyncData(`nutritionist-plan-${planId}`, () =>
      $fetch<ApiDataEnvelope<DietPlanFull>>(`${baseUrl}/plans/${planId}`)
    )
  }

  async function updatePlan(planId: string, payload: UpdatePlanRequest) {
    return $fetch<DietPlanFlat>(`${baseUrl}/plans/${planId}`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function deletePlan(planId: string) {
    return $fetch<void>(`${baseUrl}/plans/${planId}`, {
      method: 'DELETE',
    })
  }

  async function addDay(planId: string, payload: PlanDayCreateRequest) {
    return $fetch<ApiDataEnvelope<PlanDay>>(`${baseUrl}/plans/${planId}/days`, {
      method: 'POST',
      body: payload,
    })
  }

  async function deleteDay(planId: string, dayId: string) {
    return $fetch<void>(`${baseUrl}/plans/${planId}/days/${dayId}`, {
      method: 'DELETE',
    })
  }

  async function addMeal(planId: string, dayId: string, payload: PlanMealCreateRequest) {
    return $fetch<ApiDataEnvelope<PlanMeal>>(`${baseUrl}/plans/${planId}/days/${dayId}/meals`, {
      method: 'POST',
      body: payload,
    })
  }

  async function deleteMeal(planId: string, dayId: string, mealId: string) {
    return $fetch<void>(`${baseUrl}/plans/${planId}/days/${dayId}/meals/${mealId}`, {
      method: 'DELETE',
    })
  }

  async function addOption(planId: string, dayId: string, mealId: string, payload: PlanOptionCreateRequest) {
    return $fetch<ApiDataEnvelope<PlanOption>>(
      `${baseUrl}/plans/${planId}/days/${dayId}/meals/${mealId}/options`,
      {
        method: 'POST',
        body: payload,
      }
    )
  }

  async function deleteOption(planId: string, dayId: string, mealId: string, optionId: string) {
    return $fetch<void>(`${baseUrl}/plans/${planId}/days/${dayId}/meals/${mealId}/options/${optionId}`, {
      method: 'DELETE',
    })
  }

  async function addItem(
    planId: string,
    dayId: string,
    mealId: string,
    optionId: string,
    payload: PlanItemCreateRequest
  ) {
    return $fetch<ApiDataEnvelope<PlanItem>>(
      `${baseUrl}/plans/${planId}/days/${dayId}/meals/${mealId}/options/${optionId}/items`,
      {
        method: 'POST',
        body: payload,
      }
    )
  }

  async function deleteItem(planId: string, dayId: string, mealId: string, optionId: string, itemId: string) {
    return $fetch<void>(
      `${baseUrl}/plans/${planId}/days/${dayId}/meals/${mealId}/options/${optionId}/items/${itemId}`,
      {
        method: 'DELETE',
      }
    )
  }

  async function addExercise(planId: string, dayId: string, payload: PlanExerciseCreateRequest) {
    return $fetch<ApiDataEnvelope<PlanExercise>>(`${baseUrl}/plans/${planId}/days/${dayId}/exercises`, {
      method: 'POST',
      body: payload,
    })
  }

  async function removeExercise(planId: string, dayId: string, exerciseId: string) {
    return $fetch<void>(`${baseUrl}/plans/${planId}/days/${dayId}/exercises/${exerciseId}`, {
      method: 'DELETE',
    })
  }

  async function addPrescription(planId: string, dayId: string, payload: PlanPrescriptionCreateRequest) {
    return $fetch<ApiDataEnvelope<PlanPrescription>>(`${baseUrl}/plans/${planId}/days/${dayId}/prescriptions`, {
      method: 'POST',
      body: payload,
    })
  }

  async function removePrescription(planId: string, dayId: string, prescriptionId: string) {
    return $fetch<void>(`${baseUrl}/plans/${planId}/days/${dayId}/prescriptions/${prescriptionId}`, {
      method: 'DELETE',
    })
  }

  return {
    createPlan,
    listClientPlans,
    getPlan,
    updatePlan,
    deletePlan,
    addDay,
    deleteDay,
    addMeal,
    deleteMeal,
    addOption,
    deleteOption,
    addItem,
    deleteItem,
    addExercise,
    removeExercise,
    addPrescription,
    removePrescription,
  }
}
