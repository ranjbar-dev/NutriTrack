<script setup lang="ts">
import { computed, ref } from 'vue'
import AdminNutritionistCreateSheet from '~/components/admin/AdminNutritionistCreateSheet.vue'
import AdminNutritionistRosterFilters from '~/components/admin/AdminNutritionistRosterFilters.vue'
import AdminNutritionistRosterList from '~/components/admin/AdminNutritionistRosterList.vue'
import { useAdminNutritionistApi } from '~/composables/useAdminNutritionistApi'
import type {
  AdminCreateNutritionistRequest,
  AdminNutritionist,
  AdminNutritionistQuery,
} from '~/types/admin'

definePageMeta({
  layout: 'admin',
})

const api = useAdminNutritionistApi()

const listState = ref<{
  loading: boolean
  error: string
  nutritionists: AdminNutritionist[]
}>({
  loading: true,
  error: '',
  nutritionists: [],
})

const queryState = ref<Required<Pick<AdminNutritionistQuery, 'q' | 'page' | 'page_size'>>>({
  q: '',
  page: 1,
  page_size: 20,
})

const createSheetVisible = ref(false)
const createState = ref({
  loading: false,
  error: '',
  success: '',
})

async function refreshRoster() {
  listState.value.loading = true
  listState.value.error = ''

  const { data, error } = await api.listNutritionists({
    q: queryState.value.q,
    page: queryState.value.page,
    page_size: queryState.value.page_size,
  })

  listState.value.nutritionists = data.value?.data ?? []
  if (error.value) {
    listState.value.error = 'خطا در دریافت لیست متخصصان'
  }
  listState.value.loading = false
}

function applyFilters(payload: { query: string }) {
  queryState.value.q = payload.query
  queryState.value.page = 1
  refreshRoster()
}

async function handleCreate(payload: AdminCreateNutritionistRequest) {
  createState.value.loading = true
  createState.value.error = ''
  createState.value.success = ''

  try {
    await api.createNutritionist(payload)
    createState.value.success = 'حساب متخصص تغذیه با موفقیت ایجاد شد'
    createSheetVisible.value = false
    await refreshRoster()
  } catch {
    createState.value.error = 'ایجاد حساب متخصص تغذیه انجام نشد'
  } finally {
    createState.value.loading = false
  }
}

const pageTitle = computed(() => 'مدیریت متخصصان تغذیه')

onMounted(() => {
  refreshRoster()
})
</script>

<template>
  <section class="nutritionists-page">
    <header class="header">
      <div>
        <h1>{{ pageTitle }}</h1>
        <p>جستجو، مرور وضعیت و ایجاد حساب جدید برای متخصصان تغذیه.</p>
      </div>

      <button type="button" class="create-button" @click="createSheetVisible = true">
        ایجاد متخصص جدید
      </button>
    </header>

    <p v-if="createState.success" class="notice success">{{ createState.success }}</p>
    <p v-if="listState.error" class="notice error">{{ listState.error }}</p>

    <AdminNutritionistRosterFilters @apply="applyFilters" />

    <AdminNutritionistRosterList
      :nutritionists="listState.nutritionists"
      :loading="listState.loading"
      :error="listState.error"
      @open="(id) => navigateTo(`/admin/nutritionists/${id}`)"
    />

    <AdminNutritionistCreateSheet
      :visible="createSheetVisible"
      :loading="createState.loading"
      :error="createState.error"
      @close="createSheetVisible = false"
      @submit="handleCreate"
    />
  </section>
</template>

<style scoped>
.nutritionists-page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.header {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.header h1 {
  margin: 0;
  font-size: 1.2rem;
}

.header p {
  margin: 4px 0 0;
  color: #53606a;
}

.create-button {
  width: fit-content;
  border: 1px solid #173042;
  border-radius: 10px;
  background: #173042;
  color: #fff;
  padding: 10px 14px;
}

.notice {
  margin: 0;
}

.notice.success {
  color: #1b6a43;
}

.notice.error {
  color: #8b2121;
}
</style>