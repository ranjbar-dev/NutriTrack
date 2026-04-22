import { defineStore } from 'pinia'
import type { AuthRole, AuthSession } from '../types/auth'

declare const useCookie: undefined | (<T>(key: string, options?: { default?: () => T }) => { value: T })

export type LogoutReason = 'manual' | 'session-expired' | 'token-revoked' | 'unauthorized'

export interface AuthSessionState {
  accessToken: string | null
  refreshToken: string | null
  role: AuthRole | null
  userId: string | null
  hydrated: boolean
  lastLogoutReason: LogoutReason | null
}

export interface RoleScopedCleaners {
  clearClient?: () => void
  clearNutritionist?: () => void
  clearAdmin?: () => void
}

export function createAuthSessionState(): AuthSessionState {
  return {
    accessToken: null,
    refreshToken: null,
    role: null,
    userId: null,
    hydrated: false,
    lastLogoutReason: null
  }
}

export function resolveRoleHomePath(role: AuthRole): string {
  if (role === 'client') {
    return '/client'
  }

  if (role === 'nutritionist') {
    return '/nutritionist'
  }

  return '/admin'
}

export function resolveRoleAuthPath(role: AuthRole | null): string {
  if (role === 'client') {
    return '/auth/client'
  }

  if (role === 'nutritionist') {
    return '/auth/nutritionist'
  }

  if (role === 'super_admin') {
    return '/auth/admin'
  }

  return '/auth'
}

export function runRoleScopedCleanup(role: AuthRole | null, cleaners: RoleScopedCleaners): void {
  if (role === 'client') {
    cleaners.clearClient?.()
    return
  }

  if (role === 'nutritionist') {
    cleaners.clearNutritionist?.()
    return
  }

  if (role === 'super_admin') {
    cleaners.clearAdmin?.()
  }
}

function setRoleCookie(role: AuthRole | null): void {
  if (typeof useCookie !== 'function') {
    return
  }

  const roleCookie = useCookie<AuthRole | null>('nt_role', {
    default: () => null
  })
  roleCookie.value = role
}

function setSessionCookie(session: AuthSession | null): void {
  if (typeof useCookie !== 'function') {
    return
  }

  const sessionCookie = useCookie<string | null>('nt_auth_session', {
    default: () => null
  })
  sessionCookie.value = session ? JSON.stringify(session) : null
}

export const useAuthSessionStore = defineStore('auth-session', {
  state: (): AuthSessionState => createAuthSessionState(),
  actions: {
    applySession(session: AuthSession): void {
      this.accessToken = session.accessToken
      this.refreshToken = session.refreshToken
      this.role = session.role
      this.userId = session.userId
      this.hydrated = true
      this.lastLogoutReason = null
      setRoleCookie(session.role)
      setSessionCookie(session)
    },
    markHydrated(): void {
      this.hydrated = true
    },
    clearSession(reason: LogoutReason): void {
      this.accessToken = null
      this.refreshToken = null
      this.role = null
      this.userId = null
      this.lastLogoutReason = reason
      this.hydrated = true
      setRoleCookie(null)
      setSessionCookie(null)
    },
    async logoutWithCleanup(
      logoutAction: (refreshToken: string) => Promise<void>,
      cleaners: RoleScopedCleaners,
      reason: LogoutReason
    ): Promise<void> {
      const refreshToken = this.refreshToken
      const role = this.role

      if (refreshToken) {
        try {
          await logoutAction(refreshToken)
        } catch {
          // Backend logout failure should not block local secure cleanup.
        }
      }

      runRoleScopedCleanup(role, cleaners)
      this.clearSession(reason)
    }
  }
})
