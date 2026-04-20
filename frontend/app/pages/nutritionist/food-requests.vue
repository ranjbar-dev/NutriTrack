<script setup lang="ts">
definePageMeta({ layout: 'nutritionist' })

const foodRequestStore = useFoodRequestStore()
const router = useRouter()
const rejectId = ref<string | null>(null)
const rejectReason = ref('')
const processing = ref(false)

onMounted(() => foodRequestStore.fetchNutriRequests())

async function approve(id: string) {
  processing.value = true
  const fr = await foodRequestStore.approve(id)
  processing.value = false
  if (fr) {
    // Navigate to food creation pre-filled with food name
    router.push(`/nutritionist/foods/create?name=${encodeURIComponent(fr.food_name)}`)
  }
}

async function submitReject() {
  if (!rejectId.value) return
  processing.value = true
  await foodRequestStore.reject(rejectId.value, rejectReason.value.trim() || undefined)
  processing.value = false
  rejectId.value = null
  rejectReason.value = ''
}
</script>

<template>
  <div class="p-4 flex flex-col gap-4">
    <h1 class="text-xl font-bold text-gray-800">درخواست‌های غذایی مراجعین</h1>

    <div v-if="foodRequestStore.loading" class="text-center text-gray-400 text-sm">در حال بارگذاری...</div>
    <div v-else-if="foodRequestStore.requests.length === 0" class="bg-white rounded-xl p-6 shadow-sm text-center text-gray-500 text-sm">
      هیچ درخواست در انتظار بررسی وجود ندارد.
    </div>

    <div
      v-for="req in foodRequestStore.requests"
      :key="req.id"
      class="bg-white rounded-xl p-4 shadow-sm flex flex-col gap-2"
    >
      <div class="flex items-start justify-between">
        <div>
          <p class="font-semibold text-gray-800">{{ req.food_name }}</p>
          <p v-if="req.description" class="text-xs text-gray-500 mt-0.5">{{ req.description }}</p>
          <p v-if="req.client_name" class="text-xs text-gray-400 mt-0.5">مراجع: {{ req.client_name }}</p>
        </div>
        <span class="text-xs bg-yellow-100 text-yellow-700 rounded-full px-2 py-0.5">در انتظار بررسی</span>
      </div>

      <div v-if="rejectId === req.id" class="flex flex-col gap-2">
        <textarea
          v-model="rejectReason"
          rows="2"
          placeholder="دلیل رد (اختیاری)"
          class="rounded-lg border border-gray-200 p-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-red-400"
        />
        <div class="flex gap-2">
          <button
            :disabled="processing"
            class="flex-1 bg-red-500 text-white rounded-lg py-1.5 text-sm disabled:opacity-50"
            @click="submitReject"
          >
            رد کردن
          </button>
          <button
            class="flex-1 border border-gray-200 rounded-lg py-1.5 text-sm text-gray-600"
            @click="rejectId = null"
          >
            انصراف
          </button>
        </div>
      </div>

      <div v-else class="flex gap-2">
        <button
          :disabled="processing"
          class="flex-1 bg-emerald-500 text-white rounded-lg py-1.5 text-sm disabled:opacity-50"
          @click="approve(req.id)"
        >
          تأیید و ایجاد غذا
        </button>
        <button
          class="flex-1 border border-red-200 text-red-500 rounded-lg py-1.5 text-sm"
          @click="rejectId = req.id"
        >
          رد کردن
        </button>
      </div>
    </div>
  </div>
</template>
