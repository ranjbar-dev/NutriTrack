import { ref, onMounted } from 'vue'

let deferredPrompt: BeforeInstallPromptEvent | null = null
let swRegistration: ServiceWorkerRegistration | null = null

interface BeforeInstallPromptEvent extends Event {
  prompt(): Promise<void>
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>
}

export function usePWA() {
  const canInstall = ref(false)
  const needsUpdate = ref(false)

  function handleBeforeInstallPrompt(e: Event) {
    e.preventDefault()
    deferredPrompt = e as BeforeInstallPromptEvent
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
    navigator.serviceWorker?.ready.then((reg) => {
      swRegistration = reg
      if (reg.waiting) needsUpdate.value = true
      reg.addEventListener('updatefound', () => {
        const newWorker = reg.installing
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
