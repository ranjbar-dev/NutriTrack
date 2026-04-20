<script setup lang="ts">
import { LAB_TYPE_LABELS, type LabResultResponse } from '~/types/tracking.types'

definePageMeta({ middleware: ['auth'], layout: 'nutritionist' })

const route = useRoute()
const clientId = computed(() => route.params.clientId as string)
const { dateFrom, dateTo, tabs, tabPath, isActive } = useNutriTracking(clientId)
const labStore = useLabResultStore()
const downloadError = ref<string | null>(null)

async function load() {
  await labStore.fetchLabResults(clientId.value)
}

async function download(result: LabResultResponse) {
  downloadError.value = null
  try {
    await labStore.downloadLabResult(clientId.value, result)
  } catch (e: unknown) {
    const err = e as { message?: string }
    downloadError.value = err.message ?? 'خطا در دانلود فایل'
  }
}

onMounted(load)
watch(clientId, load)
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-white border-b border-gray-200 px-4 py-3 sticky top-0 z-10">
      <div class="flex gap-2 overflow-x-auto pb-2">
        <NuxtLink
          v-for="tab in tabs"
          :key="tab.key"
          :to="tabPath(tab.path)"
          :class="[
            'flex-shrink-0 rounded-xl px-3 py-2 text-sm font-medium',
            isActive(tab.path) ? 'bg-emerald-100 text-emerald-700' : 'bg-gray-100 text-gray-600',
          ]"
        >
          {{ tab.label }}
        </NuxtLink>
      </div>
      <div class="mt-3 flex gap-2">
        <input v-model="dateFrom" type="date" class="flex-1 rounded-xl border p-2 text-sm text-start" />
        <input v-model="dateTo" type="date" class="flex-1 rounded-xl border p-2 text-sm text-start" />
      </div>
    </div>
    <div class="px-4 py-4 pb-20">
      <h1 class="mb-3 text-lg font-bold text-gray-800 text-start">نتایج آزمایش بیمار</h1>
      <p v-if="downloadError" class="mb-3 text-sm text-rose-600 text-start">{{ downloadError }}</p>
      <div v-if="labStore.loading" class="py-6 text-center text-sm text-gray-400">در حال بارگذاری...</div>
      <div v-else-if="labStore.labResults.length === 0" class="rounded-2xl bg-white p-6 text-center text-sm text-gray-400 shadow-sm">داده‌ای موجود نیست</div>
      <div v-else class="space-y-2">
        <div v-for="result in labStore.labResults" :key="result.id" class="rounded-2xl bg-white p-4 shadow-sm">
          <div class="flex items-center justify-between gap-2">
            <div class="text-start">
              <p class="text-sm font-medium text-gray-800">{{ result.title }}</p>
              <p class="mt-1 text-xs text-gray-500">{{ LAB_TYPE_LABELS[result.lab_type] }}</p>
            </div>
            <span class="text-xs text-gray-400">{{ result.test_date }}</span>
          </div>
          <button v-if="result.has_file" type="button" class="mt-3 text-xs text-blue-600 underline" @click="download(result)">دانلود فایل</button>
          <a v-if="result.external_link" :href="result.external_link" target="_blank" rel="noopener noreferrer" class="mt-3 ms-3 inline-block text-xs text-emerald-700 underline">مشاهده لینک</a>
        </div>
      </div>
    </div>
  </div>
</template>
