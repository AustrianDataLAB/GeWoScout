import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  // dev proxy for local development and test
  /* 
  server: {
    proxy: {
      '/api': {
        target: 'https://gewoscout-function-app.azurewebsites.net',
        changeOrigin: true,
        secure: false,
      },
    },
  }, 
  */
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
