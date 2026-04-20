<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'nutritionist' })

const route = useRoute()
const clientId = computed(() => route.params.clientId as string)
const { dateFrom, dateTo, tabs, tabPath, isActive } = useNutriTracking(clientId)
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
    <div class="px-4 py-6">
      <div class="rounded-2xl bg-white p-5 shadow-sm text-start">
        <h1 class="text-lg font-bold text-gray-800">پیگیری بیمار</h1>
        <p class="mt-2 text-sm text-gray-500">برای مشاهده جزئیات هر بخش یکی از تب‌ها را انتخاب کنید.</p>
      </div>
    </div>
  </div>
</template>
