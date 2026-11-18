import { crearRecurso } from './recursos'

export interface Tarifa {
  costoHora: number
  vigenteDesde: string // AAAA-MM-DD
}

export interface Unidad {
  id: string
  clave: string // EQ-001-01, la genera el sistema
  numSerie: string
  fechaAlta: string // AAAA-MM-DD
  usosPrevios: number
}

// Fila del formulario: sin id es nueva (aún sin clave).
export interface DatosUnidad {
  id?: string
  clave: string
  numSerie: string
  fechaAlta: string
  usosPrevios: string
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
  numUnidades: number
  unidades: Unidad[] // activas, por clave (solo en el detalle)
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
  unidades: DatosUnidad[]
}

export const equipos = crearRecurso<Equipo, DatosEquipo>('/equipos')

export function aEnvio (d: DatosEquipo) {
  const num = (s: string) => (s.trim() === '' ? null : Number(s))
  return {
    nombre: d.nombre, marcaModelo: d.marcaModelo, requierePurga: d.requierePurga,
    cicloMax: num(d.cicloMax), costoHora: num(d.costoHora),
    purgaCada: d.requierePurga ? num(d.purgaCada) : null,
    // Estado completo: las unidades guardadas que no vengan se dan de baja.
    unidades: d.unidades.map(u => ({ id: u.id, numSerie: u.numSerie.trim(), fechaAlta: u.fechaAlta, usosPrevios: num(u.usosPrevios) })),
  }
}
