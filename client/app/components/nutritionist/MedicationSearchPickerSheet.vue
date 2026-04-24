<script setup lang="ts">
import { ref, watch } from 'vue'
import type { MedicationItem } from '~/types/catalogue'
import { useCatalogueApi } from '~/composables/useCatalogueApi'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  close: []
  select: [medication: MedicationItem]
}>()

const api = useCatalogueApi()
const q = ref('')
const medications = ref<MedicationItem[]>([])
const error = ref('')

async function search() {
  error.value = ''
  const { data, error: searchError } = await api.searchMedications({
    q: q.value.trim(),
    page: 1,
    page_size: 20,
  })
  medications.value = data.value?.data ?? []
  if (searchError.value) {
    error.value = 'جستجوی دارو با خطا مواجه شد.'
  }
}

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      search()
    }
  }
)
</script>

<template>
  <section v-if="visible" class="sheet">
    <header>
      <h4>انتخاب دارو</h4>
      <button type="button" @click="emit('close')">بستن</button>
    </header>

    <div class="controls">
      <input v-model="q" type="text" placeholder="جستجو دارو" @input="search" />
    </div>

    <p v-if="error" class="error">{{ error }}</p>
    <ul class="results">
      <li v-for="item in medications" :key="item.id">
        <button type="button" @click="emit('select', item)">
          <strong>{{ item.name }}</strong>
          <span>{{ item.unit }}</span>
        </button>
      </li>
      <li v-if="medications.length === 0">نتیجه ای یافت نشد.</li>
    </ul>
  </section>
</template>

<style scoped>
.sheet {
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 10px;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.controls {
  display: grid;
  gap: 8px;
  margin: 8px 0;
}

.results {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.results button {
  width: 100%;
  border: 1px solid #d7dfe4;
  background: #f8fbfc;
  border-radius: 8px;
  padding: 8px;
  text-align: right;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.error {
  color: #8b2121;
  margin: 0;
}
</style>
