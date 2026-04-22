<script setup lang="ts">
import AppShell from '../components/platform/AppShell.vue'
import ConnectivityBanner from '../components/platform/ConnectivityBanner.vue'
import UpdateAvailableBanner from '../components/platform/UpdateAvailableBanner.vue'
import { createRefreshAction } from '../stores/platform-pwa'
import { usePlatformPwaStore } from '../stores/platform-pwa'

const pwaStore = usePlatformPwaStore()
const { $pwaUpdate } = useNuxtApp()
const refreshAction = createRefreshAction(() => $pwaUpdate?.refresh?.())
</script>

<template>
  <AppShell role="nutritionist" title="پنل متخصص تغذیه">
    <UpdateAvailableBanner :visible="pwaStore.needRefresh" @refresh="refreshAction.refresh" @dismiss="pwaStore.setNeedRefresh(false)" />
    <ConnectivityBanner :offline="pwaStore.offline" />
    <slot />
  </AppShell>
</template>
