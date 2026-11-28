
import type { KeyboardEvent } from 'react'
import type { Event } from '../types'
import { fromUtcDateTimeToLocalDisplayTime } from '../lib/eventDateTime'
import './EventCard.css'

type EventCardProps = {
  event: Event
  hasRSVPed: boolean
  isLoading: boolean
  disabled: boolean
  onRSVP: () => void
  onRemoveRSVP: () => void
  onOpen: () => void
}

export function toWebsiteHref(website: string) {
  return /^https?:\/\//i.test(website) ? website : `https://${website}`
}

export function formatDate(value: string) {
  const [year, month, day] = value.split('-')
  if (!year || !month || !day) return value
  return `${month}/${day}/${year}`
}

function EventCard({
  event,
  hasRSVPed,
  isLoading,
  disabled,
  onRSVP,
  onRemoveRSVP,
  onOpen,
}: EventCardProps) {
  function handleKeyDown(e: KeyboardEvent<HTMLLIElement>) {
    if (e.target !== e.currentTarget) return
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      onOpen()
    }
  }

  return (
    <li
      className="event-item event-item-clickable"
      role="button"
      tabIndex={0}
      aria-label={`Open details for ${event.title}`}
      onClick={onOpen}
      onKeyDown={handleKeyDown}
    >
      <div className="event-heading">
        <div className="event-title">{event.title}</div>
        <div className="event-club">
          {event.clubName || 'Club not specified'}
        </div>
        {event.interests.length > 0 && (
          <div className="event-interest-tags" aria-label="Event interests">
            {event.interests.map((interest) => (
              <span className="event-interest-tag" key={interest}>{interest}</span>
            ))}
          </div>
        )}
      </div>
      <div className="event-section">
        <span className="event-label">When</span>
        <p className="event-value">
          {formatDate(event.date)} &middot; {event.time}
          {event.endDateTime && ` - ${fromUtcDateTimeToLocalDisplayTime(event.endDateTime)}`}
        </p>
      </div>

      <div className="event-section">
        <span className="event-label">Where</span>
        <p className="event-value">{event.location}</p>
      </div>

      {event.description && (
        <div className="event-section">
          <span className="event-label">About</span>
          <p className="event-value event-description">{event.description}</p>
        </div>
      )}

      {(event.contactEmail || event.website) && (
        <div className="event-section">
          <span className="event-label">Contact</span>
          <div className="event-value event-contact">
            {event.contactEmail && (
              <a href={`mailto:${event.contactEmail}`} onClick={(e) => e.stopPropagation()}>
                {event.contactEmail}
              </a>
            )}

            {event.contactEmail && event.website && ' · '}

            {event.website && (
              <a
                href={toWebsiteHref(event.website)}
                target="_blank"
                rel="noreferrer"
                onClick={(e) => e.stopPropagation()}
              >
                {event.website}
              </a>
            )}
          </div>
        </div>
      )}

      <div className="event-section event-action">
        <button
          type="button"
          className="button button-primary"
          onClick={(e) => {
            e.stopPropagation()
            if (hasRSVPed) {
              onRemoveRSVP()
            } else {
              onRSVP()
            }
          }}
          disabled={disabled || event.pending || isLoading}
        >
          {isLoading
            ? 'Please wait...'
            : event.pending
              ? 'Pending approval'
              : disabled
              ? 'Sign in to register'
              : hasRSVPed
                ? 'Cancel registration'
                : 'Register'}
        </button>
      </div>
    </li>
  )
}

export default EventCard
