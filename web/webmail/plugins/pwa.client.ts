/**
 * PWA plugin for service worker registration and offline support
 */

export default defineNuxtPlugin(() => {
  if (process.client && 'serviceWorker' in navigator) {
    window.addEventListener('load', async () => {
      try {
        const registration = await navigator.serviceWorker.register('/webmail/sw.js', {
          scope: '/webmail/'
        })

        console.log('Service Worker registered:', registration.scope)

        // Handle updates
        registration.addEventListener('updatefound', () => {
          const newWorker = registration.installing
          if (newWorker) {
            newWorker.addEventListener('statechange', () => {
              if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                // New service worker available
                if (confirm('New version available! Reload to update?')) {
                  window.location.reload()
                }
              }
            })
          }
        })

        // Check for updates periodically
        setInterval(() => {
          registration.update()
        }, 60 * 60 * 1000) // Check every hour

      } catch (error) {
        console.error('Service Worker registration failed:', error)
      }
    })

    // Listen for offline/online events
    window.addEventListener('offline', () => {
      console.log('App is offline')
      // Could show offline indicator in UI
    })

    window.addEventListener('online', () => {
      console.log('App is online')
      // Could trigger sync or refresh
    })
  }
})
