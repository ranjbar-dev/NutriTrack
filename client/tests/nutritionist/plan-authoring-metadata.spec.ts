import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('plan authoring metadata', () => {
  it('includes date range and water target fields in period form', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/components/nutritionist/PlanPeriodFormCard.vue'), 'utf8')

    expect(text).toContain('start_date')
    expect(text).toContain('end_date')
    expect(text).toContain('daily_water_target_ml')
    expect(text).toContain("emit('submit'")
  })

  it('uses plan edit page with period form', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/pages/nutritionist/plans/[planId]/edit.vue'), 'utf8')

    expect(text).toContain('PlanPeriodFormCard')
    expect(text).toContain('updatePlan')
    expect(text).toContain('loadPlan')
  })
})
