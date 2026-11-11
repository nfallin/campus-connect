// Unit tests written by Kaavya Jethwa

import { describe, expect, it } from 'vitest'
import {
  getEventInputErrors,
  isValidEventInput,
  MAX_DESCRIPTION_LENGTH,
  MAX_EVENT_FIELD_LENGTHS,
  type EventInput,
} from './eventValidation'

describe('event validation gaps (Kaavya Jethwa)', () => {
  const today = new Date(2026, 9, 28)
  const validEvent: EventInput = {
    title: 'Campus Hackathon',
    date: '2026-11-15',
    hour: '10',
    minute: '30',
    period: 'AM',
    endTime: '',
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

  it('limits event descriptions to 500 characters', () => {
    expect(
      isValidEventInput({
        ...validEvent,
        description: 'a'.repeat(MAX_DESCRIPTION_LENGTH),
      }, today),
    ).toBe(true)
    expect(
      getEventInputErrors({
        ...validEvent,
        description: 'a'.repeat(MAX_DESCRIPTION_LENGTH + 1),
      }, today),
    ).toContain('Keep the description to 500 characters or fewer.')
  })

  it('enforces event field character limits', () => {
    expect(getEventInputErrors({
      ...validEvent,
      title: 'a'.repeat(MAX_EVENT_FIELD_LENGTHS.title + 1),
    }, today)).toContain('Keep the event title to 75 characters or fewer.')
    expect(getEventInputErrors({
      ...validEvent,
      clubName: 'a'.repeat(MAX_EVENT_FIELD_LENGTHS.clubName + 1),
    }, today)).toContain('Keep the organization name to 50 characters or fewer.')
    expect(getEventInputErrors({
      ...validEvent,
      location: 'a'.repeat(MAX_EVENT_FIELD_LENGTHS.location + 1),
    }, today)).toContain('Keep the location to 50 characters or fewer.')
    expect(getEventInputErrors({
      ...validEvent,
      contactEmail: `${'a'.repeat(MAX_EVENT_FIELD_LENGTHS.contactEmail - 1)}@b`,
    }, today)).toContain('Keep the email to 30 characters or fewer.')
    expect(getEventInputErrors({
      ...validEvent,
      website: `https://${'a'.repeat(MAX_EVENT_FIELD_LENGTHS.website - 7)}`,
    }, today)).toContain('Keep the website to 60 characters or fewer.')
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
