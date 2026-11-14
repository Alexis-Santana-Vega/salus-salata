// Las 32 entidades federativas (dominio D_ESTADO de la base y model.Estados del back).
export const ESTADOS = [
  { value: 'AGS', title: 'Aguascalientes' },
  { value: 'BC', title: 'Baja California' },
  { value: 'BCS', title: 'Baja California Sur' },
  { value: 'CAMP', title: 'Campeche' },
  { value: 'COAH', title: 'Coahuila' },
  { value: 'COL', title: 'Colima' },
  { value: 'CHIS', title: 'Chiapas' },
  { value: 'CHIH', title: 'Chihuahua' },
  { value: 'CDMX', title: 'Ciudad de México' },
  { value: 'DGO', title: 'Durango' },
  { value: 'GTO', title: 'Guanajuato' },
  { value: 'GRO', title: 'Guerrero' },
  { value: 'HGO', title: 'Hidalgo' },
  { value: 'JAL', title: 'Jalisco' },
  { value: 'MEX', title: 'Estado de México' },
  { value: 'MICH', title: 'Michoacán' },
  { value: 'MOR', title: 'Morelos' },
  { value: 'NAY', title: 'Nayarit' },
  { value: 'NL', title: 'Nuevo León' },
  { value: 'OAX', title: 'Oaxaca' },
  { value: 'PUE', title: 'Puebla' },
  { value: 'QRO', title: 'Querétaro' },
  { value: 'QROO', title: 'Quintana Roo' },
  { value: 'SLP', title: 'San Luis Potosí' },
  { value: 'SIN', title: 'Sinaloa' },
  { value: 'SON', title: 'Sonora' },
  { value: 'TAB', title: 'Tabasco' },
  { value: 'TAMPS', title: 'Tamaulipas' },
  { value: 'TLAX', title: 'Tlaxcala' },
  { value: 'VER', title: 'Veracruz' },
  { value: 'YUC', title: 'Yucatán' },
  { value: 'ZAC', title: 'Zacatecas' },
]

export interface Direccion {
  calle: string
  numExterior: string
  numInterior: string
  colonia: string
  cp: string
  municipio: string
  estado: string | null
}

export function direccionVacia (): Direccion {
  return { calle: '', numExterior: '', numInterior: '', colonia: '', cp: '', municipio: '', estado: null }
}

// "León, GTO" para listados (vacío si no hay datos).
export function ciudad (d: Direccion) {
  return [d.municipio, d.estado].filter(Boolean).join(', ')
}
