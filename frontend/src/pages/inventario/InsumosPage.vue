<!--
  Inventario de insumos (mock: Inventario). En la Fase 1 es el listado; los
  mosaicos de resumen, filtros, "Registrar entrada" y movimientos llegan en la Fase 7.
-->
<template>
  <page-header subtitulo="Consumibles desechables de un solo uso." titulo="Inventario de insumos">
    <template #acciones>
      <btn-primario :block="mobile" prepend-icon="$agregar" size="large" to="/insumos/alta">Nuevo insumo</btn-primario>
    </template>
  </page-header>

  <v-row>
    <v-col cols="12">
      <listado-crud
        v-model:busqueda="busqueda"
        v-model:orden="orden"
        v-model:pagina="pagina"
        v-model:por-pagina="porPagina"
        campo-titulo="descripcion"
        :cargando="cargando"
        :columnas="columnas"
        :error="error"
        :hay-busqueda="hayBusqueda"
        :items="items"
        :listo="listo"
        :nuevo="{ texto: 'Nuevo insumo', to: '/insumos/alta' }"
        placeholder="Buscar por clave, descripción o código de barras"
        :subtitulo="i => `${i.clave} · ${entero(i.stockActual)} en stock${i.bajoMinimo ? ' · bajo mínimo' : ''}`"
        texto-vacio="Da de alta los consumibles que se llevan a las cirugías."
        titulo-vacio="Aún no hay insumos"
        :total="total"
        @editar="i => router.push(`/insumos/${i.id}`)"
        @eliminar="i => pedir(i)"
        @reintentar="cargar"
      >
        <template #[`item.precio`]="{ item }">{{ moneda(item.precio) }}</template>
        <template #[`item.stockActual`]="{ item }">{{ entero(item.stockActual) }}</template>
        <template #[`item.stockMinimo`]="{ item }">{{ entero(item.stockMinimo) }}</template>

        <template #[`item.estado`]="{ item }">
          <v-chip :color="item.bajoMinimo ? 'warning' : 'success'" :text="item.bajoMinimo ? 'Bajo mínimo' : 'OK'" />
        </template>
      </listado-crud>
    </v-col>
  </v-row>

  <dialogo-eliminar
    v-model="abierto"
    :eliminando="eliminando"
    :texto="`${registro?.clave} dejará de aparecer al programar cirugías. Su clave no se podrá volver a usar.`"
    titulo="¿Eliminar el insumo?"
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
  import { type Insumo, insumos } from '@/services/insumos'
  import { entero, moneda } from '@/utils/formato'

  const router = useRouter()
  const { mobile } = useDisplay()
  const { items, total, cargando, listo, error, busqueda, pagina, porPagina, orden, hayBusqueda, cargar }
    = useListado(insumos.listar, 'clave')
  const { abierto, eliminando, registro, pedir, confirmar }
    = useBaja<Insumo>(insumos.eliminar, 'Insumo eliminado', cargar)

  // Columnas del mock (stock inicial y usadas llegan con los movimientos, Fase 7).
  const columnas: DataTableHeader<Insumo>[] = [
    { key: 'clave', title: 'Clave', sortable: true },
    { key: 'descripcion', title: 'Descripción / material', sortable: true },
    { key: 'precio', title: 'Precio', sortable: true, align: 'end' },
    { key: 'stockActual', title: 'Stock actual', sortable: true, align: 'end' },
    { key: 'stockMinimo', title: 'Mínimo', sortable: true, align: 'end' },
    { key: 'estado', title: 'Estado', sortable: false },
  ]
</script>
