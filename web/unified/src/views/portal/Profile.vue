<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api/axios'

const authStore = useAuthStore()
const user = ref(null)
const editing = ref(false)
const loading = ref(false)
const error = ref('')
const success = ref('')

const formData = ref({
  name: '',
  email: ''
})

onMounted(async () => {
  await loadProfile()
})

const loadProfile = async () => {
  loading.value = true
  error.value = ''

  try {
    user.value = authStore.user
    if (user.value) {
      formData.value = {
        name: user.value.name || '',
        email: user.value.email || ''
      }
    }
  } catch (err) {
    error.value = 'Failed to load profile'
  } finally {
    loading.value = false
  }
}

const startEdit = () => {
  editing.value = true
  error.value = ''
  success.value = ''
}

const cancelEdit = () => {
  editing.value = false
  if (user.value) {
    formData.value = {
      name: user.value.name || '',
      email: user.value.email || ''
    }
  }
}

const saveProfile = async () => {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const { data } = await api.put(`/v1/users/${user.value.id}`, formData.value)
    user.value = data
    authStore.user = data
    success.value = 'Profile updated successfully'
    editing.value = false
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to update profile'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container mx-auto py-8 px-4 max-w-2xl">
    <div class="mb-6">
      <h1 class="text-3xl font-bold mb-2">Profile</h1>
      <p class="text-muted-foreground">Manage your account information</p>
    </div>

    <div v-if="error" class="bg-destructive/15 text-destructive p-4 rounded-md mb-4">
      {{ error }}
    </div>

    <div v-if="success" class="bg-green-500/15 text-green-700 p-4 rounded-md mb-4">
      {{ success }}
    </div>

    <div v-if="loading && !user" class="text-center py-12">
      <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto"></div>
      <p class="mt-4 text-muted-foreground">Loading profile...</p>
    </div>

    <div v-else class="bg-card border rounded-lg p-6">
      <div class="space-y-6">
        <!-- Name -->
        <div>
          <label class="block text-sm font-medium mb-2">Name</label>
          <input
            v-if="editing"
            v-model="formData.name"
            type="text"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="Your name"
          />
          <p v-else class="text-lg">{{ user?.name || 'Not set' }}</p>
        </div>

        <!-- Email -->
        <div>
          <label class="block text-sm font-medium mb-2">Email</label>
          <input
            v-if="editing"
            v-model="formData.email"
            type="email"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="you@example.com"
            disabled
          />
          <p v-else class="text-lg">{{ user?.email }}</p>
          <p v-if="editing" class="text-sm text-muted-foreground mt-1">Email cannot be changed</p>
        </div>

        <!-- Domain (Read-only) -->
        <div>
          <label class="block text-sm font-medium mb-2">Domain</label>
          <p class="text-lg">{{ user?.domain }}</p>
        </div>

        <!-- Created -->
        <div>
          <label class="block text-sm font-medium mb-2">Account Created</label>
          <p class="text-lg">{{ user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'Unknown' }}</p>
        </div>

        <!-- Actions -->
        <div class="flex gap-3 pt-4 border-t">
          <button
            v-if="!editing"
            @click="startEdit"
            class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
          >
            Edit Profile
          </button>
          <template v-else>
            <button
              @click="saveProfile"
              :disabled="loading"
              class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
            >
              {{ loading ? 'Saving...' : 'Save Changes' }}
            </button>
            <button
              @click="cancelEdit"
              :disabled="loading"
              class="px-4 py-2 border rounded-md hover:bg-accent disabled:opacity-50"
            >
              Cancel
            </button>
          </template>
        </div>
      </div>
    </div>

    <!-- Quick Links -->
    <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
      <router-link
        to="/portal/password"
        class="p-4 border rounded-lg hover:shadow-md transition-shadow"
      >
        <h3 class="font-semibold mb-1">Change Password</h3>
        <p class="text-sm text-muted-foreground">Update your account password</p>
      </router-link>

      <router-link
        to="/webmail"
        class="p-4 border rounded-lg hover:shadow-md transition-shadow"
      >
        <h3 class="font-semibold mb-1">Webmail</h3>
        <p class="text-sm text-muted-foreground">Access your email</p>
      </router-link>
    </div>
  </div>
</template>
