/**
 * useClientPlanApi Composable
 * Typed API composable for client plan reads (active plan, history, lookups)
 * Integrated with offline cache for plan visibility without internet
 */

import { useAsyncData } from '#app'
import type { useAuthApi } from './useAuthApi'

interface DietPlanDay {
  day_of_week: number
  meals: DietPlanMeal[]
}

interface DietPlanMeal {
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack'
  options: DietPlanOption[]
  notes?: string
}

interface DietPlanOption {
  id: string
  food_name: string
  quantity_grams: number
  calories: number
}

interface Exercise {
  id: string
  name: string
  duration_minutes: number
  notes?: string
}

interface Prescription {
  id: string
  medication_name: string
  doses_per_day: number
  notes?: string
}

export interface ActiveDietPlan {
  id: string
  client_id: string
  nutritionist_id: string
  created_at: string
  updated_at: string
  is_active: boolean
  start_date: string
  end_date?: string
  daily_water_target_ml: number
  days: DietPlanDay[]
  exercises: Exercise[]
  prescriptions: Prescription[]
  notes?: string
}

export interface ArchivedDietPlan {
  id: string
  client_id: string
  nutritionist_id: string
  created_at: string
  updated_at: string
  is_active: false
  start_date: string
  end_date: string
  daily_water_target_ml: number
  summary?: string
}

export interface FoodLookup {
  id: string
  name: string
  category: string
  calories_per_gram: number
}

export interface FoodCategoryLookup {
  id: string
  name: string
  foods: FoodLookup[]
}

export interface ExerciseLookup {
  id: string
  name: string
  category: string
}

export interface MedicationLookup {
  id: string
  name: string
  typical_doses: string
}

const baseUrl = '/api/v1'

export const useClientPlanApi = () => {
  /**
   * Fetch the active diet plan for the current client
   */
  async function getActivePlan(): Promise<ActiveDietPlan | null> {
    try {
      const { data, error } = await useAsyncData('active-plan', async () => {
        const response = await $fetch<{ data: ActiveDietPlan }>(
          `${baseUrl}/diet-plans/active`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          }
        )
        return response?.data
      })

      if (error.value || !data.value) {
        console.error('Failed to fetch active plan:', error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error('Active plan fetch exception:', err)
      return null
    }
  }

  /**
   * Fetch archived diet plans for the current client
   */
  async function getArchivedPlans(
    page = 1,
    pageSize = 10
  ): Promise<{ plans: ArchivedDietPlan[]; total: number } | null> {
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        page_size: pageSize.toString(),
      })

      const { data, error } = await useAsyncData(
        `archived-plans-${page}`,
        async () => {
          const response = await $fetch<{
            data: ArchivedDietPlan[]
            meta: { total: number }
          }>(`${baseUrl}/diet-plans/archived?${params.toString()}`, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          })
          return {
            plans: response?.data || [],
            total: response?.meta?.total || 0,
          }
        }
      )

      if (error.value || !data.value) {
        console.error('Failed to fetch archived plans:', error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error('Archived plans fetch exception:', err)
      return null
    }
  }

  /**
   * Fetch single diet plan by ID (active or archived)
   */
  async function getPlanById(planId: string): Promise<ActiveDietPlan | ArchivedDietPlan | null> {
    try {
      const { data, error } = await useAsyncData(`plan-${planId}`, async () => {
        const response = await $fetch<{ data: ActiveDietPlan | ArchivedDietPlan }>(
          `${baseUrl}/diet-plans/${planId}`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          }
        )
        return response?.data
      })

      if (error.value || !data.value) {
        console.error(`Failed to fetch plan ${planId}:`, error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error(`Plan fetch exception for ${planId}:`, err)
      return null
    }
  }

  /**
   * Fetch food categories and foods lookup for tracking UI
   */
  async function getFoodLookup(): Promise<FoodCategoryLookup[] | null> {
    try {
      const { data, error } = await useAsyncData('food-lookup', async () => {
        const response = await $fetch<{ data: FoodCategoryLookup[] }>(
          `${baseUrl}/foods/categories`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          }
        )
        return response?.data
      })

      if (error.value || !data.value) {
        console.error('Failed to fetch food lookup:', error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error('Food lookup fetch exception:', err)
      return null
    }
  }

  /**
   * Fetch exercise lookup for tracking UI
   */
  async function getExerciseLookup(): Promise<ExerciseLookup[] | null> {
    try {
      const { data, error } = await useAsyncData('exercise-lookup', async () => {
        const response = await $fetch<{ data: ExerciseLookup[] }>(
          `${baseUrl}/exercises`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          }
        )
        return response?.data
      })

      if (error.value || !data.value) {
        console.error('Failed to fetch exercise lookup:', error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error('Exercise lookup fetch exception:', err)
      return null
    }
  }

  /**
   * Fetch medication lookup for tracking UI
   */
  async function getMedicationLookup(): Promise<MedicationLookup[] | null> {
    try {
      const { data, error } = await useAsyncData('medication-lookup', async () => {
        const response = await $fetch<{ data: MedicationLookup[] }>(
          `${baseUrl}/medications`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          }
        )
        return response?.data
      })

      if (error.value || !data.value) {
        console.error('Failed to fetch medication lookup:', error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error('Medication lookup fetch exception:', err)
      return null
    }
  }

  return {
    getActivePlan,
    getArchivedPlans,
    getPlanById,
    getFoodLookup,
    getExerciseLookup,
    getMedicationLookup,
  }
}
