// Unit tests written by Kaavya Jethwa

import { describe, expect, it } from 'vitest'
import { toFormData } from './EventDetails'
import type { Event } from '../types'

describe('EventDetails toFormData (Kaavya Jethwa)', () => {
  const baseEvent: Event = {
    id: 1,
    title: 'Hackathon',
    dateTime: '2026-01-15T00:00:00.000Z',
    date: '2026-01-15',
    time: '12:00 AM',
    location: 'Student Center',
    clubName: 'ACM',
    description: 'A day of coding.',
    contactEmail: 'acm@example.com',
    website: 'https://acm.example.com',
    hostUserId: 7,
  }

  it('maps midnight to 12 AM', () => {
    const midnight = new Date(2026, 0, 15, 0, 0)
    const event = { ...baseEvent, dateTime: midnight.toISOString() }

    const formData = toFormData(event)

    expect(formData.hour).toBe('12')
    expect(formData.period).toBe('AM')
  })

  it('maps noon to 12 PM', () => {
    const noon = new Date(2026, 0, 15, 12, 0)
    const event = { ...baseEvent, dateTime: noon.toISOString() }

    const formData = toFormData(event)

    expect(formData.hour).toBe('12')
    expect(formData.period).toBe('PM')
  })

  it('maps an afternoon hour to its 12-hour equivalent with PM', () => {
    const afternoon = new Date(2026, 0, 15, 14, 5)
    const event = { ...baseEvent, dateTime: afternoon.toISOString() }

    const formData = toFormData(event)

    expect(formData.hour).toBe('2')
    expect(formData.minute).toBe('05')
    expect(formData.period).toBe('PM')
  })

  it('carries over the event\'s other fields unchanged', () => {
    const formData = toFormData(baseEvent)

    expect(formData.title).toBe('Hackathon')
    expect(formData.date).toBe('2026-01-15')
    expect(formData.location).toBe('Student Center')
    expect(formData.clubName).toBe('ACM')
    expect(formData.description).toBe('A day of coding.')
    expect(formData.contactEmail).toBe('acm@example.com')
    expect(formData.website).toBe('https://acm.example.com')
  })
})
