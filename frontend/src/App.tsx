import { useEffect, useState } from 'react'
import { api } from './lib/api'
import { sortEventsChronologically, toUtcDateTime } from './lib/eventDateTime'
import { fromEventResponse } from './lib/eventMappers'
import './App.css'

import type { Event, EventResponse, User, UserResponse, RSVPResponse, Page } from './types'
import type { EventFormData } from './components/EventForm'

import AppHeader from './components/AppHeader'
import LoginPage from './pages/LoginPage'
import EventsPage from './pages/EventsPage'
import CreateEventPage from './pages/CreateEventPage'
import MyEventsPage from './pages/MyEventsPage'


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

function App() {
  
  const [page, setPage] = useState<Page>('login')

  const [currentUser, setCurrentUser] = useState<User | null>(null)

  const [welcomeMessage, setWelcomeMessage] = useState('')
  const [events, setEvents] = useState<Event[]>([])
  
  const [rsvpedEventIds, setRsvpedEventIds] = useState<number[]>([])
  const [rsvpLoading, setRsvpLoading] = useState<number[]>([])

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
    const result = await api.post<UserResponse>('/users', {
      name,
      email,
    })

    const rsvpData = await api.get<RSVPResponse>(
      `/users/${result.user.id}/rsvps`
    )

    setRsvpedEventIds(
      rsvpData.events.map((event) => event.id)
    )

    setCurrentUser(result.user)

    setWelcomeMessage(
      result.isNewUser
        ? `Your account has been created, ${result.user.name}!`
        : `Welcome back, ${result.user.name}!`
    )

    window.location.hash = '#events'
  }

  const handleLogout = () => {
    setCurrentUser(null)
    setRsvpedEventIds([])
    setRsvpLoading([])
    setWelcomeMessage('')

    window.location.hash = '#login'
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


  const handleRSVP = async (eventId: number) => {
    if (!currentUser || rsvpLoading.includes(eventId)) return

    setRsvpLoading((prev) => [...prev, eventId])

    try {
      await api.post(
        `/users/${currentUser.id}/rsvps/${eventId}`
      )

      setRsvpedEventIds((prev) =>
        prev.includes(eventId) ? prev : [...prev, eventId]
      )
    } catch (err) {
      console.error('Failed to RSVP', err)
      alert('Unable to RSVP. Please try again.')
    } finally {
      setRsvpLoading((prev) =>
        prev.filter((id) => id !== eventId)
      )
    }
  }

  const handleRemoveRSVP = async (eventId: number) => {
    if (!currentUser || rsvpLoading.includes(eventId)) return

    setRsvpLoading((prev) => [...prev, eventId])

    try {
      await api.delete(
        `/users/${currentUser.id}/rsvps/${eventId}`
      )

      setRsvpedEventIds((prev) =>
        prev.filter((id) => id !== eventId)
      )
    } catch (err) {
      console.error('Failed to remove RSVP', err)
      alert('Unable to remove RSVP. Please try again.')
    } finally {
      setRsvpLoading((prev) =>
        prev.filter((id) => id !== eventId)
      )
    }
  }

  return (
    <div className="page">
      <AppHeader
        page={page}
        currentUser={currentUser}
        welcomeMessage={welcomeMessage}
        onLogout={handleLogout}
      />

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
      ) : page === 'rsvps' ? (
        <MyEventsPage
          events={events.filter((event) =>
            rsvpedEventIds.includes(event.id)
          )}
          onRemoveRSVP={handleRemoveRSVP}
          rsvpLoading={rsvpLoading}
          disabled={!currentUser}
        />
      ) : (
        <EventsPage
          events={events}
          rsvpedEventIds={rsvpedEventIds}
          rsvpLoading={rsvpLoading}
          onRSVP={handleRSVP}
          onRemoveRSVP={handleRemoveRSVP}
          disabled={!currentUser}
        />
      )}
    </div>
  )
}

export default App