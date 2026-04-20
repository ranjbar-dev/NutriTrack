<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'nutritionist' })

const route = useRoute()
const clientId = computed(() => route.params.clientId as string)
const { dateFrom, dateTo, tabs, tabPath, isActive } = useNutriTracking(clientId)
const bodyStore = useBodyMeasurementStore()
const saveError = ref<string | null>(null)

async function load() {
  await bodyStore.fetchHistory(clientId.value, dateFrom.value, dateTo.value)
}

async function submitMeasurement(payload: Record<string, number | undefined>) {
  saveError.value = null
  try {
    await bodyStore.logMeasurement(payload, clientId.value)
  } catch (e: unknown) {
    const err = e as { data?: { error?: string } }
    saveError.value = err.data?.error ?? 'خطا در ثبت اندازه‌گیری'
  }
}

onMounted(load)
watch([clientId, dateFrom, dateTo], load)
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
    <div class="space-y-4 px-4 py-4 pb-20">
      <h1 class="text-lg font-bold text-gray-800 text-start">اندازه‌گیری بدن</h1>
      <BodyMeasurementForm @submit="submitMeasurement" />
      <p v-if="saveError" class="text-sm text-rose-600 text-start">{{ saveError }}</p>
      <WeightChart :history="bodyStore.history" />
      <div class="rounded-2xl bg-white p-4 shadow-sm">
        <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">سوابق بیمار</h2>
        <div v-if="bodyStore.loading" class="py-6 text-center text-sm text-gray-400">در حال بارگذاری...</div>
        <div v-else-if="bodyStore.history.length === 0" class="py-6 text-center text-sm text-gray-400">داده‌ای موجود نیست</div>
        <div v-else class="space-y-2">
          <div v-for="item in bodyStore.history" :key="item.id" class="rounded-xl bg-gray-50 p-3">
            <div class="flex items-center justify-between gap-2">
              <p class="text-sm font-medium text-gray-800">{{ item.weight_kg ?? '—' }} کیلوگرم</p>
              <span class="text-xs text-gray-400">{{ item.date }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
