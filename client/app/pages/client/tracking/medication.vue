<script setup lang="ts">
import { ref } from 'vue'
import TrackingEntrySheet from '~/app/components/client/TrackingEntrySheet.vue'
import { useClientOfflineStore } from '~/app/stores/client-offline'

definePageMeta({ layout: 'client' })

const offlineStore = useClientOfflineStore()
const medicationName = ref('')
const doses = ref(1)
const loggedAt = ref(new Date().toISOString())
const notes = ref('')

function makeMedicationId(name: string): string {
  return `med-${name.trim().replace(/\s+/g, '-').toLowerCase()}`
}

async function submitMedication() {
  if (!medicationName.value.trim() || Number(doses.value) <= 0) {
    return { ok: false, message: 'نام دارو و تعداد دوز معتبر الزامی است.' }
  }

  const entry = offlineStore.enqueueDomainTrackingWrite({
    domain: 'medication',
    payload: {
      medication_id: makeMedicationId(medicationName.value),
      doses: Number(doses.value),
      logged_at: loggedAt.value,
      notes: notes.value.trim() || `ثبت دستی دارو: ${medicationName.value}`,
    },
  })

  if (!entry) {
    return { ok: false, message: 'ثبت تکراری یا نامعتبر بود.' }
  }

  return { ok: true, message: 'ثبت دارو در صف ارسال قرار گرفت.' }
}
</script>

<template>
  <div class="tracking-page">
    <TrackingEntrySheet
      title="ثبت دارو"
      description="ثبت نسخه و داروی مصرف شده با امکان نام گذاری دستی"
      submit-label="ثبت دارو"
      :on-submit="submitMedication"
    >
      <label>
        نام دارو
        <input v-model="medicationName" type="text" placeholder="مثال: ویتامین D" />
      </label>

      <label>
        تعداد دوز
        <input v-model.number="doses" type="number" min="1" />
      </label>

      <label>
        زمان مصرف
        <input v-model="loggedAt" type="datetime-local" />
      </label>

      <label>
        یادداشت
        <textarea v-model="notes" rows="2"></textarea>
      </label>
    </TrackingEntrySheet>
  </div>
</template>

<style scoped>
.tracking-page {
  padding: var(--spacing-lg);
}
</style>
