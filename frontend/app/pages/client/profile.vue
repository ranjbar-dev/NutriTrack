<script setup lang="ts">
definePageMeta({
  layout: 'client',
  middleware: ['role-guard'],
  roles: ['client'],
})

const authStore = useAuthStore()

async function handleLogout() {
  await authStore.logout()
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 py-6 pb-24">
    <div class="rounded-2xl bg-white p-5 shadow-sm">
      <p class="text-sm text-gray-500">پروفایل مراجع</p>
      <h1 class="mt-2 text-xl font-bold text-gray-800">
        {{ authStore.user?.full_name ?? 'مراجع نوتری‌ترک' }}
      </h1>
      <p class="mt-2 text-sm text-gray-500">
        {{ authStore.user?.mobile ?? 'شماره موبایل ثبت نشده است' }}
      </p>
    </div>

    <div class="mt-4 rounded-2xl bg-white p-2 shadow-sm">
      <NuxtLink
        to="/client/settings/notifications"
        class="flex items-center justify-between rounded-xl px-3 py-4 text-sm text-gray-800 transition hover:bg-gray-50"
      >
        <span class="flex items-center gap-3">
          <span class="text-xl">🔔</span>
          <span>تنظیمات اعلان‌ها</span>
        </span>
        <span class="text-gray-400">‹</span>
      </NuxtLink>
    </div>

    <div class="mt-4 rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="font-bold text-gray-800">وضعیت حساب</h2>
      <p class="mt-2 text-sm text-gray-500">
        در صورت قطع اینترنت، ثبت‌های شما در برنامه نگه‌داری می‌شود و پس از اتصال ارسال خواهد شد.
      </p>
    </div>

    <button
      type="button"
      class="mt-6 min-h-[44px] w-full rounded-2xl bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
      @click="handleLogout"
    >
      خروج از حساب
    </button>
  </div>
</template>
