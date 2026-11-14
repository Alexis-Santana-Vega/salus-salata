import { api } from './api'

export interface Asentamiento {
  id: string
  cp: string
  nombre: string
  tipo: string
  zona: string
  ciudad: string
  municipio: string
  estado: string
  estadoNombre: string
}

// Colonias de un código postal. 404 si no está en el catálogo.
export function coloniasPorCP (cp: string) {
  return api<{ cp: string, asentamientos: Asentamiento[] }>(`/codigos-postales/${encodeURIComponent(cp)}`)
}
