import type { Rol } from '@/stores/auth'

/**
 * Menú del sidebar. Una opción sola lleva directo a su página; una opción
 * con `hijos` se despliega con la página principal y sus altas.
 * Los módulos se agregan aquí conforme se construyen.
 */
export interface OpcionMenu {
  titulo: string
  icono: string
  to?: string
  roles?: Rol[] // si se omite, la ven todos los roles con sesión
  hijos?: { titulo: string, to: string }[]
}

export interface SeccionMenu {
  titulo: string
  opciones: OpcionMenu[]
}

export const menu: SeccionMenu[] = [
  {
    titulo: 'Menú',
    opciones: [
      { titulo: 'Dashboard', icono: '$inicio', to: '/' },
    ],
  },
  {
    titulo: 'Administración',
    opciones: [
      {
        titulo: 'Catálogos',
        icono: '$catalogos',
        roles: ['SYSADMIN', 'ADMIN'],
        hijos: [
          { titulo: 'Catálogos', to: '/catalogos/procedimientos' },
          { titulo: 'Alta de procedimiento', to: '/catalogos/procedimientos/alta' },
        ],
      },
    ],
  },
  {
    titulo: 'Sistema',
    opciones: [
      {
        titulo: 'Usuarios',
        icono: '$usuarios',
        roles: ['SYSADMIN'],
        hijos: [
          { titulo: 'Usuarios', to: '/usuarios' },
          { titulo: 'Alta de usuario', to: '/usuarios/alta' },
        ],
      },
    ],
  },
]
