import { getAdminToken, type MutateResult } from './client'

function getApiBase(): string {
  const stored = localStorage.getItem('dune_admin_backend')
  return (stored ? stored.replace(/\/$/, '') : 'http://localhost:8080') + '/api/v1'
}

async function req<T>(method: string, path: string, body: unknown, reason?: string): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  const admin = getAdminToken()
  if (admin) headers['X-Admin-Token'] = admin
  if (reason?.trim()) headers['X-Admin-Reason'] = reason.trim()
  const response = await fetch(`${getApiBase()}${path}`, { method, headers, body: JSON.stringify(body) })
  if (!response.ok) {
    const err = await response.json().catch(() => ({ error: response.statusText }))
    throw new Error(err.error ?? response.statusText)
  }
  return response.json()
}

export const inventoryStudioMutations = {
  setItemStackSize: (id: number, stackSize: number, reason?: string) =>
    req<MutateResult>('POST', '/players/item/stack-size', { id, stack_size: stackSize }, reason),
}
