<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMailStore } from '@/stores/mail'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'

const props = defineProps({
  replyTo: {
    type: Number,
    default: undefined
  },
  forward: {
    type: Number,
    default: undefined
  }
})

const mailStore = useMailStore()
const router = useRouter()

const to = ref('')
const cc = ref('')
const bcc = ref('')
const subject = ref('')
const attachments = ref([])
const sending = ref(false)
const showCc = ref(false)
const showBcc = ref(false)
const isPlainText = ref(false)
const draftId = ref(null)
const lastSaved = ref(null)
const autoSaveTimer = ref(null)

const editor = useEditor({
  extensions: [
    StarterKit,
    Placeholder.configure({
      placeholder: 'Write your message...'
    })
  ],
  editorProps: {
    attributes: {
      class: 'prose dark:prose-invert max-w-none focus:outline-none min-h-[200px] px-4 py-3'
    }
  }
})

// Handle reply/forward
onMounted(async () => {
  if (props.replyTo) {
    const message = await mailStore.fetchMessage(props.replyTo)
    to.value = message.from
    subject.value = `Re: ${message.subject}`
    editor.value?.commands.setContent(`<p><br></p><p>On ${new Date(message.date).toLocaleString()}, ${message.from} wrote:</p><blockquote>${message.body_html || message.body_text}</blockquote>`)
  } else if (props.forward) {
    const message = await mailStore.fetchMessage(props.forward)
    subject.value = `Fwd: ${message.subject}`
    editor.value?.commands.setContent(`<p><br></p><p>---------- Forwarded message ---------</p><p>From: ${message.from}<br>Date: ${new Date(message.date).toLocaleString()}<br>Subject: ${message.subject}</p><p>${message.body_html || message.body_text}</p>`)
  }

  // Start auto-save
  startAutoSave()
})

onUnmounted(() => {
  stopAutoSave()
})

const handleFileSelect = (event) => {
  const target = event.target
  if (target.files) {
    attachments.value.push(...Array.from(target.files))
  }
}

const removeAttachment = (index) => {
  attachments.value.splice(index, 1)
}

const togglePlainText = () => {
  isPlainText.value = !isPlainText.value
}

const sendEmail = async () => {
  try {
    sending.value = true

    const formData = new FormData()
    formData.append('to', to.value)
    if (cc.value) formData.append('cc', cc.value)
    if (bcc.value) formData.append('bcc', bcc.value)
    formData.append('subject', subject.value)

    if (isPlainText.value) {
      formData.append('body_text', editor.value?.getText() || '')
    } else {
      formData.append('body_html', editor.value?.getHTML() || '')
    }

    attachments.value.forEach(file => {
      formData.append('attachments', file)
    })

    await mailStore.sendMessage(formData)
    router.back()
  } catch (error) {
    console.error('Failed to send email:', error)
    alert('Failed to send email')
  } finally {
    sending.value = false
  }
}

const saveDraft = async () => {
  try {
    const draftData = {
      draft_id: draftId.value,
      to: to.value.split(',').map(e => e.trim()).filter(Boolean),
      cc: cc.value ? cc.value.split(',').map(e => e.trim()).filter(Boolean) : [],
      bcc: bcc.value ? bcc.value.split(',').map(e => e.trim()).filter(Boolean) : [],
      subject: subject.value,
      body_html: isPlainText.value ? '' : editor.value?.getHTML() || '',
      body_text: isPlainText.value ? editor.value?.getText() || '' : ''
    }

    const saved = await mailStore.saveDraft(draftData)
    draftId.value = saved.id
    lastSaved.value = new Date()
    console.log('Draft saved', saved.id)
  } catch (error) {
    console.error('Failed to save draft:', error)
  }
}

// Auto-save draft every 30 seconds
const startAutoSave = () => {
  autoSaveTimer.value = setInterval(() => {
    if (to.value || subject.value || editor.value?.getText()) {
      saveDraft()
    }
  }, 30000) // 30 seconds
}

const stopAutoSave = () => {
  if (autoSaveTimer.value) {
    clearInterval(autoSaveTimer.value)
    autoSaveTimer.value = null
  }
}

const discard = () => {
  if (confirm('Discard this message?')) {
    router.back()
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>

<template>
  <div class="flex-1 flex flex-col bg-background">
    <!-- Composer Header -->
    <div class="border-b bg-card px-6 py-4">
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-semibold">New Message</h2>
        <div class="flex items-center gap-2">
          <button
            @click="saveDraft"
            class="px-4 py-2 text-sm hover:bg-muted rounded-md"
          >
            Save Draft
          </button>
          <button
            @click="discard"
            class="px-4 py-2 text-sm hover:bg-muted rounded-md"
          >
            Discard
          </button>
          <button
            @click="sendEmail"
            :disabled="sending || !to || !subject"
            class="px-4 py-2 text-sm bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <Icon v-if="sending" name="lucide:loader-circle" class="w-4 h-4 animate-spin" />
            <Icon v-else name="lucide:send" class="w-4 h-4" />
            {{ sending ? 'Sending...' : 'Send' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Composer Body -->
    <div class="flex-1 overflow-y-auto">
      <div class="max-w-4xl mx-auto py-6 space-y-4">
        <!-- To -->
        <div class="flex items-center gap-3">
          <label class="w-16 text-sm font-medium text-right">To:</label>
          <input
            v-model="to"
            type="email"
            multiple
            class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="recipient@example.com"
          />
          <button
            @click="showCc = !showCc"
            class="text-sm text-primary hover:underline"
          >
            Cc
          </button>
          <button
            @click="showBcc = !showBcc"
            class="text-sm text-primary hover:underline"
          >
            Bcc
          </button>
        </div>

        <!-- Cc -->
        <div v-if="showCc" class="flex items-center gap-3">
          <label class="w-16 text-sm font-medium text-right">Cc:</label>
          <input
            v-model="cc"
            type="email"
            multiple
            class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="cc@example.com"
          />
        </div>

        <!-- Bcc -->
        <div v-if="showBcc" class="flex items-center gap-3">
          <label class="w-16 text-sm font-medium text-right">Bcc:</label>
          <input
            v-model="bcc"
            type="email"
            multiple
            class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="bcc@example.com"
          />
        </div>

        <!-- Subject -->
        <div class="flex items-center gap-3">
          <label class="w-16 text-sm font-medium text-right">Subject:</label>
          <input
            v-model="subject"
            type="text"
            class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="Email subject"
          />
        </div>

        <!-- Attachments -->
        <div v-if="attachments.length" class="flex items-start gap-3">
          <div class="w-16"></div>
          <div class="flex-1 space-y-2">
            <div
              v-for="(file, index) in attachments"
              :key="index"
              class="flex items-center gap-2 px-3 py-2 bg-muted rounded-md"
            >
              <Icon name="lucide:file" class="w-4 h-4" />
              <span class="flex-1 text-sm">{{ file.name }}</span>
              <span class="text-xs text-muted-foreground">{{ formatFileSize(file.size) }}</span>
              <button
                @click="removeAttachment(index)"
                class="p-1 hover:bg-background rounded"
              >
                <Icon name="lucide:x" class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- Editor Toolbar -->
        <div class="border-t border-b bg-muted/30 px-3 py-2 flex items-center justify-between">
          <div class="flex items-center gap-1">
            <button
              @click="editor?.chain().focus().toggleBold().run()"
              :class="{ 'bg-muted': editor?.isActive('bold') }"
              class="p-2 hover:bg-muted rounded"
              title="Bold"
            >
              <Icon name="lucide:bold" class="w-4 h-4" />
            </button>
            <button
              @click="editor?.chain().focus().toggleItalic().run()"
              :class="{ 'bg-muted': editor?.isActive('italic') }"
              class="p-2 hover:bg-muted rounded"
              title="Italic"
            >
              <Icon name="lucide:italic" class="w-4 h-4" />
            </button>
            <button
              @click="editor?.chain().focus().toggleBulletList().run()"
              :class="{ 'bg-muted': editor?.isActive('bulletList') }"
              class="p-2 hover:bg-muted rounded"
              title="Bullet List"
            >
              <Icon name="lucide:list" class="w-4 h-4" />
            </button>
            <button
              @click="editor?.chain().focus().toggleOrderedList().run()"
              :class="{ 'bg-muted': editor?.isActive('orderedList') }"
              class="p-2 hover:bg-muted rounded"
              title="Numbered List"
            >
              <Icon name="lucide:list-ordered" class="w-4 h-4" />
            </button>

            <div class="w-px h-6 bg-border mx-2"></div>

            <label class="p-2 hover:bg-muted rounded cursor-pointer" title="Attach File">
              <Icon name="lucide:paperclip" class="w-4 h-4" />
              <input
                type="file"
                multiple
                @change="handleFileSelect"
                class="hidden"
              />
            </label>
          </div>

          <button
            @click="togglePlainText"
            class="text-sm text-muted-foreground hover:text-foreground"
          >
            {{ isPlainText ? 'Rich Text' : 'Plain Text' }}
          </button>
        </div>

        <!-- Editor -->
        <div class="border rounded-md bg-card">
          <EditorContent v-if="!isPlainText" :editor="editor" />
          <textarea
            v-else
            v-model="subject"
            class="w-full min-h-[200px] px-4 py-3 focus:outline-none"
            placeholder="Write your message..."
          ></textarea>
        </div>
      </div>
    </div>
  </div>
</template>
