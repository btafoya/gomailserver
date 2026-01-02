import { defineStore } from 'pinia'
import axios from 'axios'

interface User {
  id: number
  email: string
  full_name: string
  is_admin: boolean
}

interface AuthState {
  token: string | null
  user: User | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: null,
    user: null
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user
  },

  actions: {
    async login(email: string, password: string) {
      try {
        const response = await axios.post('/api/v1/auth/login', {
          email,
          password
        })

        this.token = response.data.token
        this.user = response.data.user

        // Store in localStorage for persistence
        if (process.client) {
          localStorage.setItem('auth_token', this.token)
          localStorage.setItem('auth_user', JSON.stringify(this.user))
        }

        // Set default authorization header
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
      } catch (error) {
        console.error('Login failed:', error)
        throw error
      }
    },

    async refreshToken() {
      try {
        if (!this.token) return

        const response = await axios.post('/api/v1/auth/refresh', {
          token: this.token
        })

        this.token = response.data.token

        if (process.client) {
          localStorage.setItem('auth_token', this.token)
        }

        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
      } catch (error) {
        console.error('Token refresh failed:', error)
        this.logout()
      }
    },

    logout() {
      this.token = null
      this.user = null

      if (process.client) {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
      }

      delete axios.defaults.headers.common['Authorization']
    },

    initializeAuth() {
      if (process.client) {
        const token = localStorage.getItem('auth_token')
        const userStr = localStorage.getItem('auth_user')

        if (token && userStr) {
          this.token = token
          this.user = JSON.parse(userStr)
          axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        }
      }
    }
  }
})
