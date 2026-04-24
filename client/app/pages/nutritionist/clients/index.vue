<script setup lang="ts">
import { computed, ref } from 'vue'
import ClientRosterFilters from '~/components/nutritionist/ClientRosterFilters.vue'
import ClientRosterList from '~/components/nutritionist/ClientRosterList.vue'
import type { NutritionistClient, NutritionistClientListQuery } from '~/types/nutritionist-workspace'
import { useNutritionistClientApi } from '~/composables/useNutritionistClientApi'

definePageMeta({
  layout: 'nutritionist',
})

const api = useNutritionistClientApi()

const listState = ref<{
  loading: boolean
  error: string
  clients: NutritionistClient[]
}>({
  loading: true,
  error: '',
  clients: [],
})

const queryState = ref<Required<Pick<NutritionistClientListQuery, 'q' | 'page' | 'page_size' | 'sort'>> & { status: 'active' | 'inactive' | 'all' }>({
  q: '',
  page: 1,
  page_size: 20,
  sort: 'newest',
  status: 'all',
})

async function refreshRoster() {
  listState.value.loading = true
  listState.value.error = ''

  const { data, error } = await api.listClients({
    q: queryState.value.q,
    page: queryState.value.page,
    page_size: queryState.value.page_size,
    sort: queryState.value.sort,
    status: queryState.value.status === 'all' ? undefined : queryState.value.status,
  })

  listState.value.clients = data.value?.data ?? []
  if (error.value) {
    listState.value.error = 'خطا در دریافت لیست مراجعان'
  }
  listState.value.loading = false
}

function applyFilters(payload: { query: string; status: 'active' | 'inactive' | 'all'; sort: 'newest' | 'oldest' | 'name_asc' | 'name_desc' }) {
  queryState.value.q = payload.query
  queryState.value.status = payload.status
  queryState.value.sort = payload.sort
  queryState.value.page = 1
  refreshRoster()
}

const pageTitle = computed(() => 'لیست مراجعان')

onMounted(() => {
  refreshRoster()
})
</script>

<template>
  <section class="roster-page">
    <header class="header">
      <h1>{{ pageTitle }}</h1>
      <p>جستجو، فیلتر و ورود به پروفایل هر مرجع</p>
    </header>

    <ClientRosterFilters @apply="applyFilters" />

    <ClientRosterList
      :clients="listState.clients"
      :loading="listState.loading"
      :error="listState.error"
      @open="(id) => navigateTo(`/nutritionist/clients/${id}`)"
    />
  </section>
</template>

<style scoped>
.roster-page {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.header h1 {
  margin: 0;
  font-size: 1.15rem;
}

.header p {
  margin: 4px 0 0;
  color: #53606a;
}
</style>
