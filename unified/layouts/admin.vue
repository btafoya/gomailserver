<template>
  <div class="min-h-screen bg-background flex">
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content Area -->
    <main class="flex-1 transition-smooth">
      <!-- Header with Breadcrumb and Controls -->
      <header class="sticky top-0 z-30 flex h-16 shrink-0 items-center gap-2 px-6 border-b bg-card/95 backdrop-blur-md">
        <!-- Breadcrumb -->
        <nav class="flex items-center space-x-2 text-sm">
          <span class="text-muted-foreground">Admin</span>
          <span class="text-muted-foreground/50">/</span>
          <span
            v-for="(crumb, index) in breadcrumbs"
            :key="index"
            class="breadcrumb-item"
          >
            {{ crumb }}
          </span>
        </nav>

        <!-- Spacer -->
        <div class="flex-1" />

        <!-- Header Controls -->
        <div class="flex items-center gap-2">
          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="p-2 rounded-lg hover:bg-muted hover-effect focus-ring"
            title="Toggle theme"
            :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
          >
            <Sun v-if="isDark" class="h-5 w-5 text-foreground" />
            <Moon v-else class="h-5 w-5 text-foreground" />
          </button>

          <!-- Notifications -->
          <button
            class="p-2 rounded-lg hover:bg-muted hover-effect focus-ring"
            title="Notifications"
            aria-label="Notifications"
          >
            <Bell class="h-5 w-5 text-foreground" />
          </button>

          <!-- User Menu -->
          <div class="flex items-center space-x-2">
            <div class="h-8 w-8 rounded-full bg-primary flex items-center justify-center text-primary-foreground font-bold text-sm">
              A
            </div>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <div class="content-wrapper min-h-[calc(100vh-4rem)]">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Sun, Moon, Bell } from 'lucide-vue-next'

// Theme management
const isDark = ref(false)

const toggleTheme = () => {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

// Breadcrumb generation based on route
const breadcrumbs = computed(() => {
  const route = useRoute()
  const path = route.path.replace('/admin', '').trim()
  if (!path) return ['Dashboard']

  const parts = path.split('/').filter(Boolean)
  const result = ['Dashboard', ...parts]
  return result
})

// Initialize theme on mount
onMounted(() => {
  // Check system preference
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
})
</script>

<script>
export default {
  name: 'AdminLayout',
  layout: 'admin'
}
</script>

<style scoped>
.breadcrumb-item:last-child {
  color: var(--color-foreground);
  font-weight: 500;
}
</style>
