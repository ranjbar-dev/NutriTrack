export type SleepQuality = 'good' | 'fair' | 'poor'
export type LabType = 'blood_test' | 'urine_test' | 'thyroid' | 'hormone' | 'allergy' | 'other'

export interface FoodLogEntry {
  id: string
  client_id: string
  local_id: string
  date: string
  meal_id: string
  selected_option_id?: string | null
  is_skipped: boolean
  notes?: string | null
  created_at: string
  updated_at: string
}

export interface WaterLogEntry {
  id: string
  client_id: string
  local_id: string
  date: string
  amount_ml: number
  logged_time?: string | null
  created_at: string
}

export interface SleepLogEntry {
  id?: string
  client_id?: string
  local_id?: string
  date: string
  sleep_time: string
  wake_time: string
  quality: SleepQuality
  notes?: string | null
  created_at?: string
  updated_at?: string
}

export interface ExerciseLogEntry {
  id: string
  client_id: string
  local_id: string
  date: string
  exercise_name: string
  duration_minutes: number
  calories_burned?: number | null
  notes?: string | null
  created_at: string
}

export interface MedicationLogEntry {
  id: string
  client_id: string
  local_id: string
  date: string
  prescribed_medication_id?: string | null
  medication_name: string
  dosage?: string | null
  taken_at: string
  notes?: string | null
  is_self_reported: boolean
  created_at: string
}

export interface BodyMeasurementEntry {
  id: string
  client_id: string
  local_id: string
  date: string
  weight_kg?: number | null
  waist_cm?: number | null
  hip_cm?: number | null
  abdomen_cm?: number | null
  thigh_cm?: number | null
  chest_cm?: number | null
  wrist_cm?: number | null
  recorded_by: string
  created_at: string
  updated_at: string
}

export interface WeightPoint {
  date: string
  weight_kg: number
}

export interface LabResultResponse {
  id: string
  client_id: string
  local_id: string
  uploaded_by: string
  title: string
  lab_type: LabType
  test_date: string
  has_file: boolean
  external_link?: string | null
  original_filename?: string | null
  mime_type?: string | null
  file_size_bytes?: number | null
  created_at: string
}

export interface DailyDashboard {
  date: string
  water_total_ml: number
  water_target_ml?: number | null
  meals_logged: number
  meals_total: number
  sleep_log?: SleepLogEntry | null
  exercise_count: number
  medication_taken_count: number
  body_logged_today: boolean
  today_body_measurement?: BodyMeasurementEntry | null
  recent_lab_results: LabResultResponse[]
}

export interface LogFoodPayload {
  local_id: string
  date: string
  meal_id: string
  selected_option_id?: string
  is_skipped: boolean
  notes?: string
}

export interface LogWaterPayload {
  local_id: string
  date: string
  amount_ml: number
  logged_time?: string
}

export interface UpsertSleepPayload {
  local_id: string
  date: string
  sleep_time: string
  wake_time: string
  quality: SleepQuality
  notes?: string
}

export interface LogExercisePayload {
  local_id: string
  date: string
  exercise_name: string
  duration_minutes: number
  calories_burned?: number
  notes?: string
}

export interface LogMedicationPayload {
  local_id: string
  date: string
  prescribed_medication_id?: string
  medication_name: string
  dosage?: string
  taken_at: string
  notes?: string
  is_self_reported: boolean
}

export interface UpsertBodyMeasurementPayload {
  local_id: string
  date: string
  weight_kg?: number
  waist_cm?: number
  hip_cm?: number
  abdomen_cm?: number
  thigh_cm?: number
  chest_cm?: number
  wrist_cm?: number
}

export const LAB_TYPE_LABELS: Record<LabType, string> = {
  blood_test: 'آزمایش خون',
  urine_test: 'آزمایش ادرار',
  thyroid: 'تیروئید',
  hormone: 'هورمون',
  allergy: 'آلرژی',
  other: 'سایر',
}

export const SLEEP_QUALITY_LABELS: Record<SleepQuality, string> = {
  good: 'خوب',
  fair: 'متوسط',
  poor: 'ضعیف',
}
