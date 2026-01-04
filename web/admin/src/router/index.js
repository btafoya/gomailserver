import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('@/views/setup/Index.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue')
      },
      {
        path: 'domains',
        name: 'Domains',
        component: () => import('@/views/domains/List.vue')
      },
      {
        path: 'domains/create',
        name: 'DomainCreate',
        component: () => import('@/views/domains/Create.vue')
      },
      {
        path: 'domains/:id',
        name: 'DomainEdit',
        component: () => import('@/views/domains/Edit.vue')
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/users/List.vue')
      },
      {
        path: 'users/create',
        name: 'UserCreate',
        component: () => import('@/views/users/Create.vue')
      },
      {
        path: 'users/:id',
        name: 'UserEdit',
        component: () => import('@/views/users/Edit.vue')
      },
      {
        path: 'aliases',
        name: 'Aliases',
        component: () => import('@/views/aliases/List.vue')
      },
      {
        path: 'queue',
        name: 'Queue',
        component: () => import('@/views/Queue.vue')
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('@/views/Logs.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/Settings.vue')
      },
      {
        path: 'audit',
        name: 'Audit',
        component: () => import('@/views/Audit.vue')
      },
      {
        path: 'reputation',
        name: 'Reputation',
        component: () => import('@/views/reputation/Overview.vue')
      },
      {
        path: 'reputation/circuit-breakers',
        name: 'ReputationCircuitBreakers',
        component: () => import('@/views/reputation/CircuitBreakers.vue')
      },
      {
        path: 'reputation/warmup',
        name: 'ReputationWarmup',
        component: () => import('@/views/reputation/Warmup.vue')
      },
      {
        path: 'reputation/audit',
        name: 'ReputationAudit',
        component: () => import('@/views/reputation/Audit.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  // Check authentication directly from localStorage to avoid circular dependency
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)

  if (requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.name === 'Login' && token) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
