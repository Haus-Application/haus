// Hue integration. This one actually has good lighting, unlike the model home.

export interface HueLight {
  id: string
  name: string
  on: boolean
  brightness: number
  room_name: string
  reachable: boolean
}

export interface HueRoom {
  id: string
  name: string
  grouped_light_id: string
  light_count: number
  any_on: boolean
}

export interface HueScene {
  id: string
  name: string
  room_name: string
}

export function useHue() {
  const lights = ref<HueLight[]>([])
  const rooms = ref<HueRoom[]>([])
  const scenes = ref<HueScene[]>([])
  const connected = ref(false)
  const { on } = useWebSocket()

  on('hue:state', (payload: { lights: HueLight[]; rooms: HueRoom[]; scenes: HueScene[] }) => {
    lights.value = payload.lights ?? []
    rooms.value = payload.rooms ?? []
    scenes.value = payload.scenes ?? []
    connected.value = true
  })

  async function toggleLight(id: string, isOn: boolean) {
    await fetch(`/api/hue/lights/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ on: isOn }),
    })
  }

  async function setLightBrightness(id: string, brightness: number) {
    await fetch(`/api/hue/lights/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ brightness }),
    })
  }

  async function activateScene(id: string) {
    await fetch(`/api/hue/scenes/${encodeURIComponent(id)}/activate`, {
      method: 'POST',
    })
  }

  async function controlRoom(id: string, isOn?: boolean, brightness?: number) {
    const body: Record<string, any> = {}
    if (isOn !== undefined) body.on = isOn
    if (brightness !== undefined) body.brightness = brightness
    await fetch(`/api/hue/rooms/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
  }

  async function discoverBridges(): Promise<{ id: string; ip: string }[]> {
    const res = await fetch('/api/hue/discover')
    if (!res.ok) return []
    return await res.json()
  }

  async function pair(bridgeIP: string): Promise<boolean> {
    const res = await fetch('/api/hue/pair', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ bridge_ip: bridgeIP }),
    })
    return res.ok
  }

  async function getStatus(): Promise<{ connected: boolean; bridge_ip?: string }> {
    const res = await fetch('/api/hue/status')
    if (!res.ok) return { connected: false }
    return await res.json()
  }

  async function disconnectBridge(): Promise<void> {
    await fetch('/api/hue/disconnect', { method: 'POST' })
    connected.value = false
    lights.value = []
    rooms.value = []
    scenes.value = []
  }

  return {
    lights,
    rooms,
    scenes,
    connected,
    toggleLight,
    setLightBrightness,
    activateScene,
    controlRoom,
    discoverBridges,
    pair,
    getStatus,
    disconnect: disconnectBridge,
  }
}
