import { useState, type FormEvent } from 'react'
import { getEventDateBounds, getEventInputErrors } from '../lib/eventValidation'

const HOURS = Array.from({ length: 12 }, (_, i) => String(i + 1))
const MINUTES = Array.from({ length: 60 }, (_, i) => String(i).padStart(2, '0'))

export type EventFormData = {
  title: string
  date: string
  hour: string
  minute: string
  period: 'AM' | 'PM'
  location: string
  clubName: string
  description: string
  contactEmail: string
  website: string
}

type EventFormProps = {
  onSubmit: (data: EventFormData) => Promise<void>
  disabled?: boolean
}

function EventForm({ onSubmit, disabled = false }: EventFormProps) {
  const [title, setTitle] = useState('')
  const [date, setDate] = useState('')
  const [hour, setHour] = useState('')
  const [minute, setMinute] = useState('')
  const [period, setPeriod] = useState<'AM' | 'PM'>('AM')
  const [location, setLocation] = useState('')
  const [clubName, setClubName] = useState('')
  const [description, setDescription] = useState('')
  const [contactEmail, setContactEmail] = useState('')
  const [website, setWebsite] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState('')

  const dateBounds = getEventDateBounds()

  const formData: EventFormData = {
    title,
    date,
    hour,
    minute,
    period,
    location,
    clubName,
    description,
    contactEmail,
    website,
  }

  const inputErrors = getEventInputErrors(formData)
  const canSubmit = inputErrors.length === 0 && !disabled && !submitting

  function resetForm() {
    setTitle('')
    setDate('')
    setHour('')
    setMinute('')
    setPeriod('AM')
    setLocation('')
    setClubName('')
    setDescription('')
    setContactEmail('')
    setWebsite('')
  }

async function handleSubmit(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault()

    if (!canSubmit) return

    setSubmitting(true)
    setSubmitError('')

    try {
      await onSubmit(formData)
      resetForm()
    } catch (err) {
      console.error('Failed to create event', err)
      setSubmitError('Failed to create event. Please try again.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form className="card" onSubmit={handleSubmit}>
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
          min={dateBounds.min}
          max={dateBounds.max}
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
            <option value="" disabled>HH</option>
            {HOURS.map((h) => (
              <option key={h} value={h}>{h}</option>
            ))}
          </select>

          <span className="time-sep">:</span>

          <select
            aria-label="Minute"
            value={minute}
            onChange={(e) => setMinute(e.target.value)}
          >
            <option value="" disabled>MM</option>
            {MINUTES.map((m) => (
              <option key={m} value={m}>{m}</option>
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
        <label htmlFor="description">Description</label>
        <textarea
          id="description"
          placeholder="Tell people what your event is about..."
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={4}
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

      <button type="submit" className="create-btn" disabled={!canSubmit}>
        {submitting ? 'Creating...' : 'Create event'}
      </button>
      
      {(inputErrors.length > 0 || disabled) && (
        <div className="form-message" role="status" aria-live="polite">
          <p>Please meet the following requirements:</p>
          <ul>
            {disabled && <li>You must be signed in to create an event.</li>}
            {inputErrors.map((error) => (
              <li key={error}>{error}</li>
            ))}
          </ul>
        </div>
      )}

      {submitError && (
        <p className="form-message" role="alert">{submitError}</p>
      )}
    </form>
  )
}

export default EventForm