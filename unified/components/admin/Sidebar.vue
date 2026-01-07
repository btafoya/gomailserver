<template>
  <div class="flex">
    <aside
      :class="[
        'flex flex-col border-r bg-card transition-all duration-300',
        isOpen ? 'w-64' : 'w-0 translate-x-[-100%]'
      ]"
    >
      <div class="p-4">
        <div class="flex items-center space-x-3 mb-6">
          <div class="h-8 w-8 bg-primary rounded-lg flex items-center justify-center">
            <Mail class="h-5 w-5 text-primary-foreground" />
          </div>
          <h1 v-if="!isOpen" class="text-lg font-bold text-foreground">GoMail</h1>
        </div>

        <nav class="space-y-1">
          <NuxtLink
            v-for="item in navigationItems"
            :key="item.to"
            :to="item.to"
            :class="[
              'flex items-center space-x-3 px-3 py-2.5 rounded-lg transition-colors',
              isActive(item.to)
                ? 'bg-primary text-primary-foreground'
                : 'hover:bg-muted hover:text-foreground'
            ]"
          >
            <component :is="item.icon" class="h-5 w-5" />
            <span class="font-medium">{{ item.label }}</span>
          </NuxtLink>
        </nav>

        <div class="mt-6 pt-6 border-t">
          <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-3">
            Quick Stats
          </h3>
          <div class="space-y-3">
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">Domains</span>
              <span class="font-medium text-foreground">5</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">Users</span>
              <span class="font-medium text-foreground">23</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">Messages</span>
              <span class="font-medium text-foreground">1,234</span>
            </div>
          </div>
        </div>
      </div>
    </aside>

    <button
      @click="toggle"
      class="fixed bottom-4 right-4 z-50 p-3 rounded-full bg-primary text-primary-foreground shadow-lg hover:bg-primary/90 transition-colors lg:hidden"
    >
      <Menu class="h-6 w-6" />
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Home, Users, Globe, Mail, Settings, Menu } from 'lucide-vue-next'
import { useRoute } from 'vue-router'

const route = useRoute()
const isOpen = ref(false)

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

const toggle = () => {
  isOpen.value = !isOpen.value
}

const isActive = (path) => {
  return route.path === path
}
</script>

<script>
export default {
  name: 'AdminSidebar'
}
</script>
