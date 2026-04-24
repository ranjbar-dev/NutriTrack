<script setup lang="ts">
import { computed, ref } from 'vue'
import ClientProfileHeaderCard from '~/components/nutritionist/ClientProfileHeaderCard.vue'
import ClientProfileTabs, { type ProfileTab } from '~/components/nutritionist/ClientProfileTabs.vue'
import ClientSnapshotPanels from '~/components/nutritionist/ClientSnapshotPanels.vue'
import type { DietPlanFlat } from '~/types/diet-authoring'
import type { NutritionistClient } from '~/types/nutritionist-workspace'
import { useDietPlanAuthoringApi } from '~/composables/useDietPlanAuthoringApi'
import { useLabApi } from '~/composables/useLabApi'
import { useMessagingApi } from '~/composables/useMessagingApi'
import { useNutritionistClientApi } from '~/composables/useNutritionistClientApi'

definePageMeta({
  layout: 'nutritionist',
})

const route = useRoute()
const clientId = computed(() => String(route.params.id || ''))

const clientApi = useNutritionistClientApi()
const planApi = useDietPlanAuthoringApi()
const messagingApi = useMessagingApi()
const labApi = useLabApi()

const profile = ref<NutritionistClient | null>(null)
const plans = ref<DietPlanFlat[]>([])
const activeTab = ref<ProfileTab>('overview')
const messagePreviewDate = ref<string | null>(null)
const labPreviewDate = ref<string | null>(null)
const profileError = ref('')

const activePlan = computed(() => plans.value.find((plan) => plan.status === 'active') || null)
const archivedCount = computed(() => plans.value.filter((plan) => plan.status === 'archived').length)

async function loadProfile() {
  profileError.value = ''
  const [{ data: profileData, error: profileFetchError }, { data: plansData }, messagesResult, labsResult] = await Promise.all([
    clientApi.getClientProfile(clientId.value),
    planApi.listClientPlans(clientId.value),
    messagingApi.getNutritionistConversation(clientId.value, 1, 1),
    labApi.listLabResults(clientId.value, 1, 1),
  ])

  if (profileFetchError.value) {
    profileError.value = 'پروفایل مرجع قابل دسترس نیست.'
  }

  profile.value = profileData.value ?? null
  plans.value = plansData.value?.data ?? []

  messagePreviewDate.value = messagesResult.data.value?.data?.[0]?.created_at ?? null
  labPreviewDate.value = labsResult.data.value?.data?.[0]?.created_at ?? null
}

onMounted(() => {
  if (!clientId.value) {
    profileError.value = 'شناسه مرجع معتبر نیست.'
    return
  }
  loadProfile()
})
</script>

<template>
  <section class="profile-page">
    <p v-if="profileError" class="error">{{ profileError }}</p>
    <ClientProfileHeaderCard :client="profile" />

    <ClientProfileTabs v-model="activeTab" />

    <ClientSnapshotPanels
      v-if="activeTab === 'overview'"
      :active-plan="activePlan"
      :archived-count="archivedCount"
      :latest-message-at="messagePreviewDate"
      :latest-lab-at="labPreviewDate"
    />

    <section v-else-if="activeTab === 'plans'" class="tab-card">
      <h3>برنامه ها</h3>
      <p>برنامه فعال: {{ activePlan?.title || 'ثبت نشده' }}</p>
      <NuxtLink class="link" :to="`/nutritionist/clients/${clientId}/plans/new`">ایجاد برنامه جدید</NuxtLink>
    </section>

    <section v-else-if="activeTab === 'messages'" class="tab-card">
      <h3>پیام ها</h3>
      <p>آخرین پیام: {{ messagePreviewDate || 'موردی ثبت نشده است' }}</p>
      <NuxtLink class="link" to="/nutritionist/messages">مشاهده گفت وگوها</NuxtLink>
    </section>

    <section v-else class="tab-card">
      <h3>آزمایش ها</h3>
      <p>آخرین آزمایش: {{ labPreviewDate || 'موردی ثبت نشده است' }}</p>
      <NuxtLink class="link" :to="`/nutritionist/clients/${clientId}/labs`">ورود به فایل آزمایش</NuxtLink>
    </section>
  </section>
</template>

<style scoped>
.profile-page {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tab-card {
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 12px;
}

.tab-card h3 {
  margin: 0 0 8px;
}

.tab-card p {
  margin: 0 0 10px;
  color: #53606a;
}

.link {
  display: inline-block;
  text-decoration: none;
  color: #0f6b7a;
  font-weight: 600;
}

.error {
  margin: 0;
  color: #8b2121;
}
</style>
