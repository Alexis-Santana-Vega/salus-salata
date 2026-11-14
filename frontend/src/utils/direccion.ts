// Dirección de México: se captura calle, números y la colonia del catálogo
// (por su id); C.P., municipio y estado vienen del catálogo SEPOMEX/INEGI.
export interface Direccion {
  calle: string
  numExterior: string
  numInterior: string
  asentamientoId: string | null
  // Solo lectura (los llena el catálogo):
  cp: string
  colonia: string
  municipio: string
  estado: string // abreviatura (HGO)
  estadoNombre: string
}

export function direccionVacia (): Direccion {
  return {
    calle: '', numExterior: '', numInterior: '', asentamientoId: null,
    cp: '', colonia: '', municipio: '', estado: '', estadoNombre: '',
  }
}

// "Pachuca de Soto, HGO" para listados (vacío si no hay datos).
export function ciudad (d: Direccion) {
  return [d.municipio, d.estado].filter(Boolean).join(', ')
}
