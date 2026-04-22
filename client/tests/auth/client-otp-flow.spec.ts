import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('client otp verify flow contracts', () => {
  it('supports 6-digit input paste and per-cell keyboard behavior', () => {
    const otpInputText = readFileSync(resolve(process.cwd(), 'app/components/auth/OtpInput.vue'), 'utf8')

    expect(otpInputText).toContain("grid-template-columns: repeat(6")
    expect(otpInputText).toContain('updateByCell(index')
    expect(otpInputText).toContain('removeByCell(index)')
    expect(otpInputText).toContain('handlePaste(event')
    expect(otpInputText).toContain("pattern=\"[0-9]*\"")
  })

  it('locks verify after repeated invalid attempts until resend resets state', () => {
    const verifyPageText = readFileSync(resolve(process.cwd(), 'app/pages/auth/client/verify.vue'), 'utf8')

    expect(verifyPageText).toContain('const lockedByAttempts = computed(() => failedAttempts.value >= 3)')
    expect(verifyPageText).toContain("if (safeError.code === 'OTP_INVALID')")
    expect(verifyPageText).toContain('failedAttempts.value += 1')
    expect(verifyPageText).toContain('failedAttempts.value = 0')
    expect(verifyPageText).toContain('await authApi.sendOtp')
  })
})
