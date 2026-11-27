
import useEventImage from './useEventImage'

type EventImageProps = {
  eventId: number
  hasImage: boolean
  className?: string
}

function EventImage({
  eventId,
  hasImage,
  className = '',
}: EventImageProps) {
  const imageUrl = useEventImage(eventId, hasImage)

  if (!hasImage || !imageUrl) {
    return null
  }

  return (
    <div className={`event-image-container ${className}`}>
      <img src={imageUrl} alt="Event" className="event-image" />
    </div>
  )
}

export default EventImage
