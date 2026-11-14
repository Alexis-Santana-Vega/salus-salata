/**
 * Sesión del usuario. El token de acceso vive solo en memoria (nunca en
 * localStorage); al recargar la página se recupera con /auth/refresh usando
 * la cookie httpOnly que el navegador guarda.
 */
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api, configurarApi } from '@/services/api'

export type Rol = 'SYSADMIN' | 'ADMIN' | 'DOCTOR'

export interface Usuario {
  id: number
  usuario: string
  correo: string
  nombre: string
  rol: Rol
  activo: boolean
  debeCambiarContrasena: boolean
  ultimoAcceso: string | null
}

interface Sesion { accessToken: string, expiresIn: number, usuario: Usuario }

export const NOMBRE_ROL: Record<Rol, string> = {
  SYSADMIN: 'Sysadmin',
  ADMIN: 'Administrador',
  DOCTOR: 'Doctor',
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const usuario = ref<Usuario | null>(null)
  let iniciada: Promise<void> | null = null

  const autenticado = computed(() => !!token.value && !!usuario.value)
  const iniciales = computed(() =>
    (usuario.value?.nombre ?? '').split(/\s+/).filter(Boolean).slice(0, 2).map(p => p[0]!.toUpperCase()).join(''),
  )

  function guardar (s: Sesion) {
    token.value = s.accessToken
    usuario.value = s.usuario
  }

  function limpiar () {
    token.value = null
    usuario.value = null
  }

  async function login (ident: string, contrasena: string) {
    guardar(await api<Sesion>('/auth/login', { metodo: 'POST', cuerpo: { usuario: ident, contrasena }, sinRenovar: true }))
  }

  // Renueva con la cookie. Devuelve false si no hay sesión válida.
  async function renovar (): Promise<boolean> {
    try {
      guardar(await api<Sesion>('/auth/refresh', { metodo: 'POST', sinRenovar: true }))
      return true
    } catch {
      limpiar()
      return false
    }
  }

  // Se llama una sola vez al abrir la app (desde la guarda del router).
  function iniciar () {
    iniciada ??= renovar().then(() => undefined)
    return iniciada
  }

  async function logout () {
    try {
      await api('/auth/logout', { metodo: 'POST', sinRenovar: true })
    } finally {
      limpiar()
    }
  }

  async function cambiarContrasena (actual: string, nueva: string) {
    guardar(await api<Sesion>('/auth/contrasena', {
      metodo: 'PUT', cuerpo: { contrasenaActual: actual, contrasenaNueva: nueva },
    }))
  }

  function tieneRol (roles?: Rol[]) {
    return !roles || (!!usuario.value && roles.includes(usuario.value.rol))
  }

  configurarApi({ token: () => token.value, renovar })

  return { token, usuario, autenticado, iniciales, login, logout, renovar, iniciar, cambiarContrasena, tieneRol }
})
