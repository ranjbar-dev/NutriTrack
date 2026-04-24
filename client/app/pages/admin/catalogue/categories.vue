<script setup lang="ts">
import { ref } from 'vue'
import AdminFoodCategoryManager from '~/components/admin/AdminFoodCategoryManager.vue'
import { useAdminCatalogueApi } from '~/composables/useAdminCatalogueApi'
import type { FoodCategory } from '~/types/catalogue'

definePageMeta({
  layout: 'admin',
})

const api = useAdminCatalogueApi()
const categories = ref<FoodCategory[]>([])
const loading = ref(true)
const error = ref('')

async function refreshCategories() {
  loading.value = true
  error.value = ''

  const { data, error: requestError } = await api.listFoodCategories()
  categories.value = data.value?.data ?? []
  if (requestError.value) {
    error.value = 'دریافت دسته بندی ها انجام نشد'
  }
  loading.value = false
}

async function handleCreate(name: string) {
  try {
    await api.createFoodCategory({ name })
    await refreshCategories()
  } catch {
    error.value = 'ایجاد دسته بندی انجام نشد'
  }
}

async function handleDelete(id: string) {
  try {
    await api.deleteFoodCategory(id)
    await refreshCategories()
  } catch {
    error.value = 'حذف دسته بندی انجام نشد'
  }
}

onMounted(() => {
  refreshCategories()
})
</script>

<template>
  <section class="categories-page">
    <header>
      <h1>مدیریت دسته بندی غذا</h1>
      <p>ایجاد و حذف دسته بندی ها بدون نمایش گزارش های خارج از قرارداد API.</p>
    </header>

    <p v-if="error" class="error">{{ error }}</p>

    <AdminFoodCategoryManager
      :items="categories"
      :loading="loading"
      @create="handleCreate"
      @delete="handleDelete"
    />
  </section>
</template>

<style scoped>
.categories-page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

header h1 {
  margin: 0;
  font-size: 1.15rem;
}

header p,
.error {
  margin: 4px 0 0;
}

.error {
  color: #8b2121;
}
</style>