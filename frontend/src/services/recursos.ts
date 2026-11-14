/**
 * Cliente CRUD genérico para los catálogos. Todos siguen el mismo contrato
 * de la API (docs/09 › CRUD de catálogos), así que cada catálogo solo declara
 * su ruta y sus tipos.
 */
import { api, apiPagina, type Pagina } from './api'

export interface ParametrosListado {
  page: number
  pageSize: number
  q?: string
  sort?: string // "campo" o "-campo"
}

export interface Recurso<T, D> {
  listar: (p: ParametrosListado) => Promise<Pagina<T>>
  obtener: (id: string) => Promise<T>
  crear: (datos: D) => Promise<T>
  actualizar: (id: string, datos: D) => Promise<T>
  eliminar: (id: string) => Promise<void>
}

export function crearRecurso<T, D> (ruta: string): Recurso<T, D> {
  return {
    listar (p) {
      const qs = new URLSearchParams({ page: String(p.page), pageSize: String(p.pageSize) })
      if (p.q) {
        qs.set('q', p.q)
      }
      if (p.sort) {
        qs.set('sort', p.sort)
      }
      return apiPagina<T>(`${ruta}?${qs}`)
    },
    obtener: id => api<T>(`${ruta}/${encodeURIComponent(id)}`),
    crear: datos => api<T>(ruta, { metodo: 'POST', cuerpo: datos }),
    actualizar: (id, datos) => api<T>(`${ruta}/${encodeURIComponent(id)}`, { metodo: 'PUT', cuerpo: datos }),
    eliminar: id => api<void>(`${ruta}/${encodeURIComponent(id)}`, { metodo: 'DELETE' }),
  }
}
