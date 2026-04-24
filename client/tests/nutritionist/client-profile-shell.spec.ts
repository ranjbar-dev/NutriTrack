import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('nutritionist client profile shell', () => {
  it('renders tabs and section links for plans, messages, and labs', () => {
    const text = readFileSync(resolve(process.cwd(), 'app/pages/nutritionist/clients/[id]/index.vue'), 'utf8')

    expect(text).toContain('ClientProfileTabs')
    expect(text).toContain('/plans/new')
    expect(text).toContain('/messages')
    expect(text).toContain('/labs')
  })
})
