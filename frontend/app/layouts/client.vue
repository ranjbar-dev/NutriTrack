<script setup lang="ts">
import { useSyncManager } from '~/composables/useSyncManager'
import { usePWA } from '~/composables/usePWA.client'

const { processQueue } = useSyncManager()

const { canInstall, needsUpdate, promptInstall, applyUpdate } = usePWA()

const navItems = [
  { label: 'برنامه', icon: 'plan', to: '/client/plan' },
  { label: 'ثبت', icon: 'tracking', to: '/client/tracking' },
  { label: 'پیام‌ها', icon: 'chat', to: '/client/messages' },
  { label: 'پروفایل', icon: 'user', to: '/client/profile' },
]
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <OfflineBanner />
    <ClientSyncStatus />

    <slot />

    <UiBottomNav :items="navItems" />

    <div
      v-if="canInstall"
      class="fixed bottom-20 start-4 end-4 bg-green-600 text-white rounded-xl p-4 flex items-center gap-3 shadow-lg z-40"
    >
      <span class="text-2xl">📱</span>
      <div class="flex-1">
        <p class="text-sm font-bold">نصب نوتری‌ترک</p>
        <p class="text-xs opacity-90">برنامه را روی صفحه اصلی نصب کنید</p>
      </div>
      <button class="text-xs bg-white text-green-700 px-3 py-1.5 rounded-lg font-medium" @click="promptInstall">
        نصب
      </button>
    </div>

    <div
      v-if="needsUpdate"
      class="fixed bottom-20 start-4 end-4 bg-blue-600 text-white rounded-xl p-4 flex items-center gap-3 shadow-lg z-40"
    >
      <span class="text-2xl">🔄</span>
      <div class="flex-1">
        <p class="text-sm font-bold">بروزرسانی جدید</p>
        <p class="text-xs opacity-90">نسخه جدیدی در دسترس است</p>
      </div>
      <button class="text-xs bg-white text-blue-700 px-3 py-1.5 rounded-lg font-medium" @click="applyUpdate">
        بروزرسانی
      </button>
    </div>
  </div>
</template>
