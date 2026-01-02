// Service Worker for Webmail PWA
const CACHE_NAME = 'webmail-v1'
const RUNTIME_CACHE = 'webmail-runtime-v1'

// Assets to cache on install
const STATIC_ASSETS = [
  '/webmail/',
  '/webmail/login',
  '/webmail/_nuxt/entry.js',
  '/webmail/_nuxt/entry.css'
]

// Install event - cache static assets
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(STATIC_ASSETS)
    })
  )
  self.skipWaiting()
})

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME && name !== RUNTIME_CACHE)
          .map((name) => caches.delete(name))
      )
    })
  )
  self.clients.claim()
})

// Fetch event - network-first strategy for API, cache-first for assets
self.addEventListener('fetch', (event) => {
  const { request } = event
  const url = new URL(request.url)

  // Skip non-GET requests
  if (request.method !== 'GET') {
    return
  }

  // API requests - network-first with cache fallback
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(
      fetch(request)
        .then((response) => {
          // Clone and cache successful responses
          if (response.ok) {
            const responseClone = response.clone()
            caches.open(RUNTIME_CACHE).then((cache) => {
              cache.put(request, responseClone)
            })
          }
          return response
        })
        .catch(() => {
          // Fallback to cache if network fails
          return caches.match(request)
        })
    )
    return
  }

  // Static assets - cache-first strategy
  event.respondWith(
    caches.match(request).then((cachedResponse) => {
      if (cachedResponse) {
        return cachedResponse
      }

      return fetch(request).then((response) => {
        // Don't cache non-successful responses
        if (!response || response.status !== 200 || response.type === 'error') {
          return response
        }

        // Clone and cache for future requests
        const responseClone = response.clone()
        caches.open(RUNTIME_CACHE).then((cache) => {
          cache.put(request, responseClone)
        })

        return response
      })
    })
  )
})

// Background sync for sending emails when offline
self.addEventListener('sync', (event) => {
  if (event.tag === 'send-email') {
    event.waitUntil(syncSendEmail())
  }
})

async function syncSendEmail() {
  // Get pending emails from IndexedDB
  // Send them when back online
  // This would need IndexedDB integration
  console.log('Syncing pending emails...')
}

// Push notifications for new emails
self.addEventListener('push', (event) => {
  const data = event.data?.json() || {}
  const title = data.title || 'New Email'
  const options = {
    body: data.body || 'You have a new message',
    icon: '/webmail/icon-192.png',
    badge: '/webmail/badge-72.png',
    tag: 'new-email',
    requireInteraction: false,
    data: {
      url: data.url || '/webmail/'
    }
  }

  event.waitUntil(self.registration.showNotification(title, options))
})

// Notification click - open webmail
self.addEventListener('notificationclick', (event) => {
  event.notification.close()

  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clientList) => {
      // If webmail is already open, focus it
      for (const client of clientList) {
        if (client.url.includes('/webmail') && 'focus' in client) {
          return client.focus()
        }
      }
      // Otherwise open new window
      if (clients.openWindow) {
        return clients.openWindow(event.notification.data.url)
      }
    })
  )
})
