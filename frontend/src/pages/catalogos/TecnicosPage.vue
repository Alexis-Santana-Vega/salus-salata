<!-- Listado de técnicos (pestaña de Catálogos; mock: Catálogos › Técnicos). -->
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
    :nuevo="{ texto: 'Nuevo técnico', to: '/catalogos/tecnicos/alta' }"
    placeholder="Buscar por nombre o teléfono"
    :subtitulo="t => [nombrePuesto(t.puesto), t.telefono].filter(Boolean).join(' · ')"
    texto-vacio="Da de alta al personal que entrega el material y asiste en las cirugías."
    titulo-vacio="Aún no hay técnicos"
    :total="total"
    @editar="t => router.push(`/catalogos/tecnicos/${t.id}`)"
    @eliminar="t => pedir(t)"
    @reintentar="cargar"
  >
    <template #[`item.puesto`]="{ item }">{{ nombrePuesto(item.puesto) }}</template>
    <template #[`item.telefono`]="{ item }">{{ item.telefono || '—' }}</template>
  </listado-crud>

  <dialogo-eliminar
    v-model="abierto"
    :eliminando="eliminando"
    :texto="`${registro?.nombre} dejará de aparecer al asignar personal a las cirugías. Las cirugías donde ya participó lo conservan.`"
    titulo="¿Eliminar al técnico?"
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
  import { nombrePuesto, type Tecnico, tecnicos } from '@/services/tecnicos'

  const router = useRouter()
  const { items, total, cargando, listo, error, busqueda, pagina, porPagina, orden, hayBusqueda, cargar }
    = useListado(tecnicos.listar, 'nombre')
  const { abierto, eliminando, registro, pedir, confirmar }
    = useBaja<Tecnico>(tecnicos.eliminar, 'Técnico eliminado', cargar)

  const columnas: DataTableHeader<Tecnico>[] = [
    { key: 'nombre', title: 'Nombre', sortable: true },
    { key: 'puesto', title: 'Puesto', sortable: true },
    { key: 'telefono', title: 'Teléfono', sortable: false },
  ]
</script>
