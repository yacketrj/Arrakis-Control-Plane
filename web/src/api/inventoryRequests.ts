import { getAdminToken } from './client'

const BASE = getApiBase()

export type InventoryRequestScope = 'personal' | 'guild'
export type InventoryRequestStatus = 'open' | 'ordered' | 'fulfilled' | 'cancelled'
export type InventoryOrderStatus = 'open' | 'filled' | 'cancelled'

export type InventoryRequest = {
  id: string
  scope: InventoryRequestScope
  requester_discord_id?: string
  requester_name?: string
  guild_id?: string
  player_id?: string
  item_name: string
  item_template_id?: string
  quantity: number
  notes?: string
  status: InventoryRequestStatus
  order_id?: string
  created_at: string
  updated_at: string
}

export type InventoryOrder = {
  id: string
  scope: InventoryRequestScope
  guild_id?: string
  requester_discord_id?: string
  assignee_discord_id?: string
  assignee_name?: string
  request_ids: string[]
  status: InventoryOrderStatus
  notes?: string
  created_at: string
  updated_at: string
  completed_at?: string
}

export type CreateInventoryRequestPayload = {
  scope: InventoryRequestScope
  requester_discord_id?: string
  requester_name?: string
  guild_id?: string
  player_id?: string
  item_name: string
  item_template_id?: string
  quantity: number
  notes?: string
}

export type UpdateInventoryRequestPayload = {
  status?: InventoryRequestStatus
  order_id?: string
  notes?: string
}

export type CreateInventoryOrderPayload = {
  scope: InventoryRequestScope
  guild_id?: string
  requester_discord_id?: string
  assignee_discord_id?: string
  assignee_name?: string
  request_ids: string[]
  notes?: string
}

export type UpdateInventoryOrderPayload = {
  status?: InventoryOrderStatus
  assignee_discord_id?: string
  assignee_name?: string
  notes?: string
}

export type InventoryRequestFilters = {
  scope?: InventoryRequestScope | 'all'
  status?: InventoryRequestStatus | 'all'
  guild_id?: string
  requester_discord_id?: string
}

export type InventoryOrderFilters = {
  scope?: InventoryRequestScope | 'all'
  status?: InventoryOrderStatus | 'all'
  guild_id?: string
}

function getApiBase(): string {
  const stored = localStorage.getItem('dune_admin_backend')
  return (stored ? stored.replace(/\/$/, '') : 'http://localhost:8080') + '/api/v1'
}

function buildQuery(filters: Record<string, string | undefined>): string {
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(filters)) {
    const trimmed = value?.trim()
    if (trimmed && trimmed !== 'all') params.set(key, trimmed)
  }
  const query = params.toString()
  return query ? `?${query}` : ''
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const headers: Record<string, string> = {}
  const token = getAdminToken()
  if (token) headers['X-Admin-Token'] = token
  if (body !== undefined) headers['Content-Type'] = 'application/json'

  const response = await fetch(`${BASE}${path}`, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  })
  if (!response.ok) {
    const err = await response.json().catch(() => ({ error: response.statusText }))
    throw new Error(err.error ?? response.statusText)
  }
  return response.json()
}

export const inventoryRequestsApi = {
  listRequests: (filters: InventoryRequestFilters = {}) => request<InventoryRequest[]>('GET', `/inventory/requests${buildQuery(filters)}`),
  createRequest: (payload: CreateInventoryRequestPayload) => request<InventoryRequest>('POST', '/inventory/requests', payload),
  updateRequest: (id: string, payload: UpdateInventoryRequestPayload) => request<InventoryRequest>('PATCH', `/inventory/requests/${encodeURIComponent(id)}`, payload),
  listOrders: (filters: InventoryOrderFilters = {}) => request<InventoryOrder[]>('GET', `/inventory/orders${buildQuery(filters)}`),
  createOrder: (payload: CreateInventoryOrderPayload) => request<InventoryOrder>('POST', '/inventory/orders', payload),
  updateOrder: (id: string, payload: UpdateInventoryOrderPayload) => request<InventoryOrder>('PATCH', `/inventory/orders/${encodeURIComponent(id)}`, payload),
}
