import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  // Public routes
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/admin/Login.vue'),
    meta: { requiresAuth: false, module: 'admin' }
  },
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('@/views/admin/setup/Index.vue'),
    meta: { requiresAuth: false, module: 'admin' }
  },

  // Admin module routes
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true, module: 'admin' },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue')
      },
      {
        path: 'domains',
        name: 'AdminDomains',
        component: () => import('@/views/admin/domains/List.vue')
      },
      {
        path: 'domains/create',
        name: 'AdminDomainCreate',
        component: () => import('@/views/admin/domains/Create.vue')
      },
      {
        path: 'domains/:id',
        name: 'AdminDomainEdit',
        component: () => import('@/views/admin/domains/Edit.vue')
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/users/List.vue')
      },
      {
        path: 'users/create',
        name: 'AdminUserCreate',
        component: () => import('@/views/admin/users/Create.vue')
      },
      {
        path: 'users/:id',
        name: 'AdminUserEdit',
        component: () => import('@/views/admin/users/Edit.vue')
      },
      {
        path: 'aliases',
        name: 'AdminAliases',
        component: () => import('@/views/admin/aliases/List.vue')
      },
      {
        path: 'queue',
        name: 'AdminQueue',
        component: () => import('@/views/admin/Queue.vue')
      },
      {
        path: 'logs',
        name: 'AdminLogs',
        component: () => import('@/views/admin/Logs.vue')
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: () => import('@/views/admin/Settings.vue')
      },
      {
        path: 'audit',
        name: 'AdminAudit',
        component: () => import('@/views/admin/Audit.vue')
      }
    ]
  },

  // Webmail module routes
  {
    path: '/webmail',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: false, module: 'webmail' },
    children: [
      {
        path: '',
        name: 'WebmailIndex',
        component: () => import('@/views/webmail/Index.vue')
      },
      {
        path: 'login',
        name: 'WebmailLogin',
        component: () => import('@/views/webmail/Login.vue'),
        meta: { requiresAuth: false }
      },
      {
        path: 'mail/:mailboxId',
        name: 'WebmailMailbox',
        component: () => import('@/views/webmail/mail/Mailbox.vue'),
        children: [
          {
            path: 'message/:messageId',
            name: 'WebmailMessage',
            component: () => import('@/views/webmail/mail/message/MessageDetail.vue')
          }
        ]
      },
      {
        path: 'compose',
        name: 'WebmailCompose',
        component: () => import('@/views/webmail/mail/Compose.vue')
      },
      {
        path: 'settings/pgp',
        name: 'WebmailPGP',
        component: () => import('@/views/webmail/settings/Pgp.vue')
      }
    ]
  },

  // Portal module routes
  {
    path: '/portal',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true, module: 'portal' },
    children: [
      {
        path: '',
        name: 'PortalIndex',
        component: () => import('@/views/portal/Index.vue')
      },
      {
        path: 'profile',
        name: 'PortalProfile',
        component: () => import('@/views/portal/Profile.vue')
      },
      {
        path: 'password',
        name: 'PortalPassword',
        component: () => import('@/views/portal/PasswordReset.vue')
      }
    ]
  },

]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)

  if (requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.name === 'Login' && token) {
    next({ name: 'AdminDashboard' })
  } else {
    next()
  }
})

export default router
