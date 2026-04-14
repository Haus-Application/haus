<template>
  <div class="validation-page">
    <header class="val-header">
      <div class="header-brand">
        <NuxtLink to="/" class="back-link">← Haus</NuxtLink>
        <h1 class="title">KB Validation</h1>
        <p class="subtitle">Buster & GOB grading the 100-device knowledge base</p>
      </div>
      <div class="header-actions">
        <button class="btn-run" :disabled="running" @click="runAll">
          <svg v-if="running" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" class="spin-icon">
            <path d="M21 12a9 9 0 1 1-6.2-8.6" />
          </svg>
          {{ running ? `Validating ${progress.done}/${progress.total}...` : 'Run Validation' }}
        </button>
      </div>
    </header>

    <!-- Stats strip -->
    <section class="stats-strip">
      <div class="stat">
        <div class="stat-value">{{ devices.length }}</div>
        <div class="stat-label">Devices</div>
      </div>
      <div class="stat">
        <div class="stat-value">{{ avgScore !== null ? avgScore + '%' : '—' }}</div>
        <div class="stat-label">Avg Score</div>
      </div>
      <div class="stat passing">
        <div class="stat-value">{{ passingCount }}</div>
        <div class="stat-label">Passing ≥85%</div>
      </div>
      <div class="stat warning">
        <div class="stat-value">{{ warningCount }}</div>
        <div class="stat-label">Warning 70–84%</div>
      </div>
      <div class="stat failing">
        <div class="stat-value">{{ failingCount }}</div>
        <div class="stat-label">Failing &lt;70%</div>
      </div>
    </section>

    <!-- Filters -->
    <section class="filters">
      <label>
        <span>Category</span>
        <select v-model="filterCategory">
          <option value="">All</option>
          <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
        </select>
      </label>
      <label>
        <span>Status</span>
        <select v-model="filterStatus">
          <option value="">All</option>
          <option value="supported">Supported</option>
          <option value="planned">Planned</option>
          <option value="read_only">Read-only</option>
          <option value="detected_only">Detected only</option>
          <option value="not_feasible">Not feasible</option>
        </select>
      </label>
      <label>
        <span>Score</span>
        <select v-model="filterScore">
          <option value="">All</option>
          <option value="passing">Passing (≥85%)</option>
          <option value="warning">Warning (70–84%)</option>
          <option value="failing">Failing (&lt;70%)</option>
          <option value="not_run">Not run</option>
        </select>
      </label>
    </section>

    <!-- Grid -->
    <main class="grid" v-if="filtered.length > 0">
      <NuxtLink
        v-for="d in filtered"
        :key="d.slug"
        :to="`/validation/${d.slug}`"
        class="tile"
        :class="tileClass(d)"
      >
        <div class="tile-top">
          <span class="cat-icon" :title="d.category">{{ categoryIcon(d.category) }}</span>
          <span class="tile-score">{{ d.total_pct !== null && d.total_pct !== undefined ? d.total_pct + '%' : '—' }}</span>
        </div>
        <div class="tile-name">{{ d.name }}</div>
        <div class="tile-meta">
          <span class="tile-cat">{{ d.category }}</span>
          <span class="sep">·</span>
          <span class="tile-status">{{ d.integration_status }}</span>
        </div>
      </NuxtLink>
    </main>

    <div v-else class="empty-state">
      <p>No devices match your filters.</p>
    </div>

    <!-- Loading overlay -->
    <div v-if="running" class="progress-overlay" role="status" aria-live="polite">
      <div class="progress-card">
        <div class="progress-title">{{ progress.status === 'done' ? 'Finishing up...' : 'Validating…' }}</div>
        <div class="progress-meter">
          <div class="progress-bar" :style="{ width: progressPct + '%' }"></div>
        </div>
        <div class="progress-line">
          {{ progress.done }} / {{ progress.total }}
          <span v-if="progress.slug" class="progress-slug">· {{ progress.slug }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWebSocket } from '~/composables/useWebSocket'

type DeviceEntry = {
  slug: string
  name: string
  category: string
  integration_key: string
  integration_status: string
  total_pct: number | null
  ran: boolean
}

type Summary = {
  ran_at: string
  total_devices: number
  avg_score: number
  by_category: Record<string, number>
  failing: string[]
  passing: number
  warning: number
}

const devices = ref<DeviceEntry[]>([])
const summary = ref<Summary | null>(null)
const running = ref(false)
const progress = ref<{ job_id?: string; slug?: string; done: number; total: number; status: string }>({
  done: 0, total: 0, status: 'idle',
})

const filterCategory = ref('')
const filterStatus = ref('')
const filterScore = ref('')

const ws = useWebSocket()

const categories = computed(() => {
  const set = new Set<string>()
  devices.value.forEach((d) => d.category && set.add(d.category))
  return Array.from(set).sort()
})

const passingCount = computed(() => devices.value.filter((d) => (d.total_pct ?? -1) >= 85).length)
const warningCount = computed(() => devices.value.filter((d) => {
  const p = d.total_pct
  return p !== null && p !== undefined && p >= 70 && p < 85
}).length)
const failingCount = computed(() => devices.value.filter((d) => (d.total_pct ?? 100) < 70 && d.ran).length)
const avgScore = computed(() => summary.value?.avg_score ?? null)

const filtered = computed(() => devices.value.filter((d) => {
  if (filterCategory.value && d.category !== filterCategory.value) return false
  if (filterStatus.value && d.integration_status !== filterStatus.value) return false
  if (filterScore.value) {
    const p = d.total_pct
    if (filterScore.value === 'passing' && !(p !== null && p !== undefined && p >= 85)) return false
    if (filterScore.value === 'warning' && !(p !== null && p !== undefined && p >= 70 && p < 85)) return false
    if (filterScore.value === 'failing' && !(p !== null && p !== undefined && p < 70)) return false
    if (filterScore.value === 'not_run' && d.ran) return false
  }
  return true
}))

const progressPct = computed(() => {
  if (progress.value.total === 0) return 0
  return Math.min(100, Math.round((progress.value.done / progress.value.total) * 100))
})

function categoryIcon(cat: string): string {
  const map: Record<string, string> = {
    lighting: '💡',
    security: '🔒',
    climate: '🌡️',
    energy: '⚡',
    media: '🎵',
    smart_home: '🏠',
    network: '📡',
    compute: '💻',
  }
  return map[cat] || '•'
}

function tileClass(d: DeviceEntry) {
  if (!d.ran || d.total_pct === null || d.total_pct === undefined) return 'tile-idle'
  if (d.total_pct >= 85) return 'tile-passing'
  if (d.total_pct >= 70) return 'tile-warning'
  return 'tile-failing'
}

async function loadData() {
  try {
    const [devResp, sumResp] = await Promise.all([
      fetch('/api/validation/devices'),
      fetch('/api/validation/summary'),
    ])
    devices.value = await devResp.json()
    if (sumResp.ok) {
      summary.value = await sumResp.json()
    }
  } catch (err) {
    console.error('[validation] load failed:', err)
  }
}

async function runAll() {
  if (running.value) return
  try {
    const resp = await fetch('/api/validation/run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({}),
    })
    if (resp.status === 409) {
      running.value = true
      return
    }
    if (!resp.ok) {
      const body = await resp.json().catch(() => ({}))
      alert('Could not start validation: ' + (body.error || resp.statusText))
      return
    }
    const body = await resp.json()
    running.value = true
    progress.value = { job_id: body.job_id, done: 0, total: devices.value.length, status: 'running' }
  } catch (err) {
    console.error('[validation] run failed:', err)
    running.value = false
  }
}

function onProgress(payload: any) {
  if (!payload) return
  progress.value = {
    job_id: payload.job_id,
    slug: payload.slug,
    done: payload.done ?? progress.value.done,
    total: payload.total ?? progress.value.total,
    status: payload.status ?? 'running',
  }
  // Update tile live
  if (payload.slug && typeof payload.score === 'number') {
    const idx = devices.value.findIndex((d) => d.slug === payload.slug)
    if (idx >= 0) {
      devices.value[idx] = { ...devices.value[idx], total_pct: payload.score, ran: true }
    }
  }
  if (payload.status === 'done' && !payload.slug) {
    // final event
    running.value = false
    loadData()
  }
  if (payload.status === 'failed' && !payload.slug) {
    running.value = false
  }
}

onMounted(() => {
  loadData()
  ws.connect()
  ws.on('validation:progress', onProgress)
})

onUnmounted(() => {
  ws.off('validation:progress', onProgress)
})
</script>

<style scoped>
.validation-page {
  min-height: 100vh;
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.val-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 24px;
}

.back-link {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 14px;
  display: block;
  margin-bottom: 6px;
}

.back-link:hover { color: var(--color-text); }

.title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 14px;
  margin-top: 2px;
}

.btn-run {
  background: var(--color-accent);
  color: white;
  border: none;
  padding: 10px 18px;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: background 0.15s ease;
}
.btn-run:hover:not(:disabled) { background: var(--color-accent-hover); }
.btn-run:disabled { opacity: 0.7; cursor: not-allowed; }

.spin-icon { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.stats-strip {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.stat {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-md);
  padding: 16px;
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
}
.stat-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 4px;
}
.stat.passing .stat-value { color: var(--color-success); }
.stat.warning .stat-value { color: var(--color-warning); }
.stat.failing .stat-value { color: var(--color-error); }

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}
.filters label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--color-text-secondary);
}
.filters select {
  background: var(--color-surface);
  color: var(--color-text);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-sm);
  padding: 6px 10px;
  font-family: inherit;
  font-size: 13px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}

.tile {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-md);
  padding: 14px;
  text-decoration: none;
  color: inherit;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: transform 0.15s ease, border-color 0.15s ease;
  position: relative;
}
.tile:hover {
  transform: translateY(-2px);
  border-color: var(--color-accent);
}

.tile-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.cat-icon { font-size: 20px; }

.tile-score {
  font-weight: 700;
  font-size: 16px;
  color: var(--color-text);
}

.tile-name {
  font-weight: 600;
  font-size: 14px;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tile-meta {
  font-size: 12px;
  color: var(--color-text-tertiary);
  display: flex;
  gap: 6px;
}
.sep { opacity: 0.5; }

.tile-passing { background: rgba(34, 197, 94, 0.08); border-color: rgba(34, 197, 94, 0.25); }
.tile-passing .tile-score { color: var(--color-success); }

.tile-warning { background: rgba(245, 158, 11, 0.08); border-color: rgba(245, 158, 11, 0.25); }
.tile-warning .tile-score { color: var(--color-warning); }

.tile-failing { background: rgba(239, 68, 68, 0.08); border-color: rgba(239, 68, 68, 0.25); }
.tile-failing .tile-score { color: var(--color-error); }

.tile-idle .tile-score { color: var(--color-text-tertiary); }

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--color-text-secondary);
}

.progress-overlay {
  position: fixed;
  inset: 0;
  background: rgba(13, 13, 15, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
}
.progress-card {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-lg);
  padding: 24px;
  width: min(420px, 90%);
}
.progress-title {
  font-weight: 600;
  margin-bottom: 12px;
}
.progress-meter {
  background: var(--color-bg);
  border-radius: 999px;
  height: 8px;
  overflow: hidden;
  margin-bottom: 10px;
}
.progress-bar {
  background: var(--color-accent);
  height: 100%;
  transition: width 0.3s ease;
}
.progress-line {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-family: var(--font-mono);
}
.progress-slug { color: var(--color-accent); }

@media (max-width: 700px) {
  .stats-strip { grid-template-columns: repeat(2, 1fr); }
  .val-header { flex-direction: column; align-items: flex-start; gap: 12px; }
}
</style>
