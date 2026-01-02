// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@nuxt/icon'
  ],

  app: {
    baseURL: '/webmail/',
    head: {
      title: 'Webmail',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Modern webmail client' },
        { name: 'theme-color', content: '#3b82f6' },
        { name: 'apple-mobile-web-app-capable', content: 'yes' },
        { name: 'apple-mobile-web-app-status-bar-style', content: 'default' }
      ],
      link: [
        { rel: 'manifest', href: '/webmail/manifest.json' },
        { rel: 'icon', type: 'image/x-icon', href: '/webmail/favicon.ico' },
        { rel: 'apple-touch-icon', href: '/webmail/icon-192.png' }
      ]
    }
  },

  css: ['~/assets/css/main.css'],

  nitro: {
    devProxy: {
      '/api': {
        target: 'http://localhost:8980',
        changeOrigin: true
      }
    }
  },

  typescript: {
    strict: false,
    shim: false
  },

  vite: {
    vue: {
      script: {
        defineModel: true,
        propsDestructure: true
      }
    }
  }
})
