<template>
  <div class="notifications-settings-page">
    <AppShell title="تنظیمات اعلان" class="settings-shell">
      <div class="settings-section">
        <h3 class="section-title">اعلان‌های فوری</h3>
        <PushSubscriptionControl />
      </div>

      <div class="settings-section">
        <h3 class="section-title">تنظیمات اعلان</h3>
        <div v-if="loading" class="loading-state">
          <div class="spinner" />
        </div>
        <NotificationPreferencesForm 
          v-else-if="preferences"
          :preferences="preferences"
          @updated="preferences = $event"
        />
      </div>
    </AppShell>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { NotificationPreferences } from '~/types/notifications'
import { useNotificationApi } from '~/composables/useNotificationApi'
import AppShell from '~/components/platform/AppShell.vue'
import PushSubscriptionControl from '~/components/platform/PushSubscriptionControl.vue'
import NotificationPreferencesForm from '~/components/platform/NotificationPreferencesForm.vue'

definePageMeta({ layout: 'client' })

const notificationApi = useNotificationApi()
const preferences = ref<NotificationPreferences | null>(null)
const loading = ref(false)

async function loadPreferences() {
  loading.value = true
  try {
    const { data } = await notificationApi.getPreferences()
    if (data.value) {
      preferences.value = data.value
    }
  } catch (error) {
    console.error('Failed to load preferences:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadPreferences()
})
</script>

<style scoped>
.notifications-settings-page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--color-background, #fafafa);
}

.settings-shell {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.settings-section {
  padding: 16px;
  margin-bottom: 12px;
  background-color: var(--color-surface, #ffffff);
  border-bottom: 1px solid var(--color-border, #e0e0e0);
}

.section-title {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary, #212121);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border, #e0e0e0);
  border-top-color: var(--color-primary, #1976d2);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
