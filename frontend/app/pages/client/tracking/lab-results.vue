<script setup lang="ts">
import { LAB_TYPE_LABELS } from '~/types/tracking.types'

definePageMeta({ middleware: ['auth'], layout: 'client' })

const labStore = useLabResultStore()

onMounted(() => labStore.fetchLabResults())
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6 space-y-4">
    <h1 class="text-lg font-bold text-gray-800 text-start">نتایج آزمایش</h1>
    <LabResultUploadForm @uploaded="labStore.fetchLabResults()" />
    <div class="rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">آزمایش‌های بارگذاری‌شده</h2>
      <div v-if="labStore.loading" class="py-6 text-center text-sm text-gray-400">در حال بارگذاری...</div>
      <div v-else-if="labStore.labResults.length === 0" class="py-6 text-center text-sm text-gray-400">هنوز آزمایشی ثبت نشده است</div>
      <div v-else class="space-y-2">
        <div v-for="result in labStore.labResults" :key="result.id" class="rounded-xl bg-gray-50 p-3">
          <div class="flex items-center justify-between gap-2">
            <div>
              <p class="text-sm font-medium text-gray-800">{{ result.title }}</p>
              <p class="mt-1 text-xs text-gray-500">{{ LAB_TYPE_LABELS[result.lab_type] }}</p>
            </div>
            <span class="text-xs text-gray-400">{{ result.test_date }}</span>
          </div>
          <p v-if="result.has_file" class="mt-2 text-xs text-emerald-700">فایل آزمایش بارگذاری شده است</p>
          <a v-if="result.external_link" :href="result.external_link" target="_blank" rel="noopener noreferrer" class="mt-2 inline-block text-xs text-blue-600 underline">
            مشاهده لینک خارجی
          </a>
        </div>
      </div>
    </div>
  </div>
</template>
