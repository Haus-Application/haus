<template>
  <div class="scan-progress">
    <p class="progress-label">Scanning your network...</p>
    <ul class="stage-list">
      <li
        v-for="stage in stages"
        :key="stage.stage"
        class="stage-item"
        :class="stage.status"
      >
        <span class="stage-icon">
          <!-- Pending: dim outline circle -->
          <svg v-if="stage.status === 'pending'" class="icon-pending" viewBox="0 0 20 20">
            <circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" stroke-width="1.5" />
          </svg>
          <!-- Running: spinning arc -->
          <svg v-else-if="stage.status === 'running'" class="icon-running" viewBox="0 0 20 20">
            <circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" stroke-width="1.5" stroke-dasharray="28 22" stroke-linecap="round" />
          </svg>
          <!-- Complete: checkmark -->
          <svg v-else class="icon-complete" viewBox="0 0 20 20">
            <circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" stroke-width="1.5" />
            <polyline points="6,10 9,13 14,7" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </span>
        <span class="stage-name">{{ stage.stage }}</span>
        <span class="stage-message">{{ stage.message }}</span>
        <span v-if="stage.count != null && stage.status === 'complete'" class="stage-count">
          {{ stage.count }}
        </span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
interface ScanStage {
  stage: string
  status: 'pending' | 'running' | 'complete'
  message: string
  count?: number
}

defineProps<{
  stages: ScanStage[]
}>()
</script>

<style scoped>
.scan-progress {
  width: 100%;
  max-width: 480px;
  margin: 0 auto 32px;
}

.progress-label {
  font-size: 14px;
  color: var(--color-text-secondary);
  text-align: center;
  margin-bottom: 20px;
  letter-spacing: 0.02em;
}

.stage-list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stage-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  transition: background 0.3s ease;
  font-size: 14px;
}

.stage-item.running {
  background: var(--color-surface);
}

.stage-icon {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-pending {
  width: 18px;
  height: 18px;
  color: var(--color-text-tertiary);
}

.icon-running {
  width: 18px;
  height: 18px;
  color: var(--color-accent);
  animation: spin 0.9s linear infinite;
}

.icon-complete {
  width: 18px;
  height: 18px;
  color: var(--color-success);
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.stage-name {
  font-weight: 500;
  color: var(--color-text);
  min-width: 160px;
}

.stage-item.pending .stage-name {
  color: var(--color-text-tertiary);
}

.stage-message {
  font-size: 13px;
  color: var(--color-text-secondary);
  flex: 1;
}

.stage-item.pending .stage-message {
  color: var(--color-text-tertiary);
}

.stage-count {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-success);
  background: rgba(34, 197, 94, 0.1);
  padding: 2px 8px;
  border-radius: 99px;
}
</style>
