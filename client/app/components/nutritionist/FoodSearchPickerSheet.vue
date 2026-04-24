<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FoodCategory, FoodItem } from '~/types/catalogue'
import { useCatalogueApi } from '~/composables/useCatalogueApi'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  close: []
  select: [food: FoodItem]
}>()

const api = useCatalogueApi()
const q = ref('')
const categoryId = ref('')
const foods = ref<FoodItem[]>([])
const categories = ref<FoodCategory[]>([])
const error = ref('')

async function loadCategories() {
  const { data } = await api.getFoodCategories()
  categories.value = data.value?.data ?? []
}

async function search() {
  error.value = ''
  const { data, error: searchError } = await api.searchFoods({
    q: q.value.trim(),
    category_id: categoryId.value || undefined,
    page: 1,
    page_size: 20,
  })
  foods.value = data.value?.data ?? []
  if (searchError.value) {
    error.value = 'جستجوی غذا با خطا مواجه شد.'
  }
}

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      loadCategories()
      search()
    }
  }
)
</script>

<template>
  <section v-if="visible" class="sheet">
    <header>
      <h4>انتخاب غذا</h4>
      <button type="button" @click="emit('close')">بستن</button>
    </header>

    <div class="controls">
      <input v-model="q" type="text" placeholder="جستجو غذا" @input="search" />
      <select v-model="categoryId" @change="search">
        <option value="">همه دسته ها</option>
        <option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option>
      </select>
    </div>

    <p v-if="error" class="error">{{ error }}</p>
    <ul class="results">
      <li v-for="food in foods" :key="food.id">
        <button type="button" @click="emit('select', food)">
          <strong>{{ food.name }}</strong>
          <span>{{ food.calories }} kcal / {{ food.amount }} {{ food.unit }}</span>
        </button>
      </li>
      <li v-if="foods.length === 0">نتیجه ای یافت نشد.</li>
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
