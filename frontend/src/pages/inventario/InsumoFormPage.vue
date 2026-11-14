<!-- Alta y edición de insumo (mock: AltaInsumo). Tipo de material, proveedor y notas: Fase 7. -->
<template>
  <error-pagina v-if="noEncontrado" codigo="404" />

  <formulario-crud
    v-else
    :cargando="cargando"
    :error="error"
    :etiqueta-guardar="esEdicion ? 'Guardar cambios' : 'Guardar insumo'"
    :guardando="guardando"
    :migas="[{ title: 'Inventario', to: '/insumos' }, { title: titulo }]"
    subtitulo="Consumible desechable de un solo uso. Su stock se controla por piezas."
    :titulo="titulo"
    @cancelar="volver"
    @guardar="guardar"
  >
    <template v-if="esEdicion" #acciones>
      <v-btn :block="mobile" color="error" prepend-icon="$eliminar" @click="pedirBaja({ id: id! })">Eliminar</v-btn>
    </template>

    <seccion-formulario ayuda="La clave no puede repetirse con la de un equipo ni con la de otro insumo." primera titulo="Datos generales">
      <v-col cols="12" md="3">
        <v-text-field
          v-model="datos.clave"
          autofocus
          :error-messages="errores.clave"
          label="Clave *"
          maxlength="30"
          placeholder="SON-3-24"
          :rules="[requerido, maximo(30), codigo]"
          @update:model-value="limpiarError('clave')"
        />
      </v-col>

      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.descripcion"
          :error-messages="errores.descripcion"
          label="Descripción / material *"
          maxlength="200"
          placeholder="Sonda Foley 3 vías 24 fr"
          :rules="[requerido, maximo(200), textoGeneral]"
          @update:model-value="limpiarError('descripcion')"
        />
      </v-col>

      <v-col cols="12" md="3">
        <v-select
          v-model="datos.unidad"
          :error-messages="errores.unidad"
          :items="UNIDADES"
          label="Unidad *"
          :rules="[requerido, enLista(UNIDADES)]"
          @update:model-value="limpiarError('unidad')"
        />
      </v-col>

      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.codigoBarras"
          :error-messages="errores.codigoBarras"
          inputmode="numeric"
          label="Código de barras"
          maxlength="64"
          :rules="[maximo(64), codigo]"
          @update:model-value="limpiarError('codigoBarras')"
        />
      </v-col>
    </seccion-formulario>

    <seccion-formulario ayuda="El precio se usa en la cotización; el stock mínimo marca el insumo como bajo mínimo." titulo="Precio y existencias">
      <v-col cols="12" md="4">
        <v-text-field
          v-model="datos.precio"
          :error-messages="errores.precio"
          inputmode="decimal"
          label="Precio unitario (MXN) *"
          maxlength="13"
          prefix="$"
          :rules="[requerido, dinero]"
          @update:model-value="limpiarError('precio')"
        />
      </v-col>

      <v-col cols="6" md="4">
        <v-text-field
          v-model="datos.stockActual"
          :error-messages="errores.stockActual"
          inputmode="numeric"
          :label="esEdicion ? 'Stock actual *' : 'Stock inicial *'"
          maxlength="7"
          :rules="[requerido, enteroPositivo]"
          @update:model-value="limpiarError('stockActual')"
        />
      </v-col>

      <v-col cols="6" md="4">
        <v-text-field
          v-model="datos.stockMinimo"
          :error-messages="errores.stockMinimo"
          inputmode="numeric"
          label="Stock mínimo *"
          maxlength="7"
          :rules="[requerido, enteroPositivo]"
          @update:model-value="limpiarError('stockMinimo')"
        />
      </v-col>
    </seccion-formulario>
  </formulario-crud>

  <dialogo-eliminar
    v-model="bajaAbierta"
    :eliminando="eliminando"
    :texto="textoBaja"
    titulo="¿Eliminar el insumo?"
    @confirmar="confirmarBaja"
  />
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import DialogoEliminar from '@/components/DialogoEliminar.vue'
  import ErrorPagina from '@/components/ErrorPagina.vue'
  import FormularioCrud from '@/components/FormularioCrud.vue'
  import SeccionFormulario from '@/components/SeccionFormulario.vue'
  import { useBaja } from '@/composables/useBaja'
  import { useFormulario } from '@/composables/useFormulario'
  import { aEnvio, insumos, UNIDADES } from '@/services/insumos'
  import { codigo, dinero, enLista, enteroPositivo, maximo, requerido, textoGeneral } from '@/utils/rules'

  const route = useRoute()
  const router = useRouter()
  const { mobile } = useDisplay()
  const id = computed(() => (route.params.id as string | undefined) || undefined)

  const { datos, errores, error, cargando, guardando, noEncontrado, esEdicion, guardar, limpiarError }
    = useFormulario(insumos, {
      id,
      vacio: () => ({ clave: '', descripcion: '', unidad: 'PIEZA', precio: '', codigoBarras: '', stockActual: '0', stockMinimo: '0' }),
      desde: i => ({
        clave: i.clave, descripcion: i.descripcion, unidad: i.unidad, precio: String(i.precio),
        codigoBarras: i.codigoBarras, stockActual: String(i.stockActual), stockMinimo: String(i.stockMinimo),
      }),
      enviar: aEnvio,
      mensajes: { creado: 'Insumo dado de alta', actualizado: 'Insumo actualizado' },
      alGuardar: volver,
    })

  const titulo = computed(() => (esEdicion.value ? 'Editar insumo' : 'Alta de insumo'))
  const textoBaja = computed(() => `${datos.clave} dejará de aparecer al programar cirugías. Su clave no se podrá volver a usar.`)

  // En móvil la baja se hace desde aquí (la lista solo abre la edición).
  const { abierto: bajaAbierta, eliminando, pedir: pedirBaja, confirmar: confirmarBaja }
    = useBaja<{ id: string }>(insumos.eliminar, 'Insumo eliminado', volver)

  function volver () {
    router.push('/insumos')
  }
</script>
