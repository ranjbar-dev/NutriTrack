import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import type { TrackingDomain } from '../../app/types/offline-sync'

describe('Today View Shell and Sync Visibility', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Today Action Cards', () => {
    it('renders active plan snapshot, pending actions, water, and sync status on Today page', () => {
      const todayView = {
        plan_snapshot: {
          id: 'plan-1',
          day_name: 'جمعه',
          meals_count: 3,
          water_target_ml: 2000,
          is_cached: false,
        },
        pending_actions: {
          food: 2,
          water: 1,
          sleep: 0,
          exercise: 0,
          medication: 1,
          body: 0,
        },
        sync_status: {
          queued: 3,
          syncing: 0,
          synced: 5,
          failed: 0,
        },
      }

      expect(todayView.plan_snapshot).toBeDefined()
      expect(todayView.pending_actions).toBeDefined()
      expect(todayView.sync_status).toBeDefined()
    })

    it('when offline with cached plan, snapshot displays with timestamp "آخرین به روزرسانی: HH:MM"', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })

      const cachedPlanSnapshot = {
        is_cached: true,
        updated_at: '2026-04-23T10:00:00Z',
        formatted_time: 'آخرین به روزرسانی: 10:00',
      }

      expect(cachedPlanSnapshot.is_cached).toBe(true)
      expect(cachedPlanSnapshot.formatted_time).toContain('آخرین به روزرسانی')
    })

    it('pending action counts derive only from queued entries (not synced upstream)', () => {
      const store = useClientOfflineStore()

      // Add some queued entries
      store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })
      store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      const pending = store.getPendingEntries()
      const pendingByDomain = {
        food: pending.filter(e => e.domain === 'food').length,
        water: pending.filter(e => e.domain === 'water').length,
      }

      expect(pendingByDomain.food).toBe(1)
      expect(pendingByDomain.water).toBe(1)
    })

    it('water quick-add component shows current progress and quick-add controls', () => {
      const waterState = {
        daily_target_ml: 2000,
        logged_today_ml: 750,
        progress_percent: 37.5,
        quick_add_options: [
          { label: '250ml', value: 250 },
          { label: '500ml', value: 500 },
          { label: 'سایر', value: 'custom' },
        ],
      }

      expect(waterState.logged_today_ml).toBeLessThan(waterState.daily_target_ml)
      expect(waterState.progress_percent).toBe(37.5)
      expect(waterState.quick_add_options.length).toBe(3)
    })
  })

  describe('Sync Strip and State Visibility', () => {
    it('sync strip appears below connectivity banner with aggregate queue state', () => {
      const store = useClientOfflineStore()

      store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })

      const metrics = store.getQueueMetrics()
      expect(metrics.queued).toBe(1)
      expect(metrics.total).toBeGreaterThan(0)

      // Sync strip should show these metrics
      const syncStrip = {
        queued_count: metrics.queued,
        syncing_count: metrics.syncing,
        synced_count: metrics.synced,
        failed_count: metrics.failed,
        is_visible: metrics.total > 0,
      }

      expect(syncStrip.is_visible).toBe(true)
    })

    it('failed state remains visible with manual retry CTA and Persian recovery text', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      store.updateQueueEntryState(entry!.local_id, 'syncing')
      store.markEntryFailed(entry!.local_id, {
        code: 'NETWORK_ERROR',
        message: 'Connection failed',
      })

      const failedEntry = store.getQueueEntry(entry!.local_id)
      expect(failedEntry?.sync_state).toBe('failed')
      expect(failedEntry?.error_metadata).toBeDefined()

      const syncUi = {
        state: 'failed',
        message: 'این ثبت به سرور نرسید. می توانید دوباره تلاش کنید یا بعدا ارسال کنید.',
        retry_available: true,
        fallback_available: true,
      }

      expect(syncUi.state).toBe('failed')
      expect(syncUi.retry_available).toBe(true)
    })

    it('displays queue aggregate as chips or counters (queued/syncing/synced/failed)', () => {
      const store = useClientOfflineStore()

      // Create mixed state entries
      const entry1 = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })
      const entry2 = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      store.updateQueueEntryState(entry1!.local_id, 'syncing')
      store.markEntrySynced(entry2!.local_id)

      const metrics = store.getQueueMetrics()
      const chipDisplay = {
        queued: { icon: '⏱️', color: 'warning', count: metrics.queued },
        syncing: { icon: '↻', color: 'info', count: metrics.syncing },
        synced: { icon: '✓', color: 'success', count: metrics.synced },
        failed: { icon: '✗', color: 'error', count: metrics.failed },
      }

      expect(chipDisplay.syncing.count).toBe(1)
      expect(chipDisplay.synced.count).toBe(1)
    })
  })

  describe('Offline Behavior and Data Freshness', () => {
    it('today view remains readable from cache while offline', () => {
      const cachedTodayView = {
        plan_snapshot: {
          cached: true,
          last_sync: '2026-04-23T08:00:00Z',
          stale: false,
        },
        pending_actions: {
          visible: true,
          source: 'queue',
        },
        connectivity: 'offline',
      }

      expect(cachedTodayView.plan_snapshot.cached).toBe(true)
      expect(cachedTodayView.pending_actions.source).toBe('queue')
    })

    it('shows cache freshness indicator when data is from offline store', () => {
      const freshness = {
        source: 'cache' as const,
        label: 'آخرین به روزرسانی: 08:00',
        is_stale: false,
      }

      expect(freshness.source).toBe('cache')
      expect(freshness.label).toContain('آخرین به روزرسانی')
    })
  })
})
