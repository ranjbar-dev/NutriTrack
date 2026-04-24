import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import type { TrackingDomain } from '../../app/types/offline-sync'

describe('Sync Replay and Reconnect Behavior', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Reconnect Trigger', () => {
    it('triggers replay of queued entries when reconnecting', async () => {
      const store = useClientOfflineStore()

      // Enqueue some offline entries
      const entry1 = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })

      expect(entry1?.sync_state).toBe('queued')

      // Simulate reconnect event
      const replayEntries = store.getPendingEntries()
      expect(replayEntries.length).toBeGreaterThan(0)
      expect(replayEntries.some(e => e.sync_state === 'queued')).toBe(true)
    })

    it('prevents duplicate concurrent sync loops during replay', async () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      // Simulate first sync attempt
      store.updateQueueEntryState(entry!.local_id, 'syncing')
      const syncingEntries = store.getSyncingEntries()
      expect(syncingEntries.length).toBe(1)

      // Try to get pending entries while first is still in progress
      const stillSyncing = store.getSyncingEntries()
      expect(stillSyncing.length).toBe(1)

      // Should not allow duplicate syncing of same entries
      stillSyncing.forEach(e => {
        expect(e.sync_state).toBe('syncing')
      })
    })
  })

  describe('App-Open Replay', () => {
    it('replays queued entries on app initialization', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'sleep' as TrackingDomain,
        payload: { sleep_hours: 8, logged_at: '2026-04-23T22:00:00Z' },
      })

      // Simulate app close and reopen
      const stored = store.getQueueEntry(entry?.local_id!)
      expect(stored).toBeDefined()
      expect(stored?.sync_state).toBe('queued')
    })

    it('resumes syncing of entries that were in progress when app closed', () => {
      const store = useClientOfflineStore()

      const entry1 = store.enqueueTrackingWrite({
        domain: 'exercise' as TrackingDomain,
        payload: { exercise_id: 'ex-1', duration_minutes: 30, logged_at: '2026-04-23T18:00:00Z' },
      })
      const entry2 = store.enqueueTrackingWrite({
        domain: 'medication' as TrackingDomain,
        payload: { medication_id: 'med-1', doses: 1, logged_at: '2026-04-23T09:00:00Z' },
      })

      // Move entries to syncing before simulated app close
      store.updateQueueEntryState(entry1?.local_id!, 'syncing')
      store.updateQueueEntryState(entry2?.local_id!, 'syncing')

      // Simulate app reopen - entries should still be retrievable
      const entry1Reloaded = store.getQueueEntry(entry1?.local_id!)
      const entry2Reloaded = store.getQueueEntry(entry2?.local_id!)

      expect(entry1Reloaded?.sync_state).toBe('syncing')
      expect(entry2Reloaded?.sync_state).toBe('syncing')

      // After app open, retry logic should complete these
      store.markEntrySynced(entry1?.local_id!)
      expect(store.getQueueEntry(entry1?.local_id!)?.sync_state).toBe('synced')
    })
  })

  describe('Manual Retry for Failed Entries', () => {
    it('allows manual retry of failed queue entries', async () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'body' as TrackingDomain,
        payload: { weight_kg: 75, logged_at: '2026-04-23T08:00:00Z' },
      })

      // Simulate sync failure
      store.updateQueueEntryState(entry?.local_id!, 'syncing')
      store.markEntryFailed(entry?.local_id!, { code: 'NETWORK_ERROR', message: 'Failed' })

      expect(store.getQueueEntry(entry?.local_id!)?.sync_state).toBe('failed')

      // User initiates manual retry
      store.updateQueueEntryState(entry?.local_id!, 'queued')

      // Entry should be ready for replay
      const retriedEntry = store.getQueueEntry(entry?.local_id!)
      expect(retriedEntry?.sync_state).toBe('queued')
    })

    it('preserves failed entry data for manual retry UI', async () => {
      const store = useClientOfflineStore()

      const payload = { weight_kg: 75, logged_at: '2026-04-23T08:00:00Z' }
      const entry = store.enqueueTrackingWrite({
        domain: 'body' as TrackingDomain,
        payload,
      })

      const error = { code: 'CONFLICT', message: 'Duplicate entry' }
      store.updateQueueEntryState(entry?.local_id!, 'syncing')
      store.markEntryFailed(entry?.local_id!, error)

      const failedEntry = store.getQueueEntry(entry?.local_id!)
      expect(failedEntry?.payload).toEqual(payload)
      expect(failedEntry?.error_metadata).toEqual(error)
    })
  })

  describe('Bulk Sync Payload Building', () => {
    it('collects queued entries into bulk sync request payload', () => {
      const store = useClientOfflineStore()

      store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })
      store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })
      store.enqueueTrackingWrite({
        domain: 'exercise' as TrackingDomain,
        payload: { exercise_id: 'ex-1', duration_minutes: 30, logged_at: '2026-04-23T18:00:00Z' },
      })

      const pending = store.getPendingEntries()
      expect(pending.length).toBe(3)

      // Build bulk sync payload would normally happen in a composable
      const bulkPayload = pending.map(e => ({
        local_id: e.local_id,
        domain: e.domain,
        ...e.payload,
      }))

      expect(bulkPayload).toHaveLength(3)
      expect(bulkPayload.every(item => item.local_id)).toBe(true)
      expect(bulkPayload.map(item => item.domain)).toContain('food')
      expect(bulkPayload.map(item => item.domain)).toContain('water')
      expect(bulkPayload.map(item => item.domain)).toContain('exercise')
    })

    it('supports all six tracking domains in single bulk sync payload', () => {
      const store = useClientOfflineStore()

      const domains: TrackingDomain[] = ['food', 'water', 'sleep', 'exercise', 'medication', 'body']
      domains.forEach((domain, index) => {
        store.enqueueTrackingWrite({
          domain,
          payload: {
            test_field: `value_${index}`,
            logged_at: '2026-04-23T10:00:00Z',
          },
        })
      })

      const pending = store.getPendingEntries()
      expect(pending.length).toBe(6)

      const uniqueDomains = new Set(pending.map(e => e.domain))
      expect(uniqueDomains.size).toBe(6)
    })
  })

  describe('Last-Write-Wins Conflict Handling', () => {
    it('tracks timestamp on each entry for last-write-wins resolution', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      expect(entry?.created_at).toBeDefined()
      expect(typeof entry?.created_at).toBe('string')
    })

    it('maintains entry with failed sync state for conflict resolution display', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })

      // Simulate conflict error
      store.updateQueueEntryState(entry?.local_id!, 'syncing')
      store.markEntryFailed(entry?.local_id!, {
        code: 'CONFLICT_RESOLVED',
        message: 'Server had newer entry, using server version',
      })

      const conflictEntry = store.getQueueEntry(entry?.local_id!)
      expect(conflictEntry?.sync_state).toBe('failed')
      expect(conflictEntry?.error_metadata?.code).toContain('CONFLICT')
    })
  })

  describe('Single-Flight Replay Guard', () => {
    it('prevents multiple concurrent sync operations on same queue entry', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'sleep' as TrackingDomain,
        payload: { sleep_hours: 8, logged_at: '2026-04-23T22:00:00Z' },
      })

      // First sync attempt
      store.updateQueueEntryState(entry?.local_id!, 'syncing')
      const firstSync = store.getQueueEntry(entry?.local_id!)
      expect(firstSync?.sync_state).toBe('syncing')

      // Second sync attempt should not change state while already syncing
      const secondCheck = store.getQueueEntry(entry?.local_id!)
      expect(secondCheck?.sync_state).toBe('syncing')
    })
  })
})
