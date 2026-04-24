<script setup lang="ts">
import { computed } from 'vue'
import { useClientOfflineStore } from '~/app/stores/client-offline'
import { formatTrackingTimestampFa, trackingDomainLabel } from '~/app/lib/tracking/history'
import { retryAllFailedEntries, retryFailedEntry } from '~/app/lib/tracking/retry'
import FailedSyncList from '~/app/components/client/FailedSyncList.vue'
import TrackingProgressSummary from '~/app/components/client/TrackingProgressSummary.vue'
import SyncStateChip from '~/app/components/client/SyncStateChip.vue'

definePageMeta({ layout: 'client' })

const store = useClientOfflineStore()

const rows = computed(() => {
  return store
    .getAllQueueEntries()
    .slice()
    .sort((left, right) => {
      return new Date(right.created_at).getTime() - new Date(left.created_at).getTime()
    })
})

const failedRows = computed(() => rows.value.filter((entry) => entry.sync_state === 'failed'))

function retryOne(localId: string) {
  retryFailedEntry(store, localId)
}

function retryAll() {
  retryAllFailedEntries(store)
}
</script>

<template>
  <div class="tracking-history-page">
    <header>
      <h1>تاریخچه ثبت ها</h1>
      <p>نمایش ثبت های اخیر همراه با وضعیت همگام سازی</p>
    </header>

    <TrackingProgressSummary :entries="rows" />

    <section class="rows-section">
      <h2>ثبت های اخیر</h2>
      <p v-if="rows.length === 0" class="empty">هنوز ثبتی وجود ندارد.</p>

      <ul v-else>
        <li v-for="entry in rows" :key="entry.local_id">
          <div>
            <p class="domain">{{ trackingDomainLabel(entry.domain) }}</p>
            <p class="time">{{ formatTrackingTimestampFa(entry.created_at) }}</p>
          </div>
          <SyncStateChip :state="entry.sync_state" :count="1" />
        </li>
      </ul>
    </section>

    <FailedSyncList :entries="failedRows" @retry-one="retryOne" @retry-all="retryAll" />
  </div>
</template>

<style scoped>
.tracking-history-page {
  padding: var(--spacing-lg);
  display: grid;
  gap: var(--spacing-md);
  direction: rtl;
}

header h1 {
  margin: 0;
}

header p {
  margin: 4px 0 0;
  color: #546067;
}

.rows-section {
  border: 1px solid #d4e0e6;
  border-radius: 10px;
  padding: var(--spacing-md);
  background: #ffffff;
}

.rows-section h2 {
  margin: 0 0 var(--spacing-sm);
  font-size: 1rem;
}

ul {
  margin: 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: var(--spacing-sm);
}

li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm);
  border-radius: 8px;
  background: #f8fbfd;
}

.domain {
  margin: 0;
  font-weight: 600;
}

.time {
  margin: 4px 0 0;
  font-size: 0.82rem;
  color: #56636c;
}

.empty {
  margin: 0;
}
</style>
