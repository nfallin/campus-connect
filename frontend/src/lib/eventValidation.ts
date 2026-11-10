export type EventInput = {
  title: string
  date: string
  hour: string
  minute: string
  period: 'AM' | 'PM'
  endTime: string
  location: string
  clubName: string
  description: string
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

export function getEventInputErrorMap(
  input: EventInput,
  now = new Date(),
): Partial<Record<keyof EventInput, string>> {
  const contactEmail = input.contactEmail?.trim() ?? ''
  const website = input.website?.trim() ?? ''
  const errors: Partial<Record<keyof EventInput, string>> = {}

  if (input.title.trim() === '') errors.title = 'Enter an event title.'

  if (input.date === '') {
    errors.date = 'Choose a date.'
  } 
  else if (!isWithinEventDateRange(input.date, now)) {
    const { min, max } = getEventDateBounds(now)
    if (/^\d{4}-\d{2}-\d{2}$/.test(input.date) && input.date < min) {
      errors.date = 'Choose today or a future date.'
    }
    else if (/^\d{4}-\d{2}-\d{2}$/.test(input.date) && input.date > max) {
      errors.date = 'Choose a date within the next six months.'
    }
    else {
      errors.date = 'Enter a valid event date.'
    }
  }

  if (input.hour === '' || input.minute === '') {
    errors.hour = 'Choose a start time.'
  } else if (input.date && isWithinEventDateRange(input.date, now)) {
    let startHour = Number(input.hour) % 12
    if (input.period === 'PM') startHour += 12
    const startDateTime = new Date(
      `${input.date}T${String(startHour).padStart(2, '0')}:${input.minute}:00`,
    )
    if (startDateTime <= now) {
      errors.hour = 'Choose a future start time.'
    }
  }

  if (input.endTime !== '') {
    const [endHour, endMinute] = input.endTime.split(':').map(Number)
    if (
      !/^\d{2}:\d{2}$/.test(input.endTime) ||
      endHour > 23 ||
      endMinute > 59
    ) {
      errors.endTime = 'Choose both an hour and minute for the end time.'
    } else if (input.hour !== '' && input.minute !== '') {
      let startHour = Number(input.hour) % 12
      if (input.period === 'PM') startHour += 12
      const startMinutes = startHour * 60 + Number(input.minute)
      const endMinutes = endHour * 60 + endMinute
      if (endMinutes - startMinutes < 30) {
        errors.endTime = 'Choose an end time at least 30 minutes after the start time.'
      }
    }
  }

  if (input.location.trim() === '') {
    errors.location = 'Enter a location.'
  }
  if (input.clubName.trim() === '') {
    errors.clubName = 'Enter a club name.'
  }
  if (input.description.trim() === '') {
    errors.description = 'Enter a description.'
  }
  if (contactEmail !== '' && !EMAIL_PATTERN.test(contactEmail)) {
    errors.contactEmail = 'Enter a valid contact email or leave it blank.'
  }
  if (website !== '' && !WEBSITE_PATTERN.test(website)) {
    errors.website = 'Enter a valid website URL or leave it blank.'
  }

  return errors
}

export function getEventInputErrors(input: EventInput, now = new Date()): string[] {
  return Object.values(getEventInputErrorMap(input, now)).filter(
    (error): error is string => error !== undefined,
  )
}

export function isValidEventInput(input: EventInput, now = new Date()): boolean {
  return (getEventInputErrors(input, now).length === 0)
}
