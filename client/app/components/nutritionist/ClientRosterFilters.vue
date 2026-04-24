<script setup lang="ts">
import { ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    initialQuery?: string
    initialStatus?: 'active' | 'inactive' | 'all'
    initialSort?: 'newest' | 'oldest' | 'name_asc' | 'name_desc'
  }>(),
  {
    initialQuery: '',
    initialStatus: 'all',
    initialSort: 'newest',
  }
)

const emit = defineEmits<{
  apply: [{ query: string; status: 'active' | 'inactive' | 'all'; sort: 'newest' | 'oldest' | 'name_asc' | 'name_desc' }]
}>()

const query = ref(props.initialQuery)
const status = ref(props.initialStatus)
const sort = ref(props.initialSort)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

watch(query, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  debounceTimer = setTimeout(() => {
    emit('apply', { query: query.value.trim(), status: status.value, sort: sort.value })
  }, 300)
})

function applyNow() {
  emit('apply', { query: query.value.trim(), status: status.value, sort: sort.value })
}
</script>

<template>
  <section class="filters">
    <label>
      جستجو
      <input v-model="query" type="text" placeholder="نام یا موبایل" />
    </label>

    <label>
      وضعیت
      <select v-model="status" @change="applyNow">
        <option value="all">همه</option>
        <option value="active">فعال</option>
        <option value="inactive">غیرفعال</option>
      </select>
    </label>

    <label>
      ترتیب
      <select v-model="sort" @change="applyNow">
        <option value="newest">جدیدترین</option>
        <option value="oldest">قدیمی ترین</option>
        <option value="name_asc">نام صعودی</option>
        <option value="name_desc">نام نزولی</option>
      </select>
    </label>
  </section>
</template>

<style scoped>
.filters {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
  padding: 12px;
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
}

label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.9rem;
}

input,
select {
  min-height: 40px;
  border-radius: 8px;
  border: 1px solid #ccd5db;
  padding: 6px 10px;
}
</style>
