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
        { path: 'sin-permiso', component: () => import('@/pages/SinPermisoPage.vue') },
        // Errores dentro del layout: se pintan en el área de contenido.
        { path: ':pathMatch(.*)*', component: () => import('@/pages/NoEncontradaPage.vue') },
      ],
    },
  ],
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
