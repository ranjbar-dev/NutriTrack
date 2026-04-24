import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import {
  getFailureGuidance,
  retryAllFailedEntries,
  retryFailedEntry,
} from '../../app/lib/tracking/retry'

describe('Manual Retry Loop for Failed Sync', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('keeps failed entries visible and supports retry-single plus retry-all actions', () => {
    const store = useClientOfflineStore()

    const first = store.enqueueTrackingWrite({
      domain: 'food',
      payload: { food_id: 'food-1', quantity_grams: 120, consumed_at: '2026-04-23T07:00:00.000Z' },
    })
    const second = store.enqueueTrackingWrite({
      domain: 'water',
      payload: { amount_ml: 500, logged_at: '2026-04-23T08:00:00.000Z' },
    })

    store.markEntryFailed(first!.local_id, { code: 'NETWORK', message: 'offline' })
    store.markEntryFailed(second!.local_id, { code: 'NETWORK', message: 'offline' })

    expect(store.getPendingEntries().filter((entry) => entry.sync_state === 'failed')).toHaveLength(2)

    const retriedOne = retryFailedEntry(store, first!.local_id)
    expect(retriedOne).toBe(true)
    expect(store.getQueueEntry(first!.local_id)?.sync_state).toBe('queued')

    const retriedAllCount = retryAllFailedEntries(store)
    expect(retriedAllCount).toBe(1)
    expect(store.getQueueEntry(second!.local_id)?.sync_state).toBe('queued')
  })

  it('preserves actionable guidance for entries that remain failed after retry attempts', () => {
    const store = useClientOfflineStore()

    const entry = store.enqueueTrackingWrite({
      domain: 'body',
      payload: { weight_kg: 74.2, logged_at: '2026-04-23T08:30:00.000Z' },
    })

    store.markEntryFailed(entry!.local_id, { code: 'SERVER', message: 'temporary issue' })
    store.markEntryFailed(entry!.local_id, { code: 'SERVER', message: 'temporary issue' })
    store.markEntryFailed(entry!.local_id, { code: 'SERVER', message: 'temporary issue' })

    const failedEntry = store.getQueueEntry(entry!.local_id)
    expect(failedEntry?.sync_state).toBe('failed')

    const guidance = getFailureGuidance(failedEntry!)
    expect(guidance).toContain('اینترنت')
  })
})
