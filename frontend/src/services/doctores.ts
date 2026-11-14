import type { Direccion } from '@/utils/direccion'
import { crearRecurso } from './recursos'

export interface Doctor {
  id: string
  numero: string // número interno: lo genera el sistema (solo lectura)
  nombre: string
  correo: string
  telPersonal: string
  telConsultorio: string
  direccion: Direccion
  creadoEn: string
  actualizadoEn: string
}

export interface DatosDoctor {
  numero: string // número interno: lo genera el sistema (solo lectura)
  nombre: string
  correo: string
  telPersonal: string
  telConsultorio: string
  direccion: Direccion
}

export const doctores = crearRecurso<Doctor, DatosDoctor>('/doctores')
