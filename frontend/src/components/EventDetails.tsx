import { useState } from 'react'
import EventForm, { type EventFormData } from './EventForm'
import type { Event } from '../types'
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
}: EventDetailsProps) {
  const [isEditing, setIsEditing] = useState(false)

  return (
    <div className="event-details-backdrop">
      <section
        className={`event-details${isEditing ? ' event-details-is-editing' : ''}`}
        role="dialog"
        aria-modal="true"
        aria-labelledby="event-details-title"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="event-details-header">
          {isEditing ? (
            <button
              type="button"
              className="event-details-back"
              onClick={() => setIsEditing(false)}
            >
              Back
            </button>
          ) : (
            <span />
          )}
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
              submitLabel="Save changes"
              submittingLabel="Saving..."
            />
          </div>
        ) : (
          <>
            <div className="event-details-content">
              <p className="event-details-club">{event.clubName || 'Club not specified'}</p>
              <p className="event-details-meta">
                {event.date} &middot; {event.time} &middot; {event.location}
              </p>
              <p className="event-details-description">
                {event.description || 'No description provided.'}
              </p>

              {(event.contactEmail || event.website) && (
                <div className="event-details-contact">
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

              <div className="event-details-actions">
                <button
                  type="button"
                  className="create-btn event-details-rsvp"
                  onClick={hasRSVPed ? onRemoveRSVP : onRSVP}
                  disabled={disabled || isLoading}
                >
                  {isLoading
                    ? 'Please wait...'
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
              </div>
            </div>
          </>
        )}
      </section>
    </div>
  )
}

export default EventDetails
