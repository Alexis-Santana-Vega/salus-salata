<!-- Listado de procedimientos (pestaña de Catálogos). -->
<template>
  <listado-crud
    v-model:busqueda="busqueda"
    v-model:orden="orden"
    v-model:pagina="pagina"
    v-model:por-pagina="porPagina"
    campo-titulo="nombre"
    :cargando="cargando"
    :columnas="columnas"
    :error="error"
    :hay-busqueda="hayBusqueda"
    :items="items"
    :listo="listo"
    :nuevo="{ texto: 'Nuevo procedimiento', to: '/catalogos/procedimientos/alta' }"
    placeholder="Buscar procedimiento"
    texto-vacio="Da de alta los tipos de cirugía (Holep, Artroscopia…) para usarlos al programar."
    titulo-vacio="Aún no hay procedimientos"
    :total="total"
    @editar="p => router.push(`/catalogos/procedimientos/${p.id}`)"
    @eliminar="p => pedir(p)"
    @reintentar="cargar"
  />

  <dialogo-eliminar
    v-model="abierto"
    :eliminando="eliminando"
    :texto="`«${registro?.nombre}» dejará de aparecer al programar cirugías. Las cirugías que ya lo usan lo conservan.`"
    titulo="¿Eliminar el procedimiento?"
    @confirmar="confirmar"
  />
</template>

<script lang="ts" setup>
  import type { DataTableHeader } from 'vuetify'
  import { useRouter } from 'vue-router'
  import DialogoEliminar from '@/components/DialogoEliminar.vue'
  import ListadoCrud from '@/components/ListadoCrud.vue'
  import { useBaja } from '@/composables/useBaja'
  import { useListado } from '@/composables/useListado'
  import { type Procedimiento, procedimientos } from '@/services/procedimientos'

  const router = useRouter()
  const { items, total, cargando, listo, error, busqueda, pagina, porPagina, orden, hayBusqueda, cargar }
    = useListado(procedimientos.listar, 'nombre')
  const { abierto, eliminando, registro, pedir, confirmar }
    = useBaja<Procedimiento>(procedimientos.eliminar, 'Procedimiento eliminado', cargar)

  const columnas: DataTableHeader<Procedimiento>[] = [
    { key: 'nombre', title: 'Nombre', sortable: true },
  ]
</script>
