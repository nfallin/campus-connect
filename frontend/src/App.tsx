import { useEffect, useState } from 'react'
import { api } from './lib/api'
import { fromUtcDateTime, sortEventsChronologically, toUtcDateTime } from './lib/eventDateTime'
import './App.css'

import type { Event, EventResponse, User, UserResponse, Page } from './types'
import type { EventFormData } from './components/EventForm'

import LoginPage from './pages/LoginPage'
import EventsPage from './pages/EventsPage'
import CreateEventPage from './pages/CreateEventPage'


function getPageFromHash(): Page {
  switch (window.location.hash) {
    case '#login':
      return 'login'
    case '#events':
      return 'events'
    case '#create':
      return 'create'
    case '#rsvps':
      return 'rsvps'
    default:
      return 'login'
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
  const [page, setPage] = useState<Page>('login')
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [welcomeMessage, setWelcomeMessage] = useState('')

  useEffect(() => {
    window.location.hash = '#login'

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

    return () => window.removeEventListener('hashchange', handleHashChange)
  }, [])

  const handleLogin = async (name: string, email: string) => {
    const result: UserResponse = await api.post<UserResponse>('/users', {
      name,
      email,
    })

    setCurrentUser(result.user)
    setWelcomeMessage(
      result.isNewUser
        ? `Your account has been created, ${result.user.name}!`
        : `Welcome back, ${result.user.name}!`
    )
    window.location.hash = '#events'
  }

  const handleCreate = async (formData: EventFormData) => {
    if (!currentUser) {
      throw new Error('Please log in before creating an event.')
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
      hostUserId: currentUser.id,
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
      
      {page === 'login' ? (
        <LoginPage onLogin={handleLogin} />
      ) : page === 'create' ? (
        <CreateEventPage
          disabled={!currentUser}
          onSubmit={handleCreate}
        />
      ) : (
        <EventsPage
          welcomeMessage={welcomeMessage}
          events={events}
        />
      )}
    </div>
  )
}

export default App