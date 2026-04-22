import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
import { resolveAllowedPrefix, shouldRedirectRolePath } from '../../app/middleware/role-shell.global'

function readWorkspaceFile(path: string): string {
  return readFileSync(resolve(process.cwd(), path), 'utf8')
}

describe('platform shell role isolation baseline', () => {
  it('provides dedicated role layouts', () => {
    expect(readWorkspaceFile('app/layouts/auth.vue')).toContain('AppShell')
    expect(readWorkspaceFile('app/layouts/client.vue')).toContain('role="client"')
    expect(readWorkspaceFile('app/layouts/nutritionist.vue')).toContain('role="nutritionist"')
    expect(readWorkspaceFile('app/layouts/admin.vue')).toContain('role="admin"')
  })

  it('binds each role entry page to its dedicated layout', () => {
    expect(readWorkspaceFile('app/pages/auth/index.vue')).toContain("layout: 'auth'")
    expect(readWorkspaceFile('app/pages/client/index.vue')).toContain("layout: 'client'")
    expect(readWorkspaceFile('app/pages/nutritionist/index.vue')).toContain("layout: 'nutritionist'")
    expect(readWorkspaceFile('app/pages/admin/index.vue')).toContain("layout: 'admin'")
  })

  it('resolves allowed namespace prefixes by role', () => {
    expect(resolveAllowedPrefix('client')).toBe('/client')
    expect(resolveAllowedPrefix('nutritionist')).toBe('/nutritionist')
    expect(resolveAllowedPrefix('admin')).toBe('/admin')
    expect(resolveAllowedPrefix('auth')).toBe('/auth')
  })

  it('flags cross-role paths for redirect', () => {
    expect(shouldRedirectRolePath('client', '/nutritionist')).toBe(true)
    expect(shouldRedirectRolePath('nutritionist', '/admin/dashboard')).toBe(true)
    expect(shouldRedirectRolePath('admin', '/admin')).toBe(false)
  })

  it('renders install banner only in client layout intentional flow', () => {
    const clientLayout = readWorkspaceFile('app/layouts/client.vue')
    const authLayout = readWorkspaceFile('app/layouts/auth.vue')
    const nutritionistLayout = readWorkspaceFile('app/layouts/nutritionist.vue')
    const adminLayout = readWorkspaceFile('app/layouts/admin.vue')

    expect(clientLayout).toContain('InstallPromptBanner')
    expect(authLayout).not.toContain('InstallPromptBanner')
    expect(nutritionistLayout).not.toContain('InstallPromptBanner')
    expect(adminLayout).not.toContain('InstallPromptBanner')
  })

  it('allows update and connectivity banners across role layouts', () => {
    const authLayout = readWorkspaceFile('app/layouts/auth.vue')
    const clientLayout = readWorkspaceFile('app/layouts/client.vue')
    const nutritionistLayout = readWorkspaceFile('app/layouts/nutritionist.vue')
    const adminLayout = readWorkspaceFile('app/layouts/admin.vue')

    for (const layout of [authLayout, clientLayout, nutritionistLayout, adminLayout]) {
      expect(layout).toContain('UpdateAvailableBanner')
      expect(layout).toContain('ConnectivityBanner')
    }
  })
})
