
import { useEffect, useState } from 'react'
import { api } from '../lib/api'

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
  const [imageUrl, setImageUrl] = useState<string | null>(null)

  useEffect(() => {
    if (!hasImage) {
      setImageUrl(null)
      return
    }

    let objectUrl: string | null = null
    let cancelled = false

    async function loadImage() {
      try {
        const imageBlob = await api.getEventImage(eventId)
        const url = URL.createObjectURL(imageBlob)

        if (cancelled) {
          URL.revokeObjectURL(url)
          return
        }

        objectUrl = url
        setImageUrl(url)
      } catch (error) {
        console.error('Could not load event image:', error)
      }
    }

    setImageUrl(null)
    loadImage()

    return () => {
      cancelled = true
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl)
      }
    }
  }, [eventId, hasImage])

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
