<!-- Listado de hospitales (pestaña de Catálogos; mock: Catálogos › Hospitales). La foto llega en la Fase 9. -->
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
    :nuevo="{ texto: 'Nuevo hospital', to: '/catalogos/hospitales/alta' }"
    placeholder="Buscar por clave, nombre o municipio"
    :subtitulo="h => [h.clave, ciudad(h.direccion)].filter(Boolean).join(' · ')"
    texto-vacio="Da de alta los hospitales a los que se renta el equipo."
    titulo-vacio="Aún no hay hospitales"
    :total="total"
    @editar="h => router.push(`/catalogos/hospitales/${h.id}`)"
    @eliminar="h => pedir(h)"
    @reintentar="cargar"
  >
    <template #[`item.ciudad`]="{ item }">{{ ciudad(item.direccion) || '—' }}</template>
    <template #[`item.telefono`]="{ item }">{{ item.telefono || '—' }}</template>
    <template #[`item.contacto`]="{ item }">{{ item.contactos[0]?.nombre || '—' }}</template>
  </listado-crud>

  <dialogo-eliminar
    v-model="abierto"
    :eliminando="eliminando"
    :texto="`${registro?.nombre} dejará de aparecer al programar cirugías. Las cirugías que ya tiene lo conservan.`"
    titulo="¿Eliminar el hospital?"
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
  import { type Hospital, hospitales } from '@/services/hospitales'
  import { ciudad } from '@/utils/direccion'

  const router = useRouter()
  const { items, total, cargando, listo, error, busqueda, pagina, porPagina, orden, hayBusqueda, cargar }
    = useListado(hospitales.listar, 'nombre')
  const { abierto, eliminando, registro, pedir, confirmar }
    = useBaja<Hospital>(hospitales.eliminar, 'Hospital eliminado', cargar)

  const columnas: DataTableHeader<Hospital>[] = [
    { key: 'clave', title: 'Clave', sortable: true },
    { key: 'nombre', title: 'Nombre', sortable: true },
    { key: 'ciudad', title: 'Ciudad', sortable: false },
    { key: 'telefono', title: 'Teléfono', sortable: false },
    { key: 'contacto', title: 'Contacto', sortable: false },
  ]
</script>
