/**
 * Flujo de baja con confirmación: abrir el diálogo con el registro, eliminar,
 * avisar y recargar. El error se muestra como toast (el diálogo se cierra).
 */
import { ref, shallowRef } from 'vue'
import { ApiError } from '@/services/api'
import { useAvisosStore } from '@/stores/avisos'

export function useBaja<T extends { id: string }> (eliminar: (id: string) => Promise<void>, mensaje: string, despues: () => unknown) {
  const avisos = useAvisosStore()
  const abierto = ref(false)
  const eliminando = ref(false)
  const registro = shallowRef<T>()

  function pedir (r: T) {
    registro.value = r
    abierto.value = true
  }

  async function confirmar () {
    if (!registro.value) {
      return
    }
    eliminando.value = true
    try {
      await eliminar(registro.value.id)
      avisos.exito(mensaje)
      await despues()
    } catch (error_) {
      avisos.error(error_ instanceof ApiError ? error_.mensaje : 'Algo salió mal. Intenta de nuevo.')
    } finally {
      eliminando.value = false
      abierto.value = false
    }
  }

  return { abierto, eliminando, registro, pedir, confirmar }
}
