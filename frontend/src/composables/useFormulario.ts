/**
 * Alta y edición con el mismo formulario: carga el registro si hay id,
 * guarda (crear o actualizar), pone los errores del servidor junto a su
 * campo y avisa con un toast. La página solo define sus campos.
 */
import type { Recurso } from '@/services/recursos'
import { computed, onMounted, reactive, ref, type Ref, toRaw, type UnwrapRef } from 'vue'
import { ApiError } from '@/services/api'
import { useAvisosStore } from '@/stores/avisos'

interface Opciones<T, D> {
  enviar?: (datos: D) => unknown // transforma los datos antes de mandarlos (p. ej. texto → número)
  id: Ref<string | undefined> // sin id = alta
  vacio: () => D // valores iniciales del alta
  desde: (registro: T) => D // registro → campos del formulario
  mensajes: { creado: string, actualizado: string }
  alGuardar: (registro: T) => void // p. ej. volver al listado
}

export function useFormulario<T, D extends object> (recurso: Recurso<T, D>, op: Opciones<T, D>) {
  const avisos = useAvisosStore()
  const datos = reactive(op.vacio()) as UnwrapRef<D>
  const errores = reactive<Record<string, string[]>>({})
  const error = ref('')
  const cargando = ref(false)
  const guardando = ref(false)
  const noEncontrado = ref(false)
  const esEdicion = computed(() => !!op.id.value)

  onMounted(async () => {
    if (!op.id.value) {
      return
    }
    cargando.value = true
    try {
      Object.assign(datos as object, op.desde(await recurso.obtener(op.id.value)))
    } catch (error_) {
      if (error_ instanceof ApiError && error_.status === 404) {
        noEncontrado.value = true
      } else {
        error.value = error_ instanceof ApiError ? error_.mensaje : 'Algo salió mal. Intenta de nuevo.'
      }
    } finally {
      cargando.value = false
    }
  })

  // Se llama ya validado por v-form (FormularioCrud).
  async function guardar () {
    error.value = ''
    for (const k of Object.keys(errores)) {
      delete errores[k]
    }
    guardando.value = true
    try {
      const crudo = toRaw(datos) as D
      const d = (op.enviar ? op.enviar(crudo) : crudo) as D
      const registro = op.id.value ? await recurso.actualizar(op.id.value, d) : await recurso.crear(d)
      avisos.exito(op.id.value ? op.mensajes.actualizado : op.mensajes.creado)
      op.alGuardar(registro)
    } catch (error_) {
      // Errores por campo junto a su campo; el resto, en la alerta de arriba.
      if (error_ instanceof ApiError && error_.detalles.length > 0) {
        for (const d of error_.detalles) {
          errores[d.field] = [d.message]
        }
        error.value = 'Revisa los campos marcados.'
      } else {
        error.value = error_ instanceof ApiError ? error_.mensaje : 'Algo salió mal. Intenta de nuevo.'
      }
    } finally {
      guardando.value = false
    }
  }

  // Al escribir en un campo se quita el error que mandó el servidor.
  function limpiarError (campo: string) {
    delete errores[campo]
  }

  return { datos, errores, error, cargando, guardando, noEncontrado, esEdicion, guardar, limpiarError }
}
