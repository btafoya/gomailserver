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
        'bg-gray-900 text-gray-100',
        'w-[280px] lg:translate-x-0',
        isOpen ? 'translate-x-0' : '-translate-x-full'
      ]"
    >
      <!-- Logo / Brand Area -->
      <div class="px-6 py-4 border-b border-gray-800">
        <NuxtLink to="/admin" class="flex items-center space-x-3 hover:opacity-90 transition-opacity">
          <div class="h-10 w-10 bg-blue-600 rounded-lg flex items-center justify-center shadow-lg">
            <Mail class="h-6 w-6 text-white" />
          </div>
          <h1 class="text-xl font-bold text-white">GoMail</h1>
        </NuxtLink>
      </div>

      <!-- Navigation Groups -->
      <div class="flex-1 overflow-y-auto py-4 px-3">
        <nav class="space-y-6">
          <!-- Group: Main -->
          <div v-if="mainGroupOpen">
            <button
              @click="mainGroupOpen = false"
              class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
            >
              <span>Main</span>
              <ChevronDown class="h-4 w-4" />
            </button>
            <div class="space-y-1">
              <NuxtLink
                v-for="item in mainGroupItems"
                :key="item.to"
                :to="item.to"
                :class="[
                  'sidebar-nav-item flex items-center space-x-3 px-3 py-2.5 rounded-lg border border-transparent transition-all duration-200 cursor-pointer',
                  isActive(item.to)
                    ? 'bg-blue-600 text-white shadow-md'
                    : 'hover:bg-gray-800 text-gray-300 hover:text-white hover:border-gray-700'
                ]"
                @click="handleNavClick"
              >
                <component :is="item.icon" class="h-5 w-5 flex-shrink-0" />
                <span class="font-medium">{{ item.label }}</span>
              </NuxtLink>
            </div>
          </div>
          <div v-else>
            <button
              @click="mainGroupOpen = true"
              class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
            >
              <span>Main</span>
              <ChevronRight class="h-4 w-4" />
            </button>
          </div>

          <!-- Group: Management -->
          <div v-if="managementGroupOpen">
            <button
              @click="managementGroupOpen = false"
              class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
            >
              <span>Management</span>
              <ChevronDown class="h-4 w-4" />
            </button>
            <div class="space-y-1">
              <NuxtLink
                v-for="item in managementGroupItems"
                :key="item.to"
                :to="item.to"
                :class="[
                  'sidebar-nav-item flex items-center space-x-3 px-3 py-2.5 rounded-lg border border-transparent transition-all duration-200 cursor-pointer',
                  isActive(item.to)
                    ? 'bg-blue-600 text-white shadow-md'
                    : 'hover:bg-gray-800 text-gray-300 hover:text-white hover:border-gray-700'
                ]"
                @click="handleNavClick"
              >
                <component :is="item.icon" class="h-5 w-5 flex-shrink-0" />
                <span class="font-medium">{{ item.label }}</span>
              </NuxtLink>
            </div>
          </div>
          <div v-else>
            <button
              @click="managementGroupOpen = true"
              class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
            >
              <span>Management</span>
              <ChevronRight class="h-4 w-4" />
            </button>
           </div>

           <!-- Group: Reputation Management -->
           <div v-if="reputationGroupOpen">
             <button
               @click="reputationGroupOpen = false"
               class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
             >
               <span>Reputation Management</span>
               <ChevronDown class="h-4 w-4" />
             </button>
             <div class="space-y-1">
               <NuxtLink
                 v-for="item in reputationGroupItems"
                 :key="item.to"
                 :to="item.to"
                 :class="[
                   'sidebar-nav-item flex items-center space-x-3 px-3 py-2.5 rounded-lg border border-transparent transition-all duration-200 cursor-pointer',
                   isActive(item.to)
                     ? 'bg-blue-600 text-white shadow-md'
                     : 'hover:bg-gray-800 text-gray-300 hover:text-white hover:border-gray-700'
                 ]"
                 @click="handleNavClick"
               >
                 <component :is="item.icon" class="h-5 w-5 flex-shrink-0" />
                 <span class="font-medium">{{ item.label }}</span>
               </NuxtLink>
             </div>
           </div>
           <div v-else>
             <button
               @click="reputationGroupOpen = true"
               class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
             >
               <span>Reputation Management</span>
               <ChevronRight class="h-4 w-4" />
             </button>
           </div>

           <!-- Group: System -->
          <div v-if="systemGroupOpen">
            <button
              @click="systemGroupOpen = false"
              class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
            >
              <span>System</span>
              <ChevronDown class="h-4 w-4" />
            </button>
            <div class="space-y-1">
              <NuxtLink
                v-for="item in systemGroupItems"
                :key="item.to"
                :to="item.to"
                :class="[
                  'sidebar-nav-item flex items-center space-x-3 px-3 py-2.5 rounded-lg border border-transparent transition-all duration-200 cursor-pointer',
                  isActive(item.to)
                    ? 'bg-blue-600 text-white shadow-md'
                    : 'hover:bg-gray-800 text-gray-300 hover:text-white hover:border-gray-700'
                ]"
                @click="handleNavClick"
              >
                <component :is="item.icon" class="h-5 w-5 flex-shrink-0" />
                <span class="font-medium">{{ item.label }}</span>
              </NuxtLink>
            </div>
          </div>
          <div v-else>
            <button
              @click="systemGroupOpen = true"
              class="w-full flex items-center justify-between px-3 py-2 mb-2 text-sm font-semibold text-gray-400 hover:text-white transition-colors"
            >
              <span>System</span>
              <ChevronRight class="h-4 w-4" />
            </button>
          </div>
        </nav>
      </div>

      <!-- User Profile Section -->
      <div class="px-3 py-4 border-t border-gray-800">
        <div class="flex items-center space-x-3">
          <div class="h-10 w-10 rounded-full bg-blue-600 flex items-center justify-center text-white font-bold text-sm shadow-md">
            AD
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-white truncate">Admin User</p>
            <p class="text-xs text-gray-400 truncate">admin@localhost</p>
          </div>
          <button
            @click="handleLogout"
            class="p-2 rounded-lg hover:bg-gray-800 transition-colors cursor-pointer"
            title="Logout"
          >
            <LogOut class="h-4 w-4 text-gray-400 hover:text-white transition-colors" />
          </button>
        </div>
      </div>
    </aside>

    <!-- Mobile Toggle Button -->
    <button
      @click="toggleMobile"
      class="fixed bottom-6 right-6 z-50 p-4 rounded-full bg-blue-600 text-white shadow-lg hover:bg-blue-700 transition-colors lg:hidden cursor-pointer"
      :aria-label="isOpen ? 'Close sidebar' : 'Open sidebar'"
    >
      <Menu v-if="!isOpen" class="h-6 w-6" />
      <X v-else class="h-6 w-6" />
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
 import { Home, Users, Globe, Mail, Settings, Menu, X, LogOut, ChevronDown, ChevronRight, ShieldCheck, AlertTriangle, TrendingUp, FileText, BarChart, Sliders, Calendar, Brain } from 'lucide-vue-next'

// useRoute is auto-imported by Nuxt 3
const route = useRoute()
const isOpen = ref(true) // Visible by default on desktop

 // Collapsible groups state
 const mainGroupOpen = ref(true)
 const managementGroupOpen = ref(true)
 const reputationGroupOpen = ref(true)
 const systemGroupOpen = ref(true)

// Main group: Dashboard, Queue
const mainGroupItems = ref([
  {
    to: '/admin/',
    label: 'Dashboard',
    icon: Home
  },
  {
    to: '/admin/queue/',
    label: 'Queue',
    icon: Mail
  }
])

 // Management group: Users, Domains, Aliases
 const managementGroupItems = ref([
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
   }
 ])

 // Reputation Management group: 9 items
 const reputationGroupItems = ref([
   {
     to: '/admin/reputation',
     label: 'Overview',
     icon: Home
   },
   {
     to: '/admin/reputation/circuit-breakers',
     label: 'Circuit Breakers',
     icon: AlertTriangle
   },
   {
     to: '/admin/reputation/warmup',
     label: 'Warm-up',
     icon: TrendingUp
   },
   {
     to: '/admin/reputation/audit',
     label: 'Audit',
     icon: ShieldCheck
   },
   {
     to: '/admin/reputation/dmarc-reports',
     label: 'DMARC Reports',
     icon: FileText
   },
   {
     to: '/admin/reputation/external-metrics',
     label: 'External Metrics',
     icon: BarChart
   },
   {
     to: '/admin/reputation/provider-limits',
     label: 'Provider Limits',
     icon: Sliders
   },
   {
     to: '/admin/reputation/warmup-scheduler',
     label: 'Warmup Scheduler',
     icon: Calendar
   },
   {
     to: '/admin/reputation/predictions',
     label: 'Predictions',
     icon: Brain
   }
 ])

// System group: Settings
const systemGroupItems = ref([
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
