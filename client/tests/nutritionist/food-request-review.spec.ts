import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('nutritionist food request review flow', () => {
  it('wires list + decision sheet with approve/reject actions', () => {
    const pageText = readFileSync(resolve(process.cwd(), 'app/pages/nutritionist/food-requests/index.vue'), 'utf8')
    const listText = readFileSync(resolve(process.cwd(), 'app/components/nutritionist/FoodRequestReviewList.vue'), 'utf8')
    const sheetText = readFileSync(resolve(process.cwd(), 'app/components/nutritionist/FoodRequestDecisionSheet.vue'), 'utf8')

    expect(pageText).toContain('listPendingFoodRequests')
    expect(pageText).toContain('approveFoodRequest')
    expect(pageText).toContain('rejectFoodRequest')
    expect(listText).toContain('food_name')
    expect(sheetText).toContain("emit('approve'")
    expect(sheetText).toContain("emit('reject'")
  })
})
