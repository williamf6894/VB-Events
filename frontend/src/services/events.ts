import type { CreateEventPayload, Event, Participant } from '@/types/event'
import { authHeaders } from '@/services/api'

const API_BASE = '/api'

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, { headers: authHeaders() })
  if (!response.ok) {
    throw new Error(`Request failed: ${response.status} ${response.statusText}`)
  }
  return response.json() as Promise<T>
}

export interface EventListFilters {
  search?: string
  after?: string
  before?: string
  full?: boolean
}

export function getEvents(filters: EventListFilters = {}): Promise<Event[]> {
  const params = new URLSearchParams()
  if (filters.search) {
    params.set('q', filters.search)
  }
  if (filters.after) {
    params.set('after', filters.after)
  }
  if (filters.before) {
    params.set('before', filters.before)
  }
  if (filters.full !== undefined) {
    params.set('full', String(filters.full))
  }

  const queryString = params.toString()
  return request<Event[]>(`/events${queryString ? `?${queryString}` : ''}`)
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
  await expectOk(response, 'Failed to join event.')
}

export async function leaveEvent(id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/events/${id}/participants`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  await expectOk(response, 'Failed to leave event.')
}

export function getParticipants(): Promise<Participant[]> {
  return request<Participant[]>('/participants')
}

export async function inviteParticipant(eventId: string, participantId: string): Promise<void> {
  const response = await fetch(`${API_BASE}/events/${eventId}/participants/${participantId}`, {
    method: 'POST',
    headers: authHeaders(),
  })
  await expectOk(response, 'Failed to invite participant.')
}

async function expectOk(response: Response, fallbackMessage: string): Promise<void> {
  if (response.ok || response.status === 204) return

  let message = fallbackMessage
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
