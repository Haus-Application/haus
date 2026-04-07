<template>
  <div class="scan-button-wrapper">
    <!-- Pulse rings — the radar that actually finds your stuff -->
    <div v-if="!isRescan" class="pulse-rings">
      <span class="ring ring-1" />
      <span class="ring ring-2" />
      <span class="ring ring-3" />
    </div>
    <button class="scan-btn" @click="$emit('click')">
      {{ isRescan ? 'Scan Again' : 'Scan My Network' }}
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isRescan: boolean
}>()

defineEmits<{
  click: []
}>()
</script>

<style scoped>
.scan-button-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

/* Concentric pulse rings expanding from center */
.pulse-rings {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.ring {
  position: absolute;
  border-radius: 50%;
  border: 1.5px solid var(--color-accent);
  opacity: 0;
  animation: radar-pulse 3s ease-out infinite;
}

.ring-1 { width: 160px; height: 160px; animation-delay: 0s; }
.ring-2 { width: 240px; height: 240px; animation-delay: 0.8s; }
.ring-3 { width: 320px; height: 320px; animation-delay: 1.6s; }

@keyframes radar-pulse {
  0% {
    transform: scale(0.6);
    opacity: 0.5;
  }
  100% {
    transform: scale(1);
    opacity: 0;
  }
}

.scan-btn {
  position: relative;
  z-index: 1;
  width: 220px;
  height: 56px;
  background: var(--color-accent);
  color: #fff;
  font-family: var(--font-sans);
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.01em;
  border: none;
  border-radius: var(--radius-xl);
  cursor: pointer;
  box-shadow:
    0 0 0 0 rgba(99, 102, 241, 0),
    0 4px 24px rgba(99, 102, 241, 0.35);
  transition:
    background 0.3s ease,
    box-shadow 0.3s ease,
    transform 0.15s ease;
}

.scan-btn:hover {
  background: var(--color-accent-hover);
  box-shadow:
    0 0 0 0 rgba(99, 102, 241, 0),
    0 6px 36px rgba(129, 140, 248, 0.5);
  transform: translateY(-1px);
}

.scan-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 12px rgba(99, 102, 241, 0.3);
}
</style>
