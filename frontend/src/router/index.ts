// Rutas y guardas. Las guardas solo mejoran la experiencia: la seguridad real
// la aplica el backend en cada petición.
import type { Rol } from '@/stores/auth'
import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
  interface RouteMeta {
    publica?: boolean // no requiere sesión (login)
    roles?: Rol[] // si se omite, cualquier rol con sesión
  }
}

// Roles que administran catálogos (docs/02 › Roles y permisos).
const ADMINISTRACION: Rol[] = ['SYSADMIN', 'ADMIN']

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login', component: () => import('@/pages/LoginPage.vue'), meta: { publica: true } },
    { path: '/cambiar-contrasena', component: () => import('@/pages/CambiarContrasenaPage.vue') },
    {
      path: '/',
      component: AppLayout,
      children: [
        { path: '', component: () => import('@/pages/InicioPage.vue') },
        {
          path: 'catalogos',
          component: () => import('@/pages/catalogos/CatalogosPage.vue'),
          meta: { roles: ADMINISTRACION },
          children: [
            { path: '', redirect: '/catalogos/doctores' },
            { path: 'doctores', component: () => import('@/pages/catalogos/DoctoresPage.vue') },
            { path: 'tecnicos', component: () => import('@/pages/catalogos/TecnicosPage.vue') },
            { path: 'procedimientos', component: () => import('@/pages/catalogos/ProcedimientosPage.vue') },
          ],
        },
        // Altas y ediciones: página propia, fuera de las pestañas.
        { path: 'catalogos/doctores/alta', component: () => import('@/pages/catalogos/DoctorFormPage.vue'), meta: { roles: ADMINISTRACION } },
        { path: 'catalogos/doctores/:id', component: () => import('@/pages/catalogos/DoctorFormPage.vue'), meta: { roles: ADMINISTRACION } },
        { path: 'catalogos/tecnicos/alta', component: () => import('@/pages/catalogos/TecnicoFormPage.vue'), meta: { roles: ADMINISTRACION } },
        { path: 'catalogos/tecnicos/:id', component: () => import('@/pages/catalogos/TecnicoFormPage.vue'), meta: { roles: ADMINISTRACION } },
        { path: 'catalogos/procedimientos/alta', component: () => import('@/pages/catalogos/ProcedimientoFormPage.vue'), meta: { roles: ADMINISTRACION } },
        { path: 'catalogos/procedimientos/:id', component: () => import('@/pages/catalogos/ProcedimientoFormPage.vue'), meta: { roles: ADMINISTRACION } },
        { path: 'sin-permiso', component: () => import('@/pages/SinPermisoPage.vue') },
        // Errores dentro del layout: se pintan en el área de contenido.
        { path: ':pathMatch(.*)*', component: () => import('@/pages/NoEncontradaPage.vue') },
      ],
    },
  ],
})

// Tras un despliegue nuevo, una pestaña que ya estaba abierta pide archivos de
// pantallas que ya no existen (cambian de nombre en cada build) y la navegación
// falla en silencio. En ese caso se recarga la app directo en la pantalla
// destino. Solo una vez por destino cada 10 s para no ciclar si el servidor cae.
const ERROR_DE_CARGA = /dynamically imported module|Importing a module script failed|error loading dynamically imported module|Unable to preload CSS/i

export function recargarEn (destino: string) {
  const clave = `salus.recarga:${destino}`
  try {
    const antes = Number(sessionStorage.getItem(clave) || 0)
    if (Date.now() - antes < 10_000) {
      return false
    }
    sessionStorage.setItem(clave, String(Date.now()))
  } catch {
    // almacenamiento bloqueado: se recarga igual
  }
  window.location.assign(destino)
  return true
}

router.onError((error, to) => {
  if (ERROR_DE_CARGA.test(String(error?.message ?? error))) {
    recargarEn(to.fullPath)
  }
})

router.beforeEach(async to => {
  const auth = useAuthStore()
  await auth.iniciar() // recupera la sesión con la cookie la primera vez

  if (to.meta.publica) {
    return auth.autenticado ? '/' : true
  }
  if (!auth.autenticado) {
    return { path: '/login', query: to.fullPath === '/' ? {} : { volver: to.fullPath } }
  }
  // Primer acceso: hasta cambiar la contraseña no puede ir a otro lado.
  if (auth.usuario!.debeCambiarContrasena) {
    return to.path === '/cambiar-contrasena' ? true : '/cambiar-contrasena'
  }
  if (to.path === '/cambiar-contrasena') {
    return true
  }
  if (!auth.tieneRol(to.meta.roles)) {
    return '/sin-permiso'
  }
  return true
})

export default router
