// Tests written by Kaavya Jethwa

import { describe, expect, it } from 'vitest'
import { isValidEventInput, type EventInput } from './eventValidation'

describe('event validation optional fields (Kaavya Jethwa)', () => {
  const validEvent: EventInput = {
    title: 'Campus Hackathon',
    date: '2026-11-15',
    hour: '10',
    minute: '30',
    location: 'Student Center',
    clubName: 'ACM',
  }

  it('rejects an event with a malformed contact email', () => {
    expect(isValidEventInput({ ...validEvent, contactEmail: 'not-an-email' })).toBe(false)
  })

  it('rejects an event with a malformed website', () => {
    expect(isValidEventInput({ ...validEvent, website: 'not a url' })).toBe(false)
  })
})
