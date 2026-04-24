import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('client food request submit flow', () => {
  it('uses food request form card and submit api call', () => {
    const pageText = readFileSync(resolve(process.cwd(), 'app/pages/client/food-requests/index.vue'), 'utf8')
    const formText = readFileSync(resolve(process.cwd(), 'app/components/client/FoodRequestFormCard.vue'), 'utf8')

    expect(pageText).toContain('FoodRequestFormCard')
    expect(pageText).toContain('submitFoodRequest')
    expect(formText).toContain('food_name')
    expect(formText).toContain("emit('submit'")
  })
})
