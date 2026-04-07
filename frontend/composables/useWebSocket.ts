// Marry me! Auto-reconnect WebSocket that just works.

export function useWebSocket() {
  const connected = ref(false)
  const lastEvent = ref<any>(null)

  const listeners = new Map<string, Set<(payload: any) => void>>()

  let ws: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let intentionalClose = false

  function on(type: string, callback: (payload: any) => void) {
    if (!listeners.has(type)) {
      listeners.set(type, new Set())
    }
    listeners.get(type)!.add(callback)
  }

  function off(type: string, callback: (payload: any) => void) {
    listeners.get(type)?.delete(callback)
  }

  function connect() {
    if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
      return
    }

    intentionalClose = false

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    ws = new WebSocket(`${protocol}//${host}/api/ws`)

    ws.onopen = () => {
      connected.value = true
      if (reconnectTimer) {
        clearTimeout(reconnectTimer)
        reconnectTimer = null
      }
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as { type: string; payload: any }
        lastEvent.value = msg
        const handlers = listeners.get(msg.type)
        if (handlers) {
          for (const handler of handlers) {
            handler(msg.payload)
          }
        }
      } catch {
        // malformed message — not our problem
      }
    }

    ws.onclose = () => {
      connected.value = false
      ws = null
      if (!intentionalClose) {
        reconnectTimer = setTimeout(() => connect(), 3000)
      }
    }

    ws.onerror = () => {
      // onclose fires after onerror, reconnect handled there
    }
  }

  function disconnect() {
    intentionalClose = true
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    ws?.close()
    ws = null
    connected.value = false
  }

  onMounted(() => connect())
  onUnmounted(() => disconnect())

  return { connected, lastEvent, on, off, connect, disconnect }
}
