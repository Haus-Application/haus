// Like getting into the movie business — fake it till the devices show up.

interface ScanDevice {
  ip: string
  ipv6?: string[]
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

interface ScanStage {
  stage: string
  status: 'pending' | 'running' | 'complete'
  message: string
  count?: number
}

const STAGE_NAMES: Record<string, string> = {
  arp: 'Host Discovery',
  ipv6: 'IPv6 Discovery',
  oui: 'Manufacturer Lookup',
  ports: 'Port Scan',
  probe: 'Device Probes',
  mdns: 'mDNS Services',
  classify: 'Classification',
  // Legacy keys kept so older server builds don't render unlabeled rows.
  kasa: 'Device Probes',
  cast: 'Device Probes',
}

export function useScan() {
  const devices = ref<ScanDevice[]>([])
  const stages = ref<ScanStage[]>([])
  const status = ref<'idle' | 'scanning' | 'complete' | 'error'>('idle')
  const scanId = ref<string | null>(null)
  let eventSource: EventSource | null = null

  // Load persisted devices from DB on mount
  async function loadDevices() {
    try {
      const res = await fetch('/api/devices')
      if (res.ok) {
        const data = await res.json()
        if (Array.isArray(data) && data.length > 0) {
          devices.value = data
          status.value = 'complete'
        }
      }
    } catch { /* server not up yet */ }
  }

  onMounted(() => loadDevices())

  async function startScan() {
    status.value = 'scanning'
    stages.value = []
    // Don't clear devices — scan enriches, never wipes

    let data: { scan_id: string }
    try {
      const res = await fetch('/api/scan', { method: 'POST' })
      data = await res.json()
    } catch {
      status.value = 'error'
      return
    }
    scanId.value = data.scan_id

    eventSource = new EventSource(`/api/scan/stream?scan_id=${data.scan_id}`)

    eventSource.addEventListener('stage', (e) => {
      const stage = JSON.parse(e.data) as ScanStage
      stage.stage = STAGE_NAMES[stage.stage] || stage.stage
      const idx = stages.value.findIndex(s => s.stage === stage.stage)
      if (idx >= 0) {
        stages.value[idx] = stage
      } else {
        stages.value.push(stage)
      }
    })

    eventSource.addEventListener('device', (e) => {
      const device = JSON.parse(e.data) as ScanDevice
      const idx = devices.value.findIndex(d => d.ip === device.ip)
      if (idx >= 0) {
        devices.value[idx] = device
      } else {
        devices.value.push(device)
      }
    })

    eventSource.addEventListener('complete', async () => {
      // Marry me! This just works.
      status.value = 'complete'
      eventSource?.close()
      eventSource = null
      // Re-fetch from DB to get fully persisted results
      await loadDevices()
    })

    eventSource.onerror = () => {
      status.value = 'error'
      eventSource?.close()
      eventSource = null
    }
  }

  const devicesByCategory = computed(() => {
    const groups: Record<string, ScanDevice[]> = {}
    const order = ['lighting', 'media', 'smart_home', 'energy', 'network', 'compute', 'unknown']
    for (const cat of order) {
      groups[cat] = []
    }
    for (const device of devices.value) {
      const cat = device.category || 'unknown'
      if (!groups[cat]) groups[cat] = []
      groups[cat].push(device)
    }
    // Remove empty categories
    for (const key of Object.keys(groups)) {
      if (groups[key].length === 0) delete groups[key]
    }
    return groups
  })

  const totalDevices = computed(() => devices.value.length)

  onUnmounted(() => {
    eventSource?.close()
  })

  return { devices, stages, status, scanId, startScan, devicesByCategory, totalDevices }
}
