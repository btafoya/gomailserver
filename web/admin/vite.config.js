import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: '/admin/',
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: '0.0.0.0',  // Listen on all network interfaces for remote access
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8980',
        changeOrigin: true,
      },
    },
  },
})
