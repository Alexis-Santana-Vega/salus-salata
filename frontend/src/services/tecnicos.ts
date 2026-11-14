import { crearRecurso } from './recursos'

// Lista cerrada de puestos (igual al CHECK de la base y al oneof del back).
export const PUESTOS = [
  { value: 'AYUDANTE', title: 'Ayudante' },
  { value: 'ENFERMERIA', title: 'Enfermería quirúrgica' },
  { value: 'ANESTESIOLOGIA', title: 'Anestesiología' },
  { value: 'LOGISTICA', title: 'Logística / entrega' },
  { value: 'OTRO', title: 'Otro' },
] as const

export type Puesto = typeof PUESTOS[number]['value']

export function nombrePuesto (p: string) {
  return PUESTOS.find(x => x.value === p)?.title ?? p
}

export interface Tecnico {
  id: string
  nombre: string
  puesto: Puesto
  telefono: string
  creadoEn: string
  actualizadoEn: string
}

export interface DatosTecnico {
  nombre: string
  puesto: Puesto | null
  telefono: string
}

export const tecnicos = crearRecurso<Tecnico, DatosTecnico>('/tecnicos')
