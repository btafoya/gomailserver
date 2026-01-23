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

export interface LogEntry {
  timestamp: string
  level: 'debug' | 'info' | 'warn' | 'error'
  service: string
  message: string
  context?: any
  user_id?: number
  ip_address?: string
  request_id?: string
}

export interface LogFilters {
  level?: 'debug' | 'info' | 'warn' | 'error'
  service?: string
  user_id?: number
  ip_address?: string
  start_time?: string
  end_time?: string
  search?: string
  limit?: number
  offset?: number
}

export const useLogsApi = () => {
  const API_BASE = useApiBase()

  const getLogs = async (filters: LogFilters = {}): Promise<{ entries: LogEntry[], total: number }> => {
    const params = new URLSearchParams()
    
    if (filters.level) params.append('level', filters.level)
    if (filters.service) params.append('service', filters.service)
    if (filters.user_id) params.append('user_id', filters.user_id.toString())
    if (filters.ip_address) params.append('ip_address', filters.ip_address)
    if (filters.search) params.append('search', filters.search)
    if (filters.start_time) params.append('start_time', filters.start_time)
    if (filters.end_time) params.append('end_time', filters.end_time)
    if (filters.limit) params.append('limit', filters.limit.toString())
    if (filters.offset) params.append('offset', filters.offset.toString())
    
    const response = await fetch(`${API_BASE}/logs?${params.toString()}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch logs')
    }

    return response.json()
  }

  const getLogServices = async (): Promise<string[]> => {
    const response = await fetch(`${API_BASE}/logs/services`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch log services')
    }

    return response.json()
  }

  const downloadLogs = async (filters: LogFilters = {}): Promise<Blob> => {
    const params = new URLSearchParams()
    
    if (filters.level) params.append('level', filters.level)
    if (filters.service) params.append('service', filters.service)
    if (filters.user_id) params.append('user_id', filters.user_id.toString())
    if (filters.ip_address) params.append('ip_address', filters.ip_address)
    if (filters.start_time) params.append('start_time', filters.start_time)
    if (filters.end_time) params.append('end_time', filters.end_time)
    if (filters.search) params.append('search', filters.search)
    
    const response = await fetch(`${API_BASE}/logs/download?${params.toString()}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to download logs')
    }

    return response.blob()
  }

  const clearLogs = async (filters: LogFilters = {}): Promise<{ cleared: number }> => {
    const params = new URLSearchParams()
    
    if (filters.level) params.append('level', filters.level)
    if (filters.service) params.append('service', filters.service)
    if (filters.user_id) params.append('user_id', filters.user_id.toString())
    if (filters.ip_address) params.append('ip_address', filters.ip_address)
    if (filters.start_time) params.append('start_time', filters.start_time)
    if (filters.end_time) params.append('end_time', filters.end_time)
    
    const response = await fetch(`${API_BASE}/logs/clear?${params.toString()}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to clear logs')
    }

    return response.json()
  }

  const getLogStats = async (filters: LogFilters = {}): Promise<any> => {
    const params = new URLSearchParams()
    
    if (filters.level) params.append('level', filters.level)
    if (filters.service) params.append('service', filters.service)
    if (filters.start_time) params.append('start_time', filters.start_time)
    if (filters.end_time) params.append('end_time', filters.end_time)
    
    const response = await fetch(`${API_BASE}/logs/stats?${params.toString()}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch log statistics')
    }

    return response.json()
  }

  const tailLogs = async (level?: string, service?: string): Promise<EventSource> => {
    const params = new URLSearchParams()
    if (level) params.append('level', level)
    if (service) params.append('service', service)
    
    const url = `${API_BASE}/logs/tail?${params.toString()}`
    
    const eventSource = new EventSource(url)
    
    // Add authorization after creation if token exists
    const token = getAuthToken()
    if (token) {
      // Note: EventSource doesn't support custom headers in all browsers
      // In production, you'd implement this via query parameter or WebSocket
      console.log('Log streaming requires token authentication')
    }
    
    return eventSource
  }

  return {
    getLogs,
    getLogServices,
    downloadLogs,
    clearLogs,
    getLogStats,
    tailLogs
  }
}