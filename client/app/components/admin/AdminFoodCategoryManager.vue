<script setup lang="ts">
import { ref } from 'vue'
import type { FoodCategory } from '~/types/catalogue'

defineProps<{
  items: FoodCategory[]
  loading?: boolean
}>()

const emit = defineEmits<{
  create: [name: string]
  delete: [id: string]
}>()

const name = ref('')

function submit() {
  if (!name.value.trim()) {
    return
  }

  emit('create', name.value.trim())
  name.value = ''
}
</script>

<template>
  <section class="manager">
    <form class="create-row" @submit.prevent="submit">
      <input v-model="name" type="text" placeholder="نام دسته بندی" />
      <button type="submit">ایجاد</button>
    </form>

    <p v-if="loading" class="placeholder">در حال دریافت دسته بندی ها...</p>
    <p v-else-if="!items.length" class="placeholder">دسته بندی ثبت نشده است.</p>

    <article v-for="item in items" v-else :key="item.id" class="card">
      <strong>{{ item.name }}</strong>
      <button type="button" @click="emit('delete', item.id)">حذف</button>
    </article>
  </section>
</template>

<style scoped>
.manager {
  display: grid;
  gap: 10px;
}

.create-row {
  display: flex;
  gap: 8px;
}

input,
button {
  min-height: 40px;
  border-radius: 10px;
  border: 1px solid #c8d2d8;
  padding: 0 12px;
}

input {
  flex: 1;
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

.placeholder {
  margin: 0;
}
</style>