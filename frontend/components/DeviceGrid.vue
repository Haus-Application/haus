<template>
  <div class="device-grid-wrapper">
    <section
      v-for="(categoryDevices, category) in devicesByCategory"
      :key="category"
      class="category-section"
    >
      <div class="category-header">
        <span class="category-icon" v-html="CATEGORY_ICONS[category as string] ?? CATEGORY_ICONS.unknown" />
        <h2 class="category-name">{{ CATEGORY_NAMES[category as string] ?? category }}</h2>
        <span class="category-count">{{ categoryDevices.length }}</span>
      </div>
      <div class="device-grid">
        <DeviceCard
          v-for="device in categoryDevices"
          :key="device.ip"
          :device="device"
          :live-state="kasaState[device.ip] ?? null"
          @toggle="$emit('toggle', $event)"
          @brightness="$emit('brightness', $event)"
        />
      </div>
    </section>
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

defineProps<{
  devicesByCategory: Record<string, ScanDevice[]>
  kasaState?: Record<string, KasaLiveState>
}>()

defineEmits<{
  toggle: [payload: { ip: string; on: boolean }]
  brightness: [payload: { ip: string; brightness: number }]
}>()

const CATEGORY_ICONS: Record<string, string> = {
  lighting: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M9 21h6"/><path d="M12 3a6 6 0 0 1 6 6c0 2.2-1.2 4.1-3 5.2V18H9v-3.8A6 6 0 0 1 6 9a6 6 0 0 1 6-6z"/>
  </svg>`,

  media: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <rect x="2" y="4" width="20" height="14" rx="2"/><path d="M8 21h8"/><path d="M12 17v4"/>
    <polygon points="10,9 15,12 10,15" fill="currentColor" stroke="none"/>
  </svg>`,

  smart_home: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M3 12L12 3l9 9"/><path d="M9 21V12h6v9"/><path d="M5 10V21h14V10"/>
  </svg>`,

  energy: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <polygon points="13,2 3,14 12,14 11,22 21,10 12,10" stroke="none" fill="currentColor"/>
  </svg>`,

  network: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <circle cx="12" cy="12" r="10"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/><line x1="2" y1="12" x2="22" y2="12"/>
  </svg>`,

  compute: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <rect x="3" y="4" width="18" height="12" rx="2"/><line x1="8" y1="20" x2="16" y2="20"/><line x1="12" y1="16" x2="12" y2="20"/>
  </svg>`,

  unknown: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <circle cx="12" cy="12" r="10"/><path d="M9.1 9a3 3 0 0 1 5.8 1c0 2-3 3-3 3"/><circle cx="12" cy="17" r="0.5" fill="currentColor"/>
  </svg>`,
}

const CATEGORY_NAMES: Record<string, string> = {
  lighting: 'Lighting',
  media: 'Media & Entertainment',
  smart_home: 'Smart Home',
  energy: 'Energy',
  network: 'Network',
  compute: 'Computers & Mobile',
  unknown: 'Other Devices',
}
</script>

<style scoped>
.device-grid-wrapper {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.category-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.category-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.category-icon {
  display: flex;
  align-items: center;
  color: var(--color-text-secondary);
  line-height: 1;
  flex-shrink: 0;
}

.category-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  letter-spacing: 0.01em;
}

.category-count {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  padding: 1px 8px;
  border-radius: 99px;
}

.device-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}
</style>
