import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('persian locale and design token baseline', () => {
  it('defines palette, spacing, and typography tokens in the token file', () => {
    const tokenFile = readFileSync(resolve(process.cwd(), 'app/lib/design/tokens.css'), 'utf8')

    expect(tokenFile).toContain('--color-bg')
    expect(tokenFile).toContain('--space-4')
    expect(tokenFile).toContain('--font-base')
  })

  it('exposes mobile safe-area variables in global styles', () => {
    const globalStyles = readFileSync(resolve(process.cwd(), 'app/assets/css/main.css'), 'utf8')

    expect(globalStyles).toContain('--safe-top')
    expect(globalStyles).toContain('env(safe-area-inset-bottom)')
  })
})
