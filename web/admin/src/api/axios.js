import axios from 'axios'
import router from '@/router'

// Create axios instance with default config
// Force empty baseURL so axios uses absolute paths from document root
const api = axios.create({
  baseURL: '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    // Get token from localStorage directly to avoid circular dependency
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors and token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // Handle 401 Unauthorized - token expired
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        // Import useAuthStore dynamically to avoid circular dependency
        const { useAuthStore } = await import('@/stores/auth')
        const authStore = useAuthStore()
        await authStore.refresh()

        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${authStore.token}`
        return api(originalRequest)
      } catch (refreshError) {
        // Refresh token failed, redirect to login
        const { useAuthStore } = await import('@/stores/auth')
        const authStore = useAuthStore()
        authStore.logout()
        router.push({ name: 'Login' })
        return Promise.reject(refreshError)
      }
    }

    // Handle 429 Rate Limit
    if (error.response?.status === 429) {
      console.error('Rate limit exceeded. Please try again later.')
    }

    return Promise.reject(error)
  }
)

export default api
