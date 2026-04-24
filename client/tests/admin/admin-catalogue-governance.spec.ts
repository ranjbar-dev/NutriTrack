import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('admin catalogue governance flow', () => {
  it('uses elevated admin food and medication endpoints with explicit delete confirmation', () => {
    const foodsPage = readFileSync(resolve(process.cwd(), 'app/pages/admin/catalogue/foods.vue'), 'utf8')
    const medicationsPage = readFileSync(
      resolve(process.cwd(), 'app/pages/admin/catalogue/medications.vue'),
      'utf8',
    )
    const confirmSheet = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminDangerConfirmSheet.vue'),
      'utf8',
    )

    expect(foodsPage).toContain('searchAdminFoods')
    expect(foodsPage).toContain('forceDeleteFood')
    expect(foodsPage).toContain('activeQuery')
    expect(foodsPage).toContain('refreshFoods(activeQuery.value)')
    expect(medicationsPage).toContain('searchAdminMedications')
    expect(medicationsPage).toContain('forceDeleteMedication')
    expect(medicationsPage).toContain('activeQuery')
    expect(medicationsPage).toContain('refreshMedications(activeQuery.value)')
    expect(confirmSheet).toContain("emit('confirm')")
  })

  it('wires category governance and omits audit log ui', () => {
    const categoriesPage = readFileSync(
      resolve(process.cwd(), 'app/pages/admin/catalogue/categories.vue'),
      'utf8',
    )
    const managerText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminFoodCategoryManager.vue'),
      'utf8',
    )

    expect(categoriesPage).toContain('createFoodCategory')
    expect(categoriesPage).toContain('deleteFoodCategory')
    expect(managerText).toContain("emit('create'")
    expect(managerText).toContain("emit('delete'")
    expect(categoriesPage).not.toContain('audit')
  })
})