import { useEffect, useState } from 'react'
import { api } from './lib/api'
import {
  sortEventsChronologically,
  toUtcDateTime,
  toUtcDateTimeFromLocalTime,
} from './lib/eventDateTime'
import { fromEventResponse } from './lib/eventMappers'
import './App.css'

import type { Event, EventResponse, User, UserResponse, RSVPResponse, Page } from './types'
import type { EventFormData } from './components/EventForm'

import AppHeader from './components/AppHeader'
import LoginPage from './pages/LoginPage'
import EventsPage from './pages/EventsPage'
import CreateEventPage from './pages/CreateEventPage'
import MyEventsPage from './pages/MyEventsPage'
import EventDetails from './components/EventDetails'


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
  
  const [page, setPage] = useState<Page>(getPageFromHash)

  const [currentUser, setCurrentUser] = useState<User | null>(null)

  const [welcomeMessage, setWelcomeMessage] = useState('')
  const [events, setEvents] = useState<Event[]>([])
  
  const [rsvpedEventIds, setRsvpedEventIds] = useState<number[]>([])
  const [rsvpLoading, setRsvpLoading] = useState<number[]>([])
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null)

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

    return () => window.removeEventListener('hashchange', handleHashChange)
  }, [])

  useEffect(() => {
    if (!currentUser && (page === 'create' || page === 'rsvps')) {
      window.location.hash = '#login'
    } else if (currentUser && page === 'login') {
      window.location.hash = '#events'
    }
  }, [currentUser, page])

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
      endDateTime: formData.endTime
        ? toUtcDateTimeFromLocalTime(formData.date, formData.endTime)
        : null,
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
      console.error('Failed to register for event', err)
      alert('Unable to register. Please try again.')
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
      console.error('Failed to cancel event registration', err)
      alert('Unable to cancel registration. Please try again.')
    } finally {
      setRsvpLoading((prev) =>
        prev.filter((id) => id !== eventId)
      )
    }
  }

  const handleUpdate = async (eventId: number, formData: EventFormData) => {
    const updated = await api.patch<EventResponse>(`/events/${eventId}`, {
      name: formData.title.trim(),
      dateTime: toUtcDateTime(
        formData.date,
        formData.hour,
        formData.minute,
        formData.period,
      ),
      endDateTime: formData.endTime
        ? toUtcDateTimeFromLocalTime(formData.date, formData.endTime)
        : '',
      location: formData.location.trim(),
      clubName: formData.clubName.trim(),
      description: formData.description.trim(),
      contactEmail: formData.contactEmail.trim(),
      websiteURL: formData.website.trim(),
    })
    const event = fromEventResponse(updated)
    setEvents((prev) =>
      sortEventsChronologically(prev.map((current) =>
        current.id === event.id ? event : current
      )),
    )
    setSelectedEvent(event)
  }

  const visiblePage =
    !currentUser && (page === 'create' || page === 'rsvps')
      ? 'login'
      : currentUser && page === 'login'
        ? 'events'
        : page

  return (
    <div className="page">
      <div className="page-top">
        <AppHeader
          page={visiblePage}
          currentUser={currentUser}
          welcomeMessage={welcomeMessage}
          onLogout={handleLogout}
        />

        <nav className="page-nav" aria-label="Main navigation">
          {!currentUser && (
            <a href="#login" aria-current={visiblePage === 'login' ? 'page' : undefined}>
              Login
            </a>
          )}
          {currentUser && (
            <a href="#create" aria-current={visiblePage === 'create' ? 'page' : undefined}>
              Create event
            </a>
          )}
          <a href="#events" aria-current={visiblePage === 'events' ? 'page' : undefined}>
            View events
          </a>
          {currentUser && (
            <a href="#rsvps" aria-current={visiblePage === 'rsvps' ? 'page' : undefined}>
              My events
            </a>
          )}
        </nav>
      </div>
      
      {visiblePage === 'login' ? (
        <LoginPage onLogin={handleLogin} />
      ) : visiblePage === 'create' ? (
        <CreateEventPage
          disabled={!currentUser}
          onSubmit={handleCreate}
        />
      ) : visiblePage === 'rsvps' ? (
        <MyEventsPage
          events={events.filter((event) =>
            rsvpedEventIds.includes(event.id)
          )}
          onRemoveRSVP={handleRemoveRSVP}
          onOpenEvent={setSelectedEvent}
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
          onOpenEvent={setSelectedEvent}
          disabled={!currentUser}
        />
      )}

      {selectedEvent && (
        <EventDetails
          event={selectedEvent}
          hasRSVPed={rsvpedEventIds.includes(selectedEvent.id)}
          isLoading={rsvpLoading.includes(selectedEvent.id)}
          disabled={!currentUser}
          isOwner={currentUser?.id === selectedEvent.hostUserId}
          onRSVP={() => handleRSVP(selectedEvent.id)}
          onRemoveRSVP={() => handleRemoveRSVP(selectedEvent.id)}
          onSave={(data) => handleUpdate(selectedEvent.id, data)}
          onClose={() => setSelectedEvent(null)}
        />
      )}
    </div>
  )
}

export default App