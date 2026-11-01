//  Unit tests written by Sabrina Codipilly

import { describe, expect, it } from 'vitest'
import {
  getEventDateBounds,
  getEventInputErrors,
  isValidEventInput,
  type EventInput,
} from './eventValidation'

describe('event validation (Sabrina Codipilly)', () => {
  const today = new Date(2026, 9, 28)
  const validEvent: EventInput = {
    title: 'Campus Hackathon',
    date: '2026-11-15',
    hour: '10',
    minute: '30',
    location: 'Student Center',
    clubName: 'ACM',
    description: 'Join us for a day of coding and fun!',
  }

  it('accepts an event when every required field is provided', () => {
    expect(isValidEventInput(validEvent, today)).toBe(true)
  })

  it.each(['title', 'date', 'hour', 'minute', 'location', 'clubName'] as const)(
    'rejects an event with an empty %s',
    (field) => {
      expect(isValidEventInput({ ...validEvent, [field]: '' }, today)).toBe(false)
    },
  )

  it.each(['title', 'location', 'clubName'] as const)(
    'rejects an event when %s contains only whitespace',
    (field) => {
      expect(isValidEventInput({ ...validEvent, [field]: '   ' }, today)).toBe(false)
    },
  )

  it('allows dates from today through six calendar months from today', () => {
    const bounds = getEventDateBounds(today)

    expect(bounds).toEqual({ min: '2026-10-28', max: '2027-04-28' })
    expect(isValidEventInput({ ...validEvent, date: bounds.min }, today)).toBe(true)
    expect(isValidEventInput({ ...validEvent, date: bounds.max }, today)).toBe(true)
  })

  it('rejects past dates and dates beyond the six-month limit', () => {
    expect(isValidEventInput({ ...validEvent, date: '2026-10-27' }, today)).toBe(false)
    expect(isValidEventInput({ ...validEvent, date: '2027-04-29' }, today)).toBe(false)
  })

  it('clamps the six-month limit to the last day of the target month', () => {
    expect(getEventDateBounds(new Date(2026, 7, 31)).max).toBe('2027-02-28')
  })

  it('explains which fields need correction', () => {
    expect(
      getEventInputErrors(
        {
          ...validEvent,
          title: '',
          date: '2026-10-27',
          hour: '',
          contactEmail: 'not-an-email',
        },
        today,
      ),
    ).toEqual([
      'Enter an event title.',
      'Choose today or a future date.',
      'Choose an event time.',
      'Enter a valid contact email or leave it blank.',
    ])
  })
})
