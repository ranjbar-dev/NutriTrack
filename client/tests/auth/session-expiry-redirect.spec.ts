import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
import { resolveRoleAuthPath } from '../../app/stores/auth-session'

function read(path: string): string {
  return readFileSync(resolve(process.cwd(), path), 'utf8')
}

describe('session expiry redirect behavior', () => {
  it('routes forced logout to role-appropriate auth entry', () => {
    expect(resolveRoleAuthPath('client')).toBe('/auth/client')
    expect(resolveRoleAuthPath('nutritionist')).toBe('/auth/nutritionist')
    expect(resolveRoleAuthPath('super_admin')).toBe('/auth/admin')
  })

  it('shows session-expired notice in auth layout handoff', () => {
    const authLayout = read('app/layouts/auth.vue')

    expect(authLayout).toContain('SessionExpiredNotice')
    expect(authLayout).toContain("String(route.query.reason ?? '') === 'session-expired'")
  })

  it('keeps bootstrap plugin contract for restoring persisted auth session', () => {
    const pluginText = read('app/plugins/auth-bootstrap.client.ts')

    expect(pluginText).toContain("useCookie<string | null>('nt_auth_session'")
    expect(pluginText).toContain('authStore.applySession(parsed)')
    expect(pluginText).toContain("authStore.clearSession('session-expired')")
  })
})
