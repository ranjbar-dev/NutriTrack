import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

function read(path: string): string {
  return readFileSync(resolve(process.cwd(), path), 'utf8')
}

describe('role credential login contracts', () => {
  it('submits nutritionist and admin forms through /auth/login and lands in role roots', () => {
    const nutritionistPage = read('app/pages/auth/nutritionist/index.vue')
    const adminPage = read('app/pages/auth/admin/index.vue')

    expect(nutritionistPage).toContain('await authApi.login')
    expect(adminPage).toContain('await authApi.login')
    expect(nutritionistPage).toContain("await navigateTo('/nutritionist')")
    expect(adminPage).toContain("await navigateTo('/admin')")
  })

  it('keeps controlled persian error messages for credential failures', () => {
    const nutritionistPage = read('app/pages/auth/nutritionist/index.vue')
    const adminPage = read('app/pages/auth/admin/index.vue')

    expect(nutritionistPage).toContain('mapAuthError')
    expect(adminPage).toContain('mapAuthError')
    expect(nutritionistPage).toContain('errorMessage.value = safeError.message')
    expect(adminPage).toContain('errorMessage.value = safeError.message')
  })
})
