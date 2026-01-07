// https://nuxt.com/docs/api/configuration
export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: [
    '@nuxt/ui',
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      // API base URL - can be overridden via NUXT_PUBLIC_API_BASE environment variable
      // In development, this is set by the control script based on gomailserver.yaml
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1'
    }
  },
  app: {
    head: {
      titleTemplate: (titleChunk) => {
        return titleChunk ? `${titleChunk} - Mail Server` : 'Mail Server'
      }
    }
  },
  nitro: {
    preset: 'node-server',
    output: {
      dir: '../unified-go/.output'
    }
  },
  // Vite dev server configuration
  vite: {
    server: {
      // Proxy API requests to the Go backend
      // Target port is read from NUXT_PUBLIC_API_BASE or defaults to 8080
      proxy: {
        '/api': {
          target: process.env.NUXT_PUBLIC_API_BASE?.replace('/api/v1', '') || 'http://localhost:8080',
          changeOrigin: true
        }
      }
    }
  }
})