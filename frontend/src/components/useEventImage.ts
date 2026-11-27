
import { useEffect, useState } from 'react'
import { api } from '../lib/api'

export default function useEventImage(
  eventId: number,
  hasImage: boolean,
) {
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
        const blob = await api.getEventImage(eventId)
        const url = URL.createObjectURL(blob)

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

  return imageUrl
}
