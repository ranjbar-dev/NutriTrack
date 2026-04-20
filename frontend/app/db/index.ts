// NutriTrack offline database — Dexie.js v4 with versioned schema
// D-05: six tables; D-11: syncQueue structure
import Dexie, { type Table } from 'dexie'

// ---- Table shapes ----

export interface CachedPlan {
  id: 1           // singleton — always 1
  plan_id: string
  fetched_at: string        // ISO timestamp of last successful fetch
  updated_hint: string | null  // plan updated_at from server (for staleness check)
  data: object    // full DietPlanResponse JSON blob
}

export interface CachedMessage {
  id: string      // server message UUID (or 'local_<uuid>' for optimistic sends)
  partner_id: string
  sent_at: string  // ISO timestamp (used for sort + 50-entry window)
  is_local: boolean  // true = not yet synced to server
  payload: object    // full MessageResponse JSON blob
}

// D-11: sync queue entry shape
export interface SyncQueueEntry {
  id?: number             // auto-increment primary key
  entity_type: string     // 'food_log' | 'water_log' | 'sleep_log' | 'exercise_log' | 'medication_log' | 'body_measurement' | 'lab_result_meta' | 'message'
  request_path: string    // e.g. '/client/food-logs'
  method: string          // 'POST'
  payload: string         // JSON.stringify(body) — no Blob here
  attachment_blob?: Blob  // D-08: binary for messages/lab uploads
  attachment_filename?: string
  attachment_mime?: string
  local_id: string        // UUID — passed to server for dedup (D-07 server-side)
  created_at: string      // ISO timestamp
  status: 'pending' | 'syncing' | 'failed' | 'synced'
  retry_count: number
  last_error: string | null
  next_attempt_at: string  // ISO timestamp — when to next attempt
}

export interface SyncMetaRecord {
  key: string   // e.g. 'plan_last_fetch', 'messages_last_fetch_{partnerId}'
  value: string // ISO timestamp or string
}

export interface NotificationPrefRecord {
  id: 1          // singleton
  data: object   // NotificationPrefs JSON blob
  updated_at: string
}

export interface UiStateRecord {
  key: string
  value: string
}

// ---- Database class ----

export class NutriTrackDB extends Dexie {
  activePlan!: Table<CachedPlan>
  messages!: Table<CachedMessage>
  syncQueue!: Table<SyncQueueEntry>
  syncMeta!: Table<SyncMetaRecord>
  notificationPreferences!: Table<NotificationPrefRecord>
  uiState!: Table<UiStateRecord>

  constructor() {
    super('NutriTrackDB')
    this.version(1).stores({
      activePlan: 'id',
      messages: 'id, partner_id, sent_at',    // compound indexes for partner queries
      syncQueue: '++id, status, entity_type, created_at, next_attempt_at',
      syncMeta: 'key',
      notificationPreferences: 'id',
      uiState: 'key',
    })
  }
}

// Lazy singleton — only instantiated client-side (Pitfall 3: SSR guard)
// Do NOT call new NutriTrackDB() at module load time in SSR paths.
// Access only inside composables, onMounted, or .client.ts files.
let _db: NutriTrackDB | null = null

export function getDb(): NutriTrackDB {
  if (!_db) _db = new NutriTrackDB()
  return _db
}

// Default export for simpler imports — lazy Proxy, SSR-safe
export const db = /* @__PURE__ */ new Proxy({} as NutriTrackDB, {
  get(_target, prop) {
    return getDb()[prop as keyof NutriTrackDB]
  },
})
