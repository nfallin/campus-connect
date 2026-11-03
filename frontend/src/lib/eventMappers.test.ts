// Unit tests written by Kaavya Jethwa

import { describe, expect, it } from 'vitest'
import { fromEventResponse } from './eventMappers'
import type { EventResponse } from '../types'

describe('fromEventResponse (Kaavya Jethwa)', () => {
  const baseResponse: EventResponse = {
    id: 1,
    name: 'Hackathon',
    dateTime: '2026-01-15T18:30:00.000Z',
    location: 'Student Center',
    clubName: 'ACM',
    description: 'A day of coding.',
    contactEmail: 'acm@example.com',
    websiteURL: 'https://acm.example.com',
    attendeeCount: 3,
    hostUserId: 7,
  }

  it('maps every field from the API response onto the frontend Event shape', () => {
    const event = fromEventResponse(baseResponse)

    expect(event.id).toBe(1)
    expect(event.title).toBe('Hackathon')
    expect(event.dateTime).toBe('2026-01-15T18:30:00.000Z')
    expect(event.location).toBe('Student Center')
    expect(event.clubName).toBe('ACM')
    expect(event.description).toBe('A day of coding.')
    expect(event.contactEmail).toBe('acm@example.com')
    expect(event.website).toBe('https://acm.example.com')
    expect(event.hostUserId).toBe(7)
  })

  it('derives date and time from the dateTime field', () => {
    const event = fromEventResponse(baseResponse)
    const expectedDateTime = new Date(baseResponse.dateTime)

    expect(event.date).toBe(
      [
        expectedDateTime.getFullYear(),
        String(expectedDateTime.getMonth() + 1).padStart(2, '0'),
        String(expectedDateTime.getDate()).padStart(2, '0'),
      ].join('-'),
    )
    expect(event.time).toBe(
      expectedDateTime.toLocaleTimeString([], { hour: 'numeric', minute: '2-digit' }),
    )
  })

  it('defaults optional fields to empty strings when the backend omits them', () => {
    const sparseResponse = {
      ...baseResponse,
      clubName: undefined,
      contactEmail: undefined,
      description: undefined,
      websiteURL: undefined,
    } as unknown as EventResponse

    const event = fromEventResponse(sparseResponse)

    expect(event.clubName).toBe('')
    expect(event.contactEmail).toBe('')
    expect(event.description).toBe('')
    expect(event.website).toBe('')
  })
})
