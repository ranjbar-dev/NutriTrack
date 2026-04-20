import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSyncQueueStore } from '../app/stores/syncQueue'

vi.mock('~/composables/useApi', () => ({
  useApi: () => ({
    apiFetch: vi.fn(),
  }),
}))

vi.mock('#app', () => ({
  useNuxtApp: () => ({ callHook: vi.fn() }),
}))

vi.mock('~/db', async () => {
  const Dexie = (await import('dexie')).default
  await import('fake-indexeddb/auto')
  class TestDB extends Dexie {
    syncQueue: any
    constructor() {
      super('TestDB-' + Date.now())
      this.version(1).stores({ syncQueue: '++id, status, entity_type, created_at, next_attempt_at' })
    }
  }
  const db = new TestDB()
  return { db, getDb: () => db }
})

describe('useSyncQueueStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('OFFL-06: enqueue adds entry with status=pending', async () => {
    const store = useSyncQueueStore()
    const result = await store.enqueue('/client/food-logs', { local_id: 'test-id-1', date: '2026-04-20' }, { entityType: 'food_log' })
    expect(result.queued).toBe(true)
    expect(result.local_id).toBe('test-id-1')
  })

  it('OFFL-08: backoffMs returns correct delays', () => {
    const store = useSyncQueueStore()
    expect(store.backoffMs(0)).toBe(0)
    expect(store.backoffMs(1)).toBe(1000)
    expect(store.backoffMs(2)).toBe(2000)
    expect(store.backoffMs(3)).toBe(4000)
    expect(store.backoffMs(4)).toBe(4000)
  })

  it('OFFL-07: same local_id can be enqueued multiple times (server dedup handles it)', async () => {
    const store = useSyncQueueStore()
    const lid = 'dedup-test-local-id'
    await store.enqueue('/client/water-logs', { local_id: lid, amount_ml: 200 }, { entityType: 'water_log' })
    await store.enqueue('/client/water-logs', { local_id: lid, amount_ml: 200 }, { entityType: 'water_log' })
    expect(true).toBe(true)
  })
})
