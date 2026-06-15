import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'

// Proxy target — uses VITE_BACKEND_URL in Docker, falls back to localhost
const proxyTarget = process.env.VITE_BACKEND_URL || 'http://localhost:8080'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/be': {
        target: proxyTarget,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/be/, ''),
      },
      '/storage': {
        target: proxyTarget,
        changeOrigin: true,
      },
    },
  },
})
