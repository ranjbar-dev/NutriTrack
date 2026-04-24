<script setup lang="ts">
import type { FoodRequest } from '~/types/food-request'

defineProps<{
  requests: FoodRequest[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  review: [request: FoodRequest]
}>()
</script>

<template>
  <section class="card">
    <h3>درخواست های در انتظار</h3>
    <p v-if="loading">در حال دریافت...</p>
    <p v-else-if="error">{{ error }}</p>
    <ul v-else>
      <li v-for="request in requests" :key="request.id" class="row">
        <div>
          <strong>{{ request.food_name }}</strong>
          <p>وضعیت: {{ request.status }}</p>
        </div>
        <button type="button" @click="emit('review', request)">بررسی</button>
      </li>
      <li v-if="requests.length === 0">درخواستی برای بررسی وجود ندارد.</li>
    </ul>
  </section>
</template>

<style scoped>
.card {
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 12px;
}

ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.row {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  border: 1px solid #dde5ea;
  border-radius: 10px;
  padding: 8px;
}

p {
  margin: 0;
}

button {
  border: 1px solid #c8d2d8;
  border-radius: 8px;
  min-height: 34px;
  padding: 0 10px;
  background: #fff;
}
</style>
