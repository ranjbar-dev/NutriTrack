import { onMounted, ref } from 'vue'

let deferredPrompt: BeforeInstallPromptEvent | null = null
let swRegistration: ServiceWorkerRegistration | null = null

interface BeforeInstallPromptEvent extends Event {
  prompt(): Promise<void>
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>
}

export function useClientPWA() {
  const canInstall = ref(false)
  const needsUpdate = ref(false)

  function handleBeforeInstallPrompt(event: Event) {
    event.preventDefault()
    deferredPrompt = event as BeforeInstallPromptEvent
    canInstall.value = true
  }

  async function promptInstall() {
    if (!deferredPrompt) return
    await deferredPrompt.prompt()
    const { outcome } = await deferredPrompt.userChoice
    if (outcome === 'accepted') {
      canInstall.value = false
    }
    deferredPrompt = null
  }

  function applyUpdate() {
    swRegistration?.waiting?.postMessage({ type: 'SKIP_WAITING' })
    window.location.reload()
  }

  onMounted(() => {
    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
    navigator.serviceWorker?.ready.then((registration) => {
      swRegistration = registration
      if (registration.waiting) needsUpdate.value = true
      registration.addEventListener('updatefound', () => {
        const newWorker = registration.installing
        newWorker?.addEventListener('statechange', () => {
          if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
            needsUpdate.value = true
          }
        })
      })
    })
  })

  return { canInstall, needsUpdate, promptInstall, applyUpdate }
}
