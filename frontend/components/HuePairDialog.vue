<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="show" class="modal-backdrop" @click.self="emit('close')">
        <div class="modal" role="dialog" aria-modal="true" aria-labelledby="pair-title">
          <div class="modal-header">
            <h2 id="pair-title" class="modal-title">Pair Hue Bridge</h2>
            <button class="close-btn" aria-label="Close" @click="emit('close')">
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none">
                <path d="M2 2l14 14M16 2L2 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
              </svg>
            </button>
          </div>

          <!-- Step 1: Discover -->
          <div v-if="step === 'discover'" class="modal-body">
            <p class="step-description">
              Find Hue bridges on your network. Make sure your bridge is powered on and connected.
            </p>
            <button class="btn-primary" :disabled="discovering" @click="findBridges">
              <span v-if="discovering" class="spinner" />
              {{ discovering ? 'Scanning...' : 'Find Bridges' }}
            </button>
            <div v-if="discoverError" class="error-msg">{{ discoverError }}</div>
          </div>

          <!-- Step 2: Select bridge -->
          <div v-else-if="step === 'select'" class="modal-body">
            <p class="step-description">Found {{ bridges.length }} bridge{{ bridges.length !== 1 ? 's' : '' }}. Select one to pair.</p>
            <ul class="bridge-list">
              <li
                v-for="bridge in bridges"
                :key="bridge.id"
                class="bridge-item"
                :class="{ selected: selectedBridge?.id === bridge.id }"
                @click="selectedBridge = bridge"
              >
                <div class="bridge-id">{{ bridge.id }}</div>
                <div class="bridge-ip">{{ bridge.ip }}</div>
              </li>
            </ul>
            <div v-if="selectedBridge" class="pair-instruction">
              Press the link button on your Hue bridge, then click Pair.
            </div>
            <div class="modal-actions">
              <button class="btn-secondary" @click="step = 'discover'">Back</button>
              <button class="btn-primary" :disabled="!selectedBridge || pairing" @click="doPair">
                <span v-if="pairing" class="spinner" />
                {{ pairing ? 'Pairing...' : 'Pair' }}
              </button>
            </div>
            <div v-if="pairError" class="error-msg">{{ pairError }}</div>
          </div>

          <!-- Step 3: Success -->
          <div v-else-if="step === 'success'" class="modal-body modal-body--center">
            <div class="success-icon">
              <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
                <circle cx="20" cy="20" r="19" stroke="#22C55E" stroke-width="2" />
                <path d="M12 20l6 6 10-12" stroke="#22C55E" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </div>
            <p class="success-text">Bridge paired! Lights are coming online.</p>
            <button class="btn-primary" @click="handlePairedClose">Done</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'paired'): void
}>()

const { discoverBridges, pair } = useHue()

const step = ref<'discover' | 'select' | 'success'>('discover')
const bridges = ref<{ id: string; ip: string }[]>([])
const selectedBridge = ref<{ id: string; ip: string } | null>(null)
const discovering = ref(false)
const pairing = ref(false)
const discoverError = ref('')
const pairError = ref('')

// Reset state when dialog opens
watch(() => props.show, (val) => {
  if (val) {
    step.value = 'discover'
    bridges.value = []
    selectedBridge.value = null
    discoverError.value = ''
    pairError.value = ''
  }
})

async function findBridges() {
  discovering.value = true
  discoverError.value = ''
  try {
    const found = await discoverBridges()
    if (found.length === 0) {
      discoverError.value = 'No bridges found. Make sure your bridge is on the same network.'
    } else {
      bridges.value = found
      step.value = 'select'
    }
  } catch {
    discoverError.value = 'Discovery failed. Check your connection.'
  } finally {
    discovering.value = false
  }
}

async function doPair() {
  if (!selectedBridge.value) return
  pairing.value = true
  pairError.value = ''
  try {
    const ok = await pair(selectedBridge.value.ip)
    if (ok) {
      step.value = 'success'
    } else {
      pairError.value = 'Pairing failed. Did you press the link button?'
    }
  } catch {
    pairError.value = 'Pairing failed. Try again.'
  } finally {
    pairing.value = false
  }
}

function handlePairedClose() {
  emit('paired')
  emit('close')
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(4px);
}

.modal {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-xl);
  width: 100%;
  max-width: 440px;
  padding: 28px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.5);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text);
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: var(--color-surface-border);
  color: var(--color-text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s ease, color 0.15s ease;
}

.close-btn:hover {
  background: var(--color-surface-hover);
  color: var(--color-text);
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-body--center {
  align-items: center;
  text-align: center;
}

.step-description {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
}

.bridge-list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bridge-item {
  background: var(--color-bg);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-md);
  padding: 12px 16px;
  cursor: pointer;
  transition: border-color 0.15s ease, background 0.15s ease;
}

.bridge-item:hover {
  background: var(--color-surface-hover);
  border-color: rgba(0, 101, 211, 0.4);
}

.bridge-item.selected {
  border-color: #0065D3;
  background: rgba(0, 101, 211, 0.08);
}

.bridge-id {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text);
  font-family: var(--font-mono);
}

.bridge-ip {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 2px;
}

.pair-instruction {
  font-size: 13px;
  color: var(--color-warning);
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.2);
  border-radius: var(--radius-md);
  padding: 10px 14px;
  line-height: 1.5;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s ease, opacity 0.15s ease;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  display: inline-flex;
  align-items: center;
  background: var(--color-surface-border);
  color: var(--color-text);
  border: none;
  border-radius: var(--radius-md);
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s ease;
}

.btn-secondary:hover {
  background: var(--color-surface-hover);
}

.error-msg {
  font-size: 13px;
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: var(--radius-md);
  padding: 10px 14px;
}

.success-icon {
  margin: 8px 0;
}

.success-text {
  font-size: 15px;
  color: var(--color-text);
  font-weight: 500;
}

/* Spinner */
.spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Modal fade transition */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-active .modal,
.modal-fade-leave-active .modal {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-from .modal,
.modal-fade-leave-to .modal {
  transform: scale(0.96) translateY(8px);
  opacity: 0;
}
</style>
