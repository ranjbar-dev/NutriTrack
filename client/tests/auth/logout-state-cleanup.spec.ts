import { setActivePinia, createPinia } from 'pinia'
import { describe, expect, it } from 'vitest'
import {
  resolveRoleAuthPath,
  resolveRoleHomePath,
  runRoleScopedCleanup,
  useAuthSessionStore
} from '../../app/stores/auth-session'

describe('logout state cleanup', () => {
  it('clears local tokens and role state even if logout request fails', async () => {
    setActivePinia(createPinia())
    const store = useAuthSessionStore()

    store.applySession({
      accessToken: 'access',
      refreshToken: 'refresh',
      tokenType: 'Bearer',
      userId: 'u-1',
      role: 'client'
    })

    const cleanupCalls: string[] = []

    await store.logoutWithCleanup(
      async () => {
        throw new Error('network failed')
      },
      {
        clearClient: () => {
          cleanupCalls.push('client')
        }
      },
      'session-expired'
    )

    expect(cleanupCalls).toEqual(['client'])
    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
    expect(store.role).toBeNull()
    expect(store.lastLogoutReason).toBe('session-expired')
  })

  it('resolves deterministic role-safe home and auth routes', () => {
    expect(resolveRoleHomePath('client')).toBe('/client')
    expect(resolveRoleHomePath('nutritionist')).toBe('/nutritionist')
    expect(resolveRoleHomePath('super_admin')).toBe('/admin')

    expect(resolveRoleAuthPath('client')).toBe('/auth/client')
    expect(resolveRoleAuthPath('nutritionist')).toBe('/auth/nutritionist')
    expect(resolveRoleAuthPath('super_admin')).toBe('/auth/admin')
    expect(resolveRoleAuthPath(null)).toBe('/auth')
  })

  it('invokes only the matching role cleanup branch', () => {
    const calls: string[] = []

    runRoleScopedCleanup('nutritionist', {
      clearClient: () => {
        calls.push('client')
      },
      clearNutritionist: () => {
        calls.push('nutritionist')
      },
      clearAdmin: () => {
        calls.push('admin')
      }
    })

    expect(calls).toEqual(['nutritionist'])
  })
})
