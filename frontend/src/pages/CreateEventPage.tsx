import EventForm from '../components/EventForm'
import type { EventFormData } from '../components/EventForm'

type CreateEventPageProps = {
  onSubmit: (data: EventFormData) => Promise<void>
  disabled?: boolean
}

function CreateEventPage({ onSubmit, disabled = false }: CreateEventPageProps) {
  return (
    <EventForm
      onSubmit={onSubmit}
      disabled={disabled}
    />
  )
}

export default CreateEventPage