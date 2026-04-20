<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'client' })

const bodyStore = useBodyMeasurementStore()
const saving = ref(false)
const saveError = ref<string | null>(null)

onMounted(() => bodyStore.fetchHistory(undefined, undefined, undefined))

async function submitMeasurement(payload: { [key: string]: number | undefined }) {
  saveError.value = null
  saving.value = true
  try {
    await bodyStore.logMeasurement(payload)
  }
  catch (e: unknown) {
    const err = e as { data?: { error?: string } }
    saveError.value = err.data?.error ?? 'خطا در ذخیره اندازه‌گیری'
  }
  finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6 space-y-4">
    <h1 class="text-lg font-bold text-gray-800 text-start">اندازه‌گیری بدن</h1>
    <BodyMeasurementForm @submit="submitMeasurement" />
    <p v-if="saveError" class="text-sm text-rose-600 text-start">{{ saveError }}</p>
    <WeightChart :history="bodyStore.history" />
    <div class="rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">سوابق ثبت‌شده</h2>
      <div v-if="bodyStore.history.length === 0" class="py-6 text-center text-sm text-gray-400">هنوز اندازه‌گیری ثبت نشده است</div>
      <div v-else class="space-y-2">
        <div v-for="item in bodyStore.history" :key="item.id" class="rounded-xl bg-gray-50 p-3">
          <div class="flex items-center justify-between gap-2">
            <p class="text-sm font-medium text-gray-800">{{ item.weight_kg ?? '—' }} کیلوگرم</p>
            <span class="text-xs text-gray-400">{{ item.date }}</span>
          </div>
          <div class="mt-2 flex flex-wrap gap-2 text-xs text-gray-500">
            <span v-if="item.waist_cm">کمر: {{ item.waist_cm }}</span>
            <span v-if="item.hip_cm">باسن: {{ item.hip_cm }}</span>
            <span v-if="item.abdomen_cm">شکم: {{ item.abdomen_cm }}</span>
            <span v-if="item.thigh_cm">ران: {{ item.thigh_cm }}</span>
            <span v-if="item.chest_cm">سینه: {{ item.chest_cm }}</span>
            <span v-if="item.wrist_cm">مچ: {{ item.wrist_cm }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
