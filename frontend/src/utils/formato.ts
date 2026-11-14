// Formatos de presentación (es-MX).
const pesos = new Intl.NumberFormat('es-MX', { style: 'currency', currency: 'MXN' })
const enteros = new Intl.NumberFormat('es-MX')

export const moneda = (n: number) => pesos.format(n)
export const entero = (n: number) => enteros.format(n)
