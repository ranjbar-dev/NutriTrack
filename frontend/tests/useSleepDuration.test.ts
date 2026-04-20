import { describe, expect, it } from 'vitest'
import { computeSleepDurationMinutes, formatSleepDuration } from '~/utils/tracking'

describe('sleep duration helpers', () => {
  it('computes duration in minutes across midnight', () => {
    expect(computeSleepDurationMinutes('23:30', '06:15')).toBe(405)
    expect(formatSleepDuration(405)).toBe('۶ ساعت و ۴۵ دقیقه')
  })

  it('computes duration for same-night sleep', () => {
    expect(computeSleepDurationMinutes('22:00', '23:30')).toBe(90)
    expect(formatSleepDuration(90)).toBe('۱ ساعت و ۳۰ دقیقه')
  })

  it('returns 0 for invalid or missing times', () => {
    expect(computeSleepDurationMinutes('', '08:00')).toBe(0)
    expect(computeSleepDurationMinutes('aa', '08:00')).toBe(0)
    expect(formatSleepDuration(0)).toBe('۰ دقیقه')
  })
})
