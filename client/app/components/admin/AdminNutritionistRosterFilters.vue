<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  apply: [payload: { query: string }]
}>()

const query = ref('')

function submit() {
  emit('apply', { query: query.value.trim() })
}
</script>

<template>
  <form class="filters" @submit.prevent="submit">
    <input v-model="query" type="search" placeholder="جستجو بر اساس نام یا ایمیل" />
    <button type="submit">اعمال</button>
  </form>
</template>

<style scoped>
.filters {
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

button {
  background: #fff;
}
</style>