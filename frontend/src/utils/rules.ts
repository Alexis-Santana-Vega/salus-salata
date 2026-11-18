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

// Teléfono opcional: + inicial opcional, dígitos, espacios, guiones o paréntesis (misma regla del back).
export const telefono: Regla = v =>
  texto(v).trim().length === 0 || /^\+?\d[\d ()-]{6,19}$/.test(texto(v).trim()) || 'Escribe un teléfono válido (solo números, espacios, guiones o +)'

// Código postal opcional: exactamente 5 dígitos (misma regla del back).
export const cp: Regla = v =>
  texto(v).trim().length === 0 || /^\d{5}$/.test(texto(v).trim()) || 'El código postal tiene 5 dígitos'

// Opción de una lista cerrada (vacío permitido; para obligatorio, junto con `requerido`).
export function enLista (opciones: readonly { value: string }[]): Regla {
  return v => v == null || v === '' || opciones.some(o => o.value === v) || 'Elige una opción de la lista'
}

/*
 * Formatos de texto libre: listas BLANCAS iguales a las del back
 * (backend/internal/handler/reglas.go). Todo lo demás se rechaza, incluidos
 * < > { } [ ] \ ` $ | ^ ~ = * % " ; caracteres de control y emojis.
 */
function formato (re: RegExp, msg: string): Regla {
  return v => texto(v).trim().length === 0 || re.test(texto(v).trim()) || msg
}

// Nombres de personas: letras (con acentos y ñ), espacios, punto, apóstrofo y guion.
export const nombrePersona = formato(/^[\p{L}\p{M}][\p{L}\p{M} .'-]*$/u, 'Usa solo letras, espacios, punto, apóstrofo o guion')

// Texto general (procedimientos, calles, puestos, nombres de hospital).
export const textoGeneral = formato(/^[\p{L}\p{M}\p{N}][\p{L}\p{M}\p{N} .,#/()'&+°-]*$/u, 'Usa solo letras, números, espacios y . , # / ( ) \' & + - °')

// Número exterior o interior ("120-B", "S/N").
export const numDomicilio = formato(/^[\p{L}\p{N}][\p{L}\p{N} #/.-]*$/u, 'Usa solo letras, números, espacios y # / . -')

// Códigos y claves (DOC-001, GTO-001).
export const codigo = formato(/^[A-Z0-9][A-Z0-9-]*$/i, 'Usa solo letras sin acento, números y guion')

// Número decimal opcional dentro de un rango (latitud, longitud).
export function rango (min: number, max: number, nombre: string): Regla {
  return v => {
    const t = texto(v).trim()
    if (t === '') {
      return true
    }
    const n = Number(t)
    return (/^-?\d+(?:\.\d+)?$/.test(t) && n >= min && n <= max) || `${nombre} entre ${min} y ${max}`
  }
}

// Liga opcional: debe ser https://
export const ligaHttps: Regla = v =>
  texto(v).trim().length === 0 || /^https:\/\/[^\s<>"]+$/.test(texto(v).trim()) || 'Escribe una liga completa que empiece con https://'

// Importe en pesos: no negativo y con a lo más 2 decimales (mismo límite del back).
export const dinero: Regla = v => {
  const t = texto(v).trim()
  if (t === '') {
    return true
  }
  return (/^\d{1,10}(?:\.\d{1,2})?$/.test(t)) || 'Escribe un importe válido, con a lo más 2 decimales'
}

// Entero no negativo (piezas en stock).
export const enteroPositivo: Regla = v => {
  const t = texto(v).trim()
  return t === '' || (/^\d{1,7}$/.test(t) && Number(t) <= 1_000_000) || 'Escribe un número entero de 0 a 1,000,000'
}

// Entero mayor a cero (ciclo, purga cada cuántos usos).
export const enteroMayorACero: Regla = v => {
  const t = texto(v).trim()
  return t === '' || (/^\d{1,6}$/.test(t) && Number(t) >= 1 && Number(t) <= 100_000) || 'Escribe un número entero de 1 a 100,000'
}

// Número de serie opcional: letras sin acento, números y . / - (misma regla del back).
export const serie = formato(/^[A-Z0-9][A-Z0-9./-]*$/i, 'Usa solo letras sin acento, números y . / -')

// Fecha "AAAA-MM-DD" que no sea posterior a hoy (las ISO se comparan como texto).
export function fechaNoFutura (hoy: string): Regla {
  return v => {
    const t = texto(v).trim()
    return t === '' || (/^\d{4}-\d{2}-\d{2}$/.test(t) && t <= hoy) || 'La fecha no puede ser futura'
  }
}
