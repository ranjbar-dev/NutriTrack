<script setup lang="ts">
import { ref } from 'vue'
import TrackingEntrySheet from '~/app/components/client/TrackingEntrySheet.vue'
import { createFoodPayload } from '~/app/lib/tracking/entry'
import { useClientOfflineStore } from '~/app/stores/client-offline'

definePageMeta({ layout: 'client' })

const offlineStore = useClientOfflineStore()
const foodName = ref('')
const quantityGrams = ref(120)
const consumedAt = ref(new Date().toISOString())
const notes = ref('')

function buildFoodId(name: string): string {
  return `food-${name.trim().replace(/\s+/g, '-').toLowerCase()}`
}

async function submitFood() {
  const result = createFoodPayload({
    foodId: buildFoodId(foodName.value),
    quantityGrams: Number(quantityGrams.value),
    consumedAt: consumedAt.value,
    notes: notes.value,
  })

  if (!result.ok || !result.payload) {
    return {
      ok: false,
      message: result.error ?? 'ثبت غذا انجام نشد.',
    }
  }

  const entry = offlineStore.enqueueDomainTrackingWrite({
    domain: 'food',
    payload: result.payload,
  })

  if (!entry) {
    return { ok: false, message: 'ثبت تکراری یا نامعتبر بود.' }
  }

  consumedAt.value = new Date().toISOString()
  return { ok: true, message: 'ثبت غذا در صف ارسال قرار گرفت.' }
}
</script>

<template>
  <div class="tracking-page">
    <TrackingEntrySheet
      title="ثبت غذا"
      description="برای ثبت سریع وعده امروز، مقدار و زمان را وارد کنید."
      submit-label="ثبت غذا"
      :on-submit="submitFood"
    >
      <label>
        نام غذا
        <input v-model="foodName" type="text" placeholder="مثال: نان و پنیر" />
      </label>

      <label>
        مقدار (گرم)
        <input v-model.number="quantityGrams" type="number" min="1" step="1" />
      </label>

      <label>
        زمان مصرف
        <input v-model="consumedAt" type="datetime-local" />
      </label>

      <label>
        یادداشت
        <textarea v-model="notes" rows="2" placeholder="اختیاری"></textarea>
      </label>
    </TrackingEntrySheet>
  </div>
</template>

<style scoped>
.tracking-page {
  padding: var(--spacing-lg);
}
</style>
