<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api/axios'

const router = useRouter()
const authStore = useAuthStore()

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

const resetForm = () => {
  currentPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  error.value = ''
}

const changePassword = async () => {
  error.value = ''
  success.value = ''

  // Validation
  if (!currentPassword.value || !newPassword.value || !confirmPassword.value) {
    error.value = 'All fields are required'
    return
  }

  if (newPassword.value.length < 8) {
    error.value = 'New password must be at least 8 characters'
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    error.value = 'New passwords do not match'
    return
  }

  if (currentPassword.value === newPassword.value) {
    error.value = 'New password must be different from current password'
    return
  }

  loading.value = true

  try {
    await api.post('/v1/auth/change-password', {
      current_password: currentPassword.value,
      new_password: newPassword.value
    })

    success.value = 'Password changed successfully'
    resetForm()

    // Redirect to profile after 2 seconds
    setTimeout(() => {
      router.push('/portal/profile')
    }, 2000)
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to change password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container mx-auto py-8 px-4 max-w-2xl">
    <div class="mb-6">
      <h1 class="text-3xl font-bold mb-2">Change Password</h1>
      <p class="text-muted-foreground">Update your account password</p>
    </div>

    <div v-if="error" class="bg-destructive/15 text-destructive p-4 rounded-md mb-4">
      {{ error }}
    </div>

    <div v-if="success" class="bg-green-500/15 text-green-700 p-4 rounded-md mb-4">
      {{ success }}
    </div>

    <div class="bg-card border rounded-lg p-6">
      <form @submit.prevent="changePassword" class="space-y-6">
        <!-- Current Password -->
        <div>
          <label for="current-password" class="block text-sm font-medium mb-2">
            Current Password
          </label>
          <input
            id="current-password"
            v-model="currentPassword"
            type="password"
            required
            autocomplete="current-password"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="Enter current password"
          />
        </div>

        <!-- New Password -->
        <div>
          <label for="new-password" class="block text-sm font-medium mb-2">
            New Password
          </label>
          <input
            id="new-password"
            v-model="newPassword"
            type="password"
            required
            autocomplete="new-password"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="Enter new password"
          />
          <p class="text-sm text-muted-foreground mt-1">
            Must be at least 8 characters
          </p>
        </div>

        <!-- Confirm Password -->
        <div>
          <label for="confirm-password" class="block text-sm font-medium mb-2">
            Confirm New Password
          </label>
          <input
            id="confirm-password"
            v-model="confirmPassword"
            type="password"
            required
            autocomplete="new-password"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="Confirm new password"
          />
        </div>

        <!-- Password Requirements -->
        <div class="bg-muted/30 p-4 rounded-md">
          <h3 class="font-medium mb-2 text-sm">Password Requirements:</h3>
          <ul class="text-sm text-muted-foreground space-y-1">
            <li class="flex items-center gap-2">
              <span :class="newPassword.length >= 8 ? 'text-green-600' : 'text-muted-foreground'">
                {{ newPassword.length >= 8 ? '✓' : '○' }}
              </span>
              At least 8 characters
            </li>
            <li class="flex items-center gap-2">
              <span :class="newPassword !== confirmPassword || !newPassword ? 'text-muted-foreground' : newPassword === confirmPassword ? 'text-green-600' : 'text-destructive'">
                {{ !newPassword || !confirmPassword ? '○' : newPassword === confirmPassword ? '✓' : '✗' }}
              </span>
              Passwords match
            </li>
            <li class="flex items-center gap-2">
              <span :class="!currentPassword || !newPassword ? 'text-muted-foreground' : currentPassword !== newPassword ? 'text-green-600' : 'text-destructive'">
                {{ !currentPassword || !newPassword ? '○' : currentPassword !== newPassword ? '✓' : '✗' }}
              </span>
              Different from current password
            </li>
          </ul>
        </div>

        <!-- Actions -->
        <div class="flex gap-3 pt-4 border-t">
          <button
            type="submit"
            :disabled="loading"
            class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50 flex items-center gap-2"
          >
            <div v-if="loading" class="w-4 h-4 border-2 border-primary-foreground border-t-transparent rounded-full animate-spin"></div>
            {{ loading ? 'Changing...' : 'Change Password' }}
          </button>
          <router-link
            to="/portal/profile"
            class="px-4 py-2 border rounded-md hover:bg-accent"
          >
            Cancel
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>
