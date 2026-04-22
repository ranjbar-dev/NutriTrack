export default defineNuxtConfig({
  srcDir: 'app',
  modules: ['@pinia/nuxt', '@vite-pwa/nuxt'],
  pwa: {
    registerType: 'prompt',
    manifest: {
      name: 'NutriTrack',
      short_name: 'NutriTrack',
      lang: 'fa',
      dir: 'rtl',
      display: 'standalone',
      start_url: '/auth',
      background_color: '#f5f7f9',
      theme_color: '#0f3d5e',
      icons: [
        {
          src: '/icons/icon-192.png',
          sizes: '192x192',
          type: 'image/png'
        },
        {
          src: '/icons/icon-512.png',
          sizes: '512x512',
          type: 'image/png'
        }
      ]
    },
    workbox: {
      navigateFallback: '/auth',
      globPatterns: ['**/*.{js,css,html,ico,png,svg,woff2}'],
      runtimeCaching: [
        {
          urlPattern: ({ request }) => request.destination === 'document',
          handler: 'NetworkFirst',
          options: {
            cacheName: 'shell-documents'
          }
        },
        {
          // Keep authenticated API traffic network-only (api/** boundary).
          urlPattern: /\/api\/.*/,
          handler: 'NetworkOnly'
        }
      ]
    }
  }
})
