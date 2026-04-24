<template>
  <div class="sync-chip" :class="`state-${resolvedState}`">
    <span class="icon">{{ stateIcon }}</span>
    <span v-if="resolvedCount > 0" class="count">{{ resolvedCount }}</span>
    <span v-if="label" class="label">{{ label }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useClientOfflineStore } from '~/app/stores/client-offline'

type ChipState = 'queued' | 'syncing' | 'synced' | 'failed' | 'empty'

const props = defineProps<{
  state?: ChipState
  count?: number
  label?: string
}>()

const store = useClientOfflineStore()

const metrics = computed(() => store.getQueueMetrics())

const mainState = computed<ChipState>(() => {
  if (metrics.value.failed > 0) return 'failed'
  if (metrics.value.syncing > 0) return 'syncing'
  if (metrics.value.queued > 0) return 'queued'
  if (metrics.value.synced > 0) return 'synced'
  return 'empty'
})

const resolvedState = computed<ChipState>(() => props.state ?? mainState.value)

const resolvedCount = computed(() => {
  if (typeof props.count === 'number') {
    return props.count
  }

  return metrics.value.total
})

const stateIcon = computed(() => {
  switch (resolvedState.value) {
    case 'failed':
      return '✗'
    case 'syncing':
      return '↻'
    case 'queued':
      return '⏱️'
    case 'synced':
      return '✓'
    default:
      return '—'
  }
})
</script>

<style scoped>
.sync-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: 4px 10px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 600;
}

.state-queued {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fcd34d;
}

.state-syncing {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #93c5fd;
}

.state-synced {
  background: #d1fae5;
  color: #065f46;
  border: 1px solid #6ee7b7;
}

.state-failed {
  background: #fee2e2;
  color: #991b1b;
  border: 1px solid #fca5a5;
}

.state-empty {
  background: #f3f4f6;
  color: #6b7280;
  border: 1px solid #d1d5db;
}

.icon {
  display: inline-block;
}

.count {
  display: inline-block;
  background: rgba(0, 0, 0, 0.1);
  padding: 0 4px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.label {
  font-size: 11px;
}
</style>
