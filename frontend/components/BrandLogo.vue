<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="brand-logo"
    aria-hidden="true"
  >
    <!-- Philips / Signify / Hue -->
    <template v-if="normalized === 'philips' || normalized === 'signify' || normalized === 'hue'">
      <!-- Bridge arch shape -->
      <path d="M4 22 Q16 6 28 22" stroke="#0065D3" stroke-width="3" stroke-linecap="round" fill="none"/>
      <rect x="10" y="22" width="4" height="5" rx="1" fill="#0065D3"/>
      <rect x="18" y="22" width="4" height="5" rx="1" fill="#0065D3"/>
      <circle cx="16" cy="13" r="3" fill="#0065D3"/>
    </template>

    <!-- TP-Link -->
    <template v-else-if="normalized === 'tp-link' || normalized === 'tplink'">
      <!-- Swoosh / link chain -->
      <path d="M8 16 C8 10 14 8 18 12 L22 8" stroke="#00A4A6" stroke-width="2.5" stroke-linecap="round" fill="none"/>
      <path d="M24 16 C24 22 18 24 14 20 L10 24" stroke="#00A4A6" stroke-width="2.5" stroke-linecap="round" fill="none"/>
      <circle cx="22" cy="8" r="3" fill="#00A4A6"/>
      <circle cx="10" cy="24" r="3" fill="#00A4A6"/>
    </template>

    <!-- NVIDIA -->
    <template v-else-if="normalized === 'nvidia'">
      <!-- Eye/claw mark -->
      <path d="M4 20 L4 10 L14 4 L14 10 C18 10 22 13 22 18 C22 22 19 26 14 26 L4 26 Z" fill="#76B900"/>
      <path d="M14 12 C17 12 20 14.5 20 18 C20 21 17.5 23 14 23 L8 23 L8 12 Z" fill="#1a1a1a"/>
    </template>

    <!-- LG -->
    <template v-else-if="normalized === 'lg'">
      <!-- LG circle face mark -->
      <circle cx="16" cy="16" r="12" stroke="#A50034" stroke-width="2.5" fill="none"/>
      <path d="M10 11 L10 21 L18 21" stroke="#A50034" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
      <path d="M18 16 L22 16 L22 21" stroke="#A50034" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
    </template>

    <!-- Google -->
    <template v-else-if="normalized === 'google'">
      <!-- Simplified 4-color G mark -->
      <path d="M28 16 L16 16 L16 20 L23 20 C22 23 19 25 16 25 C11 25 7 21 7 16 C7 11 11 7 16 7 C18.5 7 20.5 8 22 9.5 L25 6.5 C22.5 4.5 19.5 3 16 3 C8.5 3 3 9 3 16 C3 23 8.5 29 16 29 C23 29 29 24 29 16 Z" fill="#4285F4"/>
      <path d="M3 16 C3 12 4.5 8.5 7 6.5 L11 10.5 C9.5 12 9 14 9 16 Z" fill="#34A853"/>
      <path d="M16 29 C13 29 10 28 8 26.5 L12 22.5 C13.5 23.5 14.5 24 16 24 Z" fill="#FBBC05"/>
      <path d="M25 6.5 L22 9.5 C20.5 8 18.5 7 16 7 L16 3 C19.5 3 22.5 4.5 25 6.5 Z" fill="#EA4335"/>
    </template>

    <!-- Apple -->
    <template v-else-if="normalized === 'apple'">
      <!-- Apple silhouette -->
      <path d="M21 10 C21 10 22 7 20 5 C18 6 17 8 18 10 Z" fill="#E8E8EA"/>
      <path d="M16 11 C13 11 10 13 10 17 C10 21 12 27 15 27 C16.5 27 17 26 18.5 26 C20 26 20.5 27 22 27 C25 27 27 21 27 17 C27 14 25 11 22 11 C20.5 11 19 12 18.5 12 C18 12 17 11 16 11 Z" fill="#E8E8EA"/>
    </template>

    <!-- Yamaha -->
    <template v-else-if="normalized === 'yamaha'">
      <!-- Tuning fork mark — three-pronged -->
      <circle cx="16" cy="16" r="12" stroke="#7C3AED" stroke-width="2" fill="none"/>
      <line x1="16" y1="8" x2="16" y2="20" stroke="#7C3AED" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="16" y1="14" x2="9" y2="23" stroke="#7C3AED" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="16" y1="14" x2="23" y2="23" stroke="#7C3AED" stroke-width="2.5" stroke-linecap="round"/>
    </template>

    <!-- Brilliant -->
    <template v-else-if="normalized === 'brilliant'">
      <!-- B in a rounded square/diamond -->
      <rect x="4" y="4" width="24" height="24" rx="6" fill="#2563EB"/>
      <path d="M11 9 L11 23 L18 23 C21 23 23 21.5 23 19 C23 17.5 22 16.5 20.5 16 C21.5 15.5 22 14.5 22 13 C22 10.5 20 9 17 9 Z" fill="white"/>
      <path d="M14 12 L17 12 C18.5 12 19 12.8 19 14 C19 15.2 18.5 16 17 16 L14 16 Z" fill="#2563EB"/>
      <path d="M14 17.5 L17.5 17.5 C19 17.5 20 18.5 20 19.8 C20 21 19 22 17.5 22 L14 22 Z" fill="#2563EB"/>
    </template>

    <!-- SunPower -->
    <template v-else-if="normalized === 'sunpower'">
      <!-- Sun icon -->
      <circle cx="16" cy="16" r="6" fill="#F59E0B"/>
      <line x1="16" y1="4" x2="16" y2="8" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="16" y1="24" x2="16" y2="28" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="4" y1="16" x2="8" y2="16" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="24" y1="16" x2="28" y2="16" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="7.5" y1="7.5" x2="10.5" y2="10.5" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="21.5" y1="21.5" x2="24.5" y2="24.5" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="24.5" y1="7.5" x2="21.5" y2="10.5" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="10.5" y1="21.5" x2="7.5" y2="24.5" stroke="#F59E0B" stroke-width="2.5" stroke-linecap="round"/>
    </template>

    <!-- Enphase -->
    <template v-else-if="normalized === 'enphase'">
      <!-- Lightning bolt energy icon -->
      <polygon points="18,3 8,18 15,18 14,29 24,14 17,14" fill="#E65100"/>
    </template>

    <!-- JellyFish -->
    <template v-else-if="normalized === 'jellyfish'">
      <!-- Jellyfish: dome + trailing tentacles -->
      <path d="M8 18 Q8 8 16 8 Q24 8 24 18 Z" fill="#06B6D4" opacity="0.9"/>
      <path d="M10 19 Q10 25 12 28" stroke="#06B6D4" stroke-width="2" stroke-linecap="round" fill="none"/>
      <path d="M14 19 Q13 25 14 28" stroke="#06B6D4" stroke-width="2" stroke-linecap="round" fill="none"/>
      <path d="M18 19 Q19 25 18 28" stroke="#06B6D4" stroke-width="2" stroke-linecap="round" fill="none"/>
      <path d="M22 19 Q22 25 20 28" stroke="#06B6D4" stroke-width="2" stroke-linecap="round" fill="none"/>
    </template>

    <!-- Sonos -->
    <template v-else-if="normalized === 'sonos'">
      <!-- Sonos: stylized S-wave in a circle -->
      <circle cx="16" cy="16" r="12" fill="#1a1a1a" stroke="#E8E8EA" stroke-width="1.5"/>
      <path d="M10 20 C10 17 12 15 16 15 C20 15 22 13 22 10" stroke="#E8E8EA" stroke-width="3" stroke-linecap="round" fill="none"/>
      <circle cx="10" cy="20" r="2" fill="#E8E8EA"/>
      <circle cx="22" cy="10" r="2" fill="#E8E8EA"/>
    </template>

    <!-- Arris -->
    <template v-else-if="normalized === 'arris'">
      <!-- Router icon: box with antenna + wifi arcs -->
      <rect x="5" y="18" width="22" height="8" rx="2" fill="#6B7280"/>
      <circle cx="9" cy="22" r="1.5" fill="#1a1a1a"/>
      <line x1="14" y1="18" x2="14" y2="12" stroke="#6B7280" stroke-width="2" stroke-linecap="round"/>
      <line x1="18" y1="18" x2="18" y2="10" stroke="#6B7280" stroke-width="2" stroke-linecap="round"/>
      <path d="M9 12 Q16 6 23 12" stroke="#9CA3AF" stroke-width="1.5" stroke-linecap="round" fill="none"/>
      <path d="M11 15 Q16 11 21 15" stroke="#9CA3AF" stroke-width="1.5" stroke-linecap="round" fill="none"/>
    </template>

    <!-- Generic / Unknown -->
    <template v-else>
      <circle cx="16" cy="16" r="11" stroke="var(--color-text-tertiary)" stroke-width="2" fill="none"/>
      <circle cx="16" cy="12" r="3.5" fill="var(--color-text-tertiary)"/>
      <path d="M9 24 C9 20 12 18 16 18 C20 18 23 20 23 24" fill="var(--color-text-tertiary)"/>
    </template>
  </svg>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  brand: string
  size?: number
}>(), {
  size: 32,
})

const normalized = computed(() => {
  const raw = props.brand?.toLowerCase().trim() ?? ''
  // Map device_type values to brand keys for logo matching
  const aliasMap: Record<string, string> = {
    'hue_bridge': 'philips',
    'hue_sync': 'philips',
    'shield_tv': 'nvidia',
    'jellyfish': 'jellyfish',
    'solar_gateway': 'sunpower',
    'brilliant_switch': 'brilliant',
    'av_receiver': 'yamaha',
    'speaker': 'sonos',
    'thread_border_router': 'google',
    'tv': 'lg',
    'lg electronics': 'lg',
  }
  return aliasMap[raw] || raw
})
</script>

<style scoped>
.brand-logo {
  flex-shrink: 0;
  display: block;
}
</style>
