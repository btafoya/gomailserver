<template>
  <div class="flex-1 p-4 md:p-8">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-3xl font-bold tracking-tight">Users</h2>
        <UButton @click="goToCreate">
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
                <button @click="goToEdit(user.id)" class="text-primary hover:underline mr-2">
                  Edit
                </button>
                <button @click="handleDelete(user.id, user.email)" class="text-red-600 hover:underline">
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from 'lucide-vue-next'
import { useUsersApi } from '~/composables/api/users'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const { getUsers, deleteUser: removeUser } = useUsersApi()

interface User {
  id: number
  email: string
  full_name: string
  display_name?: string
  domain_id: number
  domain_name?: string
  quota: number
  used_quota: number
  status: string
  forward_to?: string
  auto_reply_enabled: boolean
  auto_reply_subject?: string
  auto_reply_body?: string
  spam_threshold: number
  totp_enabled: boolean
  created_at: string
  last_login?: string
}

const logout = () => {
  localStorage.removeItem('token')
  navigateTo('/login')
}

const users = ref<User[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const loadUsers = async () => {
  loading.value = true
  error.value = null
  try {
    users.value = await getUsers()
  } catch (err: any) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id: number, email: string) => {
  if (!confirm(`Are you sure you want to delete user "${email}"? This action cannot be undone.`)) {
    return
  }

  try {
    await removeUser(id)
    users.value = users.value.filter(u => u.id !== id)
  } catch (err: any) {
    error.value = err.message
  }
}

const goToCreate = () => {
  navigateTo('/admin/users/create')
}

const goToEdit = (id: number) => {
  navigateTo(`/admin/users/${id}`)
}

onMounted(() => {
  loadUsers()
})
</script>
