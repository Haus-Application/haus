<template>
  <div class="dashboard">
    <header class="dash-header">
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

    <!-- Full-screen scanning state -->
    <div v-if="status === 'scanning'" class="scan-fullscreen">
      <div class="scan-fox">
        <img src="~/assets/fox.png" alt="Haus scanning" class="fox-avatar scanning" />
      </div>
      <ScanProgress :stages="stages" />
    </div>

    <!-- Normal content (hidden during scan) -->
    <template v-else>
      <div v-if="status === 'error'" class="error-banner">
        Could not connect to scan service. Is the server running?
      </div>

      <!-- Brand filter -->
      <div v-if="totalDevices > 0" class="brand-filters">
        <button class="brand-chip" :class="{ active: !activeBrand }" @click="activeBrand = null">All</button>
        <button v-for="brand in availableBrands" :key="brand" class="brand-chip" :class="{ active: activeBrand === brand }" @click="activeBrand = activeBrand === brand ? null : brand">{{ brand }}</button>
      </div>

      <main class="dash-content">
        <div v-if="totalDevices === 0" class="empty-state">
          <img src="~/assets/fox.png" alt="Haus" class="fox-avatar" />
          <p class="empty-title">Hi, I'm Haus</p>
          <p class="empty-sub">Scan your network and I'll find all your smart home devices</p>
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
  </div>
</template>

<script setup lang="ts">
// No tabs, no drama, just the dashboard. Marry me.
const { stages, status, startScan, devicesByCategory, totalDevices, devices } = useScan()

// Brand filtering
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
</style>
