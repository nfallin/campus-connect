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

function formatLocalDate(date: Date): string {
  return [
    date.getFullYear(),
    String(date.getMonth() + 1).padStart(2, '0'),
    String(date.getDate()).padStart(2, '0'),
  ].join('-')
}

export function getEventDateBounds(now = new Date()): { min: string; max: string } {
  const minDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const targetMonth = minDate.getMonth() + 6
  const maxDate = new Date(minDate.getFullYear(), targetMonth, 1)
  const lastDayOfTargetMonth = new Date(
    maxDate.getFullYear(),
    maxDate.getMonth() + 1,
    0,
  ).getDate()
  maxDate.setDate(Math.min(minDate.getDate(), lastDayOfTargetMonth))

  return { min: formatLocalDate(minDate), max: formatLocalDate(maxDate) }
}

function isWithinEventDateRange(value: string, now: Date): boolean {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return false

  const [year, month, day] = value.split('-').map(Number)
  const parsedDate = new Date(year, month - 1, day)
  if (formatLocalDate(parsedDate) !== value) return false

  const { min, max } = getEventDateBounds(now)
  return value >= min && value <= max
}

export function getEventInputErrors(input: EventInput, now = new Date()): string[] {
  const contactEmail = input.contactEmail?.trim() ?? ''
  const website = input.website?.trim() ?? ''
  const errors: string[] = []

  if (input.title.trim() === '') errors.push('Enter an event title.')

  if (input.date === '') {
    errors.push('Choose an event date.')
  } else if (!isWithinEventDateRange(input.date, now)) {
    const { min, max } = getEventDateBounds(now)
    if (/^\d{4}-\d{2}-\d{2}$/.test(input.date) && input.date < min) {
      errors.push('Choose today or a future date.')
    } else if (/^\d{4}-\d{2}-\d{2}$/.test(input.date) && input.date > max) {
      errors.push('Choose a date within the next six months.')
    } else {
      errors.push('Enter a valid event date.')
    }
  }

  if (input.hour === '' || input.minute === '') errors.push('Choose an event time.')
  if (input.location.trim() === '') errors.push('Enter a location.')
  if (input.clubName.trim() === '') errors.push('Enter a club name.')
  if (contactEmail !== '' && !EMAIL_PATTERN.test(contactEmail)) {
    errors.push('Enter a valid contact email or leave it blank.')
  }
  if (website !== '' && !WEBSITE_PATTERN.test(website)) {
    errors.push('Enter a valid website URL or leave it blank.')
  }

  return errors
}

export function isValidEventInput(input: EventInput, now = new Date()): boolean {
  return getEventInputErrors(input, now).length === 0
}
