export type EventInput = {
  title: string
  date: string
  hour: string
  minute: string
  location: string
  clubName: string
  contactEmail?: string
  website?: string
}

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const WEBSITE_PATTERN = /^(https?:\/\/)?[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+(\/\S*)?$/

export function isValidEventInput(input: EventInput): boolean {
  const contactEmail = input.contactEmail?.trim() ?? ''
  const website = input.website?.trim() ?? ''

  return (
    input.title.trim() !== '' &&
    input.date !== '' &&
    input.hour !== '' &&
    input.minute !== '' &&
    input.location.trim() !== '' &&
    input.clubName.trim() !== '' &&
    (contactEmail === '' || EMAIL_PATTERN.test(contactEmail)) &&
    (website === '' || WEBSITE_PATTERN.test(website))
  )
}
