<script setup lang="ts">
import AuthRolePicker from '../../components/auth/AuthRolePicker.vue'
import { resolveRoleHomePath } from '../../stores/auth-session'

definePageMeta({
  layout: 'auth'
})

const authStore = useAuthSessionStore()

const entries = [
  {
    id: 'client',
    title: 'ورود کاربر',
    description: 'دریافت کد تایید با شماره موبایل',
    href: '/auth/client'
  },
  {
    id: 'nutritionist',
    title: 'ورود متخصص تغذیه',
    description: 'ورود با ایمیل و رمز عبور',
    href: '/auth/nutritionist'
  },
  {
    id: 'admin',
    title: 'ورود ادمین',
    description: 'ورود مدیر سامانه با حساب سازمانی',
    href: '/auth/admin'
  }
] as const

onMounted(async () => {
  if (!authStore.role) {
    return
  }

  await navigateTo(resolveRoleHomePath(authStore.role))
})
</script>

<template>
  <section class="auth-gateway">
    <h2>ورود به نوتریتراک</h2>
    <p>نقش خود را انتخاب کنید.</p>
    <AuthRolePicker :entries="entries" />
  </section>
</template>

<style scoped>
.auth-gateway {
  display: grid;
  gap: var(--space-3);
}

h2,
p {
  margin: 0;
}
</style>
