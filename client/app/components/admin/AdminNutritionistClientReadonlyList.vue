<script setup lang="ts">
import type { NutritionistClient } from '~/types/nutritionist-workspace'

defineProps<{
  items: NutritionistClient[]
  loading?: boolean
}>()
</script>

<template>
  <section class="client-list">
    <header>
      <h2>مراجعان مرتبط</h2>
      <p>این بخش فقط برای مشاهده است و هیچ عملیات مدیریتی روی مراجعان ارائه نمی کند.</p>
    </header>

    <p v-if="loading" class="placeholder">در حال دریافت لیست مراجعان...</p>
    <p v-else-if="!items.length" class="placeholder">مراجعی برای این متخصص ثبت نشده است.</p>

    <article v-for="item in items" v-else :key="item.id" class="card">
      <strong>{{ item.full_name }}</strong>
      <span>{{ item.mobile }}</span>
    </article>
  </section>
</template>

<style scoped>
.client-list {
  border: 1px solid #d3dce2;
  border-radius: 14px;
  background: #fff;
  padding: 14px;
  display: grid;
  gap: 10px;
}

header h2,
header p,
.placeholder {
  margin: 0;
}

header p {
  margin-top: 4px;
  color: #53606a;
}

.card {
  border: 1px solid #e4eaee;
  border-radius: 12px;
  padding: 10px 12px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

span {
  color: #53606a;
}
</style>