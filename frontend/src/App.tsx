import { useEffect, useState } from 'react'
import { api } from './lib/api'
import { fromUtcDateTime, sortEventsChronologically, toUtcDateTime } from './lib/eventDateTime'
import './App.css'

import type { Event, EventResponse, User, Page } from './types'
import type { EventFormData } from './components/EventForm'

import EventsPage from './pages/EventsPage'
import CreateEventPage from './pages/CreateEventPage'


function getPageFromHash(): Page {
  switch (window.location.hash) {
    case '#login':
      return 'login'
    case '#create':
      return 'create'
    case '#rsvps':
      return 'rsvps'
    default:
      return 'events'
  }
}

function fromEventResponse(event: EventResponse): Event {
  const { date, time } = fromUtcDateTime(event.dateTime)

  return {
    id: event.id,
    title: event.name,
    dateTime: event.dateTime,
    date,
    time,
    location: event.location,
    clubName: event.clubName ?? '',
    contactEmail: event.contactEmail ?? '',
    description: event.description ?? '',
    website: event.websiteURL ?? '',
  }
}

function App() {
  const [events, setEvents] = useState<Event[]>([])
  const [hostUserId, setHostUserId] = useState<number | null>(null)
  const [hostUserError, setHostUserError] = useState(false)
  const [page, setPage] = useState<Page>(getPageFromHash)

  useEffect(() => {
    const handleHashChange = () => setPage(getPageFromHash())

    window.addEventListener('hashchange', handleHashChange)

    api.get<{ events: EventResponse[] }>('/events').then(
      (data) => {
        setEvents(
          sortEventsChronologically(data.events.map(fromEventResponse))
        )
      },
      (err) => console.error('Failed to load events', err),
    )

    api.post<User>('/users', {
      name: 'Guest',
      email: 'guest@example.com',
    }).then(
      (user) => setHostUserId(user.id),
      (err) => {
        console.error('Failed to provision host user', err)
        setHostUserError(true)
      },
    )

    return () => window.removeEventListener('hashchange', handleHashChange)
  }, [])

  const handleCreate = async (formData: EventFormData) => {
    if (hostUserId === null) {
      throw new Error('No host user is available.')
    }

    const created = await api.post<EventResponse>('/events', {
      name: formData.title.trim(),
      dateTime: toUtcDateTime(
        formData.date,
        formData.hour,
        formData.minute,
        formData.period,
      ),
      location: formData.location.trim(),
      clubName: formData.clubName.trim(),
      description: formData.description.trim(),
      contactEmail: formData.contactEmail.trim(),
      websiteURL: formData.website.trim(),
      hostUserId,
    })

    const event = fromEventResponse(created)

    setEvents((prev) =>
      sortEventsChronologically([...prev, event])
    )

    window.location.hash = '#events'
  }

  return (
    <div className="page">
      <header className="page-header">
        <h1>Campus Connect</h1>
        <p className="subtitle">
          {page === 'create' ? 'Create an event' : 'Browse events'}
        </p>
      </header>

      <nav className="page-nav" aria-label="Main navigation">
        <a href="#login" aria-current={page === 'login' ? 'page' : undefined}>
          Login
        </a>
        <a href="#create" aria-current={page === 'create' ? 'page' : undefined}>
          Create event
        </a>
        <a href="#events" aria-current={page === 'events' ? 'page' : undefined}>
          View events
        </a>
        <a href="#rsvps" aria-current={page === 'rsvps' ? 'page' : undefined}>
          My events
        </a>
      </nav>

      {page === 'create' ? (
        <CreateEventPage
          onSubmit={handleCreate}
          disabled={hostUserId === null || hostUserError}
        />
      ) : (
        <EventsPage events={events} />
      )}
    </div>
  )
}

export default App