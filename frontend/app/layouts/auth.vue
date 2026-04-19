<template>
  <div class="min-h-screen bg-white flex flex-col">
    <!-- Header with app title -->
    <header class="py-8 text-center">
      <h1 class="text-2xl font-bold text-emerald-600">نوتری‌ترک</h1>
      <p class="mt-1 text-sm text-gray-500">مدیریت تغذیه هوشمند</p>
    </header>

    <!-- Content area -->
    <main class="flex-1 px-6">
      <slot />
    </main>
  </div>
</template>
