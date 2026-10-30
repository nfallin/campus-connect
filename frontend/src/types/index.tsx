export type User = {
  id: number
  name: string
  email: string
}

export type Event = {
  id: number
  title: string
  dateTime: string
  date: string
  time: string
  location: string
  clubName: string
  description: string
  contactEmail: string
  website: string
}

export type EventResponse = {
  id: number
  name: string
  dateTime: string
  location: string
  clubName: string
  description: string
  host?: User
  contactEmail: string
  websiteURL: string
  attendeeCount: number
}

export type Page = 'login' | 'create' | 'events' | 'rsvps'