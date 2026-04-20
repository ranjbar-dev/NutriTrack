<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useSyncQueueStore } from '~/stores/syncQueue'

const syncStore = useSyncQueueStore()
const { pendingCount, failedCount, isProcessing } = storeToRefs(syncStore)

const isVisible = computed(() => isProcessing.value || pendingCount.value > 0 || failedCount.value > 0)
</script>

<template>
  <div v-if="isVisible" class="px-4 py-2 flex items-center gap-2 bg-white border-b border-gray-100">
    <template v-if="isProcessing">
      <span class="inline-block w-3 h-3 rounded-full bg-green-500 animate-pulse" />
      <span class="text-xs text-gray-600">همگام‌سازی در حال انجام</span>
    </template>
    <template v-else-if="pendingCount > 0">
      <span class="inline-block w-3 h-3 rounded-full bg-amber-400" />
      <span class="text-xs text-gray-600">{{ pendingCount }} مورد در انتظار ارسال</span>
    </template>
    <template v-if="failedCount > 0">
      <span class="me-auto flex items-center gap-1">
        <span class="text-xs text-red-600">{{ failedCount }} خطا</span>
        <button
          class="text-xs text-green-600 underline font-medium"
          @click="syncStore.retryFailed()"
        >
          تلاش مجدد
        </button>
      </span>
    </template>
  </div>
</template>
