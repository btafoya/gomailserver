/**
 * Setup Wizard API
 * First-run configuration endpoints (no authentication required)
 */

import { useApiBase } from '../useApiBase'

// ============================================================================
// TypeScript Interfaces
// ============================================================================

export interface SetupStatus {
  setup_complete: boolean
  current_step?: string
}

export interface SetupState {
  current_step: 'system' | 'domain' | 'admin' | 'review' | 'complete'
  system_configured: boolean
  domain_configured: boolean
  admin_configured: boolean
  started_at: string
  updated_at: string
}

export interface AdminUserRequest {
  email: string
  full_name: string
  password: string
  enable_totp?: boolean
}

// ============================================================================
// API Functions
// ============================================================================

export const useSetupApi = () => {
  const API_BASE = useApiBase()

  /**
   * Check if setup is complete
   */
  const getSetupStatus = async (): Promise<SetupStatus> => {
    const response = await fetch(`${API_BASE}/setup/status`, {
      method: 'GET'
    })

    if (!response.ok) {
      throw new Error('Failed to check setup status')
    }

    return await response.json()
  }

  /**
   * Get current setup wizard state
   */
  const getSetupState = async (): Promise<SetupState> => {
    const response = await fetch(`${API_BASE}/setup/state`, {
      method: 'GET'
    })

    if (!response.ok) {
      throw new Error('Failed to get setup state')
    }

    return await response.json()
  }

  /**
   * Create first admin user
   */
  const createAdminUser = async (
    user: AdminUserRequest
  ): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/setup/admin`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(user)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create admin user')
    }

    return await response.json()
  }

  /**
   * Mark setup as complete
   */
  const completeSetup = async (): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/setup/complete`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error('Failed to complete setup')
    }

    return await response.json()
  }

  return {
    getSetupStatus,
    getSetupState,
    createAdminUser,
    completeSetup
  }
}
