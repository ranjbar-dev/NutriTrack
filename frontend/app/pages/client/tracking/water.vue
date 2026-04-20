<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'client' })

const planStore = useClientPlanStore()
const waterLogStore = useWaterLogStore()
const customAmount = ref('')
const formError = ref<string | null>(null)

onMounted(async () => {
  if (!planStore.activePlan) await planStore.fetchActivePlan()
  await waterLogStore.fetchToday()
})

async function addPreset(amount: number) {
  formError.value = null
  await waterLogStore.addWater(amount)
}

async function addCustom() {
  const amount = Number(customAmount.value)
  if (!amount || amount <= 0) {
    formError.value = 'مقدار آب باید بیشتر از صفر باشد'
    return
  }
  formError.value = null
  await waterLogStore.addWater(amount)
  customAmount.value = ''
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6">
    <h1 class="mb-4 text-lg font-bold text-gray-800 text-start">ثبت مصرف آب</h1>
    <WaterProgressBar :total-ml="waterLogStore.totalMl" :target-ml="planStore.activePlan?.daily_water_target_ml ?? 0" />

    <div class="mt-4 rounded-2xl bg-white p-4 shadow-sm">
      <div class="grid grid-cols-3 gap-2">
        <button type="button" class="rounded-xl bg-blue-50 py-3 text-sm font-medium text-blue-700" @click="addPreset(200)">۲۰۰ ml</button>
        <button type="button" class="rounded-xl bg-blue-50 py-3 text-sm font-medium text-blue-700" @click="addPreset(250)">۲۵۰ ml</button>
        <button type="button" class="rounded-xl bg-blue-50 py-3 text-sm font-medium text-blue-700" @click="addPreset(500)">۵۰۰ ml</button>
      </div>
      <div class="mt-3 flex gap-2">
        <input v-model="customAmount" type="number" min="1" class="flex-1 rounded-xl border p-3 text-start" placeholder="مقدار دلخواه" />
        <button type="button" class="rounded-xl bg-emerald-500 px-4 py-3 text-sm font-medium text-white" @click="addCustom">ثبت</button>
      </div>
      <p v-if="formError" class="mt-2 text-sm text-rose-600 text-start">{{ formError }}</p>
    </div>

    <div class="mt-4 rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">نوشیدن‌های امروز</h2>
      <div v-if="waterLogStore.logs.length === 0" class="py-6 text-center text-sm text-gray-400">هنوز آبی ثبت نشده است</div>
      <div v-else class="space-y-2">
        <div v-for="item in waterLogStore.logs" :key="item.id" class="flex items-center justify-between rounded-xl bg-gray-50 p-3">
          <span class="text-sm text-gray-800">{{ item.amount_ml }} میلی‌لیتر</span>
          <span class="text-xs text-gray-400">{{ item.logged_time || item.created_at.slice(11, 16) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
