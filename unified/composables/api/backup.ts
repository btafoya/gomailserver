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

export interface Backup {
  id: number
  name: string
  type: 'manual' | 'scheduled'
  size: number
  created_at: string
  completed_at?: string
  status: 'pending' | 'in_progress' | 'completed' | 'failed'
  description?: string
  file_path?: string
  compressed: boolean
  includes_tables?: string[]
}

export interface BackupCreateRequest {
  name: string
  description?: string
  type: 'manual' | 'scheduled'
  schedule?: {
    enabled: boolean
    frequency: 'daily' | 'weekly' | 'monthly'
    time: string
    retention_days: number
  }
  include_tables?: string[]
  compress: boolean
}

export const useBackupApi = () => {
  const API_BASE = useApiBase()

  const getBackups = async (): Promise<Backup[]> => {
    const response = await fetch(`${API_BASE}/backups`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch backups')
    }

    const data = await response.json()
    return data.data || []
  }

  const getBackup = async (id: number): Promise<Backup> => {
    const response = await fetch(`${API_BASE}/backups/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch backup')
    }

    const data = await response.json()
    return data.data
  }

  const createBackup = async (backup: BackupCreateRequest): Promise<Backup> => {
    const response = await fetch(`${API_BASE}/backups`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(backup)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to create backup')
    }

    const data = await response.json()
    return data.data
  }

  const deleteBackup = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/backups/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete backup')
    }
  }

  const downloadBackup = async (id: number): Promise<Blob> => {
    const response = await fetch(`${API_BASE}/backups/${id}/download`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to download backup')
    }

    return response.blob()
  }

  const restoreBackup = async (id: number): Promise<{ job_id: string }> => {
    const response = await fetch(`${API_BASE}/backups/${id}/restore`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to start restore')
    }

    return response.json()
  }

  const getRestoreStatus = async (jobId: string): Promise<any> => {
    const response = await fetch(`${API_BASE}/backups/restore/${jobId}/status`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch restore status')
    }

    return response.json()
  }

  const getBackupStats = async (): Promise<any> => {
    const response = await fetch(`${API_BASE}/backups/stats`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch backup statistics')
    }

    return response.json()
  }

  const cleanupBackups = async (retentionDays: number): Promise<{ deleted: number }> => {
    const response = await fetch(`${API_BASE}/backups/cleanup`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ retention_days: retentionDays })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to cleanup backups')
    }

    return response.json()
  }

  const testBackupConfig = async (): Promise<any> => {
    const response = await fetch(`${API_BASE}/backups/test`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'Failed to test backup configuration')
    }

    return response.json()
  }

  return {
    getBackups,
    getBackup,
    createBackup,
    deleteBackup,
    downloadBackup,
    restoreBackup,
    getRestoreStatus,
    getBackupStats,
    cleanupBackups,
    testBackupConfig
  }
}