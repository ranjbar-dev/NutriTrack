import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
import { toPersianDigits } from '../../app/lib/locale/numerals'
import { formatJalaliDate } from '../../app/lib/locale/jalali'
import { usePersianFormat } from '../../app/composables/usePersianFormat'

describe('persian locale and design token baseline', () => {
  it('defines palette, spacing, and typography tokens in the token file', () => {
    const tokenFile = readFileSync(resolve(process.cwd(), 'app/lib/design/tokens.css'), 'utf8')

    expect(tokenFile).toContain('--color-bg')
    expect(tokenFile).toContain('--space-4')
    expect(tokenFile).toContain('--font-base')
  })

  it('exposes mobile safe-area variables in global styles', () => {
    const globalStyles = readFileSync(resolve(process.cwd(), 'app/assets/css/main.css'), 'utf8')

    expect(globalStyles).toContain('--safe-top')
    expect(globalStyles).toContain('env(safe-area-inset-bottom)')
  })

  it('formats numerals to Persian by default', () => {
    expect(toPersianDigits('1405')).toBe('۱۴۰۵')
  })

  it('renders known Gregorian date in Jalali output', () => {
    expect(formatJalaliDate('2026-04-22')).toContain('۱۴۰۵')
  })

  it('provides composable helpers for mixed locale display', () => {
    const formatter = usePersianFormat()

    expect(formatter.number(12)).toBe('۱۲')
    expect(formatter.otp('123456')).toBe('123456')
  })
})
