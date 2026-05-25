import { getAdminToken } from './client'

function getApiBase(): string {
  const stored = localStorage.getItem('dune_admin_backend')
  return (stored ? stored.replace(/\/$/, '') : 'http://localhost:8080') + '/api/v1'
}

async function authHeaders(): Promise<Record<string, string>> {
  const headers: Record<string, string> = {}
  const admin = getAdminToken()
  if (admin) headers['X-Admin-Token'] = admin
  return headers
}

export type PlayerProfileSectionError = { section: string; error: string }
export type PlayerProfileOnlineState = { player_id: number; name: string; map: string; status: string; last_seen?: string }
export type PlayerProfileLocation = { map: string; source?: string }
export type PlayerProfileInventorySummary = { total_items: number; total_stack_size: number; unique_templates: number; preview_items: PlayerProfileInventoryItem[] }
export type PlayerProfileJourneySummary = { total_nodes: number; complete_nodes: number; revealed_nodes: number; pending_rewards: number; preview_nodes: PlayerProfileJourneyNode[] }
export type PlayerProfileCharXP = { xp: number; level: number }

export type PlayerProfileIdentity = {
  id: number
  account_id: number
  controller_id: number
  fls_id: string
  name: string
  class: string
  map: string
  faction_id: number
  online_status: string
}

export type PlayerProfileInventoryItem = {
  id: number
  template_id: string
  name: string
  stack_size: number
  quality: number
  durability: string
  max_durability: string
}

export type PlayerProfileVehicle = {
  id: number
  class: string
  map: string
  chassis_durability: number
  vehicle_name: string
  is_recovered: boolean
  is_backup: boolean
}

export type PlayerProfileCurrency = { player_id: number; currency_id: number; balance: number }
export type PlayerProfileFaction = { actor_id: number; faction_id: number; faction_name: string; reputation: number; scrips: number }
export type PlayerProfileSpec = { player_id: number; track_type: string; xp: number; level: number }
export type PlayerProfileJourneyNode = { node_id: string; is_complete: boolean; is_revealed: boolean; has_pending_reward: boolean }
export type PlayerProfileEvent = { actor_id: number; universe_time: string; map: string; event_type: number; x: number; y: number; z: number; custom_data: string }
export type PlayerProfileDungeon = { dungeon_id: string; difficulty: string; duration_ms: number; players_num: number; completion_id: number }

export type PlayerProfile = {
  player_id: number
  identity?: PlayerProfileIdentity
  online_state?: PlayerProfileOnlineState
  location?: PlayerProfileLocation
  inventory_summary: PlayerProfileInventorySummary
  vehicles: PlayerProfileVehicle[]
  currencies: PlayerProfileCurrency[]
  factions: PlayerProfileFaction[]
  specializations: PlayerProfileSpec[]
  character_xp?: PlayerProfileCharXP
  journey_summary: PlayerProfileJourneySummary
  recent_events: PlayerProfileEvent[]
  dungeon_history: PlayerProfileDungeon[]
  section_errors: PlayerProfileSectionError[]
}

export async function getPlayerProfile(playerId: number): Promise<PlayerProfile> {
  const res = await fetch(`${getApiBase()}/players/${playerId}/profile`, { headers: await authHeaders() })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? res.statusText)
  }
  return res.json()
}
