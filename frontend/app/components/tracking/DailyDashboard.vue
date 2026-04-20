<script setup lang="ts">
import { SLEEP_QUALITY_LABELS } from '~/types/tracking.types'

const trackingStore = useTrackingStore()
const { formatShamsi } = useShamsiDate()

const quickLinks = [
  { label: 'ثبت غذا', to: '/client/tracking/food' },
  { label: 'ثبت آب', to: '/client/tracking/water' },
  { label: 'ثبت خواب', to: '/client/tracking/sleep' },
  { label: 'ثبت تمرین', to: '/client/tracking/exercise' },
  { label: 'داروها', to: '/client/tracking/medication' },
  { label: 'اندازه‌گیری', to: '/client/tracking/body' },
  { label: 'آزمایش‌ها', to: '/client/tracking/lab-results' },
]
</script>

<template>
  <div class="space-y-3">
    <div class="rounded-2xl bg-white p-4 shadow-sm">
      <p class="text-xs text-gray-500 text-start">امروز</p>
      <h1 class="mt-1 text-lg font-bold text-gray-800 text-start">{{ formatShamsi(new Date(), 'long') }}</h1>
    </div>

    <div v-if="trackingStore.loading" class="space-y-3">
      <div v-for="index in 4" :key="index" class="h-24 animate-pulse rounded-2xl bg-gray-100" />
    </div>

    <div v-else-if="trackingStore.error" class="rounded-2xl bg-white p-4 shadow-sm">
      <p class="text-sm text-rose-600 text-start">{{ trackingStore.error }}</p>
      <button class="mt-3 rounded-xl bg-emerald-500 px-4 py-2 text-sm text-white" @click="trackingStore.fetchDailyDashboard()">تلاش دوباره</button>
    </div>

    <template v-else-if="trackingStore.dashboard">
      <WaterProgressBar :total-ml="trackingStore.dashboard.water_total_ml" :target-ml="trackingStore.dashboard.water_target_ml ?? 0" />

      <div class="grid grid-cols-2 gap-3">
        <div class="rounded-2xl bg-white p-4 shadow-sm">
          <p class="text-xs text-gray-500 text-start">وعده‌ها</p>
          <p class="mt-2 text-lg font-bold text-gray-800">{{ trackingStore.dashboard.meals_logged }} / {{ trackingStore.dashboard.meals_total }}</p>
          <NuxtLink to="/client/tracking/food" class="mt-3 inline-block text-xs text-emerald-700">مشاهده و ثبت</NuxtLink>
        </div>
        <div class="rounded-2xl bg-white p-4 shadow-sm">
          <p class="text-xs text-gray-500 text-start">تمرین</p>
          <p class="mt-2 text-lg font-bold text-gray-800">{{ trackingStore.dashboard.exercise_count }}</p>
          <NuxtLink to="/client/tracking/exercise" class="mt-3 inline-block text-xs text-emerald-700">ثبت تمرین</NuxtLink>
        </div>
        <div class="rounded-2xl bg-white p-4 shadow-sm">
          <p class="text-xs text-gray-500 text-start">دارو</p>
          <p class="mt-2 text-lg font-bold text-gray-800">{{ trackingStore.dashboard.medication_taken_count }}</p>
          <NuxtLink to="/client/tracking/medication" class="mt-3 inline-block text-xs text-emerald-700">مشاهده داروها</NuxtLink>
        </div>
        <div class="rounded-2xl bg-white p-4 shadow-sm">
          <p class="text-xs text-gray-500 text-start">اندازه‌گیری</p>
          <p class="mt-2 text-sm font-semibold" :class="trackingStore.dashboard.body_logged_today ? 'text-emerald-700' : 'text-gray-700'">
            {{ trackingStore.dashboard.body_logged_today ? 'امروز ثبت شده' : 'هنوز ثبت نشده' }}
          </p>
          <NuxtLink to="/client/tracking/body" class="mt-3 inline-block text-xs text-emerald-700">ثبت اندازه‌گیری</NuxtLink>
        </div>
      </div>

      <div class="rounded-2xl bg-white p-4 shadow-sm">
        <p class="text-xs text-gray-500 text-start">خواب</p>
        <template v-if="trackingStore.dashboard.sleep_log">
          <p class="mt-2 text-sm font-semibold text-gray-800">
            {{ trackingStore.dashboard.sleep_log.sleep_time }} تا {{ trackingStore.dashboard.sleep_log.wake_time }}
          </p>
          <p class="mt-1 text-xs text-gray-500">کیفیت: {{ SLEEP_QUALITY_LABELS[trackingStore.dashboard.sleep_log.quality] }}</p>
        </template>
        <template v-else>
          <p class="mt-2 text-sm text-gray-500">خوابی برای امروز ثبت نشده است</p>
        </template>
        <NuxtLink to="/client/tracking/sleep" class="mt-3 inline-block text-xs text-emerald-700">ثبت خواب</NuxtLink>
      </div>

      <div class="rounded-2xl bg-white p-4 shadow-sm">
        <p class="mb-3 text-sm font-semibold text-gray-700 text-start">اقدام سریع</p>
        <div class="grid grid-cols-2 gap-2">
          <NuxtLink
            v-for="link in quickLinks"
            :key="link.to"
            :to="link.to"
            class="rounded-xl bg-gray-50 px-3 py-3 text-sm text-gray-700 text-center"
          >
            {{ link.label }}
          </NuxtLink>
        </div>
      </div>
    </template>
  </div>
</template>
