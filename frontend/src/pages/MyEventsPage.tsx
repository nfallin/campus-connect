import type { Event } from '../types'
import { useCallback, useMemo, useState } from 'react'
import EventCard from '../components/EventCard'
import EventSearchControls from '../components/EventSearchControls'
import {
  defaultEventCriteria,
  filterAndSortEvents,
  type EventCriteria,
} from '../lib/eventFilters'
import './EventsPage.css'

type MyEventsPageProps = {
  events: Event[]
  rsvpLoading: number[]
  disabled: boolean
  onRemoveRSVP: (eventId: number) => void
  onOpenEvent: (event: Event) => void
}

function MyEventsPage({
  events,
  rsvpLoading,
  disabled,
  onRemoveRSVP,
  onOpenEvent,
}: MyEventsPageProps) {
  const [criteria, setCriteria] = useState<EventCriteria>({
    ...defaultEventCriteria,
    rsvp: 'registered',
  })
  const handleCriteriaChange = useCallback((next: EventCriteria) => setCriteria(next), [])
  const registeredIds = useMemo(() => events.map((event) => event.id), [events])
  const visibleEvents = useMemo(
    () => filterAndSortEvents(events, criteria, registeredIds),
    [events, criteria, registeredIds],
  )

  return (
    <main className="events-page">
      <h2 className="events-heading">My Events</h2>

      {disabled ? (
        <p  className="events-empty">
          Please sign in to see your events.
        </p>
      ) : events.length === 0 ? (
        <p className="events-empty">
          You haven't registered for any events yet.
        </p>
      ) : (
        <>
          <EventSearchControls
            events={events}
            showRsvpFilter={false}
            onCriteriaChange={handleCriteriaChange}
            resultCount={visibleEvents.length}
          />
          {visibleEvents.length === 0 ? (
            <p className="events-empty">No events match your search and filters.</p>
          ) : (
            <ul className="events-list">
              {visibleEvents.map((event) => (
                <EventCard
                  key={event.id}
                  event={event}
                  hasRSVPed={true}
                  isLoading={rsvpLoading.includes(event.id)}
                  disabled={disabled}
                  onRSVP={() => {}}
                  onRemoveRSVP={() => onRemoveRSVP(event.id)}
                  onOpen={() => onOpenEvent(event)}
                />
              ))}
            </ul>
          )}
        </>
      )}
    </main>
  )
}

export default MyEventsPage