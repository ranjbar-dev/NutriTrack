<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { AdminCreateNutritionistRequest } from '~/types/admin'

const props = defineProps<{
  visible: boolean
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  close: []
  submit: [payload: AdminCreateNutritionistRequest]
}>()

const form = reactive<AdminCreateNutritionistRequest>({
  email: '',
  password: '',
  first_name: '',
  last_name: '',
  mobile: '',
})

watch(
  () => props.visible,
  (visible) => {
    if (!visible) {
      form.email = ''
      form.password = ''
      form.first_name = ''
      form.last_name = ''
      form.mobile = ''
    }
  },
)

function handleSubmit() {
  if (
    !form.email.trim() ||
    !form.password.trim() ||
    !form.first_name.trim() ||
    !form.last_name.trim() ||
    !form.mobile.trim()
  ) {
    return
  }

  emit('submit', {
    email: form.email.trim(),
    password: form.password.trim(),
    first_name: form.first_name.trim(),
    last_name: form.last_name.trim(),
    mobile: form.mobile.trim(),
  })
}
</script>

<template>
  <section v-if="visible" class="sheet">
    <header>
      <h3>ایجاد متخصص تغذیه</h3>
      <button type="button" @click="emit('close')">بستن</button>
    </header>

    <p v-if="error" class="error">{{ error }}</p>

    <div class="fields">
      <input v-model="form.first_name" type="text" placeholder="نام" />
      <input v-model="form.last_name" type="text" placeholder="نام خانوادگی" />
      <input v-model="form.email" type="email" placeholder="ایمیل" />
      <input v-model="form.mobile" type="tel" placeholder="شماره موبایل" />
      <input v-model="form.password" type="password" placeholder="رمز عبور" />
    </div>

    <button type="button" class="submit" :disabled="loading" @click="handleSubmit">
      ثبت حساب
    </button>
  </section>
</template>

<style scoped>
.sheet {
  border: 1px solid #d3dce2;
  border-radius: 14px;
  background: #fff;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

header h3,
.error {
  margin: 0;
}

.fields {
  display: grid;
  gap: 8px;
}

input,
button {
  min-height: 40px;
  border-radius: 10px;
  border: 1px solid #c8d2d8;
  padding: 0 12px;
}

.submit {
  background: #173042;
  border-color: #173042;
  color: #fff;
}

.error {
  color: #8b2121;
}
</style>