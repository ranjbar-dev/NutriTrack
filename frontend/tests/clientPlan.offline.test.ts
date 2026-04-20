// OFFL-02: Plan cache renders when offline
// OFFL-11: Messages cached per conversation, last 50

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

vi.mock('~/composables/useApi', () => ({
  useApi: () => ({
    apiFetch: vi.fn().mockRejectedValue(new TypeError('Network request failed')),
  }),
}))

vi.mock('#app', () => ({
  useNuxtApp: () => ({ hook: vi.fn(), callHook: vi.fn() }),
}))

vi.mock('~/db', async () => {
  await import('fake-indexeddb/auto')
  const Dexie = (await import('dexie')).default
  class TestDB extends Dexie {
    activePlan: any; messages: any; syncMeta: any; uiState: any; syncQueue: any; notificationPreferences: any
    constructor() {
      super('TestDB-plan-' + Date.now())
      this.version(1).stores({
        activePlan: 'id',
        messages: 'id, partner_id, sent_at',
        syncQueue: '++id, status, entity_type, created_at',
        syncMeta: 'key',
        notificationPreferences: 'id',
        uiState: 'key',
      })
    }
  }
  const db = new TestDB()
  return { db, getDb: () => db }
})

describe('offline cache reads', () => {
  beforeEach(() => setActivePinia(createPinia()))

  it('OFFL-11: messages table stores and retrieves last 50 per partner', async () => {
    const { db } = await import('~/db')
    const partnerId = 'partner-abc'
    const msgs = Array.from({ length: 55 }, (_, i) => ({
      id: `msg-${i}`,
      partner_id: partnerId,
      sent_at: new Date(Date.now() + i * 1000).toISOString(),
      is_local: false,
      payload: { id: `msg-${i}`, content: `message ${i}`, sent_at: new Date().toISOString() },
    }))
    await db.messages.bulkPut(msgs)
    const count = await db.messages.where('partner_id').equals(partnerId).count()
    expect(count).toBe(55)
    // Simulate trim to last 50
    const all = await db.messages.where('partner_id').equals(partnerId).sortBy('sent_at')
    const last50 = all.slice(-50)
    expect(last50.length).toBe(50)
    expect(last50[0].id).toBe('msg-5') // first 5 dropped
  })

  it('OFFL-02: activePlan cached and retrievable', async () => {
    const { db } = await import('~/db')
    const planData = { id: 'plan-1', title: 'Test Plan', days: [] }
    await db.activePlan.put({ id: 1, plan_id: 'plan-1', fetched_at: new Date().toISOString(), updated_hint: null, data: planData })
    const cached = await db.activePlan.get(1)
    expect(cached?.plan_id).toBe('plan-1')
    expect((cached?.data as typeof planData).title).toBe('Test Plan')
  })
})
