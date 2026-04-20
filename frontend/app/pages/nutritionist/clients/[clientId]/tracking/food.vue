<script setup lang="ts">
import type { FoodLogEntry } from '~/types/tracking.types'

definePageMeta({ middleware: ['auth'], layout: 'nutritionist' })

const route = useRoute()
const clientId = computed(() => route.params.clientId as string)
const { dateFrom, dateTo, tabs, tabPath, isActive, fetchDomain } = useNutriTracking(clientId)
const items = ref<FoodLogEntry[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    items.value = await fetchDomain('food')
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch([dateFrom, dateTo], load)
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
      <h1 class="mb-3 text-lg font-bold text-gray-800 text-start">ثبت غذا</h1>
      <div v-if="loading" class="py-6 text-center text-sm text-gray-400">در حال بارگذاری...</div>
      <div v-else-if="items.length === 0" class="rounded-2xl bg-white p-6 text-center text-sm text-gray-400 shadow-sm">داده‌ای موجود نیست</div>
      <div v-else class="space-y-2">
        <div v-for="item in items" :key="item.id" class="rounded-2xl bg-white p-4 shadow-sm">
          <div class="flex items-center justify-between gap-2">
            <div class="flex-1 min-w-0 text-start"><p class='text-sm text-gray-800'>وعده: {{ item.meal_id }}</p><p class='mt-1 text-xs text-gray-500'>{{ item.is_skipped ? 'وعده رد شده' : ('گزینه انتخابی: ' + (item.selected_option_id || 'ثبت شده')) }}</p></div>
            <span class="text-xs text-gray-400">{{ item.date }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
