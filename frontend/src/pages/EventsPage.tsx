
import type { Event } from '../types'
import EventCard from '../components/EventCard'

type EventsPageProps = {
  events: Event[]
  welcomeMessage: string
}

function EventsPage({ events, welcomeMessage }: EventsPageProps) {
  return (
    <main className="events-page">
      <h2 className="events-heading">Events</h2>
      {welcomeMessage && (
        <div className="welcome-message" role="status">
          {welcomeMessage}
        </div>
      )}

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