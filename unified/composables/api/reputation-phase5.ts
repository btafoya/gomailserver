/**
 * Reputation Management API - Phase 5
 * Advanced reputation endpoints: DMARC, ARF, External Metrics, Provider Limits, Warmup, Predictions, Alerts
 */

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

// ============================================================================
// TypeScript Interfaces
// ============================================================================

// DMARC Reports
export interface DMARCReport {
  id: number
  report_id: string
  domain: string
  org_name: string
  email_address: string
  begin_time: number
  end_time: number
  record_count: number
  spf_aligned_count: number
  dkim_aligned_count: number
  failed_count: number
  created_at: string
}

export interface DMARCStatistics {
  total_reports: number
  total_records: number
  spf_aligned_count: number
  dkim_aligned_count: number
  both_aligned_count: number
  failed_count: number
  spf_alignment_rate: number
  dkim_alignment_rate: number
  overall_alignment_rate: number
}

export interface DMARCAction {
  id: number
  domain: string
  action_type: string
  description: string
  auto_applied: boolean
  applied_at: string
}

// ARF Reports
export interface ARFReport {
  id: number
  domain: string
  reporter: string
  complaint_date: string
  complaint_type: string
  recipient_email: string
  processed: boolean
  created_at: string
}

export interface ARFStatistics {
  total_reports: number
  processed_reports: number
  pending_reports: number
  complaint_types: Record<string, number>
}

// External Metrics
export interface PostmasterMetrics {
  id: number
  domain: string
  reputation: string
  spam_rate: number
  spf_rate: number
  dkim_rate: number
  dmarc_rate: number
  encryption_rate: number
  sample_count: number
  fetched_at: string
}

export interface SNDSMetrics {
  id: number
  ip_address: string
  reputation_score: number
  spam_trap_hits: number
  complaint_rate: number
  filter_level: 'GREEN' | 'YELLOW' | 'RED'
  sample_count: number
  fetched_at: string
}

export interface ExternalMetricsTrends {
  domain: string
  dates: string[]
  reputation_scores: number[]
  spam_rates: number[]
}

// Provider Rate Limits
export interface ProviderRateLimit {
  id: number
  provider: string
  domain: string
  daily_limit: number
  current_usage: number
  utilization_percentage: number
  last_reset: string
  created_at: string
}

// Custom Warmup
export interface WarmupSchedule {
  id: number
  domain: string
  template_type: 'conservative' | 'moderate' | 'aggressive' | 'custom'
  start_date: string
  current_day: number
  daily_volumes: Record<string, number>
  progress_percentage: number
  completed: boolean
  created_at: string
}

export interface WarmupTemplate {
  name: string
  days: number
  daily_volumes: Record<string, number>
}

// Predictions
export interface Prediction {
  id: number
  domain: string
  prediction_type: string
  predicted_score: number
  current_score: number
  trend: 'improving' | 'stable' | 'declining'
  confidence: 'high' | 'medium' | 'low'
  horizon_days: number
  factors: string[]
  created_at: string
}

export interface PredictionHistory {
  id: number
  domain: string
  prediction_score: number
  actual_score: number
  accuracy: number
  horizon_days: number
  created_at: string
}

// Phase 5 Alerts
export interface Phase5Alert {
  id: number
  domain?: string
  type: 'dmarc_alignment' | 'arf_complaint' | 'external_deterioration' | 'low_reputation'
  severity: 'critical' | 'high' | 'medium' | 'low'
  message: string
  acknowledged: boolean
  resolved: boolean
  created_at: string
}

// ============================================================================
// API Functions
// ============================================================================

export const useReputationPhase5Api = () => {
  const API_BASE = useApiBase()

  // --------------------------------------------------------------------------
  // DMARC Reports
  // --------------------------------------------------------------------------
  
  /**
   * List DMARC reports with pagination
   */
  const listDMARCReports = async (page = 1, pageSize = 20): Promise<{
    data: DMARCReport[]
    total: number
    page: number
    page_size: number
  }> => {
    const response = await fetch(
      `${API_BASE}/reputation/dmarc/reports?page=${page}&page_size=${pageSize}`,
      {
        method: 'GET',
        headers: getAuthHeaders()
      }
    )

    if (!response.ok) {
      throw new Error('Failed to fetch DMARC reports')
    }

    return await response.json()
  }

  /**
   * Get DMARC report by ID
   */
  const getDMARCReport = async (id: number): Promise<DMARCReport> => {
    const response = await fetch(`${API_BASE}/reputation/dmarc/reports/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch DMARC report')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Get DMARC statistics for a domain
   */
  const getDMARCStats = async (domain: string): Promise<DMARCStatistics> => {
    const response = await fetch(`${API_BASE}/reputation/dmarc/stats/${domain}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch DMARC statistics')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Get DMARC actions
   */
  const getDMARCActions = async (): Promise<DMARCAction[]> => {
    const response = await fetch(`${API_BASE}/reputation/dmarc/actions`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch DMARC actions')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Export DMARC report
   */
  const exportDMARCReport = async (id: number): Promise<Blob> => {
    const response = await fetch(`${API_BASE}/reputation/dmarc/reports/${id}/export`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to export DMARC report')
    }

    return await response.blob()
  }

  // --------------------------------------------------------------------------
  // ARF Reports
  // --------------------------------------------------------------------------
  
  /**
   * List ARF reports
   */
  const listARFReports = async (): Promise<ARFReport[]> => {
    const response = await fetch(`${API_BASE}/reputation/arf/reports`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch ARF reports')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Get ARF statistics
   */
  const getARFStats = async (): Promise<ARFStatistics> => {
    const response = await fetch(`${API_BASE}/reputation/arf/stats`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch ARF statistics')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Process ARF report
   */
  const processARFReport = async (id: number): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/reputation/arf/reports/${id}/process`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to process ARF report')
    }

    return await response.json()
  }

  // --------------------------------------------------------------------------
  // External Metrics
  // --------------------------------------------------------------------------
  
  /**
   * Get Gmail Postmaster metrics for a domain
   */
  const getPostmasterMetrics = async (domain: string): Promise<PostmasterMetrics> => {
    const response = await fetch(`${API_BASE}/reputation/external/postmaster/${domain}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch Postmaster metrics')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Get Microsoft SNDS metrics for an IP
   */
  const getSNDSMetrics = async (ip: string): Promise<SNDSMetrics> => {
    const response = await fetch(`${API_BASE}/reputation/external/snds/${ip}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch SNDS metrics')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Get external metrics trends
   */
  const getExternalMetricsTrends = async (
    domain: string,
    days = 7
  ): Promise<ExternalMetricsTrends> => {
    const response = await fetch(
      `${API_BASE}/reputation/external/trends?domain=${domain}&days=${days}`,
      {
        method: 'GET',
        headers: getAuthHeaders()
      }
    )

    if (!response.ok) {
      throw new Error('Failed to fetch external metrics trends')
    }

    const data = await response.json()
    return data.data
  }

  // --------------------------------------------------------------------------
  // Provider Rate Limits
  // --------------------------------------------------------------------------
  
  /**
   * List all provider rate limits
   */
  const listProviderRateLimits = async (): Promise<ProviderRateLimit[]> => {
    const response = await fetch(`${API_BASE}/reputation/provider-limits`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch provider rate limits')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Update provider rate limit
   */
  const updateProviderRateLimit = async (
    id: number,
    limit: Partial<ProviderRateLimit>
  ): Promise<ProviderRateLimit> => {
    const response = await fetch(`${API_BASE}/reputation/provider-limits/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(limit)
    })

    if (!response.ok) {
      throw new Error('Failed to update provider rate limit')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Initialize provider limits for a domain
   */
  const initializeProviderLimits = async (domain: string): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/reputation/provider-limits/init/${domain}`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to initialize provider limits')
    }

    return await response.json()
  }

  /**
   * Reset provider usage
   */
  const resetProviderUsage = async (id: number): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/reputation/provider-limits/${id}/reset`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to reset provider usage')
    }

    return await response.json()
  }

  // --------------------------------------------------------------------------
  // Custom Warmup
  // --------------------------------------------------------------------------
  
  /**
   * Get custom warmup schedule for a domain
   */
  const getCustomWarmupSchedule = async (domain: string): Promise<WarmupSchedule | null> => {
    const response = await fetch(`${API_BASE}/reputation/warmup/${domain}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch custom warmup schedule')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Create custom warmup schedule
   */
  const createCustomWarmupSchedule = async (
    schedule: Omit<WarmupSchedule, 'id' | 'created_at' | 'current_day' | 'progress_percentage' | 'completed'>
  ): Promise<WarmupSchedule> => {
    const response = await fetch(`${API_BASE}/reputation/warmup`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(schedule)
    })

    if (!response.ok) {
      throw new Error('Failed to create custom warmup schedule')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Update custom warmup schedule
   */
  const updateCustomWarmupSchedule = async (
    id: number,
    schedule: Partial<WarmupSchedule>
  ): Promise<WarmupSchedule> => {
    const response = await fetch(`${API_BASE}/reputation/warmup/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(schedule)
    })

    if (!response.ok) {
      throw new Error('Failed to update custom warmup schedule')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Delete custom warmup schedule
   */
  const deleteCustomWarmupSchedule = async (id: number): Promise<void> => {
    const response = await fetch(`${API_BASE}/reputation/warmup/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to delete custom warmup schedule')
    }
  }

  /**
   * Get warmup templates
   */
  const getWarmupTemplates = async (): Promise<WarmupTemplate[]> => {
    const response = await fetch(`${API_BASE}/reputation/warmup/templates`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch warmup templates')
    }

    const data = await response.json()
    return data.data || []
  }

  // --------------------------------------------------------------------------
  // Predictions
  // --------------------------------------------------------------------------
  
  /**
   * Get latest predictions for all domains
   */
  const getLatestPredictions = async (): Promise<Prediction[]> => {
    const response = await fetch(`${API_BASE}/reputation/predictions/latest`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch latest predictions')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Get predictions for a specific domain
   */
  const getDomainPredictions = async (domain: string): Promise<Prediction> => {
    const response = await fetch(`${API_BASE}/reputation/predictions/${domain}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch domain predictions')
    }

    const data = await response.json()
    return data.data
  }

  /**
   * Generate new predictions for a domain
   */
  const generatePredictions = async (domain: string): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/reputation/predictions/generate/${domain}`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to generate predictions')
    }

    return await response.json()
  }

  /**
   * Get prediction history for a domain
   */
  const getPredictionHistory = async (
    domain: string
  ): Promise<PredictionHistory[]> => {
    const response = await fetch(`${API_BASE}/reputation/predictions/${domain}/history`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch prediction history')
    }

    const data = await response.json()
    return data.data || []
  }

  // --------------------------------------------------------------------------
  // Phase 5 Alerts
  // --------------------------------------------------------------------------
  
  /**
   * List Phase 5 alerts
   */
  const listPhase5Alerts = async (): Promise<Phase5Alert[]> => {
    const response = await fetch(`${API_BASE}/reputation/alerts/phase5`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch Phase 5 alerts')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Acknowledge alert
   */
  const acknowledgeAlert = async (id: number): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/reputation/alerts/${id}/acknowledge`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to acknowledge alert')
    }

    return await response.json()
  }

  /**
   * Resolve alert
   */
  const resolveAlert = async (id: number): Promise<{ message: string }> => {
    const response = await fetch(`${API_BASE}/reputation/alerts/${id}/resolve`, {
      method: 'POST',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to resolve alert')
    }

    return await response.json()
  }

  return {
    // DMARC Reports
    listDMARCReports,
    getDMARCReport,
    getDMARCStats,
    getDMARCActions,
    exportDMARCReport,
    // ARF Reports
    listARFReports,
    getARFStats,
    processARFReport,
    // External Metrics
    getPostmasterMetrics,
    getSNDSMetrics,
    getExternalMetricsTrends,
    // Provider Rate Limits
    listProviderRateLimits,
    updateProviderRateLimit,
    initializeProviderLimits,
    resetProviderUsage,
    // Custom Warmup
    getCustomWarmupSchedule,
    createCustomWarmupSchedule,
    updateCustomWarmupSchedule,
    deleteCustomWarmupSchedule,
    getWarmupTemplates,
    // Predictions
    getLatestPredictions,
    getDomainPredictions,
    generatePredictions,
    getPredictionHistory,
    // Phase 5 Alerts
    listPhase5Alerts,
    acknowledgeAlert,
    resolveAlert
  }
}
