<script setup lang="ts">
import type { NotificationPrefs } from '~/stores/notificationPrefs'

definePageMeta({
  layout: 'client',
  middleware: ['role-guard'],
  roles: ['client'],
})

const prefsStore = useNotificationPrefsStore()
const { status, errorMessage, refreshStatus, requestAndSubscribe } = useNotificationPermission()

const notificationItems: Array<{
  key: keyof NotificationPrefs
  label: string
  icon: string
}> = [
  { key: 'new_message', label: 'پیام جدید از متخصص تغذیه', icon: '💬' },
  { key: 'plan_activated', label: 'فعال شدن برنامه غذایی جدید', icon: '📋' },
  { key: 'food_request_decision', label: 'نتیجه درخواست غذا', icon: '✅' },
  { key: 'meal_reminder', label: 'یادآور وعده غذایی', icon: '🍽️' },
  { key: 'medication_reminder', label: 'یادآور دارو', icon: '💊' },
  { key: 'water_reminder', label: 'یادآور نوشیدن آب', icon: '💧' },
]

onMounted(async () => {
  await Promise.all([
    prefsStore.fetchPreferences(),
    refreshStatus(),
  ])
})

async function handleToggle(key: keyof NotificationPrefs) {
  await prefsStore.updatePreferences({
    [key]: !prefsStore.prefs[key],
  })
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 py-6 pb-24">
    <div class="mb-4">
      <p class="text-sm text-gray-500">پروفایل / تنظیمات</p>
      <h1 class="mt-1 text-xl font-bold text-gray-800">
        تنظیمات اعلان‌ها
      </h1>
    </div>

    <section class="rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="text-sm font-semibold text-gray-700">
        دریافت اعلان
      </h2>

      <p v-if="status === 'subscribed'" class="mt-3 rounded-xl bg-green-50 px-3 py-3 text-sm text-green-700">
        اعلان‌ها برای این دستگاه فعال است.
      </p>
      <p
        v-else-if="status === 'ios-not-installed' || status === 'denied' || status === 'error'"
        class="mt-3 rounded-xl bg-amber-50 px-3 py-3 text-sm text-amber-700"
      >
        {{ errorMessage }}
      </p>
      <p v-else class="mt-3 text-sm text-gray-600">
        برای دریافت اعلان‌های پیام، برنامه غذایی و یادآورها، اعلان‌ها را فعال کنید.
      </p>

      <button
        v-if="status !== 'subscribed'"
        type="button"
        class="mt-4 w-full rounded-xl bg-emerald-600 px-4 py-3 text-sm font-medium text-white disabled:opacity-60"
        :disabled="status === 'requesting'"
        @click="requestAndSubscribe"
      >
        {{ status === 'requesting' ? 'در حال فعال‌سازی...' : 'فعال‌سازی اعلان‌ها' }}
      </button>
    </section>

    <section class="mt-4 rounded-2xl bg-white shadow-sm">
      <div
        v-for="item in notificationItems"
        :key="item.key"
        class="flex items-center justify-between gap-3 border-b border-gray-100 px-4 py-4 last:border-b-0"
      >
        <div class="flex items-center gap-3">
          <span class="text-xl">{{ item.icon }}</span>
          <span class="text-sm text-gray-800">{{ item.label }}</span>
        </div>

        <button
          type="button"
          class="relative inline-flex h-6 w-11 rounded-full transition-colors focus:outline-none"
          :class="prefsStore.prefs[item.key] ? 'bg-emerald-500' : 'bg-gray-300'"
          :disabled="prefsStore.loading || prefsStore.saving"
          @click="handleToggle(item.key)"
        >
          <span
            class="mt-0.5 inline-block h-5 w-5 rounded-full bg-white shadow-sm transition-transform"
            :class="prefsStore.prefs[item.key] ? '-translate-x-0.5' : 'translate-x-5'"
          />
        </button>
      </div>
    </section>

    <p v-if="prefsStore.error" class="mt-4 text-center text-sm text-red-600">
      {{ prefsStore.error }}
    </p>
  </div>
</template>
