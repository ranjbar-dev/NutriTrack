import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('admin nutritionist detail flow', () => {
  it('wires detail loading, update, status change, and readonly clients', () => {
    const pageText = readFileSync(resolve(process.cwd(), 'app/pages/admin/nutritionists/[id].vue'), 'utf8')
    const detailText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminNutritionistDetailCard.vue'),
      'utf8',
    )
    const formText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminNutritionistEditForm.vue'),
      'utf8',
    )
    const statusText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminNutritionistStatusConfirmSheet.vue'),
      'utf8',
    )
    const clientsText = readFileSync(
      resolve(process.cwd(), 'app/components/admin/AdminNutritionistClientReadonlyList.vue'),
      'utf8',
    )

    expect(pageText).toContain('getNutritionist')
    expect(pageText).toContain('updateNutritionist')
    expect(pageText).toContain('setNutritionistStatus')
    expect(pageText).toContain('listNutritionistClients')
    expect(pageText).toContain('loading.value = false')
    expect(detailText).toContain("emit('change-status')")
    expect(formText).toContain("emit('submit'")
    expect(formText).toContain('const payload: AdminUpdateNutritionistRequest = {}')
    expect(formText).toContain('if (firstName)')
    expect(formText).toContain('if (email)')
    expect(statusText).toContain("emit('confirm')")
    expect(clientsText).toContain('فقط برای مشاهده')
    expect(clientsText).not.toContain('ویرایش')
  })
})