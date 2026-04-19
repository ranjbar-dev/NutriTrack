<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 gap-4">
    <h1 class="text-2xl font-bold text-red-600">دسترسی غیرمجاز</h1>
    <p class="text-gray-600">شما مجوز دسترسی به این صفحه را ندارید</p>
    <NuxtLink to="/" class="text-emerald-600 underline">بازگشت به صفحه اصلی</NuxtLink>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'auth', middleware: [] })
</script>
