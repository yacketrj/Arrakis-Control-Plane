declare global {
  interface Window { Clerk?: { session?: { getToken(): Promise<string | null> } } }
}

const ADMIN_TOKEN_SESSION_KEY = 'dune_admin_token_session'
const LEGACY_ADMIN_TOKEN_LOCAL_KEY = 'dune_admin_token'
const ADMIN_TOKEN_PATTERN = /^[A-Za-z0-9_-]{43}$/

function getApiBase(): string {
  const stored = localStorage.getItem('dune_admin_backend')
  return (stored ? stored.replace(/\/$/, '') : 'http://localhost:8080') + '/api/v1'
}

export function adminTokenValidationError(value: string): string | null {
  const token = value.trim()
  if (!token) return 'Browser Access Key is required.'
  if (token !== value || /[\x00\r\n\t ]/.test(value)) return 'Browser Access Key must not contain whitespace or control characters.'
  if (!ADMIN_TOKEN_PATTERN.test(value)) return 'Browser Access Key must be exactly 43 base64url characters: A-Z, a-z, 0-9, underscore, or dash.'
  const lower = value.toLowerCase()
  if (['changeme', 'change-me', 'password', 'admin', 'admin-token', 'replace-me', 'replace_with_token', '<your_admin_token>'].includes(lower)) {
    return 'Browser Access Key cannot use a placeholder or default value.'
  }
  return null
}

export function isAdminTokenValid(value: string): boolean {
  return adminTokenValidationError(value) === null
}

export function getAdminToken(): string {
  const sessionToken = sessionStorage.getItem(ADMIN_TOKEN_SESSION_KEY) || ''
  if (isAdminTokenValid(sessionToken)) return sessionToken
  if (sessionToken) sessionStorage.removeItem(ADMIN_TOKEN_SESSION_KEY)

  const legacyToken = localStorage.getItem(LEGACY_ADMIN_TOKEN_LOCAL_KEY) || ''
  if (isAdminTokenValid(legacyToken)) {
    sessionStorage.setItem(ADMIN_TOKEN_SESSION_KEY, legacyToken)
    localStorage.removeItem(LEGACY_ADMIN_TOKEN_LOCAL_KEY)
    return legacyToken
  }
  if (legacyToken) localStorage.removeItem(LEGACY_ADMIN_TOKEN_LOCAL_KEY)
  return ''
}

export function setAdminToken(value: string): void {
  const token = value.trim()
  localStorage.removeItem(LEGACY_ADMIN_TOKEN_LOCAL_KEY)
  if (!token) {
    sessionStorage.removeItem(ADMIN_TOKEN_SESSION_KEY)
    return
  }
  const err = adminTokenValidationError(token)
  if (err) throw new Error(err)
  sessionStorage.setItem(ADMIN_TOKEN_SESSION_KEY, token)
}

export function clearAdminToken(): void {
  sessionStorage.removeItem(ADMIN_TOKEN_SESSION_KEY)
  localStorage.removeItem(LEGACY_ADMIN_TOKEN_LOCAL_KEY)
}

export function getWsBase(): string { return getApiBase().replace(/^http/, 'ws') }

const BASE = getApiBase()

async function headers(body?: unknown, reason?: string): Promise<Record<string, string>> {
  const h: Record<string, string> = {}
  if (body) h['Content-Type'] = 'application/json'
  if (reason?.trim()) h['X-Admin-Reason'] = reason.trim()
  const admin = getAdminToken()
  if (admin) h['X-Admin-Token'] = admin
  else {
    const clerk = await window.Clerk?.session?.getToken()
    if (clerk) h.Authorization = `Bearer ${clerk}`
  }
  return h
}

function isMutatingMethod(method: string): boolean {
  return ['POST', 'PUT', 'PATCH', 'DELETE'].includes(method.toUpperCase())
}

function localMutationSafety(method: string, path: string): MutationSafetyClass {
  const lower = path.toLowerCase()
  const destructive = method.toUpperCase() === 'DELETE' || lower.includes('/wipe') || lower.includes('/delete') || lower.includes('/reset') || lower.includes('/kick') || lower.includes('/blueprints/import')
  const high = destructive || lower.includes('/give') || lower.includes('/award') || lower.includes('/grant') || lower.includes('/teleport') || lower.includes('/journey/') || lower.includes('/set-faction') || lower.includes('/repair') || lower.includes('/storage/') || lower.includes('/database/sql') || lower.includes('/battlegroup/exec')
  return {
    action: `${method.toLowerCase()}:${path.replace(/^\/api\/v1\//, '').replace(/^\//, '').replaceAll('/', '.') || 'root'}`,
    risk: destructive ? 'destructive' : high ? 'high' : 'medium',
    requires_reason: high,
    reason_enforcement_enabled: false,
    requires_preview: high,
    destructive,
    operator_warnings: high ? ['This action changes player or server state and will be written to the audit log.'] : [],
  }
}

function safetyText(safety: MutationSafetyClass): string {
  return [
    `Risk: ${safety.risk}`,
    safety.recommended_path ? `Recommended path: ${safety.recommended_path}` : '',
    safety.rollback_hint ? `Rollback hint: ${safety.rollback_hint}` : '',
    ...(safety.operator_warnings ?? []).map(w => `Warning: ${w}`),
  ].filter(Boolean).join('\n')
}

async function classifyMutation(method: string, path: string): Promise<MutationSafetyClass> {
  try {
    const res = await fetch(`${BASE}/mutation-safety/classify?method=${encodeURIComponent(method)}&path=${encodeURIComponent(path)}`, { headers: await headers() })
    if (!res.ok) return localMutationSafety(method, path)
    return res.json()
  } catch {
    return localMutationSafety(method, path)
  }
}

async function ensureAdminReasonForMutation(method: string, path: string, reason?: string): Promise<string | undefined> {
  if (reason?.trim()) return reason.trim()
  if (!isMutatingMethod(method)) return undefined

  const safety = await classifyMutation(method, path)
  if (!safety.requires_reason && !safety.requires_preview) return undefined

  if (safety.requires_preview || safety.operator_warnings?.length || safety.recommended_path || safety.rollback_hint) {
    const ok = window.confirm(`${safetyText(safety)}\n\nContinue?`)
    if (!ok) throw new Error('Mutation cancelled')
  }

  if (safety.requires_reason) {
    const entered = window.prompt('Admin reason required for audit log:', '')?.trim()
    if (!entered) throw new Error('Admin reason is required for this action')
    return entered
  }

  return undefined
}

async function req<T>(method: string, path: string, body?: unknown, reason?: string): Promise<T> {
  const finalReason = await ensureAdminReasonForMutation(method, path, reason)
  const res = await fetch(`${BASE}${path}`, { method, headers: await headers(body, finalReason), body: body ? JSON.stringify(body) : undefined })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? res.statusText)
  }
  return res.json()
}

export type TunnelStatus = { name: string; local_addr: string; remote_addr: string }
export type Status = { ssh_connected: boolean; db_connected: boolean; pod_ns: string; pod_ip?: string; ssh_host: string; tunnel_mode?: string; tunnels?: TunnelStatus[]; startup_connect_error?: string }
export type Player = { id: number; account_id: number; controller_id: number; fls_id: string; name: string; class: string; map: string; faction_id: number; online_status: string }
export type InventoryItem = { id: number; template_id: string; name: string; stack_size: number; quality: number; durability: string; max_durability: string }
export type CurrencyRow = { player_id: number; currency_id: number; balance: number }
export type FactionRep = { actor_id: number; faction_id: number; faction_name: string; reputation: number; scrips: number }
export type SpecTrack = { player_id: number; track_type: string; xp: number; level: number }
export type JourneyNode = { node_id: string; is_complete: boolean; is_revealed: boolean; has_pending_reward: boolean }
export type BlueprintRow = { id: number; owner_name: string; item_id: number; pieces: number; placeables: number }
export type LogPod = { namespace: string; name: string; display?: string }
export type MutateResult = { ok: string }
export type BGOutput = { output: string }
export type BGHealthSection = { name: string; description: string; command: string; output: string; error?: string }
export type BGHealth = { namespace: string; checked_at: string; sections: BGHealthSection[] }
export type VehicleRow = { id: number; class: string; map: string; chassis_durability: number; vehicle_name: string; is_recovered: boolean; is_backup: boolean }
export type CheatEntry = { fls_id: string; cheat_type: string; event_time: string; character_name: string }
export type GameEvent = { actor_id: number; universe_time: string; map: string; event_type: number; x: number; y: number; z: number; custom_data: string }
export type DungeonRecord = { dungeon_id: string; difficulty: string; duration_ms: number; players_num: number; completion_id: number }
export type TeleportLocation = { name: string; x: number; y: number; z: number }
export type OnlineRow = { player_id: number; name: string; map: string; status: string; last_seen: string }
export type ItemTemplate = { id: string; name: string }
export type GiveItemAugment = { name: string; grade: number; roll?: number; rolls?: number[]; roll_count?: number; effect_indices?: number[] }
export type GiveItemRow = { template: string; qty: number; quality: number; stack_size: number; augments?: GiveItemAugment[] }
export type MutationSafetyClass = { action: string; risk: string; requires_reason: boolean; reason_enforcement_enabled?: boolean; requires_preview: boolean; destructive: boolean; rollback_hint?: string; operator_warnings?: string[]; recommended_path?: string }

export const api = {
  status: () => req<Status>('GET', '/status'), reconnect: () => req<Status>('POST', '/reconnect'),
  mutationSafety: { classify: (method: string, path: string) => req<MutationSafetyClass>('GET', `/mutation-safety/classify?method=${encodeURIComponent(method)}&path=${encodeURIComponent(path)}`) },
  battlegroup: {
    status: () => req<BGOutput>('GET', '/battlegroup/status'),
    health: () => req<BGHealth>('GET', '/battlegroup/health'),
    exec: (cmd: string, reason?: string) => req<BGOutput>('POST', '/battlegroup/exec', { cmd }, reason),
    pods: () => req<{ pods: string[]; namespace: string }>('GET', '/battlegroup/pods'),
  },
  players: {
    list: () => req<Player[]>('GET', '/players'), online: () => req<OnlineRow[]>('GET', '/players/online'), currency: () => req<CurrencyRow[]>('GET', '/players/currency'), factions: () => req<FactionRep[]>('GET', '/players/factions'), specs: () => req<SpecTrack[]>('GET', '/players/specs'), templates: () => req<ItemTemplate[]>('GET', '/players/templates'), refreshTemplates: () => req<ItemTemplate[]>('POST', '/players/templates/refresh'), inventory: (id: number) => req<InventoryItem[]>('GET', `/players/${id}/inventory`), journey: (accountId: number) => req<JourneyNode[]>('GET', `/players/${accountId}/journey`), giveItem: (player_id: number, template: string, qty: number, quality: number, reason?: string) => req<MutateResult>('POST', '/players/give-item', { player_id, template, qty, quality }, reason), giveItems: (player_id: number, items: GiveItemRow[], reason?: string) => req<MutateResult>('POST', '/players/give-item', { player_id, items }, reason), grantLive: (controller_id: number, template: string, amount: number, reason?: string) => req<MutateResult>('POST', '/players/grant-live', { controller_id, template, amount }, reason), giveCurrency: (player_id: number, amount: number, reason?: string) => req<MutateResult>('POST', '/players/give-currency', { player_id, amount }, reason), giveFactionRep: (actor_id: number, faction_id: number, delta: number, reason?: string) => req<MutateResult>('POST', '/players/give-faction-rep', { actor_id, faction_id, delta }, reason), giveScrip: (actor_id: number, delta: number, reason?: string) => req<MutateResult>('POST', '/players/give-scrip', { actor_id, delta }, reason), awardXP: (player_id: number, track_type: string, delta: number, reason?: string) => req<MutateResult>('POST', '/players/award-xp', { player_id, track_type, delta }, reason), awardCharXP: (player_id: number, amount: number, reason?: string) => req<MutateResult>('POST', '/players/award-char-xp', { player_id, amount }, reason), awardIntel: (player_id: number, amount: number, reason?: string) => req<MutateResult>('POST', '/players/award-intel', { player_id, amount }, reason), kick: (player_id: number, reason?: string) => req<MutateResult>('POST', '/players/kick', { player_id }, reason), deleteItem: (id: number, reason?: string) => req<MutateResult>('DELETE', `/players/item/${id}`, undefined, reason), resetSpec: (player_id: number, track_type: string, reason?: string) => req<MutateResult>('POST', '/players/reset-spec', { player_id, track_type }, reason), setFactionTier: (actor_id: number, faction_id: number, tier: number, reason?: string) => req<MutateResult>('POST', '/players/set-faction-tier', { actor_id, faction_id, tier }, reason), journeyComplete: (account_id: number, node_id: string, reason?: string) => req<MutateResult>('POST', '/players/journey/complete', { account_id, node_id }, reason), journeyReset: (account_id: number, node_id: string, reason?: string) => req<MutateResult>('POST', '/players/journey/reset', { account_id, node_id }, reason), journeyWipe: (account_id: number, reason?: string) => req<MutateResult>('POST', '/players/journey/wipe', { account_id }, reason), deleteTutorials: (player_id: number, reason?: string) => req<MutateResult>('POST', '/players/delete-tutorials', { player_id }, reason), wipeCodex: (account_id: number, reason?: string) => req<MutateResult>('POST', '/players/wipe-codex', { account_id }, reason), charXPCurrent: (id: number) => req<{xp: number; level: number}>('GET', `/players/${id}/char-xp`), specs_for: (id: number) => req<SpecTrack[]>('GET', `/players/${id}/specs`), setSpecXP: (player_id: number, track_type: string, xp_amount: number, reason?: string) => req<MutateResult>('POST', '/players/set-spec-xp', { player_id, track_type, xp_amount }, reason), vehicles: (account_id: number) => req<VehicleRow[]>('GET', `/players/${account_id}/vehicles`), repairItem: (id: number, reason?: string) => req<MutateResult>('POST', '/players/repair-item', { id }, reason), partitions: () => req<TeleportLocation[]>('GET', '/players/partitions'), teleport: (fls_id: string, partition_label: string, reason?: string) => req<MutateResult>('POST', '/players/teleport', { fls_id, partition_label }, reason), events: (id: number) => req<GameEvent[]>('GET', `/players/${id}/events`), dungeons: (id: number) => req<DungeonRecord[]>('GET', `/players/${id}/dungeons`),
  },
  database: { tables: () => req<{name: string; row_count: number}[]>('GET', '/database/tables'), describe: (table: string) => req<{table: string; columns: {name: string; data_type: string; nullable: string}[]}>('GET', `/database/describe?table=${encodeURIComponent(table)}`), sample: (table: string, limit = 20) => req<{table: string; headers: string[]; rows: string[][]}>('GET', `/database/sample?table=${encodeURIComponent(table)}&limit=${limit}`), search: (term: string) => req<{headers: string[]; rows: string[][]}>('GET', `/database/search?term=${encodeURIComponent(term)}`), sql: (sql: string, reason?: string) => req<{result: string}>('POST', '/database/sql', { sql }, reason) },
  logs: { pods: () => req<LogPod[]>('GET', '/logs/pods'), cheats: () => req<CheatEntry[]>('GET', '/logs/cheats') },
  storage: { list: () => req<{id: number; class: string; map: string; item_count: number}[]>('GET', '/storage'), items: (id: number) => req<InventoryItem[]>('GET', `/storage/${id}/items`), giveItem: (id: number, template: string, qty: number, quality: number, reason?: string) => req<MutateResult>('POST', `/storage/${id}/give-item`, { template, qty, quality }, reason) },
  blueprints: { list: () => req<BlueprintRow[]>('GET', '/blueprints'), exportUrl: (id: number) => `${BASE}/blueprints/${id}/export`, import: async (file: File, player_id: number, reason?: string) => { const fd = new FormData(); fd.append('file', file); fd.append('player_id', String(player_id)); return fetch(`${BASE}/blueprints/import`, { method: 'POST', headers: await headers(undefined, reason), body: fd }).then(r => r.json()) } },
}
