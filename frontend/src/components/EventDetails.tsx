import { useState } from 'react'
import EventForm, { type EventFormData } from './EventForm'
import type { Event } from '../types'
import {
  fromUtcDateTimeToLocalDisplayTime,
  fromUtcDateTimeToLocalTime,
} from '../lib/eventDateTime'
import './EventDetails.css'

type EventDetailsProps = {
  event: Event
  hasRSVPed: boolean
  isLoading: boolean
  disabled: boolean
  isOwner: boolean
  onRSVP: () => void
  onRemoveRSVP: () => void
  onSave: (data: EventFormData) => Promise<void>
  onClose: () => void
  onDownloadCalendar: () => void
}

export function toFormData(event: Event): EventFormData {
  const dateTime = new Date(event.dateTime)
  const hours = dateTime.getHours()
  return {
    title: event.title,
    date: event.date,
    hour: String(hours % 12 || 12),
    minute: String(dateTime.getMinutes()).padStart(2, '0'),
    period: hours >= 12 ? 'PM' : 'AM',
    endTime: fromUtcDateTimeToLocalTime(event.endDateTime),
    location: event.location,
    clubName: event.clubName,
    description: event.description,
    contactEmail: event.contactEmail,
    website: event.website,
  }
}

function toWebsiteHref(website: string) {
  return /^https?:\/\//i.test(website) ? website : `https://${website}`
}

function EventDetails({
  event,
  hasRSVPed,
  isLoading,
  disabled,
  isOwner,
  onRSVP,
  onRemoveRSVP,
  onSave,
  onClose,
  onDownloadCalendar,
}: EventDetailsProps) {
  const [isEditing, setIsEditing] = useState(false)

  return (
    <div className="event-details-backdrop" onClick={onClose}>
      <section
        className={`event-details${isEditing ? ' event-details-is-editing' : ''}`}
        role="dialog"
        aria-modal="true"
        aria-labelledby="event-details-title"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="event-details-header">
          <span />
          <h2 id="event-details-title">{isEditing ? 'Edit event' : event.title}</h2>
          <button type="button" className="event-details-close" onClick={onClose}>
            Close
          </button>
        </div>

        {isEditing ? (
          <div className="event-details-editor-view">
            <EventForm
              key={event.id + event.dateTime}
              initialData={toFormData(event)}
              onSubmit={async (data) => {
                await onSave(data)
                setIsEditing(false)
              }}
              onCancel={() => setIsEditing(false)}
              submitLabel="Save changes"
              submittingLabel="Saving..."
            />
          </div>
        ) : (
          <>
            <div className="event-details-content">
              <p className="event-details-club">{event.clubName || 'Club not specified'}</p>

              <div className="event-details-section">
                <span className="event-details-label">When</span>
                <p className="event-details-value">
                  {event.date} &middot; {event.time}
                  {event.endDateTime && ` - ${fromUtcDateTimeToLocalDisplayTime(event.endDateTime)}`}
                </p>
              </div>

              <div className="event-details-section">
                <span className="event-details-label">Where</span>
                <p className="event-details-value">{event.location}</p>
              </div>

              <div className="event-details-section">
                <span className="event-details-label">About</span>
                <p className="event-details-value event-details-description">
                  {event.description || 'No description provided.'}
                </p>
              </div>

              {(event.contactEmail || event.website) && (
                <div className="event-details-section">
                  <span className="event-details-label">Contact</span>
                  <div className="event-details-value event-details-contact">
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
                </div>
              )}
            </div>
            <div className="event-details-actions">
              <button
                type="button"
                className="create-btn event-details-rsvp"
                onClick={hasRSVPed ? onRemoveRSVP : onRSVP}
                disabled={disabled || event.pending || isLoading}
              >
                {isLoading
                  ? 'Please wait...'
                  : event.pending
                    ? 'Pending approval'
                    : disabled
                    ? 'Sign in to register'
                    : hasRSVPed
                      ? 'Cancel registration'
                      : 'Register'}
              </button>
              {isOwner && (
                <button
                  type="button"
                  className="event-details-edit-toggle"
                  onClick={() => setIsEditing(true)}
                >
                  Edit event
                </button>
              )}
              {!event.pending && (
                <button
                  type="button"
                  className="event-details-calendar-button"
                  title="Downloads an iCalendar (.ics) file. Open it with your calendar app to add this event."
                  onClick={onDownloadCalendar}
                >
                  Add to calendar
                </button>
              )}
            </div>
          </>
        )}
      </section>
    </div>
  )
}

export default EventDetails
