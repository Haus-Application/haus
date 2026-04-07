// Kasa devices: switches, dimmers, fans. Like the Bluth family, but they actually respond.

export interface KasaDevice {
  ip: string
  alias: string
  model: string
  device_type: string  // "switch", "dimmer", "fan"
  on: boolean
  brightness: number
  fan_speed: number
}

export function useKasa() {
  const devices = ref<KasaDevice[]>([])
  const { on } = useWebSocket()

  on('kasa:state', (payload: KasaDevice[]) => {
    devices.value = payload
  })

  async function toggleDevice(ip: string, isOn: boolean) {
    // Optimistic: update UI instantly
    const dev = devices.value.find(d => d.ip === ip)
    if (dev) dev.on = isOn

    await fetch(`/api/kasa/devices/${encodeURIComponent(ip)}/state`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ on: isOn }),
    })
  }

  async function setBrightness(ip: string, brightness: number) {
    // Optimistic
    const dev = devices.value.find(d => d.ip === ip)
    if (dev) dev.brightness = brightness

    await fetch(`/api/kasa/devices/${encodeURIComponent(ip)}/brightness`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ brightness }),
    })
  }

  async function setFanSpeed(ip: string, speed: number) {
    // Optimistic
    const dev = devices.value.find(d => d.ip === ip)
    if (dev) dev.fan_speed = speed

    await fetch(`/api/kasa/devices/${encodeURIComponent(ip)}/fan-speed`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ speed }),
    })
  }

  onMounted(async () => {
    try {
      const res = await fetch('/api/kasa/devices')
      if (res.ok) devices.value = await res.json()
    } catch {
      // server not up yet, WebSocket will catch state when it arrives
    }
  })

  return { devices, toggleDevice, setBrightness, setFanSpeed }
}
