<template>
  <div class="push-control">
    <div v-if="state === 'not-asked'" class="state-content">
      <button 
        class="primary-button"
        :disabled="loading"
        @click="handleSubscribe"
      >
        <span v-if="!loading">🔔 دریافت اعلان‌ها را فعال کنید</span>
        <span v-else>{{ loadingText }}</span>
      </button>
    </div>

    <div v-else-if="state === 'subscribed'" class="state-content subscribed">
      <div class="subscribed-indicator">✓ اعلان‌ها فعال است</div>
      <button 
        class="secondary-button"
        :disabled="loading"
        @click="handleUnsubscribe"
      >
        غیرفعال کردن
      </button>
    </div>

    <div v-else-if="state === 'blocked'">
      <InlineNotice type="warning">
        اعلان‌ها مسدود شده‌اند — تنظیمات مرورگر را بررسی کنید
      </InlineNotice>
    </div>

    <div v-else-if="state === 'unsupported'">
      <InlineNotice type="info">
        مرورگر شما از اعلان‌ها پشتیبانی نمی‌کند
      </InlineNotice>
    </div>

    <InlineNotice 
      v-if="error"
      type="error"
    >
      {{ error }}
    </InlineNotice>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { PushPermissionState } from '~/lib/push/subscription'
import {
  getPushPermissionState,
  subscribeToPush,
  unsubscribeFromPush,
} from '~/lib/push/subscription'
import { useNotificationApi } from '~/composables/useNotificationApi'
import { useRuntimeConfig } from '#app'
import InlineNotice from '~/components/platform/InlineNotice.vue'

const config = useRuntimeConfig()
const notificationApi = useNotificationApi()

const state = ref<PushPermissionState>('unsupported')
const loading = ref(false)
const error = ref('')
const loadingText = ref('در حال اتصال...')

async function initializeState() {
  state.value = await getPushPermissionState()
}

async function handleSubscribe() {
  if (!config.public.vapidKey) {
    error.value = 'VAPID key not configured'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const subscription = await subscribeToPush(config.public.vapidKey as string)

    if (!subscription) {
      // User denied permission or error occurred
      state.value = await getPushPermissionState()
      return
    }

    // Register with backend
    await notificationApi.registerPushSubscription(subscription)

    state.value = 'subscribed'
  } catch (err) {
    error.value = 'ثبت اعلان‌ها ناموفق بود'
    console.error('Failed to subscribe:', err)
  } finally {
    loading.value = false
  }
}

async function handleUnsubscribe() {
  loading.value = true
  error.value = ''

  try {
    // Unregister from backend
    await notificationApi.unregisterPushSubscription()

    // Unsubscribe locally
    await unsubscribeFromPush()

    state.value = 'not-asked'
  } catch (err) {
    error.value = 'لغو اعلان‌ها ناموفق بود'
    console.error('Failed to unsubscribe:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initializeState()
})
</script>

<style scoped>
.push-control {
  padding: 16px;
  background-color: var(--color-surface, #ffffff);
  border-radius: 8px;
  border: 1px solid var(--color-border, #e0e0e0);
}

.state-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.primary-button,
.secondary-button {
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.2s;
  white-space: nowrap;
}

.primary-button {
  background-color: var(--color-primary, #1976d2);
  color: white;
}

.primary-button:hover:not(:disabled) {
  background-color: var(--color-primary-dark, #1565c0);
}

.primary-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.secondary-button {
  background-color: var(--color-surface-light, #f5f5f5);
  color: var(--color-text-primary, #212121);
  border: 1px solid var(--color-border, #e0e0e0);
}

.secondary-button:hover:not(:disabled) {
  background-color: var(--color-border-light, #f0f0f0);
}

.subscribed-indicator {
  color: var(--color-success, #4caf50);
  font-weight: 500;
  flex: 1;
}
</style>
