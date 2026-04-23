<template>
  <div class="notification-form">
    <div class="preferences-list">
      <div 
        v-for="pref in preferenceFields"
        :key="pref.field"
        class="preference-toggle"
      >
        <label class="toggle-label">
          <input
            type="checkbox"
            :checked="localPrefs[pref.field]"
            :disabled="savingField === pref.field"
            @change="handleToggle(pref.field, $event)"
          />
          <span class="label-text">{{ pref.label }}</span>
        </label>
        <div v-if="savingField === pref.field" class="saving-indicator">
          <span class="spinner" />
        </div>
      </div>
    </div>

    <InlineNotice
      v-if="updateError"
      type="error"
    >
      {{ updateError }}
    </InlineNotice>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { NotificationPreferences, UpdateNotificationPreferencesRequest } from '~/types/notifications'
import { useNotificationApi } from '~/composables/useNotificationApi'
import InlineNotice from '~/components/platform/InlineNotice.vue'

const props = defineProps<{
  preferences: NotificationPreferences
}>()

const emit = defineEmits<{
  updated: [prefs: NotificationPreferences]
}>()

const notificationApi = useNotificationApi()
const localPrefs = reactive({ ...props.preferences })
const savingField = ref<string | null>(null)
const updateError = ref('')

const preferenceFields = [
  { field: 'meal_reminders', label: 'یادآور وعده‌های غذایی' },
  { field: 'water_reminders', label: 'یادآور مصرف آب' },
  { field: 'message_alerts', label: 'اعلان پیام‌های جدید' },
  { field: 'diet_updates', label: 'اعلان برنامه غذایی جدید' },
]

async function handleToggle(
  field: keyof UpdateNotificationPreferencesRequest,
  event: Event
) {
  const checkbox = event.target as HTMLInputElement
  const newValue = checkbox.checked
  const oldValue = localPrefs[field]

  // Optimistic update
  localPrefs[field] = newValue
  savingField.value = field
  updateError.value = ''

  try {
    const result = await notificationApi.updatePreferences({ [field]: newValue })
    emit('updated', result)
  } catch (error) {
    // Rollback
    localPrefs[field] = oldValue
    updateError.value = 'ذخیره‌سازی ناموفق بود — دوباره تلاش کنید'
    console.error('Failed to update preferences:', error)
  } finally {
    savingField.value = null
  }
}
</script>

<style scoped>
.notification-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preferences-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preference-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background-color: var(--color-surface-light, #f5f5f5);
  border-radius: 6px;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  flex: 1;
}

.toggle-label input {
  width: 20px;
  height: 20px;
  cursor: pointer;
}

.toggle-label input:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.label-text {
  font-size: 14px;
  color: var(--color-text-primary, #212121);
  font-weight: 500;
}

.saving-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
}

.spinner {
  width: 16px;
  height: 16px;
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
