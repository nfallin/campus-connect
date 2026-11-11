export function toUtcDateTime(
  date: string,
  hour: string,
  minute: string,
  period: 'AM' | 'PM',
): string {
  let hour24 = Number(hour) % 12
  if (period === 'PM') hour24 += 12
  return new Date(`${date}T${String(hour24).padStart(2, '0')}:${minute}:00`).toISOString()
}

export function toUtcDateTimeFromLocalTime(date: string, time: string): string {
  return new Date(`${date}T${time}:00`).toISOString()
}

export function fromUtcDateTimeToLocalTime(value?: string | null): string {
  if (!value) return ''
  const dateTime = new Date(value)
  return `${String(dateTime.getHours()).padStart(2, '0')}:${String(dateTime.getMinutes()).padStart(2, '0')}`
}

export function fromUtcDateTimeToLocalDisplayTime(value?: string | null): string {
  if (!value) return ''
  return new Date(value).toLocaleTimeString([], {
    hour: 'numeric',
    minute: '2-digit',
  })
}

export function fromUtcDateTime(value: string): { date: string; time: string } {
  const dateTime = new Date(value)
  const date = [
    dateTime.getFullYear(),
    String(dateTime.getMonth() + 1).padStart(2, '0'),
    String(dateTime.getDate()).padStart(2, '0'),
  ].join('-')

  return {
    date,
    time: fromUtcDateTimeToLocalDisplayTime(value),
  }
}

export function sortEventsChronologically<T extends { id: number; dateTime: string }>(
  events: T[],
): T[] {
  return [...events].sort((a, b) => {
    const timeDifference = Date.parse(a.dateTime) - Date.parse(b.dateTime)
    return timeDifference || a.id - b.id
  })
}
