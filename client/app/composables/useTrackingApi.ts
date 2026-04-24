/**
 * useTrackingApi Composable
 * Typed tracking API composable for individual domain endpoints and bulk sync
 * Used by offline queue orchestration and UI tracking entry flows
 */

import { useAsyncData, useFetch } from '#app'
import type {
  BulkTrackingSyncRequest,
  BulkTrackingSyncResponse,
  TrackingHistoryResponse,
  CreateFoodTrackingRequest,
  CreateWaterTrackingRequest,
  CreateSleepTrackingRequest,
  CreateExerciseTrackingRequest,
  CreateMedicationTrackingRequest,
  CreateBodyTrackingRequest,
} from '~/app/types/tracking'
import type { TrackingDomain, TrackingQueueEntry } from '~/app/types/offline-sync'

const baseUrl = '/api/v1'

export const useTrackingApi = () => {
  /**
   * Submit a food tracking entry
   */
  async function submitFoodTracking(request: CreateFoodTrackingRequest) {
    return useFetch(`${baseUrl}/tracking/food`, {
      method: 'POST',
      body: request,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  /**
   * Submit a water tracking entry
   */
  async function submitWaterTracking(request: CreateWaterTrackingRequest) {
    return useFetch(`${baseUrl}/tracking/water`, {
      method: 'POST',
      body: request,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  /**
   * Submit a sleep tracking entry
   */
  async function submitSleepTracking(request: CreateSleepTrackingRequest) {
    return useFetch(`${baseUrl}/tracking/sleep`, {
      method: 'POST',
      body: request,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  /**
   * Submit an exercise tracking entry
   */
  async function submitExerciseTracking(request: CreateExerciseTrackingRequest) {
    return useFetch(`${baseUrl}/tracking/exercise`, {
      method: 'POST',
      body: request,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  /**
   * Submit a medication tracking entry
   */
  async function submitMedicationTracking(request: CreateMedicationTrackingRequest) {
    return useFetch(`${baseUrl}/tracking/medication`, {
      method: 'POST',
      body: request,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  /**
   * Submit a body measurement tracking entry
   */
  async function submitBodyTracking(request: CreateBodyTrackingRequest) {
    return useFetch(`${baseUrl}/tracking/body`, {
      method: 'POST',
      body: request,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  /**
   * Bulk sync tracking entries from offline queue
   * Uses idempotent local_id key for replay safety
   */
  async function bulkSyncTracking(
    entries: TrackingQueueEntry[]
  ): Promise<BulkTrackingSyncResponse | null> {
    if (entries.length === 0) return null

    const bulkRequest: BulkTrackingSyncRequest = {
      entries: entries.map(e => ({
        local_id: e.local_id,
        domain: e.domain,
        ...e.payload,
      })) as any,
    }

    try {
      const { data, error } = await useFetch<BulkTrackingSyncResponse>(
        `${baseUrl}/tracking/sync`,
        {
          method: 'POST',
          body: bulkRequest,
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )

      if (error.value) {
        console.error('Bulk sync error:', error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error('Bulk sync exception:', err)
      return null
    }
  }

  /**
   * Fetch tracking history for a domain
   */
  async function getTrackingHistory(
    domain?: TrackingDomain,
    page = 1,
    pageSize = 20
  ): Promise<TrackingHistoryResponse | null> {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    })

    if (domain) {
      params.append('domain', domain)
    }

    try {
      const { data, error } = await useFetch<TrackingHistoryResponse>(
        `${baseUrl}/tracking/history?${params.toString()}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )

      if (error.value) {
        console.error('History fetch error:', error.value)
        return null
      }

      return data.value
    } catch (err) {
      console.error('History fetch exception:', err)
      return null
    }
  }

  /**
   * Dispatch a tracking entry to the correct endpoint based on domain
   * Used for individual submissions and retry logic
   */
  async function submitTracking(entry: TrackingQueueEntry) {
    const request = {
      local_id: entry.local_id,
      ...entry.payload,
    }

    switch (entry.domain) {
      case 'food':
        return submitFoodTracking(request as CreateFoodTrackingRequest)
      case 'water':
        return submitWaterTracking(request as CreateWaterTrackingRequest)
      case 'sleep':
        return submitSleepTracking(request as CreateSleepTrackingRequest)
      case 'exercise':
        return submitExerciseTracking(request as CreateExerciseTrackingRequest)
      case 'medication':
        return submitMedicationTracking(request as CreateMedicationTrackingRequest)
      case 'body':
        return submitBodyTracking(request as CreateBodyTrackingRequest)
      default:
        throw new Error(`Unknown tracking domain: ${entry.domain}`)
    }
  }

  return {
    submitFoodTracking,
    submitWaterTracking,
    submitSleepTracking,
    submitExerciseTracking,
    submitMedicationTracking,
    submitBodyTracking,
    bulkSyncTracking,
    getTrackingHistory,
    submitTracking,
  }
}
