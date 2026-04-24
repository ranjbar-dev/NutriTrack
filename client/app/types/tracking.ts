/**
 * Tracking API Type Definitions
 * Defines all types for client tracking API endpoints (from docs/API.md)
 */

import type {
  FoodTrackingPayload,
  WaterTrackingPayload,
  SleepTrackingPayload,
  ExerciseTrackingPayload,
  MedicationTrackingPayload,
  BodyMeasurementPayload,
  TrackingDomain,
} from './offline-sync'

/**
 * POST /tracking/food
 * Creates a food tracking entry
 */
export interface CreateFoodTrackingRequest {
  local_id: string
  food_id: string
  quantity_grams: number
  consumed_at: string
  notes?: string
}

export interface TrackingResponse {
  id: string
  domain: TrackingDomain
  created_at: string
  synced_at?: string
}

/**
 * POST /tracking/water
 * Creates a water intake tracking entry
 */
export interface CreateWaterTrackingRequest {
  local_id: string
  amount_ml: number
  logged_at: string
}

/**
 * POST /tracking/sleep
 * Creates a sleep tracking entry
 */
export interface CreateSleepTrackingRequest {
  local_id: string
  sleep_hours: number
  logged_at: string
  quality?: 'poor' | 'fair' | 'good' | 'excellent'
}

/**
 * POST /tracking/exercise
 * Creates an exercise tracking entry
 */
export interface CreateExerciseTrackingRequest {
  local_id: string
  exercise_id: string
  duration_minutes: number
  logged_at: string
  intensity?: 'light' | 'moderate' | 'vigorous'
  notes?: string
}

/**
 * POST /tracking/medication
 * Creates a medication tracking entry
 */
export interface CreateMedicationTrackingRequest {
  local_id: string
  medication_id: string
  doses: number
  logged_at: string
  notes?: string
}

/**
 * POST /tracking/body
 * Creates a body measurement tracking entry
 */
export interface CreateBodyTrackingRequest {
  local_id: string
  weight_kg?: number
  waist_cm?: number
  hip_cm?: number
  chest_cm?: number
  logged_at: string
}

/**
 * POST /tracking/sync
 * Bulk sync endpoint for offline queue replay
 * Idempotent via local_id - duplicate requests with same local_id return same result
 */
export interface BulkTrackingSyncRequest {
  entries: Array<
    | (CreateFoodTrackingRequest & { domain: 'food' })
    | (CreateWaterTrackingRequest & { domain: 'water' })
    | (CreateSleepTrackingRequest & { domain: 'sleep' })
    | (CreateExerciseTrackingRequest & { domain: 'exercise' })
    | (CreateMedicationTrackingRequest & { domain: 'medication' })
    | (CreateBodyTrackingRequest & { domain: 'body' })
  >
}

/**
 * Response from POST /tracking/sync
 * Contains per-entry sync results and conflict information
 */
export interface BulkTrackingSyncResponse {
  synced: Array<{
    local_id: string
    server_id: string
    domain: TrackingDomain
    created_at: string
  }>
  conflicts?: Array<{
    local_id: string
    domain: TrackingDomain
    reason: string
    server_version?: {
      server_id: string
      created_at: string
    }
  }>
  errors?: Array<{
    local_id: string
    domain: TrackingDomain
    error_code: string
    error_message: string
  }>
}

/**
 * GET /tracking/history
 * Fetches recent tracking entries for a given domain
 * Supports pagination
 */
export interface TrackingHistoryRequest {
  domain?: TrackingDomain
  page?: number
  page_size?: number
  from_date?: string
  to_date?: string
}

export interface TrackingEntry {
  id: string
  local_id?: string
  domain: TrackingDomain
  payload: Record<string, any>
  created_at: string
  synced_at: string
}

export interface TrackingHistoryResponse {
  data: TrackingEntry[]
  meta: {
    total: number
    page: number
    page_size: number
    pages: number
  }
}

/**
 * GET /tracking/summary
 * Gets lightweight tracking summary for a date range
 */
export interface TrackingSummaryRequest {
  from_date: string
  to_date: string
}

export interface TrackingSummaryResponse {
  period: {
    from: string
    to: string
  }
  summary: {
    food_entries: number
    water_servings: number
    sleep_hours: number
    exercise_sessions: number
    medications_taken: number
    body_measurements: number
  }
}

/**
 * Unified type for any tracking API request
 */
export type AnyTrackingRequest =
  | CreateFoodTrackingRequest
  | CreateWaterTrackingRequest
  | CreateSleepTrackingRequest
  | CreateExerciseTrackingRequest
  | CreateMedicationTrackingRequest
  | CreateBodyTrackingRequest

/**
 * Type guard to ensure request has local_id (required for queue/sync)
 */
export function hasLocalId(
  request: any
): request is AnyTrackingRequest & { local_id: string } {
  return typeof request.local_id === 'string' && request.local_id.length > 0
}

/**
 * Type guard to validate payload has required logged_at or consumed_at timestamp
 */
export function hasTimestamp(payload: any): boolean {
  return !!(payload.logged_at || payload.consumed_at)
}
