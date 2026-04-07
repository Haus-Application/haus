// AI concierge. She's basically playing a movie executive and it's working out great.

export interface ToolCall {
  tool: string
  input: string
  result: string
}

export interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
  toolCalls?: ToolCall[]
}

export function useChat() {
  const messages = ref<ChatMessage[]>([])
  const loading = ref(false)
  const history = ref<any[]>([])  // anthropic message format for context

  async function sendMessage(text: string) {
    messages.value.push({ role: 'user', content: text })
    loading.value = true

    try {
      const res = await fetch('/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: text, history: history.value }),
      })

      if (!res.ok) {
        messages.value.push({
          role: 'assistant',
          content: 'Sorry, something went wrong. Try again.',
        })
        return
      }

      const data = await res.json()

      messages.value.push({
        role: 'assistant',
        content: data.text ?? '',
        toolCalls: data.tool_calls ?? [],
      })

      history.value = data.messages ?? history.value
    } catch {
      messages.value.push({
        role: 'assistant',
        content: 'Could not reach the server. Check your connection.',
      })
    } finally {
      loading.value = false
    }
  }

  return { messages, loading, sendMessage }
}
