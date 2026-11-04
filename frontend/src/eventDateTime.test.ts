import { describe, it, expect } from 'vitest'
import { toUtcDateTime } from './eventDateTime'

describe('toUtcDateTime', () => {
  it('should preserve the selected date when scheduling an event at 8 PM', () => {
    const result = toUtcDateTime(
      '2026-11-10',
      '8',
      '00',
      'PM'
    )

    const expected = new Date(2026, 10, 10, 20, 0).toISOString()

    expect(result).toBe(expected)
  })

  it('should preserve the selected date when scheduling an event at noon', () => {
    const result = toUtcDateTime(
      '2026-11-10',
      '12',
      '00',
      'PM'
    )

    const expected = new Date(2026, 10, 10, 12, 0).toISOString()

    expect(result).toBe(expected)
  })

  it('should correctly handle midnight', () => {
    const result = toUtcDateTime(
      '2026-11-10',
      '12',
      '00',
      'AM'
    )

    const expected = new Date(2026, 10, 11, 0, 0).toISOString()

    expect(result).toBe(expected)
  })
})
