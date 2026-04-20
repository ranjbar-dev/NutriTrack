<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'client' })
const trackingStore = useTrackingStore()

onMounted(() => trackingStore.fetchDailyDashboard())
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6">
    <DailyDashboard />
  </div>
</template>
