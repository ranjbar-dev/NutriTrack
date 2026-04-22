import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
import {
  createBannerStateView,
  createPlatformPwaState,
  createRefreshAction,
  shouldShowInstallPromptAtIntentionalMoment
} from '../../app/stores/platform-pwa'

describe('pwa update prompt contracts', () => {
  it('exposes typed install and refresh flags from store helper', () => {
    const state = createPlatformPwaState()

    expect(state.needRefresh).toBe(false)
    expect(state.showInstallPrompt).toBe(false)
    expect(state.offline).toBe(false)
  })

  it('shows install prompt only after an intentional moment', () => {
    expect(shouldShowInstallPromptAtIntentionalMoment('first-paint')).toBe(false)
    expect(shouldShowInstallPromptAtIntentionalMoment('role-shell-ready')).toBe(true)
  })

  it('derives update and install banner visibility from store flags', () => {
    const view = createBannerStateView({
      ...createPlatformPwaState(),
      needRefresh: true,
      installReady: true,
      showInstallPrompt: true
    })

    expect(view.showUpdateBanner).toBe(true)
    expect(view.showInstallBanner).toBe(true)
    expect(view.showConnectivityBanner).toBe(false)
  })

  it('wires explicit refresh action callback', () => {
    let refreshed = false
    const action = createRefreshAction(() => {
      refreshed = true
    })

    action.refresh()
    expect(refreshed).toBe(true)
  })

  it('contains conservative runtime cache boundaries in nuxt config', () => {
    const configPath = resolve(process.cwd(), 'nuxt.config.ts')
    const nuxtConfigText = readFileSync(configPath, 'utf8')

    expect(nuxtConfigText).toContain('@vite-pwa/nuxt')
    expect(nuxtConfigText).toContain('api/**')
    expect(nuxtConfigText).toContain('NetworkOnly')
  })
})
