import { defineStore } from 'pinia'
import api from '@/api/axios'

export const useMailStore = defineStore('mail', {
  state: () => ({
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
        const response = await api.get('/api/v1/webmail/mailboxes')
        this.mailboxes = response.data.mailboxes || []
      } catch (error) {
        console.error('Failed to fetch mailboxes:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchMessages(mailboxId, page = 1, limit = 50) {
      try {
        this.loading = true
        const response = await api.get(`/api/v1/webmail/mailboxes/${mailboxId}/messages`, {
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

    async fetchMessage(messageId) {
      try {
        this.loading = true
        const response = await api.get(`/api/v1/webmail/messages/${messageId}`)
        this.currentMessage = response.data
        return response.data
      } catch (error) {
        console.error('Failed to fetch message:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async sendMessage(data) {
      try {
        this.loading = true
        const response = await api.post('/api/v1/webmail/messages', data)
        return response.data
      } catch (error) {
        console.error('Failed to send message:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteMessage(messageId) {
      try {
        await api.delete(`/api/v1/webmail/messages/${messageId}`)
        this.messages = this.messages.filter(msg => msg.id !== messageId)
      } catch (error) {
        console.error('Failed to delete message:', error)
        throw error
      }
    },

    async moveMessage(messageId, targetMailboxId) {
      try {
        await api.post(`/api/v1/webmail/messages/${messageId}/move`, {
          mailbox_id: targetMailboxId
        })
        this.messages = this.messages.filter(msg => msg.id !== messageId)
      } catch (error) {
        console.error('Failed to move message:', error)
        throw error
      }
    },

    async markAsRead(messageId) {
      try {
        await api.post(`/api/v1/webmail/messages/${messageId}/flags`, {
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

    async markAsUnread(messageId) {
      try {
        await api.post(`/api/v1/webmail/messages/${messageId}/flags`, {
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

    async searchMessages(query) {
      try {
        this.loading = true
        const response = await api.get('/api/v1/webmail/search', {
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

    async updateFlags(messageId, action, flags) {
      try {
        await api.post(`/api/v1/webmail/messages/${messageId}/flags`, {
          flags,
          action
        })
      } catch (error) {
        console.error('Failed to update flags:', error)
        throw error
      }
    },

    async saveDraft(data) {
      try {
        const response = await api.post('/api/v1/webmail/drafts', data)
        return response.data
      } catch (error) {
        console.error('Failed to save draft:', error)
        throw error
      }
    },

    async listDrafts() {
      try {
        const response = await api.get('/api/v1/webmail/drafts')
        return response.data.drafts || []
      } catch (error) {
        console.error('Failed to list drafts:', error)
        throw error
      }
    },

    async getDraft(draftId) {
      try {
        const response = await api.get(`/api/v1/webmail/drafts/${draftId}`)
        return response.data
      } catch (error) {
        console.error('Failed to get draft:', error)
        throw error
      }
    },

    async deleteDraft(draftId) {
      try {
        await api.delete(`/api/v1/webmail/drafts/${draftId}`)
      } catch (error) {
        console.error('Failed to delete draft:', error)
        throw error
      }
    }
  }
})
