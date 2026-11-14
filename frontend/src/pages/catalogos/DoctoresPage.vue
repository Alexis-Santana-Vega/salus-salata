<!-- Listado de doctores (pestaña de Catálogos; mock: Catálogos › Doctores). -->
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
    :nuevo="{ texto: 'Nuevo doctor', to: '/catalogos/doctores/alta' }"
    placeholder="Buscar por nombre, número o correo"
    :subtitulo="d => [d.correo, ciudad(d.direccion)].filter(Boolean).join(' · ')"
    texto-vacio="Da de alta a los médicos encargados de las cirugías."
    titulo-vacio="Aún no hay doctores"
    :total="total"
    @editar="d => router.push(`/catalogos/doctores/${d.id}`)"
    @eliminar="d => pedir(d)"
    @reintentar="cargar"
  >
    <template #[`item.numero`]="{ item }">{{ item.numero || '—' }}</template>
    <template #[`item.telPersonal`]="{ item }">{{ item.telPersonal || '—' }}</template>
    <template #[`item.telConsultorio`]="{ item }">{{ item.telConsultorio || '—' }}</template>
    <template #[`item.ciudad`]="{ item }">{{ ciudad(item.direccion) || '—' }}</template>
  </listado-crud>

  <dialogo-eliminar
    v-model="abierto"
    :eliminando="eliminando"
    :texto="`${registro?.nombre} dejará de aparecer al programar cirugías. Las cirugías que ya tiene lo conservan.`"
    titulo="¿Eliminar al doctor?"
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
  import { type Doctor, doctores } from '@/services/doctores'
  import { ciudad } from '@/utils/estados'

  const router = useRouter()
  const { items, total, cargando, listo, error, busqueda, pagina, porPagina, orden, hayBusqueda, cargar }
    = useListado(doctores.listar, 'nombre')
  const { abierto, eliminando, registro, pedir, confirmar }
    = useBaja<Doctor>(doctores.eliminar, 'Doctor eliminado', cargar)

  // Columnas del mock (el acceso al portal llega en la Fase 6).
  const columnas: DataTableHeader<Doctor>[] = [
    { key: 'numero', title: 'Número', sortable: true },
    { key: 'nombre', title: 'Nombre', sortable: true },
    { key: 'telPersonal', title: 'Tel. personal', sortable: false },
    { key: 'telConsultorio', title: 'Tel. consultorio', sortable: false },
    { key: 'ciudad', title: 'Ciudad', sortable: false },
    { key: 'correo', title: 'Correo', sortable: true },
  ]
</script>
