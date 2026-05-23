import { getAdminToken } from './client'

export type DBFunctionRow = {
  oid: string
  schema: string
  name: string
  arguments: string
  result_type: string
  language: string
  volatility: string
  category: string
  references: string[]
  summary: string
}

export type DBFunctionInspection = DBFunctionRow & {
  definition: string
  risk: string
  notes: string[]
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

export const dbFunctionApi = {
  list: (term = '', category = '') => getJSON<DBFunctionRow[]>(`/database/functions?term=${encodeURIComponent(term)}&category=${encodeURIComponent(category)}`),
  inspect: (oid: string) => getJSON<DBFunctionInspection>(`/database/functions/inspect?oid=${encodeURIComponent(oid)}`),
}
