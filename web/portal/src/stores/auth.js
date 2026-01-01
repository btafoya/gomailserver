import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/axios'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('portal_token') || null)
  const refreshToken = ref(localStorage.getItem('portal_refresh_token') || null)
  const user = ref(JSON.parse(localStorage.getItem('portal_user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)

  async function login(email, password, totpCode = null) {
    const response = await api.post('/api/v1/auth/login', {
      email,
      password,
      totp_code: totpCode
    })

    token.value = response.data.token
    refreshToken.value = response.data.refresh_token
    user.value = response.data.user

    localStorage.setItem('portal_token', token.value)
    localStorage.setItem('portal_refresh_token', refreshToken.value)
    localStorage.setItem('portal_user', JSON.stringify(user.value))

    return response.data
  }

  async function refresh() {
    const response = await api.post('/api/v1/auth/refresh', {
      refresh_token: refreshToken.value
    })

    token.value = response.data.token
    localStorage.setItem('portal_token', token.value)

    return response.data
  }

  function logout() {
    token.value = null
    refreshToken.value = null
    user.value = null

    localStorage.removeItem('portal_token')
    localStorage.removeItem('portal_refresh_token')
    localStorage.removeItem('portal_user')
  }

  return {
    token,
    refreshToken,
    user,
    isAuthenticated,
    login,
    refresh,
    logout
  }
})
