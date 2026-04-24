import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('catalog medication search picker', () => {
  it('uses medications endpoint and emits selected medication', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/components/nutritionist/MedicationSearchPickerSheet.vue'), 'utf8')

    expect(text).toContain('searchMedications')
    expect(text).toContain("emit('select', item)")
    expect(text).toContain('جستجو دارو')
  })
})
