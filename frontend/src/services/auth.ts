import { authHeaders } from '@/services/api'
import type { Participant } from '@/types/event'

const API_BASE = '/api'

export async function loginRequest(email: string, password: string): Promise<string> {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })
  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Invalid email or password.')
    }
    throw new Error('Login failed. Please try again.')
  }
  const data = (await response.json()) as { token: string }
  return data.token
}

export async function registerRequest(name: string, email: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email, password }),
  })
  if (!response.ok) {
    if (response.status === 409) {
      throw new Error('That email is already registered.')
    }
    if (response.status === 400) {
      throw new Error('Please check your details — passwords must be at least 8 characters.')
    }
    throw new Error('Registration failed. Please try again.')
  }
}

export async function meRequest(): Promise<Participant> {
  const response = await fetch(`${API_BASE}/auth/me`, { headers: authHeaders() })
  if (!response.ok) {
    throw new Error(`Request failed: ${response.status} ${response.statusText}`)
  }
  return (await response.json()) as Participant
}
