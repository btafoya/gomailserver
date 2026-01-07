import { defineStore } from 'pinia'
import { useApiBase } from '../composables/useApiBase'

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
      }
    },

    async login(credentials: { email: string, password: string }) {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/auth/login`, {
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
