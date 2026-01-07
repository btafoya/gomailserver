<template>
  <div class="flex-1 p-4 md:p-8">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-3xl font-bold tracking-tight">Users</h2>
      <UButton>
        <Plus class="mr-2 h-4 w-4" />
        Create User
      </UButton>
    </div>

    <div v-if="loading" class="text-center py-12">
      <p class="text-muted-foreground">Loading users...</p>
    </div>

    <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
      Error loading users: {{ error }}
    </div>

    <div v-else-if="users.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
      <p class="text-muted-foreground">No users found. Create your first user to get started.</p>
    </div>

    <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
      <table class="w-full">
        <thead class="bg-muted/50">
          <tr>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Email</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Name</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Domain</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Created</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="border-t border-border">
            <td class="px-6 py-4 text-sm text-foreground">{{ user.email }}</td>
            <td class="px-6 py-4 text-sm text-foreground">{{ user.full_name || user.email }}</td>
            <td class="px-6 py-4 text-sm text-muted-foreground">{{ user.domain }}</td>
            <td class="px-6 py-4 text-sm text-muted-foreground">{{ new Date(user.created_at).toLocaleDateString() }}</td>
            <td class="px-6 py-4 text-sm">
              <NuxtLink :to="`/admin/users/${user.id}`" class="text-primary hover:underline">
                Edit
              </NuxtLink>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Plus } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

// TODO: Replace with actual API call once backend is configured
const users = ref([
  { id: 1, email: 'admin@example.com', full_name: 'Admin User', domain: 'example.com', created_at: new Date().toISOString() },
  { id: 2, email: 'user1@example.com', full_name: 'User One', domain: 'example.com', created_at: new Date().toISOString() },
  { id: 3, email: 'user2@example.com', full_name: 'User Two', domain: 'example.com', created_at: new Date().toISOString() }
])

const loading = ref(false)
const error = ref(null)

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}
</script>
