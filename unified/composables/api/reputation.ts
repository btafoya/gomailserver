/**
 * Reputation Management API - Phase 1-4
 * Basic reputation management endpoints: Audit, Scores, Circuit Breakers, Alerts
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

export interface AuditRequest {
  sending_ip?: string
}

export interface CheckStatus {
  passed: boolean
  message: string
  details?: Record<string, any>
}

export interface AuditResponse {
  domain: string
  timestamp: number
  spf: CheckStatus
  dkim: CheckStatus
  dmarc: CheckStatus
  rdns: CheckStatus
  fcrdns: CheckStatus
  tls: CheckStatus
  mta_sts: CheckStatus
  postmaster_ok: boolean
  abuse_ok: boolean
  overall_score: number
  issues: string[]
}

export interface ScoreResponse {
  domain: string
  reputation_score: number
  complaint_rate: number
  bounce_rate: number
  delivery_rate: number
  circuit_breaker_active: boolean
  circuit_breaker_reason?: string
  warm_up_active: boolean
  warm_up_day?: number
  last_updated: number
}

export interface CircuitBreakerResponse {
  id: number
  domain: string
  trigger_type: string
  trigger_value: number
  threshold: number
  paused_at: number
  resumed_at?: number
  auto_resumed: boolean
  admin_notes?: string
}

export interface AlertResponse {
  id: number
  domain: string
  type: string
  severity: 'critical' | 'high' | 'medium' | 'low'
  message: string
  created_at: number
  acknowledged: boolean
  resolved: boolean
}

// ============================================================================
// API Functions
// ============================================================================

export const useReputationApi = () => {
  const API_BASE = useApiBase()

  // --------------------------------------------------------------------------
  // Domain Audit
  // --------------------------------------------------------------------------
  
  /**
   * Run domain deliverability audit
   */
  const auditDomain = async (domain: string, request?: AuditRequest): Promise<AuditResponse> => {
    const response = await fetch(`${API_BASE}/reputation/audit/${domain}`, {
      method: 'POST',
      headers: getAuthHeaders(),
      ...(request ? { body: JSON.stringify(request) } : {})
    })

    if (!response.ok) {
      throw new Error('Failed to audit domain')
    }

    const data = await response.json()
    return data
  }

  // --------------------------------------------------------------------------
  // Scores
  // --------------------------------------------------------------------------
  
  /**
   * Get reputation scores for all domains
   */
  const getScores = async (): Promise<ScoreResponse[]> => {
    const response = await fetch(\`\${API_BASE}/reputation/scores\`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch reputation scores')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Get reputation score for a specific domain
   */
  const getScore = async (domain: string): Promise<ScoreResponse> => {
    const response = await fetch(\`\${API_BASE}/reputation/scores/\${domain}\`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch domain score')
    }

    const data = await response.json()
    return data
  }

  // --------------------------------------------------------------------------
  // Circuit Breakers
  // --------------------------------------------------------------------------
  
  /**
   * List all circuit breaker events
   */
  const listCircuitBreakers = async (): Promise<CircuitBreakerResponse[]> => {
    const response = await fetch(\`\${API_BASE}/reputation/circuit-breakers\`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch circuit breakers')
    }

    const data = await response.json()
    return data.data || []
  }

  /**
   * Get circuit breaker history for a specific domain
   */
  const getCircuitBreakerHistory = async (domain: string): Promise<CircuitBreakerResponse[]> => {
    const response = await fetch(\`\${API_BASE}/reputation/circuit-breakers/\${domain}/history\`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch circuit breaker history')
    }

    const data = await response.json()
    return data.data || []
  }

  // --------------------------------------------------------------------------
  // Alerts
  // --------------------------------------------------------------------------
  
  /**
   * List recent alerts
   */
  const listAlerts = async (): Promise<AlertResponse[]> => {
    const response = await fetch(\`\${API_BASE}/reputation/alerts\`, {
      method: 'GET',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      throw new Error('Failed to fetch alerts')
    }

    const data = await response.json()
    return data.data || []
  }

  return {
    // Audit
    auditDomain,
    // Scores
    getScores,
    getScore,
    // Circuit Breakers
    listCircuitBreakers,
    getCircuitBreakerHistory,
    // Alerts
    listAlerts
  }
}
