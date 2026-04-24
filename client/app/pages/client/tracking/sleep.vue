<script setup lang="ts">
import { ref } from 'vue'
import TrackingEntrySheet from '~/app/components/client/TrackingEntrySheet.vue'
import { createSleepPayload } from '~/app/lib/tracking/entry'
import { useClientOfflineStore } from '~/app/stores/client-offline'

definePageMeta({ layout: 'client' })

const offlineStore = useClientOfflineStore()
const sleepStart = ref('23:30')
const wakeTime = ref('06:30')
const quality = ref<'poor' | 'fair' | 'good' | 'excellent'>('good')

async function submitSleep() {
  const result = createSleepPayload({
    sleepStart: sleepStart.value,
    wakeTime: wakeTime.value,
    quality: quality.value,
    loggedAt: new Date().toISOString(),
  })

  if (!result.ok || !result.payload) {
    return { ok: false, message: result.error ?? 'ثبت خواب انجام نشد.' }
  }

  const entry = offlineStore.enqueueDomainTrackingWrite({
    domain: 'sleep',
    payload: result.payload,
  })

  if (!entry) {
    return { ok: false, message: 'ثبت تکراری یا نامعتبر بود.' }
  }

  return { ok: true, message: 'خواب در صف ارسال قرار گرفت.' }
}
</script>

<template>
  <div class="tracking-page">
    <TrackingEntrySheet
      title="ثبت خواب"
      description="بازه خواب شب قبل را وارد کنید تا مدت زمان خودکار محاسبه شود."
      submit-label="ثبت خواب"
      :on-submit="submitSleep"
    >
      <label>
        زمان خواب
        <input v-model="sleepStart" type="time" />
      </label>

      <label>
        زمان بیداری
        <input v-model="wakeTime" type="time" />
      </label>

      <label>
        کیفیت خواب
        <select v-model="quality">
          <option value="poor">ضعیف</option>
          <option value="fair">متوسط</option>
          <option value="good">خوب</option>
          <option value="excellent">عالی</option>
        </select>
      </label>
    </TrackingEntrySheet>
  </div>
</template>

<style scoped>
.tracking-page {
  padding: var(--spacing-lg);
}
</style>
