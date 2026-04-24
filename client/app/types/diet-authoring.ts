import type { MedicationItem } from './catalogue'

export interface NutritionTotals {
  calories: number
  protein: number
  carbs: number
  fat: number
  fiber: number
}

export interface NutritionRange {
  min: NutritionTotals
  max: NutritionTotals
}

export interface CreatePlanRequest {
  title?: string
  start_date: string
  end_date: string
  notes?: string
  daily_water_target_ml?: number
}

export interface UpdatePlanRequest {
  title?: string
  start_date?: string
  end_date?: string
  notes?: string
  daily_water_target_ml?: number
}

export interface DietPlanFlat {
  id: string
  client_id: string
  nutritionist_id: string
  title: string | null
  start_date: string
  end_date: string
  notes: string | null
  daily_water_target_ml: number | null
  status: 'draft' | 'active' | 'archived'
  created_at: string
  updated_at: string
}

export interface PlanItem {
  id: string
  option_id: string
  food_id: string
  quantity: number
  unit: string | null
  notes: string | null
  computed: NutritionTotals
  food: {
    id: string
    name: string
    unit: string
  }
  created_at: string
}

export interface PlanOption {
  id: string
  meal_id: string
  option_number: number
  totals: NutritionTotals
  items: PlanItem[]
  created_at: string
}

export interface PlanMeal {
  id: string
  day_id: string
  title: string | null
  scheduled_time: string | null
  display_order: number | null
  total_range: NutritionRange | null
  options: PlanOption[]
  created_at: string
}

export interface PlanExercise {
  id: string
  day_id: string
  exercise_name: string
  duration_minutes: number | null
  description: string | null
  calories_burn_estimate: number | null
  created_at: string
}

export interface PlanPrescription {
  id: string
  day_id: string
  medication_id: string
  dosage: string | null
  frequency: string | null
  times: string[]
  instructions: string | null
  start_date: string | null
  end_date: string | null
  medication: MedicationItem | null
  created_at: string
}

export interface PlanDay {
  id: string
  plan_id: string
  day_number: number
  total_range: NutritionRange | null
  meals: PlanMeal[]
  exercises: PlanExercise[]
  prescriptions: PlanPrescription[]
  created_at: string
}

export interface DietPlanFull extends DietPlanFlat {
  days: PlanDay[]
}

export interface PlanDayCreateRequest {
  day_number: number
}

export interface PlanMealCreateRequest {
  title?: string
  scheduled_time?: string
  display_order?: number
}

export interface PlanOptionCreateRequest {
  option_number: number
}

export interface PlanItemCreateRequest {
  food_id: string
  quantity: number
  unit?: string
  notes?: string
}

export interface PlanExerciseCreateRequest {
  exercise_name: string
  duration_minutes?: number
  description?: string
  calories_burn_estimate?: number
}

export interface PlanPrescriptionCreateRequest {
  medication_id: string
  dosage?: string
  frequency?: string
  times?: string[]
  instructions?: string
  start_date?: string
  end_date?: string
}

export interface ApiDataEnvelope<T> {
  data: T
}

export interface PaginatedEnvelope<T> {
  data: T[]
  meta: {
    page: number
    page_size: number
    total: number
  }
}
