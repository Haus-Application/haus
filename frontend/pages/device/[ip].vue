<template>
  <div class="device-page">
    <header class="device-header">
      <NuxtLink to="/" class="back-btn">← Back</NuxtLink>
      <span v-if="probe" class="status-badge" :class="isConnected ? 'connected' : 'not-connected'">
        {{ isConnected ? 'Connected' : 'Not Connected' }}
      </span>
    </header>

    <div v-if="loading && !device" class="loading-state">
      <div class="loader" />
      <p>Connecting to device...</p>
    </div>

    <template v-if="device">
      <!-- Device hero -->
      <div class="device-hero">
        <BrandLogo :brand="brandKey" :size="56" />
        <div class="hero-info">
          <h1 class="hero-name">{{ probe?.name || device.name || device.ip }}</h1>
          <div class="hero-meta">
            <code class="hero-ip">{{ device.ip }}</code>
            <span v-if="device.manufacturer" class="meta-sep">·</span>
            <span v-if="device.manufacturer" class="hero-mfr">{{ device.manufacturer }}</span>
          </div>
          <div v-if="device.model || device.device_type" class="hero-detail">
            <span v-if="device.model">{{ device.model }}</span>
            <span v-if="device.model && device.device_type" class="meta-sep">·</span>
            <span v-if="device.device_type">{{ device.device_type }}</span>
          </div>
          <div v-if="probe?.integration" class="hero-integration">
            <span class="integration-pill" :class="`int-${probe.integration}`">{{ probe.integration }}</span>
            <span v-if="!probe.reachable" class="unreachable-badge">unreachable</span>
          </div>
        </div>
      </div>

      <!-- Capabilities — what this device can do -->
      <section v-if="probe?.capabilities?.length || probe?.api" class="capabilities-section">
        <h2 class="section-title">What This Device Can Do</h2>
        <div class="cap-card">
          <!-- Capability items -->
          <div v-for="cap in capabilityList" :key="cap.id" class="cap-item">
            <span class="cap-icon" v-html="cap.icon" />
            <div class="cap-info">
              <span class="cap-name">{{ cap.name }}</span>
              <span class="cap-desc">{{ cap.description }}</span>
            </div>
            <span class="cap-status" :class="isConnected ? 'active' : 'inactive'">
              {{ isConnected ? 'Ready' : 'Needs Setup' }}
            </span>
          </div>

          <!-- API Protocol info -->
          <div v-if="probe?.api" class="cap-item cap-api">
            <span class="cap-icon">⚡</span>
            <div class="cap-info">
              <span class="cap-name">{{ probe.api.protocol }} on port {{ probe.api.port }}</span>
              <span class="cap-desc">{{ probe.api.auth_method === 'none' ? 'No authentication required' : `Auth: ${probe.api.auth_method}` }}</span>
            </div>
          </div>
        </div>

        <!-- API Endpoints (collapsible) -->
        <details v-if="probe?.api?.endpoints?.length" class="endpoints-details">
          <summary class="endpoints-summary">API Endpoints ({{ probe.api.endpoints.length }})</summary>
          <div class="endpoints-list">
            <div v-for="(ep, i) in probe.api.endpoints" :key="i" class="endpoint-row">
              <div class="ep-header">
                <code class="ep-method">{{ ep.method }}</code>
                <code class="ep-path">{{ ep.path }}</code>
              </div>
              <p class="ep-desc">{{ ep.description }}</p>
              <code v-if="ep.example" class="ep-example">{{ ep.example }}</code>
            </div>
          </div>
        </details>

        <a v-if="probe?.api?.doc_url" :href="probe.api.doc_url" target="_blank" class="doc-link">View Full Documentation →</a>
      </section>

      <!-- What we found (discovered, not yet integrated) -->
      <section v-if="probe?.status === 'discovered' && probe?.fingerprints?.length" class="discovery-section">
        <h2 class="section-title">What We Found</h2>
        <div class="discovery-card">
          <div v-for="fp in probe.fingerprints" :key="fp.port" class="fingerprint-row">
            <div class="fp-port">Port {{ fp.port }}</div>
            <div class="fp-details">
              <span v-if="fp.title" class="fp-title">{{ fp.title }}</span>
              <span v-if="fp.server" class="fp-server">{{ fp.server }}</span>
              <span v-if="fp.has_login" class="fp-login-badge">Login Required</span>
            </div>
            <a v-if="fp.url" :href="fp.url" target="_blank" class="fp-link">Open →</a>
          </div>
        </div>
      </section>

      <!-- Setup needed — link button or password -->
      <section v-if="probe?.setup_needed" class="setup-section">
        <div class="setup-card">
          <h3 class="setup-title">{{ probe.setup_needed.title }}</h3>
          <p class="setup-desc">{{ probe.setup_needed.description }}</p>

          <!-- Password auth (SunPower, etc.) -->
          <template v-if="probe.setup_needed.type === 'password'">
            <div class="auth-form">
              <input v-model="authPassword" type="password" placeholder="Password" class="auth-input" />
              <button class="setup-btn" @click="handleSetup" :disabled="setupLoading">
                {{ setupLoading ? 'Connecting...' : probe.setup_needed.action_label }}
              </button>
            </div>
          </template>

          <!-- Link button (Hue) -->
          <template v-else>
            <button class="setup-btn" @click="handleSetup" :disabled="setupLoading">
              {{ setupLoading ? 'Connecting...' : probe.setup_needed.action_label }}
            </button>
          </template>

          <p v-if="setupError" class="setup-error">{{ setupError }}</p>
          <p v-if="setupSuccess" class="setup-success">Connected! Refreshing...</p>
        </div>
      </section>

      <!-- Offline / unreachable — but still show what we know -->
      <section v-if="probe?.status === 'offline' || (probe?.status === 'discovered' && !probe?.fingerprints?.length)" class="setup-section">
        <div class="offline-card">
          <h3 class="setup-title">{{ deviceHint ? deviceHint.title : 'Device Found on Network' }}</h3>
          <p class="setup-desc">
            {{ deviceHint ? deviceHint.description : 'This device was detected on your network but isn\'t responding to standard probes. It may use Bluetooth, a cloud-only protocol, or a non-standard port.' }}
          </p>
          <div v-if="deviceHint?.suggestions" class="hint-suggestions">
            <p v-for="s in deviceHint.suggestions" :key="s" class="hint-item">{{ s }}</p>
          </div>
          <button class="setup-btn" @click="runProbe" style="margin-top:12px">Retry Probe</button>
        </div>
      </section>

      <!-- JellyFish controls -->
      <section v-if="probe?.integration === 'jellyfish'" class="controls-section">
        <h2 class="section-title">Controls</h2>
        <div class="controls-card">
          <div class="control-row">
            <span class="control-label">Power (All Zones)</span>
            <div class="control-right">
              <span class="state-text" :class="{ on: jfOn }">{{ jfOn ? 'ON' : 'OFF' }}</span>
              <div class="toggle-track" :class="{ on: jfOn }" @click="jfToggle" role="switch">
                <span class="toggle-thumb" />
              </div>
            </div>
          </div>
        </div>
      </section>

      <section v-if="probe?.integration === 'jellyfish' && probe?.state?.zone_names" class="controls-section">
        <h2 class="section-title">Zones <span class="section-count">{{ probe.state.zone_count }}</span></h2>
        <div class="controls-card">
          <div v-for="zone in probe.state.zone_names" :key="zone" class="control-row">
            <span class="control-label">{{ zone }}</span>
            <div class="control-right">
              <button class="speed-btn active" @click="jfZoneOn(zone)" style="width:auto;padding:0 12px;">On</button>
              <button class="speed-btn" @click="jfZoneOff(zone)" style="width:auto;padding:0 12px;">Off</button>
            </div>
          </div>
        </div>
      </section>

      <section v-if="probe?.integration === 'jellyfish' && jfPatterns.length" class="controls-section">
        <h2 class="section-title">Patterns <span class="section-count">{{ jfPatterns.length }}</span></h2>
        <div class="patterns-grid">
          <button v-for="p in jfPatterns" :key="p.path" class="pattern-chip pattern-btn"
            :class="{ active: jfActivePattern === p.path }"
            @click="jfPlayPattern(p.path)">{{ p.name }}</button>
        </div>
      </section>

      <!-- SunPower solar data -->
      <section v-if="probe?.integration === 'sunpower' && probe?.status === 'connected'" class="controls-section">
        <h2 class="section-title">Live Solar Data</h2>
        <div class="solar-grid">
          <div class="solar-card solar-hero">
            <span class="solar-label">Solar Production</span>
            <span class="solar-value solar-green" style="font-size:28px">{{ probe.state.production_kw }} kW</span>
          </div>
          <div class="solar-card">
            <span class="solar-label">House</span>
            <span class="solar-value">{{ probe.state.consumption_kw }} kW</span>
          </div>
          <div class="solar-card">
            <span class="solar-label">Grid</span>
            <span class="solar-value" :class="probe.state.exporting ? 'solar-green' : 'solar-red'">
              {{ probe.state.exporting ? 'Exporting' : 'Importing' }} {{ Math.abs(parseFloat(probe.state.grid_kw || '0')).toFixed(1) }} kW
            </span>
          </div>
          <div class="solar-card" v-if="probe.state.battery_soc && probe.state.battery_soc !== '0'">
            <span class="solar-label">Battery</span>
            <span class="solar-value">{{ probe.state.battery_soc }}%</span>
          </div>
          <div class="solar-card" v-if="probe.state.lifetime_kwh">
            <span class="solar-label">Lifetime Production</span>
            <span class="solar-value solar-green">{{ Number(probe.state.lifetime_kwh).toLocaleString() }} kWh</span>
          </div>
          <div class="solar-card" v-if="probe.state.panel_count">
            <span class="solar-label">Panels</span>
            <span class="solar-value">{{ probe.state.panel_count }}</span>
          </div>
        </div>
      </section>

      <!-- Yamaha receiver state -->
      <section v-if="probe?.integration === 'yamaha' && probe?.state?.model" class="controls-section">
        <h2 class="section-title">Receiver</h2>
        <div class="controls-card">
          <div class="control-row">
            <span class="control-label">Model</span>
            <span class="info-value">{{ probe.state.model }}</span>
          </div>
          <div v-if="probe.state.input" class="control-row">
            <span class="control-label">Input</span>
            <span class="info-value">{{ probe.state.input }}</span>
          </div>
          <div v-if="probe.state.volume !== undefined" class="control-row">
            <span class="control-label">Volume</span>
            <span class="info-value">{{ probe.state.volume }}</span>
          </div>
        </div>
      </section>

      <!-- Fingerprints (discovered but not yet integrated) -->
      <section v-if="probe?.fingerprints?.length && probe?.status === 'discovered'" class="controls-section">
        <h2 class="section-title">What We Found</h2>
        <div class="controls-card">
          <div v-for="fp in probe.fingerprints" :key="fp.port" class="control-row">
            <div class="fp-left">
              <code class="fp-port">:{{ fp.port }}</code>
              <span v-if="fp.title" class="fp-title">{{ fp.title }}</span>
              <span v-else-if="fp.server" class="fp-title">{{ fp.server }}</span>
              <span v-if="fp.has_login" class="login-badge">Login Required</span>
            </div>
            <a v-if="fp.url" :href="fp.url" target="_blank" class="fp-link">Open →</a>
          </div>
        </div>
      </section>

      <!-- Dynamic controls from probe (not for jellyfish — it has its own section) -->
      <section v-if="probe?.actions?.length && probe?.integration !== 'jellyfish'" class="controls-section">
        <h2 class="section-title">Controls</h2>
        <div class="controls-card">
          <template v-for="action in probe.actions" :key="action.id">
            <!-- Toggle -->
            <div v-if="action.type === 'toggle'" class="control-row">
              <span class="control-label">{{ action.label }}</span>
              <div class="control-right">
                <span class="state-text" :class="{ on: liveToggle }">{{ liveToggle ? 'ON' : 'OFF' }}</span>
                <div class="toggle-track" :class="{ on: liveToggle }" @click="doToggle" role="switch">
                  <span class="toggle-thumb" />
                </div>
              </div>
            </div>

            <!-- Slider -->
            <template v-if="action.type === 'slider'">
              <div class="control-row">
                <span class="control-label">{{ action.label }}</span>
                <span class="brightness-value">{{ liveBrightness }}%</span>
              </div>
              <div class="slider-row">
                <input type="range" :min="action.min" :max="action.max" :value="liveBrightness"
                  class="brightness-slider" @input="doBrightness" />
              </div>
            </template>

            <!-- Button group (fan speed) -->
            <div v-if="action.type === 'buttons'" class="control-row">
              <span class="control-label">{{ action.label }}</span>
              <div class="speed-buttons">
                <button v-for="opt in action.options" :key="opt.value"
                  class="speed-btn" :class="{ active: liveFanSpeed === opt.value }"
                  @click="doFanSpeed(opt.value)">{{ opt.label }}</button>
              </div>
            </div>
          </template>
        </div>
      </section>

      <!-- Hue lights (if bridge is connected) -->
      <section v-if="probe?.status === 'connected' && probe?.integration === 'hue'" class="controls-section">
        <h2 class="section-title">
          Lights
          <span class="section-count">{{ probe.state?.lights_on || 0 }} / {{ probe.state?.light_count || 0 }} on</span>
        </h2>
        <p class="section-hint">Use the chat below to control individual lights, rooms, and scenes.</p>
      </section>

      <!-- Device info -->
      <section class="info-section">
        <h2 class="section-title">Details</h2>
        <div class="info-card">
          <div class="info-row" v-if="device.mac"><span class="info-label">MAC</span><code class="info-value">{{ device.mac }}</code></div>
          <div class="info-row" v-if="device.category"><span class="info-label">Category</span><span class="info-value">{{ device.category }}</span></div>
          <div class="info-row" v-if="probe?.capabilities?.length"><span class="info-label">Capabilities</span>
            <div class="cap-pills"><span v-for="c in probe.capabilities" :key="c" class="cap-pill">{{ c }}</span></div>
          </div>
          <div class="info-row" v-if="device.protocols?.length"><span class="info-label">Protocols</span>
            <div class="protocol-pills"><span v-for="p in device.protocols" :key="p" class="pill" :class="`pill-${p}`">{{ p }}</span></div>
          </div>
          <div class="info-row" v-if="device.open_ports?.length"><span class="info-label">Open Ports</span><code class="info-value">{{ device.open_ports.join(', ') }}</code></div>
        </div>
      </section>

      <!-- AI Chat -->
      <section class="chat-section">
        <h2 class="section-title">Chat with {{ probe?.name || device.name || 'Device' }}</h2>
        <div class="chat-card">
          <div v-if="chatMessages.length === 0" class="chat-suggestions">
            <button v-for="s in suggestions" :key="s" class="suggestion-chip" @click="sendChat(s)">{{ s }}</button>
          </div>
          <div class="chat-messages" ref="messagesRef">
            <div v-for="(msg, i) in chatMessages" :key="i" class="chat-msg" :class="msg.role">
              <div class="msg-bubble">{{ msg.content }}</div>
              <div v-if="msg.toolCalls?.length" class="msg-tools">
                <span v-for="(tc, j) in msg.toolCalls" :key="j" class="tool-badge">{{ tc.tool }}</span>
              </div>
            </div>
            <div v-if="chatLoading" class="chat-msg assistant">
              <div class="msg-bubble typing"><span class="dot" /><span class="dot" /><span class="dot" /></div>
            </div>
          </div>
          <div class="chat-input-row">
            <input v-model="chatInput" @keydown.enter="sendChat(chatInput)" placeholder="Ask this device anything..." class="chat-input" :disabled="chatLoading" />
            <button @click="sendChat(chatInput)" class="chat-send" :disabled="!chatInput.trim() || chatLoading">→</button>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const deviceIP = route.params.ip as string

const device = ref<any>(null)
const probe = ref<any>(null)
const loading = ref(true)

// Live Kasa state from WebSocket
const { devices: kasaDevices, toggleDevice, setBrightness, setFanSpeed } = useKasa()
const kasaLive = computed(() => kasaDevices.value.find(d => d.ip === deviceIP))
const liveToggle = computed(() => kasaLive.value?.on ?? probe.value?.state?.on ?? false)
const liveBrightness = computed(() => kasaLive.value?.brightness ?? probe.value?.state?.brightness ?? 100)
const liveFanSpeed = computed(() => kasaLive.value?.fan_speed ?? probe.value?.state?.fan_speed ?? 0)

let brightnessTimer: ReturnType<typeof setTimeout> | null = null

function doToggle() { toggleDevice(deviceIP, !liveToggle.value) }
function doBrightness(e: Event) {
  const val = parseInt((e.target as HTMLInputElement).value, 10)
  if (brightnessTimer) clearTimeout(brightnessTimer)
  brightnessTimer = setTimeout(() => setBrightness(deviceIP, val), 300)
}
function doFanSpeed(speed: number) { setFanSpeed(deviceIP, speed) }

// Setup (e.g. Hue pairing, SunPower auth)
const setupLoading = ref(false)
const setupError = ref('')
const setupSuccess = ref(false)
const authPassword = ref('')

async function handleSetup() {
  if (!probe.value?.setup_needed) return
  setupLoading.value = true
  setupError.value = ''
  try {
    const setupType = probe.value.setup_needed.type
    let body: any

    if (setupType === 'password') {
      // Password auth (SunPower, etc.)
      body = { password: authPassword.value }
    } else {
      // Link button (Hue)
      body = { bridge_ip: deviceIP }
    }

    const res = await fetch(probe.value.setup_needed.action, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    const data = await res.json()
    if (data.error) {
      setupError.value = data.error
    } else {
      setupSuccess.value = true
      setTimeout(() => runProbe(), 2000)
    }
  } catch {
    setupError.value = 'Connection failed. Check that the device is reachable.'
  }
  setupLoading.value = false
}

// Brand key
const CAP_MAP: Record<string, { name: string; description: string; icon: string }> = {
  on_off: { name: 'Power Control', description: 'Turn on and off', icon: '⏻' },
  brightness: { name: 'Brightness', description: 'Adjust brightness 0-100%', icon: '☀' },
  fan_speed: { name: 'Fan Speed', description: 'Set speed 1-4 (low to max)', icon: '⊙' },
  zones: { name: 'Lighting Zones', description: 'Control individual LED zones', icon: '◉' },
  patterns: { name: 'Light Patterns', description: 'Play preset color patterns', icon: '✦' },
  lights: { name: 'Light Control', description: 'Control individual Hue lights', icon: '💡' },
  rooms: { name: 'Room Control', description: 'Control all lights in a room', icon: '⬡' },
  scenes: { name: 'Scene Recall', description: 'Activate saved lighting scenes', icon: '✧' },
  color: { name: 'Color Control', description: 'Change light color (XY color space)', icon: '⬮' },
  solar_production: { name: 'Solar Production', description: 'Real-time solar panel output', icon: '☀' },
  solar_consumption: { name: 'House Consumption', description: 'Real-time energy usage', icon: '⚡' },
  grid_power: { name: 'Grid Status', description: 'Importing/exporting to grid', icon: '⇄' },
  battery: { name: 'Battery', description: 'Battery state of charge', icon: '▮' },
  power: { name: 'Power', description: 'Turn on/off/standby', icon: '⏻' },
  volume: { name: 'Volume', description: 'Adjust volume level', icon: '♪' },
  input_select: { name: 'Input Select', description: 'Switch between HDMI, Bluetooth, etc.', icon: '⎆' },
  mute: { name: 'Mute', description: 'Mute/unmute audio', icon: '🔇' },
  media_info: { name: 'Media Info', description: 'View device name and status', icon: 'ⓘ' },
}

const capabilityList = computed(() => {
  if (!probe.value?.capabilities) return []
  return probe.value.capabilities.map((cap: string) => ({
    id: cap,
    ...(CAP_MAP[cap] || { name: cap, description: '', icon: '•' }),
  }))
})

// Smart hints based on hostname/manufacturer for devices we can't probe
const deviceHint = computed(() => {
  if (!device.value) return null
  const hostname = (device.value.hostname || device.value.name || '').toLowerCase()
  const mfr = (device.value.manufacturer || '').toLowerCase()

  if (hostname.includes('ge_light') || hostname.includes('cync') || mfr.includes('ge')) {
    return {
      title: 'GE / Cync Smart Light',
      description: 'This is a GE (Cync) smart light. These devices typically use Bluetooth for local control and the Cync cloud for WiFi control.',
      suggestions: [
        '• Download the Cync app to pair and control this light',
        '• Cync lights can be added to Google Home or Alexa',
        '• Some models support local control via Bluetooth mesh',
        '• Try asking the AI chat below for help connecting',
      ],
    }
  }
  if (hostname.includes('ring') || mfr.includes('ring')) {
    return {
      title: 'Ring Security Device',
      description: 'Ring devices use cloud-only communication via the Ring API. Local access requires OAuth authentication through Ring.',
      suggestions: [
        '• Control via the Ring app or Ring.com',
        '• Ring devices require cloud authentication (no local API)',
        '• Can be integrated via Ring OAuth if you have a Ring account',
      ],
    }
  }
  if (hostname.includes('nest') || hostname.includes('google') || mfr.includes('google')) {
    return {
      title: 'Google / Nest Device',
      description: 'Google and Nest devices use cloud protocols. Some support local Cast control on port 8008.',
      suggestions: [
        '• Control via the Google Home app',
        '• Cast-enabled devices can be probed on port 8008',
        '• Nest thermostats require Google SDM API access',
      ],
    }
  }
  if (hostname.includes('alexa') || hostname.includes('echo') || hostname.includes('amazon')) {
    return {
      title: 'Amazon Echo / Alexa Device',
      description: 'Amazon Echo devices use cloud-only protocols.',
      suggestions: [
        '• Control via the Alexa app',
        '• No local API available for direct control',
      ],
    }
  }
  // Generic hint for any device on the network
  return {
    title: 'Device Found on Network',
    description: 'This device was detected on your network but isn\'t responding to standard probes. It may use Bluetooth, a cloud-only protocol, or a non-standard port.',
    suggestions: [
      '• Try asking the AI chat below — it may know about this device type',
      '• Check if the device has a companion app for setup',
      '• The device might need to be on the same WiFi network',
    ],
  }
})

const isConnected = computed(() => {
  if (!probe.value) return false
  return probe.value.status === 'connected'
})

const brandKey = computed(() => {
  if (!device.value) return ''
  const mfr = device.value.manufacturer?.toLowerCase().trim() || ''
  const dt = device.value.device_type?.toLowerCase().trim() || ''
  return (mfr && mfr !== 'unknown') ? mfr : dt || ''
})

// Smart suggestions based on probe results
const suggestions = computed(() => {
  if (!probe.value) return ['What are you?']
  const caps = probe.value.capabilities || []
  const s: string[] = []
  if (caps.includes('on_off')) s.push('Are you on?', 'Turn on')
  if (caps.includes('brightness')) s.push('Set to 50%')
  if (caps.includes('fan_speed')) s.push('Set speed to 2')
  if (caps.includes('lights')) s.push('List lights', 'What scenes are available?')
  if (caps.includes('color')) s.push('Set color to warm')
  if (s.length === 0) s.push('What are you?', 'What can you do?')
  return s
})

// Load device from DB + probe it
async function runProbe() {
  try {
    const res = await fetch(`/api/devices/${encodeURIComponent(deviceIP)}/probe`)
    if (res.ok) probe.value = await res.json()
  } catch { /* server down */ }
}

onMounted(async () => {
  // Load from DB
  try {
    const res = await fetch('/api/devices')
    if (res.ok) {
      const all = await res.json()
      device.value = all.find((d: any) => d.ip === deviceIP) || { ip: deviceIP, name: deviceIP }
    }
  } catch {
    device.value = { ip: deviceIP, name: deviceIP }
  }

  // Probe the device in real-time
  await runProbe()
  loading.value = false
})

// JellyFish controls
const jfOn = ref(false)
const jfActivePattern = ref('')
const jfPatterns = computed(() => (probe.value?.state?.patterns || []).filter((p: any) => p?.name?.trim()))

async function jfCommand(action: string, zones?: string[], pattern?: string) {
  try {
    await fetch(`/api/devices/${encodeURIComponent(deviceIP)}/jellyfish`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action, zones, pattern }),
    })
  } catch { /* connection failed */ }
}

function jfToggle() {
  jfOn.value = !jfOn.value
  jfCommand(jfOn.value ? 'on' : 'off')
}

function jfZoneOn(zone: string) {
  jfCommand('on', [zone], jfActivePattern.value || 'Accent/White')
}

function jfZoneOff(zone: string) {
  jfCommand('off', [zone])
}

function jfPlayPattern(pattern: string) {
  jfActivePattern.value = pattern
  jfOn.value = true
  jfCommand('pattern', undefined, pattern)
}

// Chat logic
interface ChatMsg { role: 'user' | 'assistant'; content: string; toolCalls?: any[] }
const chatMessages = ref<ChatMsg[]>([])
const chatLoading = ref(false)
const chatInput = ref('')
const chatHistory = ref<any[]>([])
const messagesRef = ref<HTMLElement | null>(null)

async function sendChat(text: string) {
  if (!text.trim() || chatLoading.value) return
  const msg = text.trim()
  chatInput.value = ''
  chatMessages.value.push({ role: 'user', content: msg })
  chatLoading.value = true
  try {
    const res = await fetch('/api/chat/device', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        device: {
          ip: device.value.ip, name: device.value.name || device.value.ip,
          manufacturer: device.value.manufacturer || '', model: device.value.model || '',
          device_type: device.value.device_type || '', category: device.value.category || '',
          protocols: device.value.protocols || [],
          api_docs: probe.value?.api_docs || '',
        },
        message: msg, history: chatHistory.value,
      }),
    })
    const data = await res.json()
    if (data.error) {
      chatMessages.value.push({ role: 'assistant', content: data.error })
    } else {
      chatMessages.value.push({ role: 'assistant', content: data.text, toolCalls: data.tool_calls })
      chatHistory.value = data.messages || []
      // Re-probe after tool calls to update controls
      if (data.tool_calls?.length) setTimeout(() => runProbe(), 500)
    }
  } catch {
    chatMessages.value.push({ role: 'assistant', content: 'Failed to reach the server.' })
  }
  chatLoading.value = false
}

watch(chatMessages, () => {
  nextTick(() => { if (messagesRef.value) messagesRef.value.scrollTop = messagesRef.value.scrollHeight })
}, { deep: true })
</script>

<style scoped>
.device-page { min-height: 100vh; max-width: 720px; margin: 0 auto; padding: 0 24px 80px; }
.device-header { padding: 24px 0 16px; display: flex; align-items: center; justify-content: space-between; }
.back-btn { color: var(--color-text-secondary); text-decoration: none; font-size: 14px; font-weight: 500; transition: color 0.2s; }
.back-btn:hover { color: var(--color-text); }
.status-badge { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; padding: 3px 10px; border-radius: 99px; }
.status-badge.connected { color: #4ade80; background: rgba(34,197,94,0.1); border: 1px solid rgba(34,197,94,0.2); }
.status-badge.not-connected { color: #fbbf24; background: rgba(245,158,11,0.1); border: 1px solid rgba(245,158,11,0.2); }

/* Capabilities section */
.capabilities-section { padding: 28px 0; border-bottom: 1px solid var(--color-surface-border); }
.cap-card { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-lg); padding: 4px 0; }
.cap-item { display: flex; align-items: center; gap: 12px; padding: 12px 20px; }
.cap-item + .cap-item { border-top: 1px solid var(--color-surface-border); }
.cap-icon { font-size: 18px; width: 28px; text-align: center; flex-shrink: 0; color: var(--color-text-secondary); }
.cap-info { flex: 1; min-width: 0; }
.cap-name { display: block; font-size: 14px; font-weight: 500; color: var(--color-text); }
.cap-desc { display: block; font-size: 12px; color: var(--color-text-tertiary); margin-top: 1px; }
.cap-status { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; padding: 2px 8px; border-radius: 99px; white-space: nowrap; }
.cap-status.active { color: #4ade80; background: rgba(34,197,94,0.1); }
.cap-status.inactive { color: #fbbf24; background: rgba(245,158,11,0.1); }
.cap-api { background: rgba(99,102,241,0.03); }
.endpoints-details { margin-top: 12px; }
.endpoints-summary { font-size: 13px; font-weight: 600; color: var(--color-text-secondary); cursor: pointer; padding: 8px 0; }
.endpoints-summary:hover { color: var(--color-text); }
.doc-link { display: inline-block; margin-top: 12px; font-size: 13px; color: var(--color-accent); text-decoration: none; }
.doc-link:hover { text-decoration: underline; }

.loading-state { text-align: center; padding: 80px 0; color: var(--color-text-secondary); display: flex; flex-direction: column; align-items: center; gap: 16px; }
.loader { width: 32px; height: 32px; border: 3px solid var(--color-surface-border); border-top-color: var(--color-accent); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.device-hero { display: flex; align-items: flex-start; gap: 16px; padding: 24px 0 32px; border-bottom: 1px solid var(--color-surface-border); }
.hero-info { display: flex; flex-direction: column; gap: 4px; }
.hero-name { font-size: 24px; font-weight: 700; color: var(--color-text); letter-spacing: -0.02em; }
.hero-meta { display: flex; align-items: center; gap: 6px; }
.hero-ip { font-family: var(--font-mono); font-size: 13px; color: var(--color-text-secondary); background: rgba(255,255,255,0.05); padding: 2px 6px; border-radius: var(--radius-sm); }
.hero-mfr { font-size: 13px; color: var(--color-text-secondary); }
.hero-detail { font-size: 13px; color: var(--color-text-tertiary); }
.hero-integration { display: flex; align-items: center; gap: 8px; margin-top: 4px; }
.integration-pill { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; padding: 2px 8px; border-radius: 99px; }
.int-kasa { color: #4ade80; background: rgba(34,197,94,0.15); }
.int-hue { color: #60a5fa; background: rgba(0,101,211,0.2); }
.int-cast { color: #60a5fa; background: rgba(59,130,246,0.15); }
.int-generic { color: var(--color-text-secondary); background: var(--color-surface); }
.unreachable-badge { font-size: 11px; color: #ef4444; }
.meta-sep { color: var(--color-text-tertiary); opacity: 0.5; }

/* Setup */
.setup-section { padding: 28px 0; }
.setup-card { background: rgba(245,158,11,0.05); border: 1px solid rgba(245,158,11,0.2); border-radius: var(--radius-lg); padding: 24px; text-align: center; }
.setup-title { font-size: 16px; font-weight: 600; color: var(--color-text); margin-bottom: 8px; }
.setup-desc { font-size: 14px; color: var(--color-text-secondary); margin-bottom: 16px; line-height: 1.5; }
.setup-btn { background: var(--color-accent); color: white; border: none; border-radius: var(--radius-md); padding: 10px 24px; font-size: 14px; font-weight: 600; font-family: var(--font-sans); cursor: pointer; transition: background 0.2s; }
.setup-btn:hover:not(:disabled) { background: var(--color-accent-hover); }
.setup-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.setup-error { color: #ef4444; font-size: 13px; margin-top: 12px; }
.setup-success { color: #4ade80; font-size: 13px; margin-top: 12px; }

/* Sections */
.section-title { font-size: 14px; font-weight: 600; color: var(--color-text-secondary); text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 12px; display: flex; align-items: center; gap: 8px; }
.section-count { font-size: 12px; font-weight: 500; color: var(--color-text-tertiary); text-transform: none; letter-spacing: normal; }
.section-hint { font-size: 13px; color: var(--color-text-tertiary); margin-top: -4px; }

/* Controls */
.controls-section { padding: 28px 0; border-bottom: 1px solid var(--color-surface-border); }
.controls-card { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-lg); padding: 4px 0; }
.control-row { display: flex; align-items: center; justify-content: space-between; padding: 14px 20px; }
.control-label { font-size: 15px; font-weight: 500; color: var(--color-text); }
.control-right { display: flex; align-items: center; gap: 12px; }
.state-text { font-size: 13px; font-weight: 600; color: var(--color-text-tertiary); text-transform: uppercase; letter-spacing: 0.04em; }
.state-text.on { color: var(--color-success); }
.toggle-track { width: 52px; height: 28px; border-radius: 14px; background: var(--color-surface-border); position: relative; cursor: pointer; transition: background 0.2s; }
.toggle-track.on { background: var(--color-accent); }
.toggle-thumb { width: 24px; height: 24px; border-radius: 50%; background: white; position: absolute; top: 2px; left: 2px; transition: transform 0.2s; box-shadow: 0 1px 3px rgba(0,0,0,0.3); }
.toggle-track.on .toggle-thumb { transform: translateX(24px); }
.slider-row { padding: 0 20px 16px; }
.brightness-slider { width: 100%; -webkit-appearance: none; appearance: none; height: 6px; border-radius: 3px; background: var(--color-surface-border); outline: none; cursor: pointer; }
.brightness-slider::-webkit-slider-thumb { -webkit-appearance: none; width: 20px; height: 20px; border-radius: 50%; background: var(--color-accent); cursor: pointer; box-shadow: 0 1px 4px rgba(0,0,0,0.4); }
.brightness-value { font-size: 14px; font-weight: 600; color: var(--color-accent); }
.speed-buttons { display: flex; gap: 6px; }
.speed-btn { width: 48px; height: 36px; border-radius: var(--radius-sm); border: 1px solid var(--color-surface-border); background: var(--color-surface); color: var(--color-text-secondary); font-size: 12px; font-weight: 600; font-family: var(--font-sans); cursor: pointer; transition: all 0.2s; }
.speed-btn.active { background: var(--color-accent); border-color: var(--color-accent); color: white; }
.speed-btn:hover:not(.active) { border-color: var(--color-accent); color: var(--color-text); }

/* Info */
.info-section { padding: 28px 0; border-bottom: 1px solid var(--color-surface-border); }
.info-card { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-lg); padding: 4px 0; }
.info-row { display: flex; align-items: center; justify-content: space-between; padding: 12px 20px; }
.info-label { font-size: 13px; color: var(--color-text-secondary); font-weight: 500; }
.info-value { font-size: 13px; color: var(--color-text); }
.protocol-pills, .cap-pills { display: flex; gap: 4px; flex-wrap: wrap; }
.pill { font-size: 10px; font-weight: 600; padding: 2px 8px; border-radius: 99px; text-transform: uppercase; letter-spacing: 0.04em; background: rgba(255,255,255,0.07); color: var(--color-text-secondary); }
.pill-kasa { background: rgba(34,197,94,0.15); color: #4ade80; }
.pill-cast { background: rgba(59,130,246,0.15); color: #60a5fa; }
.pill-mdns { background: rgba(168,85,247,0.15); color: #c084fc; }
.pill-airplay { background: rgba(229,231,235,0.12); color: #e5e7eb; }
.pill-thread { background: rgba(245,158,11,0.15); color: #fbbf24; }
.cap-pill { font-size: 10px; font-weight: 500; padding: 2px 8px; border-radius: 99px; background: rgba(99,102,241,0.1); color: var(--color-accent); border: 1px solid rgba(99,102,241,0.2); }

/* Chat */
.chat-section { padding: 28px 0; }
.chat-card { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-lg); overflow: hidden; display: flex; flex-direction: column; }
.chat-suggestions { display: flex; flex-wrap: wrap; gap: 8px; padding: 24px 20px; }
.suggestion-chip { background: var(--color-bg); border: 1px solid var(--color-surface-border); color: var(--color-text-secondary); font-size: 13px; font-family: var(--font-sans); padding: 6px 14px; border-radius: 99px; cursor: pointer; transition: all 0.2s; }
.suggestion-chip:hover { border-color: var(--color-accent); color: var(--color-text); }
.chat-messages { max-height: 400px; overflow-y: auto; padding: 16px 20px; display: flex; flex-direction: column; gap: 10px; }
.chat-msg { display: flex; flex-direction: column; gap: 4px; }
.chat-msg.user { align-items: flex-end; }
.chat-msg.assistant { align-items: flex-start; }
.msg-bubble { max-width: 85%; padding: 10px 14px; font-size: 14px; line-height: 1.5; white-space: pre-wrap; word-break: break-word; }
.chat-msg.user .msg-bubble { background: var(--color-accent); color: #fff; border-radius: 16px 16px 4px 16px; }
.chat-msg.assistant .msg-bubble { background: var(--color-bg); color: var(--color-text); border-radius: 16px 16px 16px 4px; }
.msg-tools { display: flex; gap: 4px; flex-wrap: wrap; }
.tool-badge { font-size: 10px; color: var(--color-text-tertiary); background: rgba(255,255,255,0.05); padding: 2px 6px; border-radius: var(--radius-sm); }
.msg-bubble.typing { display: flex; gap: 4px; padding: 12px 18px; }
.dot { width: 6px; height: 6px; border-radius: 50%; background: var(--color-text-tertiary); animation: pulse-dot 1.2s ease-in-out infinite; }
.dot:nth-child(2) { animation-delay: 0.2s; }
.dot:nth-child(3) { animation-delay: 0.4s; }
@keyframes pulse-dot { 0%, 80%, 100% { opacity: 0.3; transform: scale(0.8); } 40% { opacity: 1; transform: scale(1); } }
.chat-input-row { display: flex; gap: 8px; padding: 12px 16px; border-top: 1px solid var(--color-surface-border); }
.chat-input { flex: 1; background: var(--color-bg); border: 1px solid var(--color-surface-border); border-radius: var(--radius-md); padding: 10px 14px; font-size: 14px; font-family: var(--font-sans); color: var(--color-text); outline: none; transition: border-color 0.2s; }
.chat-input:focus { border-color: var(--color-accent); }
.chat-input::placeholder { color: var(--color-text-tertiary); }
.chat-send { background: var(--color-accent); border: none; color: #fff; font-size: 18px; width: 42px; border-radius: var(--radius-md); cursor: pointer; transition: background 0.2s; }
.chat-send:hover:not(:disabled) { background: var(--color-accent-hover); }
.chat-send:disabled { opacity: 0.4; cursor: not-allowed; }

/* Discovery */
.discovery-section { padding: 28px 0; border-bottom: 1px solid var(--color-surface-border); }
.discovery-card { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-lg); padding: 4px 0; }
.fingerprint-row { display: flex; align-items: center; gap: 12px; padding: 12px 20px; border-bottom: 1px solid var(--color-surface-border); }
.fingerprint-row:last-child { border-bottom: none; }
.fp-port { font-family: var(--font-mono); font-size: 13px; font-weight: 600; color: var(--color-accent); min-width: 72px; }
.fp-details { flex: 1; display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.fp-title { font-size: 13px; color: var(--color-text); font-weight: 500; }
.fp-server { font-size: 12px; color: var(--color-text-tertiary); font-family: var(--font-mono); }
.fp-login-badge { font-size: 10px; font-weight: 600; padding: 2px 8px; border-radius: 99px; background: rgba(245,158,11,0.15); color: #fbbf24; border: 1px solid rgba(245,158,11,0.25); text-transform: uppercase; letter-spacing: 0.04em; }
.fp-link { font-size: 13px; color: var(--color-accent); text-decoration: none; font-weight: 500; white-space: nowrap; flex-shrink: 0; transition: opacity 0.2s; }
.fp-link:hover { opacity: 0.75; }

/* Auth */
.auth-section { padding: 28px 0; border-bottom: 1px solid var(--color-surface-border); }
.auth-card { background: rgba(245,158,11,0.05); border: 1px solid rgba(245,158,11,0.2); border-radius: var(--radius-lg); padding: 28px 24px; text-align: center; }
.auth-title { font-size: 16px; font-weight: 600; color: var(--color-text); margin-bottom: 8px; }
.auth-desc { font-size: 14px; color: var(--color-text-secondary); margin-bottom: 20px; line-height: 1.5; }
.auth-form { display: flex; flex-direction: column; gap: 10px; max-width: 320px; margin: 0 auto; }
.auth-input { background: var(--color-bg); border: 1px solid var(--color-surface-border); border-radius: var(--radius-md); padding: 10px 14px; font-size: 14px; font-family: var(--font-sans); color: var(--color-text); outline: none; text-align: left; transition: border-color 0.2s; }
.auth-input:focus { border-color: rgba(245,158,11,0.5); }
.auth-input::placeholder { color: var(--color-text-tertiary); }
.auth-btn { background: var(--color-accent); color: white; border: none; border-radius: var(--radius-md); padding: 10px 24px; font-size: 14px; font-weight: 600; font-family: var(--font-sans); cursor: pointer; transition: background 0.2s; margin-top: 4px; }
.auth-btn:hover:not(:disabled) { background: var(--color-accent-hover); }
.auth-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.auth-error { color: #ef4444; font-size: 13px; margin-top: 12px; }

/* Offline */
.offline-section { padding: 28px 0; border-bottom: 1px solid var(--color-surface-border); }
.offline-card { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-lg); padding: 40px 24px; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 12px; }
.offline-icon { width: 40px; height: 40px; color: #ef4444; opacity: 0.6; }
.offline-card h3 { font-size: 16px; font-weight: 600; color: var(--color-text); margin: 0; }
.offline-card p { font-size: 14px; color: var(--color-text-secondary); line-height: 1.5; max-width: 420px; margin: 0; }
.hint-suggestions { text-align: left; max-width: 420px; margin: 8px auto 0; }
.hint-item { font-size: 13px; color: var(--color-text-tertiary); line-height: 1.6; margin: 0; }
.retry-btn { background: transparent; color: var(--color-text-secondary); border: 1px solid var(--color-surface-border); border-radius: var(--radius-md); padding: 8px 20px; font-size: 14px; font-weight: 500; font-family: var(--font-sans); cursor: pointer; transition: all 0.2s; margin-top: 4px; }
.retry-btn:hover { border-color: var(--color-accent); color: var(--color-text); }

/* Auth form */
.auth-form { display: flex; gap: 8px; justify-content: center; flex-wrap: wrap; }
.auth-input { background: var(--color-bg); border: 1px solid var(--color-surface-border); border-radius: var(--radius-md); padding: 10px 14px; font-size: 14px; font-family: var(--font-sans); color: var(--color-text); outline: none; min-width: 200px; }
.auth-input:focus { border-color: var(--color-accent); }

/* Patterns grid */
.patterns-grid { display: flex; flex-wrap: wrap; gap: 6px; }
.pattern-chip { font-size: 12px; padding: 4px 12px; border-radius: 99px; background: var(--color-surface); border: 1px solid var(--color-surface-border); color: var(--color-text-secondary); }
.pattern-btn { cursor: pointer; font-family: var(--font-sans); transition: all 0.2s; }
.pattern-btn:hover { border-color: var(--color-accent); color: var(--color-text); }
.pattern-btn.active { background: var(--color-accent); border-color: var(--color-accent); color: white; }

/* Fingerprints */
.fp-left { display: flex; align-items: center; gap: 8px; flex: 1; min-width: 0; }
.fp-port { font-family: var(--font-mono); font-size: 13px; color: var(--color-accent); }
.fp-title { font-size: 13px; color: var(--color-text-secondary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.login-badge { font-size: 10px; font-weight: 600; color: #fbbf24; background: rgba(245,158,11,0.1); padding: 2px 6px; border-radius: 99px; text-transform: uppercase; letter-spacing: 0.04em; }
.fp-link { font-size: 13px; color: var(--color-accent); text-decoration: none; white-space: nowrap; }
.fp-link:hover { text-decoration: underline; }

/* Solar data grid */
.solar-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 12px; }
.solar-card { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-md); padding: 16px; text-align: center; }
.solar-label { display: block; font-size: 12px; color: var(--color-text-secondary); margin-bottom: 4px; text-transform: uppercase; letter-spacing: 0.04em; }
.solar-value { font-size: 20px; font-weight: 700; color: var(--color-text); }
.solar-green { color: #4ade80; }
.solar-red { color: #f87171; }

/* API section */
.api-proto { color: var(--color-accent); font-weight: 600; }
.api-desc { padding: 12px 20px; font-size: 13px; color: var(--color-text-secondary); line-height: 1.5; border-top: 1px solid var(--color-surface-border); }
.api-doc-link { display: block; padding: 10px 20px; font-size: 13px; color: var(--color-accent); text-decoration: none; border-top: 1px solid var(--color-surface-border); }
.api-doc-link:hover { text-decoration: underline; }
.section-subtitle { font-size: 12px; font-weight: 600; color: var(--color-text-tertiary); text-transform: uppercase; letter-spacing: 0.06em; margin: 16px 0 8px; }
.endpoints-list { display: flex; flex-direction: column; gap: 2px; }
.endpoint-row { background: var(--color-surface); border: 1px solid var(--color-surface-border); border-radius: var(--radius-md); padding: 12px 16px; }
.ep-header { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; overflow-x: auto; }
.ep-method { font-size: 11px; font-weight: 700; color: var(--color-accent); background: rgba(99,102,241,0.1); padding: 2px 6px; border-radius: 4px; white-space: nowrap; }
.ep-path { font-size: 11px; color: var(--color-text-secondary); word-break: break-all; }
.ep-desc { font-size: 12px; color: var(--color-text-tertiary); margin: 0; }
.ep-example { display: block; font-size: 11px; color: var(--color-text-tertiary); background: rgba(0,0,0,0.3); padding: 6px 8px; border-radius: 4px; margin-top: 6px; word-break: break-all; }
</style>
