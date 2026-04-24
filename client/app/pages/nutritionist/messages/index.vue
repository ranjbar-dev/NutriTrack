<template>
  <NutritionistLayout title="پیام‌ها">
    <div class="messages-page">
      <div class="messages-container">
        <EmptyState
          v-if="!isLoading && conversations.length === 0"
          icon="message-circle"
          title="بدون مکالمه"
          description="هنوز هیچ مکالمه‌ای با کلاینت‌ها ندارید"
        />

        <div v-else class="conversations-list">
          <div v-if="isLoading" class="loading-state">
            <div class="loading-skeleton" />
            <div class="loading-skeleton" />
            <div class="loading-skeleton" />
          </div>
          <template v-else>
            <ClientConversationItem
              v-for="conv in conversations"
              :key="conv.clientId"
              :client-id="conv.clientId"
              :client-name="conv.clientName"
              :last-message="conv.lastMessage"
              :message-date="conv.messageDate"
              :unread-count="conv.unreadCount"
              @click="navigateToConversation"
            />
          </template>
        </div>
      </div>
    </div>
  </NutritionistLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessagingApi } from '#composables/useMessagingApi'
import NutritionistLayout from '#components/layouts/NutritionistLayout.vue'
import EmptyState from '#components/shared/EmptyState.vue'
import ClientConversationItem from '#components/nutritionist/ClientConversationItem.vue'

interface ConversationPreview {
  clientId: string
  clientName: string
  lastMessage: string | null
  messageDate: Date | null
  unreadCount: number
}

const router = useRouter()
const { fetchConversationList } = useMessagingApi()

const conversations = ref<ConversationPreview[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)

const loadConversations = async () => {
  isLoading.value = true
  error.value = null
  try {
    const data = await fetchConversationList('nutritionist')
    conversations.value = data.map(conv => ({
      clientId: conv.clientId,
      clientName: conv.clientName,
      lastMessage: conv.lastMessageText,
      messageDate: conv.lastMessageAt ? new Date(conv.lastMessageAt) : null,
      unreadCount: conv.unreadCount
    }))
  } catch (err) {
    error.value = 'خطا در بارگیری مکالمات'
    console.error('Failed to load conversations:', err)
  } finally {
    isLoading.value = false
  }
}

const navigateToConversation = (clientId: string) => {
  router.push(/nutritionist/messages/\)
}

onMounted(() => {
  loadConversations()
})
</script>

<style scoped>
.messages-page {
  max-width: 600px;
  margin: 0 auto;
  height: 100%;
}

.messages-container {
  height: 100%;
}

.conversations-list {
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.loading-state {
  padding: 16px;
}

.loading-skeleton {
  height: 72px;
  background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%);
  background-size: 200% 100%;
  animation: loading 1.5s infinite;
  border-radius: 4px;
  margin-bottom: 12px;
}

@keyframes loading {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
