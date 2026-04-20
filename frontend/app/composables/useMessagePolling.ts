/**
 * useMessagePolling — 10-second polling interval for new chat messages.
 * Calls pollNewMessages and tracks the last-seen timestamp.
 */
export function useMessagePolling(partnerId: string) {
  const messageStore = useMessageStore()
  let intervalId: ReturnType<typeof setInterval> | null = null
  let lastSeen = new Date().toISOString()

  function start() {
    if (intervalId) return
    intervalId = setInterval(async () => {
      const newMsgs = await messageStore.pollNewMessages(partnerId, lastSeen)
      if (newMsgs.length > 0) {
        lastSeen = newMsgs[newMsgs.length - 1].sent_at
      }
    }, 10_000)
  }

  function stop() {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
  }

  onMounted(start)
  onUnmounted(stop)

  return { start, stop }
}
