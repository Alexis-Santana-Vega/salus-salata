// Rutas de la app. Las páginas se cargan bajo demanda (cada una en su archivo JS).
import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login', component: () => import('@/pages/LoginPage.vue') },
    {
      path: '/',
      component: AppLayout,
      children: [
        { path: '', component: () => import('@/pages/InicioPage.vue') },
        // Errores dentro del layout: se pintan en el área de contenido.
        { path: ':pathMatch(.*)*', component: () => import('@/pages/NoEncontradaPage.vue') },
      ],
    },
  ],
})

export default router
