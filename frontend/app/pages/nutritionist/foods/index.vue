<script setup lang="ts">
import { toPersianDigits } from '~/utils/persian-digits'

definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist'],
})

const foodStore = useFoodStore()
const searchInput = ref('')
const searchLoading = ref(false)
const toastMessage = ref('')
const showToast = ref(false)
let searchTimeout: ReturnType<typeof setTimeout> | null = null

const isInitialLoading = computed(
  () => foodStore.loading && foodStore.foods.length === 0,
)

const emptyState = computed(() => {
  if (foodStore.loading || foodStore.total > 0) return null
  if (foodStore.searchQuery) return 'search'
  if (foodStore.selectedCategory) return 'category'
  return 'empty'
})

function triggerToast(message: string) {
  toastMessage.value = message
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}

watch(searchInput, (value) => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(async () => {
    searchLoading.value = true
    try {
      await foodStore.setSearch(value)
    }
    finally {
      searchLoading.value = false
    }
  }, 300)
})

onBeforeUnmount(() => {
  if (searchTimeout) clearTimeout(searchTimeout)
})

onMounted(() => {
  foodStore.fetchFoods(true)
})

function handleEdit(id: string) {
  navigateTo(`/nutritionist/foods/${id}`)
}

async function handleDelete(id: string) {
  if (!confirm('آیا از حذف این غذا اطمینان دارید؟')) return
  try {
    await foodStore.deleteFood(id)
  }
  catch {
    triggerToast('حذف غذا با خطا مواجه شد')
  }
}

async function clearSearch() {
  searchInput.value = ''
}

function setActiveFilter(showInactive: boolean) {
  if (foodStore.showInactive !== showInactive) {
    foodStore.toggleInactive()
  }
}
</script>

<template>
  <div class="p-4 space-y-4">
    <div
      v-if="showToast"
      class="fixed top-4 start-1/2 -translate-x-1/2 bg-red-600 text-white px-4 py-2 rounded-lg shadow-lg z-50 text-sm"
    >
      {{ toastMessage }}
    </div>

    <header class="flex items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-bold text-gray-800">غذاها</h1>
        <p class="text-xs text-gray-500 mt-1">
          {{ toPersianDigits(foodStore.total) }} غذا
        </p>
      </div>
      <UiAppButton class="w-auto" size="sm" @click="navigateTo('/nutritionist/foods/new')">
        افزودن غذا
      </UiAppButton>
    </header>

    <div class="relative">
      <UiAppInput
        v-model="searchInput"
        placeholder="جستجوی غذا..."
        :disabled="foodStore.loading && foodStore.foods.length === 0"
      />
      <UiLoadingSpinner
        v-if="searchLoading"
        size="sm"
        class="absolute top-1/2 end-3 -translate-y-1/2"
      />
    </div>

    <FoodCategoryPills
      :selected="foodStore.selectedCategory"
      @select="foodStore.setCategory"
    />

    <div class="flex items-center gap-2">
      <button
        type="button"
        class="px-3 py-1.5 rounded-full text-sm transition-colors"
        :class="foodStore.showInactive ? 'bg-gray-100 text-gray-600' : 'bg-emerald-600 text-white'"
        @click="setActiveFilter(false)"
      >
        فعال
      </button>
      <button
        type="button"
        class="px-3 py-1.5 rounded-full text-sm transition-colors"
        :class="foodStore.showInactive ? 'bg-emerald-600 text-white' : 'bg-gray-100 text-gray-600'"
        @click="setActiveFilter(true)"
      >
        غیرفعال
      </button>
    </div>

    <div v-if="isInitialLoading" class="py-16">
      <UiLoadingSpinner size="lg" />
    </div>

    <div v-else-if="emptyState === 'empty'" class="text-center py-12 space-y-3">
      <p class="text-gray-600">هیچ غذایی ثبت نشده</p>
      <UiAppButton class="w-auto" size="sm" @click="navigateTo('/nutritionist/foods/new')">
        افزودن غذا
      </UiAppButton>
    </div>

    <div v-else-if="emptyState === 'search'" class="text-center py-12 space-y-3">
      <p class="text-gray-600">نتیجه‌ای یافت نشد</p>
      <UiAppButton class="w-auto" size="sm" variant="secondary" @click="clearSearch">
        پاک کردن جستجو
      </UiAppButton>
    </div>

    <div v-else-if="emptyState === 'category'" class="text-center py-12">
      <p class="text-gray-600">در این دسته غذایی وجود ندارد</p>
    </div>

    <div v-else class="space-y-3">
      <FoodFoodCard
        v-for="food in foodStore.foods"
        :key="food.id"
        :food="food"
        @edit="handleEdit"
        @delete="handleDelete"
      />

      <div v-if="foodStore.hasMore" class="pt-4">
        <UiAppButton
          variant="secondary"
          :loading="foodStore.loading"
          @click="foodStore.loadMore"
        >
          بارگذاری بیشتر
        </UiAppButton>
      </div>
    </div>
  </div>
</template>
