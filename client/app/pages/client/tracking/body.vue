<script setup lang="ts">
import { ref } from 'vue'
import TrackingEntrySheet from '~/app/components/client/TrackingEntrySheet.vue'
import { createBodyPayload } from '~/app/lib/tracking/entry'
import { useClientOfflineStore } from '~/app/stores/client-offline'

definePageMeta({ layout: 'client' })

const offlineStore = useClientOfflineStore()
const weightKg = ref<number | null>(null)
const waistCm = ref<number | null>(null)
const hipCm = ref<number | null>(null)
const abdomenCm = ref<number | null>(null)
const thighCm = ref<number | null>(null)
const chestCm = ref<number | null>(null)
const wristCm = ref<number | null>(null)

async function submitBody() {
  const result = createBodyPayload({
    loggedAt: new Date().toISOString(),
    weight_kg: weightKg.value,
    waist_cm: waistCm.value,
    hip_cm: hipCm.value,
    abdomen_cm: abdomenCm.value,
    thigh_cm: thighCm.value,
    chest_cm: chestCm.value,
    wrist_cm: wristCm.value,
  })

  if (!result.ok || !result.payload) {
    return { ok: false, message: result.error ?? 'ثبت اندازه بدن انجام نشد.' }
  }

  const entry = offlineStore.enqueueDomainTrackingWrite({
    domain: 'body',
    payload: result.payload,
  })

  if (!entry) {
    return { ok: false, message: 'ثبت تکراری یا نامعتبر بود.' }
  }

  return { ok: true, message: 'اندازه بدن در صف ارسال قرار گرفت.' }
}
</script>

<template>
  <div class="tracking-page">
    <TrackingEntrySheet
      title="ثبت اندازه بدن"
      description="حداقل یک اندازه را ثبت کنید. سایر فیلدها اختیاری هستند."
      submit-label="ثبت اندازه"
      :on-submit="submitBody"
    >
      <label>
        وزن (کیلوگرم)
        <input v-model.number="weightKg" type="number" min="0" step="0.1" />
      </label>

      <label>
        دور کمر (سانتی متر)
        <input v-model.number="waistCm" type="number" min="0" step="0.1" />
      </label>

      <label>
        دور باسن (سانتی متر)
        <input v-model.number="hipCm" type="number" min="0" step="0.1" />
      </label>

      <label>
        دور شکم (سانتی متر)
        <input v-model.number="abdomenCm" type="number" min="0" step="0.1" />
      </label>

      <label>
        دور ران (سانتی متر)
        <input v-model.number="thighCm" type="number" min="0" step="0.1" />
      </label>

      <label>
        دور سینه (سانتی متر)
        <input v-model.number="chestCm" type="number" min="0" step="0.1" />
      </label>

      <label>
        دور مچ (سانتی متر)
        <input v-model.number="wristCm" type="number" min="0" step="0.1" />
      </label>
    </TrackingEntrySheet>
  </div>
</template>

<style scoped>
.tracking-page {
  padding: var(--spacing-lg);
}
</style>
