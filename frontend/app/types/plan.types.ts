export type PlanStatus = 'draft' | 'active' | 'archived'

export interface FoodEmbedded {
  id: string
  name: string
  calories: number
  protein_g: number
  carbs_g: number
  fat_g: number
  fiber_g: number
  measurement_unit: string
  measurement_amount: number
}

export interface MealOptionItemResponse {
  id: string
  option_id: string
  food_id: string
  food: FoodEmbedded
  quantity: number
  measurement_unit: string
  notes?: string
}

export interface NutritionTotals {
  calories: number
  protein_g: number
  carbs_g: number
  fat_g: number
  fiber_g: number
}

export interface MealOptionResponse {
  id: string
  meal_id: string
  option_number: number
  label?: string
  items: MealOptionItemResponse[]
}

export interface MealResponse {
  id: string
  day_id: string
  title: string
  scheduled_time?: string
  display_order: number
  options: MealOptionResponse[]
}

export interface PlanExerciseResponse {
  id: string
  day_id: string
  exercise_name: string
  duration_minutes?: number
  description?: string
  calories_burn_estimate?: number
  display_order: number
}

export interface PlanDayResponse {
  id: string
  day_number: number
  label?: string
  meals: MealResponse[]
  exercises: PlanExerciseResponse[]
}

export interface PlanMedicationResponse {
  id: string
  plan_id: string
  medication_id: string
  medication_name: string
  medication_form: string
  dosage: string
  frequency: string
  times: string[]
  instructions?: string
  start_date?: string
  end_date?: string
}

export interface DietPlanResponse {
  id: string
  client_id: string
  nutritionist_id: string
  start_date: string
  end_date?: string
  notes?: string
  daily_water_target_ml?: number
  status: PlanStatus
  created_at: string
  days: PlanDayResponse[]
  medications: PlanMedicationResponse[]
}

export interface DietPlanSummary {
  id: string
  client_id: string
  status: PlanStatus
  start_date: string
  end_date?: string
  notes?: string
  day_count: number
  created_at: string
}

export interface DietPlanListResponse {
  data: DietPlanSummary[]
}

export interface PlanBreadcrumb {
  label: string
  to?: string
}

export interface FoodPickerState {
  open: boolean
  targetOptionId: string | null
  searchQuery: string
  searchResults: FoodEmbedded[]
  loading: boolean
}
