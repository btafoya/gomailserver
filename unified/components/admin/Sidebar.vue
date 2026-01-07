<template>
  <div class="flex">
    <!-- Mobile Backdrop -->
    <div
      v-if="isOpen"
      class="backdrop-overlay lg:hidden"
      @click="closeMobile"
    />

    <!-- Sidebar -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-50 flex flex-col border-r transition-all duration-250',
        isOpen ? 'sidebar-container lg:translate-x-0' : 'sidebar-container lg:translate-x-[-100%]'
      ]"
    >
      <!-- Logo / Brand Area -->
      <div class="p-4 border-b bg-card">
        <div class="flex items-center space-x-3">
          <div class="h-10 w-10 bg-primary rounded-lg flex items-center justify-center shadow-md">
            <Mail class="h-6 w-6 text-primary-foreground" />
          </div>
          <h1 class="text-xl font-bold text-foreground">GoMail</h1>
        </div>
      </div>

      <!-- Navigation -->
      <div class="flex-1 overflow-y-auto py-4 px-3">
        <nav class="space-y-1">
          <NuxtLink
            v-for="item in navigationItems"
            :key="item.to"
            :to="item.to"
            :class="[
              'sidebar-nav-item flex items-center space-x-3 px-4 py-3 rounded-lg border border border-transparent hover-effect',
              isActive(item.to)
                ? 'sidebar-nav-item.active'
                : 'hover:bg-muted hover:border-muted'
            ]"
            @click="handleNavClick"
          >
            <component :is="item.icon" class="h-5 w-5 flex-shrink-0" />
            <span class="font-medium">{{ item.label }}</span>
          </NuxtLink>
        </nav>

        <!-- Quick Stats -->
        <div class="mt-6 pt-6 border-t">
          <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-4">
            Quick Stats
          </h3>
          <div class="space-y-3">
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">Domains</span>
              <span class="font-medium text-foreground">{{ stats.domains }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">Users</span>
              <span class="font-medium text-foreground">{{ stats.users }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">Messages</span>
              <span class="font-medium text-foreground">{{ stats.messages }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">Queue</span>
              <span class="font-medium text-foreground">{{ stats.queue }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- User Profile Section -->
      <div class="p-4 border-t bg-card">
        <div class="flex items-center space-x-3">
          <div class="h-8 w-8 rounded-full bg-primary flex items-center justify-center text-primary-foreground font-bold text-sm">
            AD
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-foreground truncate">Admin User</p>
            <p class="text-xs text-muted-foreground truncate">admin@localhost</p>
          </div>
          <button
            @click="handleLogout"
            class="p-2 rounded-lg hover:bg-muted hover-effect focus-ring"
            title="Logout"
          >
            <LogOut class="h-4 w-4 text-muted-foreground" />
          </button>
        </div>
      </div>
    </aside>

    <!-- Mobile Toggle Button -->
    <button
      @click="toggleMobile"
      class="fixed bottom-6 right-6 z-50 p-4 rounded-full bg-primary text-primary-foreground shadow-lg hover:scale-105 transition-transform lg:hidden hover-effect focus-ring"
      :aria-label="isOpen ? 'Close sidebar' : 'Open sidebar'"
    >
      <Menu v-if="!isOpen" class="h-6 w-6" />
      <X v-else class="h-6 w-6" />
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Home, Users, Globe, Mail, Settings, Menu, X, LogOut, ShieldCheck } from 'lucide-vue-next'
import { useRoute } from 'vue-router'

const route = useRoute()
const isOpen = ref(false)
const isCollapsed = ref(false)

// Stats data (would be fetched from API)
const stats = ref({
  domains: 5,
  users: 23,
  messages: 1234,
  queue: 12
})

const navigationItems = ref([
  {
    to: '/admin/',
    label: 'Dashboard',
    icon: Home
  },
  {
    to: '/admin/users/',
    label: 'Users',
    icon: Users
  },
  {
    to: '/admin/domains/',
    label: 'Domains',
    icon: Globe
  },
  {
    to: '/admin/aliases/',
    label: 'Aliases',
    icon: Mail
  },
  {
    to: '/admin/queue/',
    label: 'Queue',
    icon: Mail
  },
  {
    to: '/admin/settings/',
    label: 'Settings',
    icon: Settings
  }
])

const toggleMobile = () => {
  isOpen.value = !isOpen.value
}

const closeMobile = () => {
  isOpen.value = false
}

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

const handleNavClick = () => {
  // Close mobile drawer on navigation
  if (window.innerWidth < 1024) {
    isOpen.value = false
  }
}

const handleLogout = () => {
  // Implement logout logic
  console.log('Logout clicked')
}

const isActive = (path) => {
  return route.path === path || route.path.startsWith(path)
}
</script>

<script>
export default {
  name: 'AdminSidebar'
}
</script>
