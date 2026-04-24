<script setup lang="ts">
import { ref } from 'vue'
import FoodRequestFormCard from '~/components/client/FoodRequestFormCard.vue'
import { useFoodRequestApi } from '~/composables/useFoodRequestApi'
import type { FoodRequest, SubmitFoodRequestRequest } from '~/types/food-request'

definePageMeta({
  layout: 'client',
})

const api = useFoodRequestApi()
const loading = ref(false)
const message = ref('')
const error = ref('')
const submitted = ref<FoodRequest[]>([])

async function submitRequest(payload: SubmitFoodRequestRequest) {
  loading.value = true
  error.value = ''
  message.value = ''
  try {
    const item = await api.submitFoodRequest(payload)
    submitted.value = [item, ...submitted.value]
    message.value = 'درخواست شما ثبت شد.'
  } catch {
    error.value = 'ارسال درخواست انجام نشد.'
  }
  loading.value = false
}
</script>

<template>
  <main class="page">
    <h2>درخواست غذا</h2>
    <FoodRequestFormCard :loading="loading" @submit="submitRequest" />

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="message" class="success">{{ message }}</p>

    <section class="card">
      <h3>درخواست های ثبت شده در این نشست</h3>
      <ul>
        <li v-for="item in submitted" :key="item.id">{{ item.food_name }} - {{ item.status }}</li>
        <li v-if="submitted.length === 0">هنوز درخواستی ثبت نشده است.</li>
      </ul>
    </section>
  </main>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
}

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
  gap: 6px;
}

.error {
  color: #8b2121;
}

.success {
  color: #176f2c;
}
</style>
