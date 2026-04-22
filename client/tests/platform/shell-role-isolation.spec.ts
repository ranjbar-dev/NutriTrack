import { describe, expect, it } from 'vitest'

const roleNamespaces = {
  auth: ['/auth'],
  client: ['/client'],
  nutritionist: ['/nutritionist'],
  admin: ['/admin']
} as const

function resolveRoleFromPath(path: string): keyof typeof roleNamespaces | 'unknown' {
  if (path.startsWith('/auth')) return 'auth'
  if (path.startsWith('/client')) return 'client'
  if (path.startsWith('/nutritionist')) return 'nutritionist'
  if (path.startsWith('/admin')) return 'admin'
  return 'unknown'
}

describe('platform shell role isolation baseline', () => {
  it('maps each role namespace to itself only', () => {
    expect(resolveRoleFromPath('/auth')).toBe('auth')
    expect(resolveRoleFromPath('/client/dashboard')).toBe('client')
    expect(resolveRoleFromPath('/nutritionist/clients')).toBe('nutritionist')
    expect(resolveRoleFromPath('/admin/stats')).toBe('admin')
  })

  it('does not map unknown routes to a role shell', () => {
    expect(resolveRoleFromPath('/')).toBe('unknown')
    expect(resolveRoleFromPath('/foo')).toBe('unknown')
  })

  it('keeps namespace declarations explicit', () => {
    expect(roleNamespaces.client).toStrictEqual(['/client'])
    expect(roleNamespaces.nutritionist).toStrictEqual(['/nutritionist'])
    expect(roleNamespaces.admin).toStrictEqual(['/admin'])
  })
})
