<script setup lang="ts">
definePageMeta({ layout: 'nutritionist' })

const clientStore = useClientManagementStore()
const router = useRouter()

const searchQuery = ref('')
const sortBy = ref('created_at')
const activeFilter = ref<'all' | 'active' | 'inactive'>('all')

const activeParam = computed(() => {
  if (activeFilter.value === 'active') return true
  if (activeFilter.value === 'inactive') return false
  return null
})

async function load(page = 1) {
  await clientStore.fetchClients({
    q: searchQuery.value,
    sort: sortBy.value,
    active: activeParam.value,
    page,
  })
}

onMounted(() => load())

let searchTimeout: ReturnType<typeof setTimeout>
watch(searchQuery, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => load(1), 400)
})

watch([sortBy, activeFilter], () => load(1))
</script>

<template>
  <div class="p-4 flex flex-col gap-4">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-bold text-gray-800">مراجعین</h1>
      <NuxtLink
        to="/nutritionist/clients/register"
        class="bg-blue-500 text-white rounded-lg px-4 py-2 text-sm"
      >
        + ثبت مراجع
      </NuxtLink>
    </div>

    <!-- Filters -->
    <div class="flex flex-col gap-2">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="جستجو بر اساس نام یا موبایل..."
        class="w-full rounded-lg border border-gray-200 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
      />
      <div class="flex gap-2">
        <select
          v-model="sortBy"
          class="flex-1 rounded-lg border border-gray-200 p-2 text-sm focus:outline-none"
        >
          <option value="created_at">مرتب‌سازی: تاریخ ثبت</option>
          <option value="name">مرتب‌سازی: نام</option>
        </select>
        <select
          v-model="activeFilter"
          class="flex-1 rounded-lg border border-gray-200 p-2 text-sm focus:outline-none"
        >
          <option value="all">همه</option>
          <option value="active">فعال</option>
          <option value="inactive">غیرفعال</option>
        </select>
      </div>
    </div>

    <div v-if="clientStore.loading" class="text-center text-gray-400 text-sm py-8">در حال بارگذاری...</div>
    <div v-else-if="clientStore.clients.length === 0" class="bg-white rounded-xl p-6 shadow-sm text-center text-gray-500 text-sm">
      هیچ مراجعی یافت نشد.
    </div>

    <div
      v-for="client in clientStore.clients"
      :key="client.id"
      class="bg-white rounded-xl p-4 shadow-sm flex items-center justify-between cursor-pointer hover:shadow-md transition"
      @click="router.push(`/nutritionist/clients/${client.id}`)"
    >
      <div>
        <p class="font-semibold text-gray-800">{{ client.full_name }}</p>
        <p v-if="client.mobile" class="text-xs text-gray-400 mt-0.5">{{ client.mobile }}</p>
      </div>
      <span
        class="text-xs rounded-full px-2 py-0.5"
        :class="client.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'"
      >
        {{ client.is_active ? 'فعال' : 'غیرفعال' }}
      </span>
    </div>

    <!-- Pagination -->
    <div v-if="clientStore.total > 20" class="flex justify-center gap-2">
      <button
        v-for="page in Math.ceil(clientStore.total / 20)"
        :key="page"
        class="w-8 h-8 rounded-lg text-sm"
        :class="page === clientStore.currentPage ? 'bg-blue-500 text-white' : 'bg-white border border-gray-200 text-gray-700'"
        @click="load(page)"
      >
        {{ page }}
      </button>
    </div>
  </div>
</template>

