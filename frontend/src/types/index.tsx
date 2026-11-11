export type User = {
  id: number
  name: string
  email: string
}

export type UserResponse = {
  user: User
  isNewUser: boolean
}

export type Event = {
  id: number
  title: string
  dateTime: string
  endDateTime?: string | null
  endTime?: string
  date: string
  time: string
  location: string
  clubName: string
  description: string
  contactEmail: string
  website: string
  hostUserId: number
}

export type EventResponse = {
  id: number
  name: string
  dateTime: string
  endDateTime?: string | null
  location: string
  clubName: string
  description: string
  host?: User
  contactEmail: string
  websiteURL: string
  attendeeCount: number
  hostUserId: number
}

export type RSVPResponse = {
  events: EventResponse[]
}

export type Page = 'login' | 'create' | 'events' | 'rsvps'