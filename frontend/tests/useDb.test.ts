import 'fake-indexeddb/auto'
import { describe, it, expect, beforeEach } from 'vitest'
import { NutriTrackDB } from '../app/db/index'

// Use a fresh in-memory DB for each test (fake-indexeddb resets between test runs)
let testDb: NutriTrackDB

beforeEach(async () => {
  testDb = new NutriTrackDB()
  // Clear all tables for isolation
  await testDb.activePlan.clear()
  await testDb.messages.clear()
  await testDb.syncQueue.clear()
  await testDb.syncMeta.clear()
  await testDb.notificationPreferences.clear()
  await testDb.uiState.clear()
})

describe('NutriTrackDB schema', () => {
  it('OFFL-05: opens with all 6 required tables', async () => {
    expect(testDb.activePlan).toBeDefined()
    expect(testDb.messages).toBeDefined()
    expect(testDb.syncQueue).toBeDefined()
    expect(testDb.syncMeta).toBeDefined()
    expect(testDb.notificationPreferences).toBeDefined()
    expect(testDb.uiState).toBeDefined()
  })

  it('activePlan stores and retrieves singleton plan', async () => {
    await testDb.activePlan.put({
      id: 1,
      plan_id: 'plan-uuid-123',
      fetched_at: new Date().toISOString(),
      updated_hint: null,
      data: { title: 'Test Plan' },
    })
    const cached = await testDb.activePlan.get(1)
    expect(cached?.plan_id).toBe('plan-uuid-123')
  })

  it('syncQueue auto-increments id', async () => {
    const id1 = await testDb.syncQueue.add({
      entity_type: 'food_log',
      request_path: '/client/food-logs',
      method: 'POST',
      payload: '{}',
      local_id: crypto.randomUUID(),
      created_at: new Date().toISOString(),
      status: 'pending',
      retry_count: 0,
      last_error: null,
      next_attempt_at: new Date().toISOString(),
    })
    const id2 = await testDb.syncQueue.add({
      entity_type: 'water_log',
      request_path: '/client/water-logs',
      method: 'POST',
      payload: '{}',
      local_id: crypto.randomUUID(),
      created_at: new Date().toISOString(),
      status: 'pending',
      retry_count: 0,
      last_error: null,
      next_attempt_at: new Date().toISOString(),
    })
    expect(id2).toBeGreaterThan(id1 as number)
  })

  it('OFFL-12: iOS eviction detection — empty activePlan after recent fetch', async () => {
    // Simulate: plan was fetched recently (syncMeta says so), but activePlan is empty
    await testDb.syncMeta.put({ key: 'plan_last_fetch', value: new Date().toISOString() })
    const planCount = await testDb.activePlan.count()
    const lastFetchRecord = await testDb.syncMeta.get('plan_last_fetch')

    if (lastFetchRecord) {
      const fetchedAt = new Date(lastFetchRecord.value)
      const oneHourAgo = new Date(Date.now() - 3600_000)
      const recentFetch = fetchedAt > oneHourAgo
      const evictionDetected = planCount === 0 && recentFetch
      expect(evictionDetected).toBe(true)
    }
  })
})
