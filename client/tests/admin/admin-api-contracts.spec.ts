import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('admin api contracts', () => {
  it('defines admin stats and nutritionist contract fields without unsupported client split metrics', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/types/admin.ts'), 'utf8')

    expect(text).toContain('export interface AdminStats')
    expect(text).toContain('total_clients: number')
    expect(text).not.toContain('active_clients')
    expect(text).not.toContain('inactive_clients')
    expect(text).toContain('export interface AdminNutritionist')
    expect(text).toContain('export interface AdminNutritionistStatusRequest')
  })

  it('defines composables for admin stats and nutritionist management endpoints', () => {
    const statsText = readFileSync(resolve(process.cwd(), 'app/composables/useAdminStatsApi.ts'), 'utf8')
    const nutritionistText = readFileSync(
      resolve(process.cwd(), 'app/composables/useAdminNutritionistApi.ts'),
      'utf8',
    )

    expect(statsText).toContain("const baseUrl = '/api/v1'")
    expect(statsText).toContain('/admin/stats')
    expect(nutritionistText).toContain('listNutritionists')
    expect(nutritionistText).toContain('createNutritionist')
    expect(nutritionistText).toContain('updateNutritionist')
    expect(nutritionistText).toContain('setNutritionistStatus')
    expect(nutritionistText).toContain('listNutritionistClients')
    expect(nutritionistText).toContain('/admin/nutritionists')
    expect(nutritionistText).toContain('/status')
  })

  it('defines elevated admin catalogue endpoints and omits audit-log methods', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/composables/useAdminCatalogueApi.ts'), 'utf8')
    const typeText = readFileSync(resolve(process.cwd(), 'app/types/catalogue.ts'), 'utf8')

    expect(typeText).toContain('export interface AdminFoodSearchQuery')
    expect(typeText).toContain('export interface AdminMedicationSearchQuery')
    expect(typeText).toContain('export interface AdminCreateFoodCategoryRequest')
    expect(text).toContain('searchAdminFoods')
    expect(text).toContain('forceDeleteFood')
    expect(text).toContain('searchAdminMedications')
    expect(text).toContain('forceDeleteMedication')
    expect(text).toContain('createFoodCategory')
    expect(text).toContain('deleteFoodCategory')
    expect(text).toContain('/admin/foods')
    expect(text).toContain('/admin/medications')
    expect(text).toContain('/admin/food-categories')
    expect(text).not.toContain('audit')
  })
})