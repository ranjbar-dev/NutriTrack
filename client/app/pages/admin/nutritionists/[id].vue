<script setup lang="ts">
import { computed, ref } from 'vue'
import AdminNutritionistClientReadonlyList from '~/components/admin/AdminNutritionistClientReadonlyList.vue'
import AdminNutritionistDetailCard from '~/components/admin/AdminNutritionistDetailCard.vue'
import AdminNutritionistEditForm from '~/components/admin/AdminNutritionistEditForm.vue'
import AdminNutritionistStatusConfirmSheet from '~/components/admin/AdminNutritionistStatusConfirmSheet.vue'
import { useAdminNutritionistApi } from '~/composables/useAdminNutritionistApi'
import type { AdminNutritionist, AdminUpdateNutritionistRequest } from '~/types/admin'
import type { NutritionistClient } from '~/types/nutritionist-workspace'

definePageMeta({
  layout: 'admin',
})

const route = useRoute()
const nutritionistId = computed(() => String(route.params.id || ''))
const api = useAdminNutritionistApi()

const profile = ref<AdminNutritionist | null>(null)
const clients = ref<NutritionistClient[]>([])
const loading = ref(true)
const saving = ref(false)
const statusUpdating = ref(false)
const error = ref('')
const success = ref('')
const statusSheetVisible = ref(false)

async function loadDetail() {
  if (!nutritionistId.value) {
    error.value = 'شناسه متخصص معتبر نیست'
    loading.value = false
    return
  }

  loading.value = true
  error.value = ''

  const [{ data: nutritionistData, error: nutritionistError }, { data: clientsData, error: clientsError }] = await Promise.all([
    api.getNutritionist(nutritionistId.value),
    api.listNutritionistClients(nutritionistId.value, { page: 1, page_size: 20 }),
  ])

  profile.value = nutritionistData.value ?? null
  clients.value = clientsData.value?.data ?? []

  if (nutritionistError.value || clientsError.value) {
    error.value = 'دریافت اطلاعات متخصص انجام نشد'
  }

  loading.value = false
}

async function handleUpdate(payload: AdminUpdateNutritionistRequest) {
  if (!nutritionistId.value) {
    return
  }

  saving.value = true
  success.value = ''
  error.value = ''

  try {
    await api.updateNutritionist(nutritionistId.value, payload)
    success.value = 'اطلاعات متخصص به روز شد'
    await loadDetail()
  } catch {
    error.value = 'به روز رسانی اطلاعات متخصص انجام نشد'
  } finally {
    saving.value = false
  }
}

async function handleStatusConfirm() {
  if (!nutritionistId.value || !profile.value) {
    return
  }

  statusUpdating.value = true
  success.value = ''
  error.value = ''

  try {
    await api.setNutritionistStatus(nutritionistId.value, { is_active: !profile.value.is_active })
    statusSheetVisible.value = false
    success.value = 'وضعیت حساب متخصص به روز شد'
    await loadDetail()
  } catch {
    error.value = 'تغییر وضعیت حساب متخصص انجام نشد'
  } finally {
    statusUpdating.value = false
  }
}

onMounted(() => {
  loadDetail()
})
</script>

<template>
  <section class="detail-page">
    <p v-if="error" class="notice error">{{ error }}</p>
    <p v-if="success" class="notice success">{{ success }}</p>

    <AdminNutritionistDetailCard :nutritionist="profile" :loading="loading" @change-status="statusSheetVisible = true" />

    <AdminNutritionistEditForm :nutritionist="profile" :loading="saving" @submit="handleUpdate" />

    <AdminNutritionistClientReadonlyList :items="clients" :loading="loading" />

    <AdminNutritionistStatusConfirmSheet
      :visible="statusSheetVisible"
      :nutritionist="profile"
      :loading="statusUpdating"
      @cancel="statusSheetVisible = false"
      @confirm="handleStatusConfirm"
    />
  </section>
</template>

<style scoped>
.detail-page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notice {
  margin: 0;
}

.notice.error {
  color: #8b2121;
}

.notice.success {
  color: #1b6a43;
}
</style>