import { defineStore } from 'pinia'
import api from '@/api/axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
    refreshToken: localStorage.getItem('refreshToken') || null
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    currentUser: (state) => state.user
  },

  actions: {
    async login(email, password) {
      try {
        const response = await api.post('/api/v1/auth/login', {
          email,
          password
        })

        this.token = response.data.data.token
        this.refreshToken = response.data.data.refresh_token
        this.user = response.data.data.user

        localStorage.setItem('token', this.token)
        localStorage.setItem('refreshToken', this.refreshToken)

        // Set default authorization header
        api.defaults.headers.common['Authorization'] = `Bearer ${this.token}`

        return response.data.data
      } catch (error) {
        throw error
      }
    },

    async refresh() {
      try {
        const response = await api.post('/api/v1/auth/refresh', {
          refresh_token: this.refreshToken
        })

        this.token = response.data.data.token
        localStorage.setItem('token', this.token)
        api.defaults.headers.common['Authorization'] = `Bearer ${this.token}`

        return response.data.data
      } catch (error) {
        this.logout()
        throw error
      }
    },

    logout() {
      this.user = null
      this.token = null
      this.refreshToken = null
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      delete api.defaults.headers.common['Authorization']
    },

    initializeAuth() {
      if (this.token) {
        api.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
      }
    }
  }
})
