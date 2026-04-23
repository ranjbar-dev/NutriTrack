/**
 * Offline Queue and Sync Type Definitions
 * Defines all types for client-side offline tracking queue and sync state management
 */

export type TrackingDomain = 'food' | 'water' | 'sleep' | 'exercise' | 'medication' | 'body' | 'message'

export type SyncState = 'queued' | 'syncing' | 'synced' | 'failed'

export interface TrackingQueueEntry {
  local_id: string
  domain: TrackingDomain
  payload: Record<string, any>
  sync_state: SyncState
  created_at: string
  synced_at?: string
  error_metadata?: {
    code: string
    message: string
  }
  retry_count: number
  last_retry_at?: string
}

export interface FoodTrackingPayload {
  food_id: string
  quantity_grams: number
  consumed_at: string
  notes?: string
}

export interface MessageQueuePayload {
  content: string
}

export interface WaterTrackingPayload {
  amount_ml: number
  logged_at: string
}

export interface SleepTrackingPayload {
  sleep_hours: number
  logged_at: string
  quality?: 'poor' | 'fair' | 'good' | 'excellent'
}

export interface ExerciseTrackingPayload {
  exercise_id: string
  duration_minutes: number
  logged_at: string
  intensity?: 'light' | 'moderate' | 'vigorous'
  notes?: string
}

export interface MedicationTrackingPayload {
  medication_id: string
  doses: number
  logged_at: string
  notes?: string
}

export interface BodyMeasurementPayload {
  weight_kg?: number
  waist_cm?: number
  hip_cm?: number
  chest_cm?: number
  logged_at: string
}

export type TrackingPayload =
  | FoodTrackingPayload
  | WaterTrackingPayload
  | SleepTrackingPayload
  | ExerciseTrackingPayload
  | MedicationTrackingPayload
  | BodyMeasurementPayload

export interface BulkSyncRequest {
  entries: Array<{
    local_id: string
    domain: TrackingDomain
    payload: TrackingPayload
  }>
}

export interface BulkSyncResponse {
  synced: Array<{
    local_id: string
    server_id: string
    domain: TrackingDomain
  }>
  conflicts?: Array<{
    local_id: string
    domain: TrackingDomain
    reason: string
  }>
  errors?: Array<{
    local_id: string
    domain: TrackingDomain
    error: string
  }>
}

export interface QueueMetrics {
  queued: number
  syncing: number
  synced: number
  failed: number
  total: number
}

export interface OfflineQueueState {
  entries: TrackingQueueEntry[]
  metrics: QueueMetrics
  lastSyncAttempt?: string
  isReplayInProgress: boolean
}
