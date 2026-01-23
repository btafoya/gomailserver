import { useApiBase } from '../useApiBase'

export const getAuthToken = () => {
  return typeof window !== 'undefined' ? localStorage.getItem('token') : null
}

export const getAuthHeaders = () => {
  const token = getAuthToken()
  return {
    'Content-Type': 'application/json',
    ...(token ? { 'Authorization': `Bearer ${token}` } : {})
  }
}

export interface RateLimitConfig {
  id: number
  type: 'smtp' | 'imap' | 'auth' | 'global'
  scope: 'per_user' | 'per_domain' | 'per_ip' | 'global'
  max_requests: number
  window_seconds: number
  penalty_seconds: number
  enabled: boolean
  description?: string
  created_at: string
  updated_at: string
}

export interface RateLimitCreateRequest {
  type: 'smtp' | 'imap' | 'auth' | 'global'
  scope: 'per_user' | 'per_domain' | 'per_ip' | 'global'
  max_requests: number
  window_seconds: number
  penalty_seconds?: number
  enabled?: boolean
  description?: string
}

export const useRateLimitsApi = () => {
  const API_BASE = useApiBase()

  const getRateLimits = async (type?: string): Promise<RateLimitConfig[]> => {
    const url = type 
      ? `${API_BASE}/rate-limits?type=${type}`
      : `${API_BASE}/rate-limits`
    
    const response = await fetch(url, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch rate limits')
    }

    const data = await response.json()
    return data.data || []
  }

  const getRateLimit = async (id: number): Promise<RateLimitConfig> => {
    const response = await fetch(`${API_BASE}/rate-limits/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch rate limit')
    }

    const data = await response.json()
    return data.data
  }

  const createRateLimit = async (config: RateLimitCreateRequest): Promise<RateLimitConfig> => {
    const response = await fetch(`${API_BASE}/rate-limits`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(config)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create rate limit')
    }

    const data = await response.json()
    return data.data
  }

  const updateRateLimit = async (id: number, config: Partial<RateLimitCreateRequest>): Promise<RateLimitConfig> => {
    const response = await fetch(`${API_BASE}/rate-limits/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(config)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to update rate limit')
    }

    const data = await response.json()
    return data.data
  }

  const deleteRateLimit = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/rate-limits/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete rate limit')
    }
  }

  const toggleRateLimit = async (id: number, enabled: boolean): Promise<void> => {
    const response = await fetch(`${API_BASE}/rate-limits/${id}/toggle`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify({ enabled })
    })

    if (!response.ok) {
      throw new Error('Failed to toggle rate limit')
    }
  }

  const getRateLimitStats = async (id: number): Promise<any> => {
    const response = await fetch(`${API_BASE}/rate-limits/${id}/stats`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch rate limit stats')
    }

    return response.json()
  }

  const getGlobalStats = async (): Promise<any> => {
    const response = await fetch(`${API_BASE}/rate-limits/stats`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch global rate limit stats')
    }

    return response.json()
  }

  const resetRateLimit = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/rate-limits/${id}/reset`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to reset rate limit')
    }
  }

  return {
    getRateLimits,
    getRateLimit,
    createRateLimit,
    updateRateLimit,
    deleteRateLimit,
    toggleRateLimit,
    getRateLimitStats,
    getGlobalStats,
    resetRateLimit
  }
}