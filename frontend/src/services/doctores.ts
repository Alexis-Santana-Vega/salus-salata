import type { Direccion } from '@/utils/direccion'
import { crearRecurso } from './recursos'

export interface Doctor {
  id: string
  numero: string
  nombre: string
  correo: string
  telPersonal: string
  telConsultorio: string
  direccion: Direccion
  creadoEn: string
  actualizadoEn: string
}

export interface DatosDoctor {
  numero: string
  nombre: string
  correo: string
  telPersonal: string
  telConsultorio: string
  direccion: Direccion
}

export const doctores = crearRecurso<Doctor, DatosDoctor>('/doctores')
