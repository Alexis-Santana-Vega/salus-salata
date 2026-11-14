/**
 * Reglas de validación para `:rules` de Vuetify (sin librerías externas).
 * Cada regla devuelve true o el mensaje para el usuario.
 *
 * Todos los campos de formulario llevan reglas. Son la primera barrera y
 * mejoran la experiencia; la seguridad real está en el backend (valida de
 * nuevo y usa SQL parametrizado) y en Vue (escapa todo lo que se pinta).
 */
type Regla = (v: unknown) => true | string

const texto = (v: unknown) => (v == null ? '' : String(v))

export const requerido: Regla = v =>
  texto(v).trim().length > 0 || 'Este campo es obligatorio'

export function minimo (n: number): Regla {
  return v =>
    texto(v).length === 0 || [...texto(v)].length >= n || `Usa al menos ${n} caracteres`
}

export function maximo (n: number): Regla {
  return v =>
    [...texto(v)].length <= n || `Usa como máximo ${n} caracteres`
}

export const correo: Regla = v =>
  texto(v).length === 0 || /^[^\s@<>]+@[^\s<>@][^\s.<>@]*\.[^\s<>@]+$/.test(texto(v)) || 'Escribe un correo válido'

// Usuario: letras, números, punto, guion y guion bajo (o un correo en el login).
export const usuarioOCorreo: Regla = v =>
  texto(v).length === 0 || /^[\w.@+-]+$/.test(texto(v).trim()) || 'Usa solo letras, números, punto, guion o @'

// Texto libre: sin etiquetas HTML.
export const textoSeguro: Regla = v =>
  !/[<>]/.test(texto(v)) || 'No uses los caracteres < ni >'

// Mismos límites que el backend (docs/09).
export const contrasena: Regla[] = [requerido, maximo(128)]
export const contrasenaNueva: Regla[] = [requerido, minimo(12), maximo(128)]

// La contraseña no debe contener el usuario ni el correo (misma regla del backend).
export function sinDatosPersonales (v: string, usuario?: string, correo?: string): true | string {
  const baja = v.toLowerCase()
  const contiene = [usuario, correo].some(d => !!d && baja.includes(d.toLowerCase()))
  return !contiene || 'No debe contener tu usuario ni tu correo'
}
