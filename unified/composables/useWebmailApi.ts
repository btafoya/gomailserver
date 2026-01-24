import { useApiBase } from './useApiBase'

export const useWebmailApi = () => {
  const API_BASE = useApiBase()

  const getAuthToken = () => {
    return typeof window !== 'undefined' ? localStorage.getItem('token') : null
  }

  const getAuthHeaders = () => {
    const token = getAuthToken()
    return {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {})
    }
  }

  // Task completion
  const completeTask = async (messageId: number, completed: boolean) => {
    const response = await fetch(`${API_BASE}/webmail/messages/${messageId}/complete`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ completed })
    })

    if (!response.ok) {
      throw new Error('Failed to update task completion')
    }

    return await response.json()
  }

  return {
    completeTask
  }
}