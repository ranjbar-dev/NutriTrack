import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import {
  buildTrackingProgressSummary,
  formatTrackingTimestampFa,
  trackingDomainLabel,
} from '../../app/lib/tracking/history'

describe('Tracking History and Progress Surface', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders row-level sync states with Persian-formatted timestamps', () => {
    const store = useClientOfflineStore()

    const first = store.enqueueTrackingWrite({
      domain: 'food',
      payload: { food_id: 'food-1', quantity_grams: 120, consumed_at: '2026-04-23T07:00:00.000Z' },
    })
    const second = store.enqueueTrackingWrite({
      domain: 'water',
      payload: { amount_ml: 250, logged_at: '2026-04-23T08:00:00.000Z' },
    })

    store.updateQueueEntryState(first!.local_id, 'syncing')
    store.markEntryFailed(second!.local_id, { code: 'NETWORK', message: 'offline' })

    const rows = store.getAllQueueEntries()
    const chips = rows.map((entry) => entry.sync_state)

    expect(chips).toContain('syncing')
    expect(chips).toContain('failed')

    const formatted = formatTrackingTimestampFa(rows[0].created_at)
    expect(/[۰-۹]/.test(formatted)).toBe(true)
    expect(trackingDomainLabel(rows[0].domain).length).toBeGreaterThan(0)
  })

  it('builds lightweight progress summary from available v1 queue data', () => {
    const store = useClientOfflineStore()

    store.enqueueTrackingWrite({
      domain: 'water',
      payload: { amount_ml: 500, logged_at: new Date().toISOString() },
    })
    store.enqueueTrackingWrite({
      domain: 'water',
      payload: { amount_ml: 750, logged_at: new Date().toISOString() },
    })
    store.enqueueTrackingWrite({
      domain: 'food',
      payload: { food_id: 'food-2', quantity_grams: 160, consumed_at: new Date().toISOString() },
    })

    const summary = buildTrackingProgressSummary(store.getAllQueueEntries(), 2000)

    expect(summary.waterTodayMl).toBe(1250)
    expect(summary.waterCompletionPercent).toBe(63)
    expect(summary.recentDaysWithEntries).toBeGreaterThanOrEqual(1)
  })
})
