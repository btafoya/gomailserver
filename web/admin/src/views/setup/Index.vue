<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'

const router = useRouter()
const currentStep = ref('welcome')
const completedSteps = ref([])
const loading = ref(false)
const error = ref(null)

// Form data
const systemConfig = ref({
  hostname: '',
  listen_smtp: '25',
  listen_imap: '143',
  listen_api: '8980'
})

const domainConfig = ref({
  name: '',
  catchall_email: ''
})

const adminConfig = ref({
  email: '',
  password: '',
  password_confirm: '',
  full_name: ''
})

const tlsConfig = ref({
  acme_email: '',
  enable_acme: true,
  production: false
})

const steps = [
  { id: 'welcome', title: 'Welcome', description: 'Initial setup for gomailserver' },
  { id: 'system', title: 'System Configuration', description: 'Configure system settings' },
  { id: 'domain', title: 'First Domain', description: 'Set up your first domain' },
  { id: 'admin', title: 'Admin Account', description: 'Create an administrator' },
  { id: 'tls', title: 'TLS Certificates', description: 'Configure ACME/Let\'s Encrypt' },
  { id: 'complete', title: 'Complete', description: 'Setup finished' }
]

const currentStepIndex = computed(() =>
  steps.findIndex(s => s.id === currentStep.value)
)

const canGoNext = computed(() => {
  switch (currentStep.value) {
    case 'welcome':
      return true
    case 'system':
      return systemConfig.value.hostname.length > 0
    case 'domain':
      return domainConfig.value.name.length > 0
    case 'admin':
      return adminConfig.value.email &&
             adminConfig.value.password &&
             adminConfig.value.password === adminConfig.value.password_confirm
    case 'tls':
      return !tlsConfig.value.enable_acme || tlsConfig.value.acme_email.length > 0
    default:
      return false
  }
})

onMounted(async () => {
  try {
    const response = await api.get('/api/v1/setup/status')
    currentStep.value = response.data.current_step
    completedSteps.value = response.data.completed_steps || []
  } catch (err) {
    console.error('Failed to load setup status:', err)
  }
})

const nextStep = async () => {
  loading.value = true
  error.value = null

  try {
    const stepData = getStepData(currentStep.value)

    await api.post('/api/v1/setup/step', {
      step: currentStep.value,
      data: stepData
    })

    if (!completedSteps.value.includes(currentStep.value)) {
      completedSteps.value.push(currentStep.value)
    }

    if (currentStepIndex.value < steps.length - 1) {
      currentStep.value = steps[currentStepIndex.value + 1].id
    }

    if (currentStep.value === 'complete') {
      await completeSetup()
    }
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to save step'
  } finally {
    loading.value = false
  }
}

const prevStep = () => {
  if (currentStepIndex.value > 0) {
    currentStep.value = steps[currentStepIndex.value - 1].id
  }
}

const getStepData = (step) => {
  switch (step) {
    case 'system':
      return systemConfig.value
    case 'domain':
      return domainConfig.value
    case 'admin':
      return adminConfig.value
    case 'tls':
      return tlsConfig.value
    default:
      return {}
  }
}

const completeSetup = async () => {
  try {
    await api.post('/api/v1/setup/complete')
    setTimeout(() => {
      router.push({ name: 'Dashboard' })
    }, 2000)
  } catch (err) {
    error.value = 'Failed to complete setup'
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-3xl mx-auto">
      <!-- Progress Steps -->
      <div class="mb-8">
        <nav aria-label="Progress">
          <ol class="flex items-center">
            <li v-for="(step, index) in steps" :key="step.id"
                :class="[
                  'relative',
                  index !== steps.length - 1 ? 'pr-8 sm:pr-20 flex-1' : ''
                ]">
              <div v-if="index !== steps.length - 1"
                   class="absolute inset-0 flex items-center"
                   aria-hidden="true">
                <div :class="[
                  'h-0.5 w-full',
                  completedSteps.includes(step.id) ? 'bg-blue-600' : 'bg-gray-200'
                ]"></div>
              </div>
              <div class="relative flex items-center justify-center">
                <span :class="[
                  'h-8 w-8 rounded-full flex items-center justify-center',
                  currentStep === step.id ? 'bg-blue-600 text-white' :
                  completedSteps.includes(step.id) ? 'bg-blue-600 text-white' :
                  'bg-gray-200 text-gray-500'
                ]">
                  {{ index + 1 }}
                </span>
              </div>
            </li>
          </ol>
        </nav>
      </div>

      <!-- Step Content -->
      <Card>
        <CardHeader>
          <CardTitle>{{ steps[currentStepIndex].title }}</CardTitle>
          <CardDescription>{{ steps[currentStepIndex].description }}</CardDescription>
        </CardHeader>
        <CardContent>
          <Alert v-if="error" variant="destructive" class="mb-4">
            <AlertDescription>{{ error }}</AlertDescription>
          </Alert>

          <!-- Welcome Step -->
          <div v-if="currentStep === 'welcome'" class="space-y-4">
            <p>Welcome to gomailserver! This wizard will guide you through the initial setup.</p>
            <p>You will configure:</p>
            <ul class="list-disc list-inside space-y-2 ml-4">
              <li>System settings (hostname, ports)</li>
              <li>Your first mail domain</li>
              <li>An administrator account</li>
              <li>TLS certificates via ACME/Let's Encrypt</li>
            </ul>
          </div>

          <!-- System Configuration -->
          <div v-if="currentStep === 'system'" class="space-y-4">
            <div>
              <Label for="hostname">Server Hostname</Label>
              <Input id="hostname" v-model="systemConfig.hostname" placeholder="mail.example.com" />
            </div>
            <div class="grid grid-cols-3 gap-4">
              <div>
                <Label for="smtp_port">SMTP Port</Label>
                <Input id="smtp_port" v-model="systemConfig.listen_smtp" type="number" />
              </div>
              <div>
                <Label for="imap_port">IMAP Port</Label>
                <Input id="imap_port" v-model="systemConfig.listen_imap" type="number" />
              </div>
              <div>
                <Label for="api_port">API Port</Label>
                <Input id="api_port" v-model="systemConfig.listen_api" type="number" />
              </div>
            </div>
          </div>

          <!-- Domain Configuration -->
          <div v-if="currentStep === 'domain'" class="space-y-4">
            <div>
              <Label for="domain_name">Domain Name</Label>
              <Input id="domain_name" v-model="domainConfig.name" placeholder="example.com" />
            </div>
            <div>
              <Label for="catchall">Catchall Email (optional)</Label>
              <Input id="catchall" v-model="domainConfig.catchall_email" placeholder="admin@example.com" />
            </div>
          </div>

          <!-- Admin Account -->
          <div v-if="currentStep === 'admin'" class="space-y-4">
            <div>
              <Label for="admin_name">Full Name</Label>
              <Input id="admin_name" v-model="adminConfig.full_name" />
            </div>
            <div>
              <Label for="admin_email">Email Address</Label>
              <Input id="admin_email" v-model="adminConfig.email" type="email" />
            </div>
            <div>
              <Label for="admin_password">Password</Label>
              <Input id="admin_password" v-model="adminConfig.password" type="password" />
            </div>
            <div>
              <Label for="admin_password_confirm">Confirm Password</Label>
              <Input id="admin_password_confirm" v-model="adminConfig.password_confirm" type="password" />
            </div>
            <p v-if="adminConfig.password && adminConfig.password !== adminConfig.password_confirm"
               class="text-sm text-red-600">
              Passwords do not match
            </p>
          </div>

          <!-- TLS Configuration -->
          <div v-if="currentStep === 'tls'" class="space-y-4">
            <div class="flex items-center space-x-2">
              <input
                id="enable_acme"
                v-model="tlsConfig.enable_acme"
                type="checkbox"
                class="rounded border-gray-300"
              />
              <Label for="enable_acme">Enable ACME/Let's Encrypt</Label>
            </div>
            <div v-if="tlsConfig.enable_acme">
              <Label for="acme_email">ACME Email</Label>
              <Input id="acme_email" v-model="tlsConfig.acme_email" type="email" />
            </div>
            <div v-if="tlsConfig.enable_acme" class="flex items-center space-x-2">
              <input
                id="production"
                v-model="tlsConfig.production"
                type="checkbox"
                class="rounded border-gray-300"
              />
              <Label for="production">Use Production (not staging)</Label>
            </div>
          </div>

          <!-- Complete -->
          <div v-if="currentStep === 'complete'" class="space-y-4 text-center">
            <svg class="mx-auto h-12 w-12 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <h3 class="text-lg font-medium">Setup Complete!</h3>
            <p>Your mail server is now configured and ready to use.</p>
            <p class="text-sm text-gray-600">Redirecting to dashboard...</p>
          </div>

          <!-- Navigation -->
          <div class="flex justify-between mt-8">
            <Button
              v-if="currentStepIndex > 0 && currentStep !== 'complete'"
              @click="prevStep"
              variant="outline"
            >
              Previous
            </Button>
            <div v-else></div>

            <Button
              v-if="currentStep !== 'complete'"
              @click="nextStep"
              :disabled="!canGoNext || loading"
            >
              {{ currentStepIndex === steps.length - 2 ? 'Complete Setup' : 'Next' }}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
