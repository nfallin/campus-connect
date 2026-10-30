import type { Event } from '../types'

type EventCardProps = {
  event: Event
}

function toWebsiteHref(website: string) {
  return /^https?:\/\//i.test(website) ? website : `https://${website}`
}

function formatDate(value: string) {
  const [year, month, day] = value.split('-')
  if (!year || !month || !day) return value
  return `${month}/${day}/${year}`
}

function EventCard({ event }: EventCardProps) {
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
    </li>
  )
}

export default EventCard