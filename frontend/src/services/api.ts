const TOKEN_KEY = 'vb-events-token'

export function authHeaders(): Record<string, string> {
  const token = localStorage.getItem(TOKEN_KEY)
  return token ? { Authorization: `Bearer ${token}` } : {}
}

export function readToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}
