import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('nutritionist client roster page', () => {
  it('wires filter and list components to client api', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/pages/nutritionist/clients/index.vue'), 'utf8')

    expect(text).toContain('ClientRosterFilters')
    expect(text).toContain('ClientRosterList')
    expect(text).toContain('useNutritionistClientApi')
    expect(text).toContain('listClients')
  })
})
