import { crearRecurso } from './recursos'

export interface Procedimiento {
  id: string
  nombre: string
  creadoEn: string
  actualizadoEn: string
}

export interface DatosProcedimiento {
  nombre: string
}

export const procedimientos = crearRecurso<Procedimiento, DatosProcedimiento>('/procedimientos')
