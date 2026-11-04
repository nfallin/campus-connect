import { describe, it, expect } from 'vitest'
import {
  toUtcDateTime,
  fromUtcDateTime,
} from './eventDateTime'

describe('Event date and time conversion', () => {
  it('preserves the selected date when scheduling an event at 8:00 PM', () => {
    const selectedDate = '2026-11-10'

    const utcDateTime = toUtcDateTime(
      selectedDate,
      '8',
      '00',
      'PM',
    )

    const result = fromUtcDateTime(utcDateTime)

    expect(result.date).toBe(selectedDate)
    expect(result.time).toMatch(/8:00\s?PM/i)
  })

  it('preserves the selected date when scheduling an event at 11:59 PM', () => {
    const selectedDate = '2026-11-10'

    const utcDateTime = toUtcDateTime(
      selectedDate,
      '11',
      '59',
      'PM',
    )

    const result = fromUtcDateTime(utcDateTime)

    expect(result.date).toBe(selectedDate)
    expect(result.time).toMatch(/11:59\s?PM/i)
  })

  it('correctly handles midnight (12:00 AM)', () => {
    const selectedDate = '2026-11-10'

    const utcDateTime = toUtcDateTime(
      selectedDate,
      '12',
      '00',
      'AM',
    )

    const result = fromUtcDateTime(utcDateTime)

    expect(result.date).toBe(selectedDate)
    expect(result.time).toMatch(/12:00\s?AM/i)
  })

  it('correctly handles noon (12:00 PM)', () => {
    const selectedDate = '2026-11-10'

    const utcDateTime = toUtcDateTime(
      selectedDate,
      '12',
      '00',
      'PM',
    )

    const result = fromUtcDateTime(utcDateTime)

    expect(result.date).toBe(selectedDate)
    expect(result.time).toMatch(/12:00\s?PM/i)
  })
})
