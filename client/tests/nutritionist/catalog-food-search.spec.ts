import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('catalog food search picker', () => {
  it('uses catalogue api and emits selected food', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/components/nutritionist/FoodSearchPickerSheet.vue'), 'utf8')

    expect(text).toContain('useCatalogueApi')
    expect(text).toContain('searchFoods')
    expect(text).toContain("emit('select', food)")
    expect(text).toContain('getFoodCategories')
  })
})
