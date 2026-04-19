<script setup lang="ts">
definePageMeta({ layout: 'auth', middleware: [] })

const mobile = ref('')
const code = ref('')
const step = ref<'request' | 'verify'>('request')
const error = ref('')
const loading = ref(false)
const authStore = useAuthStore()

async function handleRequestOTP() {
  error.value = ''
  loading.value = true
  try {
    await authStore.requestOTP(mobile.value)
    step.value = 'verify'
  }
  catch (e: any) {
    error.value = e?.message || 'خطا در ارسال کد'
  }
  finally {
    loading.value = false
  }
}

async function handleVerifyOTP() {
  error.value = ''
  loading.value = true
  try {
    await authStore.verifyOTP(mobile.value, code.value)
    navigateTo('/client/plan') // D-18: clients → /client/plan
  }
  catch (e: any) {
    error.value = e?.message || 'کد تایید اشتباه است'
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="w-full max-w-sm bg-white rounded-2xl shadow-lg p-8">
      <h1 class="text-2xl font-bold text-center mb-6">ورود مراجعین</h1>

      <!-- Step 1: Request OTP -->
      <form v-if="step === 'request'" @submit.prevent="handleRequestOTP" class="space-y-4">
        <div>
          <label for="mobile" class="block text-sm font-medium mb-1">شماره موبایل</label>
          <input
            id="mobile"
            v-model="mobile"
            type="tel"
            required
            inputmode="numeric"
            class="w-full border border-gray-300 rounded-lg px-4 py-3 text-start focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            placeholder="۰۹۱۲۳۴۵۶۷۸۹"
            dir="ltr"
          />
        </div>
        <div v-if="error" class="text-red-500 text-sm text-center">{{ error }}</div>
        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-emerald-600 text-white rounded-lg py-3 font-medium hover:bg-emerald-700 disabled:opacity-50 transition-colors"
        >
          {{ loading ? 'در حال ارسال...' : 'ارسال کد تایید' }}
        </button>
      </form>

      <!-- Step 2: Verify OTP -->
      <form v-else @submit.prevent="handleVerifyOTP" class="space-y-4">
        <p class="text-sm text-gray-600 text-center mb-2">
          کد تایید به شماره {{ toPersianDigits(mobile) }} ارسال شد
        </p>
        <div>
          <label for="code" class="block text-sm font-medium mb-1">کد تایید</label>
          <input
            id="code"
            v-model="code"
            type="text"
            required
            inputmode="numeric"
            maxlength="6"
            autocomplete="one-time-code"
            class="w-full border border-gray-300 rounded-lg px-4 py-3 text-center text-2xl tracking-widest focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            dir="ltr"
          />
        </div>
        <div v-if="error" class="text-red-500 text-sm text-center">{{ error }}</div>
        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-emerald-600 text-white rounded-lg py-3 font-medium hover:bg-emerald-700 disabled:opacity-50 transition-colors"
        >
          {{ loading ? 'در حال تایید...' : 'تایید و ورود' }}
        </button>
        <button
          type="button"
          class="w-full text-sm text-gray-500 mt-2"
          @click="step = 'request'; code = ''; error = ''"
        >
          ارسال مجدد کد
        </button>
      </form>

      <NuxtLink to="/auth/login" class="block text-center text-sm text-emerald-600 mt-4">
        ورود با ایمیل (متخصصین تغذیه)
      </NuxtLink>
    </div>
  </div>
</template>
