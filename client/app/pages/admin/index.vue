<script setup lang="ts">
import { ref } from 'vue'
import AdminStatsKpiGrid from '~/components/admin/AdminStatsKpiGrid.vue'
import type { AdminStats } from '~/types/admin'
import { useAdminStatsApi } from '~/composables/useAdminStatsApi'

definePageMeta({
  layout: 'admin',
})

const api = useAdminStatsApi()

const statsState = ref<{
  loading: boolean
  error: string
  stats: AdminStats | null
}>({
  loading: true,
  error: '',
  stats: null,
})

async function refreshStats() {
  statsState.value.loading = true
  statsState.value.error = ''

  const { data, error } = await api.getStats()
  statsState.value.stats = data.value?.data ?? null
  if (error.value) {
    statsState.value.error = 'خطا در دریافت آمار پلتفرم'
  }

  statsState.value.loading = false
}

onMounted(() => {
  refreshStats()
})
</script>

<template>
  <section class="admin-dashboard">
    <header class="hero">
      <div>
        <h1>دید کلی پلتفرم</h1>
        <p>آمار کلیدی، مدیریت متخصصان تغذیه و دسترسی سریع به راهبری کاتالوگ.</p>
      </div>

      <nav class="actions">
        <NuxtLink to="/admin/nutritionists">مدیریت متخصصان</NuxtLink>
        <NuxtLink to="/admin/catalogue/foods">کاتالوگ غذاها</NuxtLink>
      </nav>
    </header>

    <p v-if="statsState.error" class="notice error">{{ statsState.error }}</p>

    <AdminStatsKpiGrid :stats="statsState.stats" :loading="statsState.loading" />

    <section class="entry-points">
      <NuxtLink class="entry-card" to="/admin/nutritionists">
        <strong>فهرست متخصصان</strong>
        <span>جستجو، ایجاد حساب و مشاهده وضعیت هر متخصص</span>
      </NuxtLink>
      <NuxtLink class="entry-card" to="/admin/catalogue/foods">
        <strong>راهبری کاتالوگ</strong>
        <span>کنترل غذاها، داروها و دسته‌بندی‌ها با دسترسی ادمین</span>
      </NuxtLink>
    </section>
  </section>
</template>

<style scoped>
.admin-dashboard {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hero {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.hero h1 {
  margin: 0;
  font-size: 1.2rem;
}

.hero p {
  margin: 4px 0 0;
  color: #53606a;
  line-height: 1.7;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.actions a,
.entry-card {
  border: 1px solid #d3dce2;
  border-radius: 12px;
  padding: 10px 12px;
  background: #fff;
  color: #173042;
  text-decoration: none;
}

.entry-points {
  display: grid;
  gap: 10px;
}

.entry-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.entry-card span {
  color: #53606a;
  font-size: 0.93rem;
}

.notice.error {
  margin: 0;
  color: #8b2121;
}
</style>
