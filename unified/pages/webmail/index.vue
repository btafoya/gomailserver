<template>
  <div>
    <div class="border-b">
      <div class="flex h-16 items-center px-4">
        <h1 class="text-lg font-semibold">Webmail</h1>
        <div class="ml-auto">
          <UButton
            variant="outline"
            @click="logout"
          >
            Logout
          </UButton>
        </div>
      </div>
    </div>

    <div class="flex-1 space-y-4 p-4 md:p-8">
      <div class="flex items-center justify-between">
        <h2 class="text-3xl font-bold tracking-tight">Inbox</h2>
        <UButton>
          <Mail class="mr-2 h-4 w-4" />
          Compose
        </UButton>
      </div>

      <div class="rounded-md border">
        <div class="p-4">
          <p class="text-muted-foreground">Webmail interface coming soon...</p>
          <p class="text-sm text-muted-foreground mt-2">
            This will include Gmail-like email management with conversation view, rich text composer, and contact integration.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Mail } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth',
  layout: 'webmail'
})

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}
</script>