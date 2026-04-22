<script setup lang="ts">
defineProps<{
  offline: boolean
}>()
</script>

<template>
  <div v-if="offline" class="banner connectivity" role="status">
    اتصال اینترنت ضعیف یا قطع است. داده‌ها پس از اتصال دوباره همگام می‌شوند.
  </div>
</template>

<style scoped>
.banner {
  margin: var(--space-3) var(--space-4) 0;
  padding: var(--space-3);
  border-radius: var(--radius-md);
  background: #fff4e5;
  color: #7a4a07;
  border: 1px solid #f0c48f;
}
</style>
