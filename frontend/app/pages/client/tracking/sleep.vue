<script setup lang="ts">
import { SLEEP_QUALITY_LABELS, type SleepQuality } from '~/types/tracking.types'
import { computeSleepDurationMinutes, formatSleepDuration } from '~/utils/tracking'

definePageMeta({ middleware: ['auth'], layout: 'client' })

const sleepStore = useSleepLogStore()
const sleepTime = ref('')
const wakeTime = ref('')
const quality = ref<SleepQuality>('good')
const notes = ref('')
const formError = ref<string | null>(null)
const saving = ref(false)

onMounted(async () => {
  await sleepStore.fetchToday()
  if (sleepStore.todayLog) {
    sleepTime.value = sleepStore.todayLog.sleep_time.slice(0, 5)
    wakeTime.value = sleepStore.todayLog.wake_time.slice(0, 5)
    quality.value = sleepStore.todayLog.quality
    notes.value = sleepStore.todayLog.notes || ''
  }
})

const durationMinutes = computed(() => computeSleepDurationMinutes(sleepTime.value, wakeTime.value))
const durationLabel = computed(() => formatSleepDuration(durationMinutes.value))

async function save() {
  if (!sleepTime.value || !wakeTime.value) {
    formError.value = 'زمان خواب و بیداری الزامی است'
    return
  }
  formError.value = null
  saving.value = true
  try {
    await sleepStore.upsertSleep({
      date: new Date().toISOString().slice(0, 10),
      sleep_time: sleepTime.value,
      wake_time: wakeTime.value,
      quality: quality.value,
      ...(notes.value ? { notes: notes.value } : {}),
    })
  }
  finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6">
    <h1 class="mb-4 text-lg font-bold text-gray-800 text-start">ثبت خواب</h1>
    <div class="rounded-2xl bg-white p-4 shadow-sm">
      <label class="mb-1 block text-xs text-gray-500 text-start">زمان خواب</label>
      <input v-model="sleepTime" type="time" class="mb-3 w-full rounded-xl border p-3 text-start" />
      <label class="mb-1 block text-xs text-gray-500 text-start">زمان بیداری</label>
      <input v-model="wakeTime" type="time" class="mb-3 w-full rounded-xl border p-3 text-start" />
      <p class="mb-3 text-sm text-gray-600 text-start">مدت خواب: <span class="font-semibold text-gray-800">{{ durationLabel }}</span></p>
      <div class="mb-3 grid grid-cols-3 gap-2">
        <button
          v-for="(label, key) in SLEEP_QUALITY_LABELS"
          :key="key"
          type="button"
          :class="['rounded-xl border px-3 py-2 text-sm', quality === key ? 'border-blue-500 bg-blue-500 text-white' : 'border-gray-200 bg-white text-gray-700']"
          @click="quality = key as SleepQuality"
        >
          {{ label }}
        </button>
      </div>
      <textarea v-model="notes" rows="3" class="mb-3 w-full rounded-xl border p-3 text-start resize-none" placeholder="یادداشت (اختیاری)" />
      <p v-if="formError" class="mb-3 text-sm text-rose-600 text-start">{{ formError }}</p>
      <button type="button" class="w-full rounded-xl bg-blue-500 py-3 font-medium text-white disabled:opacity-50" :disabled="saving" @click="save">
        {{ saving ? 'در حال ذخیره...' : 'ذخیره خواب' }}
      </button>
    </div>
  </div>
</template>
