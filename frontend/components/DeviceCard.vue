<template>
  <div
    class="device-card"
    :class="`accent-${brandKey}`"
    @click="handleCardClick"
  >
    <!-- Header: brand logo + name block + toggle -->
    <div class="card-header">
      <BrandLogo :brand="brandKey" :size="36" class="card-brand-logo" />
      <div class="card-title-block">
        <span class="device-name">{{ device.name || device.hostname || 'Unknown Device' }}</span>
        <span class="device-meta-line">
          <code class="device-ip">{{ device.ip }}</code>
          <template v-if="device.manufacturer">
            <span class="meta-sep">·</span>
            <span class="device-manufacturer">{{ device.manufacturer }}</span>
          </template>
        </span>
        <span v-if="device.model || device.device_type" class="device-detail-line">
          <template v-if="device.model">{{ device.model }}</template>
          <template v-if="device.model && device.device_type">
            <span class="meta-sep">·</span>
          </template>
          <template v-if="device.device_type">{{ device.device_type }}</template>
        </span>
      </div>

      <!-- Toggle for Kasa switches/dimmers/fans -->
      <div
        v-if="isKasaControllable"
        class="toggle-track"
        :class="{ on: currentOn }"
        @click.stop="handleToggle"
        role="switch"
        :aria-checked="currentOn"
      >
        <span class="toggle-thumb" />
      </div>
    </div>

    <!-- Brightness slider for Kasa dimmers -->
    <div
      v-if="isKasaDimmer"
      class="dimmer-row"
      @click.stop
    >
      <input
        type="range"
        min="0"
        max="100"
        :value="currentBrightness"
        class="brightness-slider"
        @input="handleBrightnessInput"
      />
      <span class="brightness-label">{{ currentBrightness }}%</span>
    </div>

    <!-- Protocol pills -->
    <div v-if="device.protocols && device.protocols.length" class="protocol-pills">
      <span
        v-for="proto in device.protocols"
        :key="proto"
        class="pill"
        :class="`pill-${proto.toLowerCase()}`"
      >
        {{ proto }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
interface ScanDevice {
  ip: string
  mac: string
  hostname: string
  name: string
  manufacturer: string
  model: string
  device_type: string
  category: string
  protocols: string[]
  services: string[]
  open_ports: number[]
  metadata: Record<string, string>
}

interface KasaLiveState {
  ip: string
  on: boolean
  brightness: number
  device_type: string
}

const props = withDefaults(defineProps<{
  device: ScanDevice
  liveState?: KasaLiveState | null
}>(), {
  liveState: null,
})

const emit = defineEmits<{
  select: [device: ScanDevice]
  toggle: [payload: { ip: string; on: boolean }]
  brightness: [payload: { ip: string; brightness: number }]
}>()

const brandKey = computed(() => {
  const mfr = props.device.manufacturer?.toLowerCase().trim() || ''
  const dt = props.device.device_type?.toLowerCase().trim() || ''
  if (mfr && mfr !== 'unknown') return mfr
  if (dt) return dt
  return ''
})

// Is this a Kasa device with controls?
const isKasaControllable = computed(() => {
  const proto = props.device.protocols?.map(p => p.toLowerCase()) || []
  const dt = props.device.device_type?.toLowerCase() || ''
  return proto.includes('kasa') || dt === 'switch' || dt === 'dimmer' || dt === 'fan'
})

const isKasaDimmer = computed(() => {
  const dt = props.device.device_type?.toLowerCase() || ''
  return isKasaControllable.value && dt === 'dimmer'
})

// Live state wins over scan metadata
const currentOn = computed(() => {
  if (props.liveState) return props.liveState.on
  return false
})

const currentBrightness = computed(() => {
  if (props.liveState) return props.liveState.brightness ?? 100
  return 100
})

let brightnessDebounce: ReturnType<typeof setTimeout> | null = null

function handleToggle() {
  emit('toggle', { ip: props.device.ip, on: !currentOn.value })
}

function handleBrightnessInput(e: Event) {
  const val = parseInt((e.target as HTMLInputElement).value, 10)
  if (brightnessDebounce) clearTimeout(brightnessDebounce)
  brightnessDebounce = setTimeout(() => {
    emit('brightness', { ip: props.device.ip, brightness: val })
    brightnessDebounce = null
  }, 300)
}

function handleCardClick() {
  navigateTo(`/device/${encodeURIComponent(props.device.ip)}`)
}
</script>

<style scoped>
.device-card {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-left: 3px solid transparent;
  border-radius: var(--radius-lg);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  cursor: pointer;
  transition: background 0.25s ease, border-color 0.25s ease, border-left-color 0.25s ease, box-shadow 0.25s ease;
  animation: card-appear 0.35s ease forwards;
}

.device-card.selected {
  border-color: var(--color-accent);
  box-shadow: 0 0 16px rgba(99, 102, 241, 0.2);
}

.device-card:hover {
  background: var(--color-surface-hover);
  border-color: rgba(99, 102, 241, 0.25);
}

/* Brand-specific left border accent on hover */
.device-card.accent-philips:hover,
.device-card.accent-signify:hover,
.device-card.accent-hue:hover        { border-left-color: #0065D3; }
.device-card.accent-tp-link:hover,
.device-card.accent-tplink:hover      { border-left-color: #00A4A6; }
.device-card.accent-nvidia:hover      { border-left-color: #76B900; }
.device-card.accent-lg:hover          { border-left-color: #A50034; }
.device-card.accent-google:hover      { border-left-color: #4285F4; }
.device-card.accent-apple:hover       { border-left-color: #E8E8EA; }
.device-card.accent-yamaha:hover      { border-left-color: #7C3AED; }
.device-card.accent-brilliant:hover   { border-left-color: #2563EB; }
.device-card.accent-sunpower:hover    { border-left-color: #F59E0B; }
.device-card.accent-enphase:hover     { border-left-color: #E65100; }
.device-card.accent-jellyfish:hover   { border-left-color: #06B6D4; }
.device-card.accent-sonos:hover       { border-left-color: #E8E8EA; }
.device-card.accent-samsung:hover     { border-left-color: #1428A0; }
.device-card.accent-dell:hover        { border-left-color: #007DB8; }
.device-card.accent-hue_bridge:hover,
.device-card.accent-hue_sync:hover    { border-left-color: #0065D3; }
.device-card.accent-shield_tv:hover   { border-left-color: #76B900; }
.device-card.accent-brilliant_switch:hover { border-left-color: #2563EB; }
.device-card.accent-solar_gateway:hover { border-left-color: #F59E0B; }
.device-card.accent-av_receiver:hover { border-left-color: #7C3AED; }
.device-card.accent-speaker:hover     { border-left-color: #E8E8EA; }
.device-card.accent-thread_border_router:hover { border-left-color: #4285F4; }

@keyframes card-appear {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Header layout */
.card-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.card-brand-logo {
  margin-top: 2px;
  flex-shrink: 0;
}

.card-title-block {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
  flex: 1;
}

.device-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
}

.device-meta-line {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.device-ip {
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--color-text-secondary);
  background: rgba(255, 255, 255, 0.05);
  padding: 1px 6px;
  border-radius: var(--radius-sm);
  font-style: normal;
}

.device-manufacturer {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.device-detail-line {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11.5px;
  color: var(--color-text-tertiary);
}

.meta-sep {
  color: var(--color-text-tertiary);
  opacity: 0.5;
  font-size: 11px;
}

/* Toggle switch */
.toggle-track {
  width: 44px;
  height: 24px;
  border-radius: 12px;
  background: var(--color-surface-border);
  position: relative;
  cursor: pointer;
  transition: background 0.2s;
  flex-shrink: 0;
  margin-top: 4px;
}

.toggle-track.on {
  background: var(--color-accent);
}

.toggle-thumb {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.toggle-track.on .toggle-thumb {
  transform: translateX(20px);
}

/* Dimmer slider row */
.dimmer-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brightness-slider {
  flex: 1;
  -webkit-appearance: none;
  appearance: none;
  height: 4px;
  border-radius: 2px;
  background: var(--color-surface-border);
  outline: none;
  cursor: pointer;
}

.brightness-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
  transition: transform 0.15s ease;
}

.brightness-slider::-webkit-slider-thumb:hover {
  transform: scale(1.15);
}

.brightness-slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
  border: none;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
}

.brightness-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  min-width: 36px;
  text-align: right;
  font-family: var(--font-mono);
}

/* Protocol pills */
.protocol-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.pill {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 99px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  background: rgba(255, 255, 255, 0.07);
  color: var(--color-text-secondary);
}

.pill-kasa       { background: rgba(34, 197, 94, 0.15);   color: #4ade80; }
.pill-cast       { background: rgba(59, 130, 246, 0.15);  color: #60a5fa; }
.pill-mdns       { background: rgba(168, 85, 247, 0.15);  color: #c084fc; }
.pill-airplay    { background: rgba(229, 231, 235, 0.12); color: #e5e7eb; }
.pill-spotify    { background: rgba(34, 197, 94, 0.15);   color: #4ade80; }
.pill-thread     { background: rgba(245, 158, 11, 0.15);  color: #fbbf24; }
.pill-http       { background: rgba(107, 114, 128, 0.2);  color: #9ca3af; }
.pill-hue        { background: rgba(0, 101, 211, 0.2);    color: #60a5fa; }
.pill-matter     { background: rgba(99, 102, 241, 0.15);  color: #a5b4fc; }
</style>
