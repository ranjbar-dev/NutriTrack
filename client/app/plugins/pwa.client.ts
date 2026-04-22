import { useRegisterSW } from 'virtual:pwa-register/vue'
import { usePlatformPwaStore } from '../stores/platform-pwa'

export default defineNuxtPlugin(() => {
  const pwaStore = usePlatformPwaStore()

  const { needRefresh, offlineReady, updateServiceWorker } = useRegisterSW({
    immediate: true,
    onRegisteredSW() {
      pwaStore.setInstallReady(true)
    },
    onRegisterError() {
      pwaStore.setOffline(false)
    }
  })

  watch(needRefresh, (value) => {
    pwaStore.setNeedRefresh(Boolean(value))
  }, { immediate: true })

  watch(offlineReady, (value) => {
    pwaStore.setOffline(Boolean(value))
  }, { immediate: true })

  return {
    provide: {
      pwaUpdate: {
        refresh: () => updateServiceWorker(true)
      }
    }
  }
})
