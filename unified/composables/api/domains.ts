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

export interface Domain {
  id: number
  name: string
  description?: string
  status: string
  max_users?: number
  default_quota?: number
  dkim_selector?: string
  dkim_public_key?: string
  spf_record?: string
  dmarc_policy?: string
  dmarc_report_email?: string
  dkim_signing_enabled: boolean
  dkim_verify_enabled: boolean
  created_at: string
  updated_at: string
}

export interface DomainCreateRequest {
  name: string
  description?: string
  max_users?: number
  default_quota?: number
  dkim_selector?: string
  dkim_private_key?: string
  dkim_public_key?: string
  spf_record?: string
  dmarc_policy?: string
  dmarc_report_email?: string
  dkim_signing_enabled?: boolean
  dkim_verify_enabled?: boolean
}

export const useDomainsApi = () => {
  const getDomains = async (): Promise<Domain[]> => {
    const response = await fetch(`${API_BASE}/domains`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch domains')
    }

    const data = await response.json()
    return data.data || []
  }

  const getDomain = async (id: number): Promise<Domain> => {
    const response = await fetch(`${API_BASE}/domains/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch domain')
    }

    const data = await response.json()
    return data.data
  }

  const createDomain = async (domain: DomainCreateRequest): Promise<Domain> => {
    const response = await fetch(`${API_BASE}/domains`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(domain)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create domain')
    }

    const data = await response.json()
    return data.data
  }

  const updateDomain = async (id: number, domain: Partial<DomainCreateRequest>): Promise<Domain> => {
    const response = await fetch(`${API_BASE}/domains/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(domain)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to update domain')
    }

    const data = await response.json()
    return data.data
  }

  const deleteDomain = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/domains/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete domain')
    }
  }

  const generateDKIM = async (id: number): Promise<{ message: string; domain: string }> => {
    const response = await fetch(`${API_BASE}/domains/${id}/dkim`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to generate DKIM keys')
    }

    return response.json()
  }

  return {
    getDomains,
    getDomain,
    createDomain,
    updateDomain,
    deleteDomain,
    generateDKIM
  }
}
