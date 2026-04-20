<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'client' })

const planStore = useClientPlanStore()
const medicationStore = useMedicationLogStore()
const manualOpen = ref(false)
const manualForm = reactive({ medication_name: '', dosage: '', taken_at: '', notes: '' })
const formError = ref<string | null>(null)

onMounted(async () => {
  if (!planStore.activePlan) await planStore.fetchActivePlan()
  await medicationStore.fetchToday()
})

async function markTaken(item: { prescribedMedicationId: string; medicationName: string; dosage?: string | null; time: string }) {
  await medicationStore.logMedication({
    prescribed_medication_id: item.prescribedMedicationId,
    medication_name: item.medicationName,
    dosage: item.dosage || undefined,
    taken_at: item.time,
    is_self_reported: false,
  })
}

async function submitManual() {
  if (!manualForm.medication_name.trim() || !manualForm.taken_at) {
    formError.value = 'نام دارو و زمان مصرف الزامی است'
    return
  }
  formError.value = null
  await medicationStore.logMedication({
    medication_name: manualForm.medication_name.trim(),
    ...(manualForm.dosage ? { dosage: manualForm.dosage } : {}),
    taken_at: manualForm.taken_at,
    ...(manualForm.notes ? { notes: manualForm.notes } : {}),
    is_self_reported: true,
  })
  Object.assign(manualForm, { medication_name: '', dosage: '', taken_at: '', notes: '' })
  manualOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6">
    <h1 class="mb-4 text-lg font-bold text-gray-800 text-start">ثبت دارو</h1>

    <div class="rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">داروهای تجویزی</h2>
      <div v-if="!planStore.activePlan" class="py-4 text-center text-sm text-gray-400">برنامه غذایی فعالی ندارید</div>
      <div v-else-if="medicationStore.checklistItems.length === 0" class="py-4 text-center text-sm text-gray-400">دارویی تجویز نشده است</div>
      <div v-else class="space-y-2">
        <MedicationChecklistItem
          v-for="item in medicationStore.checklistItems"
          :key="`${item.prescribedMedicationId}-${item.time}`"
          :medication-name="item.medicationName"
          :dosage="item.dosage"
          :time="item.time"
          :is-taken="item.isTaken"
          @mark="markTaken(item)"
        />
      </div>
    </div>

    <div class="mt-4 rounded-2xl bg-white p-4 shadow-sm">
      <button type="button" class="text-sm font-medium text-emerald-700" @click="manualOpen = !manualOpen">
        {{ manualOpen ? 'بستن ثبت دستی' : 'افزودن داروی دستی' }}
      </button>
      <div v-if="manualOpen" class="mt-3 space-y-3">
        <input v-model="manualForm.medication_name" type="text" class="w-full rounded-xl border p-3 text-start" placeholder="نام دارو" />
        <input v-model="manualForm.dosage" type="text" class="w-full rounded-xl border p-3 text-start" placeholder="مقدار مصرف (اختیاری)" />
        <input v-model="manualForm.taken_at" type="time" class="w-full rounded-xl border p-3 text-start" />
        <textarea v-model="manualForm.notes" rows="2" class="w-full rounded-xl border p-3 text-start resize-none" placeholder="یادداشت (اختیاری)" />
        <p v-if="formError" class="text-sm text-rose-600 text-start">{{ formError }}</p>
        <button type="button" class="w-full rounded-xl bg-emerald-500 py-3 font-medium text-white" @click="submitManual">ثبت دارو</button>
      </div>
    </div>

    <div class="mt-4 rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">ثبت‌های امروز</h2>
      <div v-if="medicationStore.todayLogs.length === 0" class="py-6 text-center text-sm text-gray-400">هنوز دارویی ثبت نشده است</div>
      <div v-else class="space-y-2">
        <div v-for="log in medicationStore.todayLogs" :key="log.id" class="rounded-xl bg-gray-50 p-3">
          <div class="flex items-center justify-between gap-2">
            <p class="text-sm font-medium text-gray-800">{{ log.medication_name }}</p>
            <span class="text-xs text-gray-400">{{ log.taken_at }}</span>
          </div>
          <p class="mt-1 text-xs text-gray-500">
            <span v-if="log.dosage">{{ log.dosage }}</span>
            <span v-if="log.is_self_reported" class="ms-2 rounded-full bg-amber-100 px-2 py-0.5 text-amber-700">خوداظهاری</span>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
