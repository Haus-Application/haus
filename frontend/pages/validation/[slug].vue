<template>
  <div class="transcript-page">
    <header class="trans-header">
      <NuxtLink to="/validation" class="back-link">← Back to Grid</NuxtLink>
      <div v-if="report" class="header-row">
        <div>
          <h1 class="title">{{ report.name }}</h1>
          <div class="meta-row">
            <span class="pill">{{ report.category }}</span>
            <span class="pill">{{ report.integration_status }}</span>
            <span v-if="report.integration_key" class="pill pill-key">{{ report.integration_key }}</span>
            <span class="pill-time">ran {{ formatAgo(report.ran_at) }}</span>
          </div>
        </div>
        <div class="score-circle" :class="scoreClass(report.score.total_pct)">
          <div class="score-pct">{{ report.score.total_pct }}%</div>
          <div class="score-breakdown">
            <span>B {{ report.score.buster }}/{{ report.score.buster_max }}</span>
            <span>·</span>
            <span>G {{ report.score.gob }}/{{ report.score.gob_max }}</span>
          </div>
        </div>
      </div>
    </header>

    <div v-if="loading" class="state">Loading…</div>
    <div v-else-if="notRun" class="state">
      <p><strong>Not validated yet.</strong></p>
      <p class="state-sub">This device hasn't been graded. Kick off a run below.</p>
      <button class="btn-run" :disabled="running" @click="runOne">
        {{ running ? 'Running…' : 'Run Validation for this device' }}
      </button>
    </div>
    <div v-else-if="error" class="state">Failed to load: {{ error }}</div>

    <template v-else-if="report">
      <section v-if="report.gaps && report.gaps.length > 0" class="gaps">
        <div class="gaps-title">📋 Gaps in the documentation</div>
        <ul>
          <li v-for="(g, i) in report.gaps" :key="i">{{ g }}</li>
        </ul>
      </section>

      <div class="qa-grid">
        <section class="qa-col">
          <h2 class="col-title">🧑‍🔧 Buster — Technical</h2>
          <article v-for="(qa, i) in report.buster" :key="'b' + i" class="qa-card">
            <div class="qa-top">
              <div class="qa-q">Q{{ i + 1 }}: {{ qa.q }}</div>
              <div class="qa-score" :class="qaScoreClass(qa.score, qa.max)">
                {{ qa.score }}/{{ qa.max }}
              </div>
            </div>
            <blockquote class="qa-a">{{ qa.a }}</blockquote>
            <ul v-if="qa.notes && qa.notes.length > 0" class="qa-notes">
              <li v-for="(n, j) in qa.notes" :key="j" :class="noteClass(n)">{{ n }}</li>
            </ul>
          </article>
        </section>

        <section class="qa-col">
          <h2 class="col-title">🎭 GOB — UX &amp; UI</h2>
          <article v-for="(qa, i) in report.gob" :key="'g' + i" class="qa-card">
            <div class="qa-top">
              <div class="qa-q">Q{{ i + 1 }}: {{ qa.q }}</div>
              <div class="qa-score" :class="qaScoreClass(qa.score, qa.max)">
                {{ qa.score }}/{{ qa.max }}
              </div>
            </div>
            <blockquote class="qa-a">{{ qa.a }}</blockquote>
            <ul v-if="qa.notes && qa.notes.length > 0" class="qa-notes">
              <li v-for="(n, j) in qa.notes" :key="j" :class="noteClass(n)">{{ n }}</li>
            </ul>
          </article>
        </section>
      </div>

      <footer class="footer-actions">
        <a
          :href="`https://github.com/Haus-Application/haus/blob/main/docs/devices/${slug}.md`"
          target="_blank"
          rel="noopener"
          class="btn-ghost"
        >Edit KB file on GitHub →</a>
        <button class="btn-ghost" :disabled="running" @click="runOne">
          {{ running ? 'Running…' : 'Re-run this device' }}
        </button>
      </footer>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useWebSocket } from '~/composables/useWebSocket'

type QA = { q: string; a: string; score: number; max: number; notes: string[] }
type Report = {
  slug: string
  name: string
  category: string
  integration_key: string
  integration_status: string
  ran_at: string
  score: { buster: number; buster_max: number; gob: number; gob_max: number; total_pct: number }
  buster: QA[]
  gob: QA[]
  gaps: string[]
}

const route = useRoute()
const slug = computed(() => String(route.params.slug))
const report = ref<Report | null>(null)
const loading = ref(true)
const notRun = ref(false)
const error = ref('')
const running = ref(false)
const ws = useWebSocket()

async function loadReport() {
  loading.value = true
  notRun.value = false
  error.value = ''
  try {
    const resp = await fetch(`/api/validation/devices/${slug.value}`)
    if (resp.status === 404) {
      notRun.value = true
      loading.value = false
      return
    }
    if (!resp.ok) {
      error.value = resp.statusText
      loading.value = false
      return
    }
    const body = await resp.json()
    if (body.status === 'not_run') {
      notRun.value = true
    } else {
      report.value = body
    }
  } catch (err: any) {
    error.value = String(err)
  } finally {
    loading.value = false
  }
}

async function runOne() {
  if (running.value) return
  running.value = true
  try {
    const resp = await fetch('/api/validation/run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ only: slug.value }),
    })
    if (!resp.ok && resp.status !== 409) {
      error.value = 'Could not start run: ' + resp.statusText
      running.value = false
    }
  } catch (err: any) {
    error.value = String(err)
    running.value = false
  }
}

function onProgress(payload: any) {
  if (!payload || payload.slug !== slug.value) return
  if (payload.status === 'done' || payload.status === 'failed') {
    running.value = false
    if (payload.status === 'done') {
      loadReport()
    }
  }
}

function scoreClass(pct: number) {
  if (pct >= 85) return 'score-passing'
  if (pct >= 70) return 'score-warning'
  return 'score-failing'
}

function qaScoreClass(score: number, max: number) {
  if (max === 0) return 'qa-na'
  if (score === max) return 'qa-full'
  if (score * 2 >= max) return 'qa-partial'
  return 'qa-miss'
}

function noteClass(note: string) {
  if (note.startsWith('✗')) return 'note-miss'
  return 'note-hit'
}

function formatAgo(ts: string): string {
  try {
    const d = new Date(ts)
    const seconds = Math.floor((Date.now() - d.getTime()) / 1000)
    if (seconds < 60) return `${seconds}s ago`
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
    return `${Math.floor(seconds / 86400)}d ago`
  } catch {
    return ts
  }
}

onMounted(() => {
  loadReport()
  ws.connect()
  ws.on('validation:progress', onProgress)
})
onUnmounted(() => {
  ws.off('validation:progress', onProgress)
})
</script>

<style scoped>
.transcript-page {
  min-height: 100vh;
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.back-link {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 14px;
  display: inline-block;
  margin-bottom: 12px;
}
.back-link:hover { color: var(--color-text); }

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
  margin-bottom: 24px;
}

.title {
  font-size: 28px;
  font-weight: 700;
}

.meta-row {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  flex-wrap: wrap;
  align-items: center;
}

.pill {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  color: var(--color-text-secondary);
  border-radius: 999px;
  padding: 3px 10px;
  font-size: 12px;
  text-transform: capitalize;
}
.pill-key { color: var(--color-accent); }
.pill-time {
  color: var(--color-text-tertiary);
  font-size: 12px;
  margin-left: auto;
  padding-left: 8px;
}

.score-circle {
  min-width: 140px;
  padding: 16px 18px;
  border-radius: var(--radius-lg);
  text-align: center;
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
}
.score-circle.score-passing { background: rgba(34, 197, 94, 0.1); border-color: rgba(34, 197, 94, 0.3); }
.score-circle.score-warning { background: rgba(245, 158, 11, 0.1); border-color: rgba(245, 158, 11, 0.3); }
.score-circle.score-failing { background: rgba(239, 68, 68, 0.1); border-color: rgba(239, 68, 68, 0.3); }
.score-pct { font-size: 32px; font-weight: 700; }
.score-breakdown {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 4px;
  display: flex;
  justify-content: center;
  gap: 6px;
}
.score-passing .score-pct { color: var(--color-success); }
.score-warning .score-pct { color: var(--color-warning); }
.score-failing .score-pct { color: var(--color-error); }

.state {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-md);
  padding: 32px;
  text-align: center;
  color: var(--color-text-secondary);
}
.state-sub {
  color: var(--color-text-tertiary);
  margin-top: 4px;
  margin-bottom: 16px;
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
}
.btn-run:disabled { opacity: 0.7; cursor: not-allowed; }

.gaps {
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-radius: var(--radius-md);
  padding: 14px 18px;
  margin-bottom: 20px;
}
.gaps-title {
  font-weight: 600;
  color: var(--color-warning);
  margin-bottom: 6px;
}
.gaps ul {
  padding-left: 18px;
  color: var(--color-text);
  font-size: 14px;
}

.qa-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.qa-col { display: flex; flex-direction: column; gap: 10px; }
.col-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4px;
}

.qa-card {
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: var(--radius-md);
  padding: 14px;
}
.qa-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 8px;
}
.qa-q {
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text);
}
.qa-score {
  font-size: 12px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 999px;
  white-space: nowrap;
  background: var(--color-bg);
  border: 1px solid var(--color-surface-border);
}
.qa-full { color: var(--color-success); border-color: rgba(34, 197, 94, 0.3); }
.qa-partial { color: var(--color-warning); border-color: rgba(245, 158, 11, 0.3); }
.qa-miss { color: var(--color-error); border-color: rgba(239, 68, 68, 0.3); }
.qa-na { color: var(--color-text-tertiary); }

.qa-a {
  background: var(--color-bg);
  border-left: 3px solid var(--color-surface-border);
  padding: 10px 14px;
  font-size: 13px;
  color: var(--color-text);
  white-space: pre-wrap;
  border-radius: var(--radius-sm);
  margin-bottom: 8px;
  font-style: italic;
}

.qa-notes {
  list-style: none;
  padding: 0;
  margin: 0;
  font-size: 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.note-hit { color: var(--color-success); }
.note-miss { color: var(--color-text-tertiary); }
.note-hit::before { content: '✓ '; }

.footer-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}
.btn-ghost {
  background: transparent;
  color: var(--color-text);
  border: 1px solid var(--color-surface-border);
  padding: 8px 14px;
  border-radius: var(--radius-md);
  text-decoration: none;
  font-size: 13px;
  cursor: pointer;
  font-family: inherit;
}
.btn-ghost:hover { background: var(--color-surface-hover); }
.btn-ghost:disabled { opacity: 0.5; cursor: not-allowed; }

@media (max-width: 820px) {
  .qa-grid { grid-template-columns: 1fr; }
  .header-row { flex-direction: column; }
  .score-circle { align-self: flex-start; }
}
</style>
