import { useEffect, useState } from 'react'
import { api } from './lib/api'
import { isValidEventInput } from './lib/eventValidation'
import './App.css'

const HOURS = Array.from({ length: 12 }, (_, i) => String(i + 1))
const MINUTES = Array.from({ length: 60 }, (_, i) => String(i).padStart(2, '0'))

type Event = {
  id: number
  title: string
  date: string
  time: string
  location: string
  contactEmail: string
  website: string
}

type EventResponse = {
  id: number
  name: string
  dateTime: string
  location: string
  contactEmail?: string
  websiteURL?: string
}

type UserResponse = {
  id: number
  name: string
}

function toWebsiteHref(website: string) {
  return /^https?:\/\//i.test(website) ? website : `https://${website}`
}

function formatDate(value: string) {
  const [year, month, day] = value.split('-')
  if (!year || !month || !day) return value
  return `${month}/${day}/${year}`
}

function fromEventResponse(event: EventResponse): Event {
  const dateTime = new Date(event.dateTime)
  return {
    id: event.id,
    title: event.name,
    date: dateTime.toISOString().slice(0, 10),
    time: dateTime.toLocaleTimeString([], { hour: 'numeric', minute: '2-digit' }),
    location: event.location,
    contactEmail: event.contactEmail ?? '',
    website: event.websiteURL ?? '',
  }
}

function toDateTime(date: string, hour: string, minute: string, period: 'AM' | 'PM') {
  let hour24 = Number(hour) % 12
  if (period === 'PM') hour24 += 12
  return new Date(`${date}T${String(hour24).padStart(2, '0')}:${minute}:00`).toISOString()
}

function App() {
  const [title, setTitle] = useState('')
  const [date, setDate] = useState('')
  const [hour, setHour] = useState('')
  const [minute, setMinute] = useState('')
  const [period, setPeriod] = useState<'AM' | 'PM'>('AM')
  const [location, setLocation] = useState('')
  const [clubName, setClubName] = useState('')
  const [contactEmail, setContactEmail] = useState('')
  const [website, setWebsite] = useState('')
  const [events, setEvents] = useState<Event[]>([])
  const [hostUserId, setHostUserId] = useState<number | null>(null)

  useEffect(() => {
    api.get<{ events: EventResponse[] }>('/events').then(
      (data) => setEvents(data.events.map(fromEventResponse)),
      (err) => console.error('Failed to load events', err),
    )

    api.post<UserResponse>('/users', { name: 'Guest' }).then(
      (user) => setHostUserId(user.id),
      (err) => console.error('Failed to provision host user', err),
    )
  }, [])

  const canCreate = isValidEventInput({
    title,
    date,
    hour,
    minute,
    location,
    clubName,
    contactEmail,
    website,
  })

  const handleCreate = () => {
    if (!canCreate || hostUserId === null) return

    api
      .post<EventResponse>('/events', {
        name: title.trim(),
        dateTime: toDateTime(date, hour, minute, period),
        location: location.trim(),
        clubName: clubName.trim(),
        contactEmail: contactEmail.trim(),
        websiteURL: website.trim(),
        hostUserId,
      })
      .then((created) => {
        const event = fromEventResponse(created)
        if (!event.contactEmail) event.contactEmail = contactEmail.trim()
        if (!event.website) event.website = website.trim()
        setEvents((prev) => [...prev, event])
        setTitle('')
        setDate('')
        setHour('')
        setMinute('')
        setPeriod('AM')
        setLocation('')
        setClubName('')
        setContactEmail('')
        setWebsite('')
      })
      .catch((err) => console.error('Failed to create event', err))
  }

  return (
    <div className="page">
      <header className="page-header">
        <h1>Campus Connect</h1>
        <p className="subtitle">Create an event</p>
      </header>

      <main className="card">
        <div className="field">
          <label htmlFor="title">Event title</label>
          <input
            id="title"
            type="text"
            placeholder="e.g. Fall Hackathon Kickoff"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
        </div>

        <div className="field">
          <label htmlFor="date">Date</label>
          <input
            id="date"
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
          />
        </div>

        <div className="field">
          <label htmlFor="hour">Time</label>
          <div className="time-row">
            <select
              id="hour"
              aria-label="Hour"
              value={hour}
              onChange={(e) => setHour(e.target.value)}
            >
              <option value="" disabled>
                HH
              </option>
              {HOURS.map((h) => (
                <option key={h} value={h}>
                  {h}
                </option>
              ))}
            </select>
            <span className="time-sep">:</span>
            <select
              aria-label="Minute"
              value={minute}
              onChange={(e) => setMinute(e.target.value)}
            >
              <option value="" disabled>
                MM
              </option>
              {MINUTES.map((m) => (
                <option key={m} value={m}>
                  {m}
                </option>
              ))}
            </select>
            <div className="period-toggle">
              <button
                type="button"
                className={period === 'AM' ? 'active' : ''}
                onClick={() => setPeriod('AM')}
              >
                AM
              </button>
              <button
                type="button"
                className={period === 'PM' ? 'active' : ''}
                onClick={() => setPeriod('PM')}
              >
                PM
              </button>
            </div>
          </div>
        </div>

        <div className="field">
          <label htmlFor="location">Location</label>
          <input
            id="location"
            type="text"
            placeholder="e.g. Benson Center, Room 201"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
          />
        </div>

        <div className="field">
          <label htmlFor="clubName">Club</label>
          <input
            id="clubName"
            type="text"
            placeholder="e.g. ACM"
            value={clubName}
            onChange={(e) => setClubName(e.target.value)}
          />
        </div>

        <div className="field">
          <label htmlFor="contactEmail">Contact email (optional)</label>
          <input
            id="contactEmail"
            type="email"
            placeholder="e.g. acm@scu.edu"
            value={contactEmail}
            onChange={(e) => setContactEmail(e.target.value)}
          />
        </div>

        <div className="field">
          <label htmlFor="website">Website (optional)</label>
          <input
            id="website"
            type="text"
            placeholder="e.g. https://acm.scu.edu"
            value={website}
            onChange={(e) => setWebsite(e.target.value)}
          />
        </div>

        <button
          type="button"
          className="create-btn"
          disabled={!canCreate || hostUserId === null}
          onClick={handleCreate}
        >
          Create event
        </button>
      </main>

      <section className="events">
        <h2 className="events-heading">Events</h2>
        {events.length === 0 ? (
          <p className="events-empty">No events yet.</p>
        ) : (
          <ul className="events-list">
            {events.map((event) => (
              <li key={event.id} className="event-item">
                <div className="event-title">{event.title}</div>
                <div className="event-meta">
                  {formatDate(event.date)} &middot; {event.time} &middot; {event.location}
                </div>
                {(event.contactEmail || event.website) && (
                  <div className="event-contact">
                    {event.contactEmail && (
                      <a href={`mailto:${event.contactEmail}`}>{event.contactEmail}</a>
                    )}
                    {event.contactEmail && event.website && ' · '}
                    {event.website && (
                      <a href={toWebsiteHref(event.website)} target="_blank" rel="noreferrer">
                        {event.website}
                      </a>
                    )}
                  </div>
                )}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}

export default App
