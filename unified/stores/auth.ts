import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null as string | null,
    user: null as any,
    isAuthenticated: false
  }),

  getters: {
    getToken: (state) => state.token,
    getUser: (state) => state.user,
    isLoggedIn: (state) => state.isAuthenticated
  },

  actions: {
    initializeAuth() {
      const token = localStorage.getItem('token')
      if (token) {
        this.token = token
        this.isAuthenticated = true
        // TODO: Validate token and fetch user info
      }
    },

    async login(credentials: { email: string, password: string }) {
      try {
        const { $api } = useNuxtApp()
        const response = await $api.post('/auth/login', credentials)

        this.token = response.data.token
        this.isAuthenticated = true

        localStorage.setItem('token', this.token)

        // Redirect to portal for logged in users
        const router = useRouter()
        await router.push('/portal')

        return response.data
      } catch (error) {
        throw error
      }
    },

    logout() {
      this.token = null
      this.user = null
      this.isAuthenticated = false
      localStorage.removeItem('token')

      const router = useRouter()
      router.push('/login')
    },

    checkAuth() {
      const token = localStorage.getItem('token')
      if (token) {
        this.token = token
        this.isAuthenticated = true
        return true
      }
      return false
    }
  }
})