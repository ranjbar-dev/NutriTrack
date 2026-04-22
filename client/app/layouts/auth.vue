<script setup lang="ts">
import AppShell from '../components/platform/AppShell.vue'
import SessionExpiredNotice from '../components/auth/SessionExpiredNotice.vue'
import ConnectivityBanner from '../components/platform/ConnectivityBanner.vue'
import UpdateAvailableBanner from '../components/platform/UpdateAvailableBanner.vue'
import { createRefreshAction } from '../stores/platform-pwa'
import { usePlatformPwaStore } from '../stores/platform-pwa'

const pwaStore = usePlatformPwaStore()
const { $pwaUpdate } = useNuxtApp()
const route = useRoute()
const refreshAction = createRefreshAction(() => $pwaUpdate?.refresh?.())
const showSessionExpired = computed(() => String(route.query.reason ?? '') === 'session-expired')
</script>

<template>
  <AppShell role="auth" title="ورود">
    <UpdateAvailableBanner :visible="pwaStore.needRefresh" @refresh="refreshAction.refresh" @dismiss="pwaStore.setNeedRefresh(false)" />
    <ConnectivityBanner :offline="pwaStore.offline" />
    <SessionExpiredNotice v-if="showSessionExpired" />
    <slot />
  </AppShell>
</template>
