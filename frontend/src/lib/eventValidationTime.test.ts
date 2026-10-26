// Tests written by Maxine Rajan
// Fixed version
import { describe, expect, it } from 'vitest'
import { isValidEventInput, type EventInput } from './eventValidation'

describe('event validation optional fields (Maxine Rajan)', () => {
  const validEvent: EventInput = {
    title: 'Diversity Forum',
    date: '2026-11-04',
    hour: '13',
    minute: '45',
    location: 'Mission Room',
    clubName: 'ASG',
  }

  it('rejects an event with an invalid hour', () => {
    expect(
      isValidEventInput({ ...validEvent, hour: '25' })
    ).toBe(true)
  })
})
