import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('nutritionist workspace api contracts', () => {
  it('defines client roster endpoints and status mutation', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/composables/useNutritionistClientApi.ts'), 'utf8')

    expect(text).toContain('listClients')
    expect(text).toContain('getClientProfile')
    expect(text).toContain('setClientStatus')
    expect(text).toContain("const baseUrl = '/api/v1'")
  })

  it('defines plan authoring hierarchy endpoints', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/composables/useDietPlanAuthoringApi.ts'), 'utf8')

    expect(text).toContain('createPlan')
    expect(text).toContain('addDay')
    expect(text).toContain('addMeal')
    expect(text).toContain('addOption')
    expect(text).toContain('addItem')
    expect(text).toContain('addExercise')
    expect(text).toContain('addPrescription')
  })
})
