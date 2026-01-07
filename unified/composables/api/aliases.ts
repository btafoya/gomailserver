const API_BASE = 'http://localhost:8980/api/v1'

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

export interface Alias {
  id: number
  address: string
  destinations: string[]
  domain_id: number
  domain_name?: string
  status: string
  created_at: string
}

export interface AliasCreateRequest {
  address: string
  destinations: string[]
  domain_id: number
  status?: string
}

export const useAliasesApi = () => {
  /**
   * Get all aliases
   */
  const getAliases = async (): Promise<Alias[]> => {
    const response = await fetch(`${API_BASE}/aliases`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch aliases')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Get a specific alias by ID
   */
  const getAlias = async (id: number): Promise<Alias> => {
    const response = await fetch(`${API_BASE}/aliases/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch alias')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Create a new alias
   */
  const createAlias = async (alias: AliasCreateRequest): Promise<Alias> => {
    const response = await fetch(`${API_BASE}/aliases`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(alias)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create alias')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Update an alias
   */
  const updateAlias = async (id: number, alias: Partial<AliasCreateRequest>): Promise<Alias> => {
    const response = await fetch(`${API_BASE}/aliases/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(alias)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to update alias')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Delete an alias
   */
  const deleteAlias = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/aliases/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete alias')
    }
  }

  return {
    getAliases,
    getAlias,
    createAlias,
    updateAlias,
    deleteAlias
  }
}
