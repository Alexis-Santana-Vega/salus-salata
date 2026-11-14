import { crearRecurso } from './recursos'

// Lista cerrada de unidades (igual al CHECK de la base y al oneof del back).
export const UNIDADES = [
  { value: 'PIEZA', title: 'Pieza' },
  { value: 'CAJA', title: 'Caja' },
  { value: 'PAQUETE', title: 'Paquete' },
] as const

export function nombreUnidad (u: string) {
  return UNIDADES.find(x => x.value === u)?.title ?? u
}

export interface Insumo {
  id: string
  clave: string
  descripcion: string
  unidad: string
  precio: number
  codigoBarras: string
  stockActual: number
  stockMinimo: number
  bajoMinimo: boolean
  creadoEn: string
  actualizadoEn: string
}

// En el formulario los números son texto; se convierten al enviar.
export interface DatosInsumo {
  clave: string
  descripcion: string
  unidad: string | null
  precio: string
  codigoBarras: string
  stockActual: string
  stockMinimo: string
}

export const insumos = crearRecurso<Insumo, DatosInsumo>('/insumos')

export function aEnvio (d: DatosInsumo) {
  const num = (s: string) => (s.trim() === '' ? null : Number(s))
  return { ...d, precio: num(d.precio), stockActual: num(d.stockActual), stockMinimo: num(d.stockMinimo) }
}
