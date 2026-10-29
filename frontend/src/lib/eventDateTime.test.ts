import { describe, expect, it } from 'vitest'
import { fromUtcDateTime, sortEventsChronologically, toUtcDateTime } from './eventDateTime'

describe('event date/time conversion', () => {
  it('converts the selected local date and time to a UTC timestamp', () => {
    expect(toUtcDateTime('2026-01-15', '10', '30', 'AM')).toBe(
      new Date(2026, 0, 15, 10, 30).toISOString(),
    )
  })

  it('uses the local date as well as local time when displaying a UTC timestamp', () => {
    const dateTime = new Date('2026-01-15T02:30:00.000Z')
    const localDate = [
      dateTime.getFullYear(),
      String(dateTime.getMonth() + 1).padStart(2, '0'),
      String(dateTime.getDate()).padStart(2, '0'),
    ].join('-')

    expect(fromUtcDateTime(dateTime.toISOString())).toEqual({
      date: localDate,
      time: dateTime.toLocaleTimeString([], { hour: 'numeric', minute: '2-digit' }),
    })
  })

  it('sorts events by UTC date and time, with ID as a deterministic tie-breaker', () => {
    const events = [
      { id: 3, dateTime: '2026-01-15T10:00:00.000Z' },
      { id: 2, dateTime: '2026-01-15T09:00:00.000Z' },
      { id: 1, dateTime: '2026-01-15T09:00:00.000Z' },
    ]

    expect(sortEventsChronologically(events).map((event) => event.id)).toEqual([1, 2, 3])
  })
})
