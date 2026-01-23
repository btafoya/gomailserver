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

export interface Mailbox {
  id: number
  name: string
  email: string
  domain_id: number
  quota: number
  used: number
  message_count: number
  status: string
  auto_reply_enabled?: boolean
  auto_reply_message?: string
  forwarding_enabled?: boolean
  forwarding_address?: string
  created_at: string
  updated_at: string
}

export interface MailboxCreateRequest {
  name: string
  email: string
  domain_id: number
  quota?: number
  auto_reply_enabled?: boolean
  auto_reply_message?: string
  forwarding_enabled?: boolean
  forwarding_address?: string
}

export const useMailboxApi = () => {
  const API_BASE = useApiBase()

  const getMailboxes = async (domainId?: number): Promise<Mailbox[]> => {
    const url = domainId 
      ? `${API_BASE}/domains/${domainId}/mailboxes`
      : `${API_BASE}/mailboxes`
    
    const response = await fetch(url, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch mailboxes')
    }

    const data = await response.json()
    return data.data || []
  }

  const getMailbox = async (id: number): Promise<Mailbox> => {
    const response = await fetch(`${API_BASE}/mailboxes/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch mailbox')
    }

    const data = await response.json()
    return data.data
  }

  const createMailbox = async (mailbox: MailboxCreateRequest): Promise<Mailbox> => {
    const response = await fetch(`${API_BASE}/mailboxes`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(mailbox)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create mailbox')
    }

    const data = await response.json()
    return data.data
  }

  const updateMailbox = async (id: number, mailbox: Partial<MailboxCreateRequest>): Promise<Mailbox> => {
    const response = await fetch(`${API_BASE}/mailboxes/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(mailbox)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to update mailbox')
    }

    const data = await response.json()
    return data.data
  }

  const deleteMailbox = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/mailboxes/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete mailbox')
    }
  }

  const getMailboxStats = async (id: number): Promise<any> => {
    const response = await fetch(`${API_BASE}/mailboxes/${id}/stats`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch mailbox stats')
    }

    return response.json()
  }

  const updateMailboxQuota = async (id: number, quota: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/mailboxes/${id}/quota`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify({ quota })
    })

    if (!response.ok) {
      throw new Error('Failed to update mailbox quota')
    }
  }

  return {
    getMailboxes,
    getMailbox,
    createMailbox,
    updateMailbox,
    deleteMailbox,
    getMailboxStats,
    updateMailboxQuota
  }
}