import { describe, expect, it } from 'vitest'
import { evaluateAuthAccess } from '../../app/middleware/auth-access.global'

describe('route access control', () => {
  it('denies role-mismatch namespace access and redirects safely', () => {
    expect(evaluateAuthAccess('/nutritionist', 'client')).toEqual({
      allowed: false,
      redirectTo: '/client'
    })

    expect(evaluateAuthAccess('/admin/settings', 'nutritionist')).toEqual({
      allowed: false,
      redirectTo: '/nutritionist'
    })
  })

  it('guards protected namespaces for unauthenticated access', () => {
    expect(evaluateAuthAccess('/client', null)).toEqual({
      allowed: false,
      redirectTo: '/auth'
    })

    expect(evaluateAuthAccess('/auth', null)).toEqual({
      allowed: true,
      redirectTo: null
    })
  })

  it('redirects authenticated users away from auth routes to role roots', () => {
    expect(evaluateAuthAccess('/auth/client', 'client')).toEqual({
      allowed: false,
      redirectTo: '/client'
    })

    expect(evaluateAuthAccess('/auth', 'super_admin')).toEqual({
      allowed: false,
      redirectTo: '/admin'
    })
  })
})
