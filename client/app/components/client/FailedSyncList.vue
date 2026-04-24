<script setup lang="ts">
import type { TrackingQueueEntry } from '~/app/types/offline-sync'
import { getFailureGuidance } from '~/app/lib/tracking/retry'

const props = defineProps<{
  entries: TrackingQueueEntry[]
}>()

const emit = defineEmits<{
  retryOne: [localId: string]
  retryAll: []
}>()
</script>

<template>
  <section class="failed-list">
    <header class="failed-header">
      <h2>موارد ناموفق</h2>
      <button v-if="props.entries.length > 1" type="button" @click="emit('retryAll')">تلاش مجدد همه</button>
    </header>

    <p v-if="props.entries.length === 0" class="empty">مورد ناموفقی برای ارسال وجود ندارد.</p>

    <ul v-else>
      <li v-for="entry in props.entries" :key="entry.local_id">
        <div>
          <p class="domain">{{ entry.domain }}</p>
          <p class="guidance">{{ getFailureGuidance(entry) }}</p>
        </div>
        <button type="button" @click="emit('retryOne', entry.local_id)">ارسال مجدد</button>
      </li>
    </ul>
  </section>
</template>

<style scoped>
.failed-list {
  margin-top: var(--spacing-lg);
  padding: var(--spacing-md);
  border-radius: 10px;
  border: 1px solid #f5c1c1;
  background: #fff7f7;
}

.failed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

h2 {
  margin: 0;
  font-size: 1rem;
  color: #8b1f1f;
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
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  background: #ffffff;
  border-radius: 8px;
}

.domain {
  margin: 0;
  font-weight: 600;
}

.guidance {
  margin: 2px 0 0;
  color: #6b1f1f;
  font-size: 0.82rem;
}

button {
  border: none;
  border-radius: 8px;
  padding: 8px 10px;
  background: #bf2f2f;
  color: #ffffff;
  font-size: 0.8rem;
}

.empty {
  margin: 0;
  color: #535d64;
}
</style>
