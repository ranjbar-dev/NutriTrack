<template>
  <div class="messages-page">
    <AppShell title="پیام‌ها" class="messages-shell">
      <MessageThread
        :messages="displayMessages"
        :is-offline="isOffline"
        class="message-thread-container"
      />

      <MessageComposeBar 
        @send="handleSend"
        class="compose-section"
      />
    </AppShell>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { Message, SendMessageRequest } from '~/types/messaging'
import { useMessagingApi } from '~/composables/useMessagingApi'
import { useClientOfflineStore } from '~/stores/client-offline'
import { usePlatformPwaStore } from '~/stores/platform-pwa'
import AppShell from '~/components/platform/AppShell.vue'
import MessageThread from '~/components/client/MessageThread.vue'
import MessageComposeBar from '~/components/client/MessageComposeBar.vue'

definePageMeta({ layout: 'client' })

const messagingApi = useMessagingApi()
const offlineStore = useClientOfflineStore()
const pwaStore = usePlatformPwaStore()

const messages = ref<Message[]>([])
const isLoading = ref(false)
let pollingInterval: ReturnType<typeof setInterval> | null = null

const isOffline = computed(() => pwaStore.offline)

const displayMessages = computed(() => {
  if (isOffline.value && messages.value.length === 0) {
    return offlineStore.cachedMessages
  }
  return messages.value
})

async function loadMessages(page: number = 1, pageSize: number = 50) {
  if (isLoading.value) return

  isLoading.value = true
  try {
    const { data } = await messagingApi.getClientConversation(page, pageSize)
    if (data.value?.data) {
      // Merge new messages, deduplicate by id
      const newMessages = data.value.data
      const existingIds = new Set(messages.value.map(m => m.id))
      const uniqueNewMessages = newMessages.filter(m => !existingIds.has(m.id))
      messages.value = [...messages.value, ...uniqueNewMessages]

      // Cache messages for offline access
      if (newMessages.length > 0) {
        offlineStore.setCachedMessages(messages.value)
      }
    }
  } catch (error) {
    console.error('Failed to load messages:', error)
  } finally {
    isLoading.value = false
  }
}

function startPolling() {
  // Initial load
  loadMessages()

  // Poll every 10 seconds
  pollingInterval = setInterval(() => {
    loadMessages()
  }, 10_000)
}

function stopPolling() {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

async function handleSend(req: SendMessageRequest) {
  if (isOffline.value) {
    // Offline: text-only, queue it
    if (req.content) {
      const entry = offlineStore.enqueueTrackingWrite({
        domain: 'message',
        payload: { content: req.content },
      })

      if (entry) {
        // Add optimistic message to thread
        const optimisticMessage: Message & { local_id: string; sync_state: string } = {
          id: '',
          sender_id: '',
          receiver_id: '',
          content: req.content,
          is_mine: true,
          read_at: null,
          attachment: null,
          created_at: new Date().toISOString(),
          local_id: entry.local_id,
          sync_state: entry.sync_state,
        }
        messages.value.push(optimisticMessage)
      }
    } else if (req.file) {
      // Reject file upload offline - show error via parent or inline
      console.warn('File uploads require internet connection')
    }
  } else {
    // Online: send directly
    try {
      const sentMessage = await messagingApi.sendClientMessage(req)
      // Add to thread
      messages.value.push(sentMessage)
      offlineStore.setCachedMessages(messages.value)
    } catch (error) {
      console.error('Failed to send message:', error)
      // Show error notice
    }
  }
}

onMounted(() => {
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.messages-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--color-background, #fafafa);
}

.messages-shell {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.message-thread-container {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.compose-section {
  flex-shrink: 0;
}
</style>
