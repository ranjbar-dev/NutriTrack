import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('plan authoring structure', () => {
  it('new plan page wires day, meal, option and prescription editors', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/pages/nutritionist/clients/[id]/plans/new.vue'), 'utf8')

    expect(text).toContain('PlanDayEditor')
    expect(text).toContain('MealEditor')
    expect(text).toContain('OptionEditor')
    expect(text).toContain('ExercisePrescriptionEditor')
    expect(text).toContain('addDay')
  })
})
