<script setup lang="ts">
defineProps<{ open: boolean }>()
const emit = defineEmits<{ confirm: []; cancel: [] }>()
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-end justify-center sm:items-center p-4"
        role="dialog"
        aria-modal="true"
        aria-labelledby="activate-modal-title"
      >
        <div class="absolute inset-0 bg-black/40" @click="emit('cancel')" />
        <div class="relative bg-white rounded-2xl p-6 w-full max-w-sm shadow-xl">
          <h2 id="activate-modal-title" class="text-lg font-bold text-gray-800 mb-3">
            فعال‌سازی برنامه
          </h2>
          <p class="text-sm text-gray-600 mb-2">
            با فعال‌سازی این برنامه، برنامه فعال قبلی مشتری
            <strong>بایگانی خواهد شد</strong>.
          </p>
          <p class="text-sm text-gray-500 mb-6">آیا مطمئن هستید؟</p>
          <div class="flex gap-3">
            <button
              class="flex-1 bg-emerald-600 text-white rounded-xl py-3 font-medium hover:bg-emerald-700 transition-colors"
              @click="emit('confirm')"
            >
              بله، فعال کن
            </button>
            <button
              class="flex-1 bg-gray-100 text-gray-700 rounded-xl py-3 font-medium hover:bg-gray-200 transition-colors"
              @click="emit('cancel')"
            >
              انصراف
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
