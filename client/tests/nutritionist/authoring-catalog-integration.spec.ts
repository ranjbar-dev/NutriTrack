import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('authoring catalog integration', () => {
  it('option editor integrates with food picker and macro badge', () => {
    const optionText = readFileSync(resolve(process.cwd(), 'app/components/nutritionist/OptionEditor.vue'), 'utf8')

    expect(optionText).toContain('FoodSearchPickerSheet')
    expect(optionText).toContain('PlanItemMacroBadge')
    expect(optionText).toContain("emit('addItemFromFood'")
  })

  it('exercise prescription editor integrates medication picker and chips', () => {
    const exerciseText = readFileSync(resolve(process.cwd(), 'app/components/nutritionist/ExercisePrescriptionEditor.vue'), 'utf8')

    expect(exerciseText).toContain('MedicationSearchPickerSheet')
    expect(exerciseText).toContain('MedicationChipList')
    expect(exerciseText).toContain("emit('addPrescriptionFromMedication'")
  })
})
