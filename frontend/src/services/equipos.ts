import { crearRecurso } from './recursos'

export interface Tarifa {
  costoHora: number
  vigenteDesde: string // AAAA-MM-DD
}

export interface Equipo {
  id: string
  clave: string
  nombre: string
  marcaModelo: string
  cicloMax: number
  requierePurga: boolean
  purgaCada: number | null
  costoHora: number // vigente hoy
  tarifas: Tarifa[] // historial (solo en el detalle), la más reciente primero
  creadoEn: string
  actualizadoEn: string
}

// En el formulario los números son texto; se convierten al enviar.
export interface DatosEquipo {
  clave: string // solo lectura (la genera el sistema)
  nombre: string
  marcaModelo: string
  cicloMax: string
  requierePurga: boolean
  purgaCada: string
  costoHora: string
  tarifas: Tarifa[] // solo lectura
}

export const equipos = crearRecurso<Equipo, DatosEquipo>('/equipos')

export function aEnvio (d: DatosEquipo) {
  const num = (s: string) => (s.trim() === '' ? null : Number(s))
  return {
    nombre: d.nombre, marcaModelo: d.marcaModelo, requierePurga: d.requierePurga,
    cicloMax: num(d.cicloMax), costoHora: num(d.costoHora),
    purgaCada: d.requierePurga ? num(d.purgaCada) : null,
  }
}
