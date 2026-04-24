<template>
  <div v-if="metrics.total > 0" class="sync-strip">
    <div class="sync-content">
      <div class="sync-status">
        <button v-if="metrics.failed > 0" class="retry-btn" @click="retryFailed">
          دوباره تلاش کنید ({{ metrics.failed }})
        </button>
        <div v-else class="status-text">
          <span v-if="metrics.syncing > 0" class="syncing">
            ↻ {{ metrics.syncing }} در حال ارسال
          </span>
          <span v-else-if="metrics.queued > 0" class="queued">
            ⏱️ {{ metrics.queued }} در صف
          </span>
          <span v-else class="synced">
            ✓ {{ metrics.synced }} ارسال شده
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useClientOfflineStore } from '~/app/stores/client-offline'

const store = useClientOfflineStore()

const metrics = computed(() => store.getQueueMetrics())

const retryFailed = () => {
  const failed = store.getPendingEntries().filter(e => e.sync_state === 'failed')
  failed.forEach(entry => {
    store.updateQueueEntryState(entry.local_id, 'queued')
  })
  // Trigger replay
  if (window.$nuxt?.$clientSync) {
    window.$nuxt.$clientSync.performQueueReplay()
  }
}
</script>

<style scoped>
.sync-strip {
  background: linear-gradient(135deg, #f0f9ff 0%, #ecfdf5 100%);
  border-top: 1px solid #d1d5db;
  border-bottom: 1px solid #d1d5db;
  padding: var(--spacing-sm) var(--spacing-lg);
  min-height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sync-content {
  width: 100%;
  direction: rtl;
}

.sync-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
}

.status-text {
  font-size: 14px;
  font-weight: 500;
}

.syncing {
  color: #1e40af;
}

.queued {
  color: #92400e;
}

.synced {
  color: #065f46;
}

.retry-btn {
  background: #dc2626;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
}

.retry-btn:hover {
  background: #b91c1c;
}
</style>
