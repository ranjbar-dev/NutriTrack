<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { AdminNutritionist, AdminUpdateNutritionistRequest } from '~/types/admin'

const props = defineProps<{
  nutritionist: AdminNutritionist | null
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: AdminUpdateNutritionistRequest]
}>()

const form = reactive<AdminUpdateNutritionistRequest>({
  first_name: '',
  last_name: '',
  email: '',
  mobile: '',
})

watch(
  () => props.nutritionist,
  (nutritionist) => {
    form.first_name = nutritionist?.first_name ?? ''
    form.last_name = nutritionist?.last_name ?? ''
    form.email = nutritionist?.email ?? ''
    form.mobile = nutritionist?.mobile ?? ''
  },
  { immediate: true },
)

function handleSubmit() {
  const payload: AdminUpdateNutritionistRequest = {}

  const firstName = form.first_name?.trim()
  const lastName = form.last_name?.trim()
  const email = form.email?.trim()
  const mobile = form.mobile?.trim()

  if (firstName) {
    payload.first_name = firstName
  }
  if (lastName) {
    payload.last_name = lastName
  }
  if (email) {
    payload.email = email
  }
  if (mobile) {
    payload.mobile = mobile
  }

  emit('submit', payload)
}
</script>

<template>
  <form class="form" @submit.prevent="handleSubmit">
    <h2>ویرایش اطلاعات</h2>
    <input v-model="form.first_name" type="text" placeholder="نام" />
    <input v-model="form.last_name" type="text" placeholder="نام خانوادگی" />
    <input v-model="form.email" type="email" placeholder="ایمیل" />
    <input v-model="form.mobile" type="tel" placeholder="شماره موبایل" />
    <button type="submit" :disabled="loading">ذخیره تغییرات</button>
  </form>
</template>

<style scoped>
.form {
  border: 1px solid #d3dce2;
  border-radius: 14px;
  background: #fff;
  padding: 14px;
  display: grid;
  gap: 8px;
}

h2 {
  margin: 0;
  font-size: 1rem;
}

input,
button {
  min-height: 40px;
  border-radius: 10px;
  border: 1px solid #c8d2d8;
  padding: 0 12px;
}

button {
  background: #173042;
  border-color: #173042;
  color: #fff;
}
</style>