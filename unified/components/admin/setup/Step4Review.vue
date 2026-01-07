<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Review & Complete</h2>
      <p class="text-gray-600 mt-2">Review your configuration before completing setup</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Configuration Summary</CardTitle>
      </CardHeader>
      <CardContent class="space-y-6">
        <!-- System Configuration -->
        <div>
          <h3 class="text-lg font-semibold text-gray-800 mb-3">System Settings</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">Server Hostname:</span>
              <span class="font-medium text-gray-900">{{ config?.system?.hostname || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">API Port:</span>
              <span class="font-medium text-gray-900">{{ config?.system?.port || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">SMTP Port:</span>
              <span class="font-medium text-gray-900">{{ config?.system?.smtp_port || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">IMAP Port:</span>
              <span class="font-medium text-gray-900">{{ config?.system?.imap_port || '-' }}</span>
            </div>
            <div class="flex justify-between py-2">
              <span class="text-gray-600">TLS Enabled:</span>
              <Badge :variant="config?.system?.tls_enabled ? 'default' : 'secondary'">
                {{ config?.system?.tls_enabled ? 'Yes' : 'No' }}
              </Badge>
            </div>
          </div>
        </div>

        <!-- Domain Configuration -->
        <div class="mt-6">
          <h3 class="text-lg font-semibold text-gray-800 mb-3">Domain Configuration</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">Primary Domain:</span>
              <span class="font-medium text-gray-900">{{ config?.domain?.domain || '-' }}</span>
            </div>
            <div v-if="config?.domain?.audit_results" class="mt-3">
              <h4 class="text-sm font-semibold text-gray-800 mb-2">DNS Audit Results</h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <span class="text-gray-600">Overall Score:</span>
                  <Badge :variant="getScoreVariant(config.domain.audit_results.overall_score)">
                    {{ config.domain.audit_results.overall_score }}/100 }}
                  </Badge>
                </div>
                <div>
                  <span class="text-gray-600">SPF:</span>
                  <Badge :variant="config.domain.audit_results.spf.passed ? 'default' : 'destructive'">
                    {{ config.domain.audit_results.spf.passed ? 'PASS' : 'FAIL' }}
                  </Badge>
                </div>
                <div>
                  <span class="text-gray-600">DKIM:</span>
                  <Badge :variant="config.domain.audit_results.dkim.passed ? 'default' : 'destructive'">
                    {{ config.domain.audit_results.dkim.passed ? 'PASS' : 'FAIL' }}
                  </Badge>
                </div>
                <div>
                  <span class="text-gray-600">DMARC:</span>
                  <Badge :variant="config.domain.audit_results.dmarc.passed ? 'default' : 'destructive'">
                    {{ config.domain.audit_results.dmarc.passed ? 'PASS' : 'FAIL' }}
                  </Badge>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Admin User Configuration -->
        <div class="mt-6">
          <h3 class="text-lg font-semibold text-gray-800 mb-3">Admin User</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">Email:</span>
              <span class="font-medium text-gray-900">{{ config?.admin?.email || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">Full Name:</span>
              <span class="font-medium text-gray-900">{{ config?.admin?.full_name || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b">
              <span class="text-gray-600">2FA Enabled:</span>
              <Badge :variant="config?.admin?.enable_totp ? 'default' : 'secondary'">
                {{ config?.admin?.enable_totp ? 'Yes' : 'No' }}
              </Badge>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <Alert v-if="error" variant="destructive" class="mt-4">
      {{ error }}
    </Alert>

    <Alert class="mt-4">
      <AlertDescription>
        After completing setup, you will be redirected to the admin dashboard.
        Make sure all configuration is correct before proceeding.
      </AlertDescription>
    </Alert>

    <div class="flex justify-between mt-6">
      <Button variant="outline" :disabled="isSubmitting" @click="handlePrevious">
        Previous
      </Button>
      <Button @click="handleComplete" :disabled="isSubmitting">
        Complete Setup
        <template v-if="isSubmitting">
          <span class="ml-2">Completing...</span>
        </template>
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
const config = ref<{
  system: any
  domain: any
  admin: any
} | null>(null)

const isSubmitting = ref(false)
const error = ref('')

// Emit events
const emit = defineEmits<{
  previous: []
  complete: []
}>()

// Computed helpers
const getScoreVariant = (score: number) => {
  if (score >= 70) return 'default'
  if (score >= 50) return 'secondary'
  return 'destructive'
}

const handlePrevious = () => {
  emit('previous')
}

const handleComplete = async () => {
  isSubmitting.value = true
  error.value = ''

  try {
    const { completeSetup } = await import('~/composables/api/setup')
    await completeSetup()
    // Redirect will be handled by middleware
  } catch (err: any) {
    error.value = err.message || 'Failed to complete setup'
    isSubmitting.value = false
  }
}
</script>
