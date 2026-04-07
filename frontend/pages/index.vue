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

    <div v-if="status === 'scanning'" class="scan-progress-inline">
      <ScanProgress :stages="stages" />
    </div>

    <div v-if="status === 'error'" class="error-banner">
      Could not connect to scan service. Is the server running?
    </div>

    <main class="dash-content">
      <div v-if="totalDevices === 0 && status !== 'scanning'" class="empty-state">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="empty-icon">
          <circle cx="12" cy="12" r="10"/>
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
          <line x1="2" y1="12" x2="22" y2="12"/>
        </svg>
        <p class="empty-title">No devices yet</p>
        <p class="empty-sub">Scan your network to find smart home devices</p>
      </div>

      <DeviceGrid
        v-else
        :devices-by-category="devicesByCategory"
        :kasa-state="kasaStateByIp"
        @toggle="handleKasaToggle"
        @brightness="handleKasaBrightness"
      />
    </main>
  </div>
</template>

<script setup lang="ts">
// No tabs, no drama, just the dashboard. Marry me.
const { stages, status, startScan, devicesByCategory, totalDevices } = useScan()
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

.scan-progress-inline {
  margin-bottom: 24px;
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

.dash-content {
  flex: 1;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 80px 20px;
  text-align: center;
}

.empty-icon {
  color: var(--color-text-tertiary);
  opacity: 0.5;
}

.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text);
}

.empty-sub {
  font-size: 14px;
  color: var(--color-text-secondary);
}
</style>
