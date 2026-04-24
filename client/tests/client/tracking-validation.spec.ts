import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import {
  computeSleepDurationHours,
  createBodyPayload,
  createSleepPayload,
} from '../../app/lib/tracking/entry'

describe('Tracking Validation Rules', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('supports overnight sleep range and computes duration when explicit value is absent', () => {
    const duration = computeSleepDurationHours('23:40', '06:10')
    expect(duration).toBeCloseTo(6.5, 1)

    const payload = createSleepPayload({
      sleepStart: '23:40',
      wakeTime: '06:10',
      loggedAt: '2026-04-23T06:10:00.000Z',
      quality: 'good',
    })

    expect(payload.ok).toBe(true)
    expect(payload.payload?.sleep_hours).toBeCloseTo(6.5, 1)
  })

  it('allows partial body measurements and still creates valid queued writes', () => {
    const store = useClientOfflineStore()

    const payload = createBodyPayload({
      loggedAt: '2026-04-23T08:30:00.000Z',
      weight_kg: 74.2,
    })

    expect(payload.ok).toBe(true)

    const entry = store.enqueueDomainTrackingWrite({
      domain: 'body',
      payload: payload.payload!,
    })

    expect(entry?.sync_state).toBe('queued')
    expect(entry?.payload.weight_kg).toBe(74.2)
    expect(entry?.payload.waist_cm).toBeUndefined()
  })
})
