/**
 * Avisos breves (toasts) de toda la app. Se pintan en AppAvisos con
 * v-snackbar-queue; las páginas solo llaman exito() o error().
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Aviso {
  text: string
  prependIcon: string
  color?: string // color del icono (el fondo siempre es surface)
}

export const useAvisosStore = defineStore('avisos', () => {
  const cola = ref<Aviso[]>([])

  function exito (text: string) {
    cola.value.push({ text, prependIcon: '$success' })
  }

  function error (text: string) {
    cola.value.push({ text, prependIcon: '$error' })
  }

  return { cola, exito, error }
})
