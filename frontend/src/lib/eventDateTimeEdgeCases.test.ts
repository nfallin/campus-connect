// Unit tests written by Kaavya Jethwa

import { describe, expect, it } from 'vitest'
import { sortEventsChronologically, toUtcDateTime } from './eventDateTime'

describe('event date/time edge cases (Kaavya Jethwa)', () => {
  it('maps 12 AM to midnight (hour 0)', () => {
    expect(toUtcDateTime('2026-01-15', '12', '00', 'AM')).toBe(
      new Date(2026, 0, 15, 0, 0).toISOString(),
    )
  })

  it('maps 12 PM to noon (hour 12)', () => {
    expect(toUtcDateTime('2026-01-15', '12', '00', 'PM')).toBe(
      new Date(2026, 0, 15, 12, 0).toISOString(),
    )
  })

  it('returns an empty array when sorting no events', () => {
    expect(sortEventsChronologically([])).toEqual([])
  })

  it('leaves an already-sorted list unchanged', () => {
    const events = [
      { id: 1, dateTime: '2026-01-15T09:00:00.000Z' },
      { id: 2, dateTime: '2026-01-15T10:00:00.000Z' },
    ]

    expect(sortEventsChronologically(events).map((event) => event.id)).toEqual([1, 2])
  })
})
