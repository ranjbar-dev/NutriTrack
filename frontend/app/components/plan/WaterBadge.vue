<script setup lang="ts">
const props = defineProps<{ ml: number }>()

const liters = computed(() => (props.ml / 1000).toFixed(1))
</script>

<template>
  <div class="flex items-center gap-2 bg-blue-50 rounded-xl px-3 py-2 text-sm">
    <span class="text-blue-500">💧</span>
    <span class="text-blue-700 font-medium">{{ liters }} لیتر آب روزانه</span>
  </div>
</template>
