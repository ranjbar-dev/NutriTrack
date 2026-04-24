import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('admin dashboard and roster flow', () => {
  it('renders only api-backed stats metrics on the admin dashboard', () => {
    const pageText = readFileSync(resolve(process.cwd(), 'app/pages/admin/index.vue'), 'utf8')
    const gridText = readFileSync(resolve(process.cwd(), 'app/components/admin/AdminStatsKpiGrid.vue'), 'utf8')

    expect(pageText).toContain('useAdminStatsApi')
    expect(gridText).toContain('total_nutritionists')
    expect(gridText).toContain('active_nutritionists')
    expect(gridText).toContain('inactive_nutritionists')
    expect(gridText).toContain('total_clients')
    expect(gridText).toContain('total_foods')
    expect(gridText).toContain('active_diet_plans')
    expect(gridText).not.toContain('active_clients')
    expect(gridText).not.toContain('inactive_clients')
    expect(pageText).not.toContain('audit')
  })

  it('wires nutritionist roster search, create sheet, and detail navigation', () => {
    const pageText = readFileSync(resolve(process.cwd(), 'app/pages/admin/nutritionists/index.vue'), 'utf8')
    const filtersText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminNutritionistRosterFilters.vue'),
      'utf8',
    )
    const listText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminNutritionistRosterList.vue'),
      'utf8',
    )
    const sheetText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminNutritionistCreateSheet.vue'),
      'utf8',
    )

    expect(pageText).toContain('useAdminNutritionistApi')
    expect(pageText).toContain('listNutritionists')
    expect(pageText).toContain('createNutritionist')
    expect(pageText).toContain('/admin/nutritionists/${id}')
    expect(filtersText).toContain("emit('apply'")
    expect(listText).toContain("emit('open'")
    expect(sheetText).toContain("emit('submit'")
  })
})