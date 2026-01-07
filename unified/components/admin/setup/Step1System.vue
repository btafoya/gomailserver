<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-gray-900">System Configuration</h2>
      <p class="text-gray-600 mt-2">Configure your mail server settings</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Server Settings</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Server Hostname</label>
          <Input
            v-model="form.hostname"
            type="text"
            placeholder="mail.example.com"
            :disabled="isSubmitting"
          />
          <p class="text-xs text-gray-500 mt-1">
            The public hostname of your mail server
          </p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">API Port</label>
          <Input
            v-model="form.port"
            type="number"
            placeholder="8980"
            :disabled="isSubmitting"
          />
          <p class="text-xs text-gray-500 mt-1">
            Port for API and Web UI (default: 8980)
          </p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">SMTP Port</label>
          <Input
            v-model="form.smtp_port"
            type="number"
            placeholder="25"
            :disabled="isSubmitting"
          />
          <p class="text-xs text-gray-500 mt-1">
            Port for SMTP (default: 25)
          </p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">IMAP Port</label>
          <Input
            v-model="form.imap_port"
            type="number"
            placeholder="143"
            :disabled="isSubmitting"
          />
          <p class="text-xs text-gray-500 mt-1">
            Port for IMAP (default: 143)
          </p>
        </div>

        <div class="flex items-center space-x-2">
          <input
            id="tls"
            v-model="form.tls_enabled"
            type="checkbox"
            class="h-4 w-4 text-blue-600 rounded"
            :disabled="isSubmitting"
          />
          <label for="tls" class="text-sm font-medium text-gray-700">
            Enable TLS/SSL
          </label>
          <p class="text-xs text-gray-500 mt-1">
            Enable automatic SSL certificate management
          </p>
        </div>
      </CardContent>
    </Card>

    <div class="flex justify-between">
      <Button variant="outline" :disabled="isSubmitting">
        Previous
      </Button>
      <Button @click="handleNext" :disabled="!isFormValid || isSubmitting">
        Next: Domain Setup
        <template v-if="isSubmitting">
          <span class="ml-2">Configuring...</span>
        </template>
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
// Reactive state
const form = ref({
  hostname: '',
  port: 8980,
  smtp_port: 25,
  imap_port: 143,
  tls_enabled: true
})

const isSubmitting = ref(false)
const errors = ref<Record<string, string>>({})

// Form validation
const isFormValid = computed(() => {
  return form.value.hostname.trim().length > 0 &&
         form.value.port > 0 &&
         form.value.port < 65535 &&
         form.value.smtp_port > 0 &&
         form.value.smtp_port < 65535 &&
         form.value.imap_port > 0 &&
         form.value.imap_port < 65535
})

// Emit event to parent
const emit = defineEmits<{
  next: [config: typeof form.value]
  previous: []
}>()

const handleNext = () => {
  if (!isFormValid.value) return

  isSubmitting.value = true
  emit('next', form.value)
}

// Watch for errors and reset when form changes
watch(() => Object.values(errors.value), () => {
  if (Object.keys(errors.value).length > 0) {
    errors.value = {}
  }
})

// Expose reset method to parent
defineExpose({
  reset: () => {
    form.value = {
      hostname: '',
      port: 8980,
      smtp_port: 25,
      imap_port: 143,
      tls_enabled: true
    }
    errors.value = {}
  }
})
</script>
