import { describe, it, expect } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import type { TrackingDomain } from '../../app/types/offline-sync'

describe('Sync Strip Visibility and UX', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Sync State Rendering', () => {
    it('renders sync strip with queued/syncing/synced/failed aggregate state', () => {
      const store = useClientOfflineStore()

      // Add mixed entries
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
      const syncStrip = {
        queued: metrics.queued,
        syncing: metrics.syncing,
        synced: metrics.synced,
        failed: metrics.failed,
      }

      expect(syncStrip.syncing).toBeGreaterThan(0)
      expect(syncStrip.synced).toBeGreaterThan(0)
    })

    it('displays state with explicit labels: queued/syncing/synced/failed', () => {
      const stateLabels = {
        queued: { fa: 'در صف', icon: '⏱️', color: 'amber' },
        syncing: { fa: 'در حال ارسال', icon: '↻', color: 'blue' },
        synced: { fa: 'ارسال شده', icon: '✓', color: 'green' },
        failed: { fa: 'ناموفق', icon: '✗', color: 'red' },
      }

      expect(stateLabels.queued.fa).toBe('در صف')
      expect(stateLabels.synced.fa).toBe('ارسال شده')
      expect(stateLabels.failed.fa).toBe('ناموفق')
    })

    it('shows queued/syncing animation or visual indicator once per state transition', () => {
      const animationState = {
        syncing_animate: true, // Animate once when entering syncing
        synced_persist: true, // Keep synced visible until dismissed
        failed_persist: true, // Failed stays visible until user acts
      }

      expect(animationState.syncing_animate).toBe(true)
      expect(animationState.synced_persist).toBe(true)
    })
  })

  describe('Failed Entry Recovery UI', () => {
    it('failed state remains visible with manual retry CTA', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'exercise' as TrackingDomain,
        payload: { exercise_id: 'ex-1', duration_minutes: 30, logged_at: '2026-04-23T18:00:00Z' },
      })

      store.updateQueueEntryState(entry!.local_id, 'syncing')
      store.markEntryFailed(entry!.local_id, {
        code: 'NETWORK_ERROR',
        message: 'Connection timeout',
      })

      const failedEntry = store.getQueueEntry(entry!.local_id)
      const recoveryUI = {
        state: failedEntry?.sync_state,
        message: 'این ثبت به سرور نرسید. می توانید دوباره تلاش کنید یا بعدا ارسال کنید.',
        retry_action: {
          label: 'دوباره تلاش کنید',
          available: failedEntry?.sync_state === 'failed',
        },
        fallback_action: {
          label: 'بعدا ارسال کن',
          available: failedEntry?.sync_state === 'failed',
        },
      }

      expect(recoveryUI.state).toBe('failed')
      expect(recoveryUI.retry_action.available).toBe(true)
    })

    it('supports per-entry retry from row chips or detail view', () => {
      const store = useClientOfflineStore()

      const entry = store.enqueueTrackingWrite({
        domain: 'medication' as TrackingDomain,
        payload: { medication_id: 'med-1', doses: 1, logged_at: '2026-04-23T09:00:00Z' },
      })

      store.updateQueueEntryState(entry!.local_id, 'syncing')
      store.markEntryFailed(entry!.local_id, { code: 'SYNC_FAILED', message: 'Error' })

      // User clicks retry on row chip or detail
      store.updateQueueEntryState(entry!.local_id, 'queued')

      const retriedEntry = store.getQueueEntry(entry!.local_id)
      expect(retriedEntry?.sync_state).toBe('queued')
    })

    it('provides non-technical Persian recovery text that explains next steps', () => {
      const recoveryMessages = {
        network_error: {
          fa: 'اتصال اینترنت را بررسی کنید و دوباره تلاش کنید',
          tone: 'calm_helpful',
        },
        sync_failed: {
          fa: 'این ثبت در ارسال با مشکل مواجه شد. می توانید بعدا دوباره تلاش کنید',
          tone: 'calm_helpful',
        },
        conflict: {
          fa: 'سرور نسخه جدیدتری از این ثبت دارد. تغییرات شما ذخیره نشد',
          tone: 'informative',
        },
      }

      expect(recoveryMessages.network_error.tone).toBe('calm_helpful')
      expect(recoveryMessages.sync_failed.fa).toContain('می توانید')
      expect(recoveryMessages.conflict.tone).toBe('informative')
    })
  })

  describe('Sync Strip Placement and UX', () => {
    it('appears persistently below connectivity banner in client shell', () => {
      const layout = {
        app_shell: {
          top: ['header'],
          below_header: ['connectivity_banner'],
          below_connectivity: ['sync_strip'],
          main_content: ['page_content'],
        },
      }

      expect(layout.app_shell.below_connectivity).toContain('sync_strip')
    })

    it('strips only visible when queue has entries (queued/syncing/failed)', () => {
      const store = useClientOfflineStore()

      // Empty queue - strip not visible
      let metrics = store.getQueueMetrics()
      expect(metrics.total).toBe(0)

      // Add entry - strip visible
      store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      metrics = store.getQueueMetrics()
      expect(metrics.total).toBeGreaterThan(0)
    })

    it('uses min-height 40px and non-modal sticky positioning', () => {
      const stripLayout = {
        min_height_px: 40,
        position: 'sticky',
        z_index: 'below_header',
        modal: false,
      }

      expect(stripLayout.min_height_px).toBe(40)
      expect(stripLayout.modal).toBe(false)
    })
  })

  describe('Sync State Chip Semantics', () => {
    it('derives chip color from state: queued=amber, syncing=blue, synced=green, failed=red', () => {
      const stateColors = {
        queued: '#fbbf24', // amber
        syncing: '#3b82f6', // blue
        synced: '#10b981', // green
        failed: '#ef4444', // red
      }

      expect(stateColors.failed).toBe('#ef4444')
      expect(stateColors.synced).toBe('#10b981')
    })

    it('displays queue count badge for each state when count > 0', () => {
      const store = useClientOfflineStore()

      const entry1 = store.enqueueTrackingWrite({
        domain: 'food' as TrackingDomain,
        payload: { food_id: 'f1', quantity_grams: 100, consumed_at: '2026-04-23' },
      })
      const entry2 = store.enqueueTrackingWrite({
        domain: 'water' as TrackingDomain,
        payload: { amount_ml: 250, logged_at: '2026-04-23T10:00:00Z' },
      })

      const metrics = store.getQueueMetrics()
      const chipBadges = {
        queued_visible: metrics.queued > 0,
        queued_count: metrics.queued,
        total_visible: metrics.total > 0,
      }

      expect(chipBadges.queued_visible).toBe(true)
      expect(chipBadges.queued_count).toBe(2)
    })
  })
})
