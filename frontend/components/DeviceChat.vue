<template>
  <Teleport to="body">
    <transition name="drawer">
      <div v-if="true" class="device-chat-backdrop" @click="$emit('close')">
        <div class="device-chat-panel" @click.stop>
          <div class="chat-header">
            <BrandLogo :brand="brandKey" :size="28" />
            <div class="chat-header-info">
              <span class="chat-device-name">{{ device.name }}</span>
              <span class="chat-device-meta">{{ device.ip }} · {{ device.device_type }}</span>
            </div>
            <button class="chat-close" @click="$emit('close')">✕</button>
          </div>

          <div class="chat-messages" ref="messagesRef">
            <div v-if="messages.length === 0" class="chat-suggestions">
              <button v-for="s in suggestions" :key="s" class="suggestion-chip" @click="send(s)">{{ s }}</button>
            </div>

            <div v-for="(msg, i) in messages" :key="i" class="chat-msg" :class="msg.role">
              <div class="msg-bubble">{{ msg.content }}</div>
              <div v-if="msg.toolCalls?.length" class="msg-tools">
                <span v-for="(tc, j) in msg.toolCalls" :key="j" class="tool-badge">{{ tc.tool }}</span>
              </div>
            </div>

            <div v-if="loading" class="chat-msg assistant">
              <div class="msg-bubble typing">
                <span class="dot" /><span class="dot" /><span class="dot" />
              </div>
            </div>
          </div>

          <div class="chat-input-row">
            <input
              v-model="inputText"
              @keydown.enter="send(inputText)"
              placeholder="Ask this device..."
              class="chat-input"
              :disabled="loading"
            />
            <button @click="send(inputText)" class="chat-send" :disabled="!inputText.trim() || loading">→</button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
interface DeviceContext {
  ip: string
  name: string
  manufacturer: string
  model: string
  device_type: string
  category: string
  protocols: string[]
}

interface ChatMsg {
  role: 'user' | 'assistant'
  content: string
  toolCalls?: { tool: string; input: string; result: string }[]
}

const props = defineProps<{ device: DeviceContext }>()
defineEmits<{ close: [] }>()

const messages = ref<ChatMsg[]>([])
const loading = ref(false)
const inputText = ref('')
const history = ref<any[]>([])
const messagesRef = ref<HTMLElement | null>(null)

const brandKey = computed(() => {
  const mfr = props.device.manufacturer?.toLowerCase().trim() || ''
  const dt = props.device.device_type?.toLowerCase().trim() || ''
  return (mfr && mfr !== 'unknown') ? mfr : dt || ''
})

const suggestions = computed(() => {
  const base = ['Are you on?', 'Turn on', 'What can you do?']
  if (props.device.device_type === 'dimmer') base.push('Set to 50%')
  if (props.device.device_type === 'hue_bridge') {
    base.length = 0
    base.push('List lights', 'What scenes are available?', 'Turn on living room')
  }
  return base
})

async function send(text: string) {
  if (!text.trim() || loading.value) return
  const msg = text.trim()
  inputText.value = ''
  messages.value.push({ role: 'user', content: msg })
  loading.value = true

  try {
    const res = await fetch('/api/chat/device', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        device: props.device,
        message: msg,
        history: history.value,
      }),
    })
    const data = await res.json()
    if (data.error) {
      messages.value.push({ role: 'assistant', content: data.error })
    } else {
      messages.value.push({
        role: 'assistant',
        content: data.text,
        toolCalls: data.tool_calls,
      })
      history.value = data.messages || []
    }
  } catch {
    messages.value.push({ role: 'assistant', content: 'Failed to reach the server.' })
  }

  loading.value = false
}

watch(messages, () => {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}, { deep: true })
</script>

<style scoped>
.device-chat-backdrop {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: flex-end;
}

.device-chat-panel {
  width: 420px;
  max-width: 100%;
  height: 100vh;
  background: var(--color-bg);
  border-left: 1px solid var(--color-surface-border);
  display: flex;
  flex-direction: column;
  animation: slide-in 0.3s ease-out;
}

@keyframes slide-in {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-surface-border);
  flex-shrink: 0;
}

.chat-header-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.chat-device-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-device-meta {
  font-size: 12px;
  color: var(--color-text-tertiary);
  font-family: var(--font-mono);
}

.chat-close {
  background: none;
  border: none;
  color: var(--color-text-secondary);
  font-size: 18px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  transition: background 0.2s;
}

.chat-close:hover {
  background: var(--color-surface);
  color: var(--color-text);
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
  padding: 40px 0;
}

.suggestion-chip {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  color: var(--color-text-secondary);
  font-size: 13px;
  font-family: var(--font-sans);
  padding: 6px 14px;
  border-radius: 99px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s, border-color 0.2s;
}

.suggestion-chip:hover {
  background: var(--color-surface-hover);
  border-color: var(--color-accent);
  color: var(--color-text);
}

.chat-msg {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.chat-msg.user {
  align-items: flex-end;
}

.chat-msg.assistant {
  align-items: flex-start;
}

.msg-bubble {
  max-width: 85%;
  padding: 10px 14px;
  font-size: 14px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}

.chat-msg.user .msg-bubble {
  background: var(--color-accent);
  color: #fff;
  border-radius: 16px 16px 4px 16px;
}

.chat-msg.assistant .msg-bubble {
  background: var(--color-surface);
  color: var(--color-text);
  border-radius: 16px 16px 16px 4px;
}

.msg-tools {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.tool-badge {
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--color-text-tertiary);
  background: rgba(255, 255, 255, 0.05);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
}

.msg-bubble.typing {
  display: flex;
  gap: 4px;
  padding: 12px 18px;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-text-tertiary);
  animation: pulse-dot 1.2s ease-in-out infinite;
}

.dot:nth-child(2) { animation-delay: 0.2s; }
.dot:nth-child(3) { animation-delay: 0.4s; }

@keyframes pulse-dot {
  0%, 80%, 100% { opacity: 0.3; transform: scale(0.8); }
  40% { opacity: 1; transform: scale(1); }
}

.chat-input-row {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--color-surface-border);
  flex-shrink: 0;
}

.chat-input {
  flex: 1;
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-md);
  padding: 10px 14px;
  font-size: 14px;
  font-family: var(--font-sans);
  color: var(--color-text);
  outline: none;
  transition: border-color 0.2s;
}

.chat-input:focus {
  border-color: var(--color-accent);
}

.chat-input::placeholder {
  color: var(--color-text-tertiary);
}

.chat-send {
  background: var(--color-accent);
  border: none;
  color: #fff;
  font-size: 18px;
  width: 42px;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background 0.2s;
}

.chat-send:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.chat-send:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

@media (max-width: 480px) {
  .device-chat-panel {
    width: 100%;
  }
}
</style>
