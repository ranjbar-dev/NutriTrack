<template>
  <div class="message-thread">
    <div v-if="isOffline" class="offline-notice-container">
      <InlineNotice type="info">
        اتصال اینترنت ندارید — پیام‌های اخیر نمایش داده می‌شوند
      </InlineNotice>
    </div>

    <div v-if="messages.length === 0 && !isOffline" class="empty-container">
      <EmptyState title="هنوز پیامی وجود ندارد" />
    </div>

    <div 
      ref="threadContainer"
      class="message-list"
      :class="{ offline: isOffline }"
    >
      <MessageBubble 
        v-for="msg in reversed"
        :key="msg.id || (msg as any).local_id"
        :message="msg as Message & { local_id?: string; sync_state?: string }"
        class="message-item"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import type { Message } from '~/types/messaging'
import MessageBubble from './MessageBubble.vue'
import EmptyState from '~/components/platform/EmptyState.vue'
import InlineNotice from '~/components/platform/InlineNotice.vue'

const props = defineProps<{
  messages: (Message & { local_id?: string; sync_state?: string })[]
  isOffline: boolean
}>()

const threadContainer = ref<HTMLElement | null>(null)

const reversed = computed(() => {
  return [...props.messages].reverse()
})

watch(
  () => reversed.value.length,
  async () => {
    await nextTick()
    scrollToBottom()
  }
)

function scrollToBottom() {
  if (threadContainer.value) {
    threadContainer.value.scrollTop = threadContainer.value.scrollHeight
  }
}

// Initial scroll
watch(
  () => props.messages,
  () => {
    nextTick(() => scrollToBottom())
  },
  { immediate: true }
)
</script>

<style scoped>
.message-thread {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.offline-notice-container {
  padding: 12px;
  background-color: rgba(255, 152, 0, 0.05);
}

.empty-container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
}

.message-list {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  gap: 8px;
}

.message-list.offline {
  opacity: 0.9;
}

.message-item {
  animation: slideIn 0.2s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
