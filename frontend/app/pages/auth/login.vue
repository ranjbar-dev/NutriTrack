<script setup lang="ts">
import { ROLE_DEFAULT_ROUTES } from '~/utils/constants'

definePageMeta({ layout: 'auth', middleware: [] })

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const authStore = useAuthStore()

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await authStore.login(email.value, password.value)
    // Redirect per role (D-18)
    const route = authStore.getDefaultRoute()
    navigateTo(route)
  }
  catch (e: any) {
    error.value = e?.message || 'ایمیل یا رمز عبور اشتباه است'
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="w-full max-w-sm bg-white rounded-2xl shadow-lg p-8">
      <h1 class="text-2xl font-bold text-center mb-6">ورود</h1>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label for="email" class="block text-sm font-medium mb-1">ایمیل</label>
          <input
            id="email"
            v-model="email"
            type="email"
            required
            autocomplete="email"
            class="w-full border border-gray-300 rounded-lg px-4 py-3 text-start focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            placeholder="example@email.com"
            dir="ltr"
          />
        </div>
        <div>
          <label for="password" class="block text-sm font-medium mb-1">رمز عبور</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="current-password"
            class="w-full border border-gray-300 rounded-lg px-4 py-3 text-start focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            dir="ltr"
          />
        </div>

        <div v-if="error" class="text-red-500 text-sm text-center">{{ error }}</div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-emerald-600 text-white rounded-lg py-3 font-medium hover:bg-emerald-700 disabled:opacity-50 transition-colors"
        >
          {{ loading ? 'در حال ورود...' : 'ورود' }}
        </button>
      </form>

      <NuxtLink to="/auth/otp" class="block text-center text-sm text-emerald-600 mt-4">
        ورود با کد تایید (مراجعین)
      </NuxtLink>
    </div>
  </div>
</template>
