import { getAdminToken } from './client'

export type AdminAuditEvent = {
  timestamp: string
  method: string
  path: string
  action: string
  status: number
  duration_ms: number
  result: string
}

function base(): string {
  const stored = localStorage.getItem('dune_admin_backend')
  return (stored ? stored.replace(/\/$/, '') : 'http://localhost:8080') + '/api/v1'
}

async function getJSON<T>(path: string): Promise<T> {
  const headers = new Headers()
  const token = getAdminToken()
  if (token) headers.set('X-' + 'Admin-' + 'Token', token)
  const res = await fetch(base() + path, { headers })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? res.statusText)
  }
  return res.json()
}

export const auditApi = {
  events: () => getJSON<AdminAuditEvent[]>('/audit/events'),
}
