<script setup lang="ts">
const props = defineProps<{
  totalMl: number
  targetMl: number
}>()

const percent = computed(() => {
  if (!props.targetMl) return 0
  return Math.min(100, Math.round((props.totalMl / props.targetMl) * 100))
})
</script>

<template>
  <div class="rounded-2xl bg-blue-50 p-4">
    <div class="mb-2 flex items-center justify-between gap-3">
      <p class="text-sm font-semibold text-blue-800 text-start">مصرف آب امروز</p>
      <p class="text-xs text-blue-700 text-end">
        <template v-if="targetMl">{{ totalMl }} از {{ targetMl }} میلی‌لیتر</template>
        <template v-else>هدف تنظیم نشده</template>
      </p>
    </div>
    <div class="h-3 rounded-full bg-blue-100 overflow-hidden">
      <div class="h-full rounded-full bg-blue-500 transition-all" :style="{ width: `${percent}%` }" />
    </div>
    <p class="mt-2 text-xs text-blue-700 text-start">{{ percent }}٪ از هدف روزانه</p>
  </div>
</template>
