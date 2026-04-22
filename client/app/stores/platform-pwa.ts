import { defineStore } from 'pinia'

export type InstallPromptMoment = 'first-paint' | 'role-shell-ready' | 'user-opened-menu'

export interface PlatformPwaState {
  needRefresh: boolean
  showInstallPrompt: boolean
  offline: boolean
  installReady: boolean
}

export function createPlatformPwaState(): PlatformPwaState {
  return {
    needRefresh: false,
    showInstallPrompt: false,
    offline: false,
    installReady: false
  }
}

export function shouldShowInstallPromptAtIntentionalMoment(moment: InstallPromptMoment): boolean {
  return moment === 'role-shell-ready' || moment === 'user-opened-menu'
}

export const usePlatformPwaStore = defineStore('platform-pwa', {
  state: (): PlatformPwaState => createPlatformPwaState(),
  actions: {
    setNeedRefresh(value: boolean): void {
      this.needRefresh = value
    },
    setInstallReady(value: boolean): void {
      this.installReady = value
    },
    setOffline(value: boolean): void {
      this.offline = value
    },
    openInstallPrompt(moment: InstallPromptMoment): void {
      this.showInstallPrompt = this.installReady && shouldShowInstallPromptAtIntentionalMoment(moment)
    },
    closeInstallPrompt(): void {
      this.showInstallPrompt = false
    }
  }
})
