// Unit tests written by Kaavya Jethwa

import { describe, expect, it } from 'vitest'
import { formatDate, toWebsiteHref } from './EventCard'

describe('EventCard helpers (Kaavya Jethwa)', () => {
  it('leaves a website URL with an explicit protocol unchanged', () => {
    expect(toWebsiteHref('https://acm.scu.edu')).toBe('https://acm.scu.edu')
    expect(toWebsiteHref('http://acm.scu.edu')).toBe('http://acm.scu.edu')
  })

  it('prepends https:// to a bare domain', () => {
    expect(toWebsiteHref('acm.scu.edu')).toBe('https://acm.scu.edu')
  })

  it('formats an ISO date string as MM/DD/YYYY', () => {
    expect(formatDate('2026-11-15')).toBe('11/15/2026')
  })

  it('returns date input unchanged when it is missing a day or month', () => {
    expect(formatDate('2026-11')).toBe('2026-11')
  })
})
