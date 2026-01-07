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
        <h2 class="text-3xl font-bold tracking-tight">Admin Dashboard</h2>
      </div>

      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Domains</UCardTitle>
            <Globe class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">5</div>
            <p class="text-xs text-muted-foreground">Active domains</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Users</UCardTitle>
            <Users class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">23</div>
            <p class="text-xs text-muted-foreground">Email accounts</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Queue</UCardTitle>
            <Mail class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">12</div>
            <p class="text-xs text-muted-foreground">Messages in queue</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Reputation</UCardTitle>
            <ShieldCheck class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">94.2</div>
            <p class="text-xs text-muted-foreground">Avg sender score</p>
          </UCardContent>
        </UCard>
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <UCard>
          <UCardHeader>
            <UCardTitle>Quick Actions</UCardTitle>
          </UCardHeader>
          <UCardContent class="space-y-2">
            <NuxtLink to="/admin/domains">
              <UButton variant="outline" class="w-full justify-start">
                <Globe class="mr-2 h-4 w-4" />
                Manage Domains
              </UButton>
            </NuxtLink>
            <NuxtLink to="/admin/users">
              <UButton variant="outline" class="w-full justify-start">
                <Users class="mr-2 h-4 w-4" />
                Manage Users
              </UButton>
            </NuxtLink>
            <NuxtLink to="/admin/queue">
              <UButton variant="outline" class="w-full justify-start">
                <Mail class="mr-2 h-4 w-4" />
                View Queue
              </UButton>
            </NuxtLink>
            <NuxtLink to="/admin/settings">
              <UButton variant="outline" class="w-full justify-start">
                <Settings class="mr-2 h-4 w-4" />
                System Settings
              </UButton>
            </NuxtLink>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader>
            <UCardTitle>System Status</UCardTitle>
          </UCardHeader>
          <UCardContent>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-sm">SMTP</span>
                <span class="text-sm font-medium text-green-600">Running</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm">IMAP</span>
                <span class="text-sm font-medium text-green-600">Running</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm">ClamAV</span>
                <span class="text-sm font-medium text-green-600">Connected</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm">SpamAssassin</span>
                <span class="text-sm font-medium text-green-600">Connected</span>
              </div>
            </div>
          </UCardContent>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Globe, Users, Mail, ShieldCheck, Settings } from 'lucide-vue-next'

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}
</script>

<script>
// Page metadata
export default {
  middleware: 'auth',
  layout: 'admin'
}
</script>
