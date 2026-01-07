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
        <h2 class="text-3xl font-bold tracking-tight">Edit User</h2>
        <NuxtLink to="/admin/users">
          <UButton variant="outline">
            Back to Users
          </UButton>
        </NuxtLink>
      </div>

      <UCard class="max-w-2xl mx-auto">
        <UCardHeader>
          <UCardTitle>User Information</UCardTitle>
        </UCardHeader>
        <UCardContent class="space-y-4">
          <div v-if="loading" class="text-center py-12">
            <p class="text-muted-foreground">Loading user...</p>
          </div>

          <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
            Error loading user: {{ error }}
          </div>

          <div v-else class="space-y-4">
            <div>
              <label class="text-sm font-medium">Email</label>
              <input
                type="email"
                v-model="user.email"
                disabled
                class="w-full px-3 py-2 border rounded-md bg-muted text-foreground"
              />
              <p class="text-xs text-muted-foreground mt-1">Email address cannot be changed</p>
            </div>

            <div>
              <label class="text-sm font-medium">Full Name</label>
              <input
                type="text"
                v-model="user.full_name"
                class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
              />
            </div>

            <div>
              <label class="text-sm font-medium">Domain</label>
              <input
                type="text"
                v-model="user.domain"
                disabled
                class="w-full px-3 py-2 border rounded-md bg-muted text-foreground"
              />
              <p class="text-xs text-muted-foreground mt-1">Domain cannot be changed</p>
            </div>

            <div>
              <label class="text-sm font-medium">Password</label>
              <input
                type="password"
                v-model="password"
                placeholder="Leave blank to keep current password"
                class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
              />
            </div>

            <div>
              <label class="text-sm font-medium">Confirm Password</label>
              <input
                type="password"
                v-model="confirmPassword"
                placeholder="Re-enter new password"
                class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
              />
            </div>

            <div class="flex items-center space-x-4">
              <label class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  v-model="user.is_admin"
                  class="rounded"
                />
                <span class="text-sm font-medium">Admin Access</span>
              </label>
              <label class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  v-model="user.is_active"
                  class="rounded"
                />
                <span class="text-sm font-medium">Active</span>
              </label>
            </div>

            <div class="border-t pt-4 space-y-3">
              <h3 class="text-lg font-semibold">Two-Factor Authentication</h3>
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium">2FA Enabled</p>
                  <p class="text-xs text-muted-foreground">Add extra security to account</p>
                </div>
                <UButton variant="outline" size="sm">
                  {{ user.two_factor_enabled ? 'Disable' : 'Enable' }}
                </UButton>
              </div>
            </div>

            <div class="border-t pt-4 space-y-3">
              <h3 class="text-lg font-semibold">Account Information</h3>
              <div class="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p class="text-muted-foreground">User ID</p>
                  <p class="font-medium">{{ user.id }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Created</p>
                  <p class="font-medium">{{ new Date(user.created_at).toLocaleDateString() }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Last Login</p>
                  <p class="font-medium">{{ user.last_login ? new Date(user.last_login).toLocaleString() : 'Never' }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Emails Sent</p>
                  <p class="font-medium">{{ user.emails_sent || 0 }}</p>
                </div>
              </div>
            </div>

            <div class="flex justify-end space-x-2">
              <UButton variant="outline" @click="$router.push('/admin/users')">
                Cancel
              </UButton>
              <UButton>
                <Save class="mr-2 h-4 w-4" />
                Save Changes
              </UButton>
            </div>
          </div>
        </UCardContent>
      </UCard>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Save } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const user = ref({
  id: null,
  email: '',
  full_name: '',
  domain: '',
  is_admin: false,
  is_active: true,
  two_factor_enabled: false,
  created_at: null,
  last_login: null,
  emails_sent: 0
})

const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref(null)

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}

// TODO: Replace with actual API call
onMounted(() => {
  const userId = route.params.id
  user.value = {
    id: userId,
    email: 'admin@example.com',
    full_name: 'Admin User',
    domain: 'example.com',
    is_admin: true,
    is_active: true,
    two_factor_enabled: false,
    created_at: new Date('2024-01-01').toISOString(),
    last_login: new Date(Date.now() - 86400000).toISOString(),
    emails_sent: 1234
  }
})
</script>

<script>
export default {
  middleware: 'auth',
  layout: 'admin'
}
</script>
