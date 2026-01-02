import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Index',
    component: () => import('../pages/index.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/login.vue'),
  },
  {
    path: '/mail',
    name: 'Mail',
    component: () => import('../pages/mail/[mailboxId].vue'),
    children: [
      {
        path: 'compose',
        name: 'Compose',
        component: () => import('../pages/mail/compose.vue'),
      },
      {
        path: ':mailboxId',
        name: 'Mailbox',
        component: () => import('../pages/mail/[mailboxId].vue'),
        children: [
          {
            path: 'message/:messageId',
            name: 'Message',
            component: () => import('../pages/mail/[mailboxId]/message/[messageId].vue'),
          },
        ],
      },
    ],
  },
  {
    path: '/settings',
    name: 'Settings',
    children: [
      {
        path: 'pgp',
        name: 'PGP',
        component: () => import('../pages/settings/pgp.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory('/webmail/'),
  routes,
})

export default router
