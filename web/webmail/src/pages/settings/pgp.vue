<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const user = ref(authStore.user)
const keys = ref([])
const loading = ref(false)
const importDialog = ref(false)
const publicKeyText = ref('')
const error = ref('')
const success = ref('')

const loadKeys = async () => {
  if (!user.value?.id) return

  loading.value = true
  error.value = ''

  try {
    const response = await fetch(`/api/v1/pgp/users/${user.value.id}/keys`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to load PGP keys')
    }

    keys.value = await response.json()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const importKey = async () => {
  if (!publicKeyText.value.trim()) {
    error.value = 'Please enter a PGP public key'
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await fetch('/api/v1/pgp/keys', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        user_id: user.value.id,
        public_key: publicKeyText.value,
      }),
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || 'Failed to import key')
    }

    success.value = 'PGP key imported successfully'
    publicKeyText.value = ''
    importDialog.value = false
    await loadKeys()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const setPrimary = async (keyId) => {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await fetch(`/api/v1/pgp/keys/${keyId}/primary`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to set primary key')
    }

    success.value = 'Primary key updated'
    await loadKeys()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const deleteKey = async (keyId) => {
  if (!confirm('Are you sure you want to delete this PGP key?')) {
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await fetch(`/api/v1/pgp/keys/${keyId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to delete key')
    }

    success.value = 'PGP key deleted'
    await loadKeys()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadKeys()
})
</script>

<template>
  <div class="container mx-auto py-8 px-4 max-w-4xl">
    <div class="mb-6">
      <h1 class="text-3xl font-bold mb-2">PGP Key Management</h1>
      <p class="text-muted-foreground">Manage your PGP public keys for email encryption</p>
    </div>

    <div v-if="error" class="bg-destructive/15 text-destructive p-4 rounded-md mb-4">
      {{ error }}
    </div>

    <div v-if="success" class="bg-green-500/15 text-green-700 p-4 rounded-md mb-4">
      {{ success }}
    </div>

    <div class="flex justify-between items-center mb-6">
      <h2 class="text-xl font-semibold">Your PGP Keys</h2>
      <button
        @click="importDialog = true"
        class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
      >
        Import Key
      </button>
    </div>

    <div v-if="loading && keys.length === 0" class="text-center py-12">
      <Icon name="lucide:loader-circle" class="w-8 h-8 animate-spin mx-auto text-primary" />
      <p class="mt-4 text-muted-foreground">Loading keys...</p>
    </div>

    <div v-else-if="keys.length === 0" class="text-center py-12 border-2 border-dashed rounded-lg">
      <Icon name="lucide:key" class="w-12 h-12 mx-auto text-muted-foreground mb-4" />
      <p class="text-lg font-medium mb-2">No PGP keys</p>
      <p class="text-muted-foreground mb-4">Import a PGP public key to enable email encryption</p>
      <button
        @click="importDialog = true"
        class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
      >
        Import Your First Key
      </button>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="key in keys"
        :key="key.id"
        class="border rounded-lg p-4 hover:shadow-md transition-shadow"
      >
        <div class="flex justify-between items-start">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <h3 class="font-semibold">Key ID: {{ key.key_id }}</h3>
              <span
                v-if="key.is_primary"
                class="px-2 py-1 bg-primary/20 text-primary text-xs rounded-full"
              >
                Primary
              </span>
            </div>
            <p class="text-sm text-muted-foreground mb-1">
              Fingerprint: <code class="bg-muted px-1 py-0.5 rounded">{{ key.fingerprint }}</code>
            </p>
            <p class="text-sm text-muted-foreground mb-1">
              Created: {{ new Date(key.created_at).toLocaleDateString() }}
            </p>
            <p v-if="key.expires_at" class="text-sm text-muted-foreground">
              Expires: {{ new Date(key.expires_at).toLocaleDateString() }}
            </p>
          </div>
          <div class="flex gap-2">
            <button
              v-if="!key.is_primary"
              @click="setPrimary(key.id)"
              :disabled="loading"
              class="px-3 py-1 text-sm border rounded-md hover:bg-accent disabled:opacity-50"
            >
              Set Primary
            </button>
            <button
              @click="deleteKey(key.id)"
              :disabled="loading"
              class="px-3 py-1 text-sm border border-destructive text-destructive rounded-md hover:bg-destructive/10 disabled:opacity-50"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Import Dialog -->
    <div
      v-if="importDialog"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @click.self="importDialog = false"
    >
      <div class="bg-background border rounded-lg p-6 max-w-2xl w-full mx-4">
        <h2 class="text-xl font-semibold mb-4">Import PGP Public Key</h2>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">
            Public Key (ASCII Armored)
          </label>
          <textarea
            v-model="publicKeyText"
            class="w-full h-64 p-3 border rounded-md font-mono text-sm"
            placeholder="-----BEGIN PGP PUBLIC KEY BLOCK-----&#10;&#10;-----END PGP PUBLIC KEY BLOCK-----"
          ></textarea>
        </div>

        <div class="flex justify-end gap-2">
          <button
            @click="importDialog = false"
            :disabled="loading"
            class="px-4 py-2 border rounded-md hover:bg-accent disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            @click="importKey"
            :disabled="loading"
            class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
          >
            <Icon v-if="loading" name="lucide:loader-circle" class="w-4 h-4 animate-spin inline mr-2" />
            Import Key
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
