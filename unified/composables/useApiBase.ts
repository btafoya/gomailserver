/**
 * Global composable for API configuration
 * Provides the API base URL from runtime config or environment variables
 */

/**
 * Get the API base URL
 * - Development: Uses relative path '/api/v1' (proxied by Nuxt)
 * - Production: Uses NUXT_PUBLIC_API_BASE env var or full URL
 */
export const useApiBase = () => {
  // In development, use relative path which gets proxied by Nuxt to Go server
  if (process.env.NODE_ENV === 'development') {
    return '/api/v1'
  }

  // In production, use the full URL from environment
  if (typeof window === 'undefined') {
    // Server-side
    return process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8980/api/v1'
  }

  // Client-side - get from Nuxt runtime config
  // @ts-ignore
  const config = window.__NUXT__?.config?.public || {}
  return config.apiBase || 'http://localhost:8980/api/v1'
}

/**
 * Get the Go backend port (from NUXT_PUBLIC_API_BASE)
 * Useful for determining where the Go server is running
 */
export const useGoBackendPort = () => {
  const apiBase = useApiBase()
  // Extract port from http://localhost:PORT/api/v1
  const match = apiBase.match(/localhost:(\d+)/)
  return match ? parseInt(match[1], 10) : 8980
}
