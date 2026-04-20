<script setup lang="ts">
import type { UpsertBodyMeasurementPayload } from '~/types/tracking.types'

const emit = defineEmits<{
  submit: [payload: Omit<UpsertBodyMeasurementPayload, 'local_id' | 'date'>]
}>()

const form = reactive({
  weight_kg: '',
  waist_cm: '',
  hip_cm: '',
  abdomen_cm: '',
  thigh_cm: '',
  chest_cm: '',
  wrist_cm: '',
})
const showExtra = ref(false)
const formError = ref<string | null>(null)

function toNumber(value: string) {
  return value ? Number(value) : undefined
}

function submit() {
  if (!form.weight_kg.trim()) {
    formError.value = 'ثبت وزن الزامی است'
    return
  }
  formError.value = null
  emit('submit', {
    weight_kg: Number(form.weight_kg),
    waist_cm: toNumber(form.waist_cm),
    hip_cm: toNumber(form.hip_cm),
    abdomen_cm: toNumber(form.abdomen_cm),
    thigh_cm: toNumber(form.thigh_cm),
    chest_cm: toNumber(form.chest_cm),
    wrist_cm: toNumber(form.wrist_cm),
  })
  Object.assign(form, {
    weight_kg: '', waist_cm: '', hip_cm: '', abdomen_cm: '', thigh_cm: '', chest_cm: '', wrist_cm: '',
  })
  showExtra.value = false
}
</script>

<template>
  <div class="rounded-2xl bg-white p-4 shadow-sm">
    <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">ثبت اندازه‌گیری</h2>
    <label class="mb-1 block text-xs text-gray-500 text-start">وزن (کیلوگرم)</label>
    <input v-model="form.weight_kg" type="number" min="1" step="0.1" class="mb-3 w-full rounded-xl border p-3 text-start text-lg font-semibold" placeholder="۷۰.۵" />

    <button type="button" class="mb-3 text-sm text-emerald-700 underline underline-offset-4" @click="showExtra = !showExtra">
      {{ showExtra ? 'بستن اندازه‌گیری‌های بیشتر' : 'اندازه‌گیری‌های بیشتر' }}
    </button>

    <div v-if="showExtra" class="grid grid-cols-2 gap-3">
      <input v-model="form.waist_cm" type="number" min="1" step="0.1" class="rounded-xl border p-3 text-start" placeholder="دور کمر" />
      <input v-model="form.hip_cm" type="number" min="1" step="0.1" class="rounded-xl border p-3 text-start" placeholder="دور باسن" />
      <input v-model="form.abdomen_cm" type="number" min="1" step="0.1" class="rounded-xl border p-3 text-start" placeholder="شکم" />
      <input v-model="form.thigh_cm" type="number" min="1" step="0.1" class="rounded-xl border p-3 text-start" placeholder="ران" />
      <input v-model="form.chest_cm" type="number" min="1" step="0.1" class="rounded-xl border p-3 text-start" placeholder="سینه" />
      <input v-model="form.wrist_cm" type="number" min="1" step="0.1" class="rounded-xl border p-3 text-start" placeholder="مچ" />
    </div>

    <p v-if="formError" class="mt-3 text-sm text-rose-600 text-start">{{ formError }}</p>
    <button type="button" class="mt-4 w-full rounded-xl bg-emerald-500 py-3 font-medium text-white" @click="submit">
      ثبت اندازه‌گیری
    </button>
  </div>
</template>
