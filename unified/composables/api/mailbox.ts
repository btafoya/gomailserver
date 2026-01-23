export const useMailboxApi = () => {
  const { $fetch } = useNuxtApp()
  const config = useRuntimeConfig()

  const baseApiUrl = config.public.apiBaseUrl || '/api/v1'

  const getMailboxes = async (domainId?: string) => {
    const url = domainId 
      ? `${baseApiUrl}/domains/${domainId}/mailboxes`
      : `${baseApiUrl}/mailboxes`
    
    const response = await $fetch(url, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    return response
  }

  const getMailbox = async (id: string) => {
    const response = await $fetch(`${baseApiUrl}/mailboxes/${id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    return response
  }

  const createMailbox = async (data: any) => {
    const response = await $fetch(`${baseApiUrl}/mailboxes`, {
      method: 'POST',
      body: data,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })
    return response
  }

  const updateMailbox = async (id: string, data: any) => {
    const response = await $fetch(`${baseApiUrl}/mailboxes/${id}`, {
      method: 'PUT',
      body: data,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })
    return response
  }

  const deleteMailbox = async (id: string) => {
    const response = await $fetch(`${baseApiUrl}/mailboxes/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    return response
  }

  const getMailboxStats = async (id: string) => {
    const response = await $fetch(`${baseApiUrl}/mailboxes/${id}/stats`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    return response
  }

  const updateMailboxQuota = async (id: string, quota: number) => {
    const response = await $fetch(`${baseApiUrl}/mailboxes/${id}/quota`, {
      method: 'PUT',
      body: { quota },
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })
    return response
  }

  return {
    getMailboxes,
    getMailbox,
    createMailbox,
    updateMailbox,
    deleteMailbox,
    getMailboxStats,
    updateMailboxQuota
  }
}