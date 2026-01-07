<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-900 mb-6">Domain Deliverability Audit</h1>
      <p class="text-gray-600 mt-2">
        Run comprehensive deliverability audits for your domains to identify reputation risks
      </p>
    </div>

    <!-- Audit Input -->
    <UCard class="mb-6">
      <UCardHeader>
        <UCardTitle>Run Audit</UCardTitle>
      </UCardHeader>
      <UCardContent class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700">Domain to Audit</label>
          <UInput
            v-model="domainInput"
            type="text"
            placeholder="example.com"
            :disabled="isAuditing || isSubmitting"
          />
          <p class="text-xs text-gray-500 mt-1">
            Enter the domain you want to audit
          </p>
        </div>

        <div>
          <UButton
            @click="runAudit"
            :disabled="!auditDomain || isAuditing || isSubmitting"
            :loading="isAuditing && !isSubmitting"
            size="lg"
          >
            <template v-if="!isAuditing && !isSubmitting">
              <Shield class="mr-2 h-5 w-5" />
              Run Audit
            </template>
            <template v-else>
              <Loader2 class="mr-2 h-5 w-5 animate-spin" />
              Auditing...
            </template>
          </UButton>
        </div>
      </UCardContent>
    </UCard>

    <!-- Audit Results -->
    <div v-if="auditResults" class="mt-6">
      <UCard>
        <UCardHeader>
          <UCardTitle>Audit Results for {{ auditResults.domain }}</UCardTitle>
        </UCardHeader>
        <UCardContent>
          <!-- Overall Score Card -->
          <div class="mb-4">
            <div class="flex items-center justify-between">
              <span class="text-2xl font-bold">Overall Deliverability Score</span>
              <UBadge :variant="getScoreVariant(auditResults.overall_score)">
                {{ auditResults.overall_score }}/100
              </UBadge>
            </div>
          </div>

          <!-- DNS Configuration Checks -->
          <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h3 class="text-lg font-semibold text-gray-800 mb-3">Authentication Status</h3>

              <div class="space-y-3">
                <div class="flex items-center justify-between p-2 bg-white rounded border">
                  <span class="text-sm font-medium">SPF</span>
                  <UBadge :variant="auditResults.spf.passed ? 'default' : 'destructive'">
                    {{ auditResults.spf.passed ? 'PASS' : 'FAIL' }}
                  </UBadge>
                </div>
                <p v-if="auditResults.spf.message" class="text-xs text-gray-600 mt-1">
                  {{ auditResults.spf.message }}
                </p>

                <div class="flex items-center justify-between p-2 bg-white rounded border">
                  <span class="text-sm font-medium">DKIM</span>
                  <UBadge :variant="auditResults.dkim.passed ? 'default' : 'destructive'">
                    {{ auditResults.dkim.passed ? 'PASS' : 'FAIL' }}
                  </UBadge>
                </div>
                <p v-if="auditResults.dkim.message" class="text-xs text-gray-600 mt-1">
                  {{ auditResults.dkim.message }}
                </p>

                <div class="flex items-center justify-between p-2 bg-white rounded border">
                  <span class="text-sm font-medium">DMARC</span>
                  <UBadge :variant="auditResults.dmarc.passed ? 'default' : 'destructive'">
                    {{ auditResults.dmarc.passed ? 'PASS' : 'FAIL' }}
                  </UBadge>
                </div>
                <p v-if="auditResults.dmarc.message" class="text-xs text-gray-600 mt-1">
                  {{ auditResults.dmarc.message }}
                </p>

                <div class="flex items-center justify-between p-2 bg-white rounded border">
                  <span class="text-sm font-medium">rDNS</span>
                  <UBadge :variant="auditResults.rdns.passed ? 'default' : 'destructive'">
                    {{ auditResults.rdns.passed ? 'PASS' : 'FAIL' }}
                  </UBadge>
                </div>
                <p v-if="auditResults.rdns.message" class="text-xs text-gray-600 mt-1">
                  {{ auditResults.rdns.message }}
                </p>

                <div class="flex items-center justify-between p-2 bg-white rounded border">
                  <span class="text-sm font-medium">FCrDNS</span>
                  <UBadge :variant="auditResults.fcrdns.passed ? 'default' : 'destructive'">
                    {{ auditResults.fcrdns.passed ? 'PASS' : 'FAIL' }}
                  </UBadge>
                </div>
                <p v-if="auditResults.fcrdns.message" class="text-xs text-gray-600 mt-1">
                  {{ auditResults.fcrdns.message }}
                </p>
              </div>
            </div>

            <!-- Operational Checks -->
            <div>
              <h3 class="text-lg font-semibold text-gray-800 mb-3">Operational Status</h3>

              <div class="space-y-3">
                <div class="flex items-center justify-between p-2 bg-white rounded border">
                  <span class="text-sm font-medium">TLS Certificate</span>
                  <UBadge :variant="auditResults.tls.passed ? 'default' : 'destructive'">
                    {{ auditResults.tls.passed ? 'Valid' : 'Invalid/Expired' }}
                  </UBadge>
                </div>
                <p v-if="auditResults.tls.message" class="text-xs text-gray-600 mt-1">
                  {{ auditResults.tls.message }}
                </p>

                <div class="flex items-center justify-between p-2 bg-green-50 rounded border">
                  <span class="text-sm font-medium text-gray-700">postmaster@</span>
                  <UBadge :variant="auditResults.postmaster_ok ? 'default' : 'destructive'">
                    {{ auditResults.postmaster_ok ? 'OK' : 'MISSING' }}
                  </UBadge>
                </div>

                <div class="flex items-center justify-between p-2 bg-red-50 rounded border">
                  <span class="text-sm font-medium text-gray-700">abuse@</span>
                  <UBadge :variant="auditResults.abuse_ok ? 'default' : 'destructive'">
                    {{ auditResults.abuse_ok ? 'OK' : 'MISSING' }}
                  </UBadge>
                </div>

                <div v-if="auditResults.mta_sts" class="flex items-center justify-between p-2 bg-white rounded border">
                  <span class="text-sm font-medium">MTA-STS</span>
                  <UBadge :variant="auditResults.mta_sts.passed ? 'default' : 'destructive'">
                    {{ auditResults.mta_sts.passed ? 'PASS' : 'FAIL' }}
                  </UBadge>
                </div>
                <p v-if="auditResults.mta_sts.message" class="text-xs text-gray-600 mt-1">
                  {{ auditResults.mta_sts.message }}
                </p>
              </div>
            </div>
          </div>

          <!-- Issues Summary -->
          <div class="mt-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-3">Issues Summary</h3>
            <UAlert v-if="auditResults.issues && auditResults.issues.length > 0" variant="destructive">
              Found {{ auditResults.issues.length }} issue(s) that need attention
            </UAlert>
            <div v-else class="text-green-600">
              All checks passed! Your domain is properly configured.
            </div>

            <ul v-if="auditResults.issues && auditResults.issues.length > 0" class="list-disc list-decimal mt-2 space-y-1">
              <li v-for="(issue, idx) in auditResults.issues" :key="idx">
                <span>{{ issue }}</span>
              </li>
            </ul>
          </div>

          <!-- Export Button -->
          <div class="mt-6">
            <UButton @click="exportResults" variant="outline" :disabled="isSubmitting">
              <Download class="mr-2 h-4 w-4" />
              Export Results
            </UButton>
          </div>
        </UCardContent>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Shield, Loader2, Download } from 'lucide-vue-next'
import { useReputationApi } from '~/composables/api/reputation'

const { auditDomain: auditDomainApi } = useReputationApi()

// Reactive state
const domainInput = ref('')
const auditResults = ref<any>(null)
const isAuditing = ref(false)
const isSubmitting = ref(false)
const error = ref<string | null>(null)

const runAudit = async () => {
  if (!domainInput.value) {
    error.value = 'Please enter a domain to audit'
    return
  }

  isAuditing.value = true
  isSubmitting.value = true
  error.value = null

  try {
    const results = await auditDomainApi(domainInput.value)
    auditResults.value = results
  } catch (err: any) {
    error.value = err.message || 'Audit failed'
    isAuditing.value = false
    isSubmitting.value = false
  } finally {
    isAuditing.value = false
  }
}

const exportResults = () => {
  // Export results as PDF (client-side for now)
  const dataStr = JSON.stringify(auditResults.value, null, 2)
  const blob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `audit-${domainInput.value}.json`
  link.click()
}

const getScoreVariant = (score: number) => {
  if (score >= 70) return 'default'
  if (score >= 50) return 'secondary'
  return 'destructive'
}
</script>
