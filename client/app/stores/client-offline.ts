/**
 * Client Offline Store
 * Manages offline queue state, sync transitions, and cache lifecycle for client tracking data
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import type {
  TrackingQueueEntry,
  TrackingDomain,
  SyncState,
  QueueMetrics,
  OfflineQueueState,
} from '~/app/types/offline-sync'
import type { Message } from '~/types/messaging'

export const useClientOfflineStore = defineStore('client-offline', () => {
  // ---- Offline Queue ----
  const queueEntries = ref<Map<string, TrackingQueueEntry>>(new Map())
  const isReplayInProgress = ref(false)
  const lastGuardedSignature = ref<string | null>(null)
  const lastGuardedAt = ref<number>(0)

  // ---- Message Cache ----
  const cachedMessages = ref<Message[]>([])

  function isSupportedDomainPayload(
    domain: TrackingDomain,
    payload: Record<string, any>
  ): boolean {
    if (!payload || typeof payload !== 'object') {
      return false
    }

    switch (domain) {
      case 'food':
        return typeof payload.food_id === 'string' && Number(payload.quantity_grams) > 0
      case 'water':
        return Number(payload.amount_ml) > 0
      case 'sleep':
        return Number(payload.sleep_hours) > 0
      case 'exercise':
        return (
          typeof payload.exercise_id === 'string' &&
          Number(payload.duration_minutes) > 0
        )
      case 'medication':
        return (
          typeof payload.medication_id === 'string' &&
          Number(payload.doses) > 0
        )
      case 'body':
        return [
          payload.weight_kg,
          payload.waist_cm,
          payload.hip_cm,
          payload.abdomen_cm,
          payload.thigh_cm,
          payload.chest_cm,
          payload.wrist_cm,
        ].some((field) => Number(field) > 0)
      case 'message':
        return typeof payload.content === 'string' && payload.content.trim().length > 0
      default:
        return false
    }
  }

  // ---- Methods: Queue Mutations ----

  /**
   * Enqueue a new offline tracking write
   * Returns the created entry with generated local_id and queued state
   */
  function enqueueTrackingWrite(request: {
    domain: TrackingDomain
    payload: Record<string, any>
  }): TrackingQueueEntry | null {
    const local_id = uuidv4()
    const now = new Date().toISOString()

    const entry: TrackingQueueEntry = {
      local_id,
      domain: request.domain,
      payload: request.payload,
      sync_state: 'queued',
      created_at: now,
      retry_count: 0,
    }

    queueEntries.value.set(local_id, entry)
    return entry
  }

  /**
   * Guarded enqueue path used by tracking forms.
   * Validates domain schema and blocks accidental double-submit bursts.
   */
  function enqueueDomainTrackingWrite(request: {
    domain: TrackingDomain
    payload: Record<string, any>
  }): TrackingQueueEntry | null {
    if (!isSupportedDomainPayload(request.domain, request.payload)) {
      return null
    }

    const signature = `${request.domain}:${JSON.stringify(request.payload)}`
    const now = Date.now()
    if (lastGuardedSignature.value === signature && now - lastGuardedAt.value < 600) {
      return null
    }

    lastGuardedSignature.value = signature
    lastGuardedAt.value = now
    return enqueueTrackingWrite(request)
  }

  /**
   * Get a queue entry by local_id
   */
  function getQueueEntry(local_id: string): TrackingQueueEntry | undefined {
    return queueEntries.value.get(local_id)
  }

  /**
   * Update sync state of a queue entry
   * Performs immutable state transition
   */
  function updateQueueEntryState(local_id: string, newState: SyncState): void {
    const entry = queueEntries.value.get(local_id)
    if (!entry) return

    // Immutable update
    const updated: TrackingQueueEntry = {
      ...entry,
      sync_state: newState,
    }

    if (newState === 'syncing') {
      updated.last_retry_at = new Date().toISOString()
    } else if (newState === 'synced') {
      updated.synced_at = new Date().toISOString()
    }

    queueEntries.value.set(local_id, updated)
  }

  /**
   * Mark entry as successfully synced
   */
  function markEntrySynced(local_id: string): void {
    updateQueueEntryState(local_id, 'synced')
  }

  /**
   * Mark entry as failed with error metadata
   */
  function markEntryFailed(
    local_id: string,
    error: { code: string; message: string }
  ): void {
    const entry = queueEntries.value.get(local_id)
    if (!entry) return

    const updated: TrackingQueueEntry = {
      ...entry,
      sync_state: 'failed',
      error_metadata: error,
      retry_count: entry.retry_count + 1,
    }

    queueEntries.value.set(local_id, updated)
  }

  /**
   * Delete a queue entry (used for manual dismissal)
   */
  function deleteQueueEntry(local_id: string): void {
    queueEntries.value.delete(local_id)
  }

  /**
   * Clear all offline state on logout
   * Atomic operation - clears queue and cache
   */
  function clearAllOfflineState(): void {
    queueEntries.value.clear()
    isReplayInProgress.value = false
  }

  // ---- Methods: Query and Metrics ----

  /**
   * Get all entries in queued or failed state (pending sync)
   */
  function getPendingEntries(): TrackingQueueEntry[] {
    return Array.from(queueEntries.value.values()).filter(
      e => e.sync_state === 'queued' || e.sync_state === 'failed'
    )
  }

  /**
   * Get all entries in syncing state
   */
  function getSyncingEntries(): TrackingQueueEntry[] {
    return Array.from(queueEntries.value.values()).filter(
      e => e.sync_state === 'syncing'
    )
  }

  /**
   * Get entries by domain (for UI showing pending actions by type)
   */
  function getEntriesByDomain(domain: TrackingDomain): TrackingQueueEntry[] {
    return Array.from(queueEntries.value.values()).filter(
      e => e.domain === domain
    )
  }

  /**
   * Get queue metrics for UX sync state strip
   */
  function getQueueMetrics(): QueueMetrics {
    const entries = Array.from(queueEntries.value.values())

    return {
      queued: entries.filter(e => e.sync_state === 'queued').length,
      syncing: entries.filter(e => e.sync_state === 'syncing').length,
      synced: entries.filter(e => e.sync_state === 'synced').length,
      failed: entries.filter(e => e.sync_state === 'failed').length,
      total: entries.length,
    }
  }

  /**
   * Check if there are any failed entries requiring user attention
   */
  function hasFailedEntries(): boolean {
    return Array.from(queueEntries.value.values()).some(
      e => e.sync_state === 'failed'
    )
  }

  /**
   * Set replay in progress flag to prevent concurrent replays
   */
  function setReplayInProgress(inProgress: boolean): void {
    isReplayInProgress.value = inProgress
  }

  /**
   * Get all queue entries for debugging/export
   */
  function getAllQueueEntries(): TrackingQueueEntry[] {
    return Array.from(queueEntries.value.values())
  }

  /**
   * Set cached messages from API response, trimmed to last 50
   */
  function setCachedMessages(messages: Message[]): void {
    cachedMessages.value = messages.slice(-50)
  }

  // ---- Computed State ----

  const queueState = computed<OfflineQueueState>(() => ({
    entries: getAllQueueEntries(),
    metrics: getQueueMetrics(),
    isReplayInProgress: isReplayInProgress.value,
  }))

  return {
    // State
    queueEntries,
    isReplayInProgress,
    cachedMessages,

    // Mutations
    enqueueTrackingWrite,
    enqueueDomainTrackingWrite,
    getQueueEntry,
    updateQueueEntryState,
    markEntrySynced,
    markEntryFailed,
    deleteQueueEntry,
    clearAllOfflineState,
    setReplayInProgress,
    setCachedMessages,

    // Queries
    getPendingEntries,
    getSyncingEntries,
    getEntriesByDomain,
    getQueueMetrics,
    hasFailedEntries,
    getAllQueueEntries,

    // Computed
    queueState,
  }
})
