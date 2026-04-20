export default defineNuxtConfig({
  compatibilityDate: '2025-07-18',
  future: { compatibilityVersion: 4 },

  css: [
    './app/assets/css/main.css',
    'vazirmatn/Vazirmatn-Variable-font-face.css',
  ],

  app: {
    head: {
      htmlAttrs: { dir: 'rtl', lang: 'fa' },
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1, maximum-scale=1' },
      ],
      title: 'نوتری‌ترک',
    },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api',
      vapidPublicKey: process.env.NUXT_PUBLIC_VAPID_PUBLIC_KEY || process.env.VITE_VAPID_PUBLIC_KEY || '',
    },
  },

  postcss: {
    plugins: {
      '@tailwindcss/postcss': {},
    },
  },

  modules: ['@pinia/nuxt', '@nuxt/eslint', '@vite-pwa/nuxt'],

  // D-01: injectManifest strategy — custom SW, not zero-config preset
  pwa: {
    strategies: 'injectManifest',
    srcDir: 'service-worker',
    filename: 'sw.ts',
    registerType: 'autoUpdate',

    injectManifest: {
      injectionPoint: 'self.__WB_MANIFEST',
      globPatterns: ['**/*.{js,css,woff2,png,svg,ico,webp}'],
    },

    // D-04: Persian-only manifest
    manifest: {
      name: 'نوتری‌ترک',
      short_name: 'نوتری‌ترک',
      description: 'مدیریت برنامه غذایی',
      lang: 'fa',
      dir: 'rtl',
      display: 'standalone',
      background_color: '#f9fafb',
      theme_color: '#16a34a',
      start_url: '/client/plan',
      scope: '/',
      icons: [
        { src: '/icons/icon-192.png', sizes: '192x192', type: 'image/png' },
        { src: '/icons/icon-384.png', sizes: '384x384', type: 'image/png' },
        { src: '/icons/icon-512.png', sizes: '512x512', type: 'image/png', purpose: 'any maskable' },
      ],
      shortcuts: [
        { name: 'برنامه', short_name: 'برنامه', url: '/client/plan', icons: [{ src: '/icons/shortcut-plan.png', sizes: '96x96' }] },
        { name: 'ثبت روزانه', short_name: 'ثبت', url: '/client/tracking', icons: [{ src: '/icons/shortcut-track.png', sizes: '96x96' }] },
        { name: 'پیام‌ها', short_name: 'پیام', url: '/client/messages', icons: [{ src: '/icons/shortcut-msg.png', sizes: '96x96' }] },
      ],
    },

    devOptions: {
      enabled: true,
      type: 'module',
    },
  },
})
