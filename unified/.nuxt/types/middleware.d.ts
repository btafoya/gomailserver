import type { NavigationGuard } from 'vue-router'
export type MiddlewareKey = "auth" | "setup-redirect"
declare module 'nuxt/app' {
  interface PageMeta {
    middleware?: MiddlewareKey | NavigationGuard | Array<MiddlewareKey | NavigationGuard>
  }
}