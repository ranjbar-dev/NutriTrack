import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useClientOfflineStore } from '../../app/stores/client-offline'
import { createFoodPayload, createWaterPayload } from '../../app/lib/tracking/entry'

describe('Tracking Entry Offline Queue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('creates local_id-backed queued food entry with immediate queue feedback', () => {
    const store = useClientOfflineStore()
    const payload = createFoodPayload({
      foodId: 'food-omelette',
      quantityGrams: 140,
      consumedAt: '2026-04-23T07:00:00.000Z',
    })

    expect(payload.ok).toBe(true)

    const entry = store.enqueueDomainTrackingWrite({
      domain: 'food',
      payload: payload.payload!,
    })

    expect(entry?.local_id).toBeTruthy()
    expect(entry?.sync_state).toBe('queued')

    const metrics = store.getQueueMetrics()
    expect(metrics.queued).toBe(1)
    expect(metrics.total).toBe(1)
  })

  it('supports repeated one-thumb water quick-add writes in queue-first mode', () => {
    const store = useClientOfflineStore()

    const quickAdds = [250, 500, 250]
    quickAdds.forEach((amount, index) => {
      const payload = createWaterPayload({
        amountMl: amount,
        loggedAt: `2026-04-23T10:0${index}:00.000Z`,
      })

      expect(payload.ok).toBe(true)

      const entry = store.enqueueDomainTrackingWrite({
        domain: 'water',
        payload: payload.payload!,
      })

      expect(entry?.sync_state).toBe('queued')
    })

    const waterEntries = store.getEntriesByDomain('water')
    expect(waterEntries).toHaveLength(3)
    expect(waterEntries.every((entry) => entry.sync_state === 'queued')).toBe(true)
  })
})
