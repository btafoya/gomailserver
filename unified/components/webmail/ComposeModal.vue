<template>
  <UCard class="w-full">
    <template #header>
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold">{{ isReply ? 'Reply' : isForward ? 'Forward' : 'New Message' }}</h3>
        <UButton
          variant="ghost"
          size="sm"
          icon
          @click="$emit('close')"
        >
          <X class="h-4 w-4" />
        </UButton>
      </div>
    </template>

    <div class="space-y-4">
      <!-- To Field -->
      <div>
        <label class="mb-1.5 block text-sm font-medium text-muted-foreground">To</label>
        <UInput
          v-model="to"
          placeholder="recipient@example.com"
          :disabled="sending"
        />
      </div>

      <!-- CC Field (expandable) -->
      <div v-if="showCc">
        <label class="mb-1.5 block text-sm font-medium text-muted-foreground">CC</label>
        <UInput
          v-model="cc"
          placeholder="cc@example.com"
          :disabled="sending"
        />
      </div>

      <!-- BCC Field (expandable) -->
      <div v-if="showBcc">
        <label class="mb-1.5 block text-sm font-medium text-muted-foreground">BCC</label>
        <UInput
          v-model="bcc"
          placeholder="bcc@example.com"
          :disabled="sending"
        />
      </div>

      <!-- Toggle CC/BCC -->
      <div v-if="!showCc || !showBcc" class="flex space-x-2">
        <UButton
          v-if="!showCc"
          variant="ghost"
          size="sm"
          @click="showCc = true"
        >
          Add CC
        </UButton>
        <UButton
          v-if="!showBcc"
          variant="ghost"
          size="sm"
          @click="showBcc = true"
        >
          Add BCC
        </UButton>
      </div>

      <!-- Subject Field -->
      <div>
        <label class="mb-1.5 block text-sm font-medium text-muted-foreground">Subject</label>
        <UInput
          v-model="subject"
          placeholder="Subject"
          :disabled="sending"
        />
      </div>

      <!-- Body Field -->
      <div>
        <label class="mb-1.5 block text-sm font-medium text-muted-foreground">Message</label>
        <UTextarea
          v-model="body"
          placeholder="Write your message..."
          :rows="12"
          :disabled="sending"
          class="w-full"
        />
      </div>

      <!-- Attachments -->
      <div v-if="attachments.length > 0">
        <label class="mb-1.5 block text-sm font-medium text-muted-foreground">Attachments</label>
        <div class="flex flex-wrap gap-2">
          <div
            v-for="(attachment, index) in attachments"
            :key="index"
            class="flex items-center rounded-lg bg-muted px-3 py-1"
          >
            <Paperclip class="mr-2 h-3 w-3" />
            <span class="text-sm">{{ attachment.name }}</span>
            <UButton
              variant="ghost"
              size="sm"
              icon
              class="ml-2"
              @click="removeAttachment(index)"
            >
              <X class="h-3 w-3" />
            </UButton>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <UButton
            variant="outline"
            size="sm"
            @click="addAttachment"
            :disabled="sending"
          >
            <Paperclip class="mr-2 h-4 w-4" />
            Attach
          </UButton>
        </div>

        <div class="flex items-center space-x-2">
          <UButton
            variant="outline"
            @click="handleSaveDraft"
            :loading="saving"
            :disabled="sending"
          >
            Save Draft
          </UButton>
          <UButton
            @click="handleSend"
            :loading="sending"
            :disabled="!canSend"
          >
            <Send class="mr-2 h-4 w-4" />
            Send
          </UButton>
        </div>
      </div>
    </template>
  </UCard>
</template>

<script setup lang="ts">
import { X, Paperclip, Send } from 'lucide-vue-next'

interface ComposeData {
  to?: string[]
  cc?: string[]
  bcc?: string[]
  subject?: string
  body_text?: string
  body_html?: string
  in_reply_to?: string
  references?: string
}

interface Props {
  composeData?: ComposeData
  sending?: boolean
  saving?: boolean
}

interface Emits {
  (e: 'send', data: ComposeData): void
  (e: 'save-draft', data: ComposeData): void
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  composeData: () => ({}),
  sending: false,
  saving: false
})

const emit = defineEmits<Emits>()

// Form state
const to = ref(props.composeData?.to?.join(', ') || '')
const cc = ref(props.composeData?.cc?.join(', ') || '')
const bcc = ref(props.composeData?.bcc?.join(', ') || '')
const subject = ref(props.composeData?.subject || '')
const body = ref(props.composeData?.body_html || props.composeData?.body_text || '')
const attachments = ref<Array<{ name: string; file: File }>>([])

const showCc = ref(!!props.composeData?.cc?.length)
const showBcc = ref(!!props.composeData?.bcc?.length)

// Computed
const isReply = computed(() => !!props.composeData?.in_reply_to)
const isForward = computed(() => subject.value.startsWith('Fwd:'))

const canSend = computed(() => {
  return to.value.trim() && subject.value.trim()
})

// Methods
const parseRecipients = (value: string): string[] => {
  return value
    .split(/[,;]/)
    .map(email => email.trim())
    .filter(email => email.length > 0)
}

const getComposeData = (): ComposeData => {
  return {
    to: parseRecipients(to.value),
    cc: showCc.value ? parseRecipients(cc.value) : [],
    bcc: showBcc.value ? parseRecipients(bcc.value) : [],
    subject: subject.value,
    body_html: body.value,
    body_text: body.value.replace(/<[^>]*>/g, ''),
    in_reply_to: props.composeData?.in_reply_to,
    references: props.composeData?.references
  }
}

const handleSend = () => {
  if (canSend.value) {
    emit('send', getComposeData())
  }
}

const handleSaveDraft = () => {
  emit('save-draft', getComposeData())
}

const addAttachment = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.multiple = true
  input.onchange = (event) => {
    const files = (event.target as HTMLInputElement).files
    if (files) {
      for (const file of files) {
        attachments.value.push({ name: file.name, file })
      }
    }
  }
  input.click()
}

const removeAttachment = (index: number) => {
  attachments.value.splice(index, 1)
}

// Watch for prop changes
watch(() => props.composeData, (newData) => {
  if (newData) {
    to.value = newData.to?.join(', ') || ''
    cc.value = newData.cc?.join(', ') || ''
    bcc.value = newData.bcc?.join(', ') || ''
    subject.value = newData.subject || ''
    body.value = newData.body_html || newData.body_text || ''
    showCc.value = !!newData.cc?.length
    showBcc.value = !!newData.bcc?.length
  }
}, { deep: true })
</script>
