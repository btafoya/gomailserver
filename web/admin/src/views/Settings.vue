<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api/axios'

const authStore = useAuthStore()
const activeTab = ref('profile')
const loading = ref(false)
const saveMessage = ref('')
const saveError = ref('')

// Profile settings
const profileForm = reactive({
  fullName: '',
  email: '',
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Server settings
const serverForm = reactive({
  hostname: '',
  domain: '',
  smtp_submission_port: 587,
  smtp_relay_port: 25,
  smtps_port: 465,
  imap_port: 143,
  imaps_port: 993,
  api_port: 8980,
  max_message_size: 52428800
})

// Security settings
const securityForm = reactive({
  jwt_secret: '',
  clamav_enabled: true,
  clamav_socket_path: '/var/run/clamav/clamd.ctl',
  spamassassin_enabled: true,
  spamassassin_host: 'localhost',
  spamassassin_port: 783,
  rate_limit_enabled: true,
  rate_limit_requests: 100,
  rate_limit_window: 60
})

// TLS/Certificate settings
const tlsForm = reactive({
  acme_enabled: false,
  acme_email: '',
  acme_provider: 'cloudflare',
  acme_api_token: '',
  cert_file: '',
  key_file: ''
})

const tabs = [
  { id: 'profile', label: 'Profile & Account', icon: '👤' },
  { id: 'server', label: 'Server Configuration', icon: '⚙️' },
  { id: 'security', label: 'Security Settings', icon: '🛡️' },
  { id: 'tls', label: 'TLS/Certificates', icon: '🔒' }
]

const loadSettings = async () => {
  loading.value = true
  try {
    // Load user profile
    if (authStore.user) {
      profileForm.fullName = authStore.user.full_name || ''
      profileForm.email = authStore.user.email || ''
    }

    // Load server settings from API
    const serverResponse = await api.get('/v1/settings/server')
    Object.assign(serverForm, serverResponse.data.data)

    // Load security settings from API
    const securityResponse = await api.get('/v1/settings/security')
    Object.assign(securityForm, securityResponse.data.data)

    // Load TLS settings from API
    const tlsResponse = await api.get('/v1/settings/tls')
    Object.assign(tlsForm, tlsResponse.data.data)

  } catch (error) {
    console.error('Failed to load settings:', error)
  } finally {
    loading.value = false
  }
}

const saveProfile = async () => {
  loading.value = true
  saveError.value = ''
  saveMessage.value = ''

  try {
    // Validate password change if attempted
    if (profileForm.newPassword) {
      if (!profileForm.currentPassword) {
        saveError.value = 'Current password is required to change password'
        return
      }
      if (profileForm.newPassword !== profileForm.confirmPassword) {
        saveError.value = 'New passwords do not match'
        return
      }
      if (profileForm.newPassword.length < 8) {
        saveError.value = 'New password must be at least 8 characters'
        return
      }
    }

    // Update profile
    const updates = {
      full_name: profileForm.fullName
    }

    if (authStore.user && authStore.user.id) {
      await api.put(`/v1/users/${authStore.user.id}`, updates)

      // Update password if provided
      if (profileForm.newPassword) {
        await api.put(`/v1/users/${authStore.user.id}/password`, {
          new_password: profileForm.newPassword
        })

        // Clear password fields
        profileForm.currentPassword = ''
        profileForm.newPassword = ''
        profileForm.confirmPassword = ''
      }

      saveMessage.value = 'Profile updated successfully'

      // Update auth store
      if (authStore.user) {
        authStore.user.full_name = profileForm.fullName
      }
    }
  } catch (error) {
    saveError.value = error.response?.data?.message || 'Failed to update profile'
  } finally {
    loading.value = false
  }
}

const saveServerSettings = async () => {
  loading.value = true
  saveError.value = ''
  saveMessage.value = ''

  try {
    const response = await api.put('/v1/settings/server', serverForm)
    saveMessage.value = response.data.message || 'Server settings updated successfully'
  } catch (error) {
    saveError.value = error.response?.data?.message || 'Failed to save server settings'
  } finally {
    loading.value = false
  }
}

const saveSecuritySettings = async () => {
  loading.value = true
  saveError.value = ''
  saveMessage.value = ''

  try {
    const response = await api.put('/v1/settings/security', securityForm)
    saveMessage.value = response.data.message || 'Security settings updated successfully'
  } catch (error) {
    saveError.value = error.response?.data?.message || 'Failed to save security settings'
  } finally {
    loading.value = false
  }
}

const saveTLSSettings = async () => {
  loading.value = true
  saveError.value = ''
  saveMessage.value = ''

  try {
    const response = await api.put('/v1/settings/tls', tlsForm)
    saveMessage.value = response.data.message || 'TLS settings updated successfully'
  } catch (error) {
    saveError.value = error.response?.data?.message || 'Failed to save TLS settings'
  } finally {
    loading.value = false
  }
}

const handleSave = () => {
  switch (activeTab.value) {
    case 'profile':
      saveProfile()
      break
    case 'server':
      saveServerSettings()
      break
    case 'security':
      saveSecuritySettings()
      break
    case 'tls':
      saveTLSSettings()
      break
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<template>
  <div class="p-8">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-foreground">Settings</h1>
      <p class="text-muted-foreground mt-2">Manage your server configuration and preferences</p>
    </div>

    <!-- Alert Messages -->
    <div v-if="saveMessage" class="mb-6 bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg">
      {{ saveMessage }}
    </div>
    <div v-if="saveError" class="mb-6 bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
      {{ saveError }}
    </div>

    <!-- Tabs Navigation -->
    <div class="border-b border-border mb-6">
      <nav class="flex space-x-8">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id; saveMessage = ''; saveError = ''"
          :class="[
            'py-4 px-1 border-b-2 font-medium text-sm transition-colors',
            activeTab === tab.id
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
          ]"
        >
          <span class="mr-2">{{ tab.icon }}</span>
          {{ tab.label }}
        </button>
      </nav>
    </div>

    <!-- Tab Content -->
    <div class="bg-card rounded-lg shadow p-6">
      <!-- Profile & Account Tab -->
      <div v-show="activeTab === 'profile'" class="space-y-6">
        <div>
          <h2 class="text-xl font-semibold text-foreground mb-4">Profile Information</h2>

          <div class="space-y-4 max-w-2xl">
            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Email Address</label>
              <input
                type="email"
                v-model="profileForm.email"
                disabled
                class="w-full px-3 py-2 border border-input rounded-md bg-muted text-muted-foreground cursor-not-allowed"
              />
              <p class="text-xs text-muted-foreground mt-1">Email address cannot be changed</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Full Name</label>
              <input
                type="text"
                v-model="profileForm.fullName"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="Enter your full name"
              />
            </div>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">Change Password</h2>

          <div class="space-y-4 max-w-2xl">
            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Current Password</label>
              <input
                type="password"
                v-model="profileForm.currentPassword"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="Enter current password"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">New Password</label>
              <input
                type="password"
                v-model="profileForm.newPassword"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="Enter new password (min 8 characters)"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Confirm New Password</label>
              <input
                type="password"
                v-model="profileForm.confirmPassword"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="Confirm new password"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Server Configuration Tab -->
      <div v-show="activeTab === 'server'" class="space-y-6">
        <div>
          <h2 class="text-xl font-semibold text-foreground mb-4">General Server Settings</h2>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 max-w-4xl">
            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Hostname</label>
              <input
                type="text"
                v-model="serverForm.hostname"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="mail.example.com"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Domain</label>
              <input
                type="text"
                v-model="serverForm.domain"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="example.com"
              />
            </div>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">SMTP Configuration</h2>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-4xl">
            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Submission Port</label>
              <input
                type="number"
                v-model.number="serverForm.smtp_submission_port"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
              <p class="text-xs text-muted-foreground mt-1">Default: 587</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Relay Port</label>
              <input
                type="number"
                v-model.number="serverForm.smtp_relay_port"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
              <p class="text-xs text-muted-foreground mt-1">Default: 25</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">SMTPS Port</label>
              <input
                type="number"
                v-model.number="serverForm.smtps_port"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
              <p class="text-xs text-muted-foreground mt-1">Default: 465</p>
            </div>
          </div>

          <div class="mt-4 max-w-4xl">
            <label class="block text-sm font-medium text-foreground mb-2">Max Message Size (bytes)</label>
            <input
              type="number"
              v-model.number="serverForm.max_message_size"
              class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            />
            <p class="text-xs text-muted-foreground mt-1">Default: 52428800 (50MB)</p>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">IMAP Configuration</h2>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 max-w-4xl">
            <div>
              <label class="block text-sm font-medium text-foreground mb-2">IMAP Port</label>
              <input
                type="number"
                v-model.number="serverForm.imap_port"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
              <p class="text-xs text-muted-foreground mt-1">Default: 143</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">IMAPS Port</label>
              <input
                type="number"
                v-model.number="serverForm.imaps_port"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
              <p class="text-xs text-muted-foreground mt-1">Default: 993</p>
            </div>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">API Configuration</h2>

          <div class="max-w-4xl">
            <label class="block text-sm font-medium text-foreground mb-2">API Port</label>
            <input
              type="number"
              v-model.number="serverForm.api_port"
              class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            />
            <p class="text-xs text-muted-foreground mt-1">Default: 8980</p>
          </div>
        </div>
      </div>

      <!-- Security Settings Tab -->
      <div v-show="activeTab === 'security'" class="space-y-6">
        <div>
          <h2 class="text-xl font-semibold text-foreground mb-4">Authentication</h2>

          <div class="max-w-4xl">
            <label class="block text-sm font-medium text-foreground mb-2">JWT Secret Key</label>
            <input
              type="password"
              v-model="securityForm.jwt_secret"
              class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="Enter JWT secret key"
            />
            <p class="text-xs text-muted-foreground mt-1">Used to sign authentication tokens. Change with caution.</p>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">Antivirus (ClamAV)</h2>

          <div class="max-w-4xl space-y-4">
            <div class="flex items-center">
              <input
                type="checkbox"
                v-model="securityForm.clamav_enabled"
                id="clamav-enabled"
                class="w-4 h-4 text-primary bg-background border-input rounded focus:ring-2 focus:ring-ring"
              />
              <label for="clamav-enabled" class="ml-2 text-sm font-medium text-foreground">
                Enable ClamAV virus scanning
              </label>
            </div>

            <div v-if="securityForm.clamav_enabled">
              <label class="block text-sm font-medium text-foreground mb-2">ClamAV Socket Path</label>
              <input
                type="text"
                v-model="securityForm.clamav_socket_path"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="/var/run/clamav/clamd.ctl"
              />
            </div>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">Anti-Spam (SpamAssassin)</h2>

          <div class="max-w-4xl space-y-4">
            <div class="flex items-center">
              <input
                type="checkbox"
                v-model="securityForm.spamassassin_enabled"
                id="spam-enabled"
                class="w-4 h-4 text-primary bg-background border-input rounded focus:ring-2 focus:ring-ring"
              />
              <label for="spam-enabled" class="ml-2 text-sm font-medium text-foreground">
                Enable SpamAssassin filtering
              </label>
            </div>

            <div v-if="securityForm.spamassassin_enabled" class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-foreground mb-2">SpamAssassin Host</label>
                <input
                  type="text"
                  v-model="securityForm.spamassassin_host"
                  class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                  placeholder="localhost"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-foreground mb-2">SpamAssassin Port</label>
                <input
                  type="number"
                  v-model.number="securityForm.spamassassin_port"
                  class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">Rate Limiting</h2>

          <div class="max-w-4xl space-y-4">
            <div class="flex items-center">
              <input
                type="checkbox"
                v-model="securityForm.rate_limit_enabled"
                id="rate-limit-enabled"
                class="w-4 h-4 text-primary bg-background border-input rounded focus:ring-2 focus:ring-ring"
              />
              <label for="rate-limit-enabled" class="ml-2 text-sm font-medium text-foreground">
                Enable rate limiting
              </label>
            </div>

            <div v-if="securityForm.rate_limit_enabled" class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-foreground mb-2">Max Requests</label>
                <input
                  type="number"
                  v-model.number="securityForm.rate_limit_requests"
                  class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-foreground mb-2">Time Window (seconds)</label>
                <input
                  type="number"
                  v-model.number="securityForm.rate_limit_window"
                  class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- TLS/Certificates Tab -->
      <div v-show="activeTab === 'tls'" class="space-y-6">
        <div>
          <h2 class="text-xl font-semibold text-foreground mb-4">Let's Encrypt (ACME)</h2>

          <div class="max-w-4xl space-y-4">
            <div class="flex items-center">
              <input
                type="checkbox"
                v-model="tlsForm.acme_enabled"
                id="acme-enabled"
                class="w-4 h-4 text-primary bg-background border-input rounded focus:ring-2 focus:ring-ring"
              />
              <label for="acme-enabled" class="ml-2 text-sm font-medium text-foreground">
                Enable automatic certificate management with Let's Encrypt
              </label>
            </div>

            <div v-if="tlsForm.acme_enabled" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-foreground mb-2">Email Address</label>
                <input
                  type="email"
                  v-model="tlsForm.acme_email"
                  class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                  placeholder="admin@example.com"
                />
                <p class="text-xs text-muted-foreground mt-1">Used for certificate renewal notifications</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-foreground mb-2">DNS Provider</label>
                <select
                  v-model="tlsForm.acme_provider"
                  class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                >
                  <option value="cloudflare">Cloudflare</option>
                  <option value="route53">AWS Route 53</option>
                  <option value="digitalocean">DigitalOcean</option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-medium text-foreground mb-2">API Token</label>
                <input
                  type="password"
                  v-model="tlsForm.acme_api_token"
                  class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                  placeholder="Enter DNS provider API token"
                />
                <p class="text-xs text-muted-foreground mt-1">Required for DNS-01 challenge validation</p>
              </div>
            </div>
          </div>
        </div>

        <div class="border-t border-border pt-6">
          <h2 class="text-xl font-semibold text-foreground mb-4">Manual Certificates</h2>

          <div class="max-w-4xl space-y-4">
            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Certificate File Path</label>
              <input
                type="text"
                v-model="tlsForm.cert_file"
                :disabled="tlsForm.acme_enabled"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:bg-muted disabled:cursor-not-allowed"
                placeholder="/path/to/certificate.crt"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-foreground mb-2">Private Key File Path</label>
              <input
                type="text"
                v-model="tlsForm.key_file"
                :disabled="tlsForm.acme_enabled"
                class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:bg-muted disabled:cursor-not-allowed"
                placeholder="/path/to/private-key.key"
              />
            </div>

            <p v-if="tlsForm.acme_enabled" class="text-sm text-muted-foreground">
              Manual certificate paths are disabled when ACME is enabled
            </p>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <div class="mt-8 pt-6 border-t border-border flex justify-end">
        <button
          @click="handleSave"
          :disabled="loading"
          class="bg-primary text-primary-foreground hover:bg-primary/90 font-medium py-2 px-6 rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ loading ? 'Saving...' : 'Save Changes' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
