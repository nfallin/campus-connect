import { describe, expect, it } from 'vitest'
import type { Event } from '../types'
import { defaultEventCriteria, filterAndSortEvents } from './eventFilters'

function event(overrides: Partial<Event>): Event {
  return {
    id: 1,
    pending: false,
    title: 'Campus Event',
    dateTime: '2026-12-01T18:00:00.000Z',
    endDateTime: '2026-12-01T20:00:00.000Z',
    endTime: '8:00 PM',
    date: '2026-12-01',
    time: '6:00 PM',
    location: 'Student Center',
    clubName: 'Chess Club',
    description: 'Community gathering',
    contactEmail: 'chess@example.com',
    website: 'chess.example.com',
    hostUserId: 10,
    ...overrides,
  }
}

describe('filterAndSortEvents', () => {
  const events = [
    event({ id: 1, title: 'Beta', clubName: 'Science Club', date: '2026-12-02', dateTime: '2026-12-02T18:00:00.000Z' }),
    event({ id: 2, title: 'Alpha', clubName: 'Chess Club', date: '2026-12-01', dateTime: '2026-12-01T18:00:00.000Z', pending: true }),
  ]

  it('searches text fields and filters by club and date range', () => {
    const criteria = {
      ...defaultEventCriteria,
      search: 'chess',
      club: 'Chess Club',
      startDate: '2026-12-01',
      endDate: '2026-12-01',
    }

    expect(filterAndSortEvents(events, criteria).map(({ id }) => id)).toEqual([2])
  })

  it('filters pending and approved events', () => {
    expect(filterAndSortEvents(events, {
      ...defaultEventCriteria,
      approval: 'pending',
    }).map(({ id }) => id)).toEqual([2])

    expect(filterAndSortEvents(events, {
      ...defaultEventCriteria,
      approval: 'approved',
    }).map(({ id }) => id)).toEqual([1])
  })

  it('filters RSVP status', () => {
    expect(filterAndSortEvents(events, {
      ...defaultEventCriteria,
      rsvp: 'registered',
    }, [1]).map(({ id }) => id)).toEqual([1])

    expect(filterAndSortEvents(events, {
      ...defaultEventCriteria,
      rsvp: 'notRegistered',
    }, [1]).map(({ id }) => id)).toEqual([2])
  })

  it('sorts by title in both directions', () => {
    expect(filterAndSortEvents(events, {
      ...defaultEventCriteria,
      sortField: 'title',
      sortDirection: 'asc',
    }).map(({ id }) => id)).toEqual([2, 1])

    expect(filterAndSortEvents(events, {
      ...defaultEventCriteria,
      sortField: 'title',
      sortDirection: 'desc',
    }).map(({ id }) => id)).toEqual([1, 2])
  })

  it('supports unspecified club names as a filter value', () => {
    const withNoClub = event({ id: 3, clubName: '' })
    expect(filterAndSortEvents([withNoClub], {
      ...defaultEventCriteria,
      club: 'Club not specified',
    }).map(({ id }) => id)).toEqual([3])
  })
})
