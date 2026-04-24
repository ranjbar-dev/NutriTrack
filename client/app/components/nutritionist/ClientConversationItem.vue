<template>
  <div class="client-conversation-item" @click="handleClick">
    <div class="item-container">
      <div class="avatar-section">
        <div class="avatar">{{ initials }}</div>
        <div v-if="unreadCount > 0" class="unread-badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</div>
      </div>
      <div class="content-section" :class="{ unread: unreadCount > 0 }">
        <div class="header-row">
          <h3 class="client-name">{{ clientName }}</h3>
          <span class="date">{{ formattedDate }}</span>
        </div>
        <p class="last-message">{{ lastMessage || 'بدون پیام' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatDistanceToNow } from 'date-fns'
import { faIR } from 'date-fns/locale'

interface Props {
  clientId: string
  clientName: string
  lastMessage: string | null
  messageDate: Date | null
  unreadCount: number
}

const props = withDefaults(defineProps<Props>(), {
  lastMessage: null,
  messageDate: null,
  unreadCount: 0
})

const emit = defineEmits<{
  click: [clientId: string]
}>()

const initials = computed(() => {
  const names = props.clientName.trim().split(/\s+/)
  return names.map(n => n[0]).join('').toUpperCase().slice(0, 2)
})

const formattedDate = computed(() => {
  if (!props.messageDate) return ''
  try {
    return formatDistanceToNow(new Date(props.messageDate), {
      addSuffix: true,
      locale: faIR
    })
  } catch {
    return ''
  }
})

const handleClick = () => {
  emit('click', props.clientId)
}
</script>

<style scoped>
.client-conversation-item {
  padding: 12px 16px;
  border-bottom: 1px solid #e5e7eb;
  cursor: pointer;
  transition: background-color 0.2s;
}
.client-conversation-item:hover {
  background-color: #f9fafb;
}
.item-container {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}
.avatar-section {
  position: relative;
  flex-shrink: 0;
}
.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
}
.unread-badge {
  position: absolute;
  top: -4px;
  left: -4px;
  background-color: #ef4444;
  color: white;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  border: 2px solid white;
}
.content-section {
  flex: 1;
  min-width: 0;
}
.content-section.unread {
  font-weight: 500;
}
.header-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 4px;
}
.client-name {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}
.date {
  font-size: 12px;
  color: #9ca3af;
  flex-shrink: 0;
  margin-right: 8px;
}
.last-message {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.content-section.unread .last-message {
  color: #4b5563;
}
</style>
