<script setup lang="ts">
import { ref } from 'vue'
import TrackingEntrySheet from '~/app/components/client/TrackingEntrySheet.vue'
import { createWaterPayload } from '~/app/lib/tracking/entry'
import { useClientOfflineStore } from '~/app/stores/client-offline'

definePageMeta({ layout: 'client' })

const offlineStore = useClientOfflineStore()
const amountMl = ref(250)
const quickValues = [250, 500, 750]

async function queueWater(amount: number) {
  const result = createWaterPayload({ amountMl: Number(amount), loggedAt: new Date().toISOString() })

  if (!result.ok || !result.payload) {
    return { ok: false, message: result.error ?? 'ثبت آب انجام نشد.' }
  }

  const entry = offlineStore.enqueueDomainTrackingWrite({
    domain: 'water',
    payload: result.payload,
  })

  if (!entry) {
    return { ok: false, message: 'ثبت تکراری یا نامعتبر بود.' }
  }

  return { ok: true, message: 'ثبت آب در صف ارسال قرار گرفت.' }
}

async function submitWater() {
  return queueWater(amountMl.value)
}
</script>

<template>
  <div class="tracking-page">
    <TrackingEntrySheet
      title="ثبت آب"
      description="افزودن سریع با یک لمس برای ثبت های پیاپی روزانه"
      submit-label="ثبت دستی"
      :on-submit="submitWater"
    >
      <div class="quick-list">
        <button
          v-for="value in quickValues"
          :key="value"
          class="quick-button"
          type="button"
          @click="queueWater(value)"
        >
          {{ value }}ml
        </button>
      </div>

      <label>
        مقدار دستی (ml)
        <input v-model.number="amountMl" type="number" min="1" step="50" />
      </label>
    </TrackingEntrySheet>
  </div>
</template>

<style scoped>
.tracking-page {
  padding: var(--spacing-lg);
}

.quick-list {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--spacing-sm);
}

.quick-button {
  min-height: 42px;
  border: 1px solid #98bac4;
  border-radius: 8px;
  background: #ecf8fc;
  font-weight: 600;
}
</style>
