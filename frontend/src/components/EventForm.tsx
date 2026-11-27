
import { useEffect, useState } from 'react'
import {
  getEventDateBounds,
  getEventInputErrorMap,
  MAX_DESCRIPTION_LENGTH,
  MAX_EVENT_FIELD_LENGTHS,
} from '../lib/eventValidation'
import './Form.css'

const HOURS = Array.from({ length: 12 }, (_, i) => String(i + 1))
const MINUTES = Array.from({ length: 60 }, (_, i) => String(i).padStart(2, '0'))

function getHour12(time: string): string {
  if (!time) return ''
  const hour24 = Number(time.split(':')[0])
  return String(hour24 % 12 || 12)
}

function getPeriod(time: string): 'AM' | 'PM' {
  if (!time) return 'AM'
  return Number(time.split(':')[0]) >= 12 ? 'PM' : 'AM'
}

function to24HourTime(hour: string, minute: string, period: 'AM' | 'PM'): string {
  if (!hour || !minute) return ''
  let hour24 = Number(hour) % 12
  if (period === 'PM') hour24 += 12
  return `${String(hour24).padStart(2, '0')}:${minute}`
}

export type EventFormData = {
  title: string
  date: string
  hour: string
  minute: string
  period: 'AM' | 'PM'
  endTime: string
  location: string
  clubName: string
  description: string
  contactEmail: string
  website: string
  imageFile?: File
  removeImage?: boolean
}

type ValidatedField = Exclude<
  keyof EventFormData,
  'imageFile' | 'removeImage'
>

type EventFormProps = {
  onSubmit: (data: EventFormData) => Promise<void>
  disabled?: boolean
  initialData?: EventFormData
  initialImageUrl?: string
  submitLabel?: string
  submittingLabel?: string
  onCancel?: () => void
}

function EventForm({
  onSubmit,
  disabled = false,
  initialData,
  initialImageUrl,
  submitLabel = 'Create event',
  submittingLabel = 'Creating...',
  onCancel,
}: EventFormProps) {
  const [title, setTitle] = useState(initialData?.title ?? '')
  const [date, setDate] = useState(initialData?.date ?? '')
  const [hour, setHour] = useState(initialData?.hour ?? '')
  const [minute, setMinute] = useState(initialData?.minute ?? '')
  const [period, setPeriod] = useState<'AM' | 'PM'>(initialData?.period ?? 'AM')
  const [endHour, setEndHour] = useState(getHour12(initialData?.endTime ?? ''))
  const [endMinute, setEndMinute] = useState(initialData?.endTime?.split(':')[1] ?? '')
  const [endPeriod, setEndPeriod] = useState<'AM' | 'PM'>(
    getPeriod(initialData?.endTime ?? ''),
  )
  const [location, setLocation] = useState(initialData?.location ?? '')
  const [clubName, setClubName] = useState(initialData?.clubName ?? '')
  const [description, setDescription] = useState(initialData?.description ?? '')
  const [contactEmail, setContactEmail] = useState(initialData?.contactEmail ?? '')
  const [website, setWebsite] = useState(initialData?.website ?? '')

  const [imageFile, setImageFile] = useState<File | undefined>()
  const [removeImage, setRemoveImage] = useState(false)
  const [imageError, setImageError] = useState('')
  const [imagePreview, setImagePreview] = useState('')

  useEffect(() => {
    if (!imageFile) {
      setImagePreview('')
      return
    }

    const previewUrl = URL.createObjectURL(imageFile)
    setImagePreview(previewUrl)

    return () => URL.revokeObjectURL(previewUrl)
  }, [imageFile])

  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState('')

  const [touched, setTouched] = useState<Partial<Record<ValidatedField, boolean>>>({})
  const [submitAttempted, setSubmitAttempted] = useState(false)

  const dateBounds = getEventDateBounds()

  const endTime =
    endHour || endMinute
      ? to24HourTime(endHour, endMinute, endPeriod) || 'incomplete'
      : ''

  const formData: EventFormData = {
    title,
    date,
    hour,
    minute,
    period,
    endTime,
    location,
    clubName,
    description,
    contactEmail,
    website,
    imageFile,
    removeImage,
  }

  const inputErrors = getEventInputErrorMap(formData)
  const canSubmit = !disabled && !submitting

  function getVisibleError(field: ValidatedField) {
    return submitAttempted || touched[field] ? inputErrors[field] : undefined
  }

  function markTouched(field: ValidatedField) {
    setTouched((prev) => ({ ...prev, [field]: true }))
  }

  function resetForm() {
    setTitle('')
    setDate('')
    setHour('')
    setMinute('')
    setPeriod('AM')
    setEndHour('')
    setEndMinute('')
    setEndPeriod('AM')
    setLocation('')
    setClubName('')
    setDescription('')
    setContactEmail('')
    setWebsite('')
    setImageFile(undefined)
    setRemoveImage(false)
    setImageError('')
  }

  function handleImageChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]

    if (!file) return

    if (!['image/jpeg', 'image/png'].includes(file.type)) {
      setImageError('Please select a JPEG or PNG image.')
      e.target.value = ''
      return
    }

    if (file.size > 5 * 1024 * 1024) {
      setImageError('Image must be 5 MB or smaller.')
      e.target.value = ''
      return
    }

    setImageError('')
    setImageFile(file)
    setRemoveImage(false)
  }

  async function handleSubmit(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault()
    setSubmitAttempted(true)

    if (!canSubmit || Object.keys(inputErrors).length > 0) return

    setSubmitting(true)
    setSubmitError('')

    try {
      await onSubmit(formData)
      resetForm()
      setTouched({})
      setSubmitAttempted(false)
    } catch (err) {
      console.error('Failed to create event', err)
      setSubmitError(
        err instanceof Error ? err.message : 'Failed to create event. Please try again.',
      )
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form className="card event-form" onSubmit={handleSubmit} noValidate>
      <div className={`field field-wide${getVisibleError('title') ? ' field-invalid' : ''}`}>
        <label htmlFor="title">Event title</label>
        <input
          id="title"
          type="text"
          placeholder="e.g. Fall Hackathon Kickoff"
          value={title}
          maxLength={MAX_EVENT_FIELD_LENGTHS.title}
          onChange={(e) => setTitle(e.target.value)}
          onBlur={() => markTouched('title')}
          aria-invalid={Boolean(getVisibleError('title'))}
          aria-describedby={getVisibleError('title') ? 'title-error' : undefined}
        />
        {getVisibleError('title') && (
          <p className="field-error" id="title-error">
            {getVisibleError('title')}
          </p>
        )}
      </div>

      <div className="event-date-times field-wide">
        <div className={`field${getVisibleError('date') ? ' field-invalid' : ''}`}>
          <label htmlFor="date">Date</label>
          <input
            id="date"
            type="date"
            min={dateBounds.min}
            max={dateBounds.max}
            value={date}
            onChange={(e) => setDate(e.target.value)}
            onBlur={() => markTouched('date')}
            aria-invalid={Boolean(getVisibleError('date'))}
            aria-describedby={getVisibleError('date') ? 'date-error' : undefined}
          />
          {getVisibleError('date') && (
            <p className="field-error" id="date-error">
              {getVisibleError('date')}
            </p>
          )}
        </div>

        <div className={`field${getVisibleError('hour') || getVisibleError('minute') ? ' field-invalid' : ''}`}>
          <label htmlFor="hour">Start time</label>
          <div className="time-row">
            <select
              id="hour"
              aria-label="Hour"
              value={hour}
              onChange={(e) => setHour(e.target.value)}
              onBlur={() => {
                markTouched('hour')
                markTouched('minute')
              }}
              aria-invalid={Boolean(getVisibleError('hour') || getVisibleError('minute'))}
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
              onBlur={() => {
                markTouched('hour')
                markTouched('minute')
              }}
              aria-invalid={Boolean(getVisibleError('hour') || getVisibleError('minute'))}
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
          {(getVisibleError('hour') || getVisibleError('minute')) && (
            <p className="field-error">
              {getVisibleError('hour') || getVisibleError('minute')}
            </p>
          )}
        </div>

        <div className={`field${getVisibleError('endTime') ? ' field-invalid' : ''}`}>
          <label htmlFor="endHour">End time (optional)</label>
          <div className="time-row">
            <select
              id="endHour"
              aria-label="End time hour"
              value={endHour}
              onChange={(e) => setEndHour(e.target.value)}
              onBlur={() => markTouched('endTime')}
              aria-invalid={Boolean(getVisibleError('endTime'))}
              aria-describedby={getVisibleError('endTime') ? 'endTime-error' : undefined}
            >
              <option value="">HH</option>
              {HOURS.map((h) => (
                <option key={h} value={h}>{h}</option>
              ))}
            </select>

            <span className="time-sep">:</span>

            <select
              aria-label="End time minute"
              value={endMinute}
              onChange={(e) => setEndMinute(e.target.value)}
              onBlur={() => markTouched('endTime')}
              aria-invalid={Boolean(getVisibleError('endTime'))}
              aria-describedby={getVisibleError('endTime') ? 'endTime-error' : undefined}
            >
              <option value="">MM</option>
              {MINUTES.map((m) => (
                <option key={m} value={m}>{m}</option>
              ))}
            </select>

            <div className="period-toggle">
              <button
                type="button"
                className={endPeriod === 'AM' ? 'active' : ''}
                onClick={() => setEndPeriod('AM')}
              >
                AM
              </button>
              <button
                type="button"
                className={endPeriod === 'PM' ? 'active' : ''}
                onClick={() => setEndPeriod('PM')}
              >
                PM
              </button>
            </div>
          </div>
          {getVisibleError('endTime') && (
            <p className="field-error" id="endTime-error">
              {getVisibleError('endTime')}
            </p>
          )}
        </div>
      </div>

      <div className={`field${getVisibleError('location') ? ' field-invalid' : ''}`}>
        <label htmlFor="location">Location</label>
        <input
          id="location"
          type="text"
          placeholder="e.g. Benson Center, Room 201"
          value={location}
          maxLength={MAX_EVENT_FIELD_LENGTHS.location}
          onChange={(e) => setLocation(e.target.value)}
          onBlur={() => markTouched('location')}
          aria-invalid={Boolean(getVisibleError('location'))}
          aria-describedby={getVisibleError('location') ? 'location-error' : undefined}
        />
        {getVisibleError('location') && (
          <p className="field-error" id="location-error">
            {getVisibleError('location')}
          </p>
        )}
      </div>

      <div className={`field${getVisibleError('clubName') ? ' field-invalid' : ''}`}>
        <label htmlFor="clubName">Club</label>
        <input
          id="clubName"
          type="text"
          placeholder="e.g. ACM"
          value={clubName}
          maxLength={MAX_EVENT_FIELD_LENGTHS.clubName}
          onChange={(e) => setClubName(e.target.value)}
          onBlur={() => markTouched('clubName')}
          aria-invalid={Boolean(getVisibleError('clubName'))}
          aria-describedby={getVisibleError('clubName') ? 'clubName-error' : undefined}
        />
        {getVisibleError('clubName') && (
          <p className="field-error" id="clubName-error">
            {getVisibleError('clubName')}
          </p>
        )}
      </div>

      <div className={`field field-wide${getVisibleError('description') ? ' field-invalid' : ''}`}>
        <label htmlFor="description">Description</label>
        <p className="field-hint">
          Mention any important requirements, like a dress code, what to bring, or age restrictions.
        </p>
        <textarea
          id="description"
          placeholder="Tell people what your event is about, and include any important requirements (e.g. dress code, what to bring)..."
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          maxLength={MAX_DESCRIPTION_LENGTH}
          rows={3}
          onBlur={() => markTouched('description')}
          aria-invalid={Boolean(getVisibleError('description'))}
          aria-describedby={`description-count${getVisibleError('description') ? ' description-error' : ''}`}
        />
        {getVisibleError('description') && (
          <p className="field-error" id="description-error">
            {getVisibleError('description')}
          </p>
        )}
        <p className="description-count" id="description-count" aria-live="polite">
          {description.length}/{MAX_DESCRIPTION_LENGTH} characters
        </p>
      </div>

      <div className={`field${getVisibleError('contactEmail') ? ' field-invalid' : ''}`}>
        <label htmlFor="contactEmail">Contact email (optional)</label>
        <input
          id="contactEmail"
          type="email"
          placeholder="e.g. acm@scu.edu"
          value={contactEmail}
          maxLength={MAX_EVENT_FIELD_LENGTHS.contactEmail}
          onChange={(e) => setContactEmail(e.target.value)}
          onBlur={() => markTouched('contactEmail')}
          aria-invalid={Boolean(getVisibleError('contactEmail'))}
          aria-describedby={getVisibleError('contactEmail') ? 'contactEmail-error' : undefined}
        />
        {getVisibleError('contactEmail') && (
          <p className="field-error" id="contactEmail-error">
            {getVisibleError('contactEmail')}
          </p>
        )}
      </div>

      <div className={`field${getVisibleError('website') ? ' field-invalid' : ''}`}>
        <label htmlFor="website">Website (optional)</label>
        <input
          id="website"
          type="text"
          placeholder="e.g. https://acm.scu.edu"
          value={website}
          maxLength={MAX_EVENT_FIELD_LENGTHS.website}
          onChange={(e) => setWebsite(e.target.value)}
          onBlur={() => markTouched('website')}
          aria-invalid={Boolean(getVisibleError('website'))}
          aria-describedby={getVisibleError('website') ? 'website-error' : undefined}
        />
        {getVisibleError('website') && (
          <p className="field-error" id="website-error">
            {getVisibleError('website')}
          </p>
        )}
      </div>

      <div className="field field-wide">
        <label htmlFor="eventImage">Event image (optional)</label>
        <p className="field-hint">
          Upload a JPEG or PNG image, up to 5 MB.
        </p>

        <input
          id="eventImage"
          type="file"
          accept="image/jpeg,image/png"
          onChange={handleImageChange}
          disabled={submitting}
        />

        {imagePreview && (
          <div className="event-image-preview">
            <img src={imagePreview} alt="Selected event image preview" />
          </div>
        )}

        {!imagePreview && initialImageUrl && !removeImage && (
          <div className="event-image-preview">
            <img src={initialImageUrl} alt="Current event image" />
          </div>
        )}

        {(imageFile || (initialImageUrl && !removeImage)) && (
          <button
            type="button"
            className="button button-secondary"
            onClick={() => {
              setImageFile(undefined)
              setRemoveImage(true)
              setImageError('')
              const input = document.getElementById('eventImage') as HTMLInputElement | null
              if (input) input.value = ''
            }}
            disabled={submitting}
          >
            Remove image
          </button>
        )}

        {imageError && (
          <p className="field-error" role="alert">{imageError}</p>
        )}
      </div>

      <div className="event-form-actions">
        <button type="submit" className="button button-primary create-btn" disabled={!canSubmit}>
          {submitting ? submittingLabel : submitLabel}
        </button>
        {onCancel && (
          <button
            type="button"
            className="button button-secondary event-form-cancel"
            onClick={onCancel}
            disabled={submitting}
          >
            Cancel
          </button>
        )}
      </div>

      <p className="login-note event-approval-note">
        Note: This event will be held for approval before it appears publicly.
      </p>

      {submitAttempted && Object.keys(inputErrors).length > 0 && (
        <p className="form-message field-wide" role="status">
          {Object.keys(inputErrors).length} fields need attention.
        </p>
      )}

      {disabled && (
        <p className="form-message field-wide" role="status">
          You must be signed in to create an event.
        </p>
      )}

      {submitError && (
        <p className="form-message field-wide" role="alert">{submitError}</p>
      )}
    </form>
  )
}

export default EventForm
