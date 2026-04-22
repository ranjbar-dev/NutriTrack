import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

function normalizeMobile(raw: string): string {
  return raw.replace(/[^0-9]/g, '')
}

function isIranMobile(raw: string): boolean {
  return /^09\d{9}$/.test(normalizeMobile(raw))
}

describe('client otp validation request flow', () => {
  it('accepts latin digits and enforces iranian mobile format', () => {
    expect(isIranMobile('09123456789')).toBe(true)
    expect(isIranMobile('09-123-456-789')).toBe(true)
    expect(isIranMobile('989123456789')).toBe(false)
    expect(isIranMobile('0912345')).toBe(false)
  })

  it('disables send action during pending state and cooldown windows', () => {
    const pageText = readFileSync(resolve(process.cwd(), 'app/pages/auth/client/index.vue'), 'utf8')

    expect(pageText).toContain('const submitDisabled = computed(() => !validIranMobile.value || isPending.value || cooldownSeconds.value > 0)')
    expect(pageText).toContain('startCooldown(60)')
    expect(pageText).toContain('await authApi.sendOtp')
  })
})
