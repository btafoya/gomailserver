<template>
  <div>
    <div class="border-b">
      <div class="flex h-16 items-center px-4">
        <h1 class="text-lg font-semibold">Admin Console</h1>
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

    <div class="flex-1 p-4 md:p-8">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-3xl font-bold tracking-tight">Aliases</h2>
        <UButton @click="goToCreate">
          <Plus class="mr-2 h-4 w-4" />
          Add Alias
        </UButton>
      </div>

      <div v-if="loading" class="text-center py-12">
        <p class="text-muted-foreground">Loading aliases...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading aliases: {{ error }}
      </div>

      <div v-else-if="aliases.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No aliases found. Add your first alias to get started.</p>
      </div>

      <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
        <table class="w-full">
          <thead class="bg-muted/50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Alias Address</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Destinations</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Domain</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Created</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Status</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="alias in aliases" :key="alias.id" class="border-t border-border">
              <td class="px-6 py-4 text-sm text-foreground">{{ alias.address }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ alias.destinations.join(', ') }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ alias.domain_name || alias.domain_id }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ new Date(alias.created_at).toLocaleDateString() }}</td>
              <td class="px-6 py-4 text-sm">
                <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium" :class="getStatusClass(alias.status)">
                  {{ alias.status }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <button @click="goToEdit(alias.id)" class="text-primary hover:underline mr-2">
                  Edit
                </button>
                <button @click="handleDelete(alias.id, alias.address)" class="text-red-600 hover:underline">
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useAliasesApi } from '~/composables/api/aliases'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { getAliases, deleteAlias } = useAliasesApi()

const logout = () => {
  authStore.logout()
}

const aliases = ref([])
const loading = ref(false)
const error = ref(null)

const loadAliases = async () => {
  loading.value = true
  error.value = null
  try {
    aliases.value = await getAliases()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id, address) => {
  if (!confirm(`Are you sure you want to delete alias "${address}"? This action cannot be undone.`)) {
    return
  }

  try {
    await deleteAlias(id)
    aliases.value = aliases.value.filter(a => a.id !== id)
  } catch (err) {
    error.value = err.message
  }
}

const goToCreate = () => {
  navigateTo('/admin/aliases/create')
}

const goToEdit = (id) => {
  navigateTo(`/admin/aliases/${id}`)
}

const getStatusClass = (status) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'inactive':
      return 'bg-gray-100 text-gray-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

onMounted(() => {
  loadAliases()
})
</script>
