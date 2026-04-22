import type { AuthSession } from '../types/auth'

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

export default defineNuxtPlugin(() => {
  const authStore = useAuthSessionStore()
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
