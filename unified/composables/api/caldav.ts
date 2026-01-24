import { useApiBase } from '../useApiBase'

export const getAuthToken = () => {
  return typeof window !== 'undefined' ? localStorage.getItem('token') : null
}

export const getAuthHeaders = () => {
  const token = getAuthToken()
  return {
    'Content-Type': 'application/json',
    ...(token ? { 'Authorization': `Bearer ${token}` } : {})
  }
}

export interface Principal {
  href: string
  displayName?: string
  user?: string
}

export interface AddressBook {
  id: number
  name: string
  description?: string
  owner_id: number
  user_id?: number
  url: string
  principal: string
  etag?: string
  ctag?: string
  supported_components: string[]
  created_at: string
  updated_at: string
  sync_token?: string
  sync_status: string
}

export interface Contact {
  id: number
  addressbook_id: number
  fn: string
  ln: string
  n: string
  email: string
  organization?: string
  phone?: string
  title?: string
  company?: string
  department?: string
  job_title?: string
  office_location?: string
  home_address?: Address
  mobile?: string
  pager?: string
  notes?: string
  categories?: string[]
  urls?: string[]
  birthday?: string
  anniversary?: string
  photo?: string
  created_at: string
  updated_at: string
}

export interface Address {
  street?: string
  city?: string
  state?: string
  postal_code?: string
  country?: string
}

export interface Calendar {
  id: number
  name: string
  description?: string
  color?: string
  timezone: string
  user_id: number
  url: string
  ctag?: string
  calendar_order: number
  created_at: string
  updated_at: string
  sync_token?: string
  sync_status: string
  events: CalendarEvent[]
  components?: CalendarComponent[]
  supported_components: string[]
  free_busy: string[]
  allowed_attendees?: string[]
}

export interface CalendarEvent {
  id: number
  event_type: 'created' | 'modified' | 'deleted'
  calendar_id: number
  event_data?: any
  recurrence_id?: number
  alert?: string
  notes?: string
  uid?: string
  start: string
  end: string
  all_day?: boolean
  created_at: string
  updated_at: string
}

export interface CalendarComponent {
  id: number
  name: string
  type: 'time-zone' | 'supported-components' | 'calendar-color'
  value?: string
  description?: string
  parameter?: string
  required?: boolean
  default?: string
}

export interface Task {
  id: number
  title: string
  status: 'needs-action' | 'in-progress' | 'completed' | 'failed'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  due?: string
  created_at: string
  updated_at: string
  completed_at?: string
  assigned_to?: string
  created_by?: string
  notes?: string
}

export const CalDavService = '/caldav'

export const useCalDavApi = () => {
  const API_BASE = useApiBase()

  const getPrincipal = async (path?: string): Promise<{ principal: Principal; displayName: string }> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : '/principal'}`, {
      method: 'PROPFIND',
      headers: headers
    })

    if (!response.ok) {
      throw new Error('Failed to get principal')
    }

    const principalData = await response.json()
    
    return {
      principal: principalData.principal,
      displayName: principalData.displayName
    }
  }

  const getHomeSet = async (path?: string): Promise<{ href: string; displayname: string }> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : ''}`, {
      method: 'PROPFIND',
      headers: headers
    })

    if (!response.ok) {
      throw new Error('Failed to get home set')
    }

    const homesetData = await response.json()
    
    return {
      href: homesetData.href,
      displayname: homesetData.displayname
    }
  }

  const getUser = async (path: string, user?: string): Promise<{ user: string }> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    // First try to get user by email from LDAP or database
    let response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : '/user'}`, {
      method: 'REPORT',
      headers: headers
    })

    if (!response.ok) {
      // Try getting user by database
      const dbResponse = await fetch(`${API_BASE}/caldav/users?user=${user}`)
      
      if (dbResponse.ok) {
        const userData = await dbResponse.json()
        return {
          user: userData.name || user
        }
      }
      
      // Fallback to email authentication
      response = await fetch(`${API_BASE}/caldav/email?user=${user}`, {
        method: 'REPORT',
        headers: headers
      })
    }

    if (!response.ok) {
      throw new Error('Failed to get user')
    }

    const emailData = await response.json()
    return {
      user: emailData.email || user
    }
  }

  const getAddressBooks = async (path: string): Promise<AddressBook[]> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : '/addressbooks'}`, {
      method: 'PROPFIND',
      headers: headers
    })

    if (!response.ok) {
      throw new Error('Failed to get address books')
    }

    const booksData = await response.json()
    
    return booksData.data || []
  }

  const getAddressBook = async (path: string, bookId: string): Promise<AddressBook> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}/addressbooks/${bookId}` : `/addressbooks/${bookId}`}`, {
      method: 'PROPFIND',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to get address book ${bookId}`)
    }

    const bookData = await response.json()
    
    return bookData.data
  }

  const createAddressBook = async (path: string, data: Partial<AddressBook>): Promise<AddressBook> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}/addressbooks` : '/addressbooks'}`, {
      method: 'MKCOLLECTION',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error('Failed to create address book')
    }

    const bookData = await response.json()
    
    return bookData.data
  }

  const updateAddressBook = async (path: string, bookId: string, data: Partial<AddressBook>): Promise<AddressBook> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}/addressbooks/${bookId}` : `/addressbooks/${bookId}`}`, {
      method: 'PROPPATCH',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error(`Failed to update address book ${bookId}`)
    }

    const bookData = await response.json()
    
    return bookData.data
  }

  const deleteAddressBook = async (path: string, bookId: string): Promise<void> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}/addressbooks/${bookId}` : `/addressbooks/${bookId}`}`, {
      method: 'DELETE',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to delete address book ${bookId}`)
    }

    return
  }

  const syncAddressBook = async (path: string, bookId: string): Promise<void> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav/${path ? `/${path}/addressbooks/${bookId}` : `/addressbooks/${bookId}/sync`}`, {
      method: 'REPORT',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to sync address book ${bookId}`)
    }

    const result = await response.json()
    
    return result.data
  }

  const reportSync = async (path: string, bookId: string, report: string): Promise<void> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}/addressbooks/${bookId}` : `/addressbooks/${bookId}/report`}`, {
      method: 'REPORT',
      headers: headers,
      body: JSON.stringify({ report })
    })

    if (!response.ok) {
      throw new Error(`Failed to report sync for address book ${bookId}`)
    }

    return
  }

  const getCalendar = async (path: string): Promise<Calendar> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : '/calendar'}`, {
      method: 'REPORT',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to get calendar`)
    }

    const calendarData = await response.json()
    
    return calendarData.data
  }

  const createCalendar = async (path: string, data: Partial<Calendar>): Promise<Calendar> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : '/calendar'}`, {
      method: 'MKCALENDAR',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error(`Failed to create calendar`)
    }

    const calendarData = await response.json()
    
    return calendarData.data
  }

  const updateCalendar = async (path: string, calendarId: string, data: Partial<Calendar>): Promise<Calendar> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : `/calendar/${calendarId}`}`, {
      method: 'PROPPATCH',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error(`Failed to update calendar ${calendarId}`)
    }

    const calendarData = await response.json()
    
    return calendarData.data
  }

  const deleteCalendar = async (path: string, calendarId: string): Promise<void> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : `/calendar/${calendarId}`}`, {
      method: 'DELETE',
      headers: headers
    })

    if (!response.ok) {
      return
    }

    return
  }

  const syncCalendar = async (path: string, calendarId: string): Promise<void> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : `/calendar/${calendarId}/sync`}`, {
      method: 'REPORT',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to sync calendar ${calendarId}`)
    }

    const result = await response.json()
    
    return result.data
  }

  const getCalendarEvents = async (path: string, calendarId: string, start?: string, end?: string): Promise<CalendarEvent[]> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const params = new URLSearchParams()
    if (start) params.append('start', start)
    if (end) params.append('end', end)
    
    const response = await fetch(`${API_BASE}/caldav${path ? `/${path}` : `/calendar/${calendarId}/events`}?${params.toString()}`, {
      method: 'REPORT',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to get calendar events`)
    }

    const calendarEventsData = await response.json()
    
    return calendarEventsData.data || []
  }

  return {
    getPrincipal,
    getHomeSet,
    getUser,
    getAddressBooks,
    getAddressBook,
    createAddressBook,
    updateAddressBook,
    deleteAddressBook,
    syncAddressBook,
    reportSync,
    getCalendar,
    createCalendar,
    updateCalendar,
    deleteCalendar,
    syncCalendar,
    getCalendarEvents
  }
}

export const CardDavService = '/carddav'

export const useCardDavApi = () => {
  const API_BASE = useApiBase()

  const getPrincipal = async (path?: string): Promise<{ principal: Principal; displayName: string }> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : '/principal'}`, {
      method: 'PROPFIND',
      headers: headers
    })

    if (!response.ok) {
      throw new Error('Failed to get principal')
    }

    const principalData = await response.json()
    
    return {
      principal: principalData.principal,
      displayName: principalData.displayName
    }
  }

  const getAddressBooks = async (path?: string): Promise<AddressBook[]> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : '/addressbooks'}`, {
      method: 'PROPFIND',
      headers: headers
    })

    if (!response.ok) {
      throw new Error('Failed to get address books')
    }

    const booksData = await response.json()
    
    return booksData.data || []
  }

  const getAddressBook = async (path: string, bookId: string): Promise<AddressBook> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : `/addressbooks/${bookId}`}`, {
      method: 'PROPFIND',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to get address book ${bookId}`)
    }

    const bookData = await response.json()
    
    return bookData.data
  }

  const createAddressBook = async (path: string, data: Partial<AddressBook>): Promise<AddressBook> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : '/addressbooks'}`, {
      method: 'MKCOLLECTION',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error('Failed to create address book')
    }

    const bookData = await response.json()
    
    return bookData.data
  }

  const updateAddressBook = async (path: string, bookId: string, data: Partial<AddressBook>): Promise<AddressBook> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : `/addressbooks/${bookId}`}`, {
      method: 'PROPPATCH',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error(`Failed to update address book ${bookId}`)
    }

    const bookData = await response.json()
    
    return bookData.data
  }

  const deleteAddressBook = async (path: string, bookId: string): Promise<void> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : `/addressbooks/${bookId}`}`, {
      method: 'DELETE',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to delete address book ${bookId}`)
    }

    return
  }

  const getContact = async (path: string, contactId: string): Promise<Contact> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : `/contacts/${contactId}`}`, {
      method: 'GET',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to get contact ${contactId}`)
    }

    const contactData = await response.json()
    
    return contactData.data
  }

  const createContact = async (path: string, data: Partial<Contact>): Promise<Contact> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : '/contacts'}`, {
      method: 'PUT',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error('Failed to create contact')
    }

    const contactData = await response.json()
    
    return contactData.data
  }

  const updateContact = async (path: string, contactId: string, data: Partial<Contact>): Promise<Contact> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : `/contacts/${contactId}`}`, {
      method: 'PROPPATCH',
      headers: headers,
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error(`Failed to update contact ${contactId}`)
    }

    const contactData = await response.json()
    
    return contactData.data
  }

  const deleteContact = async (path: string, contactId: string): Promise<void> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : `/contacts/${contactId}`}`, {
      method: 'DELETE',
      headers: headers
    })

    if (!response.ok) {
      throw new Error(`Failed to delete contact ${contactId}`)
    }

    return
  }

  const getSyncToken = async (path: string): Promise<string> => {
    const token = getAuthToken()
    const headers = getAuthHeaders()
    
    // First try to get existing sync token
    const response = await fetch(`${API_BASE}/carddav${path ? `/${path}` : '/sync-token'}`, {
      method: 'GET',
      headers: headers
    })

    if (response.ok && response.status === 200) {
      const syncData = await response.json()
      return syncData.data.token || ''
    }
    
    // Fallback - create new sync token if none exists
    const syncToken = `sync-token-${Date.now()}`
    
    const createResponse = await fetch(`${API_BASE}/carddav/${path ? `/${path}` : '/sync-token'}`, {
      method: 'POST',
      headers: headers,
      body: JSON.stringify({ token: syncToken }),
    })

    if (!createResponse.ok) {
      throw new Error('Failed to create sync token')
    }

    return syncToken || ''
  }

  return {
    getPrincipal,
    getAddressBooks,
    getAddressBook,
    createAddressBook,
    updateAddressBook,
    deleteAddressBook,
    getSyncToken,
    getContact,
    createContact,
    updateContact,
    deleteContact
  }
}