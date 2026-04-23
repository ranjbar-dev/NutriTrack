<template>
  <div
    :class="[
      'message-bubble',
      message.is_mine ? 'bubble-mine' : 'bubble-other'
    ]"
  >
    <div class="bubble-content">
      <p v-if="message.content" class="bubble-text">{{ message.content }}</p>
      
      <div v-if="message.attachment" class="attachment-section">
        <span class="attachment-name">📎 {{ message.attachment.name }}</span>
        <a 
          v-if="message.attachment.url"
          :href="message.attachment.url"
          class="attachment-link"
          target="_blank"
          rel="noopener noreferrer"
        >
          دانلود
        </a>
      </div>

      <SyncStateChip 
        v-if="!message.id && 'local_id' in message"
        :state="message.sync_state || 'queued'"
        class="sync-chip"
      />

      <div v-if="message.is_mine && message.read_at" class="read-marker">
        ✓ خوانده شد
      </div>
    </div>
    <span class="bubble-time">{{ formatTime(message.created_at) }}</span>
  </div>
</template>

<script setup lang="ts">
import type { Message } from '~/types/messaging'
import { usePersianFormat } from '~/composables/usePersianFormat'

defineProps<{
  message: Message & { local_id?: string; sync_state?: string }
}>()

const { formatTime } = usePersianFormat()
</script>

<style scoped>
.message-bubble {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 16px;
  margin-bottom: 12px;
  border-radius: 12px;
  max-width: 80%;
  word-wrap: break-word;
}

.bubble-mine {
  align-self: flex-end;
  background-color: var(--color-primary-light, #e3f2fd);
  border-bottom-right-radius: 0;
  margin-right: 0;
  text-align: right;
}

.bubble-other {
  align-self: flex-start;
  background-color: var(--color-surface-light, #f5f5f5);
  border-bottom-left-radius: 0;
  margin-left: 0;
  text-align: left;
}

.bubble-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bubble-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.4;
  word-break: break-word;
}

.attachment-section {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}

.attachment-name {
  font-size: 12px;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.attachment-link {
  font-size: 12px;
  color: var(--color-primary, #1976d2);
  text-decoration: none;
  font-weight: 500;
  white-space: nowrap;
}

.sync-chip {
  align-self: flex-end;
}

.read-marker {
  font-size: 11px;
  color: var(--color-text-secondary, #757575);
  opacity: 0.7;
  text-align: left;
}

.bubble-time {
  font-size: 11px;
  color: var(--color-text-secondary, #757575);
  opacity: 0.6;
  text-align: left;
}

.bubble-mine .bubble-time {
  text-align: right;
}
</style>
