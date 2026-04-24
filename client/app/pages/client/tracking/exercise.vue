<script setup lang="ts">
import { ref } from 'vue'
import TrackingEntrySheet from '~/app/components/client/TrackingEntrySheet.vue'
import { useClientOfflineStore } from '~/app/stores/client-offline'

definePageMeta({ layout: 'client' })

const offlineStore = useClientOfflineStore()
const exerciseName = ref('پیاده روی')
const durationMinutes = ref(30)
const intensity = ref<'light' | 'moderate' | 'vigorous'>('moderate')
const notes = ref('')

function makeExerciseId(name: string): string {
  return `exercise-${name.trim().replace(/\s+/g, '-').toLowerCase()}`
}

async function submitExercise() {
  if (!exerciseName.value.trim() || Number(durationMinutes.value) <= 0) {
    return { ok: false, message: 'نام ورزش و مدت زمان معتبر الزامی است.' }
  }

  const entry = offlineStore.enqueueDomainTrackingWrite({
    domain: 'exercise',
    payload: {
      exercise_id: makeExerciseId(exerciseName.value),
      duration_minutes: Number(durationMinutes.value),
      logged_at: new Date().toISOString(),
      intensity: intensity.value,
      notes: notes.value.trim() || undefined,
    },
  })

  if (!entry) {
    return { ok: false, message: 'ثبت تکراری یا نامعتبر بود.' }
  }

  return { ok: true, message: 'ثبت ورزش در صف ارسال قرار گرفت.' }
}
</script>

<template>
  <div class="tracking-page">
    <TrackingEntrySheet
      title="ثبت ورزش"
      description="فعالیت بدنی روزانه را سریع ثبت کنید."
      submit-label="ثبت ورزش"
      :on-submit="submitExercise"
    >
      <label>
        نام ورزش
        <input v-model="exerciseName" type="text" />
      </label>

      <label>
        مدت (دقیقه)
        <input v-model.number="durationMinutes" type="number" min="1" />
      </label>

      <label>
        شدت
        <select v-model="intensity">
          <option value="light">سبک</option>
          <option value="moderate">متوسط</option>
          <option value="vigorous">سنگین</option>
        </select>
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
