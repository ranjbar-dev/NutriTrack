<script setup lang="ts">
import { ref } from 'vue'
import AdminCatalogueMedicationList from '~/components/admin/AdminCatalogueMedicationList.vue'
import AdminCatalogueSearchHeader from '~/components/admin/AdminCatalogueSearchHeader.vue'
import AdminDangerConfirmSheet from '~/components/admin/AdminDangerConfirmSheet.vue'
import { useAdminCatalogueApi } from '~/composables/useAdminCatalogueApi'
import type { MedicationItem } from '~/types/catalogue'

definePageMeta({
  layout: 'admin',
})

const api = useAdminCatalogueApi()
const medications = ref<MedicationItem[]>([])
const loading = ref(true)
const error = ref('')
const selectedMedication = ref<MedicationItem | null>(null)
const deleting = ref(false)
const activeQuery = ref('')

async function refreshMedications(query = '') {
  activeQuery.value = query
  loading.value = true
  error.value = ''

  const { data, error: requestError } = await api.searchAdminMedications({
    q: query,
    page: 1,
    page_size: 20,
  })
  medications.value = data.value?.data ?? []

  if (requestError.value) {
    error.value = 'دریافت فهرست داروها انجام نشد'
  }

  loading.value = false
}

async function confirmDelete() {
  if (!selectedMedication.value) {
    return
  }

  deleting.value = true
  try {
    await api.forceDeleteMedication(selectedMedication.value.id)
    selectedMedication.value = null
    await refreshMedications(activeQuery.value)
  } catch {
    error.value = 'حذف دارو انجام نشد'
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  refreshMedications(activeQuery.value)
})
</script>

<template>
  <section class="catalogue-page">
    <header>
      <h1>راهبری داروها</h1>
      <p>کنترل داروهای عمومی با مسیرهای ویژه ادمین.</p>
    </header>

    <AdminCatalogueSearchHeader title="جستجوی دارو" @search="refreshMedications" />
    <p v-if="error" class="error">{{ error }}</p>
    <AdminCatalogueMedicationList
      :items="medications"
      :loading="loading"
      @delete="selectedMedication = $event"
    />

    <AdminDangerConfirmSheet
      :visible="Boolean(selectedMedication)"
      :loading="deleting"
      title="حذف دارو"
      description="حذف دارو فقط بعد از تایید صریح ادمین انجام می شود."
      @cancel="selectedMedication = null"
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