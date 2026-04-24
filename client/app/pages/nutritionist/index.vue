<script setup lang="ts">
definePageMeta({
  layout: 'nutritionist'
})
</script>

<template>
  <section class="entry">
    <h2>فضای کار متخصص تغذیه</h2>
    <p>مدیریت مراجعان، مرور وضعیت برنامه، و نگارش برنامه غذایی از اینجا شروع می شود.</p>
    <NuxtLink class="entry-link" to="/nutritionist/clients">ورود به لیست مراجعان</NuxtLink>
    <NuxtLink class="entry-link secondary" to="/nutritionist/food-requests">درخواست های غذا</NuxtLink>
  </section>
</template>

<style scoped>
.entry {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  border: 1px solid #d4dce0;
  background: #fff;
}

h2 {
  margin: 0;
}

p {
  margin: 0;
  color: #53606a;
}

.entry-link {
  display: block;
  text-align: center;
  text-decoration: none;
  color: #fff;
  background: #0f6b7a;
  border-radius: 10px;
  padding: 10px;
  font-weight: 600;
}

.entry-link.secondary {
  background: #4f6870;
}
</style>
