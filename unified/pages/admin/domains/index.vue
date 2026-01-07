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
        <UButton>
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
                <NuxtLink :to="`/admin/domains/${domain.id}`" class="text-primary hover:underline">
                  Edit
                </NuxtLink>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Plus } from 'lucide-vue-next'

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}

// TODO: Replace with actual API call once backend is configured
const domains = ref([
  { id: 1, name: 'example.com', user_count: 3, alias_count: 5, created_at: new Date().toISOString() },
  { id: 2, name: 'mail.example.com', user_count: 2, alias_count: 3, created_at: new Date(Date.now() - 86400000).toISOString() },
  { id: 3, name: 'news.example.com', user_count: 1, alias_count: 2, created_at: new Date(Date.now() - 172800000).toISOString() },
  { id: 4, name: 'support.example.com', user_count: 0, alias_count: 1, created_at: new Date(Date.now() - 259200000).toISOString() },
  { id: 5, name: 'dev.example.com', user_count: 0, alias_count: 0, created_at: new Date(Date.now() - 345600000).toISOString() }
])

const loading = ref(false)
const error = ref(null)
const domainCount = computed(() => domains.value.length)
</script>

<script>
export default {
  middleware: 'auth',
  layout: 'admin'
}
</script>
