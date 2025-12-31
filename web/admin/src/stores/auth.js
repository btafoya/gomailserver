import { defineStore } from 'pinia'
import axios from 'axios'

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
        const response = await axios.post('/api/v1/auth/login', {
          email,
          password
        })

        this.token = response.data.token
        this.refreshToken = response.data.refresh_token
        this.user = response.data.user

        localStorage.setItem('token', this.token)
        localStorage.setItem('refreshToken', this.refreshToken)

        // Set default authorization header
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`

        return response.data
      } catch (error) {
        throw error
      }
    },

    async refresh() {
      try {
        const response = await axios.post('/api/v1/auth/refresh', {
          refresh_token: this.refreshToken
        })

        this.token = response.data.token
        localStorage.setItem('token', this.token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`

        return response.data
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
      delete axios.defaults.headers.common['Authorization']
    },

    initializeAuth() {
      if (this.token) {
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
      }
    }
  }
})
