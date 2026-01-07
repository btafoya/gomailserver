<template>
  <div>
    <div class="border-b">
      <div class="flex h-16 items-center px-4">
        <h1 class="text-lg font-semibold">User Portal</h1>
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

    <div class="flex-1 space-y-4 p-4 md:p-8">
      <div class="flex items-center justify-between">
        <h2 class="text-3xl font-bold tracking-tight">Dashboard</h2>
      </div>

      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Email Usage</UCardTitle>
            <Mail class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">1,234</div>
            <p class="text-xs text-muted-foreground">+20.1% from last month</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Storage Used</UCardTitle>
            <HardDrive class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">45.2 MB</div>
            <p class="text-xs text-muted-foreground">+2.5% from last month</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Active Aliases</UCardTitle>
            <AtSign class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">12</div>
            <p class="text-xs text-muted-foreground">+2 new this month</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Reputation Score</UCardTitle>
            <TrendingUp class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">98.5</div>
            <p class="text-xs text-muted-foreground">+0.3 from last month</p>
          </UCardContent>
        </UCard>
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <UCard>
          <UCardHeader>
            <UCardTitle>Quick Actions</UCardTitle>
          </UCardHeader>
          <UCardContent class="space-y-2">
            <NuxtLink to="/portal/profile">
              <UButton variant="outline" class="w-full justify-start">
                <User class="mr-2 h-4 w-4" />
                Update Profile
              </UButton>
            </NuxtLink>
            <NuxtLink to="/portal/password">
              <UButton variant="outline" class="w-full justify-start">
                <Lock class="mr-2 h-4 w-4" />
                Change Password
              </UButton>
            </NuxtLink>
            <NuxtLink to="/webmail">
              <UButton variant="outline" class="w-full justify-start">
                <Mail class="mr-2 h-4 w-4" />
                Open Webmail
              </UButton>
            </NuxtLink>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader>
            <UCardTitle>Recent Activity</UCardTitle>
          </UCardHeader>
          <UCardContent>
            <div class="space-y-4">
              <div class="flex items-center">
                <div class="ml-4 space-y-1">
                  <p class="text-sm font-medium">Password changed</p>
                  <p class="text-sm text-muted-foreground">2 hours ago</p>
                </div>
              </div>
              <div class="flex items-center">
                <div class="ml-4 space-y-1">
                  <p class="text-sm font-medium">Email sent to support</p>
                  <p class="text-sm text-muted-foreground">1 day ago</p>
                </div>
              </div>
            </div>
          </UCardContent>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Mail, HardDrive, AtSign, TrendingUp, User, Lock } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth',
  layout: 'portal'
})

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}
</script>