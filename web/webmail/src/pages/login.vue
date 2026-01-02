<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  try {
    loading.value = true
    error.value = ''
    await authStore.login(email.value, password.value)
    router.push('/mail/inbox')
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Login failed. Please check your credentials.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-muted/30">
    <div class="w-full max-w-md">
      <div class="bg-card rounded-lg shadow-lg p-8">
        <div class="text-center mb-8">
          <Icon name="lucide:mail" class="w-12 h-12 mx-auto text-primary" />
          <h1 class="text-2xl font-bold mt-4">Webmail</h1>
          <p class="text-muted-foreground mt-2">Sign in to your account</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div v-if="error" class="bg-destructive/10 text-destructive text-sm p-3 rounded-md">
            {{ error }}
          </div>

          <div>
            <label for="email" class="block text-sm font-medium mb-2">Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              autocomplete="email"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="you@example.com"
            />
          </div>

          <div>
            <label for="password" class="block text-sm font-medium mb-2">Password</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="••••••••"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-primary text-primary-foreground py-2 px-4 rounded-md hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <Icon v-if="loading" name="lucide:loader-circle" class="w-4 h-4 animate-spin" />
            <span>{{ loading ? 'Signing in...' : 'Sign in' }}</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
