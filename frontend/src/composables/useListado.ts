/**
 * Estado de un listado paginado del servidor: página, filas por página, orden
 * y búsqueda (con espera para no pedir en cada tecla). Lo usan ListadoCrud y
 * cada página de catálogo; la tabla y la lista móvil comparten este estado.
 */
import type { Pagina } from '@/services/api'
import type { ParametrosListado } from '@/services/recursos'
import { computed, onMounted, ref, shallowRef, watch } from 'vue'
import { ApiError } from '@/services/api'

export interface Orden { key: string, order?: 'asc' | 'desc' | boolean }

const ESPERA_BUSQUEDA = 300 // ms

export function useListado<T> (listar: (p: ParametrosListado) => Promise<Pagina<T>>, ordenInicial: string) {
  const items = shallowRef<T[]>([])
  const total = ref(0)
  const cargando = ref(false)
  const listo = ref(false) // ya hubo una primera respuesta (para el skeleton de página)
  const error = ref('')
  const busqueda = ref('')
  const pagina = ref(1)
  const porPagina = ref(10)
  const orden = ref<Orden[]>([{ key: ordenInicial, order: 'asc' }])

  const hayBusqueda = computed(() => busqueda.value.trim() !== '')

  // Solo cuenta la respuesta de la última petición (evita que una lenta pise a otra).
  let turno = 0

  async function cargar () {
    const mio = ++turno
    cargando.value = true
    error.value = ''
    const o = orden.value[0]
    try {
      const r = await listar({
        page: pagina.value,
        pageSize: porPagina.value,
        q: busqueda.value.trim() || undefined,
        sort: o ? (o.order === 'desc' ? '-' : '') + o.key : undefined,
      })
      if (mio !== turno) {
        return
      }
      items.value = r.data
      total.value = r.meta.total
    } catch (error_) {
      if (mio !== turno) {
        return
      }
      error.value = error_ instanceof ApiError ? error_.mensaje : 'Algo salió mal. Intenta de nuevo.'
    } finally {
      if (mio === turno) {
        cargando.value = false
        listo.value = true
      }
    }
  }

  let espera: ReturnType<typeof setTimeout> | undefined
  watch(busqueda, () => {
    clearTimeout(espera)
    espera = setTimeout(() => {
      if (pagina.value === 1) {
        cargar()
      } else {
        pagina.value = 1
      } // el watch de abajo recarga
    }, ESPERA_BUSQUEDA)
  })
  watch([pagina, porPagina, orden], cargar, { deep: true })
  onMounted(cargar)

  function limpiarBusqueda () {
    busqueda.value = ''
  }

  return { items, total, cargando, listo, error, busqueda, pagina, porPagina, orden, hayBusqueda, cargar, limpiarBusqueda }
}
