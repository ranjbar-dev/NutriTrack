import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

describe('pwa cache boundary regression guards', () => {
  it('keeps authenticated api traffic out of runtime cache allowlists', () => {
    const configText = readFileSync(resolve(process.cwd(), 'nuxt.config.ts'), 'utf8')

    expect(configText).toContain('cache-boundary::api-network-only')
    expect(configText).toContain('NetworkOnly')
    expect(configText).not.toContain('api-cache-first')
  })

  it('keeps static shell assets on explicit allowlist', () => {
    const configText = readFileSync(resolve(process.cwd(), 'nuxt.config.ts'), 'utf8')

    expect(configText).toContain('globPatterns')
    expect(configText).toContain('js,css,html,ico,png,svg,woff2')
  })
})
