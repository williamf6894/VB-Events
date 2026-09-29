import type { CreateEventPayload, Event } from '@/types/event'
import { authHeaders } from '@/services/api'

const API_BASE = '/api'

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, { headers: authHeaders() })
  if (!response.ok) {
    throw new Error(`Request failed: ${response.status} ${response.statusText}`)
  }
  return response.json() as Promise<T>
}

export function getUpcomingEvents(): Promise<Event[]> {
  const now = new Date().toISOString()
  return request<Event[]>(`/events/after?timestamp=${encodeURIComponent(now)}`)
}

export function getEvent(id: string): Promise<Event> {
  return request<Event>(`/events/${id}`)
}

export async function createEvent(payload: CreateEventPayload): Promise<Event> {
  const response = await fetch(`${API_BASE}/events`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    throw new Error(`Request failed: ${response.status} ${response.statusText}`)
  }
  return (await response.json()) as Event
}

export async function updateEvent(id: string, payload: CreateEventPayload): Promise<Event> {
  const response = await fetch(`${API_BASE}/events/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    throw new Error(`Request failed: ${response.status} ${response.statusText}`)
  }
  return (await response.json()) as Event
}

export async function deleteEvent(id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/events/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (!response.ok && response.status !== 204) {
    throw new Error(`Request failed: ${response.status} ${response.statusText}`)
  }
}

export async function joinEvent(id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/events/${id}/participants`, {
    method: 'POST',
    headers: authHeaders(),
  })
  if (!response.ok && response.status !== 204) {
    let message = 'Failed to join event.'
    try {
      const body = (await response.json()) as { message?: string }
      if (body.message) {
        message = body.message
      }
    } catch {
      // keep default message
    }
    throw new Error(message)
  }
}
