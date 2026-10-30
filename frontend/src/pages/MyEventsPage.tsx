import type { Event } from '../types'
import EventCard from '../components/EventCard'

type MyEventsPageProps = {
  events: Event[]
  rsvpLoading: number[]
  disabled: boolean
  onRemoveRSVP: (eventId: number) => void
}

function MyEventsPage({
  events,
  rsvpLoading,
  disabled,
  onRemoveRSVP,
}: MyEventsPageProps) {
  return (
    <main className="events-page">
      <h2 className="events-heading">My Events</h2>

      {disabled ? (
        <p>Please sign in to see your events.</p>
      ) : events.length === 0 ? (
        <p className="events-empty">
          You haven't RSVP'd to any events yet.
        </p>
      ) : (
        <ul className="events-list">
          {events.map((event) => (
            <EventCard
              key={event.id}
              event={event}
              hasRSVPed={true}
              isLoading={rsvpLoading.includes(event.id)}
              disabled={disabled}
              onRSVP={() => {}}
              onRemoveRSVP={() => onRemoveRSVP(event.id)}
            />
          ))}
        </ul>
      )}
    </main>
  )
}

export default MyEventsPage