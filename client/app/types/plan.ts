/**
 * Plan Types
 * Type definitions for active and archived diet plan data
 */

export interface PlanDay {
  day_of_week: number
  date?: string
  meals: PlanMeal[]
}

export interface PlanMeal {
  id: string
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack'
  options: PlanOption[]
  notes?: string
}

export interface PlanOption {
  id: string
  food_name: string
  quantity_grams: number
  calories: number
  category?: string
}

export interface PlanExercise {
  id: string
  name: string
  duration_minutes: number
  intensity?: string
  notes?: string
}

export interface PlanPrescription {
  id: string
  medication_name: string
  doses_per_day: number
  notes?: string
}

export interface ActivePlanView {
  id: string
  start_date: string
  end_date?: string
  daily_water_target_ml: number
  days: PlanDay[]
  exercises: PlanExercise[]
  prescriptions: PlanPrescription[]
  notes?: string
  is_active: true
  updated_at: string
  freshness: 'remote' | 'cache'
}

export interface ArchivedPlanView {
  id: string
  start_date: string
  end_date: string
  daily_water_target_ml: number
  summary?: string
  is_active: false
  updated_at: string
}

export interface PlanContextLabel {
  is_active: boolean
  label: string
  style: 'active' | 'archived'
}
