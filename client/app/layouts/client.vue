<script setup lang="ts">
import AppShell from '../components/platform/AppShell.vue'
import ConnectivityBanner from '../components/platform/ConnectivityBanner.vue'
import InstallPromptBanner from '../components/platform/InstallPromptBanner.vue'
import UpdateAvailableBanner from '../components/platform/UpdateAvailableBanner.vue'
import { createRefreshAction } from '../stores/platform-pwa'
import { usePlatformPwaStore } from '../stores/platform-pwa'

const pwaStore = usePlatformPwaStore()
const { $pwaUpdate } = useNuxtApp()
const refreshAction = createRefreshAction(() => $pwaUpdate?.refresh?.())

onMounted(() => {
  pwaStore.openInstallPrompt('role-shell-ready')
})
</script>

<template>
  <AppShell role="client" title="پنل کاربر">
    <UpdateAvailableBanner :visible="pwaStore.needRefresh" @refresh="refreshAction.refresh" @dismiss="pwaStore.setNeedRefresh(false)" />
    <InstallPromptBanner :visible="pwaStore.showInstallPrompt" @install="pwaStore.closeInstallPrompt" @dismiss="pwaStore.closeInstallPrompt" />
    <ConnectivityBanner :offline="pwaStore.offline" />
    <slot />
  </AppShell>
</template>
