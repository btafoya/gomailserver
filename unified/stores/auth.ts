import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null as string | null,
    user: null as any,
    isAuthenticated: false
  }),

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
        const response = await fetch('http://localhost:8980/api/v1/auth/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(credentials)
        })
        const data = await response.json()
        this.token = data.token as string
        this.isAuthenticated = true

        localStorage.setItem('token', this.token)

        // Navigate using window.location (works outside Nuxt context)
        window.location.href = '/portal'
        return data
      } catch (error) {
        throw error
      }
    },

    logout() {
      this.token = null
      this.user = null
      this.isAuthenticated = false
      localStorage.removeItem('token')

      // Navigate using window.location
      window.location.href = '/login'
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
