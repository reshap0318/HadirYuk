import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post, type IApiResponse } from '@/plugins/axios'

export interface IAiChatMessage {
  role: 'user' | 'assistant'
  content: string
  created_at: string
}

interface IAiChatMessageResponse {
  answer: string
}

interface IAiChatHistoryResponse {
  messages: IAiChatMessage[]
}

export const useAiChatStore = defineStore('ai-chat', () => {
  const messages = ref<IAiChatMessage[]>([])
  const loading = ref<Record<string, boolean>>({
    Send: false,
    History: false,
    Reset: false,
  })

  async function fetchHistory() {
    loading.value.History = true
    try {
      const { data } = await get<IApiResponse<IAiChatHistoryResponse>>('/ai-chat/history')
      messages.value = data.data?.messages || []
    } catch (error: any) {
      console.error('Failed to fetch AI chat history', error)
    } finally {
      loading.value.History = false
    }
  }

  async function sendMessage(message: string): Promise<boolean> {
    loading.value.Send = true
    messages.value.push({ role: 'user', content: message, created_at: new Date().toISOString() })
    try {
      const { data } = await post<IApiResponse<IAiChatMessageResponse>>('/ai-chat/message', {
        message,
      })
      messages.value.push({
        role: 'assistant',
        content: data.data.answer,
        created_at: new Date().toISOString(),
      })
      return true
    } catch {
      // Roll back the optimistic user message so a failed send doesn't look like it went through.
      messages.value.pop()
      return false
    } finally {
      loading.value.Send = false
    }
  }

  async function resetChat() {
    loading.value.Reset = true
    try {
      await post('/ai-chat/reset')
      messages.value = []
    } catch {
      // interceptor already toasts the error
    } finally {
      loading.value.Reset = false
    }
  }

  return {
    messages,
    loading,
    fetchHistory,
    sendMessage,
    resetChat,
  }
})
