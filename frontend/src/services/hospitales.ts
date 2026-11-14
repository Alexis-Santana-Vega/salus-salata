import type { Direccion } from '@/utils/direccion'
import { crearRecurso } from './recursos'

export interface Contacto {
  nombre: string
  puesto: string
  telefono: string
}

export interface Hospital {
  id: string
  clave: string
  nombre: string
  telefono: string
  direccion: Direccion
  latitud: number | null
  longitud: number | null
  ligaMapa: string
  contactos: Contacto[]
  creadoEn: string
  actualizadoEn: string
}

// En el formulario latitud y longitud son texto; se convierten al enviar.
export interface DatosHospital {
  clave: string
  nombre: string
  telefono: string
  direccion: Direccion
  latitud: string
  longitud: string
  ligaMapa: string
  contactos: Contacto[]
}

export const hospitales = crearRecurso<Hospital, DatosHospital>('/hospitales')

export function aEnvio (d: DatosHospital) {
  const num = (s: string) => (s.trim() === '' ? null : Number(s))
  return { ...d, latitud: num(d.latitud), longitud: num(d.longitud) }
}
