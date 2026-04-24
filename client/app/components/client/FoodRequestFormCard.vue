<script setup lang="ts">
import { reactive } from 'vue'
import type { SubmitFoodRequestRequest } from '~/types/food-request'

defineProps<{
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: SubmitFoodRequestRequest]
}>()

const form = reactive<SubmitFoodRequestRequest>({
  food_name: '',
})

function submit() {
  if (!form.food_name.trim()) {
    return
  }
  emit('submit', {
    food_name: form.food_name.trim(),
  })
}
</script>

<template>
  <section class="card">
    <h3>درخواست غذای جدید</h3>
    <label>
      نام غذا
      <input v-model="form.food_name" type="text" />
    </label>
    <button type="button" :disabled="loading || !form.food_name.trim()" @click="submit">
      {{ loading ? 'در حال ارسال...' : 'ارسال درخواست' }}
    </button>
  </section>
</template>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 12px;
}

label {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

input,
button {
  min-height: 40px;
  border-radius: 8px;
  border: 1px solid #c8d2d8;
  padding: 6px 10px;
}

button {
  border: none;
  background: #0f6b7a;
  color: #fff;
  font-weight: 700;
}
</style>
