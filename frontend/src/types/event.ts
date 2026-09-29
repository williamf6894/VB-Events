export interface Participant {
  id: string
  name: string
  email: string
}

export interface Event {
  id: string
  name: string
  description: string
  location: string
  capacity: number
  duration: number
  startTimestamp: string
  participants?: Participant[]
}
