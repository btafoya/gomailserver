/**
 * Queue API composable
 * Handles mail queue operations
 */

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

export interface QueueMessage {
  id: string
  from: string
  to: string
  subject: string
  size: number
  attempts: number
  status: string
  created_at: string
  error?: string
}

export const useQueueApi = () => {
  const getQueue = async (): Promise<QueueMessage[]> => {
    const response = await fetch(`${API_BASE}/queue`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch queue')
    }

    const data = await response.json()
    return data.data || []
  }

  const retryMessage = async (id: string): Promise<void> => {
    const response = await fetch(`${API_BASE}/queue/${id}/retry`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to retry message')
    }
  }

  const deleteMessage = async (id: string): Promise<void> => {
    const response = await fetch(`${API_BASE}/queue/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete message')
    }
  }

  const refreshQueue = async (): Promise<void> => {
    await fetch(`${API_BASE}/queue/refresh`, {
      method: 'POST',
      headers: getAuthHeaders()
    })
  }

  return {
    getQueue,
    retryMessage,
    deleteMessage,
    refreshQueue
  }
}
