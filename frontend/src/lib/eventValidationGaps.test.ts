// Unit tests written by Kaavya Jethwa

import { describe, expect, it } from 'vitest'
import { getEventInputErrors, isValidEventInput, type EventInput } from './eventValidation'

describe('event validation gaps (Kaavya Jethwa)', () => {
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

  it('rejects an event with an empty description', () => {
    expect(isValidEventInput({ ...validEvent, description: '' }, today)).toBe(false)
  })

  it('rejects an event when the description contains only whitespace', () => {
    expect(isValidEventInput({ ...validEvent, description: '   ' }, today)).toBe(false)
  })

  it('rejects an event with a malformed website', () => {
    expect(isValidEventInput({ ...validEvent, website: 'not a url' }, today)).toBe(false)
  })

  it('accepts an event with a well-formed website that omits the protocol', () => {
    expect(isValidEventInput({ ...validEvent, website: 'acm.scu.edu' }, today)).toBe(true)
  })

  it('accepts empty optional contact email and website fields', () => {
    expect(
      isValidEventInput({ ...validEvent, contactEmail: '', website: '' }, today),
    ).toBe(true)
  })

  it('reports every validation error at once rather than stopping at the first', () => {
    expect(
      getEventInputErrors(
        {
          ...validEvent,
          title: '',
          location: '',
          clubName: '',
          description: '',
          website: 'not a url',
        },
        today,
      ),
    ).toEqual([
      'Enter an event title.',
      'Enter a location.',
      'Enter a club name.',
      'Enter a description.',
      'Enter a valid website URL or leave it blank.',
    ])
  })
})
