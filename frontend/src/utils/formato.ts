// Formatos de presentación (es-MX).
const pesos = new Intl.NumberFormat('es-MX', { style: 'currency', currency: 'MXN' })
const enteros = new Intl.NumberFormat('es-MX')

export const moneda = (n: number) => pesos.format(n)
export const entero = (n: number) => enteros.format(n)

// "2026-11-14" → "14 nov 2026"
const fechas = new Intl.DateTimeFormat('es-MX', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'UTC' })
export const fecha = (iso: string) => fechas.format(new Date(`${iso}T00:00:00Z`))
