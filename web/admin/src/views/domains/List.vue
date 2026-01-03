<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api/axios'

const domains = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await api.get('/v1/domains')
    domains.value = response.data.domains || []
  } catch (err) {
    console.error('Failed to load domains:', err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="p-8">
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-foreground">Domains</h1>
        <p class="text-muted-foreground mt-2">Manage your mail domains</p>
      </div>
      <router-link
        to="/domains/create"
        class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
      >
        Add Domain
      </router-link>
    </div>

    <div v-if="loading" class="text-center py-12">
      <p class="text-muted-foreground">Loading domains...</p>
    </div>

    <div v-else-if="domains.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
      <p class="text-muted-foreground">No domains found. Create your first domain to get started.</p>
    </div>

    <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
      <table class="w-full">
        <thead class="bg-muted/50">
          <tr>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Domain</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Users</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Created</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="domain in domains" :key="domain.id" class="border-t border-border">
            <td class="px-6 py-4 text-sm text-foreground">{{ domain.name }}</td>
            <td class="px-6 py-4 text-sm text-muted-foreground">{{ domain.user_count || 0 }}</td>
            <td class="px-6 py-4 text-sm text-muted-foreground">{{ new Date(domain.created_at).toLocaleDateString() }}</td>
            <td class="px-6 py-4 text-sm">
              <router-link :to="`/domains/${domain.id}`" class="text-primary hover:underline">
                Edit
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
</style>
