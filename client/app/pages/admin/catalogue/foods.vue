<script setup lang="ts">
import { ref } from 'vue'
import AdminCatalogueFoodList from '~/components/admin/AdminCatalogueFoodList.vue'
import AdminCatalogueSearchHeader from '~/components/admin/AdminCatalogueSearchHeader.vue'
import AdminDangerConfirmSheet from '~/components/admin/AdminDangerConfirmSheet.vue'
import { useAdminCatalogueApi } from '~/composables/useAdminCatalogueApi'
import type { FoodItem } from '~/types/catalogue'

definePageMeta({
  layout: 'admin',
})

const api = useAdminCatalogueApi()
const foods = ref<FoodItem[]>([])
const loading = ref(true)
const error = ref('')
const selectedFood = ref<FoodItem | null>(null)
const deleting = ref(false)
const activeQuery = ref('')

async function refreshFoods(query = '') {
  activeQuery.value = query
  loading.value = true
  error.value = ''

  const { data, error: requestError } = await api.searchAdminFoods({ q: query, page: 1, page_size: 20 })
  foods.value = data.value?.data ?? []

  if (requestError.value) {
    error.value = 'دریافت فهرست غذاها انجام نشد'
  }

  loading.value = false
}

async function confirmDelete() {
  if (!selectedFood.value) {
    return
  }

  deleting.value = true
  try {
    await api.forceDeleteFood(selectedFood.value.id)
    selectedFood.value = null
    await refreshFoods(activeQuery.value)
  } catch {
    error.value = 'حذف غذا انجام نشد'
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  refreshFoods(activeQuery.value)
})
</script>

<template>
  <section class="catalogue-page">
    <header>
      <h1>راهبری غذاها</h1>
      <p>مدیریت فهرست غذاهای عمومی با دسترسی ویژه ادمین.</p>
    </header>

    <AdminCatalogueSearchHeader title="جستجوی غذا" @search="refreshFoods" />
    <p v-if="error" class="error">{{ error }}</p>
    <AdminCatalogueFoodList :items="foods" :loading="loading" @delete="selectedFood = $event" />

    <AdminDangerConfirmSheet
      :visible="Boolean(selectedFood)"
      :loading="deleting"
      title="حذف غذا"
      description="این حذف غیرقابل بازگشت است و باید با تایید ادمین انجام شود."
      @cancel="selectedFood = null"
      @confirm="confirmDelete"
    />
  </section>
</template>

<style scoped>
.catalogue-page {
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