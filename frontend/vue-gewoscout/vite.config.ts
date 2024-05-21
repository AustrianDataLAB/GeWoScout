import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'


export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: `https://gewoscout-app-go.azurewebsites.net`, //https://localhost:3333
        changeOrigin: true,
        secure: false
      }
    }
  },
  plugins: [
    vue()
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
