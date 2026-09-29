import type { CreateEventPayload, Event } from '@/types/event'

const API_BASE = '/api'

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`)
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
    headers: { 'Content-Type': 'application/json' },
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
    headers: { 'Content-Type': 'application/json' },
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
  })
  if (!response.ok && response.status !== 204) {
    throw new Error(`Request failed: ${response.status} ${response.statusText}`)
  }
}
