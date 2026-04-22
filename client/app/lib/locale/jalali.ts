export function formatJalaliDate(input: string | Date): string {
  const dateValue = typeof input === 'string' ? new Date(input) : input

  return new Intl.DateTimeFormat('fa-IR-u-ca-persian', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  }).format(dateValue)
}
