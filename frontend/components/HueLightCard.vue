<template>
  <div class="hue-light-card" :class="{ 'is-on': light.on, 'is-unreachable': !light.reachable }">
    <div class="card-header">
      <div class="light-info">
        <span class="light-name">{{ light.name }}</span>
        <span v-if="light.room_name" class="room-name">{{ light.room_name }}</span>
      </div>
      <div class="header-controls">
        <span v-if="!light.reachable" class="unreachable-badge">Offline</span>
        <span v-if="light.on" class="brightness-label">{{ Math.round(light.brightness) }}%</span>
        <button
          class="toggle"
          :class="{ 'toggle-on': light.on }"
          :disabled="!light.reachable"
          :aria-label="light.on ? 'Turn off' : 'Turn on'"
          :aria-pressed="light.on"
          @click="handleToggle"
        >
          <span class="toggle-thumb" />
        </button>
      </div>
    </div>

    <Transition name="slide-down">
      <div v-if="light.on && light.reachable" class="slider-row">
        <input
          type="range"
          min="1"
          max="100"
          :value="light.brightness"
          class="brightness-slider"
          @input="handleBrightnessInput"
        />
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import type { HueLight } from '~/composables/useHue'

const props = defineProps<{
  light: HueLight
}>()

const { toggleLight, setLightBrightness } = useHue()

async function handleToggle() {
  if (!props.light.reachable) return
  await toggleLight(props.light.id, !props.light.on)
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null

function handleBrightnessInput(e: Event) {
  const value = parseInt((e.target as HTMLInputElement).value, 10)
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    setLightBrightness(props.light.id, value)
  }, 300)
}
</script>

<style scoped>
.hue-light-card {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-left: 3px solid transparent;
  border-radius: var(--radius-lg);
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: background 0.2s ease, border-color 0.2s ease, border-left-color 0.2s ease;
  animation: card-appear 0.35s ease forwards;
}

.hue-light-card:hover {
  background: var(--color-surface-hover);
  border-left-color: #0065D3;
}

.hue-light-card.is-unreachable {
  opacity: 0.5;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.light-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.light-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.room-name {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.unreachable-badge {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.1);
  padding: 2px 8px;
  border-radius: 99px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.brightness-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
  min-width: 30px;
  text-align: right;
}

/* Toggle */
.toggle {
  flex-shrink: 0;
  width: 40px;
  height: 24px;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  position: relative;
  background: var(--color-surface-border);
  transition: background 0.2s ease;
  outline: none;
}

.toggle:disabled {
  cursor: not-allowed;
  opacity: 0.4;
}

.toggle:focus-visible {
  box-shadow: 0 0 0 2px #0065D3;
}

.toggle-on {
  background: #0065D3;
}

.toggle-thumb {
  position: absolute;
  top: 3px;
  left: 3px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #fff;
  transition: transform 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.4);
}

.toggle-on .toggle-thumb {
  transform: translateX(16px);
}

/* Slider */
.slider-row {
  padding: 0 2px;
}

.brightness-slider {
  width: 100%;
  -webkit-appearance: none;
  appearance: none;
  height: 3px;
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
  background: #0065D3;
  cursor: pointer;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
  transition: transform 0.15s ease;
}

.brightness-slider::-webkit-slider-thumb:hover {
  transform: scale(1.2);
}

.brightness-slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #0065D3;
  cursor: pointer;
  border: none;
}

/* Slide-down transition */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: opacity 0.2s ease, max-height 0.2s ease;
  overflow: hidden;
  max-height: 32px;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  max-height: 0;
}

@keyframes card-appear {
  from { opacity: 0; transform: translateY(6px); }
  to   { opacity: 1; transform: translateY(0); }
}
</style>
