import type { Event } from '../types'

export type EventSortField =
  | 'date'
  | 'title'
  | 'club'
  | 'location'
  | 'endDate'
  | 'approval'
  | 'rsvp'

export type EventCriteria = {
  search: string
  club: string
  location: string
  startDate: string
  endDate: string
  approval: 'all' | 'approved' | 'pending'
  rsvp: 'all' | 'registered' | 'notRegistered'
  sortField: EventSortField
  sortDirection: 'asc' | 'desc'
}

export const defaultEventCriteria: EventCriteria = {
  search: '',
  club: 'all',
  location: 'all',
  startDate: '',
  endDate: '',
  approval: 'all',
  rsvp: 'all',
  sortField: 'date',
  sortDirection: 'asc',
}

function normalized(value: string) {
  return value.trim().toLocaleLowerCase()
}

function compareText(a: string, b: string) {
  return a.localeCompare(b, undefined, { sensitivity: 'base', numeric: true })
}

export function filterAndSortEvents(
  events: Event[],
  criteria: EventCriteria,
  rsvpedEventIds: number[] = [],
): Event[] {
  const search = normalized(criteria.search)
  const rsvpIds = new Set(rsvpedEventIds)

  const filtered = events.filter((event) => {
    if (search) {
      const searchable = [
        event.title,
        event.clubName,
        event.location,
        event.description,
        event.contactEmail,
        event.website,
        event.date,
        event.time,
        event.endTime,
        event.pending ? 'pending' : 'approved',
      ].join(' ')

      if (!normalized(searchable).includes(search)) return false
    }

    if (
      criteria.club !== 'all'
      && (event.clubName || 'Club not specified') !== criteria.club
    ) return false
    if (criteria.location !== 'all' && event.location !== criteria.location) return false
    if (criteria.startDate && event.date < criteria.startDate) return false
    if (criteria.endDate && event.date > criteria.endDate) return false
    if (criteria.approval === 'approved' && event.pending) return false
    if (criteria.approval === 'pending' && !event.pending) return false

    const isRegistered = rsvpIds.has(event.id)
    if (criteria.rsvp === 'registered' && !isRegistered) return false
    if (criteria.rsvp === 'notRegistered' && isRegistered) return false

    return true
  })

  const direction = criteria.sortDirection === 'asc' ? 1 : -1

  return [...filtered].sort((a, b) => {
    let result = 0

    switch (criteria.sortField) {
      case 'title':
        result = compareText(a.title, b.title)
        break
      case 'club':
        result = compareText(a.clubName, b.clubName)
        break
      case 'location':
        result = compareText(a.location, b.location)
        break
      case 'endDate':
        result = (a.endDateTime ? Date.parse(a.endDateTime) : Number.POSITIVE_INFINITY)
          - (b.endDateTime ? Date.parse(b.endDateTime) : Number.POSITIVE_INFINITY)
        break
      case 'approval':
        result = Number(a.pending) - Number(b.pending)
        break
      case 'rsvp':
        result = Number(rsvpIds.has(a.id)) - Number(rsvpIds.has(b.id))
        break
      case 'date':
      default:
        result = (Date.parse(a.dateTime) || 0) - (Date.parse(b.dateTime) || 0)
        break
    }

    return result === 0 ? a.id - b.id : result * direction
  })
}
