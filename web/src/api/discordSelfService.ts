import { getAdminToken } from './client'
import type {
  PlayerProfileCharXP,
  PlayerProfileCurrency,
  PlayerProfileFaction,
  PlayerProfileInventorySummary,
  PlayerProfileJourneySummary,
  PlayerProfileLocation,
  PlayerProfileSectionError,
  PlayerProfileSpec,
} from './playerProfile'

function getApiBase(): string {
  const stored = localStorage.getItem('dune_admin_backend')
  return (stored ? stored.replace(/\/$/, '') : 'http://localhost:8080') + '/api/v1'
}

async function headers(body?: unknown): Promise<Record<string, string>> {
  const h: Record<string, string> = {}
  const admin = getAdminToken()
  if (admin) h['X-Admin-Token'] = admin
  if (body !== undefined) h['Content-Type'] = 'application/json'
  return h
}

async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
  const response = await fetch(`${getApiBase()}${path}`, {
    method,
    headers: await headers(body),
    credentials: 'include',
    body: body === undefined ? undefined : JSON.stringify(body),
  })
  if (!response.ok) {
    const err = await response.json().catch(() => ({ error: response.statusText }))
    throw new Error(err.error ?? response.statusText)
  }
  return response.json()
}

export type DiscordPlayerLink = {
  discord_id: string
  player_id: number
  player_name?: string
  notes?: string
  linked_by_discord_id?: string
  linked_by_auth_type?: string
  created_at: string
  updated_at: string
}

export type UpsertDiscordPlayerLinkPayload = {
  discord_id: string
  player_id: number
  player_name?: string
  notes?: string
}

export type SelfPlayerCardSummary = {
  discord_id: string
  player_id: number
  player_name?: string
  class?: string
  map?: string
  online_status?: string
  location?: PlayerProfileLocation
  inventory_summary: PlayerProfileInventorySummary
  vehicle_count: number
  currencies: PlayerProfileCurrency[]
  factions: PlayerProfileFaction[]
  specializations: PlayerProfileSpec[]
  character_xp?: PlayerProfileCharXP
  journey_summary: PlayerProfileJourneySummary
  section_errors: PlayerProfileSectionError[]
}

export const discordSelfServiceApi = {
  listLinks: () => req<DiscordPlayerLink[]>('GET', '/auth/discord/player-links'),
  upsertLink: (payload: UpsertDiscordPlayerLinkPayload) => req<DiscordPlayerLink>('POST', '/auth/discord/player-links', payload),
  deleteLink: (discordId: string) => req<{ ok: string }>('DELETE', `/auth/discord/player-links/${encodeURIComponent(discordId)}`),
  selfLink: () => req<DiscordPlayerLink>('GET', '/self/player-link'),
  selfPlayerCard: () => req<SelfPlayerCardSummary>('GET', '/self/player-card'),
}
