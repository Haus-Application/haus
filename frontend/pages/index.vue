<template>
  <div class="dashboard">
    <header class="dash-header" v-if="welcomeStep === 'done'">
      <div class="header-brand">
        <h1 class="title">Haus</h1>
        <p class="subtitle">Smart Home</p>
      </div>
      <div class="header-actions">
        <button
          class="btn-scan"
          :class="{ scanning: status === 'scanning' }"
          :disabled="status === 'scanning'"
          @click="startScan"
        >
          <svg v-if="status === 'scanning'" class="spin-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
            <path d="M21 12a9 9 0 1 1-6.2-8.6"/>
          </svg>
          {{ status === 'scanning' ? 'Scanning...' : status === 'complete' ? 'Scan Again' : 'Scan Network' }}
        </button>
      </div>
    </header>

    <!-- Step 1: First-visit welcome (auto-advances to scan) -->
    <div v-if="welcomeStep === 'welcome'" class="welcome-screen" role="status" aria-live="polite">
      <img src="~/assets/fox.png" alt="Haus" class="fox-avatar welcome-bounce" />
      <h1 class="welcome-title">Hi, I'm Haus</h1>
      <p class="welcome-sub">Give me a second — I'll find everyone on your network.</p>
      <div class="welcome-dots">
        <span /><span /><span />
      </div>
    </div>

    <!-- Step 2: Full-screen scanning state -->
    <div v-else-if="status === 'scanning'" class="scan-fullscreen">
      <div class="scan-fox">
        <img src="~/assets/fox.png" alt="Haus scanning" class="fox-avatar scanning" />
      </div>
      <ScanProgress :stages="stages" />
    </div>

    <!-- Step 3: Results -->
    <template v-else>
      <div v-if="status === 'error'" class="error-banner">
        Could not connect to scan service. Is the server running?
      </div>

      <!-- Smart suggestions: integration opportunities we detected -->
      <div v-if="suggestions.length > 0" class="suggestions" role="region" aria-label="Setup suggestions">
        <div class="suggestions-head">
          <span class="suggestions-title">
            I see {{ suggestions.length }} {{ suggestions.length === 1 ? 'thing' : 'things' }} you can set up
          </span>
        </div>
        <div class="suggestions-list">
          <div v-for="s in suggestions" :key="s.id" class="suggestion-card" :class="`sug-${s.id}`">
            <div class="sug-icon" aria-hidden="true">{{ s.icon }}</div>
            <div class="sug-body">
              <div class="sug-title">{{ s.title }}</div>
              <div class="sug-sub">{{ s.sub }}</div>
            </div>
            <div class="sug-actions">
              <button class="sug-btn primary" @click="s.action">{{ s.cta }}</button>
              <button class="sug-btn ghost" aria-label="Dismiss" @click="dismissSuggestion(s.id)">✕</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Brand filter -->
      <div v-if="totalDevices > 0" class="brand-filters">
        <button class="brand-chip" :class="{ active: !activeBrand }" @click="activeBrand = null">All</button>
        <button v-for="brand in availableBrands" :key="brand" class="brand-chip" :class="{ active: activeBrand === brand }" @click="activeBrand = activeBrand === brand ? null : brand">{{ brand }}</button>
      </div>

      <main class="dash-content">
        <div v-if="totalDevices === 0" class="empty-state">
          <img src="~/assets/fox.png" alt="Haus" class="fox-avatar" />
          <p class="empty-title">I couldn't find anyone yet</p>
          <p class="empty-sub">Make sure your devices are on the same network, then scan again.</p>
          <button class="btn-scan" @click="startScan">Scan Again</button>
        </div>

        <DeviceGrid
          v-if="totalDevices > 0"
          :devices-by-category="filteredDevicesByCategory"
          :kasa-state="kasaStateByIp"
          @toggle="handleKasaToggle"
          @brightness="handleKasaBrightness"
        />
      </main>
    </template>

    <!-- Hue pairing dialog (shown on demand from suggestion card) -->
    <HuePairDialog :show="showHueDialog" @close="onHuePairClose" />
  </div>
</template>

<script setup lang="ts">
// No tabs, no drama, just the dashboard. Marry me.
const { stages, status, startScan, devicesByCategory, totalDevices, devices } = useScan()

// -- First-visit welcome / auto-scan -----------------------------------------
// On first visit (flag not set + no devices yet) we show a 2s welcome screen,
// then auto-trigger a scan. After that first scan, we never auto-scan again —
// returning users go straight to their grid.
type WelcomeStep = 'pending' | 'welcome' | 'done'
const welcomeStep = ref<WelcomeStep>('pending')
const FIRST_SCAN_KEY = 'haus_first_scan_done'

// -- Integration status ------------------------------------------------------
const hueConnected = ref<boolean | null>(null)
const googleConnected = ref<boolean | null>(null)
const dismissed = ref<Set<string>>(new Set(JSON.parse(localStorage.getItem('haus_dismissed_suggestions') || '[]')))

async function refreshIntegrationStatus() {
  try {
    const [hue, google] = await Promise.all([
      fetch('/api/hue/status').then((r) => (r.ok ? r.json() : null)).catch(() => null),
      fetch('/api/google/status').then((r) => (r.ok ? r.json() : null)).catch(() => null),
    ])
    hueConnected.value = !!hue?.connected
    googleConnected.value = !!google?.connected
  } catch {
    // non-fatal
  }
}

onMounted(async () => {
  await refreshIntegrationStatus()

  // Wait a tick so useScan's own onMounted loadDevices can settle.
  await nextTick()
  const alreadyRan = localStorage.getItem(FIRST_SCAN_KEY) === '1'
  if (!alreadyRan && devices.value.length === 0) {
    welcomeStep.value = 'welcome'
    // Give the welcome copy a beat, then kick off the scan.
    setTimeout(() => {
      welcomeStep.value = 'done'
      startScan()
    }, 2200)
  } else {
    welcomeStep.value = 'done'
  }
})

// When a scan finishes, mark the first-run flag and re-check integrations
// (because the scan may have just discovered a Hue bridge or Nest devices).
watch(status, (s) => {
  if (s === 'complete') {
    localStorage.setItem(FIRST_SCAN_KEY, '1')
    refreshIntegrationStatus()
  }
})

// -- Brand filtering ---------------------------------------------------------
const activeBrand = ref<string | null>(null)

const BRAND_MAP: Record<string, string> = {
  'tp-link': 'TP-Link',
  'philips': 'Philips',
  'signify': 'Philips',
  'google': 'Google',
  'google/nest': 'Google',
  'nest labs': 'Google',
  'nvidia': 'NVIDIA',
  'lg electronics': 'LG',
  'brilliant': 'Brilliant',
  'sunpower': 'SunPower',
  'enphase': 'Enphase',
  'yamaha': 'Yamaha',
  'sonos': 'Sonos',
  'samsung': 'Samsung',
  'apple': 'Apple',
  'arris': 'Arris',
}

function getBrand(mfr: string): string {
  return BRAND_MAP[mfr.toLowerCase()] || mfr || ''
}

const availableBrands = computed(() => {
  const brands = new Set<string>()
  for (const d of devices.value) {
    const brand = getBrand(d.manufacturer || '')
    if (brand) brands.add(brand)
  }
  return [...brands].sort()
})

// -- Smart integration suggestions ------------------------------------------
const hasHueBridge = computed(() =>
  devices.value.some((d) => {
    const mfr = (d.manufacturer || '').toLowerCase()
    const type = (d.device_type || '').toLowerCase()
    return type === 'hue_bridge' || mfr.includes('philips') || mfr.includes('signify')
  })
)
const hasGoogleDevices = computed(() =>
  devices.value.some((d) => {
    const mfr = (d.manufacturer || '').toLowerCase()
    return mfr.includes('google') || mfr.includes('nest')
  })
)

const showHueDialog = ref(false)
function openHuePair() { showHueDialog.value = true }
function onHuePairClose() {
  showHueDialog.value = false
  refreshIntegrationStatus()
}
function startGoogleAuth() {
  window.location.href = '/api/google/auth'
}
function dismissSuggestion(id: string) {
  dismissed.value.add(id)
  dismissed.value = new Set(dismissed.value)
  localStorage.setItem('haus_dismissed_suggestions', JSON.stringify([...dismissed.value]))
}

type Suggestion = {
  id: string
  icon: string
  title: string
  sub: string
  cta: string
  action: () => void
}

const suggestions = computed<Suggestion[]>(() => {
  const out: Suggestion[] = []
  if (hasHueBridge.value && hueConnected.value === false && !dismissed.value.has('hue')) {
    out.push({
      id: 'hue',
      icon: '💡',
      title: 'Pair your Philips Hue Bridge',
      sub: "I found a Hue bridge — let's link it so I can control your lights.",
      cta: 'Pair now',
      action: openHuePair,
    })
  }
  if (hasGoogleDevices.value && googleConnected.value === false && !dismissed.value.has('google')) {
    out.push({
      id: 'google',
      icon: '🔗',
      title: 'Connect your Google account',
      sub: 'I see Google Nest devices. Sign in to control thermostats and see cameras.',
      cta: 'Sign in',
      action: startGoogleAuth,
    })
  }
  return out
})

const filteredDevicesByCategory = computed(() => {
  if (!activeBrand.value) return devicesByCategory.value
  const filtered: Record<string, any[]> = {}
  for (const [cat, devs] of Object.entries(devicesByCategory.value)) {
    const matched = devs.filter((d: any) => getBrand(d.manufacturer || '') === activeBrand.value)
    if (matched.length > 0) filtered[cat] = matched
  }
  return filtered
})
const { devices: kasaDevices, toggleDevice, setBrightness } = useKasa()

// Build a fast lookup: IP -> KasaDevice so DeviceCard can get live state
const kasaStateByIp = computed(() => {
  const map: Record<string, typeof kasaDevices.value[0]> = {}
  for (const d of kasaDevices.value) {
    map[d.ip] = d
  }
  return map
})

async function handleKasaToggle({ ip, on }: { ip: string; on: boolean }) {
  await toggleDevice(ip, on)
}

async function handleKasaBrightness({ ip, brightness }: { ip: string; brightness: number }) {
  await setBrightness(ip, brightness)
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 0 32px 80px;
  max-width: 1400px;
  margin: 0 auto;
}

.dash-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 40px 0 36px;
  gap: 16px;
}

.header-brand {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--color-text);
  line-height: 1;
}

.subtitle {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 400;
  letter-spacing: 0.02em;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn-scan {
  display: flex;
  align-items: center;
  gap: 7px;
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  padding: 9px 18px;
  font-size: 13px;
  font-weight: 600;
  font-family: var(--font-sans);
  letter-spacing: 0.01em;
  cursor: pointer;
  transition: background 0.2s ease, opacity 0.2s ease;
}

.btn-scan:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.btn-scan:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-scan.scanning {
  background: var(--color-surface);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-surface-border);
}

.spin-icon {
  animation: spin 0.9s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.scan-fullscreen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 32px;
  padding: 40px 20px;
}

.error-banner {
  font-size: 14px;
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.2);
  padding: 12px 20px;
  border-radius: var(--radius-md);
  margin-bottom: 24px;
}

.brand-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 24px;
}

.brand-chip {
  font-size: 12px;
  font-weight: 600;
  font-family: var(--font-sans);
  padding: 5px 14px;
  border-radius: 99px;
  border: 1px solid var(--color-surface-border);
  background: var(--color-surface);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.brand-chip:hover {
  border-color: var(--color-accent);
  color: var(--color-text);
}

.brand-chip.active {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: white;
}

.dash-content {
  flex: 1;
}

.fox-avatar {
  width: 280px;
  height: 280px;
  object-fit: cover;
  object-position: top;
  border-radius: 50%;
  filter: drop-shadow(0 4px 12px rgba(0,0,0,0.3));
}

.fox-avatar.scanning {
  animation: fox-pulse 2s ease-in-out infinite;
}

@keyframes fox-pulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.05); opacity: 0.9; }
}

.scan-fox {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.fox-status {
  font-size: 14px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 80px 20px;
  text-align: center;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text);
}

.empty-sub {
  font-size: 14px;
  color: var(--color-text-secondary);
  max-width: 300px;
}

/* -- First-visit welcome screen ---------------------------------------- */
.welcome-screen {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 18px;
  padding: 40px 20px;
  text-align: center;
}
.welcome-title {
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--color-text);
  margin-top: 8px;
}
.welcome-sub {
  font-size: 15px;
  color: var(--color-text-secondary);
  max-width: 380px;
}
.welcome-bounce {
  animation: fox-bounce 1.4s ease-in-out infinite;
}
@keyframes fox-bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-12px); }
}
.welcome-dots {
  display: flex;
  gap: 6px;
  margin-top: 6px;
}
.welcome-dots span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-accent);
  opacity: 0.4;
  animation: dot-pulse 1.2s ease-in-out infinite;
}
.welcome-dots span:nth-child(2) { animation-delay: 0.2s; }
.welcome-dots span:nth-child(3) { animation-delay: 0.4s; }
@keyframes dot-pulse {
  0%, 100% { opacity: 0.2; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1.1); }
}

/* -- Smart suggestions ---------------------------------------------------- */
.suggestions {
  margin: 0 0 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.suggestions-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2px 4px;
}
.suggestions-title {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
  letter-spacing: 0.02em;
}
.suggestions-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.suggestion-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-md);
  transition: border-color 0.15s ease, transform 0.15s ease;
}
.suggestion-card:hover {
  border-color: var(--color-accent);
  transform: translateY(-1px);
}
.sug-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  border-radius: var(--radius-sm);
  flex-shrink: 0;
}
.sug-body {
  flex: 1;
  min-width: 0;
}
.sug-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 2px;
}
.sug-sub {
  font-size: 12px;
  color: var(--color-text-secondary);
  line-height: 1.4;
}
.sug-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}
.sug-btn {
  border: none;
  cursor: pointer;
  font-family: inherit;
  font-size: 12px;
  font-weight: 600;
  padding: 7px 14px;
  border-radius: var(--radius-sm);
  transition: background 0.15s ease, opacity 0.15s ease;
}
.sug-btn.primary {
  background: var(--color-accent);
  color: #fff;
}
.sug-btn.primary:hover {
  background: var(--color-accent-hover);
}
.sug-btn.ghost {
  background: transparent;
  color: var(--color-text-tertiary);
  padding: 6px 8px;
}
.sug-btn.ghost:hover {
  background: var(--color-surface-hover);
  color: var(--color-text);
}

@media (max-width: 600px) {
  .suggestion-card {
    flex-wrap: wrap;
  }
  .sug-body { flex-basis: calc(100% - 56px); }
  .sug-actions { margin-left: auto; }
}
</style>
