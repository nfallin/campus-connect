import { describe, expect, it } from 'vitest'
import { isValidEventInput, type EventInput } from './eventValidation'

describe('event validation (Sabrina Codipilly)', () => {
  const validEvent: EventInput = {
    title: 'Campus Hackathon',
    date: '2026-11-15',
    hour: '10',
    minute: '30',
    location: 'Student Center',
    clubName: 'ACM',
  }

  it('accepts an event when every required field is provided', () => {
    expect(isValidEventInput(validEvent)).toBe(true)
  })

  it.each(['title', 'date', 'hour', 'minute', 'location', 'clubName'] as const)(
    'rejects an event with an empty %s',
    (field) => {
      expect(isValidEventInput({ ...validEvent, [field]: '' })).toBe(false)
    },
  )

  it.each(['title', 'location', 'clubName'] as const)(
    'rejects an event when %s contains only whitespace',
    (field) => {
      expect(isValidEventInput({ ...validEvent, [field]: '   ' })).toBe(false)
    },
  )
})
