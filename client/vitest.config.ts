import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '~': fileURLToPath(new URL('./app', import.meta.url)),
      '@': fileURLToPath(new URL('./app', import.meta.url)),
    }
  },
  test: {
    include: ['tests/**/*.spec.ts'],
    environment: 'happy-dom',
    globals: true,
    coverage: {
      provider: 'v8'
    },
    setupFiles: ['./tests/setup.ts']
  }
})
