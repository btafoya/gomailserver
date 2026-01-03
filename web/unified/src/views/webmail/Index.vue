<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

onMounted(() => {
  authStore.initializeAuth()
  if (authStore.isAuthenticated) {
    router.push('/webmail/mail/inbox')
  } else {
    router.push('/webmail/login')
  }
})
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <div class="text-center">
      <Icon name="lucide:loader-circle" class="w-8 h-8 animate-spin mx-auto text-primary" />
      <p class="mt-4 text-muted-foreground">Loading...</p>
    </div>
  </div>
</template>
