
import type { Event } from '../types'

type EventCardProps = {
  event: Event
  hasRSVPed: boolean
  isLoading: boolean
  disabled: boolean
  onRSVP: () => void
  onRemoveRSVP: () => void
}

function toWebsiteHref(website: string) {
  return /^https?:\/\//i.test(website) ? website : `https://${website}`
}

function formatDate(value: string) {
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
}: EventCardProps) {
  return (
    <li className="event-item">
      <div className="event-heading">
        <div className="event-title">{event.title}</div>
        <div className="event-club">
          {event.clubName || 'Club not specified'}
        </div>
      </div>

      <div className="event-meta">
        {formatDate(event.date)} &middot; {event.time} &middot; {event.location}
      </div>

      {event.description && (
        <p className="event-description">{event.description}</p>
      )}

      {(event.contactEmail || event.website) && (
        <div className="event-contact">
          {event.contactEmail && (
            <a href={`mailto:${event.contactEmail}`}>
              {event.contactEmail}
            </a>
          )}

          {event.contactEmail && event.website && ' · '}

          {event.website && (
            <a
              href={toWebsiteHref(event.website)}
              target="_blank"
              rel="noreferrer"
            >
              {event.website}
            </a>
          )}
        </div>
      )}

      <button
        type="button"
        onClick={hasRSVPed ? onRemoveRSVP : onRSVP}
        disabled={disabled || isLoading}
      >
        {isLoading
          ? 'Please wait...'
          : hasRSVPed
            ? 'Remove RSVP'
            : 'RSVP'}
      </button>
    </li>
  )
}

export default EventCard
