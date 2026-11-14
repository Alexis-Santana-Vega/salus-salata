// Estado de la interfaz: tema y sidebar. La preferencia de tema se guarda
// en este navegador; cuando exista el perfil se guardará por usuario.
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useTheme } from 'vuetify'

const CLAVE_TEMA = 'salus.tema'

export const useAppStore = defineStore('app', () => {
  const theme = useTheme()
  const drawer = ref<boolean | null>(null) // null: Vuetify decide según el ancho

  // Recupera el tema guardado (si el almacenamiento está bloqueado, se ignora).
  try {
    const guardado = localStorage.getItem(CLAVE_TEMA)
    if (guardado === 'salus' || guardado === 'salusOscuro') {
      theme.change(guardado)
    }
  } catch {
    // almacenamiento bloqueado: se usa el tema por defecto
  }

  function alternarTema () {
    theme.toggle(['salus', 'salusOscuro'])
  }

  watch(() => theme.global.name.value, nombre => {
    try {
      localStorage.setItem(CLAVE_TEMA, nombre)
    } catch {
    // almacenamiento bloqueado: se usa el tema por defecto
    }
  })

  return { drawer, alternarTema, esOscuro: () => theme.global.current.value.dark }
})
