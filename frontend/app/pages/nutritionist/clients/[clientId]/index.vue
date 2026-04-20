<script setup lang="ts">
definePageMeta({ layout: 'nutritionist' })

const route = useRoute()
const router = useRouter()
const clientStore = useClientManagementStore()
const clientId = route.params.clientId as string

const editMode = ref(false)
const editDob = ref('')
const editHeight = ref<number | null>(null)
const saving = ref(false)
const errorMsg = ref<string | null>(null)

onMounted(() => clientStore.fetchClientProfile(clientId))

const profile = computed(() => clientStore.selectedClient)

function startEdit() {
  editDob.value = profile.value?.date_of_birth ?? ''
  editHeight.value = profile.value?.height_cm ?? null
  editMode.value = true
}

async function save() {
  saving.value = true
  errorMsg.value = null
  try {
    await clientStore.updateProfile(clientId, {
      date_of_birth: editDob.value || undefined,
      height_cm: editHeight.value ?? undefined,
    })
    editMode.value = false
  } catch (e: any) {
    errorMsg.value = e?.data?.error ?? 'خطا در ذخیره‌سازی'
  } finally {
    saving.value = false
  }
}

async function toggleActive() {
  if (!profile.value) return
  await clientStore.setActive(clientId, !profile.value.is_active)
}

const genderLabel: Record<string, string> = {
  male: 'مرد',
  female: 'زن',
  other: 'سایر',
}
</script>

<template>
  <div class="p-4 flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <button class="text-gray-400 hover:text-gray-700" @click="router.back()">←</button>
      <h1 class="text-xl font-bold text-gray-800">پروفایل مراجع</h1>
    </div>

    <div v-if="clientStore.loading" class="text-center text-gray-400 py-8">در حال بارگذاری...</div>

    <template v-else-if="profile">
      <div class="bg-white rounded-xl p-4 shadow-sm flex flex-col gap-3">
        <div class="flex items-center justify-between">
          <div>
            <p class="font-bold text-gray-800 text-lg">{{ profile.full_name }}</p>
            <p v-if="profile.mobile" class="text-sm text-gray-500">{{ profile.mobile }}</p>
          </div>
          <span
            class="text-xs rounded-full px-2 py-0.5 cursor-pointer"
            :class="profile.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'"
            @click="toggleActive"
          >
            {{ profile.is_active ? 'فعال' : 'غیرفعال' }}
          </span>
        </div>

        <div class="grid grid-cols-2 gap-3 text-sm">
          <div v-if="profile.gender">
            <p class="text-gray-400 text-xs">جنسیت</p>
            <p class="text-gray-800">{{ genderLabel[profile.gender] ?? profile.gender }}</p>
          </div>
          <div v-if="profile.date_of_birth">
            <p class="text-gray-400 text-xs">تاریخ تولد</p>
            <p class="text-gray-800">{{ profile.date_of_birth }}</p>
          </div>
          <div v-if="profile.height_cm">
            <p class="text-gray-400 text-xs">قد (سانتی‌متر)</p>
            <p class="text-gray-800">{{ profile.height_cm }}</p>
          </div>
          <div v-if="profile.notes">
            <p class="text-gray-400 text-xs">یادداشت</p>
            <p class="text-gray-800 col-span-2">{{ profile.notes }}</p>
          </div>
        </div>

        <template v-if="editMode">
          <div class="flex flex-col gap-2 border-t border-gray-100 pt-3">
            <label class="text-xs text-gray-500">تاریخ تولد (YYYY-MM-DD)</label>
            <input v-model="editDob" type="text" placeholder="1370-01-01" class="rounded-lg border border-gray-200 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
            <label class="text-xs text-gray-500">قد (سانتی‌متر)</label>
            <input v-model.number="editHeight" type="number" placeholder="170" class="rounded-lg border border-gray-200 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
            <p v-if="errorMsg" class="text-red-500 text-xs">{{ errorMsg }}</p>
            <div class="flex gap-2">
              <button :disabled="saving" class="flex-1 bg-blue-500 text-white rounded-lg py-2 text-sm disabled:opacity-50" @click="save">ذخیره</button>
              <button class="flex-1 border border-gray-200 rounded-lg py-2 text-sm text-gray-600" @click="editMode = false">انصراف</button>
            </div>
          </div>
        </template>
        <button v-else class="text-blue-500 text-sm text-start" @click="startEdit">ویرایش اطلاعات</button>
      </div>

      <!-- Quick actions -->
      <div class="grid grid-cols-2 gap-3">
        <NuxtLink
          :to="`/nutritionist/messages/${clientId}`"
          class="bg-white rounded-xl p-4 shadow-sm text-center text-blue-600 text-sm font-medium"
        >
          💬 پیام‌ها
        </NuxtLink>
        <NuxtLink
          :to="`/nutritionist/clients/${clientId}/tracking/body`"
          class="bg-white rounded-xl p-4 shadow-sm text-center text-emerald-600 text-sm font-medium"
        >
          📊 ردیابی
        </NuxtLink>
      </div>
    </template>
  </div>
</template>
