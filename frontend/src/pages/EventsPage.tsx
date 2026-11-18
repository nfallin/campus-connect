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

type EventsPageProps = {
  events: Event[]
  rsvpedEventIds: number[]
  rsvpLoading: number[]
  disabled: boolean
  onRSVP: (eventId: number) => void
  onRemoveRSVP: (eventId: number) => void
  onOpenEvent: (event: Event) => void
}

function EventsPage({
  events,
  rsvpedEventIds,
  rsvpLoading,
  disabled,
  onRSVP,
  onRemoveRSVP,
  onOpenEvent,
}: EventsPageProps) {
  const [criteria, setCriteria] = useState<EventCriteria>(defaultEventCriteria)
  const handleCriteriaChange = useCallback((next: EventCriteria) => setCriteria(next), [])
  const visibleEvents = useMemo(
    () => filterAndSortEvents(events, criteria, rsvpedEventIds),
    [events, criteria, rsvpedEventIds],
  )

  return (
    <main className="events-page">
      <h2 className="events-heading">Events</h2>

      <EventSearchControls
        events={events}
        onCriteriaChange={handleCriteriaChange}
        resultCount={visibleEvents.length}
      />

      {events.length === 0 ? (
        <p className="events-empty">No events yet.</p>
      ) : visibleEvents.length === 0 ? (
        <p className="events-empty">No events match your search and filters.</p>
      ) : (
        <ul className="events-list">
          {visibleEvents.map((event) => (
            <EventCard
              key={event.id}
              event={event}
              hasRSVPed={rsvpedEventIds.includes(event.id)}
              isLoading={rsvpLoading.includes(event.id)}
              disabled={disabled}
              onRSVP={() => onRSVP(event.id)}
              onRemoveRSVP={() => onRemoveRSVP(event.id)}
              onOpen={() => onOpenEvent(event)}
            />
          ))}
        </ul>
      )}
    </main>
  )
}

export default EventsPage