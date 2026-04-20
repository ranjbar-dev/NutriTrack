import { describe, expect, it } from 'vitest'
import { sumWaterAmounts } from '~/utils/tracking'

describe('water log total helper', () => {
  it('sums amount_ml across multiple entries for the same day', () => {
    expect(sumWaterAmounts([
      { amount_ml: 200 },
      { amount_ml: 350 },
      { amount_ml: 500 },
    ])).toBe(1050)
  })

  it('returns 0 when no entries logged today', () => {
    expect(sumWaterAmounts([])).toBe(0)
  })
})
