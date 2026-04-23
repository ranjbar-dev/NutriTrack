/**
 * Client Sync Plugin
 * Orchestrates offline queue replay on reconnect and app-open events
 * Handles single-flight replay guard to prevent duplicate sync loops
 */

import { defineNuxtPlugin } from '#app'
import { useClientOfflineStore } from '~/app/stores/client-offline'
import { useTrackingApi } from '~/app/composables/useTrackingApi'
import { useMessagingApi } from '~/app/composables/useMessagingApi'
import { useAuthApi } from '~/app/composables/useAuthApi'
import { useAuthSessionStore } from '~/app/stores/auth-session'

export default defineNuxtPlugin((nuxtApp) => {
  // Skip plugin in SSR
  if (process.server) return

  const offlineStore = useClientOfflineStore()
  const trackingApi = useTrackingApi()
  const authSessionStore = useAuthSessionStore()

  let replayScheduled = false
  let reconnectTimeout: ReturnType<typeof setTimeout> | null = null

  /**
   * Main replay orchestrator
   * Processes pending queue entries and syncs with server
   */
  async function performQueueReplay() {
    // Single-flight guard: prevent concurrent replays
    if (offlineStore.isReplayInProgress) {
      console.debug('[Client Sync] Replay already in progress, skipping')
      return
    }

    // Check auth state - don't replay if not authenticated
    if (!authSessionStore.isAuthenticated || authSessionStore.userRole !== 'client') {
      console.debug('[Client Sync] Not authenticated as client, skipping replay')
      return
    }

    const pending = offlineStore.getPendingEntries()
    if (pending.length === 0) {
      console.debug('[Client Sync] No pending entries to replay')
      return
    }

    console.log(`[Client Sync] Starting replay of ${pending.length} entries`)
    offlineStore.setReplayInProgress(true)

    try {
      // Separate message entries from tracking entries
      const messageEntries = pending.filter(e => e.domain === 'message')
      const trackingEntries = pending.filter(e => e.domain !== 'message')

      // Handle message entries separately
      if (messageEntries.length > 0) {
        console.log(`[Client Sync] Processing ${messageEntries.length} message entries`)
        const messagingApi = useMessagingApi()
        
        for (const entry of messageEntries) {
          try {
            offlineStore.updateQueueEntryState(entry.local_id, 'syncing')
            await messagingApi.sendClientMessage({ content: entry.payload.content })
            offlineStore.markEntrySynced(entry.local_id)
            console.log(`[Client Sync] Message sent and synced: ${entry.local_id}`)
          } catch (error) {
            offlineStore.markEntryFailed(entry.local_id, {
              code: 'MESSAGE_SEND_FAILED',
              message: error instanceof Error ? error.message : 'Unknown error',
            })
            console.error(`[Client Sync] Failed to send message ${entry.local_id}:`, error)
          }
        }
      }

      // Handle tracking entries with bulk sync
      if (trackingEntries.length === 0) {
        console.log('[Client Sync] No tracking entries to sync')
        offlineStore.setReplayInProgress(false)
        return
      }

      // Mark all tracking entries as syncing
      trackingEntries.forEach(entry => {
        offlineStore.updateQueueEntryState(entry.local_id, 'syncing')
      })

      // Attempt bulk sync
      const response = await trackingApi.bulkSyncTracking(trackingEntries)

      if (!response) {
        // Network error - mark entries back to queued for later retry
        console.warn('[Client Sync] Bulk sync failed, marking entries for later retry')
        trackingEntries.forEach(entry => {
          offlineStore.updateQueueEntryState(entry.local_id, 'queued')
        })
        return
      }

      // Process sync response
      if (response.synced && response.synced.length > 0) {
        response.synced.forEach(result => {
          offlineStore.markEntrySynced(result.local_id)
          console.log(`[Client Sync] Entry synced: ${result.local_id}`)
        })
      }

      // Handle conflicts (last-write-wins already applied by server)
      if (response.conflicts && response.conflicts.length > 0) {
        response.conflicts.forEach(conflict => {
          offlineStore.markEntryFailed(conflict.local_id, {
            code: 'CONFLICT_RESOLVED',
            message: `Server version preserved: ${conflict.reason}`,
          })
          console.warn(`[Client Sync] Conflict for ${conflict.local_id}: ${conflict.reason}`)
        })
      }

      // Handle errors
      if (response.errors && response.errors.length > 0) {
        response.errors.forEach(error => {
          offlineStore.markEntryFailed(error.local_id, {
            code: error.error_code,
            message: error.error_message,
          })
          console.error(
            `[Client Sync] Error for ${error.local_id}: ${error.error_code}`
          )
        })
      }

      console.log('[Client Sync] Replay completed')
    } catch (err) {
      console.error('[Client Sync] Replay exception:', err)
      // Mark all as queued for later retry
      pending.forEach(entry => {
        offlineStore.updateQueueEntryState(entry.local_id, 'queued')
      })
    } finally {
      offlineStore.setReplayInProgress(false)
      replayScheduled = false
    }
  }

  /**
   * Schedule a replay with debounce to avoid excessive syncing
   */
  function scheduleReplay(delayMs = 1000) {
    if (replayScheduled) return

    replayScheduled = true
    if (reconnectTimeout) clearTimeout(reconnectTimeout)

    reconnectTimeout = setTimeout(() => {
      performQueueReplay()
    }, delayMs)
  }

  /**
   * Setup online/offline listeners
   */
  function setupConnectivityListeners() {
    window.addEventListener('online', () => {
      console.log('[Client Sync] Network reconnected')
      scheduleReplay()
    })

    window.addEventListener('offline', () => {
      console.log('[Client Sync] Network disconnected')
    })
  }

  /**
   * Setup app-open and app-resume listeners
   * Triggers replay on app visibility change (PWA scenario)
   */
  function setupAppLifecycleListeners() {
    // Trigger replay when app returns to foreground
    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'visible') {
        console.log('[Client Sync] App resumed from background')
        scheduleReplay(500) // Shorter delay for app resume
      }
    })

    // On page load/mount, check if there are pending entries
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => {
        console.log('[Client Sync] App initialized')
        scheduleReplay()
      })
    } else {
      console.log('[Client Sync] App initialized (DOM already ready)')
      scheduleReplay()
    }
  }

  /**
   * Cleanup on logout
   * Clear all offline state when user logs out
   */
  function setupLogoutListener() {
    nuxtApp.hook('app:created', () => {
      // Hook into auth logout (if available)
      // This would be called by the auth store when logout occurs
    })
  }

  // Initialize listeners
  if (process.client) {
    console.log('[Client Sync] Initializing offline sync orchestration')
    setupConnectivityListeners()
    setupAppLifecycleListeners()
    setupLogoutListener()
  }

  return {
    provide: {
      clientSync: {
        performQueueReplay,
        scheduleReplay,
      },
    },
  }
})
