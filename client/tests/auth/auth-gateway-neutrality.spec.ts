import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
import { resolveRoleHomePath } from '../../app/stores/auth-session'

describe('auth gateway neutrality', () => {
  it('exposes exactly three role entry options without private route leakage', () => {
    const gatewayText = readFileSync(resolve(process.cwd(), 'app/pages/auth/index.vue'), 'utf8')

    expect(gatewayText).toContain("href: '/auth/client'")
    expect(gatewayText).toContain("href: '/auth/nutritionist'")
    expect(gatewayText).toContain("href: '/auth/admin'")
    expect(gatewayText).not.toContain('/client/today')
    expect(gatewayText).not.toContain('/nutritionist/clients')
    expect(gatewayText).not.toContain('/admin/stats')
  })

  it('redirects authenticated users opening auth gateway to role roots', () => {
    expect(resolveRoleHomePath('client')).toBe('/client')
    expect(resolveRoleHomePath('nutritionist')).toBe('/nutritionist')
    expect(resolveRoleHomePath('super_admin')).toBe('/admin')
  })
})
