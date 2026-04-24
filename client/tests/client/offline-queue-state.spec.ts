import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import type { TrackingQueueEntry, TrackingDomain } from '../../app/types/offline-sync'

describe('Client Offline Queue State', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Queue Entry State Transitions', () => {
    it('creates a new offline write with state queued and local_id', () => {
      const store = useClientOfflineStore()
      const entry = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: {
          food_id: 'food-1',
          quantity_grams: 100,
          consumed_at: '2026-04-23',
        },
      })

      expect(entry).toBeDefined()
      expect(entry?.local_id).toBeDefined()
      expect(entry?.sync_state).toBe('queued')
      expect(entry?.domain).toBe('food')
    })

    it('transitions entry from queued to syncing during replay', async () => {
      const store = useClientOfflineStore()
      const entry = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: {
          amount_ml: 250,
          logged_at: '2026-04-23T10:00:00Z',
        },
      })

      expect(entry?.sync_state).toBe('queued')

      // Simulate transition to syncing
      if (entry) {
        store.updateQueueEntryState(entry.local_id, 'syncing')
        const updated = store.getQueueEntry(entry.local_id)
        expect(updated?.sync_state).toBe('syncing')
      }
    })

    it('transitions entry from syncing to synced on successful sync', async () => {
      const store = useClientOfflineStore()
      const entry = store.enqueueTrackingWrite({
        domain: 'sleep' as TrackingDomain,
        payload: {
          sleep_hours: 8,
          logged_at: '2026-04-23T22:00:00Z',
        },
      })

      if (entry) {
        store.updateQueueEntryState(entry.local_id, 'syncing')
        store.markEntrySynced(entry.local_id)
        const synced = store.getQueueEntry(entry.local_id)
        expect(synced?.sync_state).toBe('synced')
      }
    })

    it('transitions entry to failed on sync error and preserves error metadata', async () => {
      const store = useClientOfflineStore()
      const entry = store.enqueueTrackingWrite({
        domain: 'exercise' as TrackingDomain,
        payload: {
          exercise_id: 'ex-1',
          duration_minutes: 30,
          logged_at: '2026-04-23T18:00:00Z',
        },
      })

      if (entry) {
        const error = { code: 'NETWORK_ERROR', message: 'Connection failed' }
        store.updateQueueEntryState(entry.local_id, 'syncing')
        store.markEntryFailed(entry.local_id, error)
        const failed = store.getQueueEntry(entry.local_id)
        expect(failed?.sync_state).toBe('failed')
        expect(failed?.error_metadata).toBeDefined()
      }
    })
  })

  describe('Queue Metrics and Visibility', () => {
    it('computes correct aggregate sync state counts', () => {
      const store = useClientOfflineStore()

      // Add entries in different states
      const entry1 = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })
      const entry2 = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      store.updateQueueEntryState(entry2?.local_id!, 'syncing')

      const metrics = store.getQueueMetrics()
      expect(metrics.queued).toBe(1)
      expect(metrics.syncing).toBe(1)
      expect(metrics.synced).toBe(0)
      expect(metrics.failed).toBe(0)
      expect(metrics.total).toBe(2)
    })

    it('provides list of all queued and failed entries for UI rendering', () => {
      const store = useClientOfflineStore()

      store.enqueueTrackingWrite({
        domain: 'medication' as TrackingDomain,
        payload: { medication_id: 'med-1', doses: 1, logged_at: '2026-04-23T09:00:00Z' },
      })
      const entry2 = store.enqueueTrackingWrite({
        domain: 'body' as TrackingDomain,
        payload: { weight_kg: 75, logged_at: '2026-04-23T08:00:00Z' },
      })

      store.markEntryFailed(entry2?.local_id!, { code: 'SYNC_FAILED', message: 'Server error' })

      const pending = store.getPendingEntries()
      expect(pending.length).toBe(2)
      expect(pending.some(e => e.sync_state === 'queued')).toBe(true)
      expect(pending.some(e => e.sync_state === 'failed')).toBe(true)
    })
  })

  describe('Logout Cleanup', () => {
    it('deletes all queue entries on logout', () => {
      const store = useClientOfflineStore()

      store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })
      store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      expect(store.getQueueMetrics().total).toBe(2)

      store.clearAllOfflineState()

      expect(store.getQueueMetrics().total).toBe(0)
      expect(store.getPendingEntries().length).toBe(0)
    })

    it('resets sync state counters to zero on logout', () => {
      const store = useClientOfflineStore()

      const entry1 = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })
      const entry2 = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      store.updateQueueEntryState(entry1?.local_id!, 'syncing')
      store.markEntrySynced(entry2?.local_id!)

      const beforeCleanup = store.getQueueMetrics()
      expect(beforeCleanup.syncing).toBe(1)
      expect(beforeCleanup.synced).toBe(1)

      store.clearAllOfflineState()

      const afterCleanup = store.getQueueMetrics()
      expect(afterCleanup.queued).toBe(0)
      expect(afterCleanup.syncing).toBe(0)
      expect(afterCleanup.synced).toBe(0)
      expect(afterCleanup.failed).toBe(0)
    })
  })

  describe('All Six Tracking Domains Supported', () => {
    it.each(['food', 'water', 'sleep', 'exercise', 'medication', 'body'] as TrackingDomain[])(
      'accepts domain %s in queue entry',
      (domain) => {
        const store = useClientOfflineStore()
        const entry = store.enqueueTrackingWrite({
          domain,
          payload: {
            test_field: 'value',
            logged_at: '2026-04-23T10:00:00Z',
          },
        })
        expect(entry?.domain).toBe(domain)
      }
    )
  })
})
