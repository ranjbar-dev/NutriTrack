export default defineNuxtConfig({
  compatibilityDate: '2025-07-18',
  future: { compatibilityVersion: 4 },

  css: [
    '~/assets/css/main.css',
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
    },
  },

  modules: ['@pinia/nuxt', '@nuxt/eslint'],
})
