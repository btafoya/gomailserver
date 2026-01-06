export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return navigateTo('/login')
  }

  // Redirect to portal if logged in and trying to access login
  if (to.path === '/login' && authStore.isLoggedIn) {
    return navigateTo('/portal')
  }

  // If accessing root admin path and logged in, redirect to portal
  if (to.path === '/' && authStore.isLoggedIn) {
    return navigateTo('/portal')
  }
})