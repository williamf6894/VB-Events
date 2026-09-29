import type { Event } from '@/types/event'

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
