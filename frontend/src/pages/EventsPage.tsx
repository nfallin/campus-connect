import { useMemo, useState } from 'react'
import type { Event } from '../types'
import EventCard from '../components/EventCard'
import './EventsPage.css'

type RegistrationFilter = 'all' | 'registered' | 'not-registered'

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
  const [search, setSearch] = useState('')
  const [registrationFilter, setRegistrationFilter] =
    useState<RegistrationFilter>('all')

  const visibleEvents = useMemo(() => {
    const query = search.trim().toLocaleLowerCase()
    return events.filter((event) => {
      const matchesSearch = !query || [
        event.title,
        event.clubName,
        event.location,
      ].some((value) => value.toLocaleLowerCase().includes(query))
      const isRegistered = rsvpedEventIds.includes(event.id)
      const matchesFilter =
        registrationFilter === 'all' ||
        (registrationFilter === 'registered' && isRegistered) ||
        (registrationFilter === 'not-registered' && !isRegistered)

      return matchesSearch && matchesFilter
    })
  }, [events, rsvpedEventIds, search, registrationFilter])

  return (
    <main className="events-page">
      <div className="events-page-heading">
        <div>
          <h1 className="events-heading">Browse events</h1>
          <p className="events-count">
            {visibleEvents.length} {visibleEvents.length === 1 ? 'event' : 'events'}
          </p>
        </div>
      </div>

      <div className="events-controls">
        <label className="events-search">
          <span className="visually-hidden">Search events</span>
          <span className="events-search-icon" aria-hidden="true">⌕</span>
          <input
            type="search"
            placeholder="Search events"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
        </label>
        <label className="events-filter">
          <span className="visually-hidden">Filter events</span>
          <select
            value={registrationFilter}
            onChange={(event) =>
              setRegistrationFilter(event.target.value as RegistrationFilter)
            }
          >
            <option value="all">All events</option>
            <option value="registered">Registered</option>
            <option value="not-registered">Not registered</option>
          </select>
        </label>
      </div>

      {events.length === 0 ? (
        <p className="events-empty">No events yet.</p>
      ) : visibleEvents.length === 0 ? (
        <p className="events-empty">No events match your search or filter.</p>
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