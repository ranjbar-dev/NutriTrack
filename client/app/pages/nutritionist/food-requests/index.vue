<script setup lang="ts">
import { ref } from 'vue'
import FoodRequestReviewList from '~/components/nutritionist/FoodRequestReviewList.vue'
import FoodRequestDecisionSheet from '~/components/nutritionist/FoodRequestDecisionSheet.vue'
import { useFoodRequestApi } from '~/composables/useFoodRequestApi'
import type { FoodRequest } from '~/types/food-request'

definePageMeta({
  layout: 'nutritionist',
})

const api = useFoodRequestApi()
const loading = ref(false)
const acting = ref(false)
const error = ref('')
const requests = ref<FoodRequest[]>([])
const selected = ref<FoodRequest | null>(null)

async function loadRequests() {
  loading.value = true
  error.value = ''
  const { data, error: requestError } = await api.listPendingFoodRequests(1, 30)
  if (requestError.value) {
    error.value = 'دریافت لیست درخواست ها انجام نشد.'
    requests.value = []
  } else {
    requests.value = data.value?.data ?? []
  }
  loading.value = false
}

async function approve(requestId: string) {
  const request = requests.value.find((item) => item.id === requestId)
  if (!request) return
  acting.value = true
  await api.approveFoodRequest(requestId, {
    name: request.food_name,
    unit: 'گرم',
  })
  acting.value = false
  selected.value = null
  await loadRequests()
}

async function reject(payload: { requestId: string; reason: string }) {
  acting.value = true
  await api.rejectFoodRequest(payload.requestId, { reason: payload.reason })
  acting.value = false
  selected.value = null
  await loadRequests()
}

await loadRequests()
</script>

<template>
  <main class="page">
    <header>
      <h2>بررسی درخواست غذا</h2>
      <button type="button" @click="loadRequests">بازخوانی</button>
    </header>

    <FoodRequestReviewList :requests="requests" :loading="loading" :error="error" @review="selected = $event" />

    <FoodRequestDecisionSheet
      :visible="selected !== null"
      :request="selected"
      :loading="acting"
      @approve="approve"
      @reject="reject"
      @close="selected = null"
    />
  </main>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
}

header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

button {
  border: 1px solid #c8d2d8;
  border-radius: 8px;
  min-height: 36px;
  padding: 0 10px;
  background: #fff;
}
</style>
