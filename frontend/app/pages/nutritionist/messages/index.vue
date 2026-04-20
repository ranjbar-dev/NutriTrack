<script setup lang="ts">
definePageMeta({ layout: 'nutritionist' })

const messageStore = useMessageStore()
</script>

<template>
  <div class="p-4">
    <h1 class="text-xl font-bold text-gray-800 mb-4">پیام‌ها</h1>
    <div class="bg-white rounded-xl p-6 shadow-sm text-center text-gray-500">
      برای شروع گفتگو، از صفحه پروفایل مراجع وارد چت شوید.
    </div>
  </div>
</template>
