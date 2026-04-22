import { describe, expect, it } from 'vitest'
import { parseAuthSessionEnvelope } from '../../app/types/auth'
import { createSessionRefreshController } from '../../app/composables/useSessionRefresh'
import { mapAuthError } from '../../app/lib/auth/error-map'

describe('auth session refresh single-flight', () => {
  it('parses login and refresh envelopes into one shared session contract', () => {
    const session = parseAuthSessionEnvelope({
      data: {
        access_token: 'access-token',
        refresh_token: 'refresh-token',
        token_type: 'Bearer',
        user_id: 'user-1',
        role: 'client'
      }
    })

    expect(session.accessToken).toBe('access-token')
    expect(session.refreshToken).toBe('refresh-token')
    expect(session.tokenType).toBe('Bearer')
    expect(session.userId).toBe('user-1')
    expect(session.role).toBe('client')
  })

  it('maps auth errors to controlled Persian messages without exposing backend strings', () => {
    const safe = mapAuthError({
      code: 'INVALID_CREDENTIALS',
      message: 'backend internal detail should never reach UI'
    })

    expect(safe.code).toBe('INVALID_CREDENTIALS')
    expect(safe.message).toBe('ایمیل یا رمز عبور نادرست است. دوباره تلاش کنید.')
    expect(safe.message).not.toContain('backend')
  })

  it('coalesces parallel refresh attempts into one request', async () => {
    let refreshCalls = 0
    const controller = createSessionRefreshController(async () => {
      refreshCalls += 1
      await Promise.resolve()
      return {
        accessToken: 'next-access',
        refreshToken: 'next-refresh',
        tokenType: 'Bearer',
        userId: 'user-1',
        role: 'client'
      }
    })

    const [a, b, c] = await Promise.all([
      controller.refreshSingleFlight(),
      controller.refreshSingleFlight(),
      controller.refreshSingleFlight()
    ])

    expect(refreshCalls).toBe(1)
    expect(a.accessToken).toBe('next-access')
    expect(b.refreshToken).toBe('next-refresh')
    expect(c.role).toBe('client')
  })
})
