import { defineStore } from 'pinia'
import axios from 'axios'

interface Mailbox {
  id: number
  name: string
  path: string
  unread_count: number
  total_count: number
}

interface Message {
  id: number
  mailbox_id: number
  subject: string
  from: string
  to: string[]
  cc?: string[]
  bcc?: string[]
  date: string
  size: number
  flags: string[]
  thread_id?: string
  body_html?: string
  body_text?: string
  attachments?: Attachment[]
}

interface Attachment {
  id: string
  filename: string
  content_type: string
  size: number
}

interface MailState {
  mailboxes: Mailbox[]
  currentMailbox: Mailbox | null
  messages: Message[]
  currentMessage: Message | null
  loading: boolean
}

export const useMailStore = defineStore('mail', {
  state: (): MailState => ({
    mailboxes: [],
    currentMailbox: null,
    messages: [],
    currentMessage: null,
    loading: false
  }),

  getters: {
    unreadCount: (state) => {
      return state.mailboxes.reduce((sum, box) => sum + box.unread_count, 0)
    },

    inboxMailbox: (state) => {
      return state.mailboxes.find(box => box.name === 'INBOX')
    }
  },

  actions: {
    async fetchMailboxes() {
      try {
        this.loading = true
        const response = await axios.get('/api/v1/webmail/mailboxes')
        this.mailboxes = response.data.mailboxes || []
      } catch (error) {
        console.error('Failed to fetch mailboxes:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchMessages(mailboxId: number, page = 1, limit = 50) {
      try {
        this.loading = true
        const response = await axios.get(`/api/v1/webmail/mailboxes/${mailboxId}/messages`, {
          params: { page, limit }
        })
        this.messages = response.data.messages || []
        this.currentMailbox = this.mailboxes.find(box => box.id === mailboxId) || null
      } catch (error) {
        console.error('Failed to fetch messages:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchMessage(messageId: number) {
      try {
        this.loading = true
        const response = await axios.get(`/api/v1/webmail/messages/${messageId}`)
        this.currentMessage = response.data
        return response.data
      } catch (error) {
        console.error('Failed to fetch message:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async sendMessage(data: any) {
      try {
        this.loading = true
        const response = await axios.post('/api/v1/webmail/messages', data)
        return response.data
      } catch (error) {
        console.error('Failed to send message:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteMessage(messageId: number) {
      try {
        await axios.delete(`/api/v1/webmail/messages/${messageId}`)
        this.messages = this.messages.filter(msg => msg.id !== messageId)
      } catch (error) {
        console.error('Failed to delete message:', error)
        throw error
      }
    },

    async moveMessage(messageId: number, targetMailboxId: number) {
      try {
        await axios.post(`/api/v1/webmail/messages/${messageId}/move`, {
          mailbox_id: targetMailboxId
        })
        this.messages = this.messages.filter(msg => msg.id !== messageId)
      } catch (error) {
        console.error('Failed to move message:', error)
        throw error
      }
    },

    async markAsRead(messageId: number) {
      try {
        await axios.post(`/api/v1/webmail/messages/${messageId}/flags`, {
          flags: ['\\Seen'],
          action: 'add'
        })
        const message = this.messages.find(msg => msg.id === messageId)
        if (message && !message.flags.includes('\\Seen')) {
          message.flags.push('\\Seen')
        }
      } catch (error) {
        console.error('Failed to mark as read:', error)
        throw error
      }
    },

    async markAsUnread(messageId: number) {
      try {
        await axios.post(`/api/v1/webmail/messages/${messageId}/flags`, {
          flags: ['\\Seen'],
          action: 'remove'
        })
        const message = this.messages.find(msg => msg.id === messageId)
        if (message) {
          message.flags = message.flags.filter(flag => flag !== '\\Seen')
        }
      } catch (error) {
        console.error('Failed to mark as unread:', error)
        throw error
      }
    },

    async searchMessages(query: string) {
      try {
        this.loading = true
        const response = await axios.get('/api/v1/webmail/search', {
          params: { q: query }
        })
        this.messages = response.data.messages || []
      } catch (error) {
        console.error('Failed to search messages:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateFlags(messageId: number, action: 'add' | 'remove', flags: string[]) {
      try {
        await axios.post(`/api/v1/webmail/messages/${messageId}/flags`, {
          flags,
          action
        })
      } catch (error) {
        console.error('Failed to update flags:', error)
        throw error
      }
    },

    async saveDraft(data: any) {
      try {
        const response = await axios.post('/api/v1/webmail/drafts', data)
        return response.data
      } catch (error) {
        console.error('Failed to save draft:', error)
        throw error
      }
    },

    async listDrafts() {
      try {
        const response = await axios.get('/api/v1/webmail/drafts')
        return response.data.drafts || []
      } catch (error) {
        console.error('Failed to list drafts:', error)
        throw error
      }
    },

    async getDraft(draftId: number) {
      try {
        const response = await axios.get(`/api/v1/webmail/drafts/${draftId}`)
        return response.data
      } catch (error) {
        console.error('Failed to get draft:', error)
        throw error
      }
    },

    async deleteDraft(draftId: number) {
      try {
        await axios.delete(`/api/v1/webmail/drafts/${draftId}`)
      } catch (error) {
        console.error('Failed to delete draft:', error)
        throw error
      }
    }
  }
})
