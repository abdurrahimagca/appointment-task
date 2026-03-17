const startOfLocalDay = (value: Date) =>
  new Date(value.getFullYear(), value.getMonth(), value.getDate())

const addLocalDays = (value: Date, amount: number) => {
  const nextValue = new Date(value)
  nextValue.setDate(nextValue.getDate() + amount)
  return nextValue
}

const formatLongDate = (value: Date) =>
  value.toLocaleDateString(undefined, {
    day: "numeric",
    month: "long",
    weekday: "long",
    year: "numeric",
  })

const formatMediumDate = (value: Date) =>
  value.toLocaleDateString(undefined, {
    day: "numeric",
    month: "long",
    weekday: "long",
  })

const formatShortDate = (value: Date) =>
  value.toLocaleDateString(undefined, {
    day: "numeric",
    month: "short",
    weekday: "short",
  })

const formatLocalTime = (value: string) =>
  new Date(value).toLocaleTimeString(undefined, {
    hour: "numeric",
    minute: "2-digit",
  })

const formatTimeInTimezone = (value: string, timezone: string) =>
  new Date(value).toLocaleTimeString(undefined, {
    hour: "numeric",
    minute: "2-digit",
    timeZone: timezone,
  })

const formatDateParam = (date: Date) => {
  const year = date.getFullYear()
  const month = `${date.getMonth() + 1}`.padStart(2, "0")
  const day = `${date.getDate()}`.padStart(2, "0")

  return `${year}-${month}-${day}`
}

const getLocalTimezone = () => Intl.DateTimeFormat().resolvedOptions().timeZone

export {
  addLocalDays,
  formatDateParam,
  formatLocalTime,
  formatLongDate,
  formatMediumDate,
  formatShortDate,
  formatTimeInTimezone,
  getLocalTimezone,
  startOfLocalDay,
}
