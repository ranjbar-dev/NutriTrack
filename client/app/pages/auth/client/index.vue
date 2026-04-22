<script setup lang="ts">
import { mapAuthError } from '../../../lib/auth/error-map'

definePageMeta({
  layout: 'auth'
})

const authApi = useAuthApi()
const mobileInput = ref('')
const isPending = ref(false)
const errorMessage = ref('')
const cooldownSeconds = ref(0)

const normalizedMobile = computed(() => mobileInput.value.replace(/[^0-9]/g, ''))
const validIranMobile = computed(() => /^09\d{9}$/.test(normalizedMobile.value))
const submitDisabled = computed(() => !validIranMobile.value || isPending.value || cooldownSeconds.value > 0)

let cooldownTimer: ReturnType<typeof setInterval> | null = null

function startCooldown(seconds: number): void {
  cooldownSeconds.value = seconds
  if (cooldownTimer) {
    clearInterval(cooldownTimer)
  }

  cooldownTimer = setInterval(() => {
    if (cooldownSeconds.value <= 1) {
      cooldownSeconds.value = 0
      if (cooldownTimer) {
        clearInterval(cooldownTimer)
        cooldownTimer = null
      }
      return
    }

    cooldownSeconds.value -= 1
  }, 1000)
}

async function submitOtpRequest(): Promise<void> {
  if (submitDisabled.value) {
    return
  }

  isPending.value = true
  errorMessage.value = ''

  try {
    await authApi.sendOtp({
      mobile: normalizedMobile.value
    })

    startCooldown(60)
    await navigateTo({
      path: '/auth/client/verify',
      query: {
        mobile: normalizedMobile.value,
        cooldown: '60'
      }
    })
  } catch (error) {
    const safeError = mapAuthError((error as { ui?: { code?: string } })?.ui ?? null)
    errorMessage.value = safeError.message
  } finally {
    isPending.value = false
  }
}

onUnmounted(() => {
  if (cooldownTimer) {
    clearInterval(cooldownTimer)
  }
})
</script>

<template>
  <section class="auth-card">
    <h2>ورود کاربر</h2>
    <p>شماره موبایل را وارد کنید تا کد تایید ارسال شود.</p>

    <label for="mobile">شماره موبایل</label>
    <input
      id="mobile"
      v-model="mobileInput"
      dir="ltr"
      inputmode="tel"
      placeholder="09123456789"
      class="auth-input"
      :disabled="isPending"
    >

    <p v-if="cooldownSeconds > 0" class="notice" role="status">
      ارسال مجدد تا {{ cooldownSeconds }} ثانیه دیگر فعال می شود.
    </p>

    <p v-if="errorMessage" class="notice error" role="alert">{{ errorMessage }}</p>

    <button class="auth-button" :disabled="submitDisabled" @click="submitOtpRequest">
      دریافت کد تایید
    </button>
  </section>
</template>

<style scoped>
.auth-card {
  display: grid;
  gap: var(--space-3);
  padding: var(--space-4);
  border-radius: 14px;
  background: #ffffff;
}

.auth-input {
  min-height: 48px;
  border: 1px solid #c5d2da;
  border-radius: 10px;
  padding: 0 var(--space-3);
  font-size: 1rem;
}

.notice {
  margin: 0;
  color: #2d586f;
}

.notice.error {
  color: #8d2222;
}

.auth-button {
  min-height: 48px;
  border: 0;
  border-radius: 10px;
  background: #0f6b7a;
  color: #ffffff;
  font-size: 1rem;
  font-weight: 600;
}
</style>
