<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Domain Setup</h2>
      <p class="text-gray-600 mt-2">Configure your primary domain</p>
    </div>

    <UCard>
      <UCardHeader>
        <UCardTitle>Primary Domain</UCardTitle>
      </UCardHeader>
      <UCardContent class="space-y-4">
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Domain Name</label>
          <UInput
            v-model="form.domain"
            type="text"
            placeholder="example.com"
            :disabled="isSubmitting"
            @blur="validateDomain"
          />
          <p class="text-xs text-gray-500 mt-1">
            Your primary domain (e.g., mail.example.com)
          </p>
        </div>

        <div class="space-y-2">
          <UButton
            type="button"
            variant="outline"
            @click="runAudit"
            :disabled="isSubmitting || !form.domain"
            :loading="isValidating"
          >
            <template v-if="!isValidating">
              <Check class="mr-2 h-4 w-4" />
              Run DNS Audit
            </template>
            <template v-else>
              <span class="ml-2">Checking DNS...</span>
            </template>
          </UButton>
        </div>

        <!-- Audit Results -->
        <div v-if="auditResults" class="mt-4 p-4 bg-gray-50 rounded-lg">
          <h3 class="text-lg font-semibold mb-3">Audit Results</h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <!-- DNS Records -->
            <div class="flex items-center justify-between p-2 bg-white rounded">
              <span class="text-sm font-medium">SPF</span>
              <UBadge :variant="auditResults.spf.passed ? 'default' : 'destructive'">
                {{ auditResults.spf.passed ? 'PASS' : 'FAIL' }}
              </UBadge>
            </div>
            <p v-if="auditResults.spf.message" class="text-xs text-gray-600">
              {{ auditResults.spf.message }}
            </p>

            <div class="flex items-center justify-between p-2 bg-white rounded">
              <span class="text-sm font-medium">DKIM</span>
              <UBadge :variant="auditResults.dkim.passed ? 'default' : 'destructive'">
                {{ auditResults.dkim.passed ? 'PASS' : 'FAIL' }}
              </UBadge>
            </div>
            <p v-if="auditResults.dkim.message" class="text-xs text-gray-600">
              {{ auditResults.dkim.message }}
            </p>

            <div class="flex items-center justify-between p-2 bg-white rounded">
              <span class="text-sm font-medium">DMARC</span>
              <UBadge :variant="auditResults.dmarc.passed ? 'default' : 'destructive'">
                {{ auditResults.dmarc.passed ? 'PASS' : 'FAIL' }}
              </UBadge>
            </div>
            <p v-if="auditResults.dmarc.message" class="text-xs text-gray-600">
              {{ auditResults.dmarc.message }}
            </p>

            <div class="flex items-center justify-between p-2 bg-white rounded">
              <span class="text-sm font-medium">rDNS</span>
              <UBadge :variant="auditResults.rdns.passed ? 'default' : 'destructive'">
                {{ auditResults.rdns.passed ? 'PASS' : 'FAIL' }}
              </UBadge>
            </div>
            <p v-if="auditResults.rdns.message" class="text-xs text-gray-600">
              {{ auditResults.rdns.message }}
            </p>

            <div class="flex items-center justify-between p-2 bg-white rounded">
              <span class="text-sm font-medium">FCrDNS</span>
              <UBadge :variant="auditResults.fcrdns.passed ? 'default' : 'destructive'">
                {{ auditResults.fcrdns.passed ? 'PASS' : 'FAIL' }}
              </UBadge>
            </div>
            <p v-if="auditResults.fcrdns.message" class="text-xs text-gray-600">
              {{ auditResults.fcrdns.message }}
            </p>

            <div class="flex items-center justify-between p-2 bg-white rounded">
              <span class="text-sm font-medium">TLS</span>
              <UBadge :variant="auditResults.tls.passed ? 'default' : 'destructive'">
                {{ auditResults.tls.passed ? 'PASS' : 'FAIL' }}
              </UBadge>
            </div>
            <p v-if="auditResults.tls.message" class="text-xs text-gray-600">
              {{ auditResults.tls.message }}
            </p>

            <div class="flex items-center justify-between p-2 bg-white rounded">
              <span class="text-sm font-medium">MTA-STS</span>
              <UBadge :variant="auditResults.mta_sts.passed ? 'default' : 'destructive'">
                {{ auditResults.mta_sts.passed ? 'PASS' : 'FAIL' }}
              </UBadge>
            </div>
            <p v-if="auditResults.mta_sts.message" class="text-xs text-gray-600">
              {{ auditResults.mta_sts.message }}
            </p>
          </div>

          <div class="mt-3">
            <div class="flex items-center justify-between">
              <span class="font-medium">Overall Score:</span>
              <UBadge :variant="getScoreVariant(auditResults.overall_score)">
                {{ auditResults.overall_score }}/100
              </UBadge>
            </div>
            <div v-if="auditResults.issues && auditResults.issues.length > 0" class="mt-2">
              <h4 class="font-medium text-sm">Issues:</h4>
              <ul class="list-disc list-inside text-xs text-red-600">
                <li v-for="(issue, idx) in auditResults.issues" :key="idx">
                  {{ issue }}
                </li>
              </ul>
            </div>
          </div>
        </div>
      </UCardContent>
    </UCard>

    <div class="flex justify-between">
      <UButton variant="outline" :disabled="isSubmitting">
        Previous
      </UButton>
      <UButton @click="handleNext" :disabled="!isFormValid || isSubmitting">
        Next: Create Admin User
        <template v-if="isSubmitting">
          <span class="ml-2">Saving...</span>
        </template>
      </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
// Reactive state
const form = ref({
  domain: ''
})

const auditResults = ref<AuditResults | null>(null)
const isSubmitting = ref(false)
const isValidating = ref(false)
const errors = ref<Record<string, string>>({})

// Form validation
const isFormValid = computed(() => {
  return form.value.domain.trim().length > 0 && auditResults.value !== null
})

// Emit event to parent
const emit = defineEmits<{
  next: [config: typeof form.value & { domain: string, audit_results: typeof auditResults.value }]
  previous: []
}>()

const validateDomain = async () => {
  const domain = form.value.domain.trim()
  if (!domain) {
    errors.value.domain = 'Domain is required'
    return
  }

  isValidating.value = true
  errors.value.domain = ''

  try {
    // Import and use reputation API composable
    const { useReputationApi } = await import('~/composables/api/reputation')
    const { auditDomain } = useReputationApi()
    const results = await auditDomain(domain)
    auditResults.value = results
  } catch (error: any) {
    errors.value.domain = error.message || 'Failed to run DNS audit'
  } finally {
    isValidating.value = false
  }
}

const runAudit = () => {
  validateDomain()
}

const handleNext = () => {
  if (!isFormValid.value) return

  isSubmitting.value = true
  emit('next', {
    domain: form.value.domain,
    audit_results: auditResults.value
  })
}

// Computed helpers
const getScoreVariant = (score: number) => {
  if (score >= 70) return 'default'
  if (score >= 50) return 'secondary'
  return 'destructive'
}

// Watch for errors reset
watch(() => Object.values(errors.value), () => {
  if (Object.keys(errors.value).length > 0) {
    errors.value = {}
  }
})

// TypeScript interfaces
interface CheckStatus {
  passed: boolean
  message: string
  details?: Record<string, any>
}

interface AuditResults {
  domain: string
  timestamp: number
  spf: CheckStatus
  dkim: CheckStatus
  dmarc: CheckStatus
  rdns: CheckStatus
  fcrdns: CheckStatus
  tls: CheckStatus
  mta_sts: CheckStatus
  postmaster_ok: boolean
  abuse_ok: boolean
  overall_score: number
  issues: string[]
}

// Expose reset method to parent
defineExpose({
  reset: () => {
    form.value = { domain: '' }
    auditResults.value = null
    errors.value = {}
  }
})
</script>
