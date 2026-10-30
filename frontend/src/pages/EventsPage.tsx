
import type { Event } from '../types'
import EventCard from '../components/EventCard'

type EventsPageProps = {
  events: Event[]
}

function EventsPage({ events }: EventsPageProps) {
  return (
    <main className="events-page">
      <h2 className="events-heading">Events</h2>

      {events.length === 0 ? (
        <p className="events-empty">No events yet.</p>
      ) : (
        <ul className="events-list">
          {events.map((event) => (
            <EventCard key={event.id} event={event} />
          ))}
        </ul>
      )}
    </main>
  )
}

export default EventsPage