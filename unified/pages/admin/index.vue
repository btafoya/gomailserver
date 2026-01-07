<template>
  <div class="content-wrapper">
    <!-- Dashboard Title -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold tracking-tight" style="font-family: var(--font-heading)">
        Admin Dashboard
      </h1>
    </div>

    <!-- Stats Cards -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
      <UCard class="hover-effect cursor-pointer">
        <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UCardTitle class="text-sm font-medium">Domains</UCardTitle>
          <Globe class="h-4 w-4 text-muted-foreground" />
        </UCardHeader>
        <UCardContent>
          <div class="text-2xl font-bold" style="font-family: var(--font-heading)">5</div>
          <p class="text-xs text-muted-foreground">Active domains</p>
        </UCardContent>
      </UCard>

      <UCard class="hover-effect cursor-pointer">
        <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UCardTitle class="text-sm font-medium">Users</UCardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </UCardHeader>
        <UCardContent>
          <div class="text-2xl font-bold" style="font-family: var(--font-heading)">23</div>
          <p class="text-xs text-muted-foreground">Email accounts</p>
        </UCardContent>
      </UCard>

      <UCard class="hover-effect cursor-pointer">
        <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UCardTitle class="text-sm font-medium">Queue</UCardTitle>
          <Mail class="h-4 w-4 text-muted-foreground" />
        </UCardHeader>
        <UCardContent>
          <div class="text-2xl font-bold" style="font-family: var(--font-heading)">12</div>
          <p class="text-xs text-muted-foreground">Messages in queue</p>
        </UCardContent>
      </UCard>

      <UCard class="hover-effect cursor-pointer">
        <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UCardTitle class="text-sm font-medium">Reputation</UCardTitle>
          <ShieldCheck class="h-4 w-4 text-muted-foreground" />
        </UCardHeader>
        <UCardContent>
          <div class="text-2xl font-bold" style="font-family: var(--font-heading)">94.2</div>
          <p class="text-xs text-muted-foreground">Avg sender score</p>
        </UCardContent>
      </UCard>
    </div>

    <!-- Action Cards -->
    <div class="grid gap-4 md:grid-cols-2 mb-6">
      <!-- Quick Actions -->
      <UCard>
        <UCardHeader>
          <UCardTitle class="text-lg font-semibold" style="font-family: var(--font-heading)">Quick Actions</UCardTitle>
        </UCardHeader>
        <UCardContent class="space-y-3">
          <NuxtLink to="/admin/domains">
            <UButton variant="outline" class="w-full justify-start hover-effect focus-ring">
              <Globe class="mr-2 h-4 w-4" />
              Manage Domains
            </UButton>
          </NuxtLink>
          <NuxtLink to="/admin/users">
            <UButton variant="outline" class="w-full justify-start hover-effect focus-ring">
              <Users class="mr-2 h-4 w-4" />
              Manage Users
            </UButton>
          </NuxtLink>
          <NuxtLink to="/admin/queue">
            <UButton variant="outline" class="w-full justify-start hover-effect focus-ring">
              <Mail class="mr-2 h-4 w-4" />
              View Queue
            </UButton>
          </NuxtLink>
          <NuxtLink to="/admin/settings">
            <UButton variant="outline" class="w-full justify-start hover-effect focus-ring">
              <Settings class="mr-2 h-4 w-4" />
              System Settings
            </UButton>
          </NuxtLink>
        </UCardContent>
      </UCard>

      <!-- System Status -->
      <UCard>
        <UCardHeader>
          <UCardTitle class="text-lg font-semibold" style="font-family: var(--font-heading)">System Status</UCardTitle>
        </UCardHeader>
        <UCardContent class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">SMTP</span>
            <span class="text-sm font-semibold text-green-600">Running</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">IMAP</span>
            <span class="text-sm font-semibold text-green-600">Running</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">ClamAV</span>
            <span class="text-sm font-semibold text-green-600">Connected</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">SpamAssassin</span>
            <span class="text-sm font-semibold text-green-600">Connected</span>
          </div>
        </UCardContent>
      </UCard>
    </div>

    <!-- Recent Activity (Example) -->
    <UCard>
      <UCardHeader>
        <UCardTitle class="text-lg font-semibold" style="font-family: var(--font-heading)">Recent Activity</UCardTitle>
      </UCardHeader>
      <UCardContent>
        <div class="space-y-4">
          <div class="flex items-start space-x-3 p-3 border border-muted rounded-lg hover-effect cursor-pointer">
            <div class="h-2 w-2 rounded-full bg-primary mt-1 flex-shrink-0"></div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-foreground">New domain added</p>
              <p class="text-xs text-muted-foreground">example.com - 2 minutes ago</p>
            </div>
          </div>
          <div class="flex items-start space-x-3 p-3 border border-muted rounded-lg hover-effect cursor-pointer">
            <div class="h-2 w-2 rounded-full bg-secondary mt-1 flex-shrink-0"></div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-foreground">User created</p>
              <p class="text-xs text-muted-foreground">john@example.com - 5 minutes ago</p>
            </div>
          </div>
          <div class="flex items-start space-x-3 p-3 border border-muted rounded-lg hover-effect cursor-pointer">
            <div class="h-2 w-2 rounded-full bg-accent mt-1 flex-shrink-0"></div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-foreground">Queue processed</p>
              <p class="text-xs text-muted-foreground">123 messages - 10 minutes ago</p>
            </div>
          </div>
        </div>
      </UCardContent>
    </UCard>
  </div>
</template>

<script setup>
import { Globe, Users, Mail, Settings, ShieldCheck } from 'lucide-vue-next'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})
</script>
