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

export interface BlacklistEntry {
  id: number
  type: 'ip' | 'email' | 'domain'
  value: string
  reason?: string
  expires_at?: string
  created_at: string
  updated_at: string
  active: boolean
}

export interface BlacklistCreateRequest {
  type: 'ip' | 'email' | 'domain'
  value: string
  reason?: string
  expires_at?: string
}

export const useBlacklistApi = () => {
  const API_BASE = useApiBase()

  const getBlacklistEntries = async (type?: string): Promise<BlacklistEntry[]> => {
    const url = type 
      ? `${API_BASE}/blacklist?type=${type}`
      : `${API_BASE}/blacklist`
    
    const response = await fetch(url, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch blacklist entries')
    }

    const data = await response.json()
    return data.data || []
  }

  const getBlacklistEntry = async (id: number): Promise<BlacklistEntry> => {
    const response = await fetch(`${API_BASE}/blacklist/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch blacklist entry')
    }

    const data = await response.json()
    return data.data
  }

  const createBlacklistEntry = async (entry: BlacklistCreateRequest): Promise<BlacklistEntry> => {
    const response = await fetch(`${API_BASE}/blacklist`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(entry)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create blacklist entry')
    }

    const data = await response.json()
    return data.data
  }

  const updateBlacklistEntry = async (id: number, entry: Partial<BlacklistCreateRequest>): Promise<BlacklistEntry> => {
    const response = await fetch(`${API_BASE}/blacklist/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(entry)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to update blacklist entry')
    }

    const data = await response.json()
    return data.data
  }

  const deleteBlacklistEntry = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/blacklist/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete blacklist entry')
    }
  }

  const toggleBlacklistEntry = async (id: number, active: boolean): Promise<void> => {
    const response = await fetch(`${API_BASE}/blacklist/${id}/toggle`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify({ active })
    })

    if (!response.ok) {
      throw new Error('Failed to toggle blacklist entry')
    }
  }

  const importBlacklist = async (entries: BlacklistCreateRequest[]): Promise<{ imported: number; skipped: number }> => {
    const response = await fetch(`${API_BASE}/blacklist/import`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ entries })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to import blacklist entries')
    }

    return response.json()
  }

  const exportBlacklist = async (type?: string): Promise<Blob> => {
    const url = type 
      ? `${API_BASE}/blacklist/export?type=${type}`
      : `${API_BASE}/blacklist/export`
    
    const response = await fetch(url, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to export blacklist')
    }

    return response.blob()
  }

  return {
    getBlacklistEntries,
    getBlacklistEntry,
    createBlacklistEntry,
    updateBlacklistEntry,
    deleteBlacklistEntry,
    toggleBlacklistEntry,
    importBlacklist,
    exportBlacklist
  }
}