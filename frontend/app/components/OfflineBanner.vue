<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { db } from '~/db'

const isOffline = ref(false)
const evictionDetected = ref(false)

async function checkState() {
  if (typeof navigator === 'undefined') return
  isOffline.value = !navigator.onLine

  // D-23: check eviction flag written by clientPlan store
  try {
    const flag = await db.uiState.get('eviction_detected')
    evictionDetected.value = flag?.value === 'true'
  }
  catch { evictionDetected.value = false }
}

function handleOnline() { isOffline.value = false }
function handleOffline() { isOffline.value = true }

onMounted(async () => {
  await checkState()
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)
})

onUnmounted(() => {
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
})

async function dismissEviction() {
  await db.uiState.delete('eviction_detected')
  evictionDetected.value = false
}
</script>

<template>
  <!-- D-23: iOS eviction notice — shown above offline banner if eviction was detected -->
  <div
    v-if="evictionDetected"
    class="bg-amber-50 border-b border-amber-200 px-4 py-3 flex items-start gap-3"
  >
    <span class="text-amber-600 text-lg mt-0.5">⚠️</span>
    <div class="flex-1">
      <p class="text-sm text-amber-800 font-medium">داده‌های آفلاین پاک شدند</p>
      <p class="text-xs text-amber-700 mt-0.5">
        داده‌های آفلاین توسط دستگاه حذف شدند. پس از اتصال به اینترنت بازیابی می‌شوند.
      </p>
    </div>
    <button
      class="text-amber-500 hover:text-amber-700 text-lg leading-none"
      aria-label="بستن"
      @click="dismissEviction"
    >
      ×
    </button>
  </div>

  <!-- D-21: Offline mode indicator -->
  <div
    v-if="isOffline && !evictionDetected"
    class="bg-gray-800 text-white text-center text-sm py-2 px-4"
  >
    📵 حالت آفلاین — تغییرات ذخیره می‌شوند و پس از اتصال ارسال خواهند شد
  </div>
</template>
