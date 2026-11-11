import type { Event, EventResponse } from '../types'
import { fromUtcDateTime, fromUtcDateTimeToLocalTime } from './eventDateTime'

export function fromEventResponse(event: EventResponse): Event {
  const { date, time } = fromUtcDateTime(event.dateTime)

  return {
    id: event.id,
    title: event.name,
    dateTime: event.dateTime,
    endDateTime: event.endDateTime,
    endTime: fromUtcDateTimeToLocalTime(event.endDateTime),
    date,
    time,
    location: event.location,
    clubName: event.clubName ?? '',
    contactEmail: event.contactEmail ?? '',
    description: event.description ?? '',
    website: event.websiteURL ?? '',
    hostUserId: event.hostUserId,
  }
}