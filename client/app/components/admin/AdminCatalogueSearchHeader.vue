<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  title: string
}>()

const emit = defineEmits<{
  search: [query: string]
}>()

const query = ref('')

function submit() {
  emit('search', query.value.trim())
}
</script>

<template>
  <form class="search-header" @submit.prevent="submit">
    <label>
      <span>{{ title }}</span>
      <input v-model="query" type="search" :placeholder="title" />
    </label>
    <button type="submit">جستجو</button>
  </form>
</template>

<style scoped>
.search-header {
  display: flex;
  gap: 8px;
  align-items: end;
}

label {
  display: grid;
  gap: 6px;
  flex: 1;
}

input,
button {
  min-height: 40px;
  border-radius: 10px;
  border: 1px solid #c8d2d8;
  padding: 0 12px;
}
</style>