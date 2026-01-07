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

export interface User {
  id: number
  email: string
  full_name: string
  display_name?: string
  domain_id: number
  domain_name?: string
  quota: number
  used_quota: number
  status: string
  forward_to?: string
  auto_reply_enabled: boolean
  auto_reply_subject?: string
  auto_reply_body?: string
  spam_threshold: number
  totp_enabled: boolean
  created_at: string
  last_login?: string
}

export interface UserCreateRequest {
  email: string
  password: string
  full_name: string
  display_name?: string
  domain_id: number
  quota?: number
  status?: string
  forward_to?: string
  auto_reply_enabled?: boolean
  auto_reply_subject?: string
  auto_reply_body?: string
  spam_threshold?: number
}

export const useUsersApi = () => {
  /**
   * Get all users
   */
  const getUsers = async (): Promise<User[]> => {
    const response = await fetch(`${API_BASE}/users`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch users')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Get a specific user by ID
   */
  const getUser = async (id: number): Promise<User> => {
    const response = await fetch(`${API_BASE}/users/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch user')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Create a new user
   */
  const createUser = async (user: UserCreateRequest): Promise<User> => {
    const response = await fetch(`${API_BASE}/users`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(user)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create user')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Update a user
   */
  const updateUser = async (id: number, user: Partial<UserCreateRequest & { password?: string }>): Promise<User> => {
    const response = await fetch(`${API_BASE}/users/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(user)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to update user')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Delete a user
   */
  const deleteUser = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/users/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete user')
    }
  }

  /**
   * Reset user password
   */
  const resetPassword = async (id: number, newPassword: string): Promise<void> => {
    const response = await fetch(`${API_BASE}/users/${id}/password`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ new_password: newPassword })
    })

    if (!response.ok) {
      throw new Error('Failed to reset password')
    }
  }

  return {
    getUsers,
    getUser,
    createUser,
    updateUser,
    deleteUser,
    resetPassword
  }
}
