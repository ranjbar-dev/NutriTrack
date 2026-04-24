<script setup lang="ts">
import type { FoodItem } from '~/types/catalogue'

defineProps<{
  items: FoodItem[]
  loading?: boolean
}>()

const emit = defineEmits<{
  delete: [item: FoodItem]
}>()
</script>

<template>
  <section class="list">
    <p v-if="loading" class="placeholder">در حال دریافت غذاها...</p>
    <p v-else-if="!items.length" class="placeholder">غذایی پیدا نشد.</p>

    <article v-for="item in items" v-else :key="item.id" class="card">
      <div>
        <h3>{{ item.name }}</h3>
        <p>{{ item.amount }} {{ item.unit }} - {{ item.calories }} kcal</p>
      </div>
      <button type="button" @click="emit('delete', item)">حذف</button>
    </article>
  </section>
</template>

<style scoped>
.list {
  display: grid;
  gap: 10px;
}

.card,
.placeholder {
  border: 1px solid #d3dce2;
  border-radius: 12px;
  background: #fff;
  padding: 12px;
}

.card {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.card h3,
.card p,
.placeholder {
  margin: 0;
}

.card p {
  color: #53606a;
  margin-top: 4px;
}

button {
  border: 1px solid #e4b5b5;
  color: #8b2121;
  background: #fff;
  border-radius: 10px;
  padding: 8px 12px;
}
</style>