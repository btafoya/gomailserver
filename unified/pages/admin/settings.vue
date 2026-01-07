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
      <h2 class="text-3xl font-bold tracking-tight mb-6">System Settings</h2>

      <div class="max-w-2xl space-y-6">
        <UCard>
          <UCardHeader>
            <UCardTitle>Server Configuration</UCardTitle>
          </UCardHeader>
          <UCardContent>
            <div class="space-y-4">
              <div>
                <label class="text-sm font-medium">Server Hostname</label>
                <input
                  type="text"
                  value="mail.example.com"
                  disabled
                  class="w-full px-3 py-2 border rounded-md bg-muted text-foreground"
                />
              </div>
              <div>
                <label class="text-sm font-medium">Domain</label>
                <input
                  type="text"
                  value="example.com"
                  disabled
                  class="w-full px-3 py-2 border rounded-md bg-muted text-foreground"
                />
              </div>
            </div>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader>
            <UCardTitle>Security Settings</UCardTitle>
          </UCardHeader>
          <UCardContent>
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm">DKIM Signing</span>
                <span class="text-sm font-medium text-green-600">Enabled</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm">SPF Checking</span>
                <span class="text-sm font-medium text-green-600">Enabled</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm">DMARC Enforcement</span>
                <span class="text-sm font-medium text-green-600">Enabled</span>
              </div>
            </div>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader>
            <UCardTitle>Storage Settings</UCardTitle>
          </UCardHeader>
          <UCardContent>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-sm">Database Path</span>
                <span class="text-sm font-mono">./data/mailserver.db</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm">Total Size</span>
                <span class="text-sm font-medium">245.8 MB</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm">WAL Mode</span>
                <span class="text-sm font-medium text-green-600">Enabled</span>
              </div>
            </div>
          </UCardContent>
        </UCard>

        <div class="bg-card rounded-lg border border-border p-6">
          <h3 class="text-lg font-semibold mb-4">Quick Actions</h3>
          <div class="grid gap-3 md:grid-cols-2">
            <UButton>
              <RefreshCw class="mr-2 h-4 w-4" />
              Reload Configuration
            </UButton>
            <UButton variant="destructive">
              <Power class="mr-2 h-4 w-4" />
              Restart Server
            </UButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RefreshCw, Power } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'

// Page metadata
definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}
</script>
