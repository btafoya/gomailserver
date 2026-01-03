<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api/axios'

const stats = ref({
  domains: 0,
  users: 0,
  queueSize: 0,
  recentActivity: []
})

const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  try {
    const response = await api.get('/v1/stats/dashboard')
    stats.value = response.data
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="p-8">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-foreground">Dashboard</h1>
      <p class="text-muted-foreground mt-2">Overview of your mail server</p>
    </div>

    <div v-if="loading" class="text-center py-12">
      <p class="text-muted-foreground">Loading dashboard...</p>
    </div>

    <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
      Error loading dashboard: {{ error }}
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- Domains Card -->
      <div class="bg-card rounded-lg border border-border p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-muted-foreground">Domains</p>
            <p class="text-3xl font-bold text-foreground mt-2">{{ stats.domains }}</p>
          </div>
          <div class="text-primary">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Users Card -->
      <div class="bg-card rounded-lg border border-border p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-muted-foreground">Users</p>
            <p class="text-3xl font-bold text-foreground mt-2">{{ stats.users }}</p>
          </div>
          <div class="text-primary">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Queue Card -->
      <div class="bg-card rounded-lg border border-border p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-muted-foreground">Queue Size</p>
            <p class="text-3xl font-bold text-foreground mt-2">{{ stats.queueSize }}</p>
          </div>
          <div class="text-primary">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-8 bg-card rounded-lg border border-border p-6">
      <h2 class="text-xl font-semibold text-foreground mb-4">Quick Links</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <router-link
          to="/domains"
          class="p-4 rounded-lg border border-border hover:bg-accent hover:text-accent-foreground transition-colors"
        >
          <h3 class="font-medium">Manage Domains</h3>
          <p class="text-sm text-muted-foreground mt-1">Add, edit, or remove domains</p>
        </router-link>
        <router-link
          to="/users"
          class="p-4 rounded-lg border border-border hover:bg-accent hover:text-accent-foreground transition-colors"
        >
          <h3 class="font-medium">Manage Users</h3>
          <p class="text-sm text-muted-foreground mt-1">Create and manage email accounts</p>
        </router-link>
        <router-link
          to="/queue"
          class="p-4 rounded-lg border border-border hover:bg-accent hover:text-accent-foreground transition-colors"
        >
          <h3 class="font-medium">View Queue</h3>
          <p class="text-sm text-muted-foreground mt-1">Monitor outgoing mail queue</p>
        </router-link>
        <router-link
          to="/logs"
          class="p-4 rounded-lg border border-border hover:bg-accent hover:text-accent-foreground transition-colors"
        >
          <h3 class="font-medium">View Logs</h3>
          <p class="text-sm text-muted-foreground mt-1">Check system logs and activity</p>
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
