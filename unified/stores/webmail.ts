/**
 * Webmail store for managing email state, mailboxes, messages, and composition
 */

import { defineStore } from 'pinia'
import { useApiBase } from '../composables/useApiBase'

// Types
export interface Mailbox {
  id: number
  name: string
  display_name: string
  message_count: number
  unread_count: number
  type: 'inbox' | 'sent' | 'drafts' | 'trash' | 'custom'
}

export interface Message {
  id: number
  mailbox_id: number
  uid: string
  subject: string
  from: string
  to: string[]
  cc?: string[]
  bcc?: string[]
  date_sent: string
  date_received: string
  body_text?: string
  body_html?: string
  size: number
  flags: string[]
  attachments?: Attachment[]
  read: boolean
  starred: boolean
  answered: boolean
  forwarded: boolean
  task_completed: boolean
}

export interface Attachment {
  id: string
  filename: string
  content_type: string
  size: number
}

export interface ComposeData {
  to: string[]
  cc?: string[]
  bcc?: string[]
  subject: string
  body_text?: string
  body_html?: string
  attachments?: string[]
  in_reply_to?: string
  references?: string
  draft_id?: number
}

export interface Contact {
  id: number
  name: string
  email: string
  tel?: string
}

export interface Calendar {
  id: number
  name: string
  display_name: string
  color?: string
  description?: string
  timezone: string
}

export interface CalendarEvent {
  id: number
  uid: string
  summary: string
  description?: string
  location?: string
  start_time: string
  end_time: string
  all_day: boolean
  timezone?: string
}

export const useWebmailStore = defineStore('webmail', {
  state: () => ({
    // Mailboxes and messages
    mailboxes: [] as Mailbox[],
    currentMailbox: null as Mailbox | null,
    messages: [] as Message[],
    currentMessage: null as Message | null,
    
    // Message loading and pagination
    loading: false,
    loadingMessages: false,
    currentPage: 1,
    totalPages: 1,
    pageSize: 50,
    
    // Composition
    composing: false,
    composeData: null as ComposeData | null,
    sendingMessage: false,
    
    // Search
    searchQuery: '',
    searchResults: [] as Message[],
    searching: false,
    
    // Contacts and Calendar
    contacts: [] as Contact[],
    calendars: [] as Calendar[],
    upcomingEvents: [] as CalendarEvent[],
    
    // UI State
    selectedMessages: [] as number[],
    sidebarCollapsed: false,
    previewPane: true,
    
    // Drafts
    drafts: [] as Message[],
    savingDraft: false,
    
    // Errors
    error: null as string | null,
  }),

  getters: {
    // Get inbox mailbox
    inbox: (state) => state.mailboxes.find(m => m.type === 'inbox'),
    
    // Get sent mailbox
    sent: (state) => state.mailboxes.find(m => m.type === 'sent'),
    
    // Get drafts mailbox
    drafts: (state) => state.mailboxes.find(m => m.type === 'drafts'),
    
    // Get trash mailbox
    trash: (state) => state.mailboxes.find(m => m.type === 'trash'),
    
    // Get unread count across all mailboxes
    totalUnread: (state) => state.mailboxes.reduce((total, mailbox) => total + mailbox.unread_count, 0),
    
    // Get filtered messages based on current mailbox
    filteredMessages: (state) => {
      if (!state.currentMailbox) return state.messages
      return state.messages.filter(msg => msg.mailbox_id === state.currentMailbox.id)
    },
    
    // Get unread messages in current mailbox
    unreadMessages: (state, getters) => getters.filteredMessages.filter((msg: Message) => !msg.read),
    
    // Get starred messages in current mailbox
    starredMessages: (state, getters) => getters.filteredMessages.filter((msg: Message) => msg.starred),
    
    // Check if message is selected
    isMessageSelected: (state) => (messageId: number) => state.selectedMessages.includes(messageId),
    
    // Get compose recipients for display
    composeRecipients: (state) => {
      if (!state.composeData) return []
      return [
        ...state.composeData.to || [],
        ...state.composeData.cc || [],
        ...state.composeData.bcc || []
      ]
    },
  },

  actions: {
    // Initialize webmail
    async initialize() {
      try {
        await this.loadMailboxes()
        await this.loadDrafts()
        await this.loadContacts()
        await this.loadCalendars()
        await this.loadUpcomingEvents()
      } catch (error) {
        this.setError('Failed to initialize webmail')
        throw error
      }
    },

    // Mailbox operations
    async loadMailboxes() {
      this.loading = true
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/mailboxes`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to load mailboxes')
        
        const data = await response.json()
        this.mailboxes = data.mailboxes || []
        
        // Set default mailbox to inbox
        if (this.mailboxes.length > 0 && !this.currentMailbox) {
          this.currentMailbox = this.inbox || this.mailboxes[0]
        }
      } catch (error) {
        this.setError('Failed to load mailboxes')
        throw error
      } finally {
        this.loading = false
      }
    },

    // Message operations
    async loadMessages(mailboxId?: number, page = 1) {
      this.loadingMessages = true
      this.currentPage = page
      
      try {
        const API_BASE = useApiBase()
        const targetMailboxId = mailboxId || this.currentMailbox?.id
        
        if (!targetMailboxId) {
          throw new Error('No mailbox selected')
        }
        
        const url = `${API_BASE}/webmail/mailboxes/${targetMailboxId}/messages?page=${page}&limit=${this.pageSize}`
        const response = await fetch(url, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to load messages')
        
        const data = await response.json()
        
        if (page === 1) {
          this.messages = data.messages || []
        } else {
          this.messages.push(...(data.messages || []))
        }
        
        this.totalPages = Math.ceil((data.total || this.messages.length) / this.pageSize)
      } catch (error) {
        this.setError('Failed to load messages')
        throw error
      } finally {
        this.loadingMessages = false
      }
    },

    async loadMessage(messageId: number) {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/messages/${messageId}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to load message')
        
        const message = await response.json()
        this.currentMessage = message
        
        // Mark as read if not already
        if (!message.read) {
          await this.updateMessageFlags(messageId, ['\\Seen'], 'add')
          message.read = true
          message.flags = message.flags || []
          message.flags.push('\\Seen')
        }
        
        return message
      } catch (error) {
        this.setError('Failed to load message')
        throw error
      }
    },

    async sendMessage(composeData: ComposeData) {
      this.sendingMessage = true
      
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/messages`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            to: composeData.to.join(', '),
            cc: composeData.cc?.join(', ') || '',
            bcc: composeData.bcc?.join(', ') || '',
            subject: composeData.subject,
            body_text: composeData.body_text,
            body_html: composeData.body_html,
            attachments: composeData.attachments || []
          })
        })
        
        if (!response.ok) throw new Error('Failed to send message')
        
        const result = await response.json()
        
        // Clear compose form
        this.composing = false
        this.composeData = null
        
        // Remove draft if it exists
        if (composeData.draft_id) {
          await this.deleteDraft(composeData.draft_id)
        }
        
        // Reload sent mailbox
        if (this.sent) {
          await this.loadMessages(this.sent.id)
        }
        
        return result
      } catch (error) {
        this.setError('Failed to send message')
        throw error
      } finally {
        this.sendingMessage = false
      }
    },

    async deleteMessage(messageId: number) {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/messages/${messageId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to delete message')
        
        // Remove from local state
        this.messages = this.messages.filter(msg => msg.id !== messageId)
        
        // Clear current message if it's the one being deleted
        if (this.currentMessage?.id === messageId) {
          this.currentMessage = null
        }
        
        // Remove from selection
        this.selectedMessages = this.selectedMessages.filter(id => id !== messageId)
      } catch (error) {
        this.setError('Failed to delete message')
        throw error
      }
    },

    async moveMessage(messageId: number, targetMailboxId: number) {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/messages/${messageId}/move`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ mailbox_id: targetMailboxId })
        })
        
        if (!response.ok) throw new Error('Failed to move message')
        
        // Update local state
        const message = this.messages.find(msg => msg.id === messageId)
        if (message) {
          message.mailbox_id = targetMailboxId
        }
        
        // Remove from current mailbox view if we're not moving to current mailbox
        if (this.currentMailbox?.id !== targetMailboxId) {
          this.messages = this.messages.filter(msg => msg.id !== messageId)
        }
      } catch (error) {
        this.setError('Failed to move message')
        throw error
      }
    },

    async updateMessageFlags(messageId: number, flags: string[], action: 'add' | 'remove') {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/messages/${messageId}/flags`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ flags, action })
        })
        
        if (!response.ok) throw new Error('Failed to update message flags')
        
        // Update local state
        const message = this.messages.find(msg => msg.id === messageId)
        if (message) {
          if (action === 'add') {
            message.flags.push(...flags)
          } else {
            message.flags = message.flags.filter(flag => !flags.includes(flag))
          }
          
          // Update derived properties
          message.read = message.flags.includes('\\Seen')
          message.starred = message.flags.includes('\\Flagged')
          message.answered = message.flags.includes('\\Answered')
          message.forwarded = message.flags.includes('$Forwarded')
          message.task_completed = message.task_completed || false
        }
      } catch (error) {
        this.setError('Failed to update message flags')
        throw error
      }
    },

    // Search operations
    async searchMessages(query: string) {
      if (!query.trim()) {
        this.searchResults = []
        return
      }
      
      this.searching = true
      this.searchQuery = query
      
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/search?q=${encodeURIComponent(query)}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to search messages')
        
        const data = await response.json()
        this.searchResults = data.messages || []
      } catch (error) {
        this.setError('Failed to search messages')
        throw error
      } finally {
        this.searching = false
      }
    },

    // Draft operations
    async saveDraft(composeData: ComposeData) {
      this.savingDraft = true
      
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/drafts`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            draft_id: composeData.draft_id,
            to: composeData.to,
            cc: composeData.cc,
            bcc: composeData.bcc,
            subject: composeData.subject,
            body_html: composeData.body_html,
            body_text: composeData.body_text,
            in_reply_to: composeData.in_reply_to,
            references: composeData.references,
            attachments: composeData.attachments
          })
        })
        
        if (!response.ok) throw new Error('Failed to save draft')
        
        const draft = await response.json()
        
        // Update draft ID if this is a new draft
        if (!composeData.draft_id) {
          composeData.draft_id = draft.id
        }
        
        await this.loadDrafts()
        return draft
      } catch (error) {
        this.setError('Failed to save draft')
        throw error
      } finally {
        this.savingDraft = false
      }
    },

    async loadDrafts() {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/drafts`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to load drafts')
        
        const data = await response.json()
        this.drafts = data.drafts || []
      } catch (error) {
        this.setError('Failed to load drafts')
        throw error
      }
    },

    async deleteDraft(draftId: number) {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/drafts/${draftId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to delete draft')
        
        // Remove from local state
        this.drafts = this.drafts.filter(draft => draft.id !== draftId)
      } catch (error) {
        this.setError('Failed to delete draft')
        throw error
      }
    },

    // Contact operations
    async searchContacts(query: string) {
      if (!query.trim()) {
        this.contacts = []
        return
      }
      
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/contacts/search?q=${encodeURIComponent(query)}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to search contacts')
        
        this.contacts = await response.json()
      } catch (error) {
        this.setError('Failed to search contacts')
        throw error
      }
    },

    async loadContacts() {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/contacts/addressbooks`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to load addressbooks')
        
        const addressbooks = await response.json()
        
        // Load contacts from all addressbooks
        this.contacts = []
        for (const addressbook of addressbooks) {
          try {
            const contactsResponse = await fetch(`${API_BASE}/webmail/contacts/addressbooks/${addressbook.id}/contacts`, {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              }
            })
            
            if (contactsResponse.ok) {
              const contacts = await contactsResponse.json()
              this.contacts.push(...contacts)
            }
          } catch (error) {
            console.warn(`Failed to load contacts from addressbook ${addressbook.id}:`, error)
          }
        }
      } catch (error) {
        this.setError('Failed to load contacts')
        throw error
      }
    },

    // Calendar operations
    async loadCalendars() {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/calendar/calendars`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to load calendars')
        
        this.calendars = await response.json()
      } catch (error) {
        this.setError('Failed to load calendars')
        throw error
      }
    },

    async loadUpcomingEvents(days = 7) {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/calendar/upcoming?days=${days}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (!response.ok) throw new Error('Failed to load upcoming events')
        
        this.upcomingEvents = await response.json()
      } catch (error) {
        this.setError('Failed to load upcoming events')
        throw error
      }
    },

    async createCalendarEvent(eventData: Partial<CalendarEvent>) {
      try {
        const API_BASE = useApiBase()
        const response = await fetch(`${API_BASE}/webmail/calendar/events`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(eventData)
        })
        
        if (!response.ok) throw new Error('Failed to create calendar event')
        
        const event = await response.json()
        this.upcomingEvents.push(event)
        
        return event
      } catch (error) {
        this.setError('Failed to create calendar event')
        throw error
      }
    },

    // UI State management
    setCurrentMailbox(mailbox: Mailbox) {
      this.currentMailbox = mailbox
      this.messages = []
      this.currentPage = 1
      this.selectedMessages = []
      this.loadMessages(mailbox.id)
    },

    setCurrentMessage(message: Message | null) {
      this.currentMessage = message
    },

    toggleMessageSelection(messageId: number) {
      const index = this.selectedMessages.indexOf(messageId)
      if (index > -1) {
        this.selectedMessages.splice(index, 1)
      } else {
        this.selectedMessages.push(messageId)
      }
    },

    selectAllMessages() {
      this.selectedMessages = this.filteredMessages.map((msg: Message) => msg.id)
    },

    clearMessageSelection() {
      this.selectedMessages = []
    },

    startCompose(composeData?: Partial<ComposeData>) {
      this.composing = true
      this.composeData = {
        to: [],
        cc: [],
        bcc: [],
        subject: '',
        body_text: '',
        body_html: '',
        attachments: [],
        ...composeData
      } as ComposeData
    },

    stopCompose() {
      this.composing = false
      this.composeData = null
    },

    updateComposeData(data: Partial<ComposeData>) {
      if (this.composeData) {
        Object.assign(this.composeData, data)
      }
    },

    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },

    togglePreviewPane() {
      this.previewPane = !this.previewPane
      if (!this.previewPane) {
        this.currentMessage = null
      }
    },

    setError(error: string | null) {
      this.error = error
    },

    clearError() {
      this.error = null
    }
  }
})