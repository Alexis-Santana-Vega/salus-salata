<!--
  Listado de un catálogo (mock: Catálogos / CatalogosMovil, guía 05 › CRUD).
  Web: card con buscador y tabla paginada en el servidor, columnas ordenables y
  acciones por fila (editar, eliminar). Móvil: lista de tarjetas con chevron.
  Todo el diseño del listado vive aquí: cambiarlo cambia todos los catálogos.
-->
<template>
  <v-card>
    <v-card-text>
      <campo-busqueda
        v-model="busqueda"
        :aria-label="placeholder"
        clearable
        :placeholder="placeholder"
      />
    </v-card-text>

    <!-- Primera carga: un solo skeleton para todo el contenido -->
    <v-skeleton-loader v-if="!listo" :type="mobile ? 'list-item-two-line@5' : 'table-tbody'" />

    <v-fade-transition v-else mode="out-in">
      <!-- Error al consultar -->
      <v-empty-state
        v-if="error"
        key="error"
        icon="$sinConexion"
        :text="error"
        title="No pudimos cargar la información"
      >
        <template #actions>
          <btn-primario :disabled="cargando" :loading="cargando" @click="emit('reintentar')">
            {{ cargando ? 'Cargando' : 'Reintentar' }}
          </btn-primario>
        </template>
      </v-empty-state>

      <!-- Sin registros o sin coincidencias -->
      <v-empty-state
        v-else-if="total === 0"
        key="vacio"
        :icon="hayBusqueda ? '$sinResultados' : '$vacio'"
        :text="hayBusqueda ? 'Prueba con otra palabra.' : textoVacio"
        :title="hayBusqueda ? `Sin coincidencias para “${busqueda.trim()}”` : tituloVacio"
      >
        <template #actions>
          <v-btn v-if="hayBusqueda" prepend-icon="$limpiar" @click="busqueda = ''">Limpiar búsqueda</v-btn>
          <btn-primario v-else-if="nuevo" prepend-icon="$agregar" :to="nuevo.to">{{ nuevo.texto }}</btn-primario>
        </template>
      </v-empty-state>

      <!-- Móvil: tarjetas -->
      <div v-else-if="mobile" key="movil">
        <v-list lines="two">
          <v-list-item
            v-for="item in items"
            :key="(item as Fila).id"
            append-icon="$abrir"
            :subtitle="subtitulo?.(item)"
            :title="String((item as Fila)[campoTitulo])"
            @click="emit('editar', item)"
          />
        </v-list>

        <v-pagination v-model="pagina" :disabled="cargando" :length="paginas" total-visible="5" />
      </div>

      <!-- Web: tabla -->
      <v-data-table-server
        v-else
        key="tabla"
        v-model:items-per-page="porPagina"
        v-model:page="pagina"
        v-model:sort-by="orden"
        :headers="encabezados"
        :items="items"
        :items-length="total"
        :loading="cargando"
        must-sort
      >
        <template v-for="(_, nombre) in slotsDeCelda" #[nombre]="scope">
          <slot :name="nombre" v-bind="scope" />
        </template>

        <template #[`item.acciones`]="{ item }">
          <btn-icono :aria-label="`Editar ${String((item as Fila)[campoTitulo])}`" icon="$editar" @click="emit('editar', item)" />
          <btn-icono :aria-label="`Eliminar ${String((item as Fila)[campoTitulo])}`" icon="$eliminar" @click="emit('eliminar', item)" />
        </template>

        <!-- Al paginar u ordenar: skeleton de filas -->
        <template #loading>
          <v-skeleton-loader type="table-row@5" />
        </template>
      </v-data-table-server>
    </v-fade-transition>
  </v-card>
</template>

<script lang="ts" setup generic="T extends { id: string }">
  import type { Orden } from '@/composables/useListado'
  import type { DataTableHeader } from 'vuetify'
  import { computed, useSlots } from 'vue'
  import { useDisplay } from 'vuetify'

  type Fila = Record<string, unknown> & { id: string }

  const props = defineProps<{
    columnas: DataTableHeader<T>[]
    items: T[]
    total: number
    cargando: boolean
    listo: boolean
    error: string
    hayBusqueda: boolean
    campoTitulo: string // campo que titula la tarjeta en móvil
    subtitulo?: (item: T) => string
    placeholder: string
    tituloVacio: string
    textoVacio: string
    nuevo?: { texto: string, to: string }
  }>()

  const emit = defineEmits<{
    editar: [item: T]
    eliminar: [item: T]
    reintentar: []
  }>()

  const busqueda = defineModel<string>('busqueda', { required: true })
  const pagina = defineModel<number>('pagina', { required: true })
  const porPagina = defineModel<number>('porPagina', { required: true })
  const orden = defineModel<Orden[]>('orden', { required: true })

  const { mobile } = useDisplay()
  const slots = useSlots()

  // La columna de acciones se agrega sola al final.
  const encabezados = computed<DataTableHeader<T>[]>(() => [
    ...props.columnas,
    { key: 'acciones', title: 'Acciones', sortable: false, align: 'end' },
  ])

  // Celdas personalizadas que manda la página (#item.campo).
  const slotsDeCelda = computed(() =>
    Object.fromEntries(Object.entries(slots).filter(([n]) => n.startsWith('item.'))),
  )

  const paginas = computed(() => Math.max(1, Math.ceil(props.total / porPagina.value)))
</script>
