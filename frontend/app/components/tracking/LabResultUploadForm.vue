<script setup lang="ts">
import type { LabType } from '~/types/tracking.types'
import { LAB_TYPE_LABELS } from '~/types/tracking.types'

const emit = defineEmits<{ uploaded: [] }>()
const labStore = useLabResultStore()

const title = ref('')
const labType = ref<LabType | ''>('')
const testDate = ref(new Date().toISOString().slice(0, 10))
const externalLink = ref('')
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)
const formError = ref<string | null>(null)

const options = Object.entries(LAB_TYPE_LABELS) as [LabType, string][]

async function submit() {
  if (!title.value.trim()) {
    formError.value = 'عنوان آزمایش الزامی است'
    return
  }
  if (!labType.value) {
    formError.value = 'نوع آزمایش الزامی است'
    return
  }
  if (!fileInput.value?.files?.length && !externalLink.value.trim()) {
    formError.value = 'فایل یا لینک الزامی است'
    return
  }

  formError.value = null
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('local_id', crypto.randomUUID())
    formData.append('title', title.value.trim())
    formData.append('lab_type', labType.value)
    formData.append('test_date', testDate.value)
    if (externalLink.value.trim()) formData.append('external_link', externalLink.value.trim())
    if (fileInput.value?.files?.[0]) formData.append('file', fileInput.value.files[0])
    await labStore.uploadLabResult(formData)
    title.value = ''
    labType.value = ''
    externalLink.value = ''
    testDate.value = new Date().toISOString().slice(0, 10)
    if (fileInput.value) fileInput.value.value = ''
    emit('uploaded')
  }
  catch (e: unknown) {
    const err = e as { data?: { error?: string } }
    formError.value = err.data?.error ?? 'خطا در بارگذاری آزمایش'
  }
  finally {
    uploading.value = false
  }
}
</script>

<template>
  <div class="rounded-2xl bg-white p-4 shadow-sm">
    <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">بارگذاری نتیجه آزمایش</h2>
    <input v-model="title" type="text" class="mb-3 w-full rounded-xl border p-3 text-start" placeholder="عنوان آزمایش" />
    <select v-model="labType" class="mb-3 w-full rounded-xl border bg-white p-3 text-start">
      <option value="" disabled>نوع آزمایش را انتخاب کنید</option>
      <option v-for="[value, label] in options" :key="value" :value="value">{{ label }}</option>
    </select>
    <input v-model="testDate" type="date" class="mb-3 w-full rounded-xl border p-3 text-start" />
    <input v-model="externalLink" type="url" class="mb-3 w-full rounded-xl border p-3 text-start" placeholder="لینک خارجی (اختیاری)" />
    <input ref="fileInput" type="file" accept="image/*,.pdf" class="mb-3 w-full rounded-xl border p-3 text-start" />
    <p v-if="formError" class="mb-3 text-sm text-rose-600 text-start">{{ formError }}</p>
    <button type="button" class="w-full rounded-xl bg-emerald-500 py-3 font-medium text-white disabled:opacity-50" :disabled="uploading" @click="submit">
      {{ uploading ? 'در حال بارگذاری...' : 'ثبت آزمایش' }}
    </button>
  </div>
</template>
