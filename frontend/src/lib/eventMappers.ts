import type { Event, EventResponse } from '../types'
import { fromUtcDateTime } from './eventDateTime'

export function fromEventResponse(event: EventResponse): Event {
  const { date, time } = fromUtcDateTime(event.dateTime)

  return {
    id: event.id,
    title: event.name,
    dateTime: event.dateTime,
    date,
    time,
    location: event.location,
    clubName: event.clubName ?? '',
    contactEmail: event.contactEmail ?? '',
    description: event.description ?? '',
    website: event.websiteURL ?? '',
  }
}