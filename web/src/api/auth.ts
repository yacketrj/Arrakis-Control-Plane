export function getAdminToken(): string {
  return localStorage.getItem('dune_admin_token') || ''
}

export function setAdminToken(token: string): void {
  const trimmed = token.trim()
  if (trimmed) localStorage.setItem('dune_admin_token', trimmed)
  else localStorage.removeItem('dune_admin_token')
}

export function getAuthHeaders(hasBody = false): Record<string, string> {
  const headers: Record<string, string> = {}
  if (hasBody) headers['Content-Type'] = 'application/json'
  const token = getAdminToken()
  if (token) headers['X-Admin-Token'] = token
  return headers
}
