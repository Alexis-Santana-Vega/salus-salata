<!-- Alta y edición de procedimiento (mock: AltaProcedimiento). -->
<template>
  <error-pagina v-if="noEncontrado" codigo="404" />

  <formulario-crud
    v-else
    :cargando="cargando"
    :error="error"
    :etiqueta-guardar="esEdicion ? 'Guardar cambios' : 'Guardar procedimiento'"
    :guardando="guardando"
    :migas="[{ title: 'Catálogos', to: '/catalogos/procedimientos' }, { title: titulo }]"
    subtitulo="Tipo de cirugía que se elige al programar."
    :titulo="titulo"
    @cancelar="volver"
    @guardar="guardar"
  >
    <template v-if="esEdicion" #acciones>
      <v-btn :block="mobile" color="error" prepend-icon="$eliminar" @click="baja.pedir({ id: id! })">Eliminar</v-btn>
    </template>

    <seccion-formulario ayuda="Ej. Holep, Artroscopia / Morcelación, Cirugía General, Traumatología." primera titulo="Datos del procedimiento">
      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.nombre"
          autofocus
          counter="120"
          :error-messages="errores.nombre"
          label="Nombre *"
          :rules="[requerido, maximo(120), textoSeguro]"
          @update:model-value="limpiarError('nombre')"
        />
      </v-col>
    </seccion-formulario>
  </formulario-crud>

  <dialogo-eliminar
    v-model="baja.abierto.value"
    :eliminando="baja.eliminando.value"
    :texto="`«${datos.nombre}» dejará de aparecer al programar cirugías. Las cirugías que ya lo usan lo conservan.`"
    titulo="¿Eliminar el procedimiento?"
    @confirmar="baja.confirmar"
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
  import { procedimientos } from '@/services/procedimientos'
  import { maximo, requerido, textoSeguro } from '@/utils/rules'

  const route = useRoute()
  const router = useRouter()
  const { mobile } = useDisplay()
  const id = computed(() => (route.params.id as string | undefined) || undefined)

  const { datos, errores, error, cargando, guardando, noEncontrado, esEdicion, guardar, limpiarError }
    = useFormulario(procedimientos, {
      id,
      vacio: () => ({ nombre: '' }),
      desde: p => ({ nombre: p.nombre }),
      mensajes: { creado: 'Procedimiento dado de alta', actualizado: 'Procedimiento actualizado' },
      alGuardar: volver,
    })

  const titulo = computed(() => (esEdicion.value ? 'Editar procedimiento' : 'Alta de procedimiento'))

  // En móvil la baja se hace desde aquí (la lista solo abre la edición).
  const baja = useBaja<{ id: string }>(procedimientos.eliminar, 'Procedimiento eliminado', volver)

  function volver () {
    router.push('/catalogos/procedimientos')
  }
</script>
