/**
 * Setup Wizard Redirect Middleware
 * Redirects to /admin/setup if setup is not complete
 */

export default defineNuxtRouteMiddleware(async (to, from) => {
  // Skip if already going to setup
  if (to.path === '/admin/setup') {
    return
  }

  // Check if user is authenticated
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
  
  // Skip redirect if user is already logged in (setup complete or has token)
  if (token) {
    return
  }

  try {
    // Check setup status
    const API_BASE = process.env.NODE_ENV === 'development'
      ? '/api/v1'
      : process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8980/api/v1'

    const response = await fetch(`${API_BASE}/setup/status`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      console.error('Failed to check setup status:', response.status)
      return
    }

    const data = await response.json()
    
    // Redirect to setup if not complete
    if (!data.setup_complete) {
      return navigateTo('/admin/setup')
    }

    // Setup is complete, allow navigation to continue
  } catch (error: any) {
    console.error('Setup redirect error:', error)
    // Allow navigation to continue (don't block user)
  }
})
