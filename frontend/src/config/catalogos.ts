/**
 * Pestañas de la página Catálogos (mock: Catálogos). Cada catálogo se agrega
 * aquí al construirse: su ruta, y el texto y destino del botón "Nuevo".
 */
export interface PestanaCatalogo {
  titulo: string
  to: string
  nuevo: { texto: string, to: string }
}

export const pestanasCatalogos: PestanaCatalogo[] = [
  {
    titulo: 'Doctores',
    to: '/catalogos/doctores',
    nuevo: { texto: 'Nuevo doctor', to: '/catalogos/doctores/alta' },
  },
  {
    titulo: 'Técnicos',
    to: '/catalogos/tecnicos',
    nuevo: { texto: 'Nuevo técnico', to: '/catalogos/tecnicos/alta' },
  },
  {
    titulo: 'Procedimientos',
    to: '/catalogos/procedimientos',
    nuevo: { texto: 'Nuevo procedimiento', to: '/catalogos/procedimientos/alta' },
  },
]
