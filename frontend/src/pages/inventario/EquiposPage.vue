<!--
  Equipos (modelos maestros) con su número de unidades. Estado, ciclo y
  mantenimiento por unidad: Fase 4 (mock: Equipo).
-->
<template>
  <page-header subtitulo="Equipo reutilizable que se renta por hora." titulo="Equipos">
    <template #acciones>
      <btn-primario :block="mobile" prepend-icon="$agregar" size="large" to="/equipos/alta">Nuevo equipo</btn-primario>
    </template>
  </page-header>

  <v-row>
    <v-col cols="12">
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
        :nuevo="{ texto: 'Nuevo equipo', to: '/equipos/alta' }"
        placeholder="Buscar por clave, nombre o marca"
        :subtitulo="e => `${e.clave} · ${e.marcaModelo} · ${moneda(e.costoHora)} por hora · ${unidades(e.numUnidades)}`"
        texto-vacio="Da de alta los modelos de equipo que se rentan (torres, láseres, morceladores…)."
        titulo-vacio="Aún no hay equipos"
        :total="total"
        @editar="e => router.push(`/equipos/${e.id}`)"
        @eliminar="e => pedir(e)"
        @reintentar="cargar"
      >
        <template #[`item.costoHora`]="{ item }">{{ moneda(item.costoHora) }}</template>
        <template #[`item.numUnidades`]="{ item }">{{ entero(item.numUnidades) }}</template>
        <template #[`item.cicloMax`]="{ item }">{{ entero(item.cicloMax) }} usos</template>
        <template #[`item.purga`]="{ item }">{{ item.requierePurga ? `Cada ${item.purgaCada} usos` : 'No requiere' }}</template>
      </listado-crud>
    </v-col>
  </v-row>

  <dialogo-eliminar
    v-model="abierto"
    :eliminando="eliminando"
    :texto="`${registro?.clave} ${registro?.nombre} dejará de aparecer al programar cirugías. Su clave no se podrá volver a usar.`"
    titulo="¿Eliminar el equipo?"
    @confirmar="confirmar"
  />
</template>

<script lang="ts" setup>
  import type { DataTableHeader } from 'vuetify'
  import { useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import DialogoEliminar from '@/components/DialogoEliminar.vue'
  import ListadoCrud from '@/components/ListadoCrud.vue'
  import PageHeader from '@/components/PageHeader.vue'
  import { useBaja } from '@/composables/useBaja'
  import { useListado } from '@/composables/useListado'
  import { type Equipo, equipos } from '@/services/equipos'
  import { entero, moneda, unidades } from '@/utils/formato'

  const router = useRouter()
  const { mobile } = useDisplay()
  const { items, total, cargando, listo, error, busqueda, pagina, porPagina, orden, hayBusqueda, cargar }
    = useListado(equipos.listar, 'clave')
  const { abierto, eliminando, registro, pedir, confirmar }
    = useBaja<Equipo>(equipos.eliminar, 'Equipo eliminado', cargar)

  const columnas: DataTableHeader<Equipo>[] = [
    { key: 'clave', title: 'Clave', sortable: true },
    { key: 'nombre', title: 'Nombre', sortable: true },
    { key: 'marcaModelo', title: 'Marca / modelo', sortable: true },
    { key: 'costoHora', title: 'Costo por hora', sortable: false, align: 'end' },
    { key: 'numUnidades', title: 'Unidades', sortable: false, align: 'end' },
    { key: 'cicloMax', title: 'Ciclo', sortable: true, align: 'end' },
    { key: 'purga', title: 'Purga', sortable: false },
  ]
</script>
