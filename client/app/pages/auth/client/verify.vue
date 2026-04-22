<script setup lang="ts">
import OtpInput from '../../../components/auth/OtpInput.vue'
import { mapAuthError } from '../../../lib/auth/error-map'

definePageMeta({
  layout: 'auth'
})

const route = useRoute()
const authApi = useAuthApi()
const authStore = useAuthSessionStore()

const mobile = computed(() => String(route.query.mobile ?? '').replace(/[^0-9]/g, ''))
const cooldownSeed = Number.parseInt(String(route.query.cooldown ?? '0'), 10)

const otpCode = ref('')
const pendingVerify = ref(false)
const pendingResend = ref(false)
const errorMessage = ref('')
const cooldownSeconds = ref(Number.isFinite(cooldownSeed) ? Math.max(cooldownSeed, 0) : 0)
const failedAttempts = ref(0)

const lockedByAttempts = computed(() => failedAttempts.value >= 3)
const canVerify = computed(() => otpCode.value.length === 6 && !pendingVerify.value && !lockedByAttempts.value)
const canResend = computed(() => cooldownSeconds.value === 0 && !pendingResend.value)

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

async function verifyCode(): Promise<void> {
  if (!canVerify.value || !/^09\d{9}$/.test(mobile.value)) {
    return
  }

  pendingVerify.value = true
  errorMessage.value = ''

  try {
    const session = await authApi.verifyOtp({
      mobile: mobile.value,
      code: otpCode.value
    })
    authStore.applySession(session)
    await navigateTo('/client')
  } catch (error) {
    const safeError = mapAuthError((error as { ui?: { code?: string } })?.ui ?? null)
    errorMessage.value = safeError.message

    if (safeError.code === 'OTP_INVALID') {
      failedAttempts.value += 1
      otpCode.value = ''
    }

    if (safeError.code === 'OTP_MAX_ATTEMPTS') {
      failedAttempts.value = 3
    }

    if (safeError.code === 'OTP_EXPIRED') {
      failedAttempts.value = 3
    }
  } finally {
    pendingVerify.value = false
  }
}

async function resendCode(): Promise<void> {
  if (!canResend.value || !/^09\d{9}$/.test(mobile.value)) {
    return
  }

  pendingResend.value = true
  errorMessage.value = ''

  try {
    await authApi.sendOtp({
      mobile: mobile.value
    })
    failedAttempts.value = 0
    otpCode.value = ''
    startCooldown(60)
  } catch (error) {
    const safeError = mapAuthError((error as { ui?: { code?: string } })?.ui ?? null)
    errorMessage.value = safeError.message
  } finally {
    pendingResend.value = false
  }
}

onMounted(() => {
  if (cooldownSeconds.value > 0) {
    startCooldown(cooldownSeconds.value)
  }
})

onUnmounted(() => {
  if (cooldownTimer) {
    clearInterval(cooldownTimer)
  }
})
</script>

<template>
  <section class="auth-card">
    <h2>تایید کد کاربر</h2>
    <p>کد 6 رقمی ارسال شده به {{ mobile }} را وارد کنید.</p>

    <OtpInput v-model="otpCode" :disabled="pendingVerify || lockedByAttempts" />

    <p v-if="lockedByAttempts" class="notice error" role="alert">
      تعداد تلاش بیش از حد مجاز بود. کد جدید دریافت کنید.
    </p>
    <p v-else-if="errorMessage" class="notice error" role="alert">{{ errorMessage }}</p>

    <p v-if="cooldownSeconds > 0" class="notice" role="status">
      ارسال مجدد تا {{ cooldownSeconds }} ثانیه دیگر
    </p>

    <div class="actions">
      <button class="auth-button" :disabled="!canVerify" @click="verifyCode">
        تایید و ورود
      </button>
      <button class="secondary-button" :disabled="!canResend" @click="resendCode">
        دریافت کد جدید
      </button>
    </div>
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

.notice {
  margin: 0;
  color: #2d586f;
}

.notice.error {
  color: #8d2222;
}

.actions {
  display: grid;
  gap: var(--space-2);
}

.auth-button,
.secondary-button {
  min-height: 48px;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
}

.auth-button {
  border: 0;
  background: #0f6b7a;
  color: #ffffff;
}

.secondary-button {
  border: 1px solid #a9bdc8;
  background: #ffffff;
  color: #0f3d5e;
}
</style>
