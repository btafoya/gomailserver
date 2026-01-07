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
        <h2 class="text-3xl font-bold tracking-tight">Domains</h2>
        <UButton @click="goToCreate">
          <Plus class="mr-2 h-4 w-4" />
          Add Domain
        </UButton>
      </div>

      <div v-if="loading" class="text-center py-12">
        <p class="text-muted-foreground">Loading domains...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading domains: {{ error }}
      </div>

      <div v-else-if="domains.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No domains found. Add your first domain to get started.</p>
      </div>

      <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
        <table class="w-full">
          <thead class="bg-muted/50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Domain</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Users</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Aliases</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Created</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Status</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="domain in domains" :key="domain.id" class="border-t border-border">
              <td class="px-6 py-4 text-sm text-foreground">{{ domain.name }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ domain.user_count || 0 }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ domain.alias_count || 0 }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ new Date(domain.created_at).toLocaleDateString() }}</td>
              <td class="px-6 py-4 text-sm">
                <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium bg-green-100 text-green-800">
                  Active
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <button @click="goToEdit(domain.id)" class="text-primary hover:underline mr-2">
                  Edit
                </button>
                <button @click="handleDelete(domain.id)" class="text-red-600 hover:underline">
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
import { ref, computed, onMounted } from 'vue'
import { Plus } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useDomainsApi } from '~/composables/api/domains'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { getDomains, deleteDomain } = useDomainsApi()

const logout = () => {
  authStore.logout()
}

const domains = ref([])
const loading = ref(false)
const error = ref(null)
const domainCount = computed(() => domains.value.length)

const loadDomains = async () => {
  loading.value = true
  error.value = null
  try {
    domains.value = await getDomains()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id) => {
  if (!confirm('Are you sure you want to delete this domain? This action cannot be undone.')) {
    return
  }

  try {
    await deleteDomain(id)
    domains.value = domains.value.filter(d => d.id !== id)
  } catch (err) {
    error.value = err.message
  }
}

const goToCreate = () => {
  navigateTo('/admin/domains/create')
}

const goToEdit = (id) => {
  navigateTo(`/admin/domains/${id}`)
}

onMounted(() => {
  loadDomains()
})
</script>
