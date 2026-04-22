<script setup lang="ts">
import AuthFormCard from '../../../components/auth/AuthFormCard.vue'
import { mapAuthError } from '../../../lib/auth/error-map'

definePageMeta({
  layout: 'auth'
})

const authApi = useAuthApi()
const authStore = useAuthSessionStore()

const email = ref('')
const password = ref('')
const pending = ref(false)
const errorMessage = ref('')

const validForm = computed(() => email.value.includes('@') && password.value.length >= 6)

async function submitCredentials(): Promise<void> {
  if (!validForm.value || pending.value) {
    return
  }

  pending.value = true
  errorMessage.value = ''

  try {
    const session = await authApi.login({
      email: email.value,
      password: password.value
    })
    authStore.applySession(session)
    await navigateTo('/nutritionist')
  } catch (error) {
    const safeError = mapAuthError((error as { ui?: { message?: string } })?.ui ?? null)
    errorMessage.value = safeError.message
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <AuthFormCard
    title="ورود متخصص تغذیه"
    description="برای ادامه با ایمیل و رمز عبور وارد شوید."
    :pending="pending"
    :error-message="errorMessage"
    @submit="submitCredentials"
  >
    <label for="nutritionist-email">ایمیل</label>
    <input
      id="nutritionist-email"
      v-model="email"
      class="auth-input"
      type="email"
      dir="ltr"
      autocomplete="email"
      :disabled="pending"
    >

    <label for="nutritionist-password">رمز عبور</label>
    <input
      id="nutritionist-password"
      v-model="password"
      class="auth-input"
      type="password"
      dir="ltr"
      autocomplete="current-password"
      :disabled="pending"
    >
  </AuthFormCard>
</template>

<style scoped>
.auth-input {
  min-height: 48px;
  border: 1px solid #c5d2da;
  border-radius: 10px;
  padding: 0 var(--space-3);
  font-size: 1rem;
}
</style>
