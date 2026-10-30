import type { Event } from '../types'
import EventCard from '../components/EventCard'

type EventsPageProps = {
  events: Event[]
  rsvpedEventIds: number[]
  rsvpLoading: number[]
  disabled: boolean
  onRSVP: (eventId: number) => void
  onRemoveRSVP: (eventId: number) => void
}

function EventsPage({
  events,
  rsvpedEventIds,
  rsvpLoading,
  disabled,
  onRSVP,
  onRemoveRSVP,
}: EventsPageProps) {
  return (
    <main className="events-page">
      <h2 className="events-heading">Events</h2>

      {events.length === 0 ? (
        <p className="events-empty">No events yet.</p>
      ) : (
        <ul className="events-list">
          {events.map((event) => (
            <EventCard
              key={event.id}
              event={event}
              hasRSVPed={rsvpedEventIds.includes(event.id)}
              isLoading={rsvpLoading.includes(event.id)}
              disabled={disabled}
              onRSVP={() => onRSVP(event.id)}
              onRemoveRSVP={() => onRemoveRSVP(event.id)}
            />
          ))}
        </ul>
      )}
    </main>
  )
}

export default EventsPage