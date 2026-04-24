import type { AuthSession } from '../types/auth'
import { useClientOfflineStore } from '~/app/stores/client-offline'

declare const defineNuxtPlugin: <T>(factory: () => T) => T
declare const useAuthSessionStore: () => {
  hydrated: boolean
  applySession: (session: AuthSession) => void
  clearSession: (reason: 'manual' | 'session-expired' | 'token-revoked' | 'unauthorized') => void
  markHydrated: () => void
}
declare const useCookie: <T>(key: string, options?: { default?: () => T }) => { value: T }

function isAuthSession(value: unknown): value is AuthSession {
  const candidate = value as Partial<AuthSession>
  return Boolean(
    candidate
      && typeof candidate.accessToken === 'string'
      && typeof candidate.refreshToken === 'string'
      && candidate.tokenType === 'Bearer'
      && typeof candidate.userId === 'string'
      && typeof candidate.role === 'string'
  )
}

export default defineNuxtPlugin((nuxtApp) => {
  const authStore = useAuthSessionStore()
  const offlineStore = useClientOfflineStore()

  // Watch for logout and clear offline state
  const originalClearSession = authStore.clearSession
  authStore.clearSession = function(reason: 'manual' | 'session-expired' | 'token-revoked' | 'unauthorized') {
    // Clear all offline queue and cache on logout
    offlineStore.clearAllOfflineState()
    console.log('[Auth Bootstrap] Cleared offline state on logout')
    // Call original clearSession
    originalClearSession.call(this, reason)
  }

  if (authStore.hydrated) {
    return
  }

  const sessionCookie = useCookie<string | null>('nt_auth_session', {
    default: () => null
  })

  if (!sessionCookie.value) {
    authStore.markHydrated()
    return
  }

  try {
    const parsed = JSON.parse(sessionCookie.value)
    if (isAuthSession(parsed)) {
      authStore.applySession(parsed)
      return
    }
  } catch {
    // Invalid persisted snapshot must be dropped.
  }

  sessionCookie.value = null
  authStore.clearSession('session-expired')
})
