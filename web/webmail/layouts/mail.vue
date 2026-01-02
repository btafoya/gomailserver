<script setup lang="ts">
const { checkAuth, user, logout } = useAuth()
const mailStore = useMailStore()
const router = useRouter()
const darkMode = ref(false)
const showShortcutsHelp = ref(false)
const { registerShortcut } = useKeyboardShortcuts()

onMounted(() => {
  checkAuth()
  mailStore.fetchMailboxes()

  // Load dark mode preference
  if (process.client) {
    darkMode.value = localStorage.getItem('darkMode') === 'true'
    if (darkMode.value) {
      document.documentElement.classList.add('dark')
    }
  }

  // Register global keyboard shortcuts
  registerShortcut({
    key: 'c',
    handler: () => router.push('/mail/compose'),
    description: 'Compose new message'
  })

  registerShortcut({
    key: '?',
    shiftKey: true,
    handler: () => { showShortcutsHelp.value = !showShortcutsHelp.value },
    description: 'Show keyboard shortcuts'
  })

  registerShortcut({
    key: 'Escape',
    handler: () => { showShortcutsHelp.value = false },
    description: 'Close dialogs'
  })
})

const toggleDarkMode = () => {
  darkMode.value = !darkMode.value
  if (process.client) {
    if (darkMode.value) {
      document.documentElement.classList.add('dark')
      localStorage.setItem('darkMode', 'true')
    } else {
      document.documentElement.classList.remove('dark')
      localStorage.setItem('darkMode', 'false')
    }
  }
}
</script>

<template>
  <div class="h-screen flex flex-col bg-background">
    <!-- Top Header -->
    <header class="border-b bg-card px-4 py-3 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <Icon name="lucide:mail" class="w-6 h-6 text-primary" />
        <h1 class="text-xl font-semibold">Webmail</h1>
      </div>

      <div class="flex items-center gap-2">
        <!-- Search -->
        <button class="p-2 hover:bg-muted rounded-md">
          <Icon name="lucide:search" class="w-5 h-5" />
        </button>

        <!-- Dark Mode Toggle -->
        <button @click="toggleDarkMode" class="p-2 hover:bg-muted rounded-md">
          <Icon :name="darkMode ? 'lucide:sun' : 'lucide:moon'" class="w-5 h-5" />
        </button>

        <!-- User Menu -->
        <div class="flex items-center gap-3 ml-2">
          <div class="text-right text-sm">
            <div class="font-medium">{{ user?.full_name }}</div>
            <div class="text-muted-foreground text-xs">{{ user?.email }}</div>
          </div>
          <button @click="logout" class="p-2 hover:bg-muted rounded-md" title="Logout">
            <Icon name="lucide:log-out" class="w-5 h-5" />
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content Area -->
    <div class="flex-1 flex overflow-hidden">
      <slot />
    </div>

    <!-- Keyboard Shortcuts Help Modal -->
    <div v-if="showShortcutsHelp" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click="showShortcutsHelp = false">
      <div class="bg-card border rounded-lg p-6 max-w-2xl max-h-[80vh] overflow-y-auto" @click.stop>
        <h2 class="text-2xl font-bold mb-4">Keyboard Shortcuts</h2>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <h3 class="font-semibold mb-2">Navigation</h3>
            <div class="space-y-1 text-sm">
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">c</kbd><span>Compose</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">j</kbd><span>Next message</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">k</kbd><span>Previous message</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">o/Enter</kbd><span>Open message</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">Esc</kbd><span>Close/Cancel</span></div>
            </div>
          </div>
          <div>
            <h3 class="font-semibold mb-2">Actions</h3>
            <div class="space-y-1 text-sm">
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">r</kbd><span>Reply</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">a</kbd><span>Reply all</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">f</kbd><span>Forward</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">#</kbd><span>Delete</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">s</kbd><span>Star/unstar</span></div>
            </div>
          </div>
          <div>
            <h3 class="font-semibold mb-2">Go to</h3>
            <div class="space-y-1 text-sm">
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">g i</kbd><span>Inbox</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">g s</kbd><span>Sent</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">g t</kbd><span>Trash</span></div>
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">/</kbd><span>Search</span></div>
            </div>
          </div>
          <div>
            <h3 class="font-semibold mb-2">Other</h3>
            <div class="space-y-1 text-sm">
              <div class="flex justify-between"><kbd class="px-2 py-1 bg-muted rounded">?</kbd><span>Show shortcuts</span></div>
            </div>
          </div>
        </div>
        <div class="mt-4 pt-4 border-t text-center">
          <button @click="showShortcutsHelp = false" class="px-4 py-2 bg-primary text-primary-foreground rounded-md">
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
