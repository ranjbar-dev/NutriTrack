<script setup lang="ts">
definePageMeta({ layout: 'client' })

const foodRequestStore = useFoodRequestStore()
const showForm = ref(false)
const foodName = ref('')
const description = ref('')
const submitting = ref(false)
const formError = ref<string | null>(null)

onMounted(() => foodRequestStore.fetchClientRequests())

const statusLabel: Record<string, string> = {
  pending: 'در انتظار بررسی',
  approved: 'تأیید شده',
  rejected: 'رد شده',
}

const statusColor: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-700',
  approved: 'bg-green-100 text-green-700',
  rejected: 'bg-red-100 text-red-700',
}

async function submitRequest() {
  if (!foodName.value.trim()) return
  submitting.value = true
  formError.value = null
  const fr = await foodRequestStore.createRequest(
    foodName.value.trim(),
    description.value.trim() || undefined,
  )
  if (fr) {
    showForm.value = false
    foodName.value = ''
    description.value = ''
  } else {
    formError.value = foodRequestStore.error
  }
  submitting.value = false
}
</script>

<template>
  <div class="p-4 flex flex-col gap-4">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-bold text-gray-800">درخواست‌های غذایی</h1>
      <button
        class="bg-emerald-500 text-white rounded-lg px-4 py-2 text-sm"
        @click="showForm = !showForm"
      >
        {{ showForm ? 'انصراف' : '+ درخواست جدید' }}
      </button>
    </div>

    <div v-if="showForm" class="bg-white rounded-xl p-4 shadow-sm flex flex-col gap-3">
      <h2 class="font-semibold text-gray-700">درخواست اضافه شدن ماده غذایی</h2>
      <input
        v-model="foodName"
        type="text"
        placeholder="نام ماده غذایی"
        class="rounded-lg border border-gray-200 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-400"
      />
      <textarea
        v-model="description"
        rows="3"
        placeholder="توضیحات (اختیاری)"
        class="rounded-lg border border-gray-200 p-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-emerald-400"
      />
      <p v-if="formError" class="text-red-500 text-xs">{{ formError }}</p>
      <button
        :disabled="submitting || !foodName.trim()"
        class="bg-emerald-500 text-white rounded-lg py-2 text-sm disabled:opacity-50"
        @click="submitRequest"
      >
        ارسال درخواست
      </button>
    </div>

    <div v-if="foodRequestStore.loading" class="text-center text-gray-400 text-sm">در حال بارگذاری...</div>
    <div v-else-if="foodRequestStore.requests.length === 0" class="bg-white rounded-xl p-6 shadow-sm text-center text-gray-500 text-sm">
      هنوز درخواستی ثبت نشده است.
    </div>

    <div
      v-for="req in foodRequestStore.requests"
      :key="req.id"
      class="bg-white rounded-xl p-4 shadow-sm flex flex-col gap-1"
    >
      <div class="flex items-center justify-between">
        <span class="font-semibold text-gray-800">{{ req.food_name }}</span>
        <span class="text-xs rounded-full px-2 py-0.5" :class="statusColor[req.status]">
          {{ statusLabel[req.status] }}
        </span>
      </div>
      <p v-if="req.description" class="text-sm text-gray-500">{{ req.description }}</p>
      <p v-if="req.rejection_reason" class="text-xs text-red-500">دلیل رد: {{ req.rejection_reason }}</p>
    </div>
  </div>
</template>
